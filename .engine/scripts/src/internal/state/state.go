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
	"spacecraft/internal/util"
)

// AllowedStates lists valid mission states.
var AllowedStates = map[string]bool{
	"draft": true, "planned": true, "built": true,
	"ready": true, "shipped": true, "blocked": true,
}

// ValidTransitions defines allowed state transitions.
// Key is current state, value is set of allowed next states.
var ValidTransitions = map[string]map[string]bool{
	"draft":   {"planned": true, "blocked": true},
	"planned": {"built": true, "blocked": true},
	"built":   {"ready": true, "blocked": true},
	"ready":   {"shipped": true, "blocked": true},
	"blocked": {"draft": true, "planned": true, "built": true, "ready": true},
	"shipped": {}, // terminal state
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

	// Validate state transition
	currentState := m.State
	if currentState != "" && currentState != newState {
		allowed, exists := ValidTransitions[currentState]
		if !exists {
			return fmt.Errorf("unknown current state %q", currentState)
		}
		if !allowed[newState] {
			return fmt.Errorf("invalid state transition from %q to %q", currentState, newState)
		}

		// Validate prerequisites for target state (only when actually transitioning)
		if err := s.validateStatePrerequisites(id, newState); err != nil {
			return err
		}
	}

	m.State = newState
	m.UpdatedAt = s.nowFn()
	return s.store.Save(m)
}

// validateStatePrerequisites checks that prerequisites are met for a state transition.
func (s *StateSetter) validateStatePrerequisites(id, targetState string) error {
	dir := s.store.MissionDir(id)

	switch targetState {
	case "planned":
		// Requires spec.md and plan.json
		if !util.Exists(filepath.Join(dir, "spec.md")) {
			return fmt.Errorf("cannot transition to planned: spec.md is missing")
		}
		if !util.Exists(filepath.Join(dir, "plan.json")) {
			return fmt.Errorf("cannot transition to planned: plan.json is missing")
		}

	case "built":
		// Requires all tasks in plan.json to be done
		pPath := filepath.Join(dir, "plan.json")
		if !util.Exists(pPath) {
			return fmt.Errorf("cannot transition to built: plan.json is missing")
		}
		var p mission.Plan
		if err := util.ReadJson(pPath, &p); err != nil {
			return fmt.Errorf("cannot transition to built: invalid plan.json: %w", err)
		}
		for _, t := range p.Tasks {
			if !mission.TaskIsComplete(t.Status) {
				taskID := "unknown"
				if t.ID != nil {
					taskID = *t.ID
				}
				return fmt.Errorf("cannot transition to built: task %s is not complete", taskID)
			}
		}

	case "ready":
		// Requires evidence.jsonl to have entries
		evPath := filepath.Join(dir, "evidence.jsonl")
		if !util.Exists(evPath) {
			return fmt.Errorf("cannot transition to ready: evidence.jsonl is missing")
		}
		content, err := os.ReadFile(evPath)
		if err != nil {
			return fmt.Errorf("cannot transition to ready: cannot read evidence.jsonl: %w", err)
		}
		lines := strings.Split(strings.TrimSpace(string(content)), "\n")
		validLines := 0
		for _, line := range lines {
			if strings.TrimSpace(line) != "" {
				validLines++
			}
		}
		if validLines == 0 {
			return fmt.Errorf("cannot transition to ready: evidence.jsonl has no entries")
		}

	case "shipped":
		// Requires review.json with status "ready", all tasks done, evidence exists
		revPath := filepath.Join(dir, "review.json")
		if !util.Exists(revPath) {
			return fmt.Errorf("cannot transition to shipped: review.json is missing")
		}
		var r mission.Review
		if err := util.ReadJson(revPath, &r); err != nil {
			return fmt.Errorf("cannot transition to shipped: invalid review.json: %w", err)
		}
		if r.Status == nil || *r.Status != "ready" {
			return fmt.Errorf("cannot transition to shipped: review status is not ready")
		}

		// Check all tasks are done
		pPath := filepath.Join(dir, "plan.json")
		if util.Exists(pPath) {
			var p mission.Plan
			if err := util.ReadJson(pPath, &p); err == nil {
				for _, t := range p.Tasks {
					if !mission.TaskIsComplete(t.Status) {
						taskID := "unknown"
						if t.ID != nil {
							taskID = *t.ID
						}
						return fmt.Errorf("cannot transition to shipped: task %s is not complete", taskID)
					}
				}
			}
		}

		// Check evidence exists
		evPath := filepath.Join(dir, "evidence.jsonl")
		if !util.Exists(evPath) {
			return fmt.Errorf("cannot transition to shipped: evidence.jsonl is missing")
		}
	}

	return nil
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
	if !util.Exists(mPath) {
		errors = append(errors, "missing mission.json")
	} else {
		var m mission.Mission
		if err := util.ReadJson(mPath, &m); err != nil {
			errors = append(errors, "invalid JSON in mission.json: "+err.Error())
		} else {
			status := m.Clarification.Status
			if status != "" && !AllowedClarificationStatuses[status] {
				errors = append(errors, "mission.json clarification.status must be one of: open, clear, deferred")
			}
		}
	}

	// Check spec.md
	if !util.Exists(filepath.Join(dir, "spec.md")) {
		errors = append(errors, "missing spec.md")
	}

	// Check plan.json
	pPath := filepath.Join(dir, "plan.json")
	if !util.Exists(pPath) {
		errors = append(errors, "missing plan.json")
	} else {
		var p mission.Plan
		if err := util.ReadJson(pPath, &p); err != nil {
			errors = append(errors, "invalid JSON in plan.json: "+err.Error())
		} else if p.Tasks == nil {
			errors = append(errors, "plan.json must contain a tasks array")
		}
	}

	// Check evidence.jsonl
	evPath := filepath.Join(dir, "evidence.jsonl")
	if !util.Exists(evPath) {
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
	if util.Exists(revPath) {
		var r interface{}
		if err := util.ReadJson(revPath, &r); err != nil {
			errors = append(errors, "invalid JSON in review.json: "+err.Error())
		}
	}

	if len(errors) > 0 {
		return &ValidationError{Errors: errors}
	}
	return nil
}
