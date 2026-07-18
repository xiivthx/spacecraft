package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func valCmd(args []string, spaceDir string) int {
	strict := false
	var mid string
	for _, a := range args {
		if a == "--strict" {
			strict = true
			continue
		}
		if strings.HasPrefix(a, "--") {
			continue
		}
		if mid == "" {
			mid = normalizeID(a)
		}
	}
	if mid == "" {
		mid = resolveMission(os.Getenv("PWD"))
	}
	if mid == "" {
		fmt.Fprintln(os.Stderr, "spacecraft validate: no mission id - pass as argument or run from feat/<id>/ branch")
		return 1
	}

	dir := missionDir(spaceDir, mid)
	ok := true

	check := func(file, desc string) {
		path := filepath.Join(dir, file)
		if _, err := os.Stat(path); os.IsNotExist(err) {
			fmt.Printf("x %-20s missing: %s\n", desc, path)
			ok = false
		} else {
			fmt.Printf("ok %-20s %s\n", desc, path)
		}
	}

	checkJSON := func(file, desc string) {
		path := filepath.Join(dir, file)
		data, err := os.ReadFile(path)
		if err != nil {
			fmt.Printf("x %-20s missing: %s\n", desc, path)
			ok = false
			return
		}
		var v any
		if err := json.Unmarshal(data, &v); err != nil {
			fmt.Printf("x %-20s invalid JSON: %s (%v)\n", desc, path, err)
			ok = false
			return
		}
		fmt.Printf("ok %-20s valid (%d bytes)\n", desc, len(data))
	}

	check("spec.md", "spec")
	checkJSON("mission.json", "mission")
	checkJSON("plan.json", "plan")

	if !validateEvidence(filepath.Join(dir, "evidence.jsonl"), strict) {
		ok = false
	}

	if strict {
		if !validateStrictPlanEvidence(dir) {
			ok = false
		}
	}

	if !ok {
		return 1
	}
	return 0
}

// validateEvidence checks that evidence.jsonl exists and every non-empty line is a
// JSON object carrying the required fields. When strict, also requires ≥1 entry
// and exitCode as a number on every entry. Returns false on any problem.
func validateEvidence(path string, strict bool) bool {
	data, err := os.ReadFile(path)
	if err != nil {
		fmt.Printf("x %-20s missing: %s\n", "evidence", path)
		return false
	}

	required := []string{"label", "command", "output", "ts"}
	entries := 0
	ok := true
	for i, line := range strings.Split(string(data), "\n") {
		if strings.TrimSpace(line) == "" {
			continue
		}
		var entry map[string]any
		if err := json.Unmarshal([]byte(line), &entry); err != nil {
			fmt.Printf("x %-20s line %d not valid JSON: %v\n", "evidence", i+1, err)
			return false
		}
		for _, field := range required {
			if _, present := entry[field]; !present {
				fmt.Printf("x %-20s line %d missing required field %q\n", "evidence", i+1, field)
				ok = false
			}
		}
		if strict && !isJSONNumber(entry["exitCode"]) {
			fmt.Printf("x %-20s line %d missing exitCode (number)\n", "evidence", i+1)
			ok = false
		}
		entries++
	}

	if strict && entries == 0 {
		fmt.Printf("x %-20s strict mode requires ≥1 evidence entry\n", "evidence")
		return false
	}

	if !ok {
		return false
	}
	fmt.Printf("ok %-20s %d entries\n", "evidence", entries)
	return true
}

// validateStrictPlanEvidence requires each done task to have ≥1 matching evidence
// label with exitCode == 0.
func validateStrictPlanEvidence(dir string) bool {
	data, err := os.ReadFile(filepath.Join(dir, "plan.json"))
	if err != nil {
		fmt.Printf("x %-20s cannot load plan.json: %v\n", "strict", err)
		return false
	}
	var plan struct {
		Tasks []struct {
			ID       string   `json:"id"`
			Status   string   `json:"status"`
			Evidence []string `json:"evidence"`
		} `json:"tasks"`
	}
	if err := json.Unmarshal(data, &plan); err != nil {
		fmt.Printf("x %-20s plan.json invalid: %v\n", "strict", err)
		return false
	}

	evData, err := os.ReadFile(filepath.Join(dir, "evidence.jsonl"))
	if err != nil {
		fmt.Printf("x %-20s cannot load evidence.jsonl: %v\n", "strict", err)
		return false
	}
	type evEntry struct {
		label    string
		exitCode float64
		hasExit  bool
	}
	var entries []evEntry
	for _, line := range strings.Split(string(evData), "\n") {
		if strings.TrimSpace(line) == "" {
			continue
		}
		var raw map[string]any
		if err := json.Unmarshal([]byte(line), &raw); err != nil {
			continue
		}
		label, _ := raw["label"].(string)
		e := evEntry{label: label}
		if n, ok := raw["exitCode"].(float64); ok {
			e.exitCode = n
			e.hasExit = true
		}
		entries = append(entries, e)
	}

	ok := true
	for _, task := range plan.Tasks {
		if task.Status != "done" {
			continue
		}
		allowed := map[string]bool{}
		for _, lab := range task.Evidence {
			allowed[lab] = true
		}
		matched := false
		for _, e := range entries {
			if allowed[e.label] && e.hasExit && e.exitCode == 0 {
				matched = true
				break
			}
		}
		if !matched {
			tid := task.ID
			if tid == "" {
				tid = "(unnamed)"
			}
			fmt.Printf("x %-20s done task %s missing passing evidence (exitCode 0) for labels %v\n",
				"strict", tid, task.Evidence)
			ok = false
		}
	}
	return ok
}
