package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

var validStates = map[string]bool{
	"active":      true,
	"planned":     true,
	"in_progress": true,
	"ready":       true,
	"blocked":     true,
	"shipped":     true,
}

var validTransitions = map[string][]string{
	"active":      {"planned", "blocked"},
	"planned":     {"in_progress", "blocked"},
	"in_progress": {"ready", "blocked"},
	"ready":       {"shipped", "blocked"},
	"blocked":     {"active", "in_progress"},
	"shipped":     {},
}

func stateCmd(args []string, spaceDir, mid string) int {
	var missionID, newState string

	switch len(args) {
	case 1:
		newState = args[0]
		if !validStates[newState] {
			fmt.Fprintf(os.Stderr, "spacecraft state: invalid state %q\n", newState)
			return 1
		}
		missionID = resolveActive(spaceDir, mid)
		if missionID == "" {
			fmt.Fprintln(os.Stderr, "spacecraft state: no active mission - provide mission-id or select one with 'spacecraft use'")
			return 1
		}
	case 2:
		missionID, newState = args[0], args[1]
	default:
		fmt.Fprintln(os.Stderr, "Usage: spacecraft set-state [mission-id] <new-state>")
		fmt.Fprintln(os.Stderr, "Valid states: active → planned → in_progress → ready → shipped")
		return 1
	}

	if !validStates[newState] {
		fmt.Fprintf(os.Stderr, "spacecraft state: invalid state %q\n", newState)
		return 1
	}

	missionPath := filepath.Join(missionDir(spaceDir, missionID), "mission.json")
	data, err := os.ReadFile(missionPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "spacecraft state: mission not found: %s\n", missionID)
		return 1
	}

	var m map[string]interface{}
	if err := json.Unmarshal(data, &m); err != nil {
		fmt.Fprintf(os.Stderr, "spacecraft state: invalid mission.json: %v\n", err)
		return 1
	}

	oldState, _ := m["state"].(string)
	if oldState == newState {
		fmt.Printf("%s already %s — no change\n", missionID, newState)
		return 0
	}

	if oldState != "" {
		allowed := validTransitions[oldState]
		ok := false
		for _, s := range allowed {
			if s == newState {
				ok = true
				break
			}
		}
		if !ok {
			fmt.Fprintf(os.Stderr, "spacecraft state: invalid transition %s → %s\n", oldState, newState)
			fmt.Fprintf(os.Stderr, "Allowed: %v\n", allowed)
			return 1
		}
	}

	m["state"] = newState
	out, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		fmt.Fprintln(os.Stderr, "spacecraft state:", err)
		return 1
	}

	if err := os.WriteFile(missionPath, append(out, '\n'), 0644); err != nil {
		fmt.Fprintln(os.Stderr, "spacecraft state:", err)
		return 1
	}

	fmt.Printf("%s: %s → %s\n", missionID, oldState, newState)
	return 0
}
