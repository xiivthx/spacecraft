package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"spacecraft/internal/mission"
	"spacecraft/internal/roadmap"
	"spacecraft/internal/trace"
)

// withRegistryMock swaps the research/registry endpoints to a local httptest server.
func withRegistryMock(t *testing.T) func() {
	t.Helper()
	oldBrave := braveBaseURL
	oldNpm := npmRegistryURL
	oldGo := goProxyURL
	oldPypi := pypiURL
	oldCrates := cratesURL

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		p := r.URL.Path

		// Explicit failure path used by tests that need all registries to miss.
		if strings.Contains(p, "/fail") {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		// Brave Search.
		if p == "/res/v1/web/search" {
			q := r.URL.Query().Get("q")
			w.Header().Set("Content-Type", "application/json")
			if strings.Contains(q, "no-results") {
				json.NewEncoder(w).Encode(map[string]interface{}{"web": map[string]interface{}{"results": []interface{}{}}})
				return
			}
			json.NewEncoder(w).Encode(map[string]interface{}{
				"web": map[string]interface{}{
					"results": []map[string]string{
						{"title": "Example", "url": "http://example.com", "description": "An example result"},
					},
				},
			})
			return
		}

		// Go proxy.
		if strings.HasSuffix(p, "/@latest") {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{
				"Version": "v2.0.0",
				"Time":    time.Now().UTC().Format(time.RFC3339),
			})
			return
		}

		// PyPI.
		if strings.HasPrefix(p, "/pypi/") && strings.HasSuffix(p, "/json") {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{
				"info": map[string]string{
					"name":      "test-pkg",
					"version":   "2.0.0",
					"license":   "MIT",
					"home_page": "http://example.com",
				},
			})
			return
		}

		// Crates.
		if strings.HasPrefix(p, "/api/v1/crates/") {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{
				"crate": map[string]string{
					"name":               "test-pkg",
					"max_stable_version": "2.0.0",
					"license":            "MIT",
					"homepage":           "http://example.com",
					"updated_at":         time.Now().UTC().Format(time.RFC3339),
				},
			})
			return
		}

		// npm: fail paths that look like Go module URLs so the go client can match.
		if strings.Contains(p, ".com") {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		// npm package.
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"name":     "test-pkg",
			"version":  "3.0.0",
			"license":  "MIT",
			"homepage": "http://example.com",
		})
	}))

	braveBaseURL = srv.URL
	npmRegistryURL = srv.URL
	goProxyURL = srv.URL
	pypiURL = srv.URL
	cratesURL = srv.URL

	return func() {
		srv.Close()
		braveBaseURL = oldBrave
		npmRegistryURL = oldNpm
		goProxyURL = oldGo
		pypiURL = oldPypi
		cratesURL = oldCrates
	}
}

// withFakeExecutable prepends a temp directory containing an executable named name to PATH.
func withFakeExecutable(t *testing.T, name, script string) func() {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, name)
	if err := os.WriteFile(path, []byte("#!/bin/sh\n"+script), 0755); err != nil {
		t.Fatalf("write fake %s: %v", name, err)
	}
	oldPATH := os.Getenv("PATH")
	os.Setenv("PATH", dir+string(filepath.ListSeparator)+oldPATH)
	return func() { os.Setenv("PATH", oldPATH) }
}

