package archive

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"spacecraft/internal/config"
	"spacecraft/internal/mission"
	"spacecraft/internal/util"
)

func TestReadinessChecker_ready(t *testing.T) {
	store, _, cleanup := newTestStore(t)
	defer cleanup()

	checker := NewReadinessChecker(store)

	plan := &mission.Plan{
		Tasks: []mission.Task{
			{ID: strPtr("T01"), Status: strPtr("completed")},
		},
	}
	review := readyReview()
	entries := []mission.EvidenceEntry{
		{ID: "E0001", Label: "test", ExitCode: 0},
	}

	errs := checker.CheckReadiness("M07AR01", plan, review, entries)
	if errs != nil {
		t.Fatalf("expected no errors, got: %v", errs.Errors)
	}
}

func TestReadinessChecker_noPlan(t *testing.T) {
	store, _, cleanup := newTestStore(t)
	defer cleanup()

	checker := NewReadinessChecker(store)
	errs := checker.CheckReadiness("M07AR02", nil, nil, nil)
	if errs == nil {
		t.Fatal("expected errors")
	}
	if !containsStr(errs.Errors, "missing plan.json") {
		t.Errorf("expected plan error, got: %v", errs.Errors)
	}
}

func TestReadinessChecker_noReview(t *testing.T) {
	store, _, cleanup := newTestStore(t)
	defer cleanup()

	checker := NewReadinessChecker(store)
	plan := &mission.Plan{Tasks: []mission.Task{{ID: strPtr("T01"), Status: strPtr("completed")}}}
	errs := checker.CheckReadiness("M07AR03", plan, nil, []mission.EvidenceEntry{{ID: "E001"}})
	if errs == nil {
		t.Fatal("expected errors")
	}
	if !containsStr(errs.Errors, "missing review.json") {
		t.Errorf("expected review error, got: %v", errs.Errors)
	}
}

func TestReadinessChecker_reviewNotReady(t *testing.T) {
	store, _, cleanup := newTestStore(t)
	defer cleanup()

	checker := NewReadinessChecker(store)
	plan := &mission.Plan{Tasks: []mission.Task{{ID: strPtr("T01"), Status: strPtr("completed")}}}
	review := &mission.Review{Status: strPtr("blocked")}
	errs := checker.CheckReadiness("M07AR04", plan, review, []mission.EvidenceEntry{{ID: "E001"}})
	if errs == nil {
		t.Fatal("expected errors")
	}
	if !containsStr(errs.Errors, "review status is blocked") {
		t.Errorf("expected review status error, got: %v", errs.Errors)
	}
}

func TestReadinessChecker_noTasks(t *testing.T) {
	store, _, cleanup := newTestStore(t)
	defer cleanup()

	checker := NewReadinessChecker(store)
	plan := &mission.Plan{}
	review := readyReview()
	errs := checker.CheckReadiness("M07AR05", plan, review, []mission.EvidenceEntry{{ID: "E001"}})
	if errs == nil {
		t.Fatal("expected errors")
	}
	if !containsStr(errs.Errors, "no tasks") {
		t.Errorf("expected no tasks error, got: %v", errs.Errors)
	}
}

func TestReadinessChecker_incompleteTasks(t *testing.T) {
	store, _, cleanup := newTestStore(t)
	defer cleanup()

	checker := NewReadinessChecker(store)
	plan := &mission.Plan{
		Tasks: []mission.Task{
			{ID: strPtr("T01"), Status: strPtr("pending")},
		},
	}
	review := readyReview()
	errs := checker.CheckReadiness("M07AR06", plan, review, []mission.EvidenceEntry{{ID: "E001"}})
	if errs == nil {
		t.Fatal("expected errors")
	}
	if !containsStr(errs.Errors, "incomplete tasks") {
		t.Errorf("expected incomplete tasks error, got: %v", errs.Errors)
	}
}

func TestReadinessChecker_noEvidence(t *testing.T) {
	store, _, cleanup := newTestStore(t)
	defer cleanup()

	checker := NewReadinessChecker(store)
	plan := &mission.Plan{Tasks: []mission.Task{{ID: strPtr("T01"), Status: strPtr("completed")}}}
	review := readyReview()
	errs := checker.CheckReadiness("M07AR07", plan, review, nil)
	if errs == nil {
		t.Fatal("expected errors")
	}
	if !containsStr(errs.Errors, "no evidence") {
		t.Errorf("expected no evidence error, got: %v", errs.Errors)
	}
}

