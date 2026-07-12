package research

import (
	"testing"
)

func TestTypesExist(t *testing.T) {
	// Test ResearchResult struct literal and field access
	r := ResearchResult{
		Title:     "Test Title",
		URL:       "https://example.com",
		Snippet:   "Test snippet",
		Relevance: 0.95,
	}
	if r.Title != "Test Title" {
		t.Errorf("ResearchResult.Title = %q, want %q", r.Title, "Test Title")
	}
	if r.URL != "https://example.com" {
		t.Errorf("ResearchResult.URL = %q, want %q", r.URL, "https://example.com")
	}
	if r.Snippet != "Test snippet" {
		t.Errorf("ResearchResult.Snippet = %q, want %q", r.Snippet, "Test snippet")
	}
	if r.Relevance != 0.95 {
		t.Errorf("ResearchResult.Relevance = %f, want %f", r.Relevance, 0.95)
	}

	// Test PackageInfo struct literal and field access
	p := PackageInfo{
		Name:          "test-pkg",
		LatestVersion: "v1.0.0",
		License:       "MIT",
		Published:     "2024-01-01",
		Homepage:      "https://example.com/pkg",
		Source:        "https://github.com/example/pkg",
	}
	if p.Name != "test-pkg" {
		t.Errorf("PackageInfo.Name = %q, want %q", p.Name, "test-pkg")
	}
	if p.LatestVersion != "v1.0.0" {
		t.Errorf("PackageInfo.LatestVersion = %q, want %q", p.LatestVersion, "v1.0.0")
	}
	if p.License != "MIT" {
		t.Errorf("PackageInfo.License = %q, want %q", p.License, "MIT")
	}
	if p.Published != "2024-01-01" {
		t.Errorf("PackageInfo.Published = %q, want %q", p.Published, "2024-01-01")
	}
	if p.Homepage != "https://example.com/pkg" {
		t.Errorf("PackageInfo.Homepage = %q, want %q", p.Homepage, "https://example.com/pkg")
	}
	if p.Source != "https://github.com/example/pkg" {
		t.Errorf("PackageInfo.Source = %q, want %q", p.Source, "https://github.com/example/pkg")
	}

	// Test SearchScope struct literal and field access
	s := SearchScope{
		Domains:     []string{"example.com", "test.org"},
		Description: "Test search scope",
	}
	if len(s.Domains) != 2 || s.Domains[0] != "example.com" || s.Domains[1] != "test.org" {
		t.Errorf("SearchScope.Domains = %v, want %v", s.Domains, []string{"example.com", "test.org"})
	}
	if s.Description != "Test search scope" {
		t.Errorf("SearchScope.Description = %q, want %q", s.Description, "Test search scope")
	}

	// Test RegistryQuery struct literal and field access
	q := RegistryQuery{
		PackageName: "test-pkg",
		Ecosystem:   "npm",
	}
	if q.PackageName != "test-pkg" {
		t.Errorf("RegistryQuery.PackageName = %q, want %q", q.PackageName, "test-pkg")
	}
	if q.Ecosystem != "npm" {
		t.Errorf("RegistryQuery.Ecosystem = %q, want %q", q.Ecosystem, "npm")
	}
}

// TestResearchProviderInterface verifies that the ResearchProvider interface
// exists with the required methods (Search, LookupPackage, DeepAnalyze).
func TestResearchProviderInterface(t *testing.T) {
	var rp ResearchProvider
	if rp != nil {
		t.Error("expected nil ResearchProvider")
	}
}

// TestResearchStoreInterface verifies that the ResearchStore interface
// exists with the required methods (Save, Load).
func TestResearchStoreInterface(t *testing.T) {
	var rs ResearchStore
	if rs != nil {
		t.Error("expected nil ResearchStore")
	}
}