// writeTraceEntry writes a minimal trace event file for a mission.
func writeTraceFile(t *testing.T, missionID string, entries ...trace.TraceEntry) {
	t.Helper()
	dir := cfg.TraceStoreDir()
	if err := os.MkdirAll(dir, 0755); err != nil {
		t.Fatalf("mkdir traces: %v", err)
	}
	f, err := os.Create(filepath.Join(dir, missionID+".jsonl"))
	if err != nil {
		t.Fatalf("create trace file: %v", err)
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	for _, e := range entries {
		if e.TS == "" {
			e.TS = time.Now().UTC().Format(time.RFC3339Nano)
		}
		if e.ID == "" {
			e.ID = "trace-" + e.TS
		}
		if err := enc.Encode(e); err != nil {
			t.Fatalf("encode trace entry: %v", err)
		}
	}
}

func TestNormalizeDeepArgs(t *testing.T) {
	cases := []struct {
		in   []string
		want []string
	}{
		{[]string{"--deep"}, []string{"--deep=true"}},
		{[]string{"--deep", "--json", "query"}, []string{"--deep=true", "--json", "query"}},
		{[]string{"--deep", "nlm", "query"}, []string{"--deep", "nlm", "query"}},
		{[]string{"--deep=true"}, []string{"--deep=true"}},
		{[]string{"query"}, []string{"query"}},
	}
	for _, c := range cases {
		got := normalizeDeepArgs(c.in)
		if strings.Join(got, ",") != strings.Join(c.want, ",") {
			t.Errorf("normalizeDeepArgs(%v) = %v, want %v", c.in, got, c.want)
		}
	}
}

func TestResearchCmdPackageQuery(t *testing.T) {
	defer setupLifecycleTest(t)()
	defer withRegistryMock(t)()

	exit := researchCmd([]string{"example.com/foo"})
	if exit != 0 {
		t.Errorf("researchCmd(package query) = %d, want 0", exit)
	}
}

func TestResearchCmdPackageQueryNpm(t *testing.T) {
	defer setupLifecycleTest(t)()
	defer withRegistryMock(t)()

	exit := researchCmd([]string{"@scope/pkg"})
	if exit != 0 {
		t.Errorf("researchCmd(npm package query) = %d, want 0", exit)
	}
}

func TestResearchCmdBraveSearch(t *testing.T) {
	defer setupLifecycleTest(t)()
	defer withRegistryMock(t)()
	t.Setenv("SPACECRAFT_BRAVE_API_KEY", "test-key")

	exit := researchCmd([]string{"test query"})
	if exit != 0 {
		t.Errorf("researchCmd(brave search) = %d, want 0", exit)
	}
}

func TestResearchCmdBraveSearchNoResults(t *testing.T) {
	defer setupLifecycleTest(t)()
	defer withRegistryMock(t)()
	t.Setenv("SPACECRAFT_BRAVE_API_KEY", "test-key")

	exit := researchCmd([]string{"no-results query"})
	if exit != 1 {
		t.Errorf("researchCmd(no results) = %d, want 1", exit)
	}
}

func TestResearchCmdNoSave(t *testing.T) {
	defer setupLifecycleTest(t)()
	defer withRegistryMock(t)()
	t.Setenv("SPACECRAFT_BRAVE_API_KEY", "test-key")

	exit := researchCmd([]string{"--no-save", "test query"})
	if exit != 0 {
		t.Errorf("researchCmd(--no-save) = %d, want 0", exit)
	}
}

func TestResearchCmdUnknownScope(t *testing.T) {
	defer setupLifecycleTest(t)()
	defer withRegistryMock(t)()

	exit := researchCmd([]string{"--scope", "unknown-scope", "test query"})
	if exit != 1 {
		t.Errorf("researchCmd(unknown scope) = %d, want 1", exit)
	}
}

func TestRunDeepAnalysisInvalidMode(t *testing.T) {
	ctx := t.Context()
	_, err := runDeepAnalysis(ctx, "http://example.com", "query", "bad", 1*time.Second)
	if err == nil {
		t.Error("runDeepAnalysis(bad mode) expected error")
	}
}

func TestRunDeepAnalysisBrowserUse(t *testing.T) {
	defer withFakeExecutable(t, "python3", `printf '{"summary":"s","key_points":["a"],"source_url":"http://example.com"}'`)()
	ctx := t.Context()
	res, err := runDeepAnalysis(ctx, "http://example.com", "query", "true", 1*time.Second)
	if err != nil {
		t.Fatalf("runDeepAnalysis(true) = %v", err)
	}
	if res.Summary == "" {
		t.Error("runDeepAnalysis(true) returned empty summary")
	}
}

func TestRunDeepAnalysisNotebookLM(t *testing.T) {
	defer withFakeExecutable(t, "nlm", `if [ "$2" = "create" ]; then echo "nb123"; else echo "synthesis"; fi`)()
	ctx := t.Context()
	res, err := runDeepAnalysis(ctx, "http://example.com", "query", "nlm", 1*time.Second)
	if err != nil {
		t.Fatalf("runDeepAnalysis(nlm) = %v", err)
	}
	if res.Summary == "" {
		t.Error("runDeepAnalysis(nlm) returned empty summary")
	}
}

func TestLookupPackage(t *testing.T) {
	defer withRegistryMock(t)()
	ctx := t.Context()

	pkg, src := lookupPackage(ctx, "example.com/foo", 1*time.Second)
	if pkg == nil || src != "go" {
		t.Errorf("lookupPackage(go) = (%v, %q), want go package", pkg, src)
	}

	pkg, src = lookupPackage(ctx, "@scope/pkg", 1*time.Second)
	if pkg == nil || src != "npm" {
		t.Errorf("lookupPackage(npm) = (%v, %q), want npm package", pkg, src)
	}

	pkg, src = lookupPackage(ctx, "fail", 1*time.Second)
	if pkg != nil {
		t.Errorf("lookupPackage(all miss) = %v, want nil", pkg)
	}
}

func TestLookupLatestVersion(t *testing.T) {
	defer withRegistryMock(t)()
	ctx := t.Context()

	latest, err := lookupLatestVersion(ctx, researchDependency("go", "example.com/foo", "v1.0.0"), 1*time.Second)
	if err != nil || latest != "v2.0.0" {
		t.Errorf("lookupLatestVersion(go) = (%q, %v), want v2.0.0", latest, err)
	}

	_, err = lookupLatestVersion(ctx, researchDependency("unknown", "x", "1.0.0"), 1*time.Second)
	if err == nil {
		t.Error("lookupLatestVersion(unknown ecosystem) expected error")
	}
}

func researchDependency(eco, name, ver string) struct {
	Name           string
	CurrentVersion string
	Ecosystem      string
} {
	return struct {
		Name           string
		CurrentVersion string
		Ecosystem      string
	}{Name: name, CurrentVersion: ver, Ecosystem: eco}
}

func TestExtractNumericVersionAndSplitVersion(t *testing.T) {
	cases := []struct {
		in      string
		wantVer string
		wantParts []int
	}{
		{"v1.2.3", "1.2.3", []int{1, 2, 3}},
		{"^1.2.3-beta", "1.2.3", []int{1, 2, 3}},
		{">=1.2", "1.2", []int{1, 2}},
		{"  1.2.3.4  ", "1.2.3.4", []int{1, 2, 3, 4}},
	}
	for _, c := range cases {
		got := extractNumericVersion(c.in)
		if got != c.wantVer {
			t.Errorf("extractNumericVersion(%q) = %q, want %q", c.in, got, c.wantVer)
		}
		parts := splitVersion(got)
		if len(parts) != len(c.wantParts) {
			t.Errorf("splitVersion(%q) = %v, want %v", got, parts, c.wantParts)
		}
	}
}

func TestEvalCmd(t *testing.T) {
	defer setupLifecycleTest(t)()
	id := createMission(t, "Eval Mission")

	if e := evalCmd([]string{}); e != 1 {
		t.Errorf("evalCmd(empty) = %d, want 1", e)
	}
	if e := evalCmd([]string{"--help"}); e != 0 {
		t.Errorf("evalCmd(--help) = %d, want 0", e)
	}
	if e := evalCmd([]string{"init", id}); e != 0 {
		t.Errorf("evalCmd(init) = %d, want 0", e)
	}
	if e := evalCmd([]string{id}); e != 0 {
		t.Errorf("evalCmd(run) = %d, want 0", e)
	}
}

func TestEvalInitCmdNoMission(t *testing.T) {
	defer setupLifecycleTest(t)()
	if e := evalInitCmd([]string{}); e != 1 {
		t.Errorf("evalInitCmd(empty) = %d, want 1", e)
	}
}

func TestEvalRunCmdNoMission(t *testing.T) {
	defer setupLifecycleTest(t)()
	if e := evalRunCmd([]string{"nonexistent"}); e != 1 {
		t.Errorf("evalRunCmd(nonexistent) = %d, want 1", e)
	}
}

func TestCheckDepsCmdNoManifests(t *testing.T) {
	defer setupLifecycleTest(t)()
	if e := checkDepsCmd([]string{}); e != 0 {
		t.Errorf("checkDepsCmd(no manifests) = %d, want 0", e)
	}
}

func TestCheckDepsCmdGoManifest(t *testing.T) {
	defer setupLifecycleTest(t)()
	defer withRegistryMock(t)()

	goMod := `module test
go 1.22
require example.com/foo v1.0.0
`
	if err := os.WriteFile(filepath.Join(cfg.Root(), "go.mod"), []byte(goMod), 0644); err != nil {
		t.Fatalf("write go.mod: %v", err)
	}

	if e := checkDepsCmd([]string{"--registry", "go", "--timeout", "1s"}); e != 0 {
		t.Errorf("checkDepsCmd(go) = %d, want 0", e)
	}
}

func TestCheckDepsCmdJSON(t *testing.T) {
	defer setupLifecycleTest(t)()
	defer withRegistryMock(t)()

	goMod := `module test
go 1.22
require example.com/foo v1.0.0
`
	if err := os.WriteFile(filepath.Join(cfg.Root(), "go.mod"), []byte(goMod), 0644); err != nil {
		t.Fatalf("write go.mod: %v", err)
	}

	if e := checkDepsCmd([]string{"--registry", "go", "--timeout", "1s", "--json"}); e != 0 {
		t.Errorf("checkDepsCmd(json) = %d, want 0", e)
	}
}

func TestTracesCmd(t *testing.T) {
	defer setupLifecycleTest(t)()
	id := createMission(t, "Traces Mission")

	if e := tracesCmd([]string{"--help"}); e != 0 {
		t.Errorf("tracesCmd(--help) = %d, want 0", e)
	}
	if e := tracesCmd([]string{}); e != 1 {
		t.Errorf("tracesCmd(empty) = %d, want 1", e)
	}
	if e := tracesCmd([]string{id}); e != 0 {
		t.Errorf("tracesCmd(no traces) = %d, want 0", e)
	}

	model := "test-model"
	writeTraceFile(t, id, trace.TraceEntry{
		Seq:          1,
		Type:         trace.EventModelInvoke,
		Model:        &model,
		InputTokens:  100,
		OutputTokens: 50,
		LatencyMs:    10,
	})

	if e := tracesCmd([]string{"--json", id}); e != 0 {
		t.Errorf("tracesCmd(--json) = %d, want 0", e)
	}
	if e := tracesCmd([]string{"--verbose", "--flat", id}); e != 0 {
		t.Errorf("tracesCmd(verbose flat) = %d, want 0", e)
	}
}

func TestTracesCmdInvalidMission(t *testing.T) {
	defer setupLifecycleTest(t)()
	if e := tracesCmd([]string{"bad-id"}); e != 1 {
		t.Errorf("tracesCmd(bad id) = %d, want 1", e)
	}
}

func TestCostCmd(t *testing.T) {
	defer setupLifecycleTest(t)()
	id := createMission(t, "Cost Mission")

	if e := costCmd([]string{"--help"}); e != 0 {
		t.Errorf("costCmd(--help) = %d, want 0", e)
	}
	if e := costCmd([]string{}); e != 0 {
		t.Errorf("costCmd(no traces) = %d, want 0", e)
	}

	model := "test-model"
	writeTraceFile(t, id, trace.TraceEntry{
		Seq:          1,
		Type:         trace.EventModelInvoke,
		Model:        &model,
		InputTokens:  1000,
		OutputTokens: 500,
		LatencyMs:    100,
	})

	if e := costCmd([]string{}); e != 0 {
		t.Errorf("costCmd(with traces) = %d, want 0", e)
	}
	if e := costCmd([]string{"--mission", id}); e != 0 {
		t.Errorf("costCmd(--mission) = %d, want 0", e)
	}
}

func TestCostMissionCmdNoTraces(t *testing.T) {
	defer setupLifecycleTest(t)()
	id := createMission(t, "Cost No Traces")
	if e := costMissionCmd(id); e != 0 {
		t.Errorf("costMissionCmd(no traces) = %d, want 0", e)
	}
}

func TestCostMissionCmdInvalidMission(t *testing.T) {
	defer setupLifecycleTest(t)()
	if e := costMissionCmd("bad-id"); e != 1 {
		t.Errorf("costMissionCmd(bad id) = %d, want 1", e)
	}
}

func TestGitInfoCmd(t *testing.T) {
	defer setupLifecycleTest(t)()
	if e := gitInfoCmd(); e != 0 {
		t.Errorf("gitInfoCmd(no repo) = %d, want 0", e)
	}
}

func TestGitInfoCmdInRepo(t *testing.T) {
	defer setupLifecycleTest(t)()
	initGit(t)
	if e := gitInfoCmd(); e != 0 {
		t.Errorf("gitInfoCmd(repo) = %d, want 0", e)
	}
}

func TestGitSuggestCmd(t *testing.T) {
	defer setupLifecycleTest(t)()
	id := createMission(t, "Suggest Mission")

	if e := gitSuggestCmd([]string{}); e != 0 {
		t.Errorf("gitSuggestCmd() = %d, want 0", e)
	}
	if e := gitSuggestCmd([]string{"fix", "bug thing"}); e != 0 {
		t.Errorf("gitSuggestCmd(fix bug thing) = %d, want 0", e)
	}
	if e := gitSuggestCmd([]string{id}); e != 0 {
		t.Errorf("gitSuggestCmd(id slug) = %d, want 0", e)
	}
}

func TestWorkflowCmd(t *testing.T) {
	defer setupLifecycleTest(t)()
	createMission(t, "Workflow Mission")

	if e := workflowCmd([]string{}); e != 0 {
		t.Errorf("workflowCmd() = %d, want 0", e)
	}
	if e := workflowCmd([]string{"--json"}); e != 0 {
		t.Errorf("workflowCmd(--json) = %d, want 0", e)
	}

	// No mission should fail.
	restore := setupLifecycleTest(t)
	defer restore()
	if e := workflowCmd([]string{}); e != 1 {
		t.Errorf("workflowCmd(no mission) = %d, want 1", e)
	}
}

func TestRoadmapCmd(t *testing.T) {
	defer setupLifecycleTest(t)()

	if e := roadmapCmd([]string{}); e != 1 {
		t.Errorf("roadmapCmd(empty) = %d, want 1", e)
	}
	if e := roadmapCmd([]string{"bad-sub"}); e != 1 {
		t.Errorf("roadmapCmd(bad sub) = %d, want 1", e)
	}
}

func TestRoadmapNewCmd(t *testing.T) {
	defer setupLifecycleTest(t)()

	if e := roadmapNewCmd([]string{}); e != 1 {
		t.Errorf("roadmapNewCmd(empty) = %d, want 1", e)
	}
	if e := roadmapNewCmd([]string{"My Roadmap"}); e != 0 {
		t.Errorf("roadmapNewCmd(title) = %d, want 0", e)
	}
	time.Sleep(2 * time.Millisecond)
	if e := roadmapNewCmd([]string{"--desc", "A description", "Roadmap Two"}); e != 0 {
		t.Errorf("roadmapNewCmd(desc) = %d, want 0", e)
	}
}

func TestRoadmapListCmd(t *testing.T) {
	defer setupLifecycleTest(t)()

	if e := roadmapListCmd(); e != 0 {
		t.Errorf("roadmapListCmd(empty) = %d, want 0", e)
	}
	if e := roadmapNewCmd([]string{"Listed"}); e != 0 {
		t.Fatalf("roadmapNewCmd = %d", e)
	}
	time.Sleep(2 * time.Millisecond)
	if e := roadmapListCmd(); e != 0 {
		t.Errorf("roadmapListCmd(with roadmap) = %d, want 0", e)
	}
}

func TestRoadmapAddRemoveCmd(t *testing.T) {
	defer setupLifecycleTest(t)()
	mid := createMission(t, "Roadmap Mission")
	if e := roadmapNewCmd([]string{"Roadmap"}); e != 0 {
		t.Fatalf("roadmapNewCmd = %d", e)
	}
	rms, _ := roadmapStore.List()
	if len(rms) == 0 {
		t.Fatal("no roadmaps created")
	}
	rid := rms[0].ID

	if e := roadmapAddCmd([]string{rid, mid}); e != 0 {
		t.Errorf("roadmapAddCmd = %d, want 0", e)
	}
	// Duplicate add is idempotent.
	if e := roadmapAddCmd([]string{rid, mid}); e != 0 {
		t.Errorf("roadmapAddCmd(duplicate) = %d, want 0", e)
	}
	if e := roadmapAddCmd([]string{}); e != 1 {
		t.Errorf("roadmapAddCmd(empty) = %d, want 1", e)
	}

	if e := roadmapRemoveCmd([]string{rid, mid}); e != 0 {
		t.Errorf("roadmapRemoveCmd = %d, want 0", e)
	}
	if e := roadmapRemoveCmd([]string{rid, mid}); e != 1 {
		t.Errorf("roadmapRemoveCmd(again) = %d, want 1", e)
	}
}

func TestRoadmapAddAfterCmd(t *testing.T) {
	defer setupLifecycleTest(t)()
	m1 := createMission(t, "First")
	m2 := createMission(t, "Second")
	m3 := createMission(t, "Third")
	if e := roadmapNewCmd([]string{"Roadmap"}); e != 0 {
		t.Fatalf("roadmapNewCmd = %d", e)
	}
	rms, _ := roadmapStore.List()
	rid := rms[0].ID

	if e := roadmapAddCmd([]string{rid, m1}); e != 0 {
		t.Fatalf("roadmapAddCmd(first) = %d", e)
	}
	if e := roadmapAddCmd([]string{rid, m2, "--after", m1}); e != 0 {
		t.Errorf("roadmapAddCmd(after) = %d, want 0", e)
	}
	if e := roadmapAddCmd([]string{rid, m3, "--after", "missing"}); e != 1 {
		t.Errorf("roadmapAddCmd(after missing) = %d, want 1", e)
	}

	rm, _ := roadmapStore.Load(rid)
	if len(rm.Missions) != 2 || rm.Missions[1].ID != m2 {
		t.Errorf("roadmap order = %v", rm.Missions)
	}
}

func TestRoadmapAddAlreadyInOther(t *testing.T) {
	defer setupLifecycleTest(t)()
	mid := createMission(t, "Shared")
	if e := roadmapNewCmd([]string{"R1"}); e != 0 {
		t.Fatalf("roadmapNewCmd R1 = %d", e)
	}
	time.Sleep(2 * time.Millisecond)
	if e := roadmapNewCmd([]string{"R2"}); e != 0 {
		t.Fatalf("roadmapNewCmd R2 = %d", e)
	}
	rms, _ := roadmapStore.List()
	if len(rms) != 2 {
		t.Fatalf("expected 2 roadmaps, got %d", len(rms))
	}
	r1 := rms[0].ID
	r2 := rms[1].ID

	if e := roadmapAddCmd([]string{r1, mid}); e != 0 {
		t.Fatalf("add to r1 = %d", e)
	}
	if e := roadmapAddCmd([]string{r2, mid}); e != 1 {
		t.Errorf("add to r2 = %d, want 1", e)
	}
}

func TestRoadmapShowCmd(t *testing.T) {
	defer setupLifecycleTest(t)()
	mid := createMission(t, "Show Mission")
	if e := roadmapNewCmd([]string{"Show Roadmap", "--desc", "Desc"}); e != 0 {
		t.Fatalf("roadmapNewCmd = %d", e)
	}
	rms, _ := roadmapStore.List()
	rid := rms[0].ID
	if e := roadmapAddCmd([]string{rid, mid}); e != 0 {
		t.Fatalf("roadmapAddCmd = %d", e)
	}

	// Inject an issue to cover issue display.
	rm, _ := roadmapStore.Load(rid)
	rm.Issues = []roadmap.Issue{{Number: 1, Title: "Issue One", State: "open", Phase: "design", Labels: []string{"bug"}}}
	roadmapStore.Save(rm)

	if e := roadmapShowCmd([]string{}); e != 1 {
		t.Errorf("roadmapShowCmd(empty) = %d, want 1", e)
	}
	if e := roadmapShowCmd([]string{rid}); e != 0 {
		t.Errorf("roadmapShowCmd(rid) = %d, want 0", e)
	}
}

func TestRoadmapAddWithDescription(t *testing.T) {
	defer setupLifecycleTest(t)()
	mid := createMission(t, "Desc Mission")
	if e := roadmapNewCmd([]string{"Desc Roadmap"}); e != 0 {
		t.Fatalf("roadmapNewCmd = %d", e)
	}
	rms, _ := roadmapStore.List()
	rid := rms[0].ID
	if e := roadmapAddCmd([]string{rid, mid, "--desc", "This is a test description"}); e != 0 {
		t.Fatalf("roadmapAddCmd = %d", e)
	}

	rm, _ := roadmapStore.Load(rid)
	if len(rm.Missions) != 1 {
		t.Fatalf("expected 1 mission, got %d", len(rm.Missions))
	}
	if rm.Missions[0].Description != "This is a test description" {
		t.Errorf("description = %q, want %q", rm.Missions[0].Description, "This is a test description")
	}
}

func TestRoadmapContinueCmd(t *testing.T) {
	defer setupLifecycleTest(t)()
	m1 := createMission(t, "Continue One")
	_ = createMission(t, "Continue Two")
	if e := roadmapNewCmd([]string{"Continue Roadmap"}); e != 0 {
		t.Fatalf("roadmapNewCmd = %d", e)
	}
	rms, _ := roadmapStore.List()
	rid := rms[0].ID
	if e := roadmapAddCmd([]string{rid, m1}); e != 0 {
		t.Fatalf("roadmapAddCmd = %d", e)
	}

	if e := roadmapContinueCmd([]string{}); e != 1 {
		t.Errorf("roadmapContinueCmd(empty) = %d, want 1", e)
	}
	if e := roadmapContinueCmd([]string{rid}); e != 0 {
		t.Errorf("roadmapContinueCmd = %d, want 0", e)
	}
	curr, _ := store.ReadCurrent()
	if curr == nil || *curr != m1 {
		t.Errorf("current after continue = %v, want %s", curr, m1)
	}
}

func TestRoadmapArchiveCmd(t *testing.T) {
	defer setupLifecycleTest(t)()
	mid := makeShippableMission(t)
	if code := setStateCmd([]string{"shipped"}); code != 0 {
		t.Fatalf("setStateCmd(shipped) = %d, want 0", code)
	}
	if e := roadmapNewCmd([]string{"Archive Roadmap"}); e != 0 {
		t.Fatalf("roadmapNewCmd = %d", e)
	}
	rms, _ := roadmapStore.List()
	rid := rms[0].ID
	if e := roadmapAddCmd([]string{rid, mid}); e != 0 {
		t.Fatalf("roadmapAddCmd = %d", e)
	}

	if e := roadmapArchiveCmd([]string{}); e != 1 {
		t.Errorf("roadmapArchiveCmd(empty) = %d, want 1", e)
	}
	if e := roadmapArchiveCmd([]string{rid}); e != 0 {
		t.Errorf("roadmapArchiveCmd = %d, want 0", e)
	}
}

func TestRoadmapArchiveBlocked(t *testing.T) {
	defer setupLifecycleTest(t)()
	mid := createMission(t, "Active Roadmap Mission")
	if e := roadmapNewCmd([]string{"Blocked Roadmap"}); e != 0 {
		t.Fatalf("roadmapNewCmd = %d", e)
	}
	rms, _ := roadmapStore.List()
	rid := rms[0].ID
	if e := roadmapAddCmd([]string{rid, mid}); e != 0 {
		t.Fatalf("roadmapAddCmd = %d", e)
	}

	if e := roadmapArchiveCmd([]string{rid}); e != 1 {
		t.Errorf("roadmapArchiveCmd(active mission) = %d, want 1", e)
	}
}

// writeArchiveMission creates a mission.json in the archive directory.
func writeArchiveMission(t *testing.T, id, title, state string) {
	t.Helper()
	dir := filepath.Join(cfg.ArchiveDir(), id)
	if err := os.MkdirAll(dir, 0755); err != nil {
		t.Fatalf("mkdir archive: %v", err)
	}
	m := mission.Mission{
		ID:    id,
		Title: title,
		State: state,
	}
	data, err := json.Marshal(m)
	if err != nil {
		t.Fatalf("marshal mission: %v", err)
	}
	if err := os.WriteFile(filepath.Join(dir, "mission.json"), data, 0644); err != nil {
		t.Fatalf("write archive mission: %v", err)
	}
}

func TestRoadmapShowWithArchivedMission(t *testing.T) {
	defer setupLifecycleTest(t)()
	mid := createMission(t, "Live Shipped")
	if e := setStateCmd([]string{"shipped"}); e != 0 {
		t.Fatalf("setStateCmd = %d", e)
	}

	archID := "MARCHIVED01"
	writeArchiveMission(t, archID, "Archived Mission", "shipped")

	if e := roadmapNewCmd([]string{"Mixed Roadmap"}); e != 0 {
		t.Fatalf("roadmapNewCmd = %d", e)
	}
	rms, _ := roadmapStore.List()
	rid := rms[0].ID

	if e := roadmapAddCmd([]string{rid, mid}); e != 0 {
		t.Fatalf("add live = %d", e)
	}
	if e := roadmapAddCmd([]string{rid, archID}); e != 0 {
		t.Fatalf("add archived = %d", e)
	}

	if e := roadmapShowCmd([]string{rid}); e != 0 {
		t.Errorf("roadmapShowCmd with archived = %d, want 0", e)
	}
}

func TestRoadmapListCountsArchived(t *testing.T) {
	defer setupLifecycleTest(t)()
	mid := createMission(t, "Live Shipped")
	if e := setStateCmd([]string{"shipped"}); e != 0 {
		t.Fatalf("setStateCmd = %d", e)
	}

	archID := "MARCHIVED02"
	writeArchiveMission(t, archID, "Archived Mission", "shipped")

	if e := roadmapNewCmd([]string{"Count Roadmap"}); e != 0 {
		t.Fatalf("roadmapNewCmd = %d", e)
	}
	rms, _ := roadmapStore.List()
	rid := rms[0].ID

	if e := roadmapAddCmd([]string{rid, mid}); e != 0 {
		t.Fatalf("add live = %d", e)
	}
	if e := roadmapAddCmd([]string{rid, archID}); e != 0 {
		t.Fatalf("add archived = %d", e)
	}

	// roadmapListCmd prints [shipped/total] — just verify it doesn't error
	if e := roadmapListCmd(); e != 0 {
		t.Errorf("roadmapListCmd with archived = %d, want 0", e)
	}
}

func TestRoadmapContinueSkipsArchived(t *testing.T) {
	defer setupLifecycleTest(t)()
	archID := "MARCHIVED03"
	writeArchiveMission(t, archID, "Already Archived", "archived")

	mid := createMission(t, "Still Active")

	if e := roadmapNewCmd([]string{"Continue Roadmap"}); e != 0 {
		t.Fatalf("roadmapNewCmd = %d", e)
	}
	rms, _ := roadmapStore.List()
	rid := rms[0].ID

	if e := roadmapAddCmd([]string{rid, archID}); e != 0 {
		t.Fatalf("add archived = %d", e)
	}
	if e := roadmapAddCmd([]string{rid, mid}); e != 0 {
		t.Fatalf("add active = %d", e)
	}

	if e := roadmapContinueCmd([]string{rid}); e != 0 {
		t.Errorf("roadmapContinueCmd = %d, want 0", e)
	}
	curr, _ := store.ReadCurrent()
	if curr == nil || *curr != mid {
		t.Errorf("current = %v, want %s (must skip archived)", curr, mid)
	}
}

func TestRoadmapArchiveAllArchived(t *testing.T) {
	defer setupLifecycleTest(t)()
	archID := "MARCHIVED04"
	writeArchiveMission(t, archID, "Archived Mission", "shipped")

	if e := roadmapNewCmd([]string{"All Archived Roadmap"}); e != 0 {
		t.Fatalf("roadmapNewCmd = %d", e)
	}
	rms, _ := roadmapStore.List()
	rid := rms[0].ID

	if e := roadmapAddCmd([]string{rid, archID}); e != 0 {
		t.Fatalf("roadmapAddCmd = %d", e)
	}

	if e := roadmapArchiveCmd([]string{rid}); e != 0 {
		t.Errorf("roadmapArchiveCmd(all archived) = %d, want 0", e)
	}
}

func TestRoadmapAddArchivedMission(t *testing.T) {
	defer setupLifecycleTest(t)()
	archID := "MARCHIVED05"
	writeArchiveMission(t, archID, "Archived Mission", "shipped")

	if e := roadmapNewCmd([]string{"Add Roadmap"}); e != 0 {
		t.Fatalf("roadmapNewCmd = %d", e)
	}
	rms, _ := roadmapStore.List()
	rid := rms[0].ID

	if e := roadmapAddCmd([]string{rid, archID}); e != 0 {
		t.Errorf("roadmapAddCmd(archived mid) = %d, want 0", e)
	}
}

func TestRoadmapPrefersLiveOverArchived(t *testing.T) {
	defer setupLifecycleTest(t)()
	mid := createMission(t, "Live Shipped Title")
	if e := setStateCmd([]string{"shipped"}); e != 0 {
		t.Fatalf("setStateCmd = %d", e)
	}

	// Write stale archive version with different title
	archDir := filepath.Join(cfg.ArchiveDir(), mid)
	if err := os.MkdirAll(archDir, 0755); err != nil {
		t.Fatalf("mkdir archive: %v", err)
	}
	stale := mission.Mission{ID: mid, Title: "Stale Archived Title", State: "shipped"}
	data, _ := json.Marshal(stale)
	os.WriteFile(filepath.Join(archDir, "mission.json"), data, 0644)

	if e := roadmapNewCmd([]string{"Prefer Live Roadmap"}); e != 0 {
		t.Fatalf("roadmapNewCmd = %d", e)
	}
	rms, _ := roadmapStore.List()
	rid := rms[0].ID

	if e := roadmapAddCmd([]string{rid, mid}); e != 0 {
		t.Fatalf("roadmapAddCmd = %d", e)
	}

	// roadmapShowCmd should use the live title, not the stale archive title
	if e := roadmapShowCmd([]string{rid}); e != 0 {
		t.Errorf("roadmapShowCmd = %d, want 0", e)
	}
}

func TestRoadmapCorruptArchiveTreatedAsNotFound(t *testing.T) {
	defer setupLifecycleTest(t)()
	mid := createMission(t, "Only Live")
	if e := setStateCmd([]string{"shipped"}); e != 0 {
		t.Fatalf("setStateCmd = %d", e)
	}

	badID := "MBADCORRUPT"
	archDir := filepath.Join(cfg.ArchiveDir(), badID)
	if err := os.MkdirAll(archDir, 0755); err != nil {
		t.Fatalf("mkdir archive: %v", err)
	}
	os.WriteFile(filepath.Join(archDir, "mission.json"), []byte("{not json"), 0644)

	if e := roadmapNewCmd([]string{"Corrupt Roadmap"}); e != 0 {
		t.Fatalf("roadmapNewCmd = %d", e)
	}
	rms, _ := roadmapStore.List()
	rid := rms[0].ID

	if e := roadmapAddCmd([]string{rid, mid}); e != 0 {
		t.Fatalf("add live = %d", e)
	}
	if e := roadmapAddCmd([]string{rid, badID}); e != 1 {
		t.Errorf("roadmapAddCmd(corrupt archive) = %d, want 1", e)
	}
}
