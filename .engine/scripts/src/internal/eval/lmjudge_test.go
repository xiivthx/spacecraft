package eval

import (
	"os"
	"testing"

	"spacecraft/internal/mission"
)

func TestLMJudgeFallbackUnconfigured(t *testing.T) {
	os.Unsetenv("SPACECRAFT_JUDGE_MODEL")
	result := RunLMJudge(nil)
	if !result.Fallback {
		t.Error("expected fallback when SPACECRAFT_JUDGE_MODEL is not set")
	}
	if result.FallbackReason == "" {
		t.Error("fallback reason should not be empty")
	}
	if result.Model != "" {
		t.Errorf("model should be empty in fallback, got %q", result.Model)
	}
}

func TestLMJudgeFallbackNoCLI(t *testing.T) {
	savedPath := os.Getenv("PATH")
	os.Setenv("PATH", "/nonexistent-empty-dir")
	os.Setenv("SPACECRAFT_JUDGE_MODEL", "test-model")
	defer func() {
		os.Setenv("PATH", savedPath)
		os.Unsetenv("SPACECRAFT_JUDGE_MODEL")
	}()
	result := RunLMJudge(nil)
	if !result.Fallback {
		t.Error("expected fallback when no judge CLI is available")
	}
}

func TestIsLMJudgeAvailable(t *testing.T) {
	savedPath := os.Getenv("PATH")
	os.Setenv("PATH", "/nonexistent-empty-dir")
	defer os.Setenv("PATH", savedPath)
	if isLMJudgeAvailable() {
		t.Error("should return false when no CLI in PATH")
	}
}

func TestLMJudgeResultStructure(t *testing.T) {
	result := LMJudgeResult{
		Score:    3,
		MaxScore: 4,
		Reasoning: "Good quality output",
		Model:    "test-model",
		Fallback: false,
	}
	if result.Score != 3 || result.MaxScore != 4 {
		t.Errorf("result fields: %+v", result)
	}
	if result.Fallback {
		t.Error("non-fallback result should have Fallback=false")
	}
}

func TestExtractScoreNumeric(t *testing.T) {
	tests := []struct {
		output string
		expect int
	}{
		{`{"score": 4, "reasoning": "perfect"}`, 4},
		{`{"score": 3, "reasoning": "good"}`, 3},
		{`{"score": 2, "reasoning": "ok"}`, 2},
		{`{"score": 1, "reasoning": "bad"}`, 1},
		{`{"score": 0, "reasoning": "terrible"}`, 0},
		{`{"score":4}`, 4},
		{`some text without score`, 0},
	}
	for _, tc := range tests {
		got := extractScore(tc.output)
		if got != tc.expect {
			t.Errorf("extractScore(%q) = %d, want %d", tc.output, got, tc.expect)
		}
	}
}

func TestExtractScoreTextual(t *testing.T) {
	tests := []struct {
		output string
		min int
	}{
		{"excellent output", 4},
		{"good job", 3},
		{"acceptable work", 2},
		{"poor quality", 0},
		{"bad result, failed", 0},
		{"well done", 3},
		{"this is an acceptable and excellent answer", 4},
	}
	for _, tc := range tests {
		got := extractScore(tc.output)
		if got < tc.min {
			t.Errorf("extractScore(%q) = %d, want >= %d", tc.output, got, tc.min)
		}
	}
}

func TestBuildJudgePrompt(t *testing.T) {
	entries := []mission.EvidenceEntry{
		{ID: "e1", Label: "build", Command: "go build ./...", ExitCode: 0, Stdout: "ok"},
		{ID: "e2", Label: "test", Command: "go test ./...", ExitCode: 0, Stdout: "ok"},
	}
	prompt := buildJudgePrompt(entries)
	if prompt == "" {
		t.Error("judge prompt should not be empty")
	}
	if !containsStr(prompt, "go build") {
		t.Error("prompt should contain evidence command")
	}
	if !containsStr(prompt, "go test") {
		t.Error("prompt should contain test command")
	}
	if !containsStr(prompt, "score") {
		t.Error("prompt should ask for score")
	}
}

func TestBuildJudgePromptEmpty(t *testing.T) {
	prompt := buildJudgePrompt(nil)
	if prompt == "" {
		t.Error("judge prompt should not be empty even with no entries")
	}
	if !containsStr(prompt, "no evidence") {
		t.Error("empty prompt should indicate no evidence")
	}
}

func containsStr(s, sub string) bool {
	return len(s) >= len(sub) && subIndex(s, sub)
}

func subIndex(s, sub string) bool {
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
