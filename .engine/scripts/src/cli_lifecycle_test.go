package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"spacecraft/internal/archive"
	"spacecraft/internal/closeout"
	"spacecraft/internal/config"
	"spacecraft/internal/gitutil"
	"spacecraft/internal/hooks"
	"spacecraft/internal/mission"
	"spacecraft/internal/resolver"
	"spacecraft/internal/roadmap"
	"spacecraft/internal/state"
	"spacecraft/internal/trace"
	"spacecraft/internal/workflow"
)

// setupLifecycleTest swaps the CLI globals to an isolated temp directory for the
// duration of the test and restores them afterwards.
func setupLifecycleTest(t *testing.T) func() {
	oldCfg := cfg
	oldStore := store
	oldTraceStore := traceStore
	oldRoadmapStore := roadmapStore
	oldR := r
	oldSS := ss
	oldWS := ws
	oldCC := cc
	oldArc := arc
	oldAr := ar
	oldHooksCfg := hooksCfg

	tmp := t.TempDir()
	t.Chdir(tmp)

	c, err := config.NewConfig(tmp)
	if err != nil {
		t.Fatalf("config: %v", err)
	}
	s := mission.NewFSStore(c)
	rr := resolver.New(s, gitutil.OSCommandRunner{}, nil)
	setter := state.NewSetter(s)
	snap := workflow.NewSnapshot(s)
	snap.SetCommandsDir(filepath.Join(c.Root(), ".opencode", "commands"))
	checker := closeout.NewChecker(s, gitutil.OSCommandRunner{})
	rc := archive.NewReadinessChecker(s)
	arch := archive.NewArchiver(s)
	ts := trace.NewFSTraceStore(c)
	rs := roadmap.NewFSStore(c)
	hc, _ := hooks.LoadConfig(filepath.Join(c.SpaceDir(), "hooks.json"))

	cfg = c
	store = s
	traceStore = ts
	roadmapStore = rs
	r = rr
	ss = setter
	ws = snap
	cc = checker
	arc = rc
	ar = arch
	hooksCfg = hc

	return func() {
		cfg = oldCfg
		store = oldStore
		traceStore = oldTraceStore
		roadmapStore = oldRoadmapStore
		r = oldR
		ss = oldSS
		ws = oldWS
		cc = oldCC
		arc = oldArc
		ar = oldAr
		hooksCfg = oldHooksCfg
	}
}

func runCmd(t *testing.T, name string, args ...string) {
	t.Helper()
	cmd := exec.Command(name, args...)
	cmd.Dir = cfg.Root()
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("%s %v: %v\n%s", name, args, err, out)
	}
}

func initGit(t *testing.T) {
	t.Helper()
	runCmd(t, "git", "init")
	runCmd(t, "git", "config", "user.email", "test@example.com")
	runCmd(t, "git", "config", "user.name", "Test")
	keep := filepath.Join(cfg.Root(), ".gitkeep")
	if f, err := os.Create(keep); err == nil {
		f.Close()
	}
	runCmd(t, "git", "add", ".gitkeep")
	runCmd(t, "git", "commit", "-m", "feat: initial commit")
}

func createMission(t *testing.T, title string) string {
	t.Helper()
	if code := newCmd([]string{title}); code != 0 {
		t.Fatalf("newCmd(%q) = %d, want 0", title, code)
	}
	id, err := store.ReadCurrent()
	if err != nil || id == nil {
		t.Fatalf("no current mission after newCmd: %v", err)
	}
	return *id
}

