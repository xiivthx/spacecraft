// Package workflow provides mission workflow snapshot and next-command logic.
package workflow

import (
	"fmt"

	"spacecraft/internal/mission"
)

// NextTask returns the first open task in a plan, or nil.
func NextTask(tasks []mission.Task) *mission.Task {
	for _, task := range tasks {
		if taskIsOpen(task) {
			t := task
			return &t
		}
	}
	return nil
}

// taskIsOpen returns true if the task is not completed/done/cancelled.
func taskIsOpen(task mission.Task) bool {
	if task.Status == nil {
		return true
	}
	switch *task.Status {
	case "completed", "done", "cancelled":
		return false
	}
	return true
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
		return "/sc-clarify"
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
	store mission.MissionStore
}

// NewSnapshot creates a workflow snapshot from mission data.
func NewSnapshot(store mission.MissionStore) *Snapshot {
	return &Snapshot{store: store}
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
		next = "/sc-clarify"
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
		if m.State == "ready" {
			next = "/sc-ship"
		} else {
			next = "/sc-review"
		}
	}

	completedCount := 0
	for _, t := range tasks {
		if t.Status != nil && (*t.Status == "completed" || *t.Status == "done") {
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
	if onMainBlock && gitInfo.IsRepo && res.Git.Branch != "" {
		// Check dirty state isn't from res (no direct access to dirty in GitInfo)
	}

	return blockers
}

func getString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
