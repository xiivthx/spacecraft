package research

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

// FormatOptions controls how research results are formatted for output.
type FormatOptions struct {
	JSON         bool
	Query        string
	Scope        string // scope name, empty if none
	Source       string // "brave-search", "npm", "pypi", "go", "crates"
	Deep         bool
	Method       string // "search" (default) or "browser-use" or "nlm" for --deep
	DeepAnalysis *DeepResult
	Timestamp    string
	PersistDir   string // empty means --no-save (caller handles this)
}

// FormatResults writes results to w using the format selected in opts.
// Supported result types are []ResearchResult and *PackageInfo.
func FormatResults(w io.Writer, results interface{}, opts FormatOptions) error {
	if opts.JSON {
		return formatJSON(w, results, opts)
	}
	return formatHuman(w, results, opts)
}

func formatHuman(w io.Writer, results interface{}, opts FormatOptions) error {
	switch v := results.(type) {
	case []ResearchResult:
		return formatResearchResultsHuman(w, v, opts)
	case *PackageInfo:
		return formatPackageInfoHuman(w, v, opts)
	default:
		return fmt.Errorf("unsupported result type %T for human formatting", results)
	}
}

func formatResearchResultsHuman(w io.Writer, results []ResearchResult, opts FormatOptions) error {
	fmt.Fprintf(w, "Query: %s\n", opts.Query)
	fmt.Fprintln(w)
	if opts.Scope != "" {
		fmt.Fprintf(w, "Scope: %s (%s)\n", opts.Scope, opts.Source)
		fmt.Fprintln(w)
	}
	for i, r := range results {
		if i > 0 {
			fmt.Fprintln(w)
		}
		fmt.Fprintf(w, "%d. %s\n", i+1, r.Title)
		fmt.Fprintf(w, "   %s\n", r.URL)
		fmt.Fprintf(w, "   %s\n", r.Snippet)
	}
	return nil
}

func formatPackageInfoHuman(w io.Writer, pkg *PackageInfo, opts FormatOptions) error {
	fmt.Fprintf(w, "Query: %s\n", opts.Query)
	fmt.Fprintln(w)
	fmt.Fprintf(w, "%s (%s)\n", pkg.Name, opts.Source)
	fmt.Fprintf(w, "  Latest:    %s\n", pkg.LatestVersion)
	fmt.Fprintf(w, "  License:   %s\n", pkg.License)
	fmt.Fprintf(w, "  Published: %s\n", pkg.Published)
	fmt.Fprintf(w, "  Homepage:  %s\n", pkg.Homepage)
	fmt.Fprintf(w, "  Source:    %s\n", pkg.Source)
	return nil
}

func formatJSON(w io.Writer, results interface{}, opts FormatOptions) error {
	timestamp := opts.Timestamp
	if timestamp == "" {
		timestamp = time.Now().UTC().Format("2006-01-02T15:04:05Z")
	}

	method := opts.Method
	if method == "" {
		method = "search"
	}

	switch v := results.(type) {
	case []ResearchResult:
		env := struct {
			Query        string           `json:"query"`
			Timestamp    string           `json:"timestamp"`
			Source       string           `json:"source"`
			Scope        string           `json:"scope"`
			Method       string           `json:"method"`
			Results      []ResearchResult `json:"results"`
			DeepAnalysis *DeepResult      `json:"deep_analysis"`
		}{
			Query:        opts.Query,
			Timestamp:    timestamp,
			Source:       opts.Source,
			Scope:        opts.Scope,
			Method:       method,
			Results:      v,
			DeepAnalysis: opts.DeepAnalysis,
		}
		return encodeJSON(w, env)
	case *PackageInfo:
		env := struct {
			Query        string       `json:"query"`
			Timestamp    string       `json:"timestamp"`
			Source       string       `json:"source"`
			Method       string       `json:"method"`
			Package      *PackageInfo `json:"package"`
			DeepAnalysis *DeepResult  `json:"deep_analysis"`
		}{
			Query:        opts.Query,
			Timestamp:    timestamp,
			Source:       opts.Source,
			Method:       method,
			Package:      v,
			DeepAnalysis: opts.DeepAnalysis,
		}
		return encodeJSON(w, env)
	default:
		return fmt.Errorf("unsupported result type %T for JSON formatting", results)
	}
}

func encodeJSON(w io.Writer, v interface{}) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(v)
}

// SaveResults persists the research output as indented JSON to dir with an
// auto-generated filename.  The file is written as the same full envelope that
// --json would produce (query, timestamp, source, scope, method, results,
// deep_analysis / package).
func SaveResults(dir string, prefix string, results interface{}, opts FormatOptions) (string, error) {
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", err
	}

	timestamp := time.Now().UTC().Format("20060102T150405")
	filename := fmt.Sprintf("%s-%s.json", prefix, timestamp)
	path := filepath.Join(dir, filename)

	// Build the envelope inline so the file matches --json output.
	ts := opts.Timestamp
	if ts == "" {
		ts = time.Now().UTC().Format("2006-01-02T15:04:05Z")
	}
	method := opts.Method
	if method == "" {
		method = "search"
	}

	var env interface{}
	switch v := results.(type) {
	case []ResearchResult:
		env = struct {
			Query        string           `json:"query"`
			Timestamp    string           `json:"timestamp"`
			Source       string           `json:"source"`
			Scope        string           `json:"scope"`
			Method       string           `json:"method"`
			Results      []ResearchResult `json:"results"`
			DeepAnalysis *DeepResult      `json:"deep_analysis"`
		}{
			Query:        opts.Query,
			Timestamp:    ts,
			Source:       opts.Source,
			Scope:        opts.Scope,
			Method:       method,
			Results:      v,
			DeepAnalysis: opts.DeepAnalysis,
		}
	case *PackageInfo:
		env = struct {
			Query        string       `json:"query"`
			Timestamp    string       `json:"timestamp"`
			Source       string       `json:"source"`
			Method       string       `json:"method"`
			Package      *PackageInfo `json:"package"`
			DeepAnalysis *DeepResult  `json:"deep_analysis"`
		}{
			Query:        opts.Query,
			Timestamp:    ts,
			Source:       opts.Source,
			Method:       method,
			Package:      v,
			DeepAnalysis: opts.DeepAnalysis,
		}
	default:
		return "", fmt.Errorf("unsupported result type %T for saving", results)
	}

	raw, err := json.MarshalIndent(env, "", "  ")
	if err != nil {
		return "", err
	}

	if err := os.WriteFile(path, raw, 0644); err != nil {
		return "", err
	}

	return path, nil
}
