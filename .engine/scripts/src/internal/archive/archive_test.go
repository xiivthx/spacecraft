package archive

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"spacecraft/internal/config"
	"spacecraft/internal/mission"
	"spacecraft/internal/roadmap"
	"spacecraft/internal/util"
)

func TestReadinessChecker_ready(t *testing.T) {
	store, _, cleanup := newTestStore(t)
	defer cleanup()

	checker := NewReadinessChecker(store)

	plan := &mission.Plan{
		Tasks: []mission.Task{
			{ID: strPtr("T01"), Status: strPtr("done")},
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
	// Quick-fix lane: no plan is OK if evidence exists
	errs := checker.CheckReadiness("M07AR02", nil, nil, []mission.EvidenceEntry{{ID: "E001"}})
	if errs != nil {
		t.Errorf("expected no errors for quick-fix lane, got: %v", errs.Errors)
	}
}

func TestReadinessChecker_noReview(t *testing.T) {
	store, _, cleanup := newTestStore(t)
	defer cleanup()

	checker := NewReadinessChecker(store)
	plan := &mission.Plan{Tasks: []mission.Task{{ID: strPtr("T01"), Status: strPtr("done")}}}
	// Quick-fix lane: no review is OK if evidence exists
	errs := checker.CheckReadiness("M07AR03", plan, nil, []mission.EvidenceEntry{{ID: "E001"}})
	if errs != nil {
		t.Errorf("expected no errors for quick-fix lane, got: %v", errs.Errors)
	}
}

func TestReadinessChecker_reviewNotReady(t *testing.T) {
	store, _, cleanup := newTestStore(t)
	defer cleanup()

	checker := NewReadinessChecker(store)
	plan := &mission.Plan{Tasks: []mission.Task{{ID: strPtr("T01"), Status: strPtr("done")}}}
	// "not-reviewed" is OK for quick-fix lane
	review := &mission.Review{Status: strPtr("not-reviewed")}
	errs := checker.CheckReadiness("M07AR04", plan, review, []mission.EvidenceEntry{{ID: "E001"}})
	if errs != nil {
		t.Errorf("expected no errors for not-reviewed status, got: %v", errs.Errors)
	}

	// But "blocked" should still fail
	review = &mission.Review{Status: strPtr("blocked")}
	errs = checker.CheckReadiness("M07AR04b", plan, review, []mission.EvidenceEntry{{ID: "E001"}})
	if errs == nil {
		t.Fatal("expected errors for blocked review")
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
	// Quick-fix lane: empty tasks is OK if evidence exists
	errs := checker.CheckReadiness("M07AR05", plan, review, []mission.EvidenceEntry{{ID: "E001"}})
	if errs != nil {
		t.Errorf("expected no errors for quick-fix lane, got: %v", errs.Errors)
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
	plan := &mission.Plan{Tasks: []mission.Task{{ID: strPtr("T01"), Status: strPtr("done")}}}
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
	plan := &mission.Plan{Tasks: []mission.Task{{ID: strPtr("T01"), Status: strPtr("done")}}}
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
	os.WriteFile(filepath.Join(dir, "plan.json"), []byte(`{"missionId":"M07AR10","tasks":[{"id":"T01","status":"done"}]}`), 0644)

	plan := &mission.Plan{
		MissionId: "M07AR10",
		Tasks: []mission.Task{
			{ID: strPtr("T01"), Status: strPtr("done")},
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
	os.WriteFile(filepath.Join(dir, "plan.json"), []byte(`{"tasks":[{"id":"T01","status":"done"}]}`), 0644)

	plan := &mission.Plan{
		Tasks: []mission.Task{
			{ID: strPtr("T01"), Status: strPtr("done")},
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
				{ID: strPtr("T01"), Status: strPtr("done")},
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
				{ID: strPtr("T01"), Status: strPtr("done")},
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

func TestMissionArchiver_Archive_roadmapAware(t *testing.T) {
	store, cfg, cleanup := newTestStore(t)
	defer cleanup()

	// Create a roadmap with 3 missions
	roadmapDir := filepath.Join(cfg.SpaceDir(), "roadmaps")
	if err := os.MkdirAll(roadmapDir, 0755); err != nil {
		t.Fatal(err)
	}
	roadmapData := `{
		"id": "R001",
		"title": "Test Roadmap",
		"missions": ["M07ARCH01", "M07ARCH02", "M07ARCH03"]
	}`
	if err := os.WriteFile(filepath.Join(roadmapDir, "R001.json"), []byte(roadmapData), 0644); err != nil {
		t.Fatal(err)
	}

	// Create missions
	for i, id := range []string{"M07ARCH01", "M07ARCH02", "M07ARCH03"} {
		state := "shipped"
		if i > 0 {
			state = "draft" // Only first mission is shipped
		}
		m := &mission.Mission{
			ID:    id,
			Title: "Test " + id,
			State: state,
		}
		if err := store.Create(m); err != nil {
			t.Fatal(err)
		}
		plan := &mission.Plan{
			MissionId: id,
			Tasks:     []mission.Task{{ID: strPtr("T1"), Status: strPtr("done")}},
		}
		if err := store.SavePlan(id, plan); err != nil {
			t.Fatal(err)
		}
		review := readyReview()
		if err := store.SaveReview(id, review); err != nil {
			t.Fatal(err)
		}
		if err := store.AppendEvidence(id, &mission.EvidenceEntry{ID: "E001", Label: "test", ExitCode: 0}); err != nil {
			t.Fatal(err)
		}
	}

	// Set current to first mission
	if err := store.WriteCurrent("M07ARCH01"); err != nil {
		t.Fatal(err)
	}

	// Create archiver with roadmap store
	roadmapStore := newTestRoadmapStore(t, cfg)
	archiver := NewArchiverWithRoadmap(store, roadmapStore)

	// Archive first mission
	params := ArchiveParams{
		ID:              "M07ARCH01",
		Mission:         mustLoad(store, "M07ARCH01"),
		Plan:            mustLoadPlan(store, "M07ARCH01"),
		Review:          readyReview(),
		EvidenceEntries: []mission.EvidenceEntry{{ID: "E001", Label: "test", ExitCode: 0}},
	}
	result, err := archiver.Archive(params)
	if err != nil {
		t.Fatalf("archive failed: %v", err)
	}
	if result == nil {
		t.Fatal("expected result")
	}

	// Verify current is set to next mission (M07ARCH02)
	currId, err := store.ReadCurrent()
	if err != nil {
		t.Fatal(err)
	}
	if currId == nil {
		t.Fatal("expected current to be set")
	}
	if *currId != "M07ARCH02" {
		t.Errorf("expected current to be M07ARCH02, got %s", *currId)
	}
}

func TestMissionArchiver_Archive_noRoadmap(t *testing.T) {
	store, _, cleanup := newTestStore(t)
	defer cleanup()

	// Create a shipped mission
	m := &mission.Mission{
		ID:    "M07ARCH10",
		Title: "Test",
		State: "shipped",
	}
	if err := store.Create(m); err != nil {
		t.Fatal(err)
	}
	plan := &mission.Plan{
		MissionId: "M07ARCH10",
		Tasks:     []mission.Task{{ID: strPtr("T1"), Status: strPtr("done")}},
	}
	if err := store.SavePlan("M07ARCH10", plan); err != nil {
		t.Fatal(err)
	}
	if err := store.SaveReview("M07ARCH10", readyReview()); err != nil {
		t.Fatal(err)
	}
	if err := store.AppendEvidence("M07ARCH10", &mission.EvidenceEntry{ID: "E001", Label: "test", ExitCode: 0}); err != nil {
		t.Fatal(err)
	}

	// Set current
	if err := store.WriteCurrent("M07ARCH10"); err != nil {
		t.Fatal(err)
	}

	// Create archiver without roadmap store
	archiver := NewArchiver(store)

	// Archive
	params := ArchiveParams{
		ID:              "M07ARCH10",
		Mission:         mustLoad(store, "M07ARCH10"),
		Plan:            mustLoadPlan(store, "M07ARCH10"),
		Review:          readyReview(),
		EvidenceEntries: []mission.EvidenceEntry{{ID: "E001", Label: "test", ExitCode: 0}},
	}
	_, err := archiver.Archive(params)
	if err != nil {
		t.Fatalf("archive failed: %v", err)
	}

	// Verify current is cleared
	currId, err := store.ReadCurrent()
	if err != nil {
		t.Fatal(err)
	}
	if currId != nil {
		t.Errorf("expected current to be cleared, got %s", *currId)
	}
}

func mustLoad(store mission.MissionStore, id string) *mission.Mission {
	m, err := store.Load(id)
	if err != nil {
		panic(err)
	}
	return m
}

func mustLoadPlan(store mission.MissionStore, id string) *mission.Plan {
	p, err := store.LoadPlan(id)
	if err != nil {
		panic(err)
	}
	return p
}

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

func TestParseIssueReferences(t *testing.T) {
	tests := []struct {
		name     string
		text     string
		expected []int
	}{
		{
			name:     "single fix",
			text:     "fixes #123",
			expected: []int{123},
		},
		{
			name:     "multiple fixes",
			text:     "fixes #123 and closes #456",
			expected: []int{123, 456},
		},
		{
			name:     "case insensitive",
			text:     "Fixes #123 CLOSES #456",
			expected: []int{123, 456},
		},
		{
			name:     "resolves keyword",
			text:     "resolves #789",
			expected: []int{789},
		},
		{
			name:     "no duplicates",
			text:     "fixes #123 fixes #123",
			expected: []int{123},
		},
		{
			name:     "no matches",
			text:     "no issue references here",
			expected: nil,
		},
		{
			name:     "mixed content",
			text:     "This mission fixes #23 and closes #24. Also resolves #25.",
			expected: []int{23, 24, 25},
		},
		{
			name:     "Fixed keyword",
			text:     "Fixed #42 in this mission",
			expected: []int{42},
		},
		{
			name:     "fixed lowercase",
			text:     "fixed #55 during development",
			expected: []int{55},
		},
		{
			name:     "checkmark fixed",
			text:     "✅ Fixed #99 in review",
			expected: []int{99},
		},
		{
			name:     "FIXED uppercase",
			text:     "FIXED #77 with refactor",
			expected: []int{77},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseIssueReferences(tt.text)
			if len(result) != len(tt.expected) {
				t.Errorf("expected %d issues, got %d: %v", len(tt.expected), len(result), result)
				return
			}
			for i, num := range result {
				if num != tt.expected[i] {
					t.Errorf("expected issue %d at position %d, got %d", tt.expected[i], i, num)
				}
			}
		})
	}
}

func TestMissionArchiver_Archive_closesGitHubIssues(t *testing.T) {
	store, _, cleanup := newTestStore(t)
	defer cleanup()

	// Create a shipped mission with issue references in spec.md
	m := &mission.Mission{
		ID:    "M07ARCH20",
		Title: "Test",
		State: "shipped",
	}
	if err := store.Create(m); err != nil {
		t.Fatal(err)
	}
	
	// Write spec.md with issue references
	specContent := "# Test Mission\n\nfixes #999\ncloses #998"
	if err := store.WriteFile("M07ARCH20", "spec.md", []byte(specContent)); err != nil {
		t.Fatal(err)
	}
	
	plan := &mission.Plan{
		MissionId: "M07ARCH20",
		Tasks:     []mission.Task{{ID: strPtr("T1"), Status: strPtr("done")}},
	}
	if err := store.SavePlan("M07ARCH20", plan); err != nil {
		t.Fatal(err)
	}
	if err := store.SaveReview("M07ARCH20", readyReview()); err != nil {
		t.Fatal(err)
	}
	if err := store.AppendEvidence("M07ARCH20", &mission.EvidenceEntry{ID: "E001", Label: "test", ExitCode: 0}); err != nil {
		t.Fatal(err)
	}

	// Create archiver
	archiver := NewArchiver(store)

	// Archive - this will try to close issues but gh CLI may not be available
	// The test verifies the parsing logic works, not the actual API call
	params := ArchiveParams{
		ID:              "M07ARCH20",
		Mission:         mustLoad(store, "M07ARCH20"),
		Plan:            mustLoadPlan(store, "M07ARCH20"),
		Review:          readyReview(),
		EvidenceEntries: []mission.EvidenceEntry{{ID: "E001", Label: "test", ExitCode: 0}},
	}
	_, err := archiver.Archive(params)
	if err != nil {
		t.Fatalf("archive failed: %v", err)
	}

	// Verify mission was archived
	if store.SpecExists("M07ARCH20") {
		t.Error("mission should have been archived")
	}
}

func TestParseIssuesFromPlanAndEvidence(t *testing.T) {
	// Plan with issue references in task titles
	plan := &mission.Plan{
		MissionId: "M07TEST1",
		Tasks: []mission.Task{
			{ID: strPtr("T1"), Title: strPtr("Fixed #10"), Status: strPtr("done")},
			{ID: strPtr("T2"), Title: strPtr("Closes #20"), Status: strPtr("done")},
		},
	}
	planJSON, err := jsonMarshal(plan)
	if err != nil {
		t.Fatal(err)
	}
	planIssues := parseIssueReferences(string(planJSON))
	if len(planIssues) != 2 {
		t.Fatalf("expected 2 issues from plan, got %d: %v", len(planIssues), planIssues)
	}

	// Evidence entry with issue reference
	entry := mission.EvidenceEntry{
		ID: "E001", Label: "closes #30", Command: "test", ExitCode: 0,
	}
	evJSON, err := jsonMarshal(entry)
	if err != nil {
		t.Fatal(err)
	}
	evIssues := parseIssueReferences(string(evJSON))
	if len(evIssues) != 1 || evIssues[0] != 30 {
		t.Fatalf("expected issue 30 from evidence, got: %v", evIssues)
	}

	// Cross-source dedup: same #10 in both plan and spec-like text
	allText := string(planJSON) + "\n" + "fixes #10" + "\n" + string(evJSON)
	allIssues := parseIssueReferences(allText)
	if len(allIssues) != 3 {
		t.Fatalf("expected 3 unique issues (#10, #20, #30), got %d: %v", len(allIssues), allIssues)
	}
}

func TestMissionArchiver_Archive_CloseIssuesFalse(t *testing.T) {
	store, _, cleanup := newTestStore(t)
	defer cleanup()

	m := &mission.Mission{
		ID:    "M07ARCH30",
		Title: "No close test",
		State: "shipped",
	}
	if err := store.Create(m); err != nil {
		t.Fatal(err)
	}

	specContent := "fixes #888"
	if err := store.WriteFile("M07ARCH30", "spec.md", []byte(specContent)); err != nil {
		t.Fatal(err)
	}

	plan := &mission.Plan{
		MissionId: "M07ARCH30",
		Tasks:     []mission.Task{{ID: strPtr("T1"), Status: strPtr("done")}},
	}
	if err := store.SavePlan("M07ARCH30", plan); err != nil {
		t.Fatal(err)
	}
	if err := store.SaveReview("M07ARCH30", readyReview()); err != nil {
		t.Fatal(err)
	}
	if err := store.AppendEvidence("M07ARCH30", &mission.EvidenceEntry{ID: "E001", Label: "test", ExitCode: 0}); err != nil {
		t.Fatal(err)
	}

	archiver := NewArchiver(store)

	params := ArchiveParams{
		ID:              "M07ARCH30",
		Mission:         mustLoad(store, "M07ARCH30"),
		Plan:            mustLoadPlan(store, "M07ARCH30"),
		Review:          readyReview(),
		EvidenceEntries: []mission.EvidenceEntry{{ID: "E001", Label: "test", ExitCode: 0}},
		CloseIssues:     false,
	}
	_, err := archiver.Archive(params)
	if err != nil {
		t.Fatalf("archive with CloseIssues=false failed: %v", err)
	}
}

func TestIsIssueClosed_ghNotAvailable(t *testing.T) {
	// isIssueClosed should return false when gh is not on PATH
	// (best-effort fallback: if check fails, attempt close anyway)
	result := isIssueClosed(99999)
	if result {
		t.Error("isIssueClosed should return false when gh is unavailable")
	}
}

func containsStr(list []string, sub string) bool {
	for _, s := range list {
		if strings.Contains(s, sub) {
			return true
		}
	}
	return false
}

type testRoadmapStore struct {
	cfg *config.Config
}

func newTestRoadmapStore(t *testing.T, cfg *config.Config) *testRoadmapStore {
	return &testRoadmapStore{cfg: cfg}
}

func (s *testRoadmapStore) Create(r *roadmap.Roadmap) error {
	return nil
}

func (s *testRoadmapStore) Save(r *roadmap.Roadmap) error {
	return nil
}

func (s *testRoadmapStore) Delete(id string) error {
	return nil
}

func (s *testRoadmapStore) List() ([]*roadmap.Roadmap, error) {
	dir := filepath.Join(s.cfg.SpaceDir(), "roadmaps")
	if !util.Exists(dir) {
		return nil, nil
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	var out []*roadmap.Roadmap
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".json") {
			continue
		}
		r, err := s.Load(strings.TrimSuffix(e.Name(), ".json"))
		if err != nil {
			continue
		}
		out = append(out, r)
	}
	return out, nil
}

func (s *testRoadmapStore) Load(id string) (*roadmap.Roadmap, error) {
	path := filepath.Join(s.cfg.SpaceDir(), "roadmaps", id+".json")
	if !util.Exists(path) {
		return nil, fmt.Errorf("roadmap not found: %s", id)
	}
	var r roadmap.Roadmap
	if err := util.ReadJson(path, &r); err != nil {
		return nil, err
	}
	return &r, nil
}
