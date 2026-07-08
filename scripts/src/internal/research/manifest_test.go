package research

import (
	"os"
	"path/filepath"
	"testing"
)

// TestDependencyStruct verifies Dependency struct literal and field access.
func TestDependencyStruct(t *testing.T) {
	d := Dependency{
		Name:           "test-pkg",
		CurrentVersion: "v1.0.0",
		Ecosystem:      "go",
	}
	if d.Name != "test-pkg" {
		t.Errorf("Dependency.Name = %q, want %q", d.Name, "test-pkg")
	}
	if d.CurrentVersion != "v1.0.0" {
		t.Errorf("Dependency.CurrentVersion = %q, want %q", d.CurrentVersion, "v1.0.0")
	}
	if d.Ecosystem != "go" {
		t.Errorf("Dependency.Ecosystem = %q, want %q", d.Ecosystem, "go")
	}
}

func TestParseGoMod(t *testing.T) {
	content := `module example.com/myapp

go 1.26

require (
	github.com/charmbracelet/bubbletea v1.3.0
	github.com/mattn/go-isatty v0.0.20
	golang.org/x/sys v0.30.0 // indirect
)
`
	dir := t.TempDir()
	path := filepath.Join(dir, "go.mod")
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	deps, err := ParseManifest(path)
	if err != nil {
		t.Fatalf("ParseManifest(%q) unexpected error: %v", path, err)
	}

	if len(deps) != 2 {
		t.Fatalf("got %d deps, want 2", len(deps))
	}

	checks := []struct {
		name    string
		version string
	}{
		{"github.com/charmbracelet/bubbletea", "v1.3.0"},
		{"github.com/mattn/go-isatty", "v0.0.20"},
	}
	for i, c := range checks {
		if deps[i].Name != c.name {
			t.Errorf("deps[%d].Name = %q, want %q", i, deps[i].Name, c.name)
		}
		if deps[i].CurrentVersion != c.version {
			t.Errorf("deps[%d].CurrentVersion = %q, want %q", i, deps[i].CurrentVersion, c.version)
		}
		if deps[i].Ecosystem != "go" {
			t.Errorf("deps[%d].Ecosystem = %q, want %q", i, deps[i].Ecosystem, "go")
		}
	}
}

func TestParsePackageJSON(t *testing.T) {
	content := `{
	"name": "my-app",
	"dependencies": {
		"express": "^4.18.0",
		"lodash": "~4.17.21"
	},
	"devDependencies": {
		"jest": "29.7.0",
		"@types/node": "^20.0.0"
	}
}`
	dir := t.TempDir()
	path := filepath.Join(dir, "package.json")
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	deps, err := ParseManifest(path)
	if err != nil {
		t.Fatalf("ParseManifest(%q) unexpected error: %v", path, err)
	}

	if len(deps) != 4 {
		t.Fatalf("got %d deps, want 4", len(deps))
	}

	checks := []struct {
		name    string
		version string
	}{
		{"express", "^4.18.0"},
		{"lodash", "~4.17.21"},
		{"jest", "29.7.0"},
		{"@types/node", "^20.0.0"},
	}
	for i, c := range checks {
		if deps[i].Name != c.name {
			t.Errorf("deps[%d].Name = %q, want %q", i, deps[i].Name, c.name)
		}
		if deps[i].CurrentVersion != c.version {
			t.Errorf("deps[%d].CurrentVersion = %q, want %q", i, deps[i].CurrentVersion, c.version)
		}
		if deps[i].Ecosystem != "npm" {
			t.Errorf("deps[%d].Ecosystem = %q, want %q", i, deps[i].Ecosystem, "npm")
		}
	}
}