func makeShippableMission(t *testing.T) string {
	t.Helper()
	id := createMission(t, "Shippable Mission")

	done := "done"
	tid := "T01"
	title := "do the thing"
	if err := store.SavePlan(id, &mission.Plan{
		MissionId: id,
		Tasks:     []mission.Task{{ID: &tid, Title: &title, Status: &done}},
	}); err != nil {
		t.Fatalf("SavePlan: %v", err)
	}

	if code := evidenceCmd([]string{"tests pass", "--", "echo", "hello"}); code != 0 {
		t.Fatalf("evidenceCmd = %d, want 0", code)
	}

	ready := "ready"
	doneStatus := "done"
	if err := store.SaveReview(id, &mission.Review{
		Status: &ready,
		ReleaseReadiness: mission.ReleaseReadiness{
			Version:                &mission.ReleaseGate{Status: &doneStatus},
			Changelog:              &mission.ReleaseGate{Status: &doneStatus},
			SpecNote:               &mission.ReleaseGate{Status: &doneStatus},
			TagPlan:                &mission.ReleaseGate{Status: &doneStatus},
			PostRebaseVerification: &mission.ReleaseGate{Status: &doneStatus},
			EvalCoverage:           &mission.ReleaseGate{Status: &doneStatus},
		},
	}); err != nil {
		t.Fatalf("SaveReview: %v", err)
	}

	if code := clarifyStatusCmd([]string{"clear"}); code != 0 {
		t.Fatalf("clarifyStatusCmd = %d, want 0", code)
	}
	if code := setStateCmd([]string{"shipped"}); code != 0 {
		t.Fatalf("setStateCmd = %d, want 0", code)
	}
	return id
}

type fakeGitRunner struct{}

func (fakeGitRunner) Run(name string, args ...string) (int, string, string) {
	cmd := strings.Join(append([]string{name}, args...), " ")
	switch cmd {
	case "git rev-parse --is-inside-work-tree":
		return 0, "true", ""
	case "git branch --show-current":
		return 0, "feat/test", ""
	case "git status --short":
		return 0, "", ""
	case "git merge-base --is-ancestor main HEAD":
		return 0, "", ""
	case "git rev-list --count main..HEAD":
		return 0, "1", ""
	case "git log --format=%s main..HEAD":
		return 0, "feat: do thing", ""
	case "git log --oneline main..HEAD -- CHANGELOG.md":
		return 0, "abc123 chore: bump version", ""
	}
	return 1, "", ""
}

func TestInitCmd(t *testing.T) {
	defer setupLifecycleTest(t)()
	if code := initCmd([]string{}); code != 0 {
		t.Errorf("initCmd() = %d, want 0", code)
	}
	if _, err := os.Stat(cfg.MissionsDir()); err != nil {
		t.Errorf("missions dir missing: %v", err)
	}
	if _, err := os.Stat(cfg.CurrentFile()); err != nil {
		t.Errorf("current file missing: %v", err)
	}
	if code := initCmd([]string{}); code != 0 {
		t.Errorf("initCmd() idempotent = %d, want 0", code)
	}
}

func TestNewCmd(t *testing.T) {
	defer setupLifecycleTest(t)()
	if code := newCmd([]string{"Test Mission"}); code != 0 {
		t.Errorf("newCmd(title) = %d, want 0", code)
	}
	id, _ := store.ReadCurrent()
	if id == nil {
		t.Error("no current mission after newCmd")
	}
	if code := newCmd([]string{}); code != 1 {
		t.Errorf("newCmd(empty) = %d, want 1", code)
	}
}

func TestCurrentCmd(t *testing.T) {
	defer setupLifecycleTest(t)()
	if code := currentCmd(); code != 0 {
		t.Errorf("currentCmd() empty = %d, want 0", code)
	}
	id := createMission(t, "Current")
	if code := currentCmd(); code != 0 {
		t.Errorf("currentCmd() with mission = %d, want 0", code)
	}
	curr, _ := store.ReadCurrent()
	if curr == nil || *curr != id {
		t.Errorf("current = %v, want %s", curr, id)
	}
}

