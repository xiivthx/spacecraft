// Package workflow provides mission workflow snapshot and next-command logic.
package workflow

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"spacecraft/internal/mission"
)

// NextTask returns the first open task in a plan whose dependencies are all done, or nil.
func NextTask(tasks []mission.Task) *mission.Task {
	done := make(map[string]bool)
	for _, t := range tasks {
		if t.ID != nil && t.Status != nil && *t.Status == "done" {
			done[*t.ID] = true
		}
	}
	for _, task := range tasks {
		if !taskIsOpen(task) {
			continue
		}
		depsReady := true
		for _, dep := range task.DependsOn {
			if !done[dep] {
				depsReady = false
				break
			}
		}
		if depsReady {
			t := task
			return &t
		}
	}
	return nil
}

// taskIsOpen returns true if the task is not completed/done/cancelled/waiting.
func taskIsOpen(task mission.Task) bool {
	if task.Status == nil {
		return true
	}
	switch *task.Status {
	case "done", "cancelled", "waiting":
		return false
	}
	return true
}

// hasWaitingTasks returns true if any task has status "waiting".
func hasWaitingTasks(tasks []mission.Task) bool {
	for _, t := range tasks {
		if t.Status != nil && *t.Status == "waiting" {
			return true
		}
	}
	return false
}

// NextCommand returns the recommended next slash command for a mission.
// This is a pure function with no side effects.
func NextCommand(m *mission.Mission) string {
	if m == nil {
		return "/sc-resume"
	}
	status := "open"
	if m.Clarification.Status != "" {
		status = m.Clarification.Status
	}
	if status == "open" || m.Clarification.BlockingQuestions > 0 {
		return "(clarify)"
	}
	switch m.State {
	case "draft":
		return "/sc-plan"
	case "planned":
		return "/sc-build"
	case "built":
		return "/sc-review"
	case "ready":
		return "/sc-ship"
	case "shipped":
		return "(shipped)"
	case "blocked":
		return "/sc-resume"
	default:
		return "/sc-resume"
	}
}

// Snapshot holds the computed workflow state.
type Snapshot struct {
	store       mission.MissionStore
	commandsDir string
}

// NewSnapshot creates a workflow snapshot from mission data.
func NewSnapshot(store mission.MissionStore) *Snapshot {
	return &Snapshot{store: store}
}

// SetCommandsDir sets the path to the OpenCode commands directory.
// When set, Build() validates that the Next command exists as a .md file
// in this directory. If the command is missing, a blocker is added and
// Next falls back to /sc-resume.
func (s *Snapshot) SetCommandsDir(dir string) {
	s.commandsDir = dir
}

