package mission

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"spacecraft/internal/config"
)

func newTestConfig(t *testing.T) (*config.Config, *FSStore, func()) {
	t.Helper()
	dir, err := os.MkdirTemp("", "spacecraft-store-test-")
	if err != nil {
		t.Fatal(err)
	}
	cfg, err := config.NewConfig(dir)
	if err != nil {
		os.RemoveAll(dir)
		t.Fatal(err)
	}
	store := NewFSStore(cfg)
	cleanup := func() { os.RemoveAll(dir) }
	return cfg, store, cleanup
}

func writeTestMission(t *testing.T, store *FSStore, id, title, state string) {
	t.Helper()
	m := &Mission{
		ID:    id,
		Title: title,
		State: state,
	}
	if err := store.Create(m); err != nil {
		t.Fatal(err)
	}
}

// --- Current file ---

func TestFSStore_ReadCurrent_none(t *testing.T) {
	_, store, cleanup := newTestConfig(t)
	defer cleanup()

	id, err := store.ReadCurrent()
	if err != nil {
		t.Fatal(err)
	}
	if id != nil {
		t.Errorf("expected nil, got %v", *id)
	}
}

func TestFSStore_WriteReadCurrent(t *testing.T) {
	_, store, cleanup := newTestConfig(t)
	defer cleanup()

	if err := store.WriteCurrent("M07TEST01"); err != nil {
		t.Fatal(err)
	}
	id, err := store.ReadCurrent()
	if err != nil {
		t.Fatal(err)
	}
	if id == nil || *id != "M07TEST01" {
		t.Errorf("expected M07TEST01, got %v", id)
	}
}

func TestFSStore_ClearCurrent(t *testing.T) {
	_, store, cleanup := newTestConfig(t)
	defer cleanup()

	store.WriteCurrent("M07TEST02")
	store.ClearCurrent()

	id, err := store.ReadCurrent()
	if err != nil {
		t.Fatal(err)
	}
	if id != nil {
		t.Errorf("expected nil after clear, got %v", *id)
	}
}

// --- Session file ---

func TestFSStore_SessionRoundTrip(t *testing.T) {
	_, store, cleanup := newTestConfig(t)
	defer cleanup()

	key := "test-session-key-123"
	if err := store.WriteSession(key, "M07SESSION"); err != nil {
		t.Fatal(err)
	}

	id, err := store.ReadSession(key)
	if err != nil {
		t.Fatal(err)
	}
	if id == nil || *id != "M07SESSION" {
		t.Errorf("expected M07SESSION, got %v", id)
	}

	// Different key should not find
	other, err := store.ReadSession("other-key")
	if err != nil {
		t.Fatal(err)
	}
	if other != nil {
		t.Errorf("expected nil for different key, got %v", *other)
	}
}

func TestFSStore_ClearSession(t *testing.T) {
	_, store, cleanup := newTestConfig(t)
	defer cleanup()

	store.WriteSession("test-key", "M07SESS_CLEAR")
	store.ClearSession("test-key")

	id, err := store.ReadSession("test-key")
	if err != nil {
		t.Fatal(err)
	}
	if id != nil {
		t.Errorf("expected nil after clear, got %v", *id)
	}
}

// --- List ---

func TestFSStore_List_empty(t *testing.T) {
	_, store, cleanup := newTestConfig(t)
	defer cleanup()

	records, err := store.List()
	if err != nil {
		t.Fatal(err)
	}
	if len(records) != 0 {
		t.Errorf("expected 0 records, got %d", len(records))
	}
}

func TestFSStore_List_single(t *testing.T) {
	_, store, cleanup := newTestConfig(t)
	defer cleanup()

	writeTestMission(t, store, "M07TEST10", "Test mission", "draft")

	records, err := store.List()
	if err != nil {
		t.Fatal(err)
	}
	if len(records) != 1 {
		t.Fatalf("expected 1 record, got %d", len(records))
	}
	if records[0].ID != "M07TEST10" {
		t.Errorf("expected M07TEST10, got %s", records[0].ID)
	}
	if !records[0].Active {
		t.Error("expected active mission")
	}
}

func TestFSStore_List_skips_snipped(t *testing.T) {
	_, store, cleanup := newTestConfig(t)
	defer cleanup()

	writeTestMission(t, store, "M07TEST11", "Active", "draft")
	writeTestMission(t, store, "M07TEST12", "Shipped", "shipped")

	records, err := store.List()
	if err != nil {
		t.Fatal(err)
	}
	if len(records) != 2 {
		t.Fatalf("expected 2 records, got %d", len(records))
	}
	// Second should be shipped / inactive
	if records[1].ID == "M07TEST12" && records[1].Active {
		t.Error("shipped mission should not be active")
	}
}

// --- Load / Save / Create ---

