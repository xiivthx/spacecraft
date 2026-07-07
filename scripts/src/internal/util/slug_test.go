package util

import (
	"testing"
)

func TestSlugify(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"Hello World", "hello-world"},
		{"Go API", "go-api"},
		{"  Trim  Me  ", "trim-me"},
		{"Special!@#Chars", "special-chars"},
		{"UPPERCASE", "uppercase"},
		{"a", "a"},
		{"", "mission"},
		{"xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx", "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"},
	}
	for _, tt := range tests {
		got := Slugify(tt.input)
		if got != tt.want {
			t.Errorf("Slugify(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func strings64(ch string, n int) string {
	s := ""
	for i := 0; i < n; i++ {
		s += ch
	}
	return s
}

func TestCommandToString(t *testing.T) {
	tests := []struct {
		parts []string
		want  string
	}{
		{[]string{"make", "test"}, "make test"},
		{[]string{"echo", "hello world"}, "echo 'hello world'"},
		{[]string{"cat", "file.txt"}, "cat file.txt"},
		{[]string{"sh", "-c", "echo 'hi'"}, "sh -c 'echo '\\''hi'\\'''"},
	}
	for _, tt := range tests {
		got := CommandToString(tt.parts)
		if got != tt.want {
			t.Errorf("CommandToString(%v) = %q, want %q", tt.parts, got, tt.want)
		}
	}
}

func TestNormalizeMissionId(t *testing.T) {
	tests := []struct {
		input string
		want  *string
	}{
		{"M07ABCDEF", ptr("M07ABCDEF")},
		{"m07abcdef", ptr("M07ABCDEF")},
		{"M-20260707-141230", ptr("M-20260707-141230")},
		{"m-20260707-141230", ptr("M-20260707-141230")},
		{"branch/feat/m07abcdef/title", ptr("M07ABCDEF")},
		{"invalid", nil},
		{"", nil},
		{"   ", nil},
		{"M07", nil}, // too short
	}
	for _, tt := range tests {
		got := NormalizeMissionId(tt.input)
		if tt.want == nil {
			if got != nil {
				t.Errorf("NormalizeMissionId(%q) = %v, want nil", tt.input, *got)
			}
		} else if got == nil {
			t.Errorf("NormalizeMissionId(%q) = nil, want %q", tt.input, *tt.want)
		} else if *got != *tt.want {
			t.Errorf("NormalizeMissionId(%q) = %q, want %q", tt.input, *got, *tt.want)
		}
	}
}

func TestContainsStr(t *testing.T) {
	slice := []string{"a", "b", "c"}
	if !ContainsStr(slice, "b") {
		t.Error("ContainsStr should find 'b'")
	}
	if ContainsStr(slice, "z") {
		t.Error("ContainsStr should not find 'z'")
	}
	if ContainsStr(nil, "a") {
		t.Error("ContainsStr on nil should not find anything")
	}
}

func ptr(s string) *string { return &s }

func TestRegexpReplace(t *testing.T) {
	got := RegexpReplace(`\s+`, "_", "hello   world")
	if got != "hello_world" {
		t.Errorf("RegexpReplace = %q, want %q", got, "hello_world")
	}
}

func TestIsoNowFormat(t *testing.T) {
	now := IsoNow()
	if len(now) < 20 {
		t.Errorf("IsoNow too short: %q", now)
	}
	// Check it looks like ISO: 2026-07-07T11:21:27.586Z
	if now[len(now)-1] != 'Z' {
		t.Errorf("IsoNow should end with Z: %q", now)
	}
}