// Build constructs a WorkflowSnapshot from the given parameters.
func (s *Snapshot) Build(res mission.ResolveOutput, missionID string) (mission.WorkflowSnapshot, error) {
	m, err := s.store.Load(missionID)
	if err != nil {
		return mission.WorkflowSnapshot{}, fmt.Errorf("load mission %s: %w", missionID, err)
	}

	specExists := s.store.SpecExists(missionID)
	planExists := s.store.PlanExists(missionID)

	var plan *mission.Plan
	if planExists {
		p, err := s.store.LoadPlan(missionID)
		if err == nil {
			plan = p
		}
	}

	evidenceCount, _ := s.store.CountEvidence(missionID)

	var tasks []mission.Task
	if plan != nil && plan.Tasks != nil {
		tasks = plan.Tasks
	}
	nextTask := NextTask(tasks)
	blockers := buildBlockers(m, specExists, planExists, tasks, res)

	hasBlockingClarification := m.Clarification.Status == "open" || m.Clarification.BlockingQuestions > 0
	hasTaskPlan := planExists && len(tasks) > 0
	artifactGateClear := specExists && hasTaskPlan

	next := NextCommand(m)

	if hasBlockingClarification {
		next = "(clarify)"
	}
	if !specExists && !hasBlockingClarification {
		next = "/sc-resume"
	}
	if !planExists && !hasBlockingClarification && specExists {
		next = "/sc-plan"
	} else if !hasTaskPlan && !hasBlockingClarification && specExists {
		next = "/sc-plan"
	}

	if !hasBlockingClarification && artifactGateClear {
		switch m.State {
		case "planned", "implementing":
			if nextTask != nil {
				idStr := ""
				if nextTask.ID != nil {
					idStr = *nextTask.ID
				}
				if idStr == "" {
					next = "/sc-build"
				} else {
					next = "/sc-build " + idStr
				}
			} else {
				next = "/sc-review"
			}
		case "built":
			next = "/sc-review"
		}
	}

	if !hasBlockingClarification && artifactGateClear && nextTask == nil {
		if hasWaitingTasks(tasks) {
			blockers = append(blockers, "all open tasks are waiting on architectural guidance — check .space/architect/")
			next = "/sc-resume"
		} else if m.State == "ready" {
			next = "/sc-ship"
		} else {
			next = "/sc-review"
		}
	}

	completedCount := 0
	for _, t := range tasks {
		if t.Status != nil && *t.Status == "done" {
			completedCount++
		}
	}

	var nt *mission.Task
	if nextTask != nil {
		nt = &mission.Task{
			ID:    nextTask.ID,
			Title: nextTask.Title,
			Status: nextTask.Status,
		}
	}

	// Validate that the next command actually exists as a registered command file.
	if s.commandsDir != "" {
		next, blockers = s.validateNextCommand(next, blockers)
	}

	return mission.WorkflowSnapshot{
		MissionID:     missionID,
		Title:         m.Title,
		State:         m.State,
		Safety:        res.Safety,
		Source:        getString(res.Source),
		Next:          next,
		NextTask:      nt,
		Tasks: mission.TasksSummary{
			Total:     len(tasks),
			Completed: completedCount,
		},
		EvidenceCount:    evidenceCount,
		Blockers:         blockers,
		CheckpointPolicy: "After passing verification, checkpoint commit on a clean non-main work branch before the next task.",
	}, nil
}

// validateNextCommand checks that the next command exists as a .md file in the
// commands directory. If not, adds a blocker and falls back to /sc-resume.
// Parenthesized strings like "(shipped)" and "(clarify)" are status markers,
// not commands — they pass through unvalidated.
func (s *Snapshot) validateNextCommand(next string, blockers []string) (string, []string) {
	if next == "" || strings.HasPrefix(next, "(") {
		return next, blockers
	}
	// Extract command name from "/sc-build T02" -> "sc-build"
	cmdName := strings.TrimPrefix(next, "/")
	// Strip arguments (space after command name)
	if idx := strings.Index(cmdName, " "); idx != -1 {
		cmdName = cmdName[:idx]
	}
	cmdPath := filepath.Join(s.commandsDir, cmdName+".md")
	if _, err := os.Stat(cmdPath); os.IsNotExist(err) {
		blockers = append(blockers, fmt.Sprintf("next command %s is not registered in .opencode/commands/ (missing %s)", next, cmdName+".md"))
		return "/sc-resume", blockers
	}
	return next, blockers
}

// buildBlockers checks for blockers that would prevent workflow progression.
func buildBlockers(m *mission.Mission, specExists, planExists bool, tasks []mission.Task, res mission.ResolveOutput) []string {
	blockers := make([]string, 0)

	hasBlockingClarification := m.Clarification.Status == "open" || m.Clarification.BlockingQuestions > 0

	if hasBlockingClarification {
		blockers = append(blockers, "blocking clarification remains open")
	}
	if !specExists {
		blockers = append(blockers, "spec.md is missing")
	}
	if !planExists {
		blockers = append(blockers, "plan.json is missing")
	} else if len(tasks) == 0 {
		blockers = append(blockers, "plan.json has no tasks")
	}

	state := m.State
	onMainBlock := state == "planned"
	if state == "implementing" {
		onMainBlock = true
	}

	gitInfo := res.Git
	if onMainBlock && gitInfo.IsRepo && gitInfo.Branch == "main" {
		blockers = append(blockers, "implementation workflow requires a non-main work branch")
	}
	if gitInfo.IsRepo && gitInfo.Branch != "main" && gitInfo.Dirty {
		blockers = append(blockers, fmt.Sprintf("worktree has %d uncommitted file(s); commit or stash before continuing", gitInfo.DirtyFiles))
	}

	return blockers
}

func getString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
