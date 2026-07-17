package resolver

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"spacecraft/internal/mission"
)

// mockStore implements mission.MissionStore for testing.
type mockStore struct {
	current  string
	sessions map[string]string
	missions map[string]*mission.Mission
}

func newMockStore() *mockStore {
	return &mockStore{
		sessions: make(map[string]string),
		missions: make(map[string]*mission.Mission),
	}
}

func (m *mockStore) ReadCurrent() (*string, error) {
	if m.current == "" {
		return nil, nil
	}
	return &m.current, nil
}
func (m *mockStore) WriteCurrent(id string) error {
	m.current = id
	return nil
}
func (m *mockStore) ClearCurrent() error {
	m.current = ""
	return nil
}
func (m *mockStore) ReadSession(sessionKey string) (*string, error) {
	id, ok := m.sessions[sessionKey]
	if !ok || id == "" {
		return nil, nil
	}
	return &id, nil
}
func (m *mockStore) WriteSession(sessionKey, id string) error {
	m.sessions[sessionKey] = id
	return nil
}
func (m *mockStore) ClearSession(sessionKey string) error {
	delete(m.sessions, sessionKey)
	return nil
}
func (m *mockStore) SessionDir() string                { return "" }
func (m *mockStore) SessionFilePath(string) string      { return "" }
func (m *mockStore) List() ([]mission.MissionRecord, error) {
	var records []mission.MissionRecord
	for id, mission := range m.missions {
		records = append(records, makeRecord(id, mission))
	}
	return records, nil
}
func (m *mockStore) Load(id string) (*mission.Mission, error) {
	mm, ok := m.missions[id]
	if !ok {
		return nil, os.ErrNotExist
	}
	return mm, nil
}
func (m *mockStore) Save(mm *mission.Mission) error {
	m.missions[mm.ID] = mm
	return nil
}
func (m *mockStore) Create(mm *mission.Mission) error {
	m.missions[mm.ID] = mm
	return nil
}
func (m *mockStore) LoadPlan(id string) (*mission.Plan, error) {
	return nil, os.ErrNotExist
}
func (m *mockStore) SavePlan(id string, p *mission.Plan) error { return nil }
func (m *mockStore) LoadReview(id string) (*mission.Review, error) {
	return nil, os.ErrNotExist
}
func (m *mockStore) SaveReview(id string, r *mission.Review) error { return nil }
func (m *mockStore) CountEvidence(id string) (int, error)          { return 0, nil }
func (m *mockStore) ReserveEvidencePath(id string) (string, string, string, error) {
	return "", "", "", nil
}
func (m *mockStore) AppendEvidence(id string, entry *mission.EvidenceEntry) error { return nil }
func (m *mockStore) ReadEvidenceEntries(id string) ([]mission.EvidenceEntry, error) {
	return nil, nil
}
func (m *mockStore) SpecExists(id string) bool          { return false }
func (m *mockStore) PlanExists(id string) bool          { return false }
func (m *mockStore) QuestionsExists(id string) bool     { return false }
func (m *mockStore) DecisionsExists(id string) bool     { return false }
func (m *mockStore) DesignExists(id string) bool        { return false }
func (m *mockStore) ReviewJSONExists(id string) bool    { return false }
func (m *mockStore) ReviewMDExists(id string) bool      { return false }
func (m *mockStore) MissionDir(id string) string        { return filepath.Join(".space", "missions", id) }
func (m *mockStore) ReadFile(id string, relPath string) ([]byte, error) {
	return nil, os.ErrNotExist
}
func (m *mockStore) WriteFile(id string, relPath string, data []byte) error { return nil }
func (m *mockStore) RemoveAll(id string) error                              { return nil }
func (m *mockStore) ArchiveMission(id, archiveDir string, compactM mission.CompactMission, compactP mission.CompactPlan, compactEvidence []mission.CompactEvidenceEntry, review *mission.Review) error {
	return nil
}

func addMission(m *mockStore, id, title, state string, branches ...string) {
	mm := &mission.Mission{
		ID:    id,
		Title: title,
		State: state,
	}
	if len(branches) > 0 {
		mm.Git.WorkBranch = &branches[0]
	}
	m.missions[id] = mm
}

func makeRecord(id string, mm *mission.Mission) mission.MissionRecord {
	br := &mm.Git
	var branches []string
	if br.WorkBranch != nil && *br.WorkBranch != "" {
		branches = append(branches, *br.WorkBranch)
	}
	if mm.WorkBranch != nil && *mm.WorkBranch != "" {
		branches = append(branches, *mm.WorkBranch)
	}
	if mm.Branch != nil && *mm.Branch != "" {
		branches = append(branches, *mm.Branch)
	}
	return mission.MissionRecord{
		ID:       id,
		Mission:  mm,
		Dir:      filepath.Join(".space", "missions", id),
		Active:   mm.State != "shipped",
		Branches: branches,
	}
}