func TestResolveCmd(t *testing.T) {
	defer setupLifecycleTest(t)()
	if code := resolveCmd([]string{}); code != 1 {
		t.Errorf("resolveCmd() empty = %d, want 1", code)
	}
	id := createMission(t, "Resolve")
	if code := resolveCmd([]string{}); code != 0 {
		t.Errorf("resolveCmd() with current = %d, want 0", code)
	}
	if code := resolveCmd([]string{id}); code != 0 {
		t.Errorf("resolveCmd(id) = %d, want 0", code)
	}
	if code := resolveCmd([]string{id, "--json"}); code != 0 {
		t.Errorf("resolveCmd(id --json) = %d, want 0", code)
	}
	if code := resolveCmd([]string{"--json"}); code != 0 {
		t.Errorf("resolveCmd(--json) = %d, want 0", code)
	}
	if code := resolveCmd([]string{"bad"}); code != 1 {
		t.Errorf("resolveCmd(bad) = %d, want 1", code)
	}
	if code := resolveCmd([]string{"bad", "--json"}); code != 1 {
		t.Errorf("resolveCmd(bad --json) = %d, want 1", code)
	}
}

func TestMissionsCmd(t *testing.T) {
	defer setupLifecycleTest(t)()
	if code := missionsCmd(); code != 0 {
		t.Errorf("missionsCmd() empty = %d, want 0", code)
	}
	createMission(t, "Listed")
	if code := missionsCmd(); code != 0 {
		t.Errorf("missionsCmd() with mission = %d, want 0", code)
	}
}

func TestUseCmd(t *testing.T) {
	defer setupLifecycleTest(t)()
	id1 := createMission(t, "First")
	if code := useCmd([]string{"1"}); code != 0 {
		t.Errorf("useCmd(1) = %d, want 0", code)
	}
	curr, _ := store.ReadCurrent()
	if curr == nil || *curr != id1 {
		t.Errorf("use by number = %v, want %s", curr, id1)
	}
	if code := useCmd([]string{id1}); code != 0 {
		t.Errorf("useCmd(id) = %d, want 0", code)
	}
	if code := useCmd([]string{"First"}); code != 0 {
		t.Errorf("useCmd(title) = %d, want 0", code)
	}

	id2 := createMission(t, "Second")
	if code := useCmd([]string{id2}); code != 0 {
		t.Errorf("useCmd(second id) = %d, want 0", code)
	}
	if code := useCmd([]string{"Seco"}); code != 0 {
		t.Errorf("useCmd(substring) = %d, want 0", code)
	}
	curr, _ = store.ReadCurrent()
	if curr == nil || *curr != id2 {
		t.Errorf("use by substring = %v, want %s", curr, id2)
	}

	if code := useCmd([]string{}); code != 1 {
		t.Errorf("useCmd(empty) = %d, want 1", code)
	}
	if code := useCmd([]string{"zzzzzzzz"}); code != 1 {
		t.Errorf("useCmd(notfound) = %d, want 1", code)
	}
}

func TestBindBranchCmd(t *testing.T) {
	defer setupLifecycleTest(t)()
	initGit(t)
	id := createMission(t, "Bind Branch")
	if code := bindBranchCmd([]string{}); code != 0 {
		t.Errorf("bindBranchCmd() = %d, want 0", code)
	}
	m, _ := store.Load(id)
	if m.WorkBranch == nil || *m.WorkBranch == "" {
		t.Error("WorkBranch not bound")
	}
	if code := bindBranchCmd([]string{id}); code != 0 {
		t.Errorf("bindBranchCmd(id) = %d, want 0", code)
	}
	if code := bindBranchCmd([]string{"nope"}); code != 1 {
		t.Errorf("bindBranchCmd(notfound) = %d, want 1", code)
	}
}

func TestBindBranchCmdNotRepo(t *testing.T) {
	defer setupLifecycleTest(t)()
	id := createMission(t, "No Git")
	if code := bindBranchCmd([]string{id}); code != 1 {
		t.Errorf("bindBranchCmd(not repo) = %d, want 1", code)
	}
}

func TestStatusCmd(t *testing.T) {
	defer setupLifecycleTest(t)()
	if code := statusCmd([]string{}); code != 0 {
		t.Errorf("statusCmd() empty = %d, want 0", code)
	}
	createMission(t, "Status")
	if code := statusCmd([]string{}); code != 0 {
		t.Errorf("statusCmd() with mission = %d, want 0", code)
	}
}

