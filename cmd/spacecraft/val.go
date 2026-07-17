package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func valCmd(args []string, spaceDir string) int {
	var mid string
	for _, a := range args {
		if !strings.HasPrefix(a, "--") {
			mid = normalizeID(a)
			break
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

	if !validateEvidence(filepath.Join(dir, "evidence.jsonl")) {
		ok = false
	}

	if !ok {
		return 1
	}
	return 0
}

// validateEvidence checks that evidence.jsonl exists and every non-empty line is a
// JSON object carrying the required fields. Returns false on any problem.
func validateEvidence(path string) bool {
	data, err := os.ReadFile(path)
	if err != nil {
		fmt.Printf("x %-20s missing: %s\n", "evidence", path)
		return false
	}

	required := []string{"label", "command", "output", "ts"}
	entries := 0
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
				return false
			}
		}
		entries++
	}

	fmt.Printf("ok %-20s %d entries\n", "evidence", entries)
	return true
}
