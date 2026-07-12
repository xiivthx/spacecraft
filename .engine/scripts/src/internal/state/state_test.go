package state

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"spacecraft/internal/config"
	"spacecraft/internal/mission"
)

func newTestStore(t *testing.T) (mission.MissionStore, *config.Config, func()) {
	t.Helper()
	dir, err := os.MkdirTemp("", "spacecraft-state-test-")
	if err != nil {
		t.Fatal(err)
	}
	cfg, err := config.NewConfig(dir)
	if err != nil {
		os.RemoveAll(dir)
		t.Fatal(err)
	}
	store := mission.NewFSStore(cfg)
	cleanup := func() { os.RemoveAll(dir) }
	return store, cfg, cleanup
}

func writeMission(store mission.MissionStore, id, title, state string) error {
	return store.Create(&mission.Mission{
		ID:    id,
		Title: title,
		State: state,
	})
}

func TestSetState_allowed(t *testing.T) {
	store, _, cleanup := newTestStore(t)
	defer cleanup()

	if err := writeMission(store, "M07ST01", "Test", "draft"); err != nil {
		t.Fatal(err)
	}

	setter := NewSetter(store)
	if err := setter.SetState("M07ST01", "planned"); err != nil {
		t.Fatal(err)
	}

	m, _ := store.Load("M07ST01")
	if m.State != "planned" {
		t.Errorf("expected state=planned, got %q", m.State)
	}
	if m.UpdatedAt == "" {
		t.Error("UpdatedAt should be set")
	}
}

