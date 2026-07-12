package research

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestGoProxyLookup verifies that GoProxyClient.Lookup fetches and parses
// the Go proxy /@latest endpoint correctly.
func TestGoProxyLookup(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{
			"Version": "v1.3.0",
			"Time":    "2026-07-01T00:00:00Z",
		})
	}))
	defer srv.Close()

	client := NewGoProxyClient(srv.URL)
	info, err := client.Lookup(context.Background(), "github.com/charmbracelet/bubbletea")
	if err != nil {
		t.Fatalf("GoProxyClient.Lookup returned error: %v", err)
	}
	if info.Name != "github.com/charmbracelet/bubbletea" {
		t.Errorf("Name = %q, want %q", info.Name, "github.com/charmbracelet/bubbletea")
	}
	if info.LatestVersion != "v1.3.0" {
		t.Errorf("LatestVersion = %q, want %q", info.LatestVersion, "v1.3.0")
	}
	if info.Source != "proxy.golang.org" {
		t.Errorf("Source = %q, want %q", info.Source, "proxy.golang.org")
	}
	if info.Published != "2026-07-01T00:00:00Z" {
		t.Errorf("Published = %q, want %q", info.Published, "2026-07-01T00:00:00Z")
	}
}

// TestNpmLookup verifies that NpmClient.Lookup fetches and parses
// the npm registry install-v1 response correctly.
func TestNpmLookup(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.Header.Get("Accept") != "application/vnd.npm.install-v1+json" {
			t.Errorf("expected Accept header %q, got %q",
				"application/vnd.npm.install-v1+json", r.Header.Get("Accept"))
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"name":        "express",
			"version":     "5.1.0",
			"license":     "MIT",
			"description": "Fast, unopinionated, minimalist web framework",
			"homepage":    "https://expressjs.com",
		})
	}))
	defer srv.Close()

	client := NewNpmClient(srv.URL)
	info, err := client.Lookup(context.Background(), "express")
	if err != nil {
		t.Fatalf("NpmClient.Lookup returned error: %v", err)
	}
	if info.Name != "express" {
		t.Errorf("Name = %q, want %q", info.Name, "express")
	}
	if info.LatestVersion != "5.1.0" {
		t.Errorf("LatestVersion = %q, want %q", info.LatestVersion, "5.1.0")
	}
	if info.License != "MIT" {
		t.Errorf("License = %q, want %q", info.License, "MIT")
	}
	if info.Homepage != "https://expressjs.com" {
		t.Errorf("Homepage = %q, want %q", info.Homepage, "https://expressjs.com")
	}
	if info.Source != "registry.npmjs.org" {
		t.Errorf("Source = %q, want %q", info.Source, "registry.npmjs.org")
	}
}

// TestPypiLookup verifies that PypiClient.Lookup fetches and parses
// the PyPI JSON API response correctly.
func TestPypiLookup(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"info": map[string]interface{}{
				"name":      "requests",
				"version":   "2.32.0",
				"license":   "Apache-2.0",
				"home_page": "https://requests.readthedocs.io",
				"summary":   "HTTP library for Python",
			},
		})
	}))
	defer srv.Close()

	client := NewPypiClient(srv.URL)
	info, err := client.Lookup(context.Background(), "requests")
	if err != nil {
		t.Fatalf("PypiClient.Lookup returned error: %v", err)
	}
	if info.Name != "requests" {
		t.Errorf("Name = %q, want %q", info.Name, "requests")
	}
	if info.LatestVersion != "2.32.0" {
		t.Errorf("LatestVersion = %q, want %q", info.LatestVersion, "2.32.0")
	}
	if info.License != "Apache-2.0" {
		t.Errorf("License = %q, want %q", info.License, "Apache-2.0")
	}
	if info.Homepage != "https://requests.readthedocs.io" {
		t.Errorf("Homepage = %q, want %q", info.Homepage, "https://requests.readthedocs.io")
	}
	if info.Source != "pypi.org" {
		t.Errorf("Source = %q, want %q", info.Source, "pypi.org")
	}
}

// TestCargoLookup verifies that CargoClient.Lookup fetches and parses
// the crates.io API response correctly.
func TestCargoLookup(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"crate": map[string]interface{}{
				"name":               "serde",
				"max_stable_version": "1.0.217",
				"homepage":           "https://serde.rs",
				"description":        "Serialization framework",
				"license":            "MIT OR Apache-2.0",
				"updated_at":         "2024-11-15T00:00:00Z",
			},
		})
	}))
	defer srv.Close()

	client := NewCargoClient(srv.URL)
	info, err := client.Lookup(context.Background(), "serde")
	if err != nil {
		t.Fatalf("CargoClient.Lookup returned error: %v", err)
	}
	if info.Name != "serde" {
		t.Errorf("Name = %q, want %q", info.Name, "serde")
	}
	if info.LatestVersion != "1.0.217" {
		t.Errorf("LatestVersion = %q, want %q", info.LatestVersion, "1.0.217")
	}
	if info.Homepage != "https://serde.rs" {
		t.Errorf("Homepage = %q, want %q", info.Homepage, "https://serde.rs")
	}
	if info.License != "MIT OR Apache-2.0" {
		t.Errorf("License = %q, want %q", info.License, "MIT OR Apache-2.0")
	}
	if info.Published != "2024-11-15T00:00:00Z" {
		t.Errorf("Published = %q, want %q", info.Published, "2024-11-15T00:00:00Z")
	}
	if info.Source != "crates.io" {
		t.Errorf("Source = %q, want %q", info.Source, "crates.io")
	}
}

// TestRegistryLookupError verifies that registry clients return an error
// when the server responds with a non-200 status.
func TestRegistryLookupError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer srv.Close()

	client := NewGoProxyClient(srv.URL)
	_, err := client.Lookup(context.Background(), "some/module")
	if err == nil {
		t.Error("GoProxyClient.Lookup should return error on 500, got nil")
	}
}
