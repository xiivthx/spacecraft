package compact

import (
	"os"
	"testing"
)

func TestLoadDSLFilter_MissingFile(t *testing.T) {
	filter, err := LoadDSLFilter("/nonexistent/path/to/config.yaml", &CommandInfo{Exe: "git", Arg1: "status"})
	if filter != nil {
		t.Errorf("expected nil filter, got %v", filter)
	}
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
}

func TestLoadDSLFilter_InvalidJSON(t *testing.T) {
	f, err := os.CreateTemp(t.TempDir(), "invalid*.json")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	if _, err := f.WriteString("{bad json {{{"); err != nil {
		t.Fatal(err)
	}

	filter, err := LoadDSLFilter(f.Name(), &CommandInfo{Exe: "git", Arg1: "status"})
	if filter != nil {
		t.Errorf("expected nil filter for invalid JSON, got %v", filter)
	}
	if err == nil {
		t.Error("expected non-nil error for invalid JSON config, got nil")
	}
}

func TestLoadDSLFilter_IncludeStage(t *testing.T) {
	configJSON := `{"rules":[{"exe":"go","stages":[{"include":{"pattern":"error"}}]}]}`

	f, err := os.CreateTemp(t.TempDir(), "include*.json")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	if _, err := f.WriteString(configJSON); err != nil {
		t.Fatal(err)
	}

	ci := &CommandInfo{Exe: "go"}
	filter, err := LoadDSLFilter(f.Name(), ci)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if filter == nil {
		t.Fatal("expected non-nil filter for matching rule, got nil")
	}

	input := "ok: package\nwarning: error found\nall good\nanother error here"
	expected := "warning: error found\nanother error here"
	got := filter.Apply(input)

	if got != expected {
		t.Fatalf("IncludeStage.Apply mismatch:\n  got:  %q\n  want: %q", got, expected)
	}
}

func TestLoadDSLFilter_ExcludeStage(t *testing.T) {
	configJSON := `{"rules":[{"exe":"go","stages":[{"exclude":{"pattern":"^#"}}]}]}`

	f, err := os.CreateTemp(t.TempDir(), "exclude*.json")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	if _, err := f.WriteString(configJSON); err != nil {
		t.Fatal(err)
	}

	ci := &CommandInfo{Exe: "go"}
	filter, err := LoadDSLFilter(f.Name(), ci)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if filter == nil {
		t.Fatal("expected non-nil filter for matching rule, got nil")
	}

	input := "keep this\n# drop comment\ncode here\n# another comment\nmore code"
	expected := "keep this\ncode here\nmore code"
	got := filter.Apply(input)

	if got != expected {
		t.Fatalf("ExcludeStage.Apply mismatch:\n  got:  %q\n  want: %q", got, expected)
	}
}

func TestLoadDSLFilter_DedupStage(t *testing.T) {
	configJSON := `{"rules":[{"exe":"echo","stages":[{"dedup":{}}]}]}`

	f, err := os.CreateTemp(t.TempDir(), "dedup*.json")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	if _, err := f.WriteString(configJSON); err != nil {
		t.Fatal(err)
	}

	ci := &CommandInfo{Exe: "echo"}
	filter, err := LoadDSLFilter(f.Name(), ci)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if filter == nil {
		t.Fatal("expected non-nil filter for matching rule, got nil")
	}

	input := "hello\nhello\nworld\nworld\nworld\ntest"
	expected := "hello [x2]\nworld [x3]\ntest"
	got := filter.Apply(input)

	if got != expected {
		t.Fatalf("DedupStage.Apply mismatch:\n  got:  %q\n  want: %q", got, expected)
	}
}

func TestLoadDSLFilter_TruncateStage(t *testing.T) {
	configJSON := `{"rules":[{"exe":"ls","stages":[{"truncate":{"head":2,"tail":1}}]}]}`

	f, err := os.CreateTemp(t.TempDir(), "truncate*.json")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	if _, err := f.WriteString(configJSON); err != nil {
		t.Fatal(err)
	}

	ci := &CommandInfo{Exe: "ls"}
	filter, err := LoadDSLFilter(f.Name(), ci)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if filter == nil {
		t.Fatal("expected non-nil filter for matching rule, got nil")
	}

	input := "line1\nline2\nline3\nline4\nline5\nline6"
	expected := "line1\nline2\n--- 3 lines skipped (total: 6) ---\nline6"
	got := filter.Apply(input)

	if got != expected {
		t.Fatalf("TruncateStage.Apply mismatch:\n  got:  %q\n  want: %q", got, expected)
	}
}

