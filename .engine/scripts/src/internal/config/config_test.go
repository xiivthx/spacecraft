package config

import (
	"path/filepath"
	"testing"
)

func TestNewConfigValid(t *testing.T) {
	cfg, err := NewConfig("/home/user/project")
	if err != nil {
		t.Fatalf("NewConfig should not error for valid root: %v", err)
	}
	if cfg.Root() != "/home/user/project" {
		t.Errorf("Root() = %q, want %q", cfg.Root(), "/home/user/project")
	}
}

func TestNewConfigDerivesPaths(t *testing.T) {
	cfg, err := NewConfig("/repo")
	if err != nil {
		t.Fatal(err)
	}
	wantSpace := "/repo/.space"
	wantMissions := "/repo/.space/missions"
	wantArchive := "/repo/.space/archive"
	wantCurrent := "/repo/.space/current"

	if cfg.SpaceDir() != wantSpace {
		t.Errorf("SpaceDir = %q, want %q", cfg.SpaceDir(), wantSpace)
	}
	if cfg.MissionsDir() != wantMissions {
		t.Errorf("MissionsDir = %q, want %q", cfg.MissionsDir(), wantMissions)
	}
	if cfg.ArchiveDir() != wantArchive {
		t.Errorf("ArchiveDir = %q, want %q", cfg.ArchiveDir(), wantArchive)
	}
	if cfg.CurrentFile() != wantCurrent {
		t.Errorf("CurrentFile = %q, want %q", cfg.CurrentFile(), wantCurrent)
	}
}

func TestNewConfigEmptyRoot(t *testing.T) {
	_, err := NewConfig("")
	if err == nil {
		t.Fatal("NewConfig with empty root should error")
	}
}

func TestNewConfigRelativeRoot(t *testing.T) {
	_, err := NewConfig("relative/path")
	if err == nil {
		t.Fatal("NewConfig with relative root should error")
	}
}

func TestNewConfigWithOptions(t *testing.T) {
	cfg, err := NewConfig("/root", WithSpaceDir("/custom/space"))
	if err != nil {
		t.Fatal(err)
	}
	if cfg.SpaceDir() != "/custom/space" {
		t.Errorf("SpaceDir = %q, want %q", cfg.SpaceDir(), "/custom/space")
	}
	// Other paths should still be derived from root
	wantMissions := filepath.Join("/root", ".space", "missions")
	if cfg.MissionsDir() != wantMissions {
		t.Errorf("MissionsDir = %q, want %q", cfg.MissionsDir(), wantMissions)
	}
	wantArchive := filepath.Join("/root", ".space", "archive")
	if cfg.ArchiveDir() != wantArchive {
		t.Errorf("ArchiveDir = %q, want %q", cfg.ArchiveDir(), wantArchive)
	}
	if cfg.Root() != "/root" {
		t.Errorf("Root = %q, want %q", cfg.Root(), "/root")
	}
}

func TestConfigMissionDir(t *testing.T) {
	cfg, _ := NewConfig("/base")
	md := cfg.MissionDir("M07ABCDEF")
	want := filepath.Join("/base", ".space", "missions", "M07ABCDEF")
	if md != want {
		t.Errorf("MissionDir = %q, want %q", md, want)
	}
}

func TestConfigTraceStoreDir(t *testing.T) {
	cfg, _ := NewConfig("/repo")
	want := "/repo/.space/traces"
	if cfg.TraceStoreDir() != want {
		t.Errorf("TraceStoreDir = %q, want %q", cfg.TraceStoreDir(), want)
	}
}

func TestConfigTraceStoreDirOverride(t *testing.T) {
	cfg, _ := NewConfig("/repo", WithTraceStoreDir("/custom/traces"))
	want := "/custom/traces"
	if cfg.TraceStoreDir() != want {
		t.Errorf("TraceStoreDir = %q, want %q", cfg.TraceStoreDir(), want)
	}
	// Other paths still derived from root
	if cfg.Root() != "/repo" {
		t.Errorf("Root = %q, want %q", cfg.Root(), "/repo")
	}
}

func TestConfigOptionOverrides(t *testing.T) {
	cfg, err := NewConfig("/root",
		WithMissionsDir("/custom/missions"),
		WithArchiveDir("/custom/archive"),
		WithCurrentFile("/custom/current"),
		WithRoadmapsDir("/custom/roadmaps"),
	)
	if err != nil {
		t.Fatal(err)
	}
	if cfg.MissionsDir() != "/custom/missions" {
		t.Errorf("MissionsDir = %q, want %q", cfg.MissionsDir(), "/custom/missions")
	}
	if cfg.ArchiveDir() != "/custom/archive" {
		t.Errorf("ArchiveDir = %q, want %q", cfg.ArchiveDir(), "/custom/archive")
	}
	if cfg.CurrentFile() != "/custom/current" {
		t.Errorf("CurrentFile = %q, want %q", cfg.CurrentFile(), "/custom/current")
	}
	if cfg.RoadmapsDir() != "/custom/roadmaps" {
		t.Errorf("RoadmapsDir = %q, want %q", cfg.RoadmapsDir(), "/custom/roadmaps")
	}
}

func TestConfigRoadmapsDir(t *testing.T) {
	cfg, _ := NewConfig("/base")
	want := filepath.Join("/base", ".space", "roadmaps")
	if cfg.RoadmapsDir() != want {
		t.Errorf("RoadmapsDir = %q, want %q", cfg.RoadmapsDir(), want)
	}
}

func TestConfigEvalsDir(t *testing.T) {
	cfg, _ := NewConfig("/base")
	want := filepath.Join("/base", ".space", "evals")
	if cfg.EvalsDir() != want {
		t.Errorf("EvalsDir = %q, want %q", cfg.EvalsDir(), want)
	}
}

func TestConfigEvalMissionDir(t *testing.T) {
	cfg, _ := NewConfig("/base")
	want := filepath.Join("/base", ".space", "evals", "M123")
	if cfg.EvalMissionDir("M123") != want {
		t.Errorf("EvalMissionDir = %q, want %q", cfg.EvalMissionDir("M123"), want)
	}
}
