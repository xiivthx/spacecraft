package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"
)

type TasksSummary struct {
	Total     int `json:"total"`
	Completed int `json:"completed"`
}

type WorkflowSnapshot struct {
	MissionID        string       `json:"missionId"`
	Title            string       `json:"title"`
	State            string       `json:"state"`
	Safety           string       `json:"safety"`
	Source           string       `json:"source"`
	Next             string       `json:"next"`
	NextTask         *Task        `json:"nextTask"`
	Tasks            TasksSummary `json:"tasks"`
	EvidenceCount    int          `json:"evidenceCount"`
	Blockers         []string     `json:"blockers"`
	CheckpointPolicy string       `json:"checkpointPolicy"`
}

type Task struct {
	ID     *string `json:"id"`
	Title  *string `json:"title"`
	Status *string `json:"status"`
}

func taskIsOpen(task Task) bool {
	if task.Status == nil {
		return true
	}
	switch *task.Status {
	case "completed", "done", "cancelled":
		return false
	}
	return true
}

func nextOpenTask(tasks []Task) *Task {
	for _, task := range tasks {
		if taskIsOpen(task) {
			t := task
			return &t
		}
	}
	return nil
}

func nextCommandForMission(m *Mission) string {
	if m == nil {
		return "/sc-status"
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
		return "/sc-work"
	case "verifying":
		return "/sc-verify"
	case "reviewing":
		return "/sc-review"
	case "ready":
		return "/sc-ship"
	case "shipped":
		return "(shipped)"
	case "blocked":
		return "/sc-status"
	default:
		return "/sc-status"
	}
}

