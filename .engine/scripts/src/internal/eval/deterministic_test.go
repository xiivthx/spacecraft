package eval

import (
	"testing"

	"spacecraft/internal/mission"
)

func TestDeterministicAllPass(t *testing.T) {
	entries := []mission.EvidenceEntry{
		{ID: "e1", Label: "T1-compile", Command: "go build ./...", ExitCode: 0},
		{ID: "e2", Label: "T1-test", Command: "go test ./...", ExitCode: 0},
		{ID: "e3", Label: "T1-lint", Command: "go vet ./...", ExitCode: 0},
	}
	result := RunDeterministic(entries)
	if !result.AllPassed {
		t.Errorf("expected all passed, got AllPassed=%v", result.AllPassed)
	}
	if len(result.Checks) != 3 {
		t.Errorf("expected 3 checks, got %d", len(result.Checks))
	}
	for _, c := range result.Checks {
		if !c.Passed {
			t.Errorf("check %s should pass, got passed=%v: %s", c.Name, c.Passed, c.Details)
		}
	}
	if result.TrajectorySummary == "" {
		t.Error("trajectory summary should not be empty")
	}
}

func TestDeterministicCompileFail(t *testing.T) {
	entries := []mission.EvidenceEntry{
		{ID: "e1", Label: "build-check", Command: "go build ./...", ExitCode: 1},
		{ID: "e2", Label: "test-check", Command: "go test ./...", ExitCode: 0},
	}
	result := RunDeterministic(entries)
	if result.AllPassed {
		t.Errorf("expected AllPassed=false since compile failed")
	}
	compile := result.Checks[0]
	if compile.Passed {
		t.Errorf("compile should fail with exit code 1")
	}
	if compile.Name != "compile" {
		t.Errorf("first check name should be compile, got %s", compile.Name)
	}
}

func TestDeterministicTestFail(t *testing.T) {
	entries := []mission.EvidenceEntry{
		{ID: "e1", Label: "test", Command: "go test ./...", ExitCode: 1},
	}
	result := RunDeterministic(entries)
	if result.AllPassed {
		t.Errorf("test failure should set AllPassed=false")
	}
}

func TestDeterministicEmpty(t *testing.T) {
	result := RunDeterministic(nil)
	if len(result.Checks) != 3 {
		t.Errorf("expected 3 checks for empty input, got %d", len(result.Checks))
	}
	for _, c := range result.Checks {
		if !c.Passed {
			t.Errorf("empty input: check %s should pass (no evidence), got passed=%v", c.Name, c.Passed)
		}
	}
}

func TestDeterministicLintDetection(t *testing.T) {
	entries := []mission.EvidenceEntry{
		{ID: "e1", Label: "lint", Command: "go vet ./...", ExitCode: 0},
	}
	result := RunDeterministic(entries)
	if !result.AllPassed {
		t.Errorf("lint pass should produce AllPassed=true")
	}
	lint := result.Checks[2]
	if lint.Name != "lint" {
		t.Errorf("third check should be lint, got %s", lint.Name)
	}
	if !lint.Passed {
		t.Errorf("lint check should pass")
	}
}

func TestDeterministicLintFail(t *testing.T) {
	entries := []mission.EvidenceEntry{
		{ID: "e1", Label: "lint", Command: "go vet ./...", ExitCode: 1},
	}
	result := RunDeterministic(entries)
	if result.AllPassed {
		t.Error("expected AllPassed=false when lint fails")
	}
	lint := result.Checks[2]
	if lint.Name != "lint" {
		t.Errorf("third check should be lint, got %s", lint.Name)
	}
	if lint.Passed {
		t.Error("lint check should fail with exit code 1")
	}
}

func TestTrajectorySummary(t *testing.T) {
	entries := []mission.EvidenceEntry{
		{ID: "e1", Label: "build", Command: "go build ./...", ExitCode: 0},
		{ID: "e2", Label: "test", Command: "go test ./...", ExitCode: 1},
		{ID: "e3", Label: "lint", Command: "go vet ./...", ExitCode: 0},
	}
	result := RunDeterministic(entries)
	if result.TrajectorySummary == "" {
		t.Error("trajectory summary should not be empty")
	}
	if !contains(result.TrajectorySummary, "3 evidence entries") {
		t.Errorf("trajectory summary should mention entry count: %s", result.TrajectorySummary)
	}
	if !contains(result.TrajectorySummary, "2 passed") {
		t.Errorf("trajectory summary should mention passed count: %s", result.TrajectorySummary)
	}
	if !contains(result.TrajectorySummary, "1 failed") {
		t.Errorf("trajectory summary should mention failed count: %s", result.TrajectorySummary)
	}
}

func contains(s, sub string) bool {
	return len(s) >= len(sub) && searchString(s, sub)
}

func searchString(s, sub string) bool {
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
