package roadmap

import (
	"encoding/json"
	"os"
	"testing"
	"time"

	"spacecraft/internal/config"
)

func newTestConfig(t *testing.T) (*config.Config, func()) {
	t.Helper()
	dir, err := os.MkdirTemp("", "spacecraft-roadmap-test-")
	if err != nil {
		t.Fatal(err)
	}
	cfg, err := config.NewConfig(dir)
	if err != nil {
		os.RemoveAll(dir)
		t.Fatal(err)
	}
	cleanup := func() { os.RemoveAll(dir) }
	return cfg, cleanup
}

func TestCreateAndLoad(t *testing.T) {
	cfg, cleanup := newTestConfig(t)
	defer cleanup()
	s := NewFSStore(cfg)

	r := &Roadmap{
		ID:          "M07TEST01",
		Title:       "Test Road",
		Description: "a test",
		Missions:    []MissionEntry{{ID: "M07A"}, {ID: "M07B"}},
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	if err := s.Create(r); err != nil {
		t.Fatal(err)
	}

	got, err := s.Load("M07TEST01")
	if err != nil {
		t.Fatal(err)
	}
	if got.ID != r.ID {
		t.Errorf("ID = %q, want %q", got.ID, r.ID)
	}
	if got.Title != r.Title {
		t.Errorf("Title = %q, want %q", got.Title, r.Title)
	}
	if got.Description != r.Description {
		t.Errorf("Description = %q, want %q", got.Description, r.Description)
	}
	if len(got.Missions) != 2 || got.Missions[0].ID != "M07A" || got.Missions[1].ID != "M07B" {
		t.Errorf("Missions = %v, want [M07A M07B]", got.Missions)
	}
}

func TestLoadNonExistent(t *testing.T) {
	cfg, cleanup := newTestConfig(t)
	defer cleanup()
	s := NewFSStore(cfg)

	_, err := s.Load("M07NOPE")
	if err == nil {
		t.Fatal("expected error for non-existent roadmap")
	}
}

func TestSaveAndLoadRoundtrip(t *testing.T) {
	cfg, cleanup := newTestConfig(t)
	defer cleanup()
	s := NewFSStore(cfg)

	r := &Roadmap{
		ID:          "M07TEST02",
		Title:       "Roundtrip",
		Missions:    []MissionEntry{},
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	if err := s.Create(r); err != nil {
		t.Fatal(err)
	}

	r.Title = "Updated Title"
	r.Missions = []MissionEntry{{ID: "M07X"}, {ID: "M07Y"}}
	r.UpdatedAt = time.Now()
	if err := s.Save(r); err != nil {
		t.Fatal(err)
	}

	got, err := s.Load("M07TEST02")
	if err != nil {
		t.Fatal(err)
	}
	if got.Title != "Updated Title" {
		t.Errorf("Title = %q, want %q", got.Title, "Updated Title")
	}
	if len(got.Missions) != 2 || got.Missions[0].ID != "M07X" {
		t.Errorf("Missions = %v, want [M07X M07Y]", got.Missions)
	}
}

func TestListEmpty(t *testing.T) {
	cfg, cleanup := newTestConfig(t)
	defer cleanup()
	s := NewFSStore(cfg)

	rms, err := s.List()
	if err != nil {
		t.Fatal(err)
	}
	if len(rms) != 0 {
		t.Errorf("expected 0 roadmaps, got %d", len(rms))
	}
}

func TestListMultiple(t *testing.T) {
	cfg, cleanup := newTestConfig(t)
	defer cleanup()
	s := NewFSStore(cfg)

	for _, id := range []string{"M07A01", "M07A02", "M07A03"} {
		if err := s.Create(&Roadmap{
			ID:        id,
			Title:     "R " + id,
			Missions:  []MissionEntry{},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}); err != nil {
			t.Fatal(err)
		}
	}

	rms, err := s.List()
	if err != nil {
		t.Fatal(err)
	}
	if len(rms) != 3 {
		t.Errorf("expected 3 roadmaps, got %d", len(rms))
	}
}

func TestDelete(t *testing.T) {
	cfg, cleanup := newTestConfig(t)
	defer cleanup()
	s := NewFSStore(cfg)

	r := &Roadmap{
		ID:        "M07DEL01",
		Title:     "To Delete",
		Missions:  []MissionEntry{},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	if err := s.Create(r); err != nil {
		t.Fatal(err)
	}
	if err := s.Delete("M07DEL01"); err != nil {
		t.Fatal(err)
	}
	_, err := s.Load("M07DEL01")
	if err == nil {
		t.Fatal("expected error after delete")
	}
}

func TestDeleteNonExistent(t *testing.T) {
	cfg, cleanup := newTestConfig(t)
	defer cleanup()
	s := NewFSStore(cfg)

	err := s.Delete("M07NOPE")
	if err == nil {
		t.Fatal("expected error for non-existent roadmap")
	}
}

func TestInsertOrdering(t *testing.T) {
	cfg, cleanup := newTestConfig(t)
	defer cleanup()
	s := NewFSStore(cfg)

	r := &Roadmap{
		ID:          "M07ORD01",
		Title:       "Ordered",
		Missions:    []MissionEntry{{ID: "A"}, {ID: "B"}, {ID: "C"}},
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	if err := s.Create(r); err != nil {
		t.Fatal(err)
	}

	loaded, err := s.Load("M07ORD01")
	if err != nil {
		t.Fatal(err)
	}
	if len(loaded.Missions) != 3 {
		t.Fatalf("expected 3 missions, got %d", len(loaded.Missions))
	}
	if loaded.Missions[0].ID != "A" || loaded.Missions[1].ID != "B" || loaded.Missions[2].ID != "C" {
		t.Errorf("order mismatch: %v", loaded.Missions)
	}
}

func TestIssueJSONRoundtrip(t *testing.T) {
	issue := Issue{
		Number: 27,
		Title:  "Test issue",
		URL:    "https://github.com/test/27",
		State:  "open",
		Labels: []string{"bug", "priority"},
		Phase:  "phase-1",
	}

	data, err := json.Marshal(issue)
	if err != nil {
		t.Fatal(err)
	}

	var got Issue
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatal(err)
	}

	if got.Number != issue.Number {
		t.Errorf("Number = %d, want %d", got.Number, issue.Number)
	}
	if got.Title != issue.Title {
		t.Errorf("Title = %q, want %q", got.Title, issue.Title)
	}
	if got.URL != issue.URL {
		t.Errorf("URL = %q, want %q", got.URL, issue.URL)
	}
	if got.State != issue.State {
		t.Errorf("State = %q, want %q", got.State, issue.State)
	}
	if len(got.Labels) != 2 || got.Labels[0] != "bug" {
		t.Errorf("Labels = %v, want [bug priority]", got.Labels)
	}
	if got.Phase != issue.Phase {
		t.Errorf("Phase = %q, want %q", got.Phase, issue.Phase)
	}
}

func TestRoadmapWithIssues(t *testing.T) {
	r := Roadmap{
		ID:          "M07ISSUES",
		Title:       "Roadmap with issues",
		Description: "testing issue serialization",
		Missions:    []MissionEntry{{ID: "M07A"}},
		Issues: []Issue{
			{Number: 1, Title: "First issue", URL: "https://github.com/test/1", State: "open", Labels: []string{"bug"}, Phase: "phase-1"},
			{Number: 2, Title: "Second issue", URL: "https://github.com/test/2", State: "closed", Phase: "phase-2"},
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	data, err := json.Marshal(r)
	if err != nil {
		t.Fatal(err)
	}

	var got Roadmap
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatal(err)
	}

	if len(got.Issues) != 2 {
		t.Fatalf("Issues length = %d, want 2", len(got.Issues))
	}
	if got.Issues[0].Number != 1 || got.Issues[0].Title != "First issue" || got.Issues[0].Phase != "phase-1" {
		t.Errorf("Issues[0] = %+v, want first issue with phase-1", got.Issues[0])
	}
	if got.Issues[1].Number != 2 || got.Issues[1].State != "closed" || got.Issues[1].Phase != "phase-2" {
		t.Errorf("Issues[1] = %+v, want second closed issue with phase-2", got.Issues[1])
	}
}

func TestBackwardCompat(t *testing.T) {
	oldJSON := `{"id":"M07OLD","title":"Old roadmap","description":"legacy","missions":["M07A"],"createdAt":"2026-01-01T00:00:00Z","updatedAt":"2026-01-01T00:00:00Z"}`

	var got Roadmap
	if err := json.Unmarshal([]byte(oldJSON), &got); err != nil {
		t.Fatalf("failed to unmarshal legacy roadmap: %v", err)
	}

	if got.ID != "M07OLD" {
		t.Errorf("ID = %q, want %q", got.ID, "M07OLD")
	}
	if got.Title != "Old roadmap" {
		t.Errorf("Title = %q, want %q", got.Title, "Old roadmap")
	}
	if len(got.Missions) != 1 || got.Missions[0].ID != "M07A" {
		t.Errorf("Missions = %v, want [M07A]", got.Missions)
	}
	if len(got.Issues) != 0 {
		t.Errorf("Issues = %v, want nil or empty", got.Issues)
	}
}
