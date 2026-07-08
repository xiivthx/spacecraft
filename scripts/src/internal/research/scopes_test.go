package research

import (
	"os"
	"path/filepath"
	"testing"
)

// TestScopeConfigStruct verifies the ScopeConfig type can be instantiated.
func TestScopeConfigStruct(t *testing.T) {
	sc := ScopeConfig{
		Scopes: map[string]SearchScope{
			"test": {Domains: []string{"example.com"}, Description: "test"},
		},
	}
	if sc.Scopes == nil {
		t.Error("ScopeConfig.Scopes is nil")
	}
	if len(sc.Scopes) != 1 {
		t.Errorf("ScopeConfig.Scopes has length %d, want 1", len(sc.Scopes))
	}
}

// TestDefaultScopes verifies that DefaultScopes() returns all 9 built-in scopes.
func TestDefaultScopes(t *testing.T) {
	cfg := DefaultScopes()

	if len(cfg.Scopes) != 9 {
		t.Fatalf("DefaultScopes() returned %d scopes, want 9", len(cfg.Scopes))
	}

	// Verify all expected scope keys are present.
	expected := []string{
		"react",
		"npm",
		"pypi",
		"tailwindcss",
		"nextjs",
		"storybook",
		"postgresql",
		"go",
		"rust",
	}
	for _, name := range expected {
		if _, ok := cfg.Scopes[name]; !ok {
			t.Errorf("DefaultScopes() missing scope %q", name)
		}
	}
}

// TestDefaultScopesDomains checks domain values for at least 2 scopes.
func TestDefaultScopesDomains(t *testing.T) {
	cfg := DefaultScopes()

	// Check "react" scope domains.
	react, ok := cfg.Scopes["react"]
	if !ok {
		t.Fatal("DefaultScopes() missing 'react' scope")
	}
	expectedReact := []string{"react.dev", "legacy.reactjs.org"}
	if len(react.Domains) != len(expectedReact) {
		t.Fatalf("react domains length = %d, want %d", len(react.Domains), len(expectedReact))
	}
	for i, d := range expectedReact {
		if react.Domains[i] != d {
			t.Errorf("react domains[%d] = %q, want %q", i, react.Domains[i], d)
		}
	}

	// Check "go" scope domains.
	goScope, ok := cfg.Scopes["go"]
	if !ok {
		t.Fatal("DefaultScopes() missing 'go' scope")
	}
	expectedGo := []string{"go.dev/doc", "pkg.go.dev"}
	if len(goScope.Domains) != len(expectedGo) {
		t.Fatalf("go domains length = %d, want %d", len(goScope.Domains), len(expectedGo))
	}
	for i, d := range expectedGo {
		if goScope.Domains[i] != d {
			t.Errorf("go domains[%d] = %q, want %q", i, goScope.Domains[i], d)
		}
	}
}

// TestSmartScope verifies SmartScope auto-detection from query keywords and manifests.
func TestSmartScope(t *testing.T) {
	tests := []struct {
		name      string
		query     string
		manifests []string
		want      string
	}{
		{
			name:  "query keyword react",
			query: "react server components",
			want:  "react",
		},
		{
			name:  "query keyword tailwind",
			query: "tailwind flex center",
			want:  "tailwindcss",
		},
		{
			name:  "manifest go.mod",
			query: "",
			manifests: []string{"go.mod"},
			want:  "go",
		},
		{
			name:  "manifest Cargo.toml",
			query: "",
			manifests: []string{"Cargo.toml"},
			want:  "rust",
		},
		{
			name:  "query keyword golang",
			query: "golang generics",
			want:  "go",
		},
		{
			name:  "query keyword rust",
			query: "rust async",
			want:  "rust",
		},
		{
			name:  "no match at all",
			query: "",
			want:  "",
		},
		{
			name:  "query keyword next",
			query: "next app",
			want:  "nextjs",
		},
		{
			name:  "query keyword postgres",
			query: "postgres query",
			want:  "postgresql",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SmartScope(tt.query, tt.manifests)
			if got != tt.want {
				t.Errorf("SmartScope(%q, %v) = %q, want %q", tt.query, tt.manifests, got, tt.want)
			}
		})
	}
}

// TestLoadScopes verifies LoadScopes loads/merges .space/scopes.json.
func TestLoadScopes(t *testing.T) {
	// Test 1: No file exists -> returns DefaultScopes(), no error.
	t.Run("no file returns defaults", func(t *testing.T) {
		cfg, err := LoadScopes("/nonexistent/path/scopes.json")
		if err != nil {
			t.Errorf("LoadScopes returned error for missing file: %v", err)
		}
		defaults := DefaultScopes()
		if len(cfg.Scopes) != len(defaults.Scopes) {
			t.Errorf("got %d scopes, want %d", len(cfg.Scopes), len(defaults.Scopes))
		}
	})

	// Test 2: Valid file merges scopes: user overrides built-in by name,
	// new scopes are added, unmentioned defaults preserved.
	t.Run("valid file merges scopes", func(t *testing.T) {
		dir := t.TempDir()
		path := filepath.Join(dir, "scopes.json")
		content := `{
			"scopes": {
				"react": {
					"domains": ["react.dev", "legacy.reactjs.org", "beta.reactjs.org"],
					"description": "React docs"
				},
				"svelte": {
					"domains": ["svelte.dev"],
					"description": "Svelte documentation"
				}
			}
		}`
		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			t.Fatal(err)
		}

		cfg, err := LoadScopes(path)
		if err != nil {
			t.Fatalf("LoadScopes returned error: %v", err)
		}

		// Built-in scope "react" should be overridden (description changed).
		react, ok := cfg.Scopes["react"]
		if !ok {
			t.Fatal("react scope missing after merge")
		}
		if react.Description != "React docs" {
			t.Errorf("react description = %q, want %q", react.Description, "React docs")
		}
		if len(react.Domains) != 3 {
			t.Errorf("react domains length = %d, want 3", len(react.Domains))
		}

		// New scope "svelte" should be added.
		if _, ok := cfg.Scopes["svelte"]; !ok {
			t.Error("svelte scope missing after merge (should be added)")
		}

		// Unmentioned default "tailwindcss" should be preserved.
		if _, ok := cfg.Scopes["tailwindcss"]; !ok {
			t.Error("tailwindcss scope missing after merge (should be preserved from defaults)")
		}
	})

	// Test 3: Invalid JSON file returns an error.
	t.Run("invalid json returns error", func(t *testing.T) {
		dir := t.TempDir()
		path := filepath.Join(dir, "scopes.json")
		content := `{invalid json}`
		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			t.Fatal(err)
		}

		_, err := LoadScopes(path)
		if err == nil {
			t.Error("LoadScopes should return error for invalid JSON")
		}
	})
}
