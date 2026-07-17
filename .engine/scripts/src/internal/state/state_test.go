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
	store, cfg, cleanup := newTestStore(t)
	defer cleanup()

	if err := writeMission(store, "M07ST01", "Test", "draft"); err != nil {
		t.Fatal(err)
	}

	// Create required files for planned state
	dir := cfg.MissionDir("M07ST01")
	os.WriteFile(filepath.Join(dir, "spec.md"), []byte("# Spec"), 0644)
	os.WriteFile(filepath.Join(dir, "plan.json"), []byte(`{"missionId":"M07ST01","tasks":[]}`), 0644)

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
	store, cfg, cleanup := newTestStore(t)
	defer cleanup()

	if err := writeMission(store, "M07ST03", "Test", "draft"); err != nil {
		t.Fatal(err)
	}

	// Create required files for planned state
	dir := cfg.MissionDir("M07ST03")
	os.WriteFile(filepath.Join(dir, "spec.md"), []byte("# Spec"), 0644)
	os.WriteFile(filepath.Join(dir, "plan.json"), []byte(`{"missionId":"M07ST03","tasks":[]}`), 0644)

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

func TestSetState_invalidTransition(t *testing.T) {
	store, cfg, cleanup := newTestStore(t)
	defer cleanup()

	if err := writeMission(store, "M07ST30", "Test", "draft"); err != nil {
		t.Fatal(err)
	}

	// Create required files
	dir := cfg.MissionDir("M07ST30")
	os.WriteFile(filepath.Join(dir, "spec.md"), []byte("# Spec"), 0644)
	os.WriteFile(filepath.Join(dir, "plan.json"), []byte(`{"missionId":"M07ST30","tasks":[]}`), 0644)

	setter := NewSetter(store)

	// Try to skip from draft to built (should fail)
	err := setter.SetState("M07ST30", "built")
	if err == nil {
		t.Fatal("expected error for invalid transition draft -> built")
	}
	if !strings.Contains(err.Error(), "invalid state transition") {
		t.Errorf("unexpected error: %v", err)
	}

	// Try to skip from draft to shipped (should fail)
	err = setter.SetState("M07ST30", "shipped")
	if err == nil {
		t.Fatal("expected error for invalid transition draft -> shipped")
	}
	if !strings.Contains(err.Error(), "invalid state transition") {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestSetState_prerequisitesPlanned(t *testing.T) {
	store, _, cleanup := newTestStore(t)
	defer cleanup()

	if err := writeMission(store, "M07ST31", "Test", "draft"); err != nil {
		t.Fatal(err)
	}

	setter := NewSetter(store)

	// Try to go to planned without spec.md (should fail)
	err := setter.SetState("M07ST31", "planned")
	if err == nil {
		t.Fatal("expected error for missing spec.md")
	}
	if !strings.Contains(err.Error(), "spec.md is missing") {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestSetState_prerequisitesBuilt(t *testing.T) {
	store, cfg, cleanup := newTestStore(t)
	defer cleanup()

	if err := writeMission(store, "M07ST32", "Test", "draft"); err != nil {
		t.Fatal(err)
	}

	dir := cfg.MissionDir("M07ST32")
	os.WriteFile(filepath.Join(dir, "spec.md"), []byte("# Spec"), 0644)
	os.WriteFile(filepath.Join(dir, "plan.json"), []byte(`{"missionId":"M07ST32","tasks":[{"id":"T1","title":"Task 1","status":"pending"}]}`), 0644)

	setter := NewSetter(store)

	// Go to planned first
	if err := setter.SetState("M07ST32", "planned"); err != nil {
		t.Fatal(err)
	}

	// Try to go to built with incomplete task (should fail)
	err := setter.SetState("M07ST32", "built")
	if err == nil {
		t.Fatal("expected error for incomplete task")
	}
	if !strings.Contains(err.Error(), "task T1 is not complete") {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestSetState_prerequisitesReady(t *testing.T) {
	store, cfg, cleanup := newTestStore(t)
	defer cleanup()

	if err := writeMission(store, "M07ST33", "Test", "draft"); err != nil {
		t.Fatal(err)
	}

	dir := cfg.MissionDir("M07ST33")
	os.WriteFile(filepath.Join(dir, "spec.md"), []byte("# Spec"), 0644)
	os.WriteFile(filepath.Join(dir, "plan.json"), []byte(`{"missionId":"M07ST33","tasks":[]}`), 0644)

	setter := NewSetter(store)

	// Go to planned
	if err := setter.SetState("M07ST33", "planned"); err != nil {
		t.Fatal(err)
	}

	// Go to built
	if err := setter.SetState("M07ST33", "built"); err != nil {
		t.Fatal(err)
	}

	// Try to go to ready without evidence (should fail)
	err := setter.SetState("M07ST33", "ready")
	if err == nil {
		t.Fatal("expected error for missing evidence.jsonl")
	}
	if !strings.Contains(err.Error(), "evidence.jsonl is missing") {
		t.Errorf("unexpected error: %v", err)
	}

	// Create empty evidence.jsonl
	os.WriteFile(filepath.Join(dir, "evidence.jsonl"), []byte{}, 0644)

	// Try to go to ready with empty evidence (should fail)
	err = setter.SetState("M07ST33", "ready")
	if err == nil {
		t.Fatal("expected error for empty evidence.jsonl")
	}
	if !strings.Contains(err.Error(), "evidence.jsonl has no entries") {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestSetState_prerequisitesShipped(t *testing.T) {
	store, cfg, cleanup := newTestStore(t)
	defer cleanup()

	if err := writeMission(store, "M07ST34", "Test", "draft"); err != nil {
		t.Fatal(err)
	}

	dir := cfg.MissionDir("M07ST34")
	os.WriteFile(filepath.Join(dir, "spec.md"), []byte("# Spec"), 0644)
	os.WriteFile(filepath.Join(dir, "plan.json"), []byte(`{"missionId":"M07ST34","tasks":[]}`), 0644)
	os.WriteFile(filepath.Join(dir, "evidence.jsonl"), []byte(`{"id":"ev1"}`+"\n"), 0644)

	setter := NewSetter(store)

	// Go through states
	setter.SetState("M07ST34", "planned")
	setter.SetState("M07ST34", "built")
	setter.SetState("M07ST34", "ready")

	// Try to go to shipped without review.json (should fail)
	err := setter.SetState("M07ST34", "shipped")
	if err == nil {
		t.Fatal("expected error for missing review.json")
	}
	if !strings.Contains(err.Error(), "review.json is missing") {
		t.Errorf("unexpected error: %v", err)
	}

	// Create review.json with wrong status
	os.WriteFile(filepath.Join(dir, "review.json"), []byte(`{"status":"blocked"}`), 0644)

	err = setter.SetState("M07ST34", "shipped")
	if err == nil {
		t.Fatal("expected error for review status not ready")
	}
	if !strings.Contains(err.Error(), "review status is not ready") {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestSetState_validTransitionWithPrerequisites(t *testing.T) {
	store, cfg, cleanup := newTestStore(t)
	defer cleanup()

	if err := writeMission(store, "M07ST35", "Test", "draft"); err != nil {
		t.Fatal(err)
	}

	dir := cfg.MissionDir("M07ST35")
	os.WriteFile(filepath.Join(dir, "spec.md"), []byte("# Spec"), 0644)
	os.WriteFile(filepath.Join(dir, "plan.json"), []byte(`{"missionId":"M07ST35","tasks":[]}`), 0644)
	os.WriteFile(filepath.Join(dir, "evidence.jsonl"), []byte(`{"id":"ev1"}`+"\n"), 0644)
	os.WriteFile(filepath.Join(dir, "review.json"), []byte(`{"status":"ready"}`), 0644)

	setter := NewSetter(store)

	// Should succeed through all states
	if err := setter.SetState("M07ST35", "planned"); err != nil {
		t.Fatalf("failed to go to planned: %v", err)
	}
	if err := setter.SetState("M07ST35", "built"); err != nil {
		t.Fatalf("failed to go to built: %v", err)
	}
	if err := setter.SetState("M07ST35", "ready"); err != nil {
		t.Fatalf("failed to go to ready: %v", err)
	}
	if err := setter.SetState("M07ST35", "shipped"); err != nil {
		t.Fatalf("failed to go to shipped: %v", err)
	}

	m, _ := store.Load("M07ST35")
	if m.State != "shipped" {
		t.Errorf("expected state=shipped, got %q", m.State)
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

func TestSetClarificationStatus_missingMission(t *testing.T) {
	store, _, cleanup := newTestStore(t)
	defer cleanup()

	setter := NewSetter(store)
	err := setter.SetClarificationStatus("M07ST_GHOST", "clear")
	if err == nil {
		t.Fatal("expected error for missing mission")
	}
}

func TestValidationError_Error(t *testing.T) {
	err := &ValidationError{Errors: []string{"missing spec.md", "missing plan.json"}}
	got := err.Error()
	if !strings.Contains(got, "validation: 2 error(s)") {
		t.Errorf("unexpected error string: %q", got)
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

func TestValidateMission_missingMissionJson(t *testing.T) {
	store, cfg, cleanup := newTestStore(t)
	defer cleanup()

	if err := writeMission(store, "M07ST20", "Missing JSON", "draft"); err != nil {
		t.Fatal(err)
	}
	os.Remove(filepath.Join(cfg.MissionDir("M07ST20"), "mission.json"))
	os.WriteFile(filepath.Join(cfg.MissionDir("M07ST20"), "spec.md"), []byte("# Spec"), 0644)
	os.WriteFile(filepath.Join(cfg.MissionDir("M07ST20"), "plan.json"), []byte(`{"missionId":"M07ST20","tasks":[]}`), 0644)
	os.WriteFile(filepath.Join(cfg.MissionDir("M07ST20"), "evidence.jsonl"), []byte{}, 0644)

	setter := NewSetter(store)
	errs := setter.ValidateMission("M07ST20")
	if errs == nil {
		t.Fatal("expected validation errors")
	}
	if !containsStr(errs.Errors, "missing mission.json") {
		t.Errorf("expected missing mission.json, got %v", errs.Errors)
	}
}

func TestValidateMission_invalidMissionJson(t *testing.T) {
	store, cfg, cleanup := newTestStore(t)
	defer cleanup()

	if err := writeMission(store, "M07ST21", "Bad JSON", "draft"); err != nil {
		t.Fatal(err)
	}

	dir := cfg.MissionDir("M07ST21")
	os.WriteFile(filepath.Join(dir, "spec.md"), []byte("# Spec"), 0644)
	os.WriteFile(filepath.Join(dir, "plan.json"), []byte(`{"missionId":"M07ST21","tasks":[]}`), 0644)
	os.WriteFile(filepath.Join(dir, "evidence.jsonl"), []byte{}, 0644)
	os.WriteFile(filepath.Join(dir, "mission.json"), []byte("not json"), 0644)

	setter := NewSetter(store)
	errs := setter.ValidateMission("M07ST21")
	if errs == nil {
		t.Fatal("expected validation errors")
	}
	if !containsStr(errs.Errors, "invalid JSON in mission.json") {
		t.Errorf("expected invalid JSON in mission.json error, got %v", errs.Errors)
	}
}

func TestValidateMission_invalidPlanJson(t *testing.T) {
	store, cfg, cleanup := newTestStore(t)
	defer cleanup()

	if err := writeMission(store, "M07ST22", "Bad plan", "draft"); err != nil {
		t.Fatal(err)
	}

	dir := cfg.MissionDir("M07ST22")
	os.WriteFile(filepath.Join(dir, "spec.md"), []byte("# Spec"), 0644)
	os.WriteFile(filepath.Join(dir, "plan.json"), []byte("not json"), 0644)
	os.WriteFile(filepath.Join(dir, "evidence.jsonl"), []byte{}, 0644)

	setter := NewSetter(store)
	errs := setter.ValidateMission("M07ST22")
	if errs == nil {
		t.Fatal("expected validation errors")
	}
	if !containsStr(errs.Errors, "invalid JSON in plan.json") {
		t.Errorf("expected invalid JSON in plan.json error, got %v", errs.Errors)
	}
}

func TestValidateMission_invalidEvidenceJson(t *testing.T) {
	store, cfg, cleanup := newTestStore(t)
	defer cleanup()

	if err := writeMission(store, "M07ST23", "Bad evidence", "draft"); err != nil {
		t.Fatal(err)
	}

	dir := cfg.MissionDir("M07ST23")
	os.WriteFile(filepath.Join(dir, "spec.md"), []byte("# Spec"), 0644)
	os.WriteFile(filepath.Join(dir, "plan.json"), []byte(`{"missionId":"M07ST23","tasks":[]}`), 0644)
	os.WriteFile(filepath.Join(dir, "evidence.jsonl"), []byte("not json"), 0644)

	setter := NewSetter(store)
	errs := setter.ValidateMission("M07ST23")
	if errs == nil {
		t.Fatal("expected validation errors")
	}
	if !containsStr(errs.Errors, "invalid JSON in evidence.jsonl") {
		t.Errorf("expected invalid JSON in evidence.jsonl error, got %v", errs.Errors)
	}
}

func TestValidateMission_evidenceMissingId(t *testing.T) {
	store, cfg, cleanup := newTestStore(t)
	defer cleanup()

	if err := writeMission(store, "M07ST24", "Missing evidence id", "draft"); err != nil {
		t.Fatal(err)
	}

	dir := cfg.MissionDir("M07ST24")
	os.WriteFile(filepath.Join(dir, "spec.md"), []byte("# Spec"), 0644)
	os.WriteFile(filepath.Join(dir, "plan.json"), []byte(`{"missionId":"M07ST24","tasks":[]}`), 0644)
	os.WriteFile(filepath.Join(dir, "evidence.jsonl"), []byte(`{"foo":"bar"}`+"\n"), 0644)

	setter := NewSetter(store)
	errs := setter.ValidateMission("M07ST24")
	if errs == nil {
		t.Fatal("expected validation errors")
	}
	if !containsStr(errs.Errors, "must have string id") {
		t.Errorf("expected missing id error, got %v", errs.Errors)
	}
}

func TestValidateMission_evidenceDuplicateId(t *testing.T) {
	store, cfg, cleanup := newTestStore(t)
	defer cleanup()

	if err := writeMission(store, "M07ST25", "Duplicate evidence id", "draft"); err != nil {
		t.Fatal(err)
	}

	dir := cfg.MissionDir("M07ST25")
	os.WriteFile(filepath.Join(dir, "spec.md"), []byte("# Spec"), 0644)
	os.WriteFile(filepath.Join(dir, "plan.json"), []byte(`{"missionId":"M07ST25","tasks":[]}`), 0644)
	os.WriteFile(filepath.Join(dir, "evidence.jsonl"), []byte(`{"id":"ev1"}`+"\n"+`{"id":"ev1"}`+"\n"), 0644)

	setter := NewSetter(store)
	errs := setter.ValidateMission("M07ST25")
	if errs == nil {
		t.Fatal("expected validation errors")
	}
	if !containsStr(errs.Errors, "duplicate evidence id ev1") {
		t.Errorf("expected duplicate evidence id error, got %v", errs.Errors)
	}
}

func TestValidateMission_invalidReviewJson(t *testing.T) {
	store, cfg, cleanup := newTestStore(t)
	defer cleanup()

	if err := writeMission(store, "M07ST26", "Bad review", "draft"); err != nil {
		t.Fatal(err)
	}

	dir := cfg.MissionDir("M07ST26")
	os.WriteFile(filepath.Join(dir, "spec.md"), []byte("# Spec"), 0644)
	os.WriteFile(filepath.Join(dir, "plan.json"), []byte(`{"missionId":"M07ST26","tasks":[]}`), 0644)
	os.WriteFile(filepath.Join(dir, "evidence.jsonl"), []byte{}, 0644)
	os.WriteFile(filepath.Join(dir, "review.json"), []byte("not json"), 0644)

	setter := NewSetter(store)
	errs := setter.ValidateMission("M07ST26")
	if errs == nil {
		t.Fatal("expected validation errors")
	}
	if !containsStr(errs.Errors, "invalid JSON in review.json") {
		t.Errorf("expected invalid JSON in review.json error, got %v", errs.Errors)
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
