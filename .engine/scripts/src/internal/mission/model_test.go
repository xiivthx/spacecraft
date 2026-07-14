package mission

import (
	"encoding/json"
	"testing"
)

func TestMissionJSONRoundTrip(t *testing.T) {
	m := Mission{
		ID:    "M07ABCDEF",
		Title: "Test Mission",
		State: "draft",
	}
	data, err := json.Marshal(m)
	if err != nil {
		t.Fatalf("Marshal Mission: %v", err)
	}

	var got Mission
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("Unmarshal Mission: %v", err)
	}
	if got.ID != "M07ABCDEF" || got.Title != "Test Mission" || got.State != "draft" {
		t.Errorf("round-trip got %+v", got)
	}
}

func TestPlanRoundTrip(t *testing.T) {
	id := "T01"
	title := "Do something"
	status := "pending"
	p := Plan{
		MissionId: "M07ABCDEF",
		Tasks: []Task{
			{ID: &id, Title: &title, Status: &status},
		},
	}
	data, err := json.Marshal(p)
	if err != nil {
		t.Fatal(err)
	}
	var got Plan
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatal(err)
	}
	if got.MissionId != "M07ABCDEF" {
		t.Errorf("MissionId = %q", got.MissionId)
	}
	if len(got.Tasks) != 1 {
		t.Fatalf("expected 1 task, got %d", len(got.Tasks))
	}
	if *got.Tasks[0].ID != "T01" || *got.Tasks[0].Title != "Do something" {
		t.Errorf("task round-trip: %+v", got.Tasks[0])
	}
}

func TestEvidenceEntryRoundTrip(t *testing.T) {
	e := EvidenceEntry{
		ID:       "E07ABCDEF",
		Label:    "test",
		Command:  "echo hi",
		ExitCode: 0,
		Stdout:   "outputs/E07ABCDEF.stdout.txt",
	}
	data, err := json.Marshal(e)
	if err != nil {
		t.Fatal(err)
	}
	var got EvidenceEntry
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatal(err)
	}
	if got.ID != "E07ABCDEF" || got.Label != "test" || got.ExitCode != 0 {
		t.Errorf("round-trip: %+v", got)
	}
}

func TestGitInfoDataDefaults(t *testing.T) {
	var g GitInfoData
	if g.Available || g.IsRepo || g.Branch != "" {
		t.Error("zero value GitInfoData should have empty fields")
	}
}

func TestResolveOutputRoundTrip(t *testing.T) {
	mid := "M07ABCDEF"
	safety := "safe"
	out := ResolveOutput{
		Selected: &MissionInfo{ID: mid, Title: "Test", State: "draft"},
		Safety:   safety,
	}
	data, err := json.Marshal(out)
	if err != nil {
		t.Fatal(err)
	}
	var got ResolveOutput
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatal(err)
	}
	if got.Selected == nil || got.Selected.ID != mid {
		t.Errorf("selected.ID = %v, want %s", got.Selected, mid)
	}
	if got.Safety != safety {
		t.Errorf("Safety = %q, want %q", got.Safety, safety)
	}
}

func TestWorkflowSnapshotRoundTrip(t *testing.T) {
	ws := WorkflowSnapshot{
		MissionID: "M07ABCDEF",
		Title:     "Test",
		State:     "planned",
		Next:      "/sc-build",
	}
	data, err := json.Marshal(ws)
	if err != nil {
		t.Fatal(err)
	}
	var got WorkflowSnapshot
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatal(err)
	}
	if got.MissionID != "M07ABCDEF" || got.State != "planned" {
		t.Errorf("round-trip: %+v", got)
	}
}

func TestReviewRoundTrip(t *testing.T) {
	status := "ready"
	r := Review{
		Status: &status,
		Findings: []Finding{
			{Summary: stringPtr("All good")},
		},
		ReleaseReadiness: ReleaseReadiness{
			Version: &ReleaseGate{Status: stringPtr("ok")},
		},
	}
	data, err := json.Marshal(r)
	if err != nil {
		t.Fatal(err)
	}
	var got Review
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatal(err)
	}
	if got.Status == nil || *got.Status != "ready" {
		t.Errorf("Status = %v, want 'ready'", got.Status)
	}
	if len(got.Findings) != 1 {
		t.Errorf("expected 1 finding, got %d", len(got.Findings))
	}
}

func TestCompactMissionRoundTrip(t *testing.T) {
	cm := CompactMission{
		ID:    "M07ABCDEF",
		Title: "Archived",
		State: "shipped",
	}
	data, err := json.Marshal(cm)
	if err != nil {
		t.Fatal(err)
	}
	var got CompactMission
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatal(err)
	}
	if got.ID != "M07ABCDEF" || got.State != "shipped" {
		t.Errorf("round-trip: %+v", got)
	}
}

func TestMissionRecordFields(t *testing.T) {
	r := MissionRecord{
		ID:   "M07ABCDEF",
		Dir:  "/tmp/.space/missions/M07ABCDEF",
		Active: true,
	}
	if r.ID != "M07ABCDEF" || !r.Active {
		t.Errorf("MissionRecord fields: %+v", r)
	}
}

func TestTaskIsComplete(t *testing.T) {
	cases := []struct {
		status *string
		want   bool
	}{
		{nil, false},
		{stringPtr("pending"), false},
		{stringPtr("done"), true},
		{stringPtr("cancelled"), true},
		{stringPtr("in_progress"), false},
	}
	for _, c := range cases {
		if got := TaskIsComplete(c.status); got != c.want {
			t.Errorf("TaskIsComplete(%v) = %v, want %v", c.status, got, c.want)
		}
	}
}

func TestBlockingFindings(t *testing.T) {
	if got := BlockingFindings(nil); got != nil {
		t.Errorf("expected nil for nil review, got %v", got)
	}

	summary := "ok"
	info := "info"
	blocksFalse := false
	r := &Review{
		Findings: []Finding{
			{Summary: &summary, Severity: &info, BlocksShip: &blocksFalse},
		},
	}
	if got := BlockingFindings(r); len(got) != 0 {
		t.Errorf("expected no blocking findings, got %d", len(got))
	}

	blocksTrue := true
	r2 := &Review{
		Findings: []Finding{
			{Summary: &summary, BlocksShip: &blocksTrue},
		},
	}
	if got := BlockingFindings(r2); len(got) != 1 {
		t.Errorf("expected 1 blocking finding (blocksShip), got %d", len(got))
	}

	critical := "critical"
	r3 := &Review{
		Findings: []Finding{
			{Summary: &summary, Severity: &critical},
		},
	}
	if got := BlockingFindings(r3); len(got) != 1 {
		t.Errorf("expected 1 blocking finding (critical), got %d", len(got))
	}
}

func stringPtr(s string) *string { return &s }