func TestSetStateCmd(t *testing.T) {
	defer setupLifecycleTest(t)()
	id := createMission(t, "Set State")
	if code := setStateCmd([]string{"planned"}); code != 0 {
		t.Errorf("setStateCmd(planned) = %d, want 0", code)
	}
	m, _ := store.Load(id)
	if m.State != "planned" {
		t.Errorf("state = %q, want planned", m.State)
	}
	if code := setStateCmd([]string{}); code != 1 {
		t.Errorf("setStateCmd(empty) = %d, want 1", code)
	}
	if code := setStateCmd([]string{"invalid"}); code != 1 {
		t.Errorf("setStateCmd(invalid) = %d, want 1", code)
	}
}

func TestClarifyStatusCmd(t *testing.T) {
	defer setupLifecycleTest(t)()
	id := createMission(t, "Clarify")
	if code := clarifyStatusCmd([]string{"deferred"}); code != 0 {
		t.Errorf("clarifyStatusCmd(deferred) = %d, want 0", code)
	}
	m, _ := store.Load(id)
	if m.Clarification.Status != "deferred" {
		t.Errorf("clarification = %q, want deferred", m.Clarification.Status)
	}
	if code := clarifyStatusCmd([]string{}); code != 1 {
		t.Errorf("clarifyStatusCmd(empty) = %d, want 1", code)
	}
	if code := clarifyStatusCmd([]string{"bad"}); code != 1 {
		t.Errorf("clarifyStatusCmd(bad) = %d, want 1", code)
	}
}

func TestEvidenceCmd(t *testing.T) {
	defer setupLifecycleTest(t)()
	id := createMission(t, "Evidence")
	if code := evidenceCmd([]string{"tests pass", "--", "echo", "hello"}); code != 0 {
		t.Errorf("evidenceCmd = %d, want 0", code)
	}
	count, _ := store.CountEvidence(id)
	if count != 1 {
		t.Errorf("evidence count = %d, want 1", count)
	}
	if code := evidenceCmd([]string{"label", "echo", "hello"}); code != 1 {
		t.Errorf("evidenceCmd(missing --) = %d, want 1", code)
	}
	if code := evidenceCmd([]string{"", "--", "echo", "hi"}); code != 1 {
		t.Errorf("evidenceCmd(empty label) = %d, want 1", code)
	}
	if code := evidenceCmd([]string{"label", "--"}); code != 1 {
		t.Errorf("evidenceCmd(empty command) = %d, want 1", code)
	}
	if code := evidenceCmd([]string{"empty", "--", "true"}); code != 1 {
		t.Errorf("evidenceCmd(empty stdout) = %d, want 1", code)
	}
}

func TestEvidenceCmdNoMission(t *testing.T) {
	defer setupLifecycleTest(t)()
	if code := evidenceCmd([]string{"label", "--", "echo", "hi"}); code != 1 {
		t.Errorf("evidenceCmd(no mission) = %d, want 1", code)
	}
}

func TestExecCmd(t *testing.T) {
	defer setupLifecycleTest(t)()
	res := execCmd([]string{"echo", "hello"})
	if res.code != 0 || res.stdout != "hello\n" {
		t.Errorf("execCmd(echo) = (%d, %q), want (0, hello\\n)", res.code, res.stdout)
	}
	res = execCmd([]string{})
	if res.code != 1 {
		t.Errorf("execCmd(empty) = %d, want 1", res.code)
	}
	res = execCmd([]string{"false"})
	if res.code != 1 {
		t.Errorf("execCmd(false) = %d, want 1", res.code)
	}
	res = execCmd([]string{"/nonexistent-binary-xyz"})
	if res.code != 127 {
		t.Errorf("execCmd(missing) = %d, want 127", res.code)
	}
}

func TestValidateCmd(t *testing.T) {
	defer setupLifecycleTest(t)()
	id := createMission(t, "Validate")
	if code := validateCmd(); code != 0 {
		t.Errorf("validateCmd() = %d, want 0", code)
	}
	os.Remove(filepath.Join(store.MissionDir(id), "spec.md"))
	if code := validateCmd(); code != 1 {
		t.Errorf("validateCmd(invalid) = %d, want 1", code)
	}
}

