// Package state provides mission state management, clarification status,
// and artifact validation — all returning errors instead of calling os.Exit.
package state

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"spacecraft/internal/mission"
)

// AllowedStates lists valid mission states.
var AllowedStates = map[string]bool{
	"draft": true, "planned": true, "built": true,
	"ready": true, "shipped": true, "blocked": true,
}

// AllowedClarificationStatuses lists valid clarification statuses.
var AllowedClarificationStatuses = map[string]bool{
	"open": true, "clear": true, "deferred": true,
}

// StateSetter provides state mutation operations on a mission.
type StateSetter struct {
	store  mission.MissionStore
	nowFn  func() string
}

// NewSetter creates a new StateSetter backed by the given store.
func NewSetter(store mission.MissionStore) *StateSetter {
	return &StateSetter{
		store: store,
		nowFn: func() string {
			return time.Now().UTC().Format("2006-01-02T15:04:05.000Z")
		},
	}
}

// SetState updates the state of a mission. Returns an error if the state is invalid.
// Legacy state names (specified, implementing) are mapped to current values.
func (s *StateSetter) SetState(id, newState string) error {
	if !AllowedStates[newState] {
		// Backward compat: map removed states
		switch newState {
		case "specified":
			newState = "draft"
		case "implementing":
			newState = "planned"
		case "verifying":
			newState = "built"
		case "reviewing":
			newState = "built"
		default:
			return fmt.Errorf("invalid mission state %q: allowed: draft, planned, built, ready, shipped, blocked", newState)
		}
	}

	m, err := s.store.Load(id)
	if err != nil {
		return fmt.Errorf("load mission %s: %w", id, err)
	}
	m.State = newState
	m.UpdatedAt = s.nowFn()
	return s.store.Save(m)
}

// SetClarificationStatus updates the clarification status of a mission.
func (s *StateSetter) SetClarificationStatus(id, status string) error {
	if !AllowedClarificationStatuses[status] {
		return fmt.Errorf("invalid clarification status %q: allowed: open, clear, deferred", status)
	}

	m, err := s.store.Load(id)
	if err != nil {
		return fmt.Errorf("load mission %s: %w", id, err)
	}
	m.Clarification.Status = status
	if status == "clear" {
		m.Clarification.BlockingQuestions = 0
		m.Clarification.LastQuestion = nil
	}
	m.UpdatedAt = s.nowFn()
	return s.store.Save(m)
}

// ValidationError holds a collection of validation errors.
type ValidationError struct {
	Errors []string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation: %d error(s)", len(e.Errors))
}

// ValidateMission validates all artifacts for a mission.
func (s *StateSetter) ValidateMission(id string) *ValidationError {
	var errors []string

	dir := s.store.MissionDir(id)

	// Check mission.json
	mPath := filepath.Join(dir, "mission.json")
	if !fileExists(mPath) {
		errors = append(errors, "missing mission.json")
	} else {
		var m mission.Mission
		if err := readJSON(mPath, &m); err != nil {
			errors = append(errors, "invalid JSON in mission.json: "+err.Error())
		} else {
			status := m.Clarification.Status
			if status != "" && !AllowedClarificationStatuses[status] {
				errors = append(errors, "mission.json clarification.status must be one of: open, clear, deferred")
			}
		}
	}

	// Check spec.md
	if !fileExists(filepath.Join(dir, "spec.md")) {
		errors = append(errors, "missing spec.md")
	}

	// Check plan.json
	pPath := filepath.Join(dir, "plan.json")
	if !fileExists(pPath) {
		errors = append(errors, "missing plan.json")
	} else {
		var p mission.Plan
		if err := readJSON(pPath, &p); err != nil {
			errors = append(errors, "invalid JSON in plan.json: "+err.Error())
		} else if p.Tasks == nil {
			errors = append(errors, "plan.json must contain a tasks array")
		}
	}

	// Check evidence.jsonl
	evPath := filepath.Join(dir, "evidence.jsonl")
	if !fileExists(evPath) {
		errors = append(errors, "missing evidence.jsonl")
	} else {
		content, err := os.ReadFile(evPath)
		if err == nil {
			lines := strings.Split(string(content), "\n")
			evIds := make(map[string]int)
			for i, line := range lines {
				if strings.TrimSpace(line) == "" {
					continue
				}
				var entry map[string]interface{}
				if err := json.Unmarshal([]byte(line), &entry); err != nil {
					errors = append(errors, fmt.Sprintf("invalid JSON in evidence.jsonl line %d: %s", i+1, err.Error()))
					continue
				}
				idVal, ok := entry["id"].(string)
				if !ok || idVal == "" {
					errors = append(errors, fmt.Sprintf("evidence.jsonl line %d must have string id", i+1))
				} else if prev, exists := evIds[idVal]; exists {
					errors = append(errors, fmt.Sprintf("duplicate evidence id %s on lines %d and %d", idVal, prev, i+1))
				} else {
					evIds[idVal] = i + 1
				}
			}
		}
	}

	// Check review.json (optional but validate if present)
	revPath := filepath.Join(dir, "review.json")
	if fileExists(revPath) {
		var r interface{}
		if err := readJSON(revPath, &r); err != nil {
			errors = append(errors, "invalid JSON in review.json: "+err.Error())
		}
	}

	if len(errors) > 0 {
		return &ValidationError{Errors: errors}
	}
	return nil
}

// --- helpers ---

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
