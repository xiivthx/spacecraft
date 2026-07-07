package util

import (
	"os"
	"path/filepath"
	"testing"
)

func TestExists(t *testing.T) {
	// Should return false for non-existent path
	if Exists("/nonexistent/path/xyz123") {
		t.Error("Exists should return false for non-existent path")
	}

	// Should return true for current directory
	if !Exists(".") {
		t.Error("Exists should return true for '.'")
	}
}

func TestDisplayPath(t *testing.T) {
	tests := []struct {
		root, filePath, want string
	}{
		{"/a/b", "/a/b/c/d.txt", "c/d.txt"},
		{"/a/b", "/other/path.txt", "../../other/path.txt"},
		{"/a/b", "/a/b", "."},
	}
	for _, tt := range tests {
		got := DisplayPath(tt.root, tt.filePath)
		if got != tt.want {
			t.Errorf("DisplayPath(%q, %q) = %q, want %q", tt.root, tt.filePath, got, tt.want)
		}
	}
}

func TestWriteReadJson(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.json")

	type Data struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}

	orig := Data{Name: "test", Age: 42}
	if err := WriteJson(path, &orig); err != nil {
		t.Fatalf("WriteJson failed: %v", err)
	}

	if !Exists(path) {
		t.Fatal("WriteJson should create file")
	}

	var got Data
	if err := ReadJson(path, &got); err != nil {
		t.Fatalf("ReadJson failed: %v", err)
	}
	if got.Name != "test" || got.Age != 42 {
		t.Errorf("Read back = %+v, want {test 42}", got)
	}
}

func TestWriteJsonIndent(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "indent.json")

	data := map[string]int{"a": 1}
	if err := WriteJson(path, data); err != nil {
		t.Fatal(err)
	}

	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}

	if string(content) != "{\n  \"a\": 1\n}\n" {
		t.Errorf("unexpected content: %q", string(content))
	}
}

func TestCountEvidence(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "evidence.jsonl")

	// Empty file
	if n := CountEvidence(path); n != 0 {
		t.Errorf("CountEvidence empty = %d, want 0", n)
	}

	// Write some lines
	content := "{\"id\":\"E01\"}\n{\"id\":\"E02\"}\n\n{\"id\":\"E03\"}\n"
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	if n := CountEvidence(path); n != 3 {
		t.Errorf("CountEvidence = %d, want 3", n)
	}
}

func TestEnsureFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "newfile")

	if err := EnsureFile(path); err != nil {
		t.Fatalf("EnsureFile failed: %v", err)
	}
	if !Exists(path) {
		t.Error("EnsureFile should create file")
	}

	// Calling again should be a no-op
	if err := EnsureFile(path); err != nil {
		t.Fatalf("EnsureFile on existing file failed: %v", err)
	}
}