func workflowSnapshot(res ResolveOutput, mission *Mission, specExists bool, planExists bool, plan *Plan, evidenceCount int, git GitInfoData) WorkflowSnapshot {
	var tasks []Task
	if plan != nil && plan.Tasks != nil {
		tasks = plan.Tasks
	}
	nextTask := nextOpenTask(tasks)
	blockers := []string{}
	
	hasBlockingClarification := false
	if mission != nil && (mission.Clarification.Status == "open" || mission.Clarification.BlockingQuestions > 0) {
		hasBlockingClarification = true
	}
	
	hasTaskPlan := planExists && len(tasks) > 0
	artifactGateClear := specExists && hasTaskPlan
	
	next := nextCommandForMission(mission)
	
	if hasBlockingClarification {
		blockers = append(blockers, "blocking clarification remains open")
		next = "/sc-clarify"
	}
	if !specExists {
		blockers = append(blockers, "spec.md is missing")
		if !hasBlockingClarification {
			next = "/sc-status"
		}
	}
	if !planExists {
		blockers = append(blockers, "plan.json is missing")
		if !hasBlockingClarification && specExists {
			next = "/sc-plan"
		}
	} else if !hasTaskPlan {
		blockers = append(blockers, "plan.json has no tasks")
		if !hasBlockingClarification && specExists {
			next = "/sc-plan"
		}
	}
	
	state := ""
	if mission != nil {
		state = mission.State
	}
	
	onMainBlock := state == "planned" || state == "verifying"
	// Backward compat for missions with legacy 'implementing' state
	if state == "implementing" {
		onMainBlock = true
	}
	
	if onMainBlock && git.IsRepo && (git.Branch == "main") {
		blockers = append(blockers, "implementation workflow requires a non-main work branch")
	}
	if onMainBlock && git.IsRepo && git.Dirty {
		blockers = append(blockers, fmt.Sprintf("worktree is dirty (%d files); inspect before automated workflow", git.DirtyFiles))
	}
	
	if !hasBlockingClarification && artifactGateClear && (state == "planned" || state == "implementing") {
		if nextTask != nil {
			idStr := ""
			if nextTask.ID != nil {
				idStr = *nextTask.ID
			}
			if idStr == "" {
				next = "/sc-work"
			} else {
				next = "/sc-work " + idStr
			}
		} else {
			next = "/sc-review"
		}
	}
	if !hasBlockingClarification && artifactGateClear && state == "verifying" {
		if nextTask != nil {
			idStr := ""
			if nextTask.ID != nil {
				idStr = *nextTask.ID
			}
			if idStr == "" {
				next = "/sc-verify"
			} else {
				next = "/sc-verify " + idStr
			}
		} else {
			next = "/sc-review"
		}
	}
	if !hasBlockingClarification && artifactGateClear && nextTask == nil {
		if state == "ready" {
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

	var nt *Task
	if nextTask != nil {
		nt = &Task{
			ID: nextTask.ID,
			Title: nextTask.Title,
			Status: nextTask.Status,
		}
	}

	title := ""
	if mission != nil {
		title = mission.Title
	}
	mID := ""
	if mission != nil {
		mID = mission.ID
	}

	return WorkflowSnapshot{
		MissionID:        mID,
		Title:            title,
		State:            state,
		Safety:           res.Safety,
		Source:           getString(res.Source),
		Next:             next,
		NextTask:         nt,
		Tasks: TasksSummary{
			Total:     len(tasks),
			Completed: completedCount,
		},
		EvidenceCount:    evidenceCount,
		Blockers:         blockers,
		CheckpointPolicy: "After passing verification, checkpoint commit on a clean non-main work branch before the next task.",
	}
}

func printWorkflow(args []string) {
	jsonOut := len(args) > 0 && args[0] == "--json"
	res := resolveMission("")
	if res.Safety != "safe" || res.Selected == nil {
		fail(formatResolutionBlock(res, "flow"))
	}

	id := res.Selected.ID
	dir := missionDir(id)
	
	var mission Mission
	readJson(filepath.Join(dir, "mission.json"), &mission)
	
	specExists := exists(filepath.Join(dir, "spec.md"))
	planPath := filepath.Join(dir, "plan.json")
	planExists := exists(planPath)
	
	var plan Plan
	if planExists {
		readJson(planPath, &plan)
	}
	
	evCount := countEvidence(filepath.Join(dir, "evidence.jsonl"))
	git := gitInfo()
	
	snapshot := workflowSnapshot(res, &mission, specExists, planExists, &plan, evCount, git)

	if jsonOut {
		var buf bytes.Buffer
		enc := json.NewEncoder(&buf)
		enc.SetEscapeHTML(false)
		enc.SetIndent("", "  ")
		enc.Encode(snapshot)
		fmt.Print(buf.String())
		return
	}

	readyStr := "ready"
	if len(snapshot.Blockers) > 0 {
		readyStr = "blocked"
	}
	fmt.Printf("Workflow: %s\n", readyStr)
	fmt.Printf("Mission: %s (%s)\n", snapshot.Title, snapshot.MissionID)
	fmt.Printf("State: %s\n", snapshot.State)
	fmt.Printf("Tasks: %d/%d completed\n", snapshot.Tasks.Completed, snapshot.Tasks.Total)
	fmt.Printf("Evidence: %d\n", snapshot.EvidenceCount)
	if snapshot.NextTask != nil {
		display := "(unnamed task)"
		if snapshot.NextTask.ID != nil && *snapshot.NextTask.ID != "" && snapshot.NextTask.Title != nil {
			display = fmt.Sprintf("%s %s", *snapshot.NextTask.ID, *snapshot.NextTask.Title)
		} else if snapshot.NextTask.ID != nil && *snapshot.NextTask.ID != "" {
			display = *snapshot.NextTask.ID
		} else if snapshot.NextTask.Title != nil {
			display = *snapshot.NextTask.Title
		}
		fmt.Printf("Next task: %s\n", display)
	}
	fmt.Printf("Next: %s\n", snapshot.Next)
	fmt.Println("Loop: /sc-work Txx -> /sc-verify Txx -> checkpoint commit -> next task, until a gate blocks.")
	fmt.Printf("Checkpoint: %s\n", snapshot.CheckpointPolicy)
	if len(snapshot.Blockers) > 0 {
		fmt.Println("Blockers:")
		for _, b := range snapshot.Blockers {
			fmt.Printf("- %s\n", b)
		}
	}
}

func printStatus() {
	res := resolveMission("")
	if res.Safety != "safe" || res.Selected == nil {
		fmt.Println("No selected Spacecraft mission. Start one with /sc-start <title>.")
		if len(res.Candidates) > 0 {
			fmt.Println("Candidates:")
			for i, c := range res.Candidates {
				num := i + 1
				if c.Number != nil {
					num = *c.Number
				}
				branchHint := ""
				if len(c.Branches) > 0 {
					branchHint = " branch:" + strings.Join(c.Branches, ",")
				}
				signal := ""
				if c.Signal != nil {
					signal = " signal:" + *c.Signal
				}
				fmt.Printf("%d. %s (%s) - state:%s%s%s\n", num, c.Title, c.ID, c.State, signal, branchHint)
			}
		}
		return
	}

	id := res.Selected.ID
	dir := missionDir(id)
	missionPath := filepath.Join(dir, "mission.json")
	if !exists(missionPath) {
		fail(fmt.Sprintf("Selected mission %s is missing mission.json.", id))
	}

	var mission Mission
	readJson(missionPath, &mission)

	var plan Plan
	planPath := filepath.Join(dir, "plan.json")
	planExists := exists(planPath)
	if planExists {
		readJson(planPath, &plan)
	}

	var review map[string]interface{}
	reviewPath := filepath.Join(dir, "review.json")
	reviewExists := exists(reviewPath)
	if reviewExists {
		readJson(reviewPath, &review)
	}

	taskCount := 0
	if plan.Tasks != nil {
		taskCount = len(plan.Tasks)
	}

	evCount := countEvidence(filepath.Join(dir, "evidence.jsonl"))

	src := "unknown"
	if res.Source != nil {
		src = *res.Source
	}

	fmt.Printf("Mission: %s\n", mission.ID)
	fmt.Printf("Selected by: %s\n", src)
	if res.Safety != "safe" {
		fmt.Printf("Mission safety: %s\n", res.Safety)
	}
	if res.CurrentMissionId != nil {
		fmt.Printf("Current: %s\n", *res.CurrentMissionId)
	}
	if len(res.Conflicts) > 0 {
		fmt.Println("Conflicts:")
		for _, c := range res.Conflicts {
			fmt.Printf("- %s\n", c.Type)
		}
	}
	if res.Safety != "safe" && len(res.Candidates) > 0 {
		fmt.Println("Candidates:")
		for _, c := range res.Candidates {
			fmt.Printf("- %s (%s)\n", c.Title, c.ID)
		}
	}
	fmt.Printf("Title: %s\n", mission.Title)
	fmt.Printf("State: %s\n", mission.State)
	if mission.Clarification.Status != "" {
		fmt.Printf("Clarification: %s\n", mission.Clarification.Status)
		fmt.Printf("Blocking questions: %d\n", mission.Clarification.BlockingQuestions)
	}

	git := gitInfo()
	if git.IsRepo {
		branch := "(detached)"
		if git.Branch != "" {
			branch = git.Branch
		}
		sha := "(no commit)"
		if git.Sha != "" && len(git.Sha) >= 12 {
			sha = git.Sha[:12]
		}
		status := " clean"
		if git.Dirty {
			status = fmt.Sprintf(" dirty:%d", git.DirtyFiles)
		}
		fmt.Printf("Git: %s %s%s\n", branch, sha, status)
		if mission.BaseSha != nil && git.Sha != "" && *mission.BaseSha != git.Sha {
			fmt.Printf("Mission base: %s\n", (*mission.BaseSha)[:12])
		}
	} else {
		fmt.Println("Git: not a git worktree")
	}

	fmt.Println("Artifacts:")
	fmt.Printf("  spec: %s\n", displayPath(filepath.Join(dir, "spec.md")))
	fmt.Printf("  plan: %s\n", displayPath(planPath))
	fmt.Printf("  evidence: %s\n", displayPath(filepath.Join(dir, "evidence.jsonl")))
	fmt.Printf("  review: %s\n", displayPath(filepath.Join(dir, "review.md")))
	fmt.Printf("  reviewJson: %s\n", displayPath(reviewPath))
	if exists(filepath.Join(dir, "questions.md")) {
		fmt.Printf("  questions: %s\n", displayPath(filepath.Join(dir, "questions.md")))
	}
	if exists(filepath.Join(dir, "decisions.md")) {
		fmt.Printf("  decisions: %s\n", displayPath(filepath.Join(dir, "decisions.md")))
	}
	if exists(filepath.Join(dir, "design")) {
		fmt.Printf("  design: %s\n", displayPath(filepath.Join(dir, "design")))
	}
	fmt.Printf("Tasks: %d\n", taskCount)
	fmt.Printf("Evidence: %d\n", evCount)
	reviewStatus := "not-reviewed"
	if reviewExists {
		if s, ok := review["status"].(string); ok && s != "" {
			reviewStatus = s
		}
	}
	fmt.Printf("Review: %s\n", reviewStatus)
}

func getString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
