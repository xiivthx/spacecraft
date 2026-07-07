package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"path/filepath"
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

func nextOpenTask(tasks []Task) *Task {
	for _, task := range tasks {
		if task.Status == nil || *task.Status != "completed" {
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
	case "draft", "specified":
		return "/sc-plan"
	case "planned", "implementing":
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
	
	if (state == "planned" || state == "implementing" || state == "verifying") && git.IsRepo && (git.Branch == "main") {
		blockers = append(blockers, "implementation workflow requires a non-main work branch")
	}
	if (state == "planned" || state == "implementing" || state == "verifying") && git.IsRepo && git.Dirty {
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
		if t.Status != nil && *t.Status == "completed" {
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
	// Stub until T04 is fully implemented, but flow is now implemented
	fmt.Println("Status printed")
}

func getString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
