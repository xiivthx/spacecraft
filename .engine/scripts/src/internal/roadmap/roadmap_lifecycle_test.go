package roadmap

import (
	"testing"
	"time"
)

func makeShipped(ids ...string) func(string) bool {
	m := make(map[string]bool)
	for _, id := range ids {
		m[id] = true
	}
	return func(id string) bool {
		return m[id]
	}
}

func TestDeriveStateActive(t *testing.T) {
	r := &Roadmap{
		ID:       "M07L01",
		Title:    "Test",
		Missions: []MissionEntry{{ID: "A"}, {ID: "B"}, {ID: "C"}},
	}
	state := DeriveState(r, makeShipped("A"))
	if state != "active" {
		t.Errorf("state = %q, want %q", state, "active")
	}
}

func TestDeriveStateDone(t *testing.T) {
	r := &Roadmap{
		ID:       "M07L02",
		Title:    "Test",
		Missions: []MissionEntry{{ID: "A"}, {ID: "B"}},
	}
	state := DeriveState(r, makeShipped("A", "B"))
	if state != "done" {
		t.Errorf("state = %q, want %q", state, "done")
	}
}

func TestDeriveStateEmpty(t *testing.T) {
	r := &Roadmap{
		ID:       "M07L03",
		Title:    "Empty",
		Missions: []MissionEntry{},
	}
	state := DeriveState(r, makeShipped())
	if state != "active" {
		t.Errorf("state = %q, want %q", state, "active")
	}
}

func TestSaveRoundtripPreservesMissions(t *testing.T) {
	cfg, cleanup := newTestConfig(t)
	defer cleanup()
	s := NewFSStore(cfg)

	r := &Roadmap{
		ID:        "M07L04",
		Title:     "Test",
		Missions:  []MissionEntry{{ID: "X"}, {ID: "Y"}, {ID: "Z"}},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	if err := s.Create(r); err != nil {
		t.Fatal(err)
	}

	r.Missions = []MissionEntry{{ID: "X"}, {ID: "Y"}, {ID: "Z"}, {ID: "W"}}
	r.UpdatedAt = time.Now()
	if err := s.Save(r); err != nil {
		t.Fatal(err)
	}

	got, err := s.Load("M07L04")
	if err != nil {
		t.Fatal(err)
	}
	if len(got.Missions) != 4 {
		t.Errorf("missions count = %d, want 4", len(got.Missions))
	}
}
