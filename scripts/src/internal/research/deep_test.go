package research

import (
	"context"
	"errors"
	"testing"
)

// mockExecutor provides a controllable implementation of the Executor interface
// for testing runner availability and install instructions without real subprocess calls.
type mockExecutor struct {
	lookPathFn func(file string) (string, error)
	runFn      func(ctx context.Context, name string, args ...string) ([]byte, error)
}

func (m *mockExecutor) LookPath(file string) (string, error) {
	return m.lookPathFn(file)
}

func (m *mockExecutor) Run(ctx context.Context, name string, args ...string) ([]byte, error) {
	return m.runFn(ctx, name, args...)
}

// ---------------------------------------------------------------------------
// DeepResult struct
// ---------------------------------------------------------------------------

// TestDeepResultStruct verifies that the DeepResult struct exists with the
// expected fields: Summary, KeyPoints, SourceURL, FetchedAt.
func TestDeepResultStruct(t *testing.T) {
	r := DeepResult{
		Summary:   "test summary",
		KeyPoints: []string{"point 1", "point 2"},
		SourceURL: "https://example.com",
		FetchedAt: "2026-07-08T12:00:00Z",
	}
	if r.Summary != "test summary" {
		t.Errorf("DeepResult.Summary = %q, want %q", r.Summary, "test summary")
	}
	if len(r.KeyPoints) != 2 || r.KeyPoints[0] != "point 1" || r.KeyPoints[1] != "point 2" {
		t.Errorf("DeepResult.KeyPoints = %v, want %v", r.KeyPoints, []string{"point 1", "point 2"})
	}
	if r.SourceURL != "https://example.com" {
		t.Errorf("DeepResult.SourceURL = %q, want %q", r.SourceURL, "https://example.com")
	}
	if r.FetchedAt != "2026-07-08T12:00:00Z" {
		t.Errorf("DeepResult.FetchedAt = %q, want %q", r.FetchedAt, "2026-07-08T12:00:00Z")
	}
}

// ---------------------------------------------------------------------------
// BrowserUseRunner
// ---------------------------------------------------------------------------

// TestBrowserUseRunnerNotAvailable verifies that IsAvailable returns false
// when the executor cannot find python3 in PATH.
func TestBrowserUseRunnerNotAvailable(t *testing.T) {
	exec := &mockExecutor{
		lookPathFn: func(file string) (string, error) {
			if file == "python3" {
				return "", errors.New("not found")
			}
			return "", errors.New("not found")
		},
		runFn: func(ctx context.Context, name string, args ...string) ([]byte, error) {
			return nil, errors.New("not executed")
		},
	}
	runner := NewBrowserUseRunner(exec)
	if runner.IsAvailable() {
		t.Error("BrowserUseRunner.IsAvailable() = true, want false")
	}
}

// TestBrowserUseRunnerAvailable verifies that IsAvailable returns true when
// the executor finds python3 in PATH.
func TestBrowserUseRunnerAvailable(t *testing.T) {
	exec := &mockExecutor{
		lookPathFn: func(file string) (string, error) {
			if file == "python3" {
				return "/usr/bin/python3", nil
			}
			return "", errors.New("not found")
		},
		runFn: func(ctx context.Context, name string, args ...string) ([]byte, error) {
			return nil, errors.New("not executed")
		},
	}
	runner := NewBrowserUseRunner(exec)
	if !runner.IsAvailable() {
		t.Error("BrowserUseRunner.IsAvailable() = false, want true")
	}
}

// TestBrowserUseRunnerInstallInstructions verifies that InstallInstructions
// contains the expected install command for browser-use.
func TestBrowserUseRunnerInstallInstructions(t *testing.T) {
	exec := &mockExecutor{
		lookPathFn: func(file string) (string, error) {
			return "", errors.New("not found")
		},
		runFn: func(ctx context.Context, name string, args ...string) ([]byte, error) {
			return nil, errors.New("not executed")
		},
	}
	runner := NewBrowserUseRunner(exec)
	instructions := runner.InstallInstructions()
	if instructions == "" {
		t.Fatal("BrowserUseRunner.InstallInstructions() returned empty string")
	}
	if !contains(instructions, "pip install browser-use") {
		t.Errorf("BrowserUseRunner.InstallInstructions() = %q, want it to contain %q",
			instructions, "pip install browser-use")
	}
}

