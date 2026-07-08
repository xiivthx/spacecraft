package research

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestFormatResearchResultsHuman verifies that FormatResearchResults
// produces human-readable output for []ResearchResult with a scope.
func TestFormatResearchResultsHuman(t *testing.T) {
	results := []ResearchResult{
		{
			Title:     "React Docs",
			URL:       "https://react.dev",
			Snippet:   "Official React documentation",
			Relevance: 0.95,
		},
		{
			Title:     "React GitHub",
			URL:       "https://github.com/facebook/react",
			Snippet:   "React source code",
			Relevance: 0.90,
		},
	}

	var buf bytes.Buffer
	opts := FormatOptions{
		Query: "react hooks",
		Scope: "react",
		Source: "brave-search",
	}

	err := FormatResults(&buf, results, opts)
	if err != nil {
		t.Fatalf("FormatResults returned error: %v", err)
	}

	output := buf.String()

	if !strings.Contains(output, "react hooks") {
		t.Errorf("output missing query %q", "react hooks")
	}
	if !strings.Contains(output, "react") {
		t.Errorf("output missing scope name %q", "react")
	}
	if !strings.Contains(output, "1.") {
		t.Errorf("output missing result numbering")
	}
	if !strings.Contains(output, "https://react.dev") {
		t.Errorf("output missing URL https://react.dev")
	}
	if !strings.Contains(output, "https://github.com/facebook/react") {
		t.Errorf("output missing URL https://github.com/facebook/react")
	}
	if !strings.Contains(output, "React Docs") {
		t.Errorf("output missing result title %q", "React Docs")
	}
}

// TestFormatResearchResultsJSON verifies that FormatResults with JSON=true
// produces valid JSON with the expected structure for []ResearchResult.
func TestFormatResearchResultsJSON(t *testing.T) {
	results := []ResearchResult{
		{
			Title:     "Test Result",
			URL:       "https://example.com",
			Snippet:   "A test result",
			Relevance: 0.85,
		},
	}

	var buf bytes.Buffer
	opts := FormatOptions{
		JSON:   true,
		Query:  "test query",
		Scope:  "test-scope",
		Source: "brave-search",
	}

	err := FormatResults(&buf, results, opts)
	if err != nil {
		t.Fatalf("FormatResults returned error: %v", err)
	}

	var parsed struct {
		Query        string           `json:"query"`
		Timestamp    string           `json:"timestamp"`
		Source       string           `json:"source"`
		Scope        string           `json:"scope"`
		Method       string           `json:"method"`
		Results      []ResearchResult `json:"results"`
		DeepAnalysis *DeepResult      `json:"deep_analysis"`
	}
	if err := json.Unmarshal(buf.Bytes(), &parsed); err != nil {
		t.Fatalf("invalid JSON output: %v", err)
	}

	if parsed.Query != "test query" {
		t.Errorf("query = %q, want %q", parsed.Query, "test query")
	}
	if parsed.Source != "brave-search" {
		t.Errorf("source = %q, want %q", parsed.Source, "brave-search")
	}
	if parsed.Scope != "test-scope" {
		t.Errorf("scope = %q, want %q", parsed.Scope, "test-scope")
	}
	if parsed.Timestamp == "" {
		t.Error("timestamp is empty")
	}
	if len(parsed.Results) != 1 {
		t.Fatalf("got %d results, want 1", len(parsed.Results))
	}
	if parsed.Results[0].Title != "Test Result" {
		t.Errorf("result title = %q, want %q", parsed.Results[0].Title, "Test Result")
	}
}

// TestFormatPackageInfoHuman verifies human-readable output for *PackageInfo.
func TestFormatPackageInfoHuman(t *testing.T) {
	pkg := &PackageInfo{
		Name:          "lodash",
		LatestVersion: "4.17.21",
		License:       "MIT",
		Published:     "2021-02-20",
		Homepage:      "https://lodash.com",
		Source:        "registry.npmjs.org",
	}

	var buf bytes.Buffer
	opts := FormatOptions{
		Query:  "lodash",
		Source: "npm",
	}

	err := FormatResults(&buf, pkg, opts)
	if err != nil {
		t.Fatalf("FormatResults returned error: %v", err)
	}

	output := buf.String()

	if !strings.Contains(output, "lodash") {
		t.Errorf("output missing package name %q", "lodash")
	}
	if !strings.Contains(output, "4.17.21") {
		t.Errorf("output missing version %q", "4.17.21")
	}
	if !strings.Contains(output, "MIT") {
		t.Errorf("output missing license %q", "MIT")
	}
	if !strings.Contains(output, "https://lodash.com") {
		t.Errorf("output missing homepage %q", "https://lodash.com")
	}
	if !strings.Contains(output, "registry.npmjs.org") {
		t.Errorf("output missing source %q", "registry.npmjs.org")
	}
}