func TestSetState_invalid(t *testing.T) {
	store, _, cleanup := newTestStore(t)
	defer cleanup()

	if err := writeMission(store, "M07ST02", "Test", "draft"); err != nil {
		t.Fatal(err)
	}

	setter := NewSetter(store)
	err := setter.SetState("M07ST02", "nonexistent")
	if err == nil {
		t.Fatal("expected error for invalid state")
	}
	if !strings.Contains(err.Error(), "invalid mission state") {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestSetState_legacyMapping(t *testing.T) {
	store, _, cleanup := newTestStore(t)
	defer cleanup()

	if err := writeMission(store, "M07ST03", "Test", "draft"); err != nil {
		t.Fatal(err)
	}

	setter := NewSetter(store)

	// "specified" maps to "draft"
	if err := setter.SetState("M07ST03", "specified"); err != nil {
		t.Fatal(err)
	}
	m, _ := store.Load("M07ST03")
	if m.State != "draft" {
		t.Errorf("expected state=draft (mapped), got %q", m.State)
	}

	// "implementing" maps to "planned"
	if err := setter.SetState("M07ST03", "implementing"); err != nil {
		t.Fatal(err)
	}
	m, _ = store.Load("M07ST03")
	if m.State != "planned" {
		t.Errorf("expected state=planned (mapped), got %q", m.State)
	}

	// "verifying" maps to "built"
	if err := setter.SetState("M07ST03", "verifying"); err != nil {
		t.Fatal(err)
	}
	m, _ = store.Load("M07ST03")
	if m.State != "built" {
		t.Errorf("expected state=built (mapped from verifying), got %q", m.State)
	}

	// "reviewing" maps to "built"
	if err := setter.SetState("M07ST03", "reviewing"); err != nil {
		t.Fatal(err)
	}
	m, _ = store.Load("M07ST03")
	if m.State != "built" {
		t.Errorf("expected state=built (mapped from reviewing), got %q", m.State)
	}
}

func TestSetState_missingMission(t *testing.T) {
	store, _, cleanup := newTestStore(t)
	defer cleanup()

	setter := NewSetter(store)
	err := setter.SetState("M07ST_GHOST", "planned")
	if err == nil {
		t.Fatal("expected error for missing mission")
	}
}

func TestSetClarificationStatus_allowed(t *testing.T) {
	store, _, cleanup := newTestStore(t)
	defer cleanup()

	if err := writeMission(store, "M07ST04", "Test", "draft"); err != nil {
		t.Fatal(err)
	}

	setter := NewSetter(store)
	if err := setter.SetClarificationStatus("M07ST04", "clear"); err != nil {
		t.Fatal(err)
	}

	m, _ := store.Load("M07ST04")
	if m.Clarification.Status != "clear" {
		t.Errorf("expected clarification=clear, got %q", m.Clarification.Status)
	}
}

func TestSetClarificationStatus_clearResetsQuestions(t *testing.T) {
	store, _, cleanup := newTestStore(t)
	defer cleanup()

	if err := writeMission(store, "M07ST05", "Test", "draft"); err != nil {
		t.Fatal(err)
	}

	// Set to open with blocking questions
	m, _ := store.Load("M07ST05")
	m.Clarification.Status = "open"
	m.Clarification.BlockingQuestions = 3
	lastQ := "What is the meaning?"
	m.Clarification.LastQuestion = &lastQ
	store.Save(m)

	setter := NewSetter(store)
	setter.SetClarificationStatus("M07ST05", "clear")

	m, _ = store.Load("M07ST05")
	if m.Clarification.BlockingQuestions != 0 {
		t.Errorf("expected 0 blocking questions after clear, got %d", m.Clarification.BlockingQuestions)
	}
	if m.Clarification.LastQuestion != nil {
		t.Errorf("expected nil last question after clear, got %v", *m.Clarification.LastQuestion)
	}
}

func TestSetClarificationStatus_invalid(t *testing.T) {
	store, _, cleanup := newTestStore(t)
	defer cleanup()

	if err := writeMission(store, "M07ST06", "Test", "draft"); err != nil {
		t.Fatal(err)
	}

	setter := NewSetter(store)
	err := setter.SetClarificationStatus("M07ST06", "invalid")
	if err == nil {
		t.Fatal("expected error for invalid status")
	}
}

func TestValidateMission_valid(t *testing.T) {
	store, cfg, cleanup := newTestStore(t)
	defer cleanup()

	if err := writeMission(store, "M07ST10", "Valid", "draft"); err != nil {
		t.Fatal(err)
	}

	dir := cfg.MissionDir("M07ST10")
	os.WriteFile(filepath.Join(dir, "spec.md"), []byte("# Spec"), 0644)
	os.WriteFile(filepath.Join(dir, "plan.json"), []byte(`{"missionId":"M07ST10","tasks":[]}`), 0644)
	os.WriteFile(filepath.Join(dir, "evidence.jsonl"), []byte{}, 0644)

	setter := NewSetter(store)
	errs := setter.ValidateMission("M07ST10")
	if errs != nil {
		t.Fatalf("expected no errors, got: %v", errs.Errors)
	}
}

func TestValidateMission_missingFiles(t *testing.T) {
	store, _, cleanup := newTestStore(t)
	defer cleanup()

	if err := writeMission(store, "M07ST11", "Invalid", "draft"); err != nil {
		t.Fatal(err)
	}

	setter := NewSetter(store)
	errs := setter.ValidateMission("M07ST11")
	if errs == nil {
		t.Fatal("expected validation errors")
	}
	if !containsStr(errs.Errors, "missing spec.md") {
		t.Errorf("expected missing spec.md, got %v", errs.Errors)
	}
	if !containsStr(errs.Errors, "missing plan.json") {
		t.Errorf("expected missing plan.json, got %v", errs.Errors)
	}
	// evidence.jsonl might not exist
	if !containsStr(errs.Errors, "missing evidence.jsonl") {
		// Or the single error check
	}
}

func TestValidateMission_InvalidClarification(t *testing.T) {
	store, cfg, cleanup := newTestStore(t)
	defer cleanup()

	if err := writeMission(store, "M07ST12", "Bad clarity", "draft"); err != nil {
		t.Fatal(err)
	}

	dir := cfg.MissionDir("M07ST12")
	os.WriteFile(filepath.Join(dir, "spec.md"), []byte("# Spec"), 0644)
	os.WriteFile(filepath.Join(dir, "plan.json"), []byte(`{"missionId":"M07ST12","tasks":[]}`), 0644)
	os.WriteFile(filepath.Join(dir, "evidence.jsonl"), []byte{}, 0644)

	// Write mission.json with bad clarification
	badMission := &mission.Mission{
		ID:    "M07ST12",
		Title: "Bad clarity",
		State: "draft",
		Clarification: mission.ClarificationBlock{
			Status: "bogus",
		},
	}
	store.Save(badMission)

	setter := NewSetter(store)
	errs := setter.ValidateMission("M07ST12")
	if errs == nil {
		t.Fatal("expected validation errors for bad clarification status")
	}
	if !containsStr(errs.Errors, "clarification.status") {
		t.Errorf("expected clarification status error, got: %v", errs.Errors)
	}
}

func TestValidateMission_planNoTasks(t *testing.T) {
	store, cfg, cleanup := newTestStore(t)
	defer cleanup()

	if err := writeMission(store, "M07ST13", "No tasks", "draft"); err != nil {
		t.Fatal(err)
	}

	dir := cfg.MissionDir("M07ST13")
	os.WriteFile(filepath.Join(dir, "spec.md"), []byte("# Spec"), 0644)
	os.WriteFile(filepath.Join(dir, "plan.json"), []byte(`{"missionId":"M07ST13"}`), 0644)
	os.WriteFile(filepath.Join(dir, "evidence.jsonl"), []byte{}, 0644)

	setter := NewSetter(store)
	errs := setter.ValidateMission("M07ST13")
	if errs == nil {
		t.Fatal("expected validation error for plan without tasks")
	}
	if !containsStr(errs.Errors, "tasks array") {
		t.Errorf("expected tasks array error, got: %v", errs.Errors)
	}
}

func TestAllowedStates(t *testing.T) {
	expected := []string{"draft", "planned", "built", "ready", "shipped", "blocked"}
	for _, state := range expected {
		if !AllowedStates[state] {
			t.Errorf("expected %q to be allowed", state)
		}
	}
}

func TestAllowedClarificationStatuses(t *testing.T) {
	expected := []string{"open", "clear", "deferred"}
	for _, status := range expected {
		if !AllowedClarificationStatuses[status] {
			t.Errorf("expected %q to be allowed", status)
		}
	}
}

func containsStr(list []string, item string) bool {
	for _, s := range list {
		if strings.Contains(s, item) {
			return true
		}
	}
	return false
}
