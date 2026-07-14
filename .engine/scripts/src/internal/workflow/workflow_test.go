package workflow

import (
	"os"
	"testing"

	"spacecraft/internal/mission"
)

func containsStr(list []string, sub string) bool {
	for _, s := range list {
		if len(s) >= len(sub) {
			for i := 0; i <= len(s)-len(sub); i++ {
				if s[i:i+len(sub)] == sub {
					return true
				}
			}
		}
	}
	return false
}

func TestNextTask_nil(t *testing.T) {
	nt := NextTask(nil)
	if nt != nil {
		t.Errorf("expected nil, got %v", nt)
	}
}

func TestNextTask_empty(t *testing.T) {
	nt := NextTask([]mission.Task{})
	if nt != nil {
		t.Errorf("expected nil, got %v", nt)
	}
}

func TestNextTask_returnsFirstOpen(t *testing.T) {
	tasks := []mission.Task{
		{ID: strPtr("T01"), Status: strPtr("done")},
		{ID: strPtr("T02"), Status: strPtr("pending")},
		{ID: strPtr("T03"), Status: strPtr("pending")},
	}
	nt := NextTask(tasks)
	if nt == nil || *nt.ID != "T02" {
		t.Errorf("expected T02, got %v", nt)
	}
}

func TestNextTask_allDone(t *testing.T) {
	tasks := []mission.Task{
		{ID: strPtr("T01"), Status: strPtr("done")},
		{ID: strPtr("T02"), Status: strPtr("done")},
		{ID: strPtr("T03"), Status: strPtr("cancelled")},
	}
	nt := NextTask(tasks)
	if nt != nil {
		t.Errorf("expected nil, got %v", nt)
	}
}

func TestNextTask_nilStatusIsOpen(t *testing.T) {
	tasks := []mission.Task{
		{ID: strPtr("T01"), Status: nil},
	}
	nt := NextTask(tasks)
	if nt == nil || *nt.ID != "T01" {
		t.Errorf("expected T01, got %v", nt)
	}
}

func TestNextTask_skipsWaiting(t *testing.T) {
	tasks := []mission.Task{
		{ID: strPtr("T01"), Status: strPtr("waiting")},
		{ID: strPtr("T02"), Status: strPtr("pending")},
	}
	nt := NextTask(tasks)
	if nt == nil || *nt.ID != "T02" {
		t.Errorf("expected T02 (skipped waiting T01), got %v", nt)
	}
}

func TestNextTask_allWaiting(t *testing.T) {
	tasks := []mission.Task{
		{ID: strPtr("T01"), Status: strPtr("waiting")},
		{ID: strPtr("T02"), Status: strPtr("waiting")},
	}
	nt := NextTask(tasks)
	if nt != nil {
		t.Errorf("expected nil (all waiting), got %v", nt)
	}
}

func TestNextCommand_states(t *testing.T) {
	tests := []struct {
		state     string
		clarified bool
		want      string
	}{
		{"draft", true, "/sc-plan"},
		{"planned", true, "/sc-build"},
		{"built", true, "/sc-review"},
		{"ready", true, "/sc-ship"},
		{"shipped", true, "(shipped)"},
		{"blocked", true, "/sc-resume"},
		{"", true, "/sc-resume"},
	}

	for _, tc := range tests {
		m := &mission.Mission{
			State: tc.state,
			Clarification: mission.ClarificationBlock{
				Status: "clear",
			},
		}
		got := NextCommand(m)
		if got != tc.want {
			t.Errorf("NextCommand(state=%q, clarified=%v) = %q, want %q", tc.state, tc.clarified, got, tc.want)
		}
	}
}

func TestNextCommand_nilMission(t *testing.T) {
	got := NextCommand(nil)
	if got != "/sc-resume" {
		t.Errorf("expected /sc-resume, got %q", got)
	}
}

func TestNextCommand_clarificationPriority(t *testing.T) {
	m := &mission.Mission{
		State: "planned",
		Clarification: mission.ClarificationBlock{
			Status:            "open",
			BlockingQuestions: 2,
		},
	}
	got := NextCommand(m)
	if got != "(clarify)" {
		t.Errorf("expected (clarify), got %q", got)
	}
}

func TestNextCommand_nonDefaultStatus(t *testing.T) {
	m := &mission.Mission{
		State: "draft",
		Clarification: mission.ClarificationBlock{
			Status: "deferred",
		},
	}
	got := NextCommand(m)
	if got != "/sc-plan" {
		t.Errorf("expected /sc-plan, got %q", got)
	}
}

func TestSnapshot_Build_basic(t *testing.T) {
	store := newMockStore(t)
	snapper := NewSnapshot(store)

	res := mission.ResolveOutput{
		Safety: "safe",
		Source: strPtr(".space/current"),
		Git: mission.GitInfoData{
			Branch: "main",
			IsRepo: true,
		},
	}

	snap, err := snapper.Build(res, "M07WF01")
	if err != nil {
		t.Fatal(err)
	}
	if snap.MissionID != "M07WF01" {
		t.Errorf("expected M07WF01, got %s", snap.MissionID)
	}
	if snap.State != "planned" {
		t.Errorf("expected planned, got %s", snap.State)
	}
	if snap.Safety != "safe" {
		t.Errorf("expected safe, got %s", snap.Safety)
	}
}