// TestBrowserUseRunnerAnalyzeNotAvailable verifies that Analyze returns an
// error containing install instructions when the tool is not available.
func TestBrowserUseRunnerAnalyzeNotAvailable(t *testing.T) {
	exec := &mockExecutor{
		lookPathFn: func(file string) (string, error) {
			return "", errors.New("not found")
		},
		runFn: func(ctx context.Context, name string, args ...string) ([]byte, error) {
			return nil, errors.New("not executed")
		},
	}
	runner := NewBrowserUseRunner(exec)
	result, err := runner.Analyze(context.Background(), "https://example.com")
	if err == nil {
		t.Fatal("BrowserUseRunner.Analyze() should return error when tool not available, got nil")
	}
	if result != nil {
		t.Errorf("BrowserUseRunner.Analyze() result = %v, want nil when tool not available", result)
	}
	if !contains(err.Error(), "pip install browser-use") {
		t.Errorf("error message should contain install instructions, got: %v", err)
	}
}

// ---------------------------------------------------------------------------
// NotebookLMRunner
// ---------------------------------------------------------------------------

// TestNotebookLMRunnerNotAvailable verifies that IsAvailable returns false
// when the executor cannot find nlm in PATH.
func TestNotebookLMRunnerNotAvailable(t *testing.T) {
	exec := &mockExecutor{
		lookPathFn: func(file string) (string, error) {
			if file == "nlm" {
				return "", errors.New("not found")
			}
			return "", errors.New("not found")
		},
		runFn: func(ctx context.Context, name string, args ...string) ([]byte, error) {
			return nil, errors.New("not executed")
		},
	}
	runner := NewNotebookLMRunner(exec)
	if runner.IsAvailable() {
		t.Error("NotebookLMRunner.IsAvailable() = true, want false")
	}
}

// TestNotebookLMRunnerAvailable verifies that IsAvailable returns true when
// the executor finds nlm in PATH.
func TestNotebookLMRunnerAvailable(t *testing.T) {
	exec := &mockExecutor{
		lookPathFn: func(file string) (string, error) {
			if file == "nlm" {
				return "/usr/local/bin/nlm", nil
			}
			return "", errors.New("not found")
		},
		runFn: func(ctx context.Context, name string, args ...string) ([]byte, error) {
			return nil, errors.New("not executed")
		},
	}
	runner := NewNotebookLMRunner(exec)
	if !runner.IsAvailable() {
		t.Error("NotebookLMRunner.IsAvailable() = false, want true")
	}
}

// TestNotebookLMRunnerInstallInstructions verifies that InstallInstructions
// contains the expected install command for notebooklm-mcp-cli.
func TestNotebookLMRunnerInstallInstructions(t *testing.T) {
	exec := &mockExecutor{
		lookPathFn: func(file string) (string, error) {
			return "", errors.New("not found")
		},
		runFn: func(ctx context.Context, name string, args ...string) ([]byte, error) {
			return nil, errors.New("not executed")
		},
	}
	runner := NewNotebookLMRunner(exec)
	instructions := runner.InstallInstructions()
	if instructions == "" {
		t.Fatal("NotebookLMRunner.InstallInstructions() returned empty string")
	}
	if !contains(instructions, "notebooklm-mcp-cli") {
		t.Errorf("NotebookLMRunner.InstallInstructions() = %q, want it to contain %q",
			instructions, "notebooklm-mcp-cli")
	}
}

// TestNotebookLMRunnerAnalyzeNotAvailable verifies that Analyze returns an
// error containing install instructions when the tool is not available.
func TestNotebookLMRunnerAnalyzeNotAvailable(t *testing.T) {
	exec := &mockExecutor{
		lookPathFn: func(file string) (string, error) {
			return "", errors.New("not found")
		},
		runFn: func(ctx context.Context, name string, args ...string) ([]byte, error) {
			return nil, errors.New("not executed")
		},
	}
	runner := NewNotebookLMRunner(exec)
	result, err := runner.Analyze(context.Background(), "test query")
	if err == nil {
		t.Fatal("NotebookLMRunner.Analyze() should return error when tool not available, got nil")
	}
	if result != nil {
		t.Errorf("NotebookLMRunner.Analyze() result = %v, want nil when tool not available", result)
	}
	if !contains(err.Error(), "notebooklm-mcp-cli") {
		t.Errorf("error message should contain install instructions, got: %v", err)
	}
}

// contains is a small helper to check substring presence without importing strings.
func contains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
