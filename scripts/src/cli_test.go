package main

import (
	"os"
	"testing"

	"spacecraft/internal/config"
	"spacecraft/internal/gitutil"
	"spacecraft/internal/mission"
	"spacecraft/internal/resolver"
)

// TestMain initializes the global config and dependencies needed by CLI tests.
func TestMain(m *testing.M) {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	cfg, err = config.NewConfig(cwd)
	if err != nil {
		panic(err)
	}
	store = mission.NewFSStore(cfg)
	r = resolver.New(store, gitutil.OSCommandRunner{}, nil)
	os.Exit(m.Run())
}

// TestResearchCmdHelp verifies --help returns exit code 0.
func TestResearchCmdHelp(t *testing.T) {
	exit := researchCmd([]string{"--help"})
	if exit != 0 {
		t.Errorf("researchCmd() with --help returned %d, want 0", exit)
	}
}

// TestResearchCmdMissingQuery verifies exit code 1 for missing query.
func TestResearchCmdMissingQuery(t *testing.T) {
	exit := researchCmd([]string{})
	if exit != 1 {
		t.Errorf("researchCmd() with empty args returned %d, want 1", exit)
	}
}

// TestResearchCmdBlankQuery verifies exit code 1 for blank query.
func TestResearchCmdBlankQuery(t *testing.T) {
	exit := researchCmd([]string{"   "})
	if exit != 1 {
		t.Errorf("researchCmd() with blank query returned %d, want 1", exit)
	}
}

// TestResearchCmdMissingAPIKey verifies exit code 2 for missing API key.
func TestResearchCmdMissingAPIKey(t *testing.T) {
	os.Unsetenv("SPACECRAFT_BRAVE_API_KEY")
	exit := researchCmd([]string{"test query"})
	if exit != 2 {
		t.Errorf("researchCmd() without SPACECRAFT_BRAVE_API_KEY returned %d, want 2", exit)
	}
}

// TestResearchCmdInvalidDeep verifies exit code 2 for invalid --deep value.
func TestResearchCmdInvalidDeep(t *testing.T) {
	exit := researchCmd([]string{"--deep", "invalid", "test"})
	if exit != 2 {
		t.Errorf("researchCmd() with --deep=invalid returned %d, want 2", exit)
	}
}

// TestResearchCmdHelpWithJSON verifies --json --help still returns 0.
func TestResearchCmdHelpWithJSON(t *testing.T) {
	exit := researchCmd([]string{"--json", "--help"})
	if exit != 0 {
		t.Errorf("researchCmd() --json --help returned %d, want 0", exit)
	}
}

// TestResearchCmdExitCodesMapping verifies the three exit code classes:
// 0 = success/help, 1 = usage/no results, 2 = error.
func TestResearchCmdExitCodesMapping(t *testing.T) {
	// Exit 0: --help
	if e := researchCmd([]string{"--help"}); e != 0 {
		t.Errorf("--help: got %d, want 0", e)
	}
	// Exit 2: missing API key (operational error)
	os.Unsetenv("SPACECRAFT_BRAVE_API_KEY")
	if e := researchCmd([]string{"query"}); e != 2 {
		t.Errorf("missing API key: got %d, want 2", e)
	}
	// Exit 2: invalid --deep value
	if e := researchCmd([]string{"--deep", "bad", "query"}); e != 2 {
		t.Errorf("invalid --deep: got %d, want 2", e)
	}
}

// TestCheckDepsCmdHelp verifies --help returns exit code 0.
func TestCheckDepsCmdHelp(t *testing.T) {
	exit := checkDepsCmd([]string{"--help"})
	if exit != 0 {
		t.Errorf("checkDepsCmd() with --help returned %d, want 0", exit)
	}
}

// TestCheckDepsCmdUnknownRegistry verifies exit code 2 for unknown --registry.
func TestCheckDepsCmdUnknownRegistry(t *testing.T) {
	exit := checkDepsCmd([]string{"--registry", "nonexistent", "--timeout", "1s"})
	if exit != 2 {
		t.Errorf("checkDepsCmd() with unknown registry returned %d, want 2", exit)
	}
}

// TestIsPackageQuery verifies the heuristic for detecting package queries.
// After M4 fix: requires dot, slash, or @ for registry trigger.
func TestIsPackageQuery(t *testing.T) {
	tests := []struct {
		query string
		want  bool
	}{
		{"express", false},            // single word, no special chars
		{"react", false},              // single word
		{"@scope/pkg", true},          // scoped npm package
		{"golang.org/x/tools", true},  // Go module with path
		{"react-dom", false},          // hyphenated but no dot/slash/@ — M4 fix: stricter
		{"how to use react", false},   // multi-word
		{"", false},                   // empty
		{"lodash", false},             // single word
		{"@babel/core", true},         // scoped npm
		{"github.com/foo/bar", true},  // Go-like module path
	}
	for _, tt := range tests {
		got := isPackageQuery(tt.query)
		if got != tt.want {
			t.Errorf("isPackageQuery(%q) = %v, want %v", tt.query, got, tt.want)
		}
	}
}

// TestCompareNumericVersions verifies version comparison including pre-releases.
func TestCompareNumericVersions(t *testing.T) {
	tests := []struct {
		current, latest string
		want            string
	}{
		{"1.0.0", "2.0.0", "MAJOR UPGRADE"},
		{"1.0.0", "1.1.0", "MINOR UPGRADE"},
		{"1.0.0", "1.0.1", "PATCH"},
		{"1.0.0", "1.0.0", "current"},
		{"2.0.0", "1.0.0", "current"},
		{"v1.0.0", "v2.0.0", "MAJOR UPGRADE"},
		{"1.0.0", "1.0.0-beta", "current"},      // latest is pre-release
		{"1.0.0-beta", "1.0.0", "PATCH"},        // current is pre-release, latest is release
		{"1.2.3.4", "1.2.3.5", "PATCH"},
		{"^1.0.0", "1.1.0", "MINOR UPGRADE"},
		{"~1.0.0", "1.0.1", "PATCH"},
		{">=1.0.0", "1.0.5", "PATCH"},
	}
	for _, tt := range tests {
		got := compareNumericVersions(tt.current, tt.latest)
		if got != tt.want {
			t.Errorf("compareNumericVersions(%q, %q) = %q, want %q",
				tt.current, tt.latest, got, tt.want)
		}
	}
}

// TestFatalErr verifies fatalErr returns exit code 2.
func TestFatalErr(t *testing.T) {
	code := fatalErr("test error")
	if code != 2 {
		t.Errorf("fatalErr() = %d, want 2", code)
	}
}

// TestPrintErr verifies printErr returns exit code 1.
func TestPrintErr(t *testing.T) {
	code := printErr("test error")
	if code != 1 {
		t.Errorf("printErr() = %d, want 1", code)
	}
}