// fakeGitRunner returns a runner that pretends to be on a given branch.
type fakeGitRunner struct {
	branch string
	repo   bool
}

func (f fakeGitRunner) Run(name string, args ...string) (int, string, string) {
	if name == "git" && len(args) > 0 {
		switch args[0] {
		case "rev-parse":
			if len(args) > 1 && args[1] == "--is-inside-work-tree" {
				if f.repo {
					return 0, "true", ""
				}
				return 1, "", ""
			}
		case "branch":
			return 0, f.branch, ""
		case "status":
			return 0, "", ""
		}
	}
	return 0, "", ""
}

func noEnv(string) string { return "" }

func TestResolver_NoMissions(t *testing.T) {
	store := newMockStore()
	rr := New(store, fakeGitRunner{}, noEnv)
	out := rr.Resolve("")
	if out.Safety != "none" {
		t.Errorf("expected safety=none, got %s", out.Safety)
	}
	if out.Selected != nil {
		t.Errorf("expected nil selected, got %v", out.Selected)
	}
}

func TestResolver_SelectByExplicitId(t *testing.T) {
	store := newMockStore()
	addMission(store, "M07TEST01", "Alpha", "draft")
	addMission(store, "M07TEST02", "Beta", "draft")

	rr := New(store, fakeGitRunner{}, noEnv)
	out := rr.Resolve("M07TEST02")
	if out.Selected == nil || out.Selected.ID != "M07TEST02" {
		t.Errorf("expected M07TEST02, got %v", out.Selected)
	}
	if out.Source == nil || *out.Source != "selector" {
		t.Errorf("expected selector source, got %v", out.Source)
	}
	if out.Safety != "safe" {
		t.Errorf("expected safe, got %s", out.Safety)
	}
}

func TestResolver_SelectByNumber(t *testing.T) {
	store := newMockStore()
	addMission(store, "M07TEST10", "Alpha", "draft")
	addMission(store, "M07TEST20", "Beta", "draft")

	rr := New(store, fakeGitRunner{}, noEnv)
	out := rr.Resolve("1")
	if out.Selected == nil {
		t.Fatal("expected a selected mission")
	}
	if out.Source == nil || *out.Source != "selector" {
		t.Errorf("expected selector source, got %v", out.Source)
	}
}

func TestResolver_SelectByTitle(t *testing.T) {
	store := newMockStore()
	addMission(store, "M07TEST30", "Unique Title Alpha", "draft")
	addMission(store, "M07TEST31", "Beta", "draft")

	rr := New(store, fakeGitRunner{}, noEnv)
	out := rr.Resolve("Unique Title Alpha")
	if out.Selected == nil || out.Selected.ID != "M07TEST30" {
		t.Errorf("expected M07TEST30, got %v", out.Selected)
	}
}

func TestResolver_SelectBySubstringTitle(t *testing.T) {
	store := newMockStore()
	addMission(store, "M07TEST40", "Long Mission Name Alpha", "draft")
	addMission(store, "M07TEST41", "Beta", "draft")

	rr := New(store, fakeGitRunner{}, noEnv)
	out := rr.Resolve("Alpha")
	if out.Selected == nil || out.Selected.ID != "M07TEST40" {
		t.Errorf("expected M07TEST40, got %v", out.Selected)
	}
}

func TestResolver_CurrentFallback(t *testing.T) {
	store := newMockStore()
	addMission(store, "M07TEST50", "Current", "draft")
	store.WriteCurrent("M07TEST50")

	rr := New(store, fakeGitRunner{}, noEnv)
	out := rr.Resolve("")
	if out.Selected == nil || out.Selected.ID != "M07TEST50" {
		t.Errorf("expected M07TEST50, got %v", out.Selected)
	}
	if out.Source == nil || *out.Source != ".space/current" {
		t.Errorf("expected .space/current source, got %v", out.Source)
	}
	if out.Safety != "safe" {
		t.Errorf("expected safe, got %s", out.Safety)
	}
}

func TestResolver_CurrentOverridesBranch(t *testing.T) {
	store := newMockStore()
	addMission(store, "M07TEST60", "Current", "draft")
	addMission(store, "M07TEST61", "Branch", "draft")
	store.WriteCurrent("M07TEST60")

	rr := New(store, fakeGitRunner{branch: "feat/m07test61-feature", repo: true}, noEnv)
	out := rr.Resolve("")
	if out.Selected == nil || out.Selected.ID != "M07TEST60" {
		t.Errorf("expected M07TEST60 (current overrides branch), got %v", out.Selected)
	}
	if out.Source == nil || *out.Source != ".space/current" {
		t.Errorf("expected .space/current source, got %v", out.Source)
	}
}

