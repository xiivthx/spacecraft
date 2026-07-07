// Package mission defines the core data model and store abstractions.
package mission

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"spacecraft/internal/config"
	"spacecraft/internal/id"
)

// MissionStore defines the data-access interface for mission artifacts.
// Implementations must be safe for concurrent read-only use.
type MissionStore interface {
	// ReadCurrent returns the mission id stored in .space/current, or nil.
	ReadCurrent() (*string, error)
	// WriteCurrent stores the mission id in .space/current.
	WriteCurrent(id string) error
	// ClearCurrent empties .space/current.
	ClearCurrent() error

	// ReadSession returns the session-bound mission id, or nil.
	ReadSession(sessionKey string) (*string, error)
	// WriteSession stores the mission id for the given session key.
	WriteSession(sessionKey, id string) error
	// ClearSession empties the session binding for the given key.
	ClearSession(sessionKey string) error

	// SessionDir returns the sessions directory path.
	SessionDir() string
	// SessionFilePath returns the session binding file path for a key.
	SessionFilePath(sessionKey string) string

	// List returns all mission records found on disk.
	List() ([]MissionRecord, error)

	// Load reads a mission.json for the given id.
	Load(id string) (*Mission, error)
	// Save writes a mission.json for the given id.
	Save(m *Mission) error
	// Create writes a new mission.json and ensures output/design dirs.
	Create(m *Mission) error

	// LoadPlan reads a plan.json for the given mission id.
	LoadPlan(id string) (*Plan, error)
	// SavePlan writes a plan.json for the given mission id.
	SavePlan(id string, p *Plan) error

	// LoadReview reads a review.json for the given mission id.
	LoadReview(id string) (*Review, error)
	// SaveReview writes a review.json for the given mission id.
	SaveReview(id string, r *Review) error

	// CountEvidence counts non-empty lines in evidence.jsonl.
	CountEvidence(id string) (int, error)
	// ReserveEvidencePath creates unique stdout/stderr output files and returns the
	// evidence id, stdout path, and stderr path. Caller must write to these files.
	ReserveEvidencePath(id string) (evidenceId string, stdoutPath string, stderrPath string, err error)
	// AppendEvidence appends a JSON evidence entry to evidence.jsonl.
	AppendEvidence(id string, entry *EvidenceEntry) error
	// ReadEvidenceEntries reads all entries from evidence.jsonl.
	ReadEvidenceEntries(id string) ([]EvidenceEntry, error)

	// Artifact existence checks:
	SpecExists(id string) bool
	PlanExists(id string) bool
	QuestionsExists(id string) bool
	DecisionsExists(id string) bool
	DesignExists(id string) bool
	ReviewJSONExists(id string) bool
	ReviewMDExists(id string) bool

	// MissionDir returns the filesystem path for a mission directory.
	MissionDir(id string) string
	// ReadFile reads the raw bytes of any file under the mission directory.
	ReadFile(id string, relPath string) ([]byte, error)
	// WriteFile writes raw bytes to a file under the mission directory.
	WriteFile(id string, relPath string, data []byte) error
	// RemoveAll removes the entire mission directory.
	RemoveAll(id string) error
	// CopyToArchive copies the mission directory to the archive dir.
	ArchiveMission(id, archiveDir string, compactM CompactMission, compactP CompactPlan, compactEvidence []CompactEvidenceEntry, review *Review) error
}

// FSStore implements MissionStore using the local filesystem.
type FSStore struct {
	cfg *config.Config
}

// NewFSStore creates a new FSStore backed by the given config.
func NewFSStore(cfg *config.Config) *FSStore {
	return &FSStore{cfg: cfg}
}

// --- current file ---