func TestLoadDSLFilter_StripPrefixStage(t *testing.T) {
	configJSON := `{"rules":[{"exe":"find","stages":[{"stripPrefix":{"prefix":"/home/user/"}}]}]}`

	f, err := os.CreateTemp(t.TempDir(), "stripprefix*.json")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	if _, err := f.WriteString(configJSON); err != nil {
		t.Fatal(err)
	}

	ci := &CommandInfo{Exe: "find"}
	filter, err := LoadDSLFilter(f.Name(), ci)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if filter == nil {
		t.Fatal("expected non-nil filter for matching rule, got nil")
	}

	input := "/home/user/src/main.go\n/home/user/docs/readme.md"
	expected := "src/main.go\ndocs/readme.md"
	got := filter.Apply(input)

	if got != expected {
		t.Fatalf("StripPrefixStage.Apply mismatch:\n  got:  %q\n  want: %q", got, expected)
	}
}

func TestLoadDSLFilter_FirstMatchWins(t *testing.T) {
	// Two rules match exe:"go". The first (include "aaa") should win over the second (passthrough).
	configJSON := `{"rules":[{"exe":"go","stages":[{"include":{"pattern":"aaa"}}]},{"exe":"go","stages":[{"passthrough":{}}]}]}`

	f, err := os.CreateTemp(t.TempDir(), "firstmatch*.json")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	if _, err := f.WriteString(configJSON); err != nil {
		t.Fatal(err)
	}

	ci := &CommandInfo{Exe: "go"}
	filter, err := LoadDSLFilter(f.Name(), ci)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if filter == nil {
		t.Fatal("expected non-nil filter for matching rule, got nil")
	}

	input := "aaa line\nbbb line\nccc aaa"
	expected := "aaa line\nccc aaa" // first rule (include "aaa") wins, only lines with "aaa" kept
	got := filter.Apply(input)

	if got != expected {
		t.Fatalf("FirstMatchWins mismatch:\n  got:  %q\n  want: %q", got, expected)
	}
}

func TestLoadDSLFilter_EmptyStage(t *testing.T) {
	configJSON := `{"rules":[{"exe":"echo","stages":[{}]}]}`

	f, err := os.CreateTemp(t.TempDir(), "emptystage*.json")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	if _, err := f.WriteString(configJSON); err != nil {
		t.Fatal(err)
	}

	ci := &CommandInfo{Exe: "echo"}
	filter, err := LoadDSLFilter(f.Name(), ci)
	if err == nil {
		t.Error("expected non-nil error for empty stage, got nil")
	}
	if filter != nil {
		t.Errorf("expected nil filter for empty stage, got %v", filter)
	}
}

func TestLoadDSLFilter_InvalidRegex(t *testing.T) {
	configJSON := `{"rules":[{"exe":"echo","stages":[{"include":{"pattern":"[invalid"}}]}]}`

	f, err := os.CreateTemp(t.TempDir(), "invalidregex*.json")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	if _, err := f.WriteString(configJSON); err != nil {
		t.Fatal(err)
	}

	ci := &CommandInfo{Exe: "echo"}
	filter, err := LoadDSLFilter(f.Name(), ci)
	if err == nil {
		t.Error("expected non-nil error for invalid regex, got nil")
	}
	if filter != nil {
		t.Errorf("expected nil filter for invalid regex, got %v", filter)
	}
}

func TestLoadDSLFilter_ValidConfig(t *testing.T) {
	configJSON := `{"rules":[{"exe":"echo","arg1":"hello","stages":[{"passthrough":{}}]}]}`

	f, err := os.CreateTemp(t.TempDir(), "valid*.json")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	if _, err := f.WriteString(configJSON); err != nil {
		t.Fatal(err)
	}

	ci := &CommandInfo{Exe: "echo", Arg1: "hello"}
	filter, err := LoadDSLFilter(f.Name(), ci)
	if err != nil {
		t.Fatalf("expected nil error for valid config, got %v", err)
	}
	if filter == nil {
		t.Fatal("expected non-nil filter for matching rule, got nil")
	}
}