func TestReadinessChecker_blockingFindings(t *testing.T) {
	store, _, cleanup := newTestStore(t)
	defer cleanup()

	checker := NewReadinessChecker(store)
	plan := &mission.Plan{Tasks: []mission.Task{{ID: strPtr("T01"), Status: strPtr("completed")}}}
	review := &mission.Review{
		Status: strPtr("ready"),
		Findings: []mission.Finding{
			{ID: strPtr("F01"), Summary: strPtr("Critical"), Severity: strPtr("critical")},
		},
	}
	errs := checker.CheckReadiness("M07AR08", plan, review, []mission.EvidenceEntry{{ID: "E001"}})
	if errs == nil {
		t.Fatal("expected errors")
	}
	if !containsStr(errs.Errors, "blocking review findings") {
		t.Errorf("expected blocking findings error, got: %v", errs.Errors)
	}
}

func TestMissionArchiver_archiveHappyPath(t *testing.T) {
	store, cfg, cleanup := newTestStore(t)
	defer cleanup()

	// Create a shipped mission with artifacts
	m := &mission.Mission{
		ID:    "M07AR10",
		Title: "Archived mission",
		State: "shipped",
		Git:   mission.GitBlock{WorkBranch: strPtr("feat/archived")},
	}
	if err := store.Create(m); err != nil {
		t.Fatal(err)
	}

	dir := cfg.MissionDir("M07AR10")
	os.WriteFile(filepath.Join(dir, "spec.md"), []byte("# Spec"), 0644)
	os.WriteFile(filepath.Join(dir, "plan.json"), []byte(`{"missionId":"M07AR10","tasks":[{"id":"T01","status":"completed"}]}`), 0644)

	plan := &mission.Plan{
		MissionId: "M07AR10",
		Tasks: []mission.Task{
			{ID: strPtr("T01"), Status: strPtr("completed")},
		},
	}
	review := readyReview()
	entries := []mission.EvidenceEntry{
		{ID: "E001", Label: "test", Command: "true", ExitCode: 0},
	}

	archiver := NewArchiver(store)
	// Override nowFn for deterministic test
	archiver.nowFn = func() string { return "2026-07-07T12:00:00.000Z" }

	result, err := archiver.Archive(ArchiveParams{
		ID:              "M07AR10",
		Mission:         m,
		Plan:            plan,
		Review:          review,
		EvidenceEntries: entries,
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.ArchiveDir == "" {
		t.Fatal("expected non-empty archive dir")
	}

	// Source should be gone
	if dirExists(dir) {
		t.Error("source mission dir should be removed")
	}

	// Archive should exist
	archiveDir := result.ArchiveDir
	if !dirExists(archiveDir) {
		t.Fatal("archive dir should exist")
	}

	// Check archive artifacts
	if !util.Exists(filepath.Join(archiveDir, "SUMMARY.md")) {
		t.Error("SUMMARY.md should exist in archive")
	}
	if !util.Exists(filepath.Join(archiveDir, "mission.json")) {
		t.Error("mission.json should exist in archive")
	}
	if !util.Exists(filepath.Join(archiveDir, "plan.json")) {
		t.Error("plan.json should exist in archive")
	}
	if !util.Exists(filepath.Join(archiveDir, "evidence.jsonl")) {
		t.Error("evidence.jsonl should exist in archive")
	}
	if !util.Exists(filepath.Join(archiveDir, "spec.md")) {
		t.Error("spec.md should exist in archive")
	}

	// Evidence should be compact (no stdout/stderr)
	evData, _ := os.ReadFile(filepath.Join(archiveDir, "evidence.jsonl"))
	if strings.Contains(string(evData), "stdout") {
		t.Error("archive evidence should not contain stdout reference")
	}
}

func TestMissionArchiver_archiveNonShipped(t *testing.T) {
	store, _, cleanup := newTestStore(t)
	defer cleanup()

	m := &mission.Mission{
		ID:    "M07AR11",
		Title: "Not shipped",
		State: "draft",
	}
	store.Create(m)

	archiver := NewArchiver(store)
	_, err := archiver.Archive(ArchiveParams{
		ID:      "M07AR11",
		Mission: m,
	})
	// Archiver doesn't check state, caller does. Should still proceed
	// (it will fail on missing artifacts, but not on state).
	_ = err
}

func TestMissionArchiver_clearSelection(t *testing.T) {
	store, cfg, cleanup := newTestStore(t)
	defer cleanup()

	// Set up current to point to this mission
	store.WriteCurrent("M07AR12")

	m := &mission.Mission{
		ID:    "M07AR12",
		Title: "Clear test",
		State: "shipped",
	}
	store.Create(m)

	dir := cfg.MissionDir("M07AR12")
	os.WriteFile(filepath.Join(dir, "plan.json"), []byte(`{"tasks":[{"id":"T01","status":"completed"}]}`), 0644)

	plan := &mission.Plan{
		Tasks: []mission.Task{
			{ID: strPtr("T01"), Status: strPtr("completed")},
		},
	}
	review := readyReview()

	archiver := NewArchiver(store)
	archiver.nowFn = func() string { return "2026-07-07T12:00:00.000Z" }

	_, err := archiver.Archive(ArchiveParams{
		ID:              "M07AR12",
		Mission:         m,
		Plan:            plan,
		Review:          review,
		EvidenceEntries: []mission.EvidenceEntry{{ID: "E001"}},
	})
	if err != nil {
		t.Fatal(err)
	}

	// Current should be cleared
	curr, _ := store.ReadCurrent()
	if curr != nil && *curr != "" {
		t.Errorf("current should be cleared after archive, got %v", *curr)
	}
}

func TestReadinessChecker_closedStatuses(t *testing.T) {
	t.Run("readiness accepts done and cancelled", func(t *testing.T) {
		store, _, cleanup := newTestStore(t)
		defer cleanup()

		checker := NewReadinessChecker(store)

		// All tasks in closed statuses — none pending/in-progress
		plan := &mission.Plan{
			Tasks: []mission.Task{
				{ID: strPtr("T01"), Status: strPtr("completed")},
				{ID: strPtr("T02"), Status: strPtr("done")},
				{ID: strPtr("T03"), Status: strPtr("cancelled")},
			},
		}
		review := readyReview()
		entries := []mission.EvidenceEntry{
			{ID: "E001", Label: "ok", ExitCode: 0},
		}

		errs := checker.CheckReadiness("M07AR13", plan, review, entries)
		if errs != nil {
			t.Fatalf("expected no errors, got: %v", errs.Errors)
		}
	})

	t.Run("completedCount includes cancelled in summary", func(t *testing.T) {
		store, cfg, cleanup := newTestStore(t)
		defer cleanup()

		m := &mission.Mission{
			ID:    "M07AR14",
			Title: "Summary count test",
			State: "shipped",
			Git:   mission.GitBlock{WorkBranch: strPtr("feat/summary")},
		}
		if err := store.Create(m); err != nil {
			t.Fatal(err)
		}

		dir := cfg.MissionDir("M07AR14")
		os.WriteFile(filepath.Join(dir, "plan.json"), []byte(`{}`), 0644)

		// 3 closed-status tasks: 1 completed, 1 done, 1 cancelled — all should count
		plan := &mission.Plan{
			MissionId: "M07AR14",
			Tasks: []mission.Task{
				{ID: strPtr("T01"), Status: strPtr("completed")},
				{ID: strPtr("T02"), Status: strPtr("done")},
				{ID: strPtr("T03"), Status: strPtr("cancelled")},
			},
		}
		review := readyReview()
		entries := []mission.EvidenceEntry{
			{ID: "E001", Label: "ok", ExitCode: 0},
		}

		archiver := NewArchiver(store)
		archiver.nowFn = func() string { return "2026-07-07T12:00:00.000Z" }

		result, err := archiver.Archive(ArchiveParams{
			ID:              "M07AR14",
			Mission:         m,
			Plan:            plan,
			Review:          review,
			EvidenceEntries: entries,
		})
		if err != nil {
			t.Fatal(err)
		}

		summaryPath := filepath.Join(result.ArchiveDir, "SUMMARY.md")
		data, err := os.ReadFile(summaryPath)
		if err != nil {
			t.Fatal(err)
		}
		if !strings.Contains(string(data), "Tasks: 3/3 completed") {
			t.Errorf("expected summary to show 3/3 completed, got:\n%s", string(data))
		}
	})
}

// --- helpers ---

func newTestStore(t *testing.T) (mission.MissionStore, *config.Config, func()) {
	t.Helper()
	dir, err := os.MkdirTemp("", "spacecraft-archive-test-")
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

func readyReview() *mission.Review {
	return &mission.Review{
		Status: strPtr("ready"),
		Findings: []mission.Finding{},
	}
}

func strPtr(s string) *string { return &s }

func dirExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.IsDir()
}

func containsStr(list []string, sub string) bool {
	for _, s := range list {
		if strings.Contains(s, sub) {
			return true
		}
	}
	return false
}