func TestParseRequirementsTxt(t *testing.T) {
	content := `requests==2.32.0
flask>=2.0.0
# this is a comment
numpy==1.26.0

pandas>=2.1.0
`
	dir := t.TempDir()
	path := filepath.Join(dir, "requirements.txt")
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	deps, err := ParseManifest(path)
	if err != nil {
		t.Fatalf("ParseManifest(%q) unexpected error: %v", path, err)
	}

	if len(deps) != 4 {
		t.Fatalf("got %d deps, want 4", len(deps))
	}

	checks := []struct {
		name    string
		version string
	}{
		{"requests", "2.32.0"},
		{"flask", "2.0.0"},
		{"numpy", "1.26.0"},
		{"pandas", "2.1.0"},
	}
	for i, c := range checks {
		if deps[i].Name != c.name {
			t.Errorf("deps[%d].Name = %q, want %q", i, deps[i].Name, c.name)
		}
		if deps[i].CurrentVersion != c.version {
			t.Errorf("deps[%d].CurrentVersion = %q, want %q", i, deps[i].CurrentVersion, c.version)
		}
		if deps[i].Ecosystem != "pypi" {
			t.Errorf("deps[%d].Ecosystem = %q, want %q", i, deps[i].Ecosystem, "pypi")
		}
	}
}

func TestParsePyprojectToml(t *testing.T) {
	content := `[build-system]
requires = ["setuptools>=64"]
build-backend = "setuptools.backends._legacy:_Backend"

[project]
name = "my-package"
version = "1.0.0"
dependencies = [
	"requests>=2.32.0",
	"flask>=2.0.0",
	"numpy==1.26.0",
]
`
	dir := t.TempDir()
	path := filepath.Join(dir, "pyproject.toml")
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	deps, err := ParseManifest(path)
	if err != nil {
		t.Fatalf("ParseManifest(%q) unexpected error: %v", path, err)
	}

	if len(deps) != 3 {
		t.Fatalf("got %d deps, want 3", len(deps))
	}

	checks := []struct {
		name    string
		version string
	}{
		{"requests", "2.32.0"},
		{"flask", "2.0.0"},
		{"numpy", "1.26.0"},
	}
	for i, c := range checks {
		if deps[i].Name != c.name {
			t.Errorf("deps[%d].Name = %q, want %q", i, deps[i].Name, c.name)
		}
		if deps[i].CurrentVersion != c.version {
			t.Errorf("deps[%d].CurrentVersion = %q, want %q", i, deps[i].CurrentVersion, c.version)
		}
		if deps[i].Ecosystem != "pypi" {
			t.Errorf("deps[%d].Ecosystem = %q, want %q", i, deps[i].Ecosystem, "pypi")
		}
	}
}

func TestParseCargoToml(t *testing.T) {
	content := `[package]
name = "my-crate"
version = "0.1.0"

[dependencies]
serde = "1.0.200"
tokio = { version = "1.0", features = ["full"] }
regex = "1.10.0"
`
	dir := t.TempDir()
	path := filepath.Join(dir, "Cargo.toml")
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	deps, err := ParseManifest(path)
	if err != nil {
		t.Fatalf("ParseManifest(%q) unexpected error: %v", path, err)
	}

	if len(deps) != 3 {
		t.Fatalf("got %d deps, want 3", len(deps))
	}

	checks := []struct {
		name    string
		version string
	}{
		{"serde", "1.0.200"},
		{"tokio", "1.0"},
		{"regex", "1.10.0"},
	}
	for i, c := range checks {
		if deps[i].Name != c.name {
			t.Errorf("deps[%d].Name = %q, want %q", i, deps[i].Name, c.name)
		}
		if deps[i].CurrentVersion != c.version {
			t.Errorf("deps[%d].CurrentVersion = %q, want %q", i, deps[i].CurrentVersion, c.version)
		}
		if deps[i].Ecosystem != "crates" {
			t.Errorf("deps[%d].Ecosystem = %q, want %q", i, deps[i].Ecosystem, "crates")
		}
	}
}

func TestParseManifestUnknown(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "unknown.lock")

	content := []byte("some random content\n")
	if err := os.WriteFile(path, content, 0644); err != nil {
		t.Fatal(err)
	}

	_, err := ParseManifest(path)
	if err == nil {
		t.Error("expected error for unknown manifest type, got nil")
	}
}

// TestParseManifestNonExistent verifies error for a non-existent file.
func TestParseManifestNonExistent(t *testing.T) {
	path := "/nonexistent/path/go.mod"
	_, err := ParseManifest(path)
	if err == nil {
		t.Error("expected error for non-existent file, got nil")
	}
}