func TestSnapshot_Build_nextTask(t *testing.T) {
	store := newMockStore(t)
	store.writePlan("M07WF02", &mission.Plan{
		MissionId: "M07WF02",
		Tasks: []mission.Task{
			{ID: strPtr("T01"), Title: strPtr("First"), Status: strPtr("done")},
			{ID: strPtr("T02"), Title: strPtr("Second"), Status: strPtr("pending")},
		},
	})
	store.writeEvidence("M07WF02")
	store.writeSpec("M07WF02")

	snapper := NewSnapshot(store)
	res := mission.ResolveOutput{Safety: "safe", Source: strPtr(".space/current")}

	snap, err := snapper.Build(res, "M07WF02")
	if err != nil {
		t.Fatal(err)
	}
	if snap.NextTask == nil || *snap.NextTask.ID != "T02" {
		t.Errorf("expected next task T02, got %v", snap.NextTask)
	}
	if snap.Next != "/sc-build T02" {
		t.Errorf("expected /sc-build T02, got %q", snap.Next)
	}
}

func TestSnapshot_Build_clarificationBlock(t *testing.T) {
	store := newMockStore(t)
	store.mission.Clarification.Status = "open"
	store.mission.Clarification.BlockingQuestions = 1

	snapper := NewSnapshot(store)
	res := mission.ResolveOutput{Safety: "safe"}

	snap, err := snapper.Build(res, "M07WF03")
	if err != nil {
		t.Fatal(err)
	}
	if snap.Next != "(clarify)" {
		t.Errorf("expected (clarify), got %q", snap.Next)
	}
	if !containsStr(snap.Blockers, "blocking clarification") {
		t.Errorf("expected blocker about clarification, got %v", snap.Blockers)
	}
}

func TestSnapshot_Build_missingSpec(t *testing.T) {
	store := newMockStore(t)
	store.mission.State = "planned"
	store.writePlan("M07WF04", &mission.Plan{MissionId: "M07WF04", Tasks: []mission.Task{{ID: strPtr("T01")}}})

	snapper := NewSnapshot(store)
	res := mission.ResolveOutput{Safety: "safe"}

	snap, err := snapper.Build(res, "M07WF04")
	if err != nil {
		t.Fatal(err)
	}
	if !containsStr(snap.Blockers, "spec.md is missing") {
		t.Errorf("expected spec.md blocker, got %v", snap.Blockers)
	}
}

func TestSnapshot_Build_missingPlan(t *testing.T) {
	store := newMockStore(t)
	store.mission.State = "planned"

	snapper := NewSnapshot(store)
	res := mission.ResolveOutput{Safety: "safe"}

	snap, err := snapper.Build(res, "M07WF05")
	if err != nil {
		t.Fatal(err)
	}
	if !containsStr(snap.Blockers, "plan.json is missing") {
		t.Errorf("expected plan.json blocker, got %v", snap.Blockers)
	}
}

func TestSnapshot_Build_mainBranchBlock(t *testing.T) {
	store := newMockStore(t)
	store.mission.State = "planned"
	store.writeSpec("M07WF06")
	store.writePlan("M07WF06", &mission.Plan{
		MissionId: "M07WF06",
		Tasks:     []mission.Task{{ID: strPtr("T01"), Status: strPtr("pending")}},
	})

	snapper := NewSnapshot(store)
	res := mission.ResolveOutput{
		Safety: "safe",
		Git: mission.GitInfoData{
			Branch: "main",
			IsRepo: true,
		},
	}

	snap, err := snapper.Build(res, "M07WF06")
	if err != nil {
		t.Fatal(err)
	}
	if !containsStr(snap.Blockers, "non-main work branch") {
		t.Errorf("expected non-main branch blocker, got %v", snap.Blockers)
	}
}

func TestSnapshot_Build_allDone(t *testing.T) {
	store := newMockStore(t)
	store.mission.State = "ready"
	store.writeSpec("M07WF07")
	store.writePlan("M07WF07", &mission.Plan{
		MissionId: "M07WF07",
		Tasks:     []mission.Task{{ID: strPtr("T01"), Status: strPtr("done")}},
	})

	snapper := NewSnapshot(store)
	res := mission.ResolveOutput{Safety: "safe"}

	snap, err := snapper.Build(res, "M07WF07")
	if err != nil {
		t.Fatal(err)
	}
	if snap.Next != "/sc-ship" {
		t.Errorf("expected /sc-ship, got %q", snap.Next)
	}
	if snap.Tasks.Completed != 1 {
		t.Errorf("expected 1 completed, got %d", snap.Tasks.Completed)
	}
}