func (s *FSStore) ReadCurrent() (*string, error) {
	path := s.cfg.CurrentFile()
	if !fileExists(path) {
		return nil, nil
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	val := strings.TrimSpace(string(data))
	if val == "" {
		return nil, nil
	}
	return &val, nil
}

func (s *FSStore) WriteCurrent(id string) error {
	if err := os.MkdirAll(filepath.Dir(s.cfg.CurrentFile()), 0755); err != nil {
		return err
	}
	return os.WriteFile(s.cfg.CurrentFile(), []byte(id+"\n"), 0644)
}

func (s *FSStore) ClearCurrent() error {
	if err := os.MkdirAll(filepath.Dir(s.cfg.CurrentFile()), 0755); err != nil {
		return err
	}
	return os.WriteFile(s.cfg.CurrentFile(), []byte(""), 0644)
}

// --- session file ---

func (s *FSStore) SessionDir() string {
	return filepath.Join(s.cfg.SpaceDir(), "sessions")
}

func (s *FSStore) SessionFilePath(key string) string {
	safe := slugifySimple(key)
	if len(safe) > 80 {
		safe = safe[:80]
	}
	if safe == "" {
		return ""
	}
	return filepath.Join(s.SessionDir(), safe+".current")
}

func (s *FSStore) ReadSession(sessionKey string) (*string, error) {
	path := s.SessionFilePath(sessionKey)
	if path == "" || !fileExists(path) {
		return nil, nil
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	val := strings.TrimSpace(string(data))
	if val == "" {
		return nil, nil
	}
	return &val, nil
}

func (s *FSStore) WriteSession(sessionKey, id string) error {
	path := s.SessionFilePath(sessionKey)
	if path == "" {
		return fmt.Errorf("store: cannot derive session file path for key %q", sessionKey)
	}
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	return os.WriteFile(path, []byte(id+"\n"), 0644)
}

func (s *FSStore) ClearSession(sessionKey string) error {
	path := s.SessionFilePath(sessionKey)
	if path == "" {
		return nil
	}
	return os.WriteFile(path, []byte(""), 0644)
}

// --- list missions ---

func (s *FSStore) List() ([]MissionRecord, error) {
	dir := s.cfg.MissionsDir()
	if !fileExists(dir) {
		return nil, nil
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	var records []MissionRecord
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		id := entry.Name()
		mDir := filepath.Join(dir, id)
		mPath := filepath.Join(mDir, "mission.json")
		if !fileExists(mPath) {
			continue
		}
		var m Mission
		if err := readJSON(mPath, &m); err != nil {
			continue
		}
		records = append(records, MissionRecord{
			ID:       id,
			Mission:  &m,
			Dir:      mDir,
			Active:   missionActive(&m),
			Branches: missionBranchNames(&m),
		})
	}
	sort.Slice(records, func(i, j int) bool {
		return records[i].ID < records[j].ID
	})
	return records, nil
}

// --- load / save mission ---

func (s *FSStore) Load(id string) (*Mission, error) {
	path := filepath.Join(s.MissionDir(id), "mission.json")
	var m Mission
	if err := readJSON(path, &m); err != nil {
		return nil, err
	}
	return &m, nil
}

func (s *FSStore) Save(m *Mission) error {
	path := filepath.Join(s.MissionDir(m.ID), "mission.json")
	return writeJSON(path, m)
}

func (s *FSStore) Create(m *Mission) error {
	dir := s.MissionDir(m.ID)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Join(dir, "outputs"), 0755); err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Join(dir, "design"), 0755); err != nil {
		return err
	}
	return s.Save(m)
}

// --- plan ---

func (s *FSStore) LoadPlan(id string) (*Plan, error) {
	path := filepath.Join(s.MissionDir(id), "plan.json")
	var p Plan
	if err := readJSON(path, &p); err != nil {
		return nil, err
	}
	return &p, nil
}

func (s *FSStore) SavePlan(id string, p *Plan) error {
	path := filepath.Join(s.MissionDir(id), "plan.json")
	return writeJSON(path, p)
}

// --- review ---

func (s *FSStore) LoadReview(id string) (*Review, error) {
	path := filepath.Join(s.MissionDir(id), "review.json")
	var r Review
	if err := readJSON(path, &r); err != nil {
		return nil, err
	}
	return &r, nil
}

func (s *FSStore) SaveReview(id string, r *Review) error {
	path := filepath.Join(s.MissionDir(id), "review.json")
	return writeJSON(path, r)
}

// --- evidence ---

func (s *FSStore) CountEvidence(id string) (int, error) {
	path := filepath.Join(s.MissionDir(id), "evidence.jsonl")
	if !fileExists(path) {
		return 0, nil
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return 0, err
	}
	lines := strings.Split(string(data), "\n")
	count := 0
	for _, line := range lines {
		if strings.TrimSpace(line) != "" {
			count++
		}
	}
	return count, nil
}

func (s *FSStore) ReserveEvidencePath(missionID string) (string, string, string, error) {
	dir := s.MissionDir(missionID)
	outputsDir := filepath.Join(dir, "outputs")

	now := time.Now()
	base, err := id.EvidenceId(now)
	if err != nil {
		return "", "", "", err
	}
	candidate := base
	offset := 2

	for {
		stdoutPath := filepath.Join(outputsDir, candidate+".stdout.txt")
		stderrPath := filepath.Join(outputsDir, candidate+".stderr.txt")

		f, err := os.OpenFile(stdoutPath, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0644)
		if err == nil {
			f.Close()
			return candidate, stdoutPath, stderrPath, nil
		}
		if !os.IsExist(err) {
			return "", "", "", err
		}
		// Collision — increment by a few milliseconds to generate a new id
		candidate, err = id.EvidenceId(now.Add(time.Duration(offset) * time.Millisecond))
		if err != nil {
			return "", "", "", err
		}
		offset++
	}
}

func (s *FSStore) AppendEvidence(id string, entry *EvidenceEntry) error {
	path := filepath.Join(s.MissionDir(id), "evidence.jsonl")
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(false)
	if err := enc.Encode(entry); err != nil {
		return err
	}
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.Write(buf.Bytes())
	return err
}

