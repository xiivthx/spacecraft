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
	if len(args) >= 1 {
		mid = args[0]
	} else {
		mid = resolveMission(os.Getenv("PWD"))
	}
	if mid == "" {
		fmt.Fprintln(os.Stderr, "spacecraft val: no mission id — pass as argument or run from feat/<id>/ branch")
		return 1
	}

	dir := missionDir(spaceDir, mid)
	ok := true

	check := func(file, desc string) {
		path := filepath.Join(dir, file)
		if _, err := os.Stat(path); os.IsNotExist(err) {
			fmt.Printf("✗ %-20s missing: %s\n", desc, path)
			ok = false
		} else {
			fmt.Printf("✓ %-20s %s\n", desc, path)
		}
	}

	checkJSON := func(file, desc string) {
		path := filepath.Join(dir, file)
		data, err := os.ReadFile(path)
		if err != nil {
			if os.IsNotExist(err) {
				fmt.Printf("✗ %-20s missing: %s\n", desc, path)
				ok = false
			} else {
				fmt.Printf("✗ %-20s read error: %s\n", desc, path)
				ok = false
			}
			return
		}
		var v interface{}
		if err := json.Unmarshal(data, &v); err != nil {
			fmt.Printf("✗ %-20s invalid JSON: %s (%v)\n", desc, path, err)
			ok = false
			return
		}
		fmt.Printf("✓ %-20s valid (%d bytes)\n", desc, len(data))
	}

	check("spec.md", "spec")
	checkJSON("mission.json", "mission")
	checkJSON("plan.json", "plan")

	evidencePath := filepath.Join(dir, "evidence.jsonl")
	if data, err := os.ReadFile(evidencePath); err == nil {
		lines := strings.Count(strings.TrimSpace(string(data)), "\n") + 1
		fmt.Printf("✓ %-20s %d entries\n", "evidence", lines)
	} else {
		fmt.Printf("✗ %-20s missing: %s\n", "evidence", evidencePath)
		ok = false
	}

	if !ok {
		return 1
	}
	return 0
}