// TestFormatPackageInfoJSON verifies JSON output for *PackageInfo.
func TestFormatPackageInfoJSON(t *testing.T) {
	pkg := &PackageInfo{
		Name:          "express",
		LatestVersion: "4.18.2",
		License:       "MIT",
		Published:     "2022-10-20",
		Homepage:      "https://expressjs.com",
		Source:        "registry.npmjs.org",
	}

	var buf bytes.Buffer
	opts := FormatOptions{
		JSON:   true,
		Query:  "express",
		Source: "npm",
	}

	err := FormatResults(&buf, pkg, opts)
	if err != nil {
		t.Fatalf("FormatResults returned error: %v", err)
	}

	var parsed struct {
		Query        string       `json:"query"`
		Timestamp    string       `json:"timestamp"`
		Source       string       `json:"source"`
		Method       string       `json:"method"`
		Package      *PackageInfo `json:"package"`
		DeepAnalysis *DeepResult  `json:"deep_analysis"`
	}
	if err := json.Unmarshal(buf.Bytes(), &parsed); err != nil {
		t.Fatalf("invalid JSON output: %v", err)
	}

	if parsed.Query != "express" {
		t.Errorf("query = %q, want %q", parsed.Query, "express")
	}
	if parsed.Source != "npm" {
		t.Errorf("source = %q, want %q", parsed.Source, "npm")
	}
	if parsed.Package == nil {
		t.Fatal("package field is nil")
	}
	if parsed.Package.Name != "express" {
		t.Errorf("package.name = %q, want %q", parsed.Package.Name, "express")
	}
	if parsed.Package.LatestVersion != "4.18.2" {
		t.Errorf("package.latest_version = %q, want %q", parsed.Package.LatestVersion, "4.18.2")
	}
	if parsed.Package.License != "MIT" {
		t.Errorf("package.license = %q, want %q", parsed.Package.License, "MIT")
	}
	if parsed.Package.Homepage != "https://expressjs.com" {
		t.Errorf("package.homepage = %q, want %q", parsed.Package.Homepage, "https://expressjs.com")
	}
	if parsed.Package.Source != "registry.npmjs.org" {
		t.Errorf("package.source = %q, want %q", parsed.Package.Source, "registry.npmjs.org")
	}
}

// TestSaveResultsToDir verifies SaveResults writes JSON to an existing dir.
func TestSaveResultsToDir(t *testing.T) {
	dir := t.TempDir()

	results := []ResearchResult{
		{
			Title:     "React Docs",
			URL:       "https://react.dev",
			Snippet:   "Official React documentation",
			Relevance: 0.95,
		},
	}
	opts := FormatOptions{
		Query:  "react hooks",
		Source: "brave-search",
		Scope:  "react",
	}

	path, err := SaveResults(dir, "research", results, opts)
	if err != nil {
		t.Fatalf("SaveResults returned error: %v", err)
	}

	if path == "" {
		t.Fatal("SaveResults returned empty path")
	}

	// Verify file exists.
	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("saved file stat error: %v", err)
	}
	if info.Size() == 0 {
		t.Error("saved file is empty")
	}

	// Verify parent directory matches.
	if filepath.Dir(path) != dir {
		t.Errorf("file parent = %q, want %q", filepath.Dir(path), dir)
	}

	// Verify filename pattern: <prefix>-<timestamp>.json
	base := filepath.Base(path)
	if !strings.HasPrefix(base, "research-") {
		t.Errorf("filename %q does not have prefix %q", base, "research-")
	}
	if !strings.HasSuffix(base, ".json") {
		t.Errorf("filename %q does not end with %q", base, ".json")
	}

	// Verify JSON content with full envelope.
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("reading saved file: %v", err)
	}

	var decoded struct {
		Query        string           `json:"query"`
		Timestamp    string           `json:"timestamp"`
		Source       string           `json:"source"`
		Scope        string           `json:"scope"`
		Method       string           `json:"method"`
		Results      []ResearchResult `json:"results"`
		DeepAnalysis *DeepResult      `json:"deep_analysis"`
	}
	if err := json.Unmarshal(raw, &decoded); err != nil {
		t.Fatalf("saved file is not valid JSON: %v", err)
	}
	if decoded.Query != "react hooks" {
		t.Errorf("saved query = %q, want %q", decoded.Query, "react hooks")
	}
	if decoded.Source != "brave-search" {
		t.Errorf("saved source = %q, want %q", decoded.Source, "brave-search")
	}
	if decoded.Scope != "react" {
		t.Errorf("saved scope = %q, want %q", decoded.Scope, "react")
	}
	if decoded.Method != "search" {
		t.Errorf("saved method = %q, want %q", decoded.Method, "search")
	}
	if len(decoded.Results) != 1 {
		t.Fatalf("got %d results, want 1", len(decoded.Results))
	}
	if decoded.Results[0].Title != "React Docs" {
		t.Errorf("result title = %q, want %q", decoded.Results[0].Title, "React Docs")
	}
}

// TestSaveResultsCreatesDir verifies SaveResults creates a non-existent dir.
func TestSaveResultsCreatesDir(t *testing.T) {
	base := t.TempDir()
	subdir := filepath.Join(base, "new", "nested", "dir")

	results := []ResearchResult{
		{Title: "Test", URL: "https://example.com", Snippet: "test"},
	}
	opts := FormatOptions{Query: "hello", Source: "test"}

	path, err := SaveResults(subdir, "test", results, opts)
	if err != nil {
		t.Fatalf("SaveResults returned error: %v", err)
	}

	if path == "" {
		t.Fatal("SaveResults returned empty path")
	}

	// Verify directory was created.
	info, err := os.Stat(subdir)
	if err != nil {
		t.Fatalf("expected directory %q to exist: %v", subdir, err)
	}
	if !info.IsDir() {
		t.Fatalf("%q is not a directory", subdir)
	}

	// Verify file exists in the created directory.
	if filepath.Dir(path) != subdir {
		t.Errorf("file parent = %q, want %q", filepath.Dir(path), subdir)
	}

	info, err = os.Stat(path)
	if err != nil {
		t.Fatalf("saved file stat error: %v", err)
	}
	if info.Size() == 0 {
		t.Error("saved file is empty")
	}
}
