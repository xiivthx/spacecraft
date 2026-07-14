package eval

import (
	"os"
	"path/filepath"
	"testing"

	"spacecraft/internal/mission"
)

func TestRubricScoring(t *testing.T) {
	rubric := &EvalRubric{Dimensions: StdDimensions()}
	entries := []mission.EvidenceEntry{
		{ID: "e1", Label: "build-check", Command: "go build ./...", ExitCode: 0, Stdout: "Build successful with no errors"},
		{ID: "e2", Label: "test-check", Command: "go test ./...", ExitCode: 0, Stdout: "PASS: all tests passed. 10 tests ran."},
		{ID: "e3", Label: "lint-check", Command: "go vet ./...", ExitCode: 0, Stdout: "No issues found"},
	}

	scorecard, err := ScoreRubric(rubric, entries)
	if err != nil {
		t.Fatalf("ScoreRubric: %v", err)
	}
	if scorecard == nil {
		t.Fatal("scorecard is nil")
	}
	if len(scorecard.Scores) != 5 {
		t.Errorf("expected 5 dimension scores, got %d", len(scorecard.Scores))
	}

	taskSuccess := scorecard.Scores[0]
	if taskSuccess.Dimension != "task_success" {
		t.Errorf("first dimension should be task_success, got %s", taskSuccess.Dimension)
	}
	if taskSuccess.Score != 4 {
		t.Errorf("task_success: expected 4 (all entries pass), got %d", taskSuccess.Score)
	}
}

func TestRubricPartialFail(t *testing.T) {
	rubric := &EvalRubric{Dimensions: StdDimensions()}
	entries := []mission.EvidenceEntry{
		{ID: "e1", Label: "build", Command: "go build ./...", ExitCode: 0},
		{ID: "e2", Label: "test", Command: "go test ./...", ExitCode: 1},
		{ID: "e3", Label: "fail", Command: "unknown-cmd", ExitCode: 1},
	}

	scorecard, err := ScoreRubric(rubric, entries)
	if err != nil {
		t.Fatalf("ScoreRubric: %v", err)
	}

	taskSuccess := scorecard.Scores[0]
	if taskSuccess.Score >= 4 {
		t.Errorf("task_success with 1/3 passing should be < 4, got %d", taskSuccess.Score)
	}

	if scorecard.Average >= 4.0 {
		t.Errorf("average should be below 4.0 with failures, got %.1f", scorecard.Average)
	}
}

func TestRubricEmptyEntries(t *testing.T) {
	rubric := &EvalRubric{Dimensions: StdDimensions()}
	scorecard, err := ScoreRubric(rubric, nil)
	if err != nil {
		t.Fatalf("ScoreRubric on empty: %v", err)
	}
	for _, s := range scorecard.Scores {
		if s.Dimension == "hallucination" {
			if s.Score != 4 {
				t.Errorf("hallucination with empty entries should default to 4 (clean), got %d", s.Score)
			}
		} else {
			if s.Score != 0 {
				t.Errorf("dimension %s with empty entries should be 0, got %d", s.Dimension, s.Score)
			}
		}
	}
}

func TestRubricNilRubric(t *testing.T) {
	_, err := ScoreRubric(nil, []mission.EvidenceEntry{})
	if err == nil {
		t.Error("expected error for nil rubric")
	}
}

func TestRubricValidation(t *testing.T) {
	badRubric := &EvalRubric{
		Dimensions: []RubricDimension{
			{Name: "task_success", Description: "desc"},
			// missing 4 dimensions
		},
	}
	_, err := ScoreRubric(badRubric, nil)
	if err == nil {
		t.Error("expected error for rubric with wrong dimension count")
	}
}

func TestRubricTooManyDimensions(t *testing.T) {
	dims := StdDimensions()
	dims = append(dims, RubricDimension{Name: "extra", Description: "extra"})
	badRubric := &EvalRubric{Dimensions: dims}
	_, err := ScoreRubric(badRubric, nil)
	if err == nil {
		t.Error("expected error for rubric with too many dimensions")
	}
}

func TestRubricResolveContentFromFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "output.txt")
	content := "Build succeeded with no errors"
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("write output file: %v", err)
	}

	rubric := &EvalRubric{Dimensions: StdDimensions()}
	entries := []mission.EvidenceEntry{
		{ID: "e1", Label: "build-check", Command: "go build ./...", ExitCode: 0, Stdout: path},
	}

	scorecard, err := ScoreRubric(rubric, entries)
	if err != nil {
		t.Fatalf("ScoreRubric: %v", err)
	}
	if len(scorecard.Scores) != 5 {
		t.Errorf("expected 5 scores, got %d", len(scorecard.Scores))
	}
}