func TestCloseoutCmdBlocked(t *testing.T) {
	defer setupLifecycleTest(t)()
	createMission(t, "Closeout")
	if code := closeoutCmd(); code != 1 {
		t.Errorf("closeoutCmd() = %d, want 1", code)
	}
}

func TestCloseoutCmdReady(t *testing.T) {
	defer setupLifecycleTest(t)()
	makeShippableMission(t)
	cc = closeout.NewChecker(store, fakeGitRunner{})
	if code := closeoutCmd(); code != 0 {
		t.Errorf("closeoutCmd() ready = %d, want 0", code)
	}
}

func TestArchiveCmdBlocked(t *testing.T) {
	defer setupLifecycleTest(t)()
	id := createMission(t, "Archive")
	if code := archiveCmd([]string{id}); code != 1 {
		t.Errorf("archiveCmd(not shipped) = %d, want 1", code)
	}
}

func TestArchiveCmdNoArgsNoCurrent(t *testing.T) {
	defer setupLifecycleTest(t)()
	if code := archiveCmd([]string{}); code != 1 {
		t.Errorf("archiveCmd(no current) = %d, want 1", code)
	}
}

func TestArchiveCmdHappy(t *testing.T) {
	defer setupLifecycleTest(t)()
	id := makeShippableMission(t)
	if code := archiveCmd([]string{id}); code != 0 {
		t.Errorf("archiveCmd() = %d, want 0", code)
	}
	if _, err := store.Load(id); err == nil {
		t.Error("mission source not removed after archive")
	}
}

func TestRequireResolved(t *testing.T) {
	defer setupLifecycleTest(t)()
	id := createMission(t, "Resolved")
	res := requireResolved("test")
	if res.Selected == nil || res.Selected.ID != id {
		t.Errorf("requireResolved selected = %v, want %s", res.Selected, id)
	}
}

func TestMissionDisplayRecords(t *testing.T) {
	recs := []mission.MissionRecord{
		{ID: "M00000002", Active: false},
		{ID: "M00000001", Active: true},
		{ID: "M00000003", Active: true},
	}
	got := missionDisplayRecords(recs)
	if len(got) != 3 {
		t.Fatalf("len = %d, want 3", len(got))
	}
	if !got[0].Active || got[0].ID != "M00000003" {
		t.Errorf("first = %+v, want active M00000003", got[0])
	}
	if !got[1].Active || got[1].ID != "M00000001" {
		t.Errorf("second = %+v, want active M00000001", got[1])
	}
	if got[2].Active || got[2].ID != "M00000002" {
		t.Errorf("third = %+v, want shipped M00000002", got[2])
	}
}

func TestDisplayTitle(t *testing.T) {
	if got := displayTitle(""); got != "(untitled)" {
		t.Errorf("displayTitle(empty) = %q", got)
	}
	if got := displayTitle("  trim  "); got != "trim" {
		t.Errorf("displayTitle(trim) = %q", got)
	}
	long := strings.Repeat("a", 90)
	if got := displayTitle(long); got != long[:85]+"..." {
		t.Errorf("displayTitle(long) = %q", got)
	}
}

func TestFindSelector(t *testing.T) {
	records := []mission.MissionRecord{
		{ID: "M07FYB5W5", Mission: &mission.Mission{Title: "Alpha"}},
		{ID: "M07FYB5W6", Mission: &mission.Mission{Title: "Beta"}},
	}
	ordered := missionDisplayRecords(records)

	if got := findSelector(records, "1", ordered); got == nil || got.ID != "M07FYB5W6" {
		t.Errorf("findSelector(1) = %v, want M07FYB5W6", got)
	}
	if got := findSelector(records, "M07FYB5W5", ordered); got == nil || got.ID != "M07FYB5W5" {
		t.Errorf("findSelector(id) = %v", got)
	}
	if got := findSelector(records, "Beta", ordered); got == nil || got.ID != "M07FYB5W6" {
		t.Errorf("findSelector(title) = %v", got)
	}
	if got := findSelector(records, "Alp", ordered); got == nil || got.ID != "M07FYB5W5" {
		t.Errorf("findSelector(substring) = %v", got)
	}
	if got := findSelector(records, "", ordered); got != nil {
		t.Errorf("findSelector(empty) = %v, want nil", got)
	}
	if got := findSelector(records, "99", ordered); got != nil {
		t.Errorf("findSelector(bad number) = %v, want nil", got)
	}
}

