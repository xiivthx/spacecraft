package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestValidStatesAccepted(t *testing.T) {
	valid := []string{"active", "planned", "in_progress", "ready", "blocked", "shipped"}
	for _, s := range valid {
		if !validStates[s] {
			t.Errorf("validStates missing %q", s)
		}
	}
}

func TestStateTransitionHappyPath(t *testing.T) {
	dir := spaceRoot(t)
	id := "M07STH01"
	writeMission(t, dir, id, "active")

	steps := []string{"planned", "in_progress", "ready", "shipped"}
	for _, next := range steps {
		res := runCLI(t, dir, "state", id, next)
		if res.code != 0 {
			t.Fatalf("state %s → %s exit=%d\nstderr=%s", id, next, res.code, res.stderr)
		}
		got := readMissionState(t, dir, id)
		if got != next {
			t.Fatalf("after state %s: got %q", next, got)
		}
	}
}

func TestSetStateAliasMatchesState(t *testing.T) {
	dir := spaceRoot(t)
	id := "M07STH02"
	writeMission(t, dir, id, "active")

	res := runCLI(t, dir, "set-state", id, "planned")
	if stringsContainsUnknown(res.stderr) {
		t.Fatalf("set-state unknown command: %s", res.stderr)
	}
	if res.code != 0 {
		t.Fatalf("set-state exit=%d stderr=%s", res.code, res.stderr)
	}
	if got := readMissionState(t, dir, id); got != "planned" {
		t.Fatalf("state=%q, want planned", got)
	}
}

func TestSetStateSingleArgUsesCurrentMission(t *testing.T) {
	dir := spaceRoot(t)
	id := "M07STAR1"
	writeMission(t, dir, id, "active")

	use := runCLI(t, dir, "use", id)
	if use.code != 0 {
		t.Fatalf("use exit=%d stderr=%s", use.code, use.stderr)
	}

	res := runCLI(t, dir, "set-state", "planned")
	if res.code != 0 {
		t.Fatalf("set-state planned (single arg) exit=%d stderr=%s", res.code, res.stderr)
	}
	if got := readMissionState(t, dir, id); got != "planned" {
		t.Fatalf("state=%q, want planned", got)
	}
}

func TestStateSingleArgUsesCurrentMission(t *testing.T) {
	dir := spaceRoot(t)
	id := "M07STAR2"
	writeMission(t, dir, id, "active")

	use := runCLI(t, dir, "use", id)
	if use.code != 0 {
		t.Fatalf("use exit=%d stderr=%s", use.code, use.stderr)
	}

	res := runCLI(t, dir, "state", "planned")
	if res.code != 0 {
		t.Fatalf("state planned (single arg) exit=%d stderr=%s", res.code, res.stderr)
	}
	if got := readMissionState(t, dir, id); got != "planned" {
		t.Fatalf("state=%q, want planned", got)
	}
}

func TestSetStateExplicitTwoArgStillWorks(t *testing.T) {
	dir := spaceRoot(t)
	id := "M07STAR3"
	writeMission(t, dir, id, "active")

	res := runCLI(t, dir, "set-state", id, "planned")
	if res.code != 0 {
		t.Fatalf("set-state %s planned exit=%d stderr=%s", id, res.code, res.stderr)
	}
	if got := readMissionState(t, dir, id); got != "planned" {
		t.Fatalf("state=%q, want planned", got)
	}
}

func TestSetStateInvalidSingleArgFailsClearly(t *testing.T) {
	dir := spaceRoot(t)
	id := "M07STAR4"
	writeMission(t, dir, id, "active")
	use := runCLI(t, dir, "use", id)
	if use.code != 0 {
		t.Fatalf("use exit=%d stderr=%s", use.code, use.stderr)
	}

	res := runCLI(t, dir, "set-state", "not-a-state")
	if res.code == 0 {
		t.Fatal("invalid single-arg state must be rejected")
	}
	errOut := strings.ToLower(res.stderr)
	if !strings.Contains(errOut, "invalid state") {
		t.Fatalf("expected clear invalid-state error, got stderr=%s", res.stderr)
	}
	if got := readMissionState(t, dir, id); got != "active" {
		t.Fatalf("state mutated to %q after rejected single-arg", got)
	}
}

func TestInvalidStateRejected(t *testing.T) {
	dir := spaceRoot(t)
	id := "M07STI01"
	writeMission(t, dir, id, "active")

	res := runCLI(t, dir, "state", id, "nonexistent")
	if res.code == 0 {
		t.Fatal("invalid state must be rejected")
	}
}

func TestInvalidTransitionsRejected(t *testing.T) {
	tests := []struct {
		name string
		from string
		to   string
	}{
		{"skip active to in_progress", "active", "in_progress"},
		{"skip active to ready", "active", "ready"},
		{"skip active to shipped", "active", "shipped"},
		{"skip planned to ready", "planned", "ready"},
		{"skip planned to shipped", "planned", "shipped"},
		{"skip in_progress to shipped", "in_progress", "shipped"},
		{"shipped is terminal", "shipped", "active"},
		{"shipped to ready", "shipped", "ready"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := spaceRoot(t)
			id := "M07STI02"
			writeMission(t, dir, id, tt.from)
			res := runCLI(t, dir, "state", id, tt.to)
			if res.code == 0 {
				t.Fatalf("expected reject %s → %s", tt.from, tt.to)
			}
			if got := readMissionState(t, dir, id); got != tt.from {
				t.Fatalf("state mutated to %q after rejected transition", got)
			}
		})
	}
}

func TestBlockedTransitions(t *testing.T) {
	dir := spaceRoot(t)
	id := "M07STB01"
	writeMission(t, dir, id, "in_progress")

	res := runCLI(t, dir, "state", id, "blocked")
	if res.code != 0 {
		t.Fatalf("in_progress → blocked exit=%d stderr=%s", res.code, res.stderr)
	}
	res = runCLI(t, dir, "state", id, "in_progress")
	if res.code != 0 {
		t.Fatalf("blocked → in_progress exit=%d stderr=%s", res.code, res.stderr)
	}
}

func readMissionState(t *testing.T, root, id string) string {
	t.Helper()
	data, err := os.ReadFile(filepath.Join(root, ".space", "missions", id, "mission.json"))
	if err != nil {
		t.Fatal(err)
	}
	var m map[string]any
	if err := json.Unmarshal(data, &m); err != nil {
		t.Fatal(err)
	}
	s, _ := m["state"].(string)
	return s
}

func stringsContainsUnknown(s string) bool {
	return strings.Contains(strings.ToLower(s), "unknown command")
}
