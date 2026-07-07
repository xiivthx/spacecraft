package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

var (
	ROOT         string
	SPACE_DIR    string
	MISSIONS_DIR string
	ARCHIVE_DIR  string
	CURRENT_FILE string
)

func init() {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	ROOT = cwd
	SPACE_DIR = filepath.Join(ROOT, ".space")
	MISSIONS_DIR = filepath.Join(SPACE_DIR, "missions")
	ARCHIVE_DIR = filepath.Join(SPACE_DIR, "archive")
	CURRENT_FILE = filepath.Join(SPACE_DIR, "current")
}

func fail(message string, code ...int) {
	exitCode := 1
	if len(code) > 0 {
		exitCode = code[0]
	}
	fmt.Fprintln(os.Stderr, message)
	os.Exit(exitCode)
}

func exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func ensureCurrentFile() error {
	if !exists(CURRENT_FILE) {
		f, err := os.OpenFile(CURRENT_FILE, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0644)
		if err != nil {
			if os.IsExist(err) {
				return nil
			}
			return err
		}
		f.Close()
	}
	return nil
}

func displayPath(filePath string) string {
	rel, err := filepath.Rel(ROOT, filePath)
	if err != nil || rel == "" {
		return "."
	}
	return rel
}

func readJson(filePath string, target interface{}) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, target)
}

func writeJson(filePath string, data interface{}) error {
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", "  ")
	if err := enc.Encode(data); err != nil {
		return err
	}
	// enc.Encode adds a trailing newline, so we don't need to append one
	return os.WriteFile(filePath, buf.Bytes(), 0644)
}

func isoNow() string {
	return time.Now().UTC().Format("2006-01-02T15:04:05.000Z")
}

var legacyMissionRegex = regexp.MustCompile(`\b[Mm]-(\d{8}-\d{6})\b`)
var compactMissionRegex = regexp.MustCompile(`(?:^|[^A-Za-z0-9])([Mm][0-9A-Za-z]{8})(?:$|[^A-Za-z0-9])`)

func normalizeMissionId(value string) *string {
	text := strings.TrimSpace(value)
	if legacy := legacyMissionRegex.FindStringSubmatch(text); legacy != nil {
		res := "M-" + legacy[1]
		return &res
	}
	if compact := compactMissionRegex.FindStringSubmatch(text); compact != nil {
		res := strings.ToUpper(compact[1])
		return &res
	}
	return nil
}

func regexpReplace(pattern, repl, src string) string {
	re := regexp.MustCompile(pattern)
	return re.ReplaceAllString(src, repl)
}

func slugify(value string) string {
	// Simple standard lib ASCII slugify mapping JS behavior broadly
	text := strings.ToLower(strings.TrimSpace(value))
	var sb strings.Builder
	for _, r := range text {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') {
			sb.WriteRune(r)
		} else {
			sb.WriteRune('-')
		}
	}
	slug := regexp.MustCompile(`-+`).ReplaceAllString(sb.String(), "-")
	slug = strings.Trim(slug, "-")
	if len(slug) > 60 {
		slug = slug[:60]
		slug = strings.TrimRight(slug, "-")
	}
	if slug == "" {
		return "mission"
	}
	return slug
}

func commandToString(parts []string) string {
	safeRegex := regexp.MustCompile(`^[A-Za-z0-9_./:=@%+-]+$`)
	var result []string
	for _, part := range parts {
		if safeRegex.MatchString(part) {
			result = append(result, part)
		} else {
			escaped := strings.ReplaceAll(part, "'", "'\\''")
			result = append(result, "'"+escaped+"'")
		}
	}
	return strings.Join(result, " ")
}