func TestResolver_SessionOverride(t *testing.T) {
	store := newMockStore()
	addMission(store, "M07TEST70", "Current", "draft")
	addMission(store, "M07TEST71", "Session", "draft")
	store.WriteCurrent("M07TEST70")
	store.WriteSession("test-session", "M07TEST71")

	rr := New(store, fakeGitRunner{branch: "feat/m07test70-feature", repo: true}, func(s string) string {
		if s == "SPACECRAFT_SESSION" {
			return "test-session"
		}
		return ""
	})
	out := rr.Resolve("")
	if out.Selected == nil || out.Selected.ID != "M07TEST71" {
		t.Errorf("expected M07TEST71 (session override), got %v", out.Selected)
	}
	if out.Source == nil || *out.Source != "session" {
		t.Errorf("expected session source, got %v", out.Source)
	}
}

func TestResolver_SingleActiveFallback(t *testing.T) {
	store := newMockStore()
	addMission(store, "M07TEST80", "Only Active", "planned")

	rr := New(store, fakeGitRunner{}, noEnv)
	out := rr.Resolve("")
	if out.Selected == nil || out.Selected.ID != "M07TEST80" {
		t.Errorf("expected M07TEST80, got %v", out.Selected)
	}
	if out.Source == nil || *out.Source != "single-active" {
		t.Errorf("expected single-active source, got %v", out.Source)
	}
	if out.Safety != "safe" {
		t.Errorf("expected safe, got %s", out.Safety)
	}
}

func TestResolver_MultipleActiveAmbiguous(t *testing.T) {
	store := newMockStore()
	addMission(store, "M07TEST90", "Alpha", "draft")
	addMission(store, "M07TEST91", "Beta", "draft")

	rr := New(store, fakeGitRunner{}, noEnv)
	out := rr.Resolve("")
	if out.Selected != nil {
		t.Errorf("expected nil selected, got %v", out.Selected)
	}
	if out.Safety != "ambiguous" {
		t.Errorf("expected ambiguous, got %s", out.Safety)
	}
	if len(out.Candidates) < 2 {
		t.Errorf("expected at least 2 candidates, got %d", len(out.Candidates))
	}
}

func TestResolver_RequireResolved_fails_on_conflict(t *testing.T) {
	store := newMockStore()
	addMission(store, "M07TEST95", "Current", "draft")
	addMission(store, "M07TEST96", "Branch", "draft")
	store.WriteCurrent("M07TEST95")

	rr := New(store, fakeGitRunner{branch: "feat/m07test96-feature", repo: true}, noEnv)
	_, err := rr.RequireResolved("test-command")
	if err == nil {
		t.Error("expected error for conflict")
	}
	if err != nil && !strings.Contains(err.Error(), "resolution") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestResolver_RequireResolved_ok(t *testing.T) {
	store := newMockStore()
	addMission(store, "M07TEST97", "Solo", "draft")
	store.WriteCurrent("M07TEST97")

	rr := New(store, fakeGitRunner{}, noEnv)
	out, err := rr.RequireResolved("test-command")
	if err != nil {
		t.Fatal(err)
	}
	if out.Selected == nil || out.Selected.ID != "M07TEST97" {
		t.Errorf("expected M07TEST97, got %v", out.Selected)
	}
}

func TestResolver_BranchMetadataMatch(t *testing.T) {
	store := newMockStore()
	addMission(store, "M07TEST98", "Meta match", "draft", "feature/my-feature")

	rr := New(store, fakeGitRunner{branch: "feature/my-feature", repo: true}, noEnv)
	out := rr.Resolve("")
	if out.Selected == nil || out.Selected.ID != "M07TEST98" {
		t.Errorf("expected M07TEST98, got %v", out.Selected)
	}
	if out.Source == nil || *out.Source != "branch-metadata" {
		t.Errorf("expected branch-metadata source, got %v", out.Source)
	}
	if out.Safety != "safe" {
		t.Errorf("expected safe, got %s", out.Safety)
	}
}

func TestResolver_SPACECRAFT_MISSION_env(t *testing.T) {
	store := newMockStore()
	addMission(store, "M07TEST99", "Env mission", "draft")

	rr := New(store, fakeGitRunner{}, func(s string) string {
		if s == "SPACECRAFT_MISSION" {
			return "M07TEST99"
		}
		return ""
	})
	out := rr.Resolve("")
	if out.Selected == nil || out.Selected.ID != "M07TEST99" {
		t.Errorf("expected M07TEST99, got %v", out.Selected)
	}
	if out.Source == nil || *out.Source != "SPACECRAFT_MISSION" {
		t.Errorf("expected SPACECRAFT_MISSION source, got %v", out.Source)
	}
}
