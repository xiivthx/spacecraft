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

func TestConfigMissionDir(t *testing.T) {
	cfg, _ := NewConfig("/base")
	md := cfg.MissionDir("M07ABCDEF")
	want := filepath.Join("/base", ".space", "missions", "M07ABCDEF")
	if md != want {
		t.Errorf("MissionDir = %q, want %q", md, want)
	}
}