func TestFSStore_CreateAndLoad(t *testing.T) {
	_, store, cleanup := newTestConfig(t)
	defer cleanup()

	m := &Mission{
		ID:    "M07TEST20",
		Title: "Created mission",
		State: "draft",
	}
	if err := store.Create(m); err != nil {
		t.Fatal(err)
	}

	loaded, err := store.Load("M07TEST20")
	if err != nil {
		t.Fatal(err)
	}
	if loaded.Title != "Created mission" {
		t.Errorf("expected 'Created mission', got %q", loaded.Title)
	}
	if loaded.State != "draft" {
		t.Errorf("expected 'draft', got %q", loaded.State)
	}

	// Ensure output dirs exist
	dir := store.MissionDir("M07TEST20")
	if !dirExists(filepath.Join(dir, "outputs")) {
		t.Error("outputs dir should exist after Create")
	}
	if !dirExists(filepath.Join(dir, "design")) {
		t.Error("design dir should exist after Create")
	}
}

func TestFSStore_Save_updates(t *testing.T) {
	_, store, cleanup := newTestConfig(t)
	defer cleanup()

	writeTestMission(t, store, "M07TEST21", "Original", "draft")

	m, _ := store.Load("M07TEST21")
	m.State = "planned"
	m.Title = "Updated"
	if err := store.Save(m); err != nil {
		t.Fatal(err)
	}

	loaded, _ := store.Load("M07TEST21")
	if loaded.State != "planned" {
		t.Errorf("expected 'planned', got %q", loaded.State)
	}
	if loaded.Title != "Updated" {
		t.Errorf("expected 'Updated', got %q", loaded.Title)
	}
}

// --- Plan ---

func TestFSStore_PlanRoundTrip(t *testing.T) {
	_, store, cleanup := newTestConfig(t)
	defer cleanup()

	writeTestMission(t, store, "M07TEST30", "Plan test", "planned")

	p := &Plan{
		MissionId: "M07TEST30",
		Tasks: []Task{
			{ID: strPtr("T01"), Title: strPtr("First task"), Status: strPtr("pending")},
		},
	}
	if err := store.SavePlan("M07TEST30", p); err != nil {
		t.Fatal(err)
	}

	loaded, err := store.LoadPlan("M07TEST30")
	if err != nil {
		t.Fatal(err)
	}
	if loaded.MissionId != "M07TEST30" {
		t.Errorf("expected M07TEST30, got %s", loaded.MissionId)
	}
	if len(loaded.Tasks) != 1 {
		t.Fatalf("expected 1 task, got %d", len(loaded.Tasks))
	}
	if *loaded.Tasks[0].ID != "T01" {
		t.Errorf("expected T01, got %s", *loaded.Tasks[0].ID)
	}
}

func TestFSStore_LoadPlan_missing(t *testing.T) {
	_, store, cleanup := newTestConfig(t)
	defer cleanup()

	_, err := store.LoadPlan("M07TEST31")
	if err == nil {
		t.Error("expected error for missing plan")
	}
}

// --- Review ---

func TestFSStore_ReviewRoundTrip(t *testing.T) {
	_, store, cleanup := newTestConfig(t)
	defer cleanup()

	writeTestMission(t, store, "M07TEST40", "Review test", "built")

	r := &Review{
		Status: strPtr("ready"),
		Findings: []Finding{
			{ID: strPtr("F01"), Summary: strPtr("All good"), Severity: strPtr("info")},
		},
	}
	if err := store.SaveReview("M07TEST40", r); err != nil {
		t.Fatal(err)
	}

	loaded, err := store.LoadReview("M07TEST40")
	if err != nil {
		t.Fatal(err)
	}
	if loaded.Status == nil || *loaded.Status != "ready" {
		t.Errorf("expected 'ready', got %v", loaded.Status)
	}
	if len(loaded.Findings) != 1 {
		t.Fatalf("expected 1 finding, got %d", len(loaded.Findings))
	}
}

// --- Evidence ---

func TestFSStore_CountEvidence_empty(t *testing.T) {
	_, store, cleanup := newTestConfig(t)
	defer cleanup()

	writeTestMission(t, store, "M07TEST50", "Evidence test", "planned")

	count, err := store.CountEvidence("M07TEST50")
	if err != nil {
		t.Fatal(err)
	}
	if count != 0 {
		t.Errorf("expected 0, got %d", count)
	}
}

