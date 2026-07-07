package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func validateMission() {
	res := resolveMission("")
	if res.Safety != "safe" || res.Selected == nil {
		fail(formatResolutionBlock(res, "validate"))
	}

	id := res.Selected.ID
	dir := missionDir(id)
	var errors []string

	requireFile := func(relPath string) *string {
		fp := filepath.Join(dir, relPath)
		if !exists(fp) {
			errors = append(errors, "Missing "+relPath)
			return nil
		}
		return &fp
	}

	mPath := requireFile("mission.json")
	if mPath != nil {
		var m Mission
		err := readJson(*mPath, &m)
		if err != nil {
			errors = append(errors, "Invalid JSON in mission.json: "+err.Error())
		} else {
			status := m.Clarification.Status
			if status != "" && status != "open" && status != "clear" && status != "deferred" {
				errors = append(errors, "mission.json clarification.status must be one of: open, clear, deferred")
			}
		}
	}

	requireFile("spec.md")

	pPath := requireFile("plan.json")
	if pPath != nil {
		var p Plan
		err := readJson(*pPath, &p)
		if err != nil {
			errors = append(errors, "Invalid JSON in plan.json: "+err.Error())
		} else if p.Tasks == nil {
			errors = append(errors, "plan.json must contain a tasks array")
		}
	}

	evPath := requireFile("evidence.jsonl")
	if evPath != nil {
		content, err := os.ReadFile(*evPath)
		if err == nil {
			lines := strings.Split(string(content), "\n")
			evIds := make(map[string]int)
			evPaths := make(map[string]int)
			for i, line := range lines {
				if strings.TrimSpace(line) == "" {
					continue
				}
				var entry map[string]interface{}
				if err := json.Unmarshal([]byte(line), &entry); err != nil {
					errors = append(errors, fmt.Sprintf("Invalid JSON in evidence.jsonl line %d: %s", i+1, err.Error()))
					continue
				}
				idVal, ok := entry["id"].(string)
				if !ok || idVal == "" {
					errors = append(errors, fmt.Sprintf("evidence.jsonl line %d must have string id", i+1))
				} else if prev, exists := evIds[idVal]; exists {
					errors = append(errors, fmt.Sprintf("Duplicate evidence id %s on lines %d and %d", idVal, prev, i+1))
				} else {
					evIds[idVal] = i + 1
				}

				for _, field := range []string{"stdout", "stderr"} {
					pathVal, ok := entry[field].(string)
					if !ok || pathVal == "" {
						errors = append(errors, fmt.Sprintf("evidence.jsonl line %d must have %s path", i+1, field))
						continue
					}
					absPath := pathVal
					if !filepath.IsAbs(pathVal) {
						absPath = filepath.Join(ROOT, pathVal)
					}
					if prev, exists := evPaths[absPath]; exists {
						errors = append(errors, fmt.Sprintf("Duplicate evidence %s path %s on lines %d and %d", field, pathVal, prev, i+1))
					} else {
						evPaths[absPath] = i + 1
					}
					if !exists(absPath) {
						errors = append(errors, fmt.Sprintf("Missing evidence %s file for line %d: %s", field, i+1, pathVal))
					}
				}
			}
		}
	}

	revPath := filepath.Join(dir, "review.json")
	if exists(revPath) {
		var r interface{}
		err := readJson(revPath, &r)
		if err != nil {
			errors = append(errors, "Invalid JSON in review.json: "+err.Error())
		}
	}

	if len(errors) > 0 {
		fmt.Fprintf(os.Stderr, "Spacecraft mission %s is invalid:\n", id)
		for _, errStr := range errors {
			fmt.Fprintf(os.Stderr, "- %s\n", errStr)
		}
		os.Exit(1)
	}

	fmt.Printf("Spacecraft mission %s is valid.\n", id)
}

func setState(state string) {
	allowed := map[string]bool{
		"draft": true, "planned": true,
		"verifying": true, "reviewing": true,
		"ready": true, "shipped": true, "blocked": true,
	}
	// Backward compat: map removed states to their replacement
	if !allowed[state] {
		switch state {
		case "specified":
			state = "draft"
		case "implementing":
			state = "planned"
		default:
			fail("Invalid state. Allowed states: draft, planned, verifying, reviewing, ready, shipped, blocked")
		}
	}

	res := resolveMission("")
	if res.Safety != "safe" || res.Selected == nil {
		fail(formatResolutionBlock(res, "set-state"))
	}
	
	id := res.Selected.ID
	mPath := filepath.Join(missionDir(id), "mission.json")
	var m Mission
	readJson(mPath, &m)
	m.State = state
	m.UpdatedAt = isoNow()
	writeJson(mPath, m)
	fmt.Printf("Spacecraft mission %s state: %s\n", id, state)
}

func setClarificationStatus(status string) {
	allowed := map[string]bool{
		"open": true, "clear": true, "deferred": true,
	}
	if !allowed[status] {
		fail("Invalid clarification status. Allowed statuses: open, clear, deferred")
	}

	res := resolveMission("")
	if res.Safety != "safe" || res.Selected == nil {
		fail(formatResolutionBlock(res, "clarify-status"))
	}
	
	id := res.Selected.ID
	mPath := filepath.Join(missionDir(id), "mission.json")
	var m Mission
	readJson(mPath, &m)
	m.Clarification.Status = status
	if status == "clear" {
		m.Clarification.BlockingQuestions = 0
		m.Clarification.LastQuestion = nil
	}
	m.UpdatedAt = isoNow()
	writeJson(mPath, m)
	fmt.Printf("Spacecraft mission %s clarification: %s\n", id, status)
}