func (s *FSStore) ReadEvidenceEntries(id string) ([]EvidenceEntry, error) {
	path := filepath.Join(s.MissionDir(id), "evidence.jsonl")
	if !fileExists(path) {
		return nil, nil
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var entries []EvidenceEntry
	lines := strings.Split(string(data), "\n")
	for i, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		var entry EvidenceEntry
		if err := json.Unmarshal([]byte(line), &entry); err != nil {
			entry.ID = fmt.Sprintf("invalid-line-%d", i+1)
			entry.Label = "Invalid evidence entry"
		}
		entries = append(entries, entry)
	}
	return entries, nil
}

// --- artifact checks ---

func (s *FSStore) SpecExists(id string) bool {
	return fileExists(filepath.Join(s.MissionDir(id), "spec.md"))
}

func (s *FSStore) PlanExists(id string) bool {
	return fileExists(filepath.Join(s.MissionDir(id), "plan.json"))
}

func (s *FSStore) QuestionsExists(id string) bool {
	return fileExists(filepath.Join(s.MissionDir(id), "questions.md"))
}

func (s *FSStore) DecisionsExists(id string) bool {
	return fileExists(filepath.Join(s.MissionDir(id), "decisions.md"))
}

func (s *FSStore) DesignExists(id string) bool {
	return fileExists(filepath.Join(s.MissionDir(id), "design"))
}

func (s *FSStore) ReviewJSONExists(id string) bool {
	return fileExists(filepath.Join(s.MissionDir(id), "review.json"))
}

func (s *FSStore) ReviewMDExists(id string) bool {
	return fileExists(filepath.Join(s.MissionDir(id), "review.md"))
}

// --- file operations ---

func (s *FSStore) MissionDir(id string) string {
	return s.cfg.MissionDir(id)
}

func (s *FSStore) ReadFile(id string, relPath string) ([]byte, error) {
	return os.ReadFile(filepath.Join(s.MissionDir(id), relPath))
}

func (s *FSStore) WriteFile(id string, relPath string, data []byte) error {
	return os.WriteFile(filepath.Join(s.MissionDir(id), relPath), data, 0644)
}

func (s *FSStore) RemoveAll(id string) error {
	return os.RemoveAll(s.MissionDir(id))
}

func (s *FSStore) ArchiveMission(id, archiveDir string, compactM CompactMission, compactP CompactPlan, compactEvidence []CompactEvidenceEntry, review *Review) error {
	if err := os.MkdirAll(archiveDir, 0755); err != nil {
		return err
	}
	dest := filepath.Join(archiveDir, id)
	if fileExists(dest) {
		return fmt.Errorf("archive already exists: %s", dest)
	}
	if err := os.MkdirAll(dest, 0755); err != nil {
		return err
	}

	// Write summary
	if err := writeJSON(filepath.Join(dest, "mission.json"), compactM); err != nil {
		return err
	}
	if err := writeJSON(filepath.Join(dest, "plan.json"), compactP); err != nil {
		return err
	}

	// Write compact evidence
	var evLines []string
	for _, e := range compactEvidence {
		b, _ := json.Marshal(e)
		evLines = append(evLines, string(b))
	}
	evOut := strings.Join(evLines, "\n")
	if len(evLines) > 0 {
		evOut += "\n"
	}
	if err := os.WriteFile(filepath.Join(dest, "evidence.jsonl"), []byte(evOut), 0644); err != nil {
		return err
	}

	if review != nil {
		if err := writeJSON(filepath.Join(dest, "review.json"), review); err != nil {
			return err
		}
	}
	return nil
}

// --- internal helpers ---

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func readJSON(path string, target interface{}) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, target)
}

func writeJSON(path string, data interface{}) error {
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", "  ")
	if err := enc.Encode(data); err != nil {
		return err
	}
	return os.WriteFile(path, buf.Bytes(), 0644)
}

func missionActive(m *Mission) bool {
	return m.State != "shipped"
}

func missionBranchNames(m *Mission) []string {
	var res []string
	if m.Branch != nil && *m.Branch != "" {
		res = append(res, *m.Branch)
	}
	if m.WorkBranch != nil && *m.WorkBranch != "" {
		res = append(res, *m.WorkBranch)
	}
	if m.Git.WorkBranch != nil && *m.Git.WorkBranch != "" {
		res = append(res, *m.Git.WorkBranch)
	}
	return res
}

// slugifySimple is a basic ASCII slug function (no external deps).
func slugifySimple(value string) string {
	text := strings.ToLower(strings.TrimSpace(value))
	var sb strings.Builder
	for _, r := range text {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') {
			sb.WriteRune(r)
		} else {
			sb.WriteRune('-')
		}
	}
	slug := strings.Join(strings.FieldsFunc(sb.String(), func(r rune) bool { return false }), "")
	// Collapse multiple hyphens
	parts := strings.FieldsFunc(sb.String(), func(r rune) bool {
		return r == '-'
	})
	if len(parts) == 0 {
		return ""
	}
	slug = strings.Join(parts, "-")
	slug = strings.Trim(slug, "-")
	if len(slug) > 80 {
		slug = slug[:80]
		slug = strings.TrimRight(slug, "-")
	}
	return slug
}
