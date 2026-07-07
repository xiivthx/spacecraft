// Package util provides shared pure helpers for path, JSON, and string operations.
package util

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
)

// ReadJson reads a JSON file and unmarshals it into target.
func ReadJson(filePath string, target interface{}) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, target)
}

// WriteJson marshals data as indented JSON and writes to filePath.
func WriteJson(filePath string, data interface{}) error {
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", "  ")
	if err := enc.Encode(data); err != nil {
		return err
	}
	return os.WriteFile(filePath, buf.Bytes(), 0644)
}

// Exists returns true if the path exists on disk.
func Exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// DisplayPath returns a human-readable relative path from root.
// If filePath is under root, returns the relative path; otherwise returns filePath.
func DisplayPath(root, filePath string) string {
	rel, err := filepath.Rel(root, filePath)
	if err != nil || rel == "" {
		return filePath
	}
	return rel
}

// CountEvidence counts non-empty lines in the evidence file.
func CountEvidence(filePath string) int {
	if !Exists(filePath) {
		return 0
	}
	content, err := os.ReadFile(filePath)
	if err != nil {
		return 0
	}
	lines := strings.Split(string(content), "\n")
	count := 0
	for _, line := range lines {
		if strings.TrimSpace(line) != "" {
			count++
		}
	}
	return count
}

// EnsureFile creates an empty file at path if it does not exist.
func EnsureFile(path string) error {
	if Exists(path) {
		return nil
	}
	f, err := os.OpenFile(path, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0644)
	if err != nil {
		if os.IsExist(err) {
			return nil
		}
		return err
	}
	return f.Close()
}