func TestSnapshot_Build_allWaitingTasks(t *testing.T) {
	store := newMockStore(t)
	store.mission.State = "implementing"
	store.writeSpec("M07WF08")
	store.writePlan("M07WF08", &mission.Plan{
		MissionId: "M07WF08",
		Tasks: []mission.Task{
			{ID: strPtr("T01"), Status: strPtr("done")},
			{ID: strPtr("T02"), Status: strPtr("waiting")},
			{ID: strPtr("T03"), Status: strPtr("waiting")},
		},
	})

	snapper := NewSnapshot(store)
	res := mission.ResolveOutput{Safety: "safe"}

	snap, err := snapper.Build(res, "M07WF08")
	if err != nil {
		t.Fatal(err)
	}
	if !containsStr(snap.Blockers, "architectural guidance") {
		t.Errorf("expected architect-tasks blocker, got %v", snap.Blockers)
	}
	if snap.Next != "/sc-resume" {
		t.Errorf("expected /sc-resume when all open tasks are waiting, got %q", snap.Next)
	}
	if snap.Tasks.Completed != 1 {
		t.Errorf("expected 1 completed, got %d", snap.Tasks.Completed)
	}
}

// --- mock store ---

type mockStore struct {
	mission   *mission.Mission
	spec      bool
	plan      *mission.Plan
	planExist bool
	evidence  int
}

func newMockStore(t *testing.T) *mockStore {
	t.Helper()
	return &mockStore{
		mission: &mission.Mission{
			ID:    "M07WF01",
			Title: "Workflow test mission",
			State: "planned",
		},
	}
}

func (m *mockStore) readMission() *mission.Mission { return m.mission }

func (m *mockStore) writeSpec(id string)            { m.spec = true }
func (m *mockStore) writeEvidence(id string)        { m.evidence = 1 }
func (m *mockStore) writePlan(id string, p *mission.Plan) {
	m.plan = p
	m.planExist = true
}

// --- MissionStore interface ---

func (m *mockStore) ReadCurrent() (*string, error)              { return strPtr("M07WF01"), nil }
func (m *mockStore) WriteCurrent(id string) error               { return nil }
func (m *mockStore) ClearCurrent() error                        { return nil }
func (m *mockStore) ReadSession(string) (*string, error)        { return nil, nil }
func (m *mockStore) WriteSession(string, string) error          { return nil }
func (m *mockStore) ClearSession(string) error                  { return nil }
func (m *mockStore) SessionDir() string                         { return "" }
func (m *mockStore) SessionFilePath(string) string              { return "" }
func (m *mockStore) List() ([]mission.MissionRecord, error)     { return nil, nil }
func (m *mockStore) Load(id string) (*mission.Mission, error)   { return m.mission, nil }
func (m *mockStore) Save(mm *mission.Mission) error             { m.mission = mm; return nil }
func (m *mockStore) Create(mm *mission.Mission) error           { return nil }
func (m *mockStore) LoadPlan(id string) (*mission.Plan, error)  { return m.plan, nil }
func (m *mockStore) SavePlan(id string, p *mission.Plan) error  { m.plan = p; return nil }
func (m *mockStore) LoadReview(id string) (*mission.Review, error) {
	return nil, os.ErrNotExist
}
func (m *mockStore) SaveReview(id string, r *mission.Review) error { return nil }
func (m *mockStore) CountEvidence(id string) (int, error)          { return m.evidence, nil }
func (m *mockStore) ReserveEvidencePath(id string) (string, string, string, error) {
	return "", "", "", nil
}
func (m *mockStore) AppendEvidence(id string, entry *mission.EvidenceEntry) error { return nil }
func (m *mockStore) ReadEvidenceEntries(id string) ([]mission.EvidenceEntry, error) {
	return nil, nil
}
func (m *mockStore) SpecExists(id string) bool          { return m.spec }
func (m *mockStore) PlanExists(id string) bool          { return m.planExist }
func (m *mockStore) QuestionsExists(id string) bool     { return false }
func (m *mockStore) DecisionsExists(id string) bool     { return false }
func (m *mockStore) DesignExists(id string) bool        { return false }
func (m *mockStore) ReviewJSONExists(id string) bool    { return false }
func (m *mockStore) ReviewMDExists(id string) bool      { return false }
func (m *mockStore) MissionDir(id string) string        { return "" }
func (m *mockStore) ReadFile(id string, relPath string) ([]byte, error) {
	return nil, os.ErrNotExist
}
func (m *mockStore) WriteFile(id string, relPath string, data []byte) error { return nil }
func (m *mockStore) RemoveAll(id string) error                              { return nil }
func (m *mockStore) ArchiveMission(id, archiveDir string, compactM mission.CompactMission, compactP mission.CompactPlan, compactEvidence []mission.CompactEvidenceEntry, review *mission.Review) error {
	return nil
}

// misc

func strPtr(s string) *string { return &s }
