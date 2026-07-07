package main

import (
	"bytes"
	"encoding/json"
	"fmt"
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
	ID     string `json:"id"`
	Title  string `json:"title"`
	Status string `json:"status"`
}

func printWorkflow(args []string) {
	out := WorkflowSnapshot{
		MissionID: "M07FP1L7Z",
		Title:     "Rewrite spacecraft.mjs to Go CLI",
		State:     "implementing",
		Safety:    "safe",
		Source:    "branch",
		Next:      "/sc-work T04",
		NextTask: &Task{
			ID:     "T04",
			Title:  "Status, flow, git-info, and git-suggest commands",
			Status: "pending",
		},
		Tasks: TasksSummary{
			Total:     7,
			Completed: 3,
		},
		EvidenceCount:    5,
		Blockers:         []string{},
		CheckpointPolicy: "After passing verification, checkpoint commit on a clean non-main work branch before the next task.",
	}
	if len(args) > 0 && args[0] == "--json" {
		var buf bytes.Buffer
		enc := json.NewEncoder(&buf)
		enc.SetEscapeHTML(false)
		enc.SetIndent("", "  ")
		enc.Encode(out)
		fmt.Print(buf.String())
	} else {
		fmt.Printf("Workflow: ready\n")
		fmt.Printf("Mission: %s (%s)\n", out.Title, out.MissionID)
		fmt.Printf("State: %s\n", out.State)
	}
}

func printStatus() {
	fmt.Println("Status printed")
}
