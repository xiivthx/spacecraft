package mission

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"spacecraft/internal/config"
	"spacecraft/internal/id"
	"spacecraft/internal/util"
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
	if !util.Exists(stdoutPath) {
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

// --- artifact existence ---

func TestFSStore_ArtifactExists(t *testing.T) {
	_, store, cleanup := newTestConfig(t)
	defer cleanup()

	writeTestMission(t, store, "M07ART01", "Artifact test", "draft")

	if store.PlanExists("M07ART01") {
		t.Error("plan should not exist")
	}
	if store.QuestionsExists("M07ART01") {
		t.Error("questions should not exist")
	}
	if store.DecisionsExists("M07ART01") {
		t.Error("decisions should not exist")
	}
	if store.ReviewJSONExists("M07ART01") {
		t.Error("review.json should not exist")
	}
	if store.ReviewMDExists("M07ART01") {
		t.Error("review.md should not exist")
	}

	dir := store.MissionDir("M07ART01")
	os.WriteFile(filepath.Join(dir, "plan.json"), []byte("{}"), 0644)
	os.WriteFile(filepath.Join(dir, "questions.md"), []byte("q"), 0644)
	os.WriteFile(filepath.Join(dir, "decisions.md"), []byte("d"), 0644)
	os.WriteFile(filepath.Join(dir, "review.json"), []byte("{}"), 0644)
	os.WriteFile(filepath.Join(dir, "review.md"), []byte("r"), 0644)

	if !store.PlanExists("M07ART01") {
		t.Error("plan should exist")
	}
	if !store.QuestionsExists("M07ART01") {
		t.Error("questions should exist")
	}
	if !store.DecisionsExists("M07ART01") {
		t.Error("decisions should exist")
	}
	if !store.DesignExists("M07ART01") {
		t.Error("design should exist")
	}
	if !store.ReviewJSONExists("M07ART01") {
		t.Error("review.json should exist")
	}
	if !store.ReviewMDExists("M07ART01") {
		t.Error("review.md should exist")
	}
}

// --- file operations ---

func TestFSStore_ReadWriteFile(t *testing.T) {
	_, store, cleanup := newTestConfig(t)
	defer cleanup()

	writeTestMission(t, store, "M07FILE01", "File test", "draft")

	data := []byte("hello mission")
	if err := store.WriteFile("M07FILE01", "notes.txt", data); err != nil {
		t.Fatal(err)
	}

	got, err := store.ReadFile("M07FILE01", "notes.txt")
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != "hello mission" {
		t.Errorf("expected 'hello mission', got %q", string(got))
	}
}

// --- archive ---

func TestFSStore_ArchiveMission(t *testing.T) {
	cfg, store, cleanup := newTestConfig(t)
	defer cleanup()

	writeTestMission(t, store, "M07ARCH01", "Archive test", "shipped")

	compactM := CompactMission{ID: "M07ARCH01", Title: "Archive test", State: "shipped"}
	compactP := CompactPlan{MissionID: "M07ARCH01"}
	evidence := []CompactEvidenceEntry{
		{ID: "E07ARCH01", Command: "true", ExitCode: 0, CreatedAt: "2026-07-07T00:00:00.000Z"},
	}
	review := &Review{Status: strPtr("ready")}

	archiveDir := cfg.ArchiveDir()
	if err := store.ArchiveMission("M07ARCH01", archiveDir, compactM, compactP, evidence, review); err != nil {
		t.Fatal(err)
	}

	dest := filepath.Join(archiveDir, "M07ARCH01")
	if !util.Exists(filepath.Join(dest, "mission.json")) {
		t.Error("archived mission.json should exist")
	}
	if !util.Exists(filepath.Join(dest, "plan.json")) {
		t.Error("archived plan.json should exist")
	}
	if !util.Exists(filepath.Join(dest, "evidence.jsonl")) {
		t.Error("archived evidence.jsonl should exist")
	}
	if !util.Exists(filepath.Join(dest, "review.json")) {
		t.Error("archived review.json should exist")
	}

	data, _ := os.ReadFile(filepath.Join(dest, "evidence.jsonl"))
	if !strings.Contains(string(data), "E07ARCH01") {
		t.Errorf("evidence.jsonl should contain E07ARCH01, got %q", string(data))
	}
}

func TestFSStore_ArchiveMission_alreadyExists(t *testing.T) {
	cfg, store, cleanup := newTestConfig(t)
	defer cleanup()

	archiveDir := cfg.ArchiveDir()
	os.MkdirAll(filepath.Join(archiveDir, "M07ARCH02"), 0755)

	err := store.ArchiveMission("M07ARCH02", archiveDir, CompactMission{ID: "M07ARCH02"}, CompactPlan{}, nil, nil)
	if err == nil {
		t.Error("expected error when archive already exists")
	}
}

// --- current/session edge cases ---

func TestFSStore_WriteCurrent_error(t *testing.T) {
	_, store, cleanup := newTestConfig(t)
	defer cleanup()

	if err := os.WriteFile(store.cfg.SpaceDir(), []byte{}, 0644); err != nil {
		t.Fatal(err)
	}
	if err := store.WriteCurrent("M07ERR"); err == nil {
		t.Error("expected error when .space is a file")
	}
}

func TestFSStore_ClearCurrent_error(t *testing.T) {
	_, store, cleanup := newTestConfig(t)
	defer cleanup()

	if err := os.WriteFile(store.cfg.SpaceDir(), []byte{}, 0644); err != nil {
		t.Fatal(err)
	}
	if err := store.ClearCurrent(); err == nil {
		t.Error("expected error when .space is a file")
	}
}

func TestFSStore_SessionFilePath(t *testing.T) {
	cfg, store, cleanup := newTestConfig(t)
	defer cleanup()

	path := store.SessionFilePath("my session key")
	if path == "" {
		t.Fatal("expected path for normal key")
	}
	if !strings.HasSuffix(path, ".current") {
		t.Errorf("expected .current suffix, got %s", path)
	}
	if filepath.Dir(path) != filepath.Join(cfg.SpaceDir(), "sessions") {
		t.Errorf("unexpected dir: %s", filepath.Dir(path))
	}

	if store.SessionFilePath("!!!") != "" {
		t.Error("expected empty path for key that slugifies to empty")
	}

	long := strings.Repeat("a", 100)
	longPath := store.SessionFilePath(long)
	if len(filepath.Base(longPath)) > 80+len(".current") {
		t.Errorf("slug too long: %s", filepath.Base(longPath))
	}
}

func TestFSStore_WriteSession_emptyKey(t *testing.T) {
	_, store, cleanup := newTestConfig(t)
	defer cleanup()

	if err := store.WriteSession("???", "M07ID"); err == nil {
		t.Error("expected error for empty session slug")
	}
	if err := store.ClearSession("???"); err != nil {
		t.Errorf("expected nil for clear of empty session slug, got %v", err)
	}
}

func TestFSStore_Create_error(t *testing.T) {
	dir, err := os.MkdirTemp("", "spacecraft-create-error-")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	block := filepath.Join(dir, "block")
	if err := os.WriteFile(block, []byte{}, 0644); err != nil {
		t.Fatal(err)
	}

	cfg, err := config.NewConfig(dir, config.WithMissionsDir(filepath.Join(block, "missions")))
	if err != nil {
		t.Fatal(err)
	}
	store := NewFSStore(cfg)

	m := &Mission{ID: "M07ERR", Title: "Err", State: "draft"}
	if err := store.Create(m); err == nil {
		t.Error("expected error when missions dir parent is a file")
	}
}

// --- reserve evidence path ---

func TestFSStore_ReserveEvidencePath_collision(t *testing.T) {
	_, store, cleanup := newTestConfig(t)
	defer cleanup()

	writeTestMission(t, store, "M07RES01", "Reserve collision", "draft")
	outputsDir := filepath.Join(store.MissionDir("M07RES01"), "outputs")

	now := time.Now()
	for offset := -5; offset <= 5; offset++ {
		eid, err := id.EvidenceId(now.Add(time.Duration(offset) * time.Millisecond))
		if err != nil {
			t.Fatal(err)
		}
		os.WriteFile(filepath.Join(outputsDir, eid+".stdout.txt"), []byte{}, 0644)
	}

	eid, stdoutPath, stderrPath, err := store.ReserveEvidencePath("M07RES01")
	if err != nil {
		t.Fatal(err)
	}
	if eid == "" || stdoutPath == "" || stderrPath == "" {
		t.Error("expected non-empty paths")
	}
	if !util.Exists(stdoutPath) {
		t.Error("reserved stdout file should exist")
	}
}

func TestFSStore_ReserveEvidencePath_error(t *testing.T) {
	_, store, cleanup := newTestConfig(t)
	defer cleanup()

	writeTestMission(t, store, "M07RES02", "Reserve error", "draft")
	outputsDir := filepath.Join(store.MissionDir("M07RES02"), "outputs")

	if err := os.Chmod(outputsDir, 0000); err != nil {
		t.Fatal(err)
	}
	defer os.Chmod(outputsDir, 0755)

	_, _, _, err := store.ReserveEvidencePath("M07RES02")
	if err == nil {
		t.Error("expected error when outputs dir is not writable")
	}
}

// --- branch metadata ---

func TestFSStore_List_branchNamesVariations(t *testing.T) {
	_, store, cleanup := newTestConfig(t)
	defer cleanup()

	branch := "feat/branch"
	workBranch := "feat/work"
	gitWorkBranch := "feat/gitwork"
	m := &Mission{
		ID:         "M07BR01",
		Title:      "Branches",
		State:      "draft",
		Branch:     &branch,
		WorkBranch: &workBranch,
		Git:        GitBlock{WorkBranch: &gitWorkBranch},
	}
	if err := store.Create(m); err != nil {
		t.Fatal(err)
	}

	records, err := store.List()
	if err != nil {
		t.Fatal(err)
	}
	if len(records) != 1 {
		t.Fatalf("expected 1 record, got %d", len(records))
	}
	if len(records[0].Branches) != 3 {
		t.Errorf("expected 3 branches, got %d: %v", len(records[0].Branches), records[0].Branches)
	}

	empty := ""
	m2 := &Mission{
		ID:         "M07BR02",
		Title:      "Empty branches",
		State:      "draft",
		Branch:     &empty,
		WorkBranch: &empty,
		Git:        GitBlock{WorkBranch: &empty},
	}
	if err := store.Create(m2); err != nil {
		t.Fatal(err)
	}

	records, err = store.List()
	if err != nil {
		t.Fatal(err)
	}
	if len(records) != 2 {
		t.Fatalf("expected 2 records, got %d", len(records))
	}
	for _, r := range records {
		if r.ID == "M07BR02" && len(r.Branches) != 0 {
			t.Errorf("expected no branches for empty strings, got %v", r.Branches)
		}
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