func TestWriteSessionBinding(t *testing.T) {
	defer setupLifecycleTest(t)()
	id := createMission(t, "Session")
	t.Setenv("SPACECRAFT_SESSION", "sess-key")
	path := writeSessionBinding(id)
	if path == "" {
		t.Fatal("writeSessionBinding returned empty")
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read session file: %v", err)
	}
	if strings.TrimSpace(string(data)) != id {
		t.Errorf("session file = %q, want %s", data, id)
	}
}

func TestRel(t *testing.T) {
	defer setupLifecycleTest(t)()
	p := filepath.Join(cfg.Root(), ".space", "current")
	if got := rel(p); got != filepath.Join(".space", "current") {
		t.Errorf("rel(under root) = %q", got)
	}
	if got := rel("/tmp/other/path"); !strings.HasPrefix(got, "..") {
		t.Errorf("rel(outside root) = %q, expected relative with ..", got)
	}
}

func TestIsoNow(t *testing.T) {
	s := isoNow()
	if !strings.HasSuffix(s, "Z") || !strings.Contains(s, "T") {
		t.Errorf("isoNow() = %q", s)
	}
	if _, err := time.Parse("2006-01-02T15:04:05.000Z", s); err != nil {
		t.Errorf("isoNow() parse: %v", err)
	}
}

func TestStrPtr(t *testing.T) {
	if strPtr("") != nil {
		t.Error("strPtr(empty) != nil")
	}
	p := strPtr("x")
	if p == nil || *p != "x" {
		t.Errorf("strPtr(x) = %v", p)
	}
}

func TestIfStr(t *testing.T) {
	if ifStr(true, "a", "b") != "a" {
		t.Error("ifStr(true)")
	}
	if ifStr(false, "a", "b") != "b" {
		t.Error("ifStr(false)")
	}
}

func TestNormMissionID(t *testing.T) {
	if got := normMissionID("M-20250714-123456"); got == nil || *got != "M-20250714-123456" {
		t.Errorf("normMissionID(legacy) = %v", got)
	}
	if got := normMissionID("m07fyb5w5"); got == nil || *got != "M07FYB5W5" {
		t.Errorf("normMissionID(compact) = %v", got)
	}
	if got := normMissionID("prefix M07FYB5W5 suffix"); got == nil || *got != "M07FYB5W5" {
		t.Errorf("normMissionID(embedded) = %v", got)
	}
	if got := normMissionID("abc"); got != nil {
		t.Errorf("normMissionID(none) = %v, want nil", got)
	}
}

func TestSplitFlags(t *testing.T) {
	flags, pos := splitFlags([]string{"--json", "id", "--verbose"})
	if len(flags) != 2 || len(pos) != 1 || pos[0] != "id" {
		t.Errorf("splitFlags = (%v, %v)", flags, pos)
	}
}

func TestSlugify(t *testing.T) {
	if got := slugify("Hello World!"); got != "hello-world" {
		t.Errorf("slugify = %q", got)
	}
	if got := slugify("  "); got != "mission" {
		t.Errorf("slugify(empty) = %q", got)
	}
	long := strings.Repeat("a", 70)
	if got := slugify(long); len(got) != 60 {
		t.Errorf("slugify(long) len = %d", len(got))
	}
}

func TestCmdToStr(t *testing.T) {
	if got := cmdToStr([]string{"echo", "hello"}); got != "echo hello" {
		t.Errorf("cmdToStr = %q", got)
	}
	if got := cmdToStr([]string{"echo", "hello world"}); got != "echo 'hello world'" {
		t.Errorf("cmdToStr(quote) = %q", got)
	}
}
