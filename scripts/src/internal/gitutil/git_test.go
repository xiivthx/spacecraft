package gitutil

import (
	"testing"
)

// mockRunner simulates git responses for testing.
type mockRunner struct {
	responses map[string]mockResponse
}

type mockResponse struct {
	exitCode int
	stdout   string
	stderr   string
}

func (m mockRunner) Run(name string, args ...string) (int, string, string) {
	key := name
	for _, a := range args {
		key += " " + a
	}
	if resp, ok := m.responses[key]; ok {
		return resp.exitCode, resp.stdout, resp.stderr
	}
	return 127, "", "command not found"
}

func TestGitInfoNotRepo(t *testing.T) {
	m := mockRunner{responses: map[string]mockResponse{
		"git rev-parse --is-inside-work-tree": {exitCode: 128, stdout: "false\n"},
	}}
	info := GitInfo(m)
	if info.IsRepo {
		t.Error("should not be a repo")
	}
	if !info.Available {
		t.Error("git should be available (exit code 128 is not 127)")
	}
}

func TestGitInfoGitUnavailable(t *testing.T) {
	m := mockRunner{responses: map[string]mockResponse{
		"git rev-parse --is-inside-work-tree": {exitCode: 127},
	}}
	info := GitInfo(m)
	if info.IsRepo {
		t.Error("should not be a repo")
	}
	if info.Available {
		t.Error("git should be unavailable with exit code 127")
	}
}

func TestGitInfoRepo(t *testing.T) {
	m := mockRunner{responses: map[string]mockResponse{
		"git rev-parse --is-inside-work-tree": {exitCode: 0, stdout: "true\n"},
		"git rev-parse --show-toplevel":       {exitCode: 0, stdout: "/repo\n"},
		"git branch --show-current":           {exitCode: 0, stdout: "main\n"},
		"git rev-parse HEAD":                  {exitCode: 0, stdout: "abc123\n"},
		"git status --short":                  {exitCode: 0, stdout: ""},
	}}
	info := GitInfo(m)
	if !info.IsRepo {
		t.Error("should be a repo")
	}
	if !info.Available {
		t.Error("git should be available")
	}
	if info.Root != "/repo" {
		t.Errorf("Root = %q, want /repo", info.Root)
	}
	if info.Branch != "main" {
		t.Errorf("Branch = %q, want main", info.Branch)
	}
	if info.Sha != "abc123" {
		t.Errorf("Sha = %q, want abc123", info.Sha)
	}
	if info.Dirty {
		t.Error("should be clean")
	}
	if info.DirtyFiles != 0 {
		t.Errorf("DirtyFiles = %d, want 0", info.DirtyFiles)
	}
}

func TestGitInfoDirty(t *testing.T) {
	m := mockRunner{responses: map[string]mockResponse{
		"git rev-parse --is-inside-work-tree": {exitCode: 0, stdout: "true\n"},
		"git rev-parse --show-toplevel":       {exitCode: 0, stdout: "/repo\n"},
		"git branch --show-current":           {exitCode: 0, stdout: "feat\n"},
		"git rev-parse HEAD":                  {exitCode: 0, stdout: "def456\n"},
		"git status --short":                  {exitCode: 0, stdout: " M file.go\n"},
	}}
	info := GitInfo(m)
	if !info.Dirty {
		t.Error("should be dirty")
	}
	if info.DirtyFiles != 1 {
		t.Errorf("DirtyFiles = %d, want 1", info.DirtyFiles)
	}
}

func TestGitInfoMultipleDirtyFiles(t *testing.T) {
	m := mockRunner{responses: map[string]mockResponse{
		"git rev-parse --is-inside-work-tree": {exitCode: 0, stdout: "true\n"},
		"git rev-parse --show-toplevel":       {exitCode: 0, stdout: "/repo\n"},
		"git branch --show-current":           {exitCode: 0, stdout: "feat\n"},
		"git rev-parse HEAD":                  {exitCode: 0, stdout: "def456\n"},
		"git status --short":                  {exitCode: 0, stdout: " M a.go\n M b.go\n"},
	}}
	info := GitInfo(m)
	if !info.Dirty {
		t.Error("should be dirty")
	}
	if info.DirtyFiles != 2 {
		t.Errorf("DirtyFiles = %d, want 2", info.DirtyFiles)
	}
}
