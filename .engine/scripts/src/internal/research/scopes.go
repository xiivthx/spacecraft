package research

import (
	"encoding/json"
	"os"
	"strings"
)

// ScopeConfig holds a named collection of research scopes.
type ScopeConfig struct {
	Scopes map[string]SearchScope
}

// DefaultScopes returns the built-in research scopes.
func DefaultScopes() ScopeConfig {
	return ScopeConfig{
		Scopes: map[string]SearchScope{
			"react": {
				Domains:     []string{"react.dev", "legacy.reactjs.org"},
				Description: "React documentation and API reference",
			},
			"npm": {
				Domains:     []string{"docs.npmjs.com"},
				Description: "npm package manager documentation",
			},
			"pypi": {
				Domains:     []string{"docs.python.org", "pypi.org"},
				Description: "Python and PyPI documentation",
			},
			"tailwindcss": {
				Domains:     []string{"tailwindcss.com/docs"},
				Description: "Tailwind CSS utility class documentation",
			},
			"nextjs": {
				Domains:     []string{"nextjs.org/docs"},
				Description: "Next.js framework documentation",
			},
			"storybook": {
				Domains:     []string{"storybook.js.org/docs"},
				Description: "Storybook component development documentation",
			},
			"postgresql": {
				Domains:     []string{"postgresql.org/docs"},
				Description: "PostgreSQL database documentation",
			},
			"go": {
				Domains:     []string{"go.dev/doc", "pkg.go.dev"},
				Description: "Go language and standard library documentation",
			},
			"rust": {
				Domains:     []string{"doc.rust-lang.org", "docs.rs"},
				Description: "Rust language and crate documentation",
			},
		},
	}
}

// LoadScopes loads scopes from path and merges them over the built-in defaults.
// If path does not exist, it returns DefaultScopes().
func LoadScopes(path string) (ScopeConfig, error) {
	defaults := DefaultScopes()

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return defaults, nil
		}
		return ScopeConfig{}, err
	}

	var override ScopeConfig
	if err := json.Unmarshal(data, &override); err != nil {
		return ScopeConfig{}, err
	}

	for name, scope := range override.Scopes {
		defaults.Scopes[name] = scope
	}

	return defaults, nil
}

// SmartScope auto-detects a research scope name from query keywords and manifests.
func SmartScope(query string, manifests []string) string {
	keywordToScope := map[string]string{
		"react":       "react",
		"tailwind":    "tailwindcss",
		"tailwindcss": "tailwindcss",
		"next":        "nextjs",
		"nextjs":      "nextjs",
		"storybook":   "storybook",
		"postgres":    "postgresql",
		"postgresql":  "postgresql",
		"go":          "go",
		"golang":      "go",
		"rust":        "rust",
		"cargo":       "rust",
	}

	for _, word := range strings.Fields(strings.ToLower(query)) {
		if scope, ok := keywordToScope[word]; ok {
			return scope
		}
	}

	manifestToScope := map[string]string{
		"go.mod":           "go",
		"Cargo.toml":       "rust",
		"package.json":     "npm",
		"requirements.txt": "pypi",
		"pyproject.toml":   "pypi",
	}

	for _, m := range manifests {
		if scope, ok := manifestToScope[m]; ok {
			return scope
		}
	}

	return ""
}