func TestFSStore_AppendAndCountEvidence(t *testing.T) {
	_, store, cleanup := newTestConfig(t)
	defer cleanup()

	writeTestMission(t, store, "M07TEST51", "Evidence append", "planned")

	entry := &EvidenceEntry{
		ID:        "E07TEST01",
		Label:     "test evidence",
		Command:   "echo hello",
		ExitCode:  0,
		Stdout:    "path/to/stdout.txt",
		Stderr:    "path/to/stderr.txt",
		CreatedAt: "2026-07-07T00:00:00.000Z",
	}
	if err := store.AppendEvidence("M07TEST51", entry); err != nil {
		t.Fatal(err)
	}

	count, err := store.CountEvidence("M07TEST51")
	if err != nil {
		t.Fatal(err)
	}
	if count != 1 {
		t.Errorf("expected 1, got %d", count)
	}

	// Append another
	entry2 := &EvidenceEntry{
		ID:        "E07TEST02",
		Label:     "second",
		Command:   "true",
		ExitCode:  0,
		CreatedAt: "2026-07-07T00:00:01.000Z",
	}
	store.AppendEvidence("M07TEST51", entry2)

	count, _ = store.CountEvidence("M07TEST51")
	if count != 2 {
		t.Errorf("expected 2, got %d", count)
	}
}

func TestFSStore_ReadEvidenceEntries(t *testing.T) {
	_, store, cleanup := newTestConfig(t)
	defer cleanup()

	writeTestMission(t, store, "M07TEST52", "Evidence read", "planned")

	entry := &EvidenceEntry{
		ID:        "E07TEST10",
		Label:     "read test",
		Command:   "true",
		ExitCode:  0,
		CreatedAt: "2026-07-07T00:00:00.000Z",
	}
	store.AppendEvidence("M07TEST52", entry)

	entries, err := store.ReadEvidenceEntries("M07TEST52")
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1, got %d", len(entries))
	}
	if entries[0].ID != "E07TEST10" {
		t.Errorf("expected E07TEST10, got %s", entries[0].ID)
	}
}

// --- Artifact existence ---

func TestFSStore_SpecExists(t *testing.T) {
	_, store, cleanup := newTestConfig(t)
	defer cleanup()

	writeTestMission(t, store, "M07TEST60", "Spec test", "draft")

	if store.SpecExists("M07TEST60") {
		t.Error("spec should not exist yet")
	}

	dir := store.MissionDir("M07TEST60")
	os.WriteFile(filepath.Join(dir, "spec.md"), []byte("# Spec"), 0644)

	if !store.SpecExists("M07TEST60") {
		t.Error("spec should exist now")
	}
}

// --- ReserveEvidencePath ---

func TestFSStore_ReserveEvidencePath(t *testing.T) {
	_, store, cleanup := newTestConfig(t)
	defer cleanup()

	writeTestMission(t, store, "M07TEST70", "Reserve test", "draft")

	eid, stdoutPath, stderrPath, err := store.ReserveEvidencePath("M07TEST70")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(eid, "E") {
		t.Errorf("evidence id should start with E, got %s", eid)
	}
	if len(eid) != 9 {
		t.Errorf("expected 9-char id, got %d: %s", len(eid), eid)
	}
	if stdoutPath == "" || stderrPath == "" {
		t.Error("paths should not be empty")
	}
	// Files should have been created (empty)
	if !fileExists(stdoutPath) {
		t.Errorf("stdout file should exist: %s", stdoutPath)
	}
}

// --- MissionDir ---

func TestFSStore_MissionDir(t *testing.T) {
	cfg, store, cleanup := newTestConfig(t)
	defer cleanup()

	dir := store.MissionDir("M07TEST80")
	expected := filepath.Join(cfg.MissionsDir(), "M07TEST80")
	if dir != expected {
		t.Errorf("expected %s, got %s", expected, dir)
	}
}

// --- RemoveAll ---

func TestFSStore_RemoveAll(t *testing.T) {
	_, store, cleanup := newTestConfig(t)
	defer cleanup()

	writeTestMission(t, store, "M07TEST90", "Remove test", "draft")

	dir := store.MissionDir("M07TEST90")
	if !dirExists(dir) {
		t.Fatal("mission dir should exist")
	}

	if err := store.RemoveAll("M07TEST90"); err != nil {
		t.Fatal(err)
	}
	if dirExists(dir) {
		t.Error("mission dir should be removed")
	}
}

// --- List with branch metadata ---

func TestFSStore_List_branchMetadata(t *testing.T) {
	_, store, cleanup := newTestConfig(t)
	defer cleanup()

	branch := "feat/m07test100-feature"
	m := &Mission{
		ID:    "M07TEST100",
		Title: "Branch mission",
		State: "implementing",
		Git: GitBlock{
			WorkBranch: &branch,
		},
	}
	store.Create(m)

	records, err := store.List()
	if err != nil {
		t.Fatal(err)
	}
	if len(records) != 1 {
		t.Fatalf("expected 1, got %d", len(records))
	}
	if len(records[0].Branches) != 1 || records[0].Branches[0] != branch {
		t.Errorf("expected branch %s, got %v", branch, records[0].Branches)
	}
}

// --- helpers ---

func strPtr(s string) *string { return &s }

func dirExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.IsDir()
}

// Ensure store.go compiled without errors — build test
func TestFSStore_Compiles(t *testing.T) {
	_, store, cleanup := newTestConfig(t)
	defer cleanup()
	if store == nil {
		t.Fatal("store should not be nil")
	}
}
