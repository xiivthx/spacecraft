package roadmap

import (
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
		Missions:    []string{"M07A", "M07B"},
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
	if len(got.Missions) != 2 || got.Missions[0] != "M07A" || got.Missions[1] != "M07B" {
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
		Missions:    []string{},
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	if err := s.Create(r); err != nil {
		t.Fatal(err)
	}

	r.Title = "Updated Title"
	r.Missions = []string{"M07X", "M07Y"}
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
	if len(got.Missions) != 2 || got.Missions[0] != "M07X" {
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
			Missions:  []string{},
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
		Missions:  []string{},
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
		Missions:    []string{"A", "B", "C"},
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
	if loaded.Missions[0] != "A" || loaded.Missions[1] != "B" || loaded.Missions[2] != "C" {
		t.Errorf("order mismatch: %v", loaded.Missions)
	}
}
