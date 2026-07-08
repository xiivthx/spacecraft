package research

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// TestBraveSearchSuccess verifies that BraveClient.Search correctly parses a
// valid Brave Search API response and returns ResearchResult items.
func TestBraveSearchSuccess(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.Header.Get("X-Subscription-Token") == "" {
			t.Error("expected X-Subscription-Token header to be set")
		}
		if r.URL.Query().Get("q") == "" {
			t.Error("expected 'q' query parameter to be set")
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"web": map[string]interface{}{
				"results": []map[string]interface{}{
					{
						"title":       "Flex — Tailwind CSS",
						"url":         "https://tailwindcss.com/docs/flex",
						"description": "Utilities for controlling how flex items grow and shrink.",
					},
				},
			},
		})
	}))
	defer srv.Close()

	client := NewBraveClient("test-api-key", srv.URL)
	results, err := client.Search(context.Background(), "tailwind flex center", nil, 0)
	if err != nil {
		t.Fatalf("BraveClient.Search returned error: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("got %d results, want 1", len(results))
	}
	if results[0].Title != "Flex — Tailwind CSS" {
		t.Errorf("Title = %q, want %q", results[0].Title, "Flex — Tailwind CSS")
	}
	if results[0].URL != "https://tailwindcss.com/docs/flex" {
		t.Errorf("URL = %q, want %q", results[0].URL, "https://tailwindcss.com/docs/flex")
	}
	if results[0].Snippet != "Utilities for controlling how flex items grow and shrink." {
		t.Errorf("Snippet = %q, want %q", results[0].Snippet, "Utilities for controlling how flex items grow and shrink.")
	}
	if results[0].Relevance != 0.0 {
		t.Errorf("Relevance = %f, want 0.0", results[0].Relevance)
	}
}

// TestBraveSearchWithDomainScoping verifies that when domains are provided,
// BraveClient.Search prepends "site:<domain>" to the query string.
func TestBraveSearchWithDomainScoping(t *testing.T) {
	var capturedQuery string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedQuery = r.URL.Query().Get("q")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"web": map[string]interface{}{
				"results": []map[string]interface{}{},
			},
		})
	}))
	defer srv.Close()

	client := NewBraveClient("test-api-key", srv.URL)
	_, err := client.Search(context.Background(), "tailwind flex center", []string{"react.dev"}, 0)
	if err != nil {
		t.Fatalf("BraveClient.Search returned error: %v", err)
	}
	expected := "site:react.dev tailwind flex center"
	if capturedQuery != expected {
		t.Errorf("query = %q, want %q", capturedQuery, expected)
	}
}

// TestBraveSearchAuthFailure verifies that a 401 response returns an error.
func TestBraveSearchAuthFailure(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer srv.Close()

	client := NewBraveClient("invalid-key", srv.URL)
	_, err := client.Search(context.Background(), "test query", nil, 0)
	if err == nil {
		t.Error("BraveClient.Search should return error on 401, got nil")
	}
}

// TestBraveSearchServerError verifies that a 500 response returns an error.
func TestBraveSearchServerError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer srv.Close()

	client := NewBraveClient("test-api-key", srv.URL)
	_, err := client.Search(context.Background(), "test query", nil, 0)
	if err == nil {
		t.Error("BraveClient.Search should return error on 500, got nil")
	}
}

// TestBraveSearchTimeout verifies that a slow server triggers a context deadline error
// when the client has a very short timeout.
func TestBraveSearchTimeout(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(100 * time.Millisecond)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"web": map[string]interface{}{
				"results": []map[string]interface{}{},
			},
		})
	}))
	defer srv.Close()

	client := NewBraveClient("test-api-key", srv.URL)
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
	defer cancel()

	_, err := client.Search(ctx, "test query", nil, 0)
	if err == nil {
		t.Fatal("BraveClient.Search should return error on timeout, got nil")
	}
	if ctx.Err() != context.DeadlineExceeded {
		t.Errorf("expected context deadline exceeded, got: %v", err)
	}
}
