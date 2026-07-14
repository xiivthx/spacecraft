package eval

import (
	"os"
	"path/filepath"
	"testing"

	"spacecraft/internal/config"
	"spacecraft/internal/mission"
)

func ptr(s string) *string { return &s }

func setupRunnerTest(t *testing.T) (store mission.MissionStore, evalsDir, missionID string) {
	t.Helper()
	root := t.TempDir()
	cfg, err := config.NewConfig(root)
	if err != nil {
		t.Fatalf("config: %v", err)
	}
	store = mission.NewFSStore(cfg)
	missionID = "M123"

	if err := os.MkdirAll(cfg.MissionDir(missionID), 0755); err != nil {
		t.Fatalf("mkdir mission: %v", err)
	}
	if err := os.MkdirAll(filepath.Join(cfg.MissionDir(missionID), "outputs"), 0755); err != nil {
		t.Fatalf("mkdir outputs: %v", err)
	}

	plan := &mission.Plan{
		MissionId: missionID,
		Tasks:     []mission.Task{{ID: ptr("t1")}, {ID: ptr("t2")}},
	}
	if err := store.SavePlan(missionID, plan); err != nil {
		t.Fatalf("save plan: %v", err)
	}

	evalsDir = cfg.EvalsDir()
	if err := Init(evalsDir, missionID); err != nil {
		t.Fatalf("init evals: %v", err)
	}

	return store, evalsDir, missionID
}

func TestNewRunner(t *testing.T) {
	var store mission.MissionStore = mission.NewFSStore(nil)
	r := NewRunner(store, "/tmp/evals")
	if r.Store != store {
		t.Error("runner store mismatch")
	}
	if r.EvalsDir != "/tmp/evals" {
		t.Errorf("expected evals dir /tmp/evals, got %s", r.EvalsDir)
	}
}

func TestFilterEvalEntries(t *testing.T) {
	evalType := "eval"
	otherType := "test"
	entries := []mission.EvidenceEntry{
		{ID: "e1", Type: &evalType},
		{ID: "e2", Type: &otherType},
		{ID: "e3"},
	}
	filtered := filterEvalEntries(entries)
	if len(filtered) != 2 {
		t.Errorf("expected 2 filtered entries, got %d", len(filtered))
	}
	for _, e := range filtered {
		if e.Type != nil && *e.Type == "eval" {
			t.Errorf("eval entry %s should be filtered out", e.ID)
		}
	}
}

func TestEnsureFileExists(t *testing.T) {
	path := filepath.Join(t.TempDir(), "new.txt")
	if err := ensureFileExists(path); err != nil {
		t.Fatalf("ensureFileExists: %v", err)
	}
	if _, err := os.Stat(path); err != nil {
		t.Errorf("expected file to exist: %v", err)
	}
	if err := ensureFileExists(path); err != nil {
		t.Fatalf("ensureFileExists idempotent: %v", err)
	}
}

func TestRunnerRun(t *testing.T) {
	store, evalsDir, missionID := setupRunnerTest(t)

	entries := []mission.EvidenceEntry{
		{ID: "e1", Label: "build", Command: "go build ./...", ExitCode: 0},
		{ID: "e2", Label: "test", Command: "go test ./...", ExitCode: 0},
	}
	for _, e := range entries {
		if err := store.AppendEvidence(missionID, &e); err != nil {
			t.Fatalf("append evidence: %v", err)
		}
	}

	if err := os.WriteFile(filepath.Join(evalsDir, missionID, "config.json"), []byte(`{"coverageThreshold":0.25}`), 0644); err != nil {
		t.Fatalf("write config: %v", err)
	}
	if err := os.WriteFile(filepath.Join(evalsDir, missionID, "dataset.json"), []byte(`{"examples":[{"id":"ex1","label":"example","expectedPass":true}]}`), 0644); err != nil {
		t.Fatalf("write dataset: %v", err)
	}

	runner := NewRunner(store, evalsDir)
	result, err := runner.Run(missionID)
	if err != nil {
		t.Fatalf("Run failed: %v", err)
	}
	if result == nil || result.EvalResult == nil {
		t.Fatal("expected non-nil result")
	}
	if result.EvalResult.MissionID != missionID {
		t.Errorf("expected mission ID %s, got %s", missionID, result.EvalResult.MissionID)
	}
	if result.Entry == nil {
		t.Error("expected evidence entry to be created")
	}
	if !result.EvalResult.CoverageSatisfied {
		t.Errorf("coverage should be satisfied, got %v", result.EvalResult.Coverage)
	}
}

func TestRunnerRunReservePathError(t *testing.T) {
	root := t.TempDir()
	cfg, err := config.NewConfig(root)
	if err != nil {
		t.Fatalf("config: %v", err)
	}
	store := mission.NewFSStore(cfg)
	missionID := "M123"
	if err := os.MkdirAll(cfg.MissionDir(missionID), 0755); err != nil {
		t.Fatalf("mkdir mission: %v", err)
	}

	runner := NewRunner(store, cfg.EvalsDir())
	_, err = runner.Run(missionID)
	if err == nil {
		t.Error("expected error when outputs directory does not exist")
	}
}
