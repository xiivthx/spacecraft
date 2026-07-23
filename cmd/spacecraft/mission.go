package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
)

// missionExists reports whether a mission directory exists on disk.
func missionExists(spaceDir, id string) bool {
	fi, err := os.Stat(missionDir(spaceDir, id))
	return err == nil && fi.IsDir()
}

// listMissionIDs returns mission IDs found under .space/missions, sorted descending.
func listMissionIDs(spaceDir string) []string {
	entries, err := os.ReadDir(filepath.Join(spaceDir, "missions"))
	if err != nil {
		return nil
	}
	var ids []string
	for _, e := range entries {
		if e.IsDir() {
			ids = append(ids, e.Name())
		}
	}
	sort.Sort(sort.Reverse(sort.StringSlice(ids)))
	return ids
}

// readMission loads a mission.json into a generic map, preserving unknown fields.
func readMission(spaceDir, id string) (map[string]any, error) {
	data, err := os.ReadFile(filepath.Join(missionDir(spaceDir, id), "mission.json"))
	if err != nil {
		return nil, err
	}
	var m map[string]any
	if err := json.Unmarshal(data, &m); err != nil {
		return nil, err
	}
	return m, nil
}

// currentFile is the path to the pointer tracking the selected mission.
func currentFile(spaceDir string) string { return filepath.Join(spaceDir, "current") }

// readCurrent returns the current mission ID, if any.
func readCurrent(spaceDir string) string {
	data, err := os.ReadFile(currentFile(spaceDir))
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(data))
}

// writeCurrent records the current mission ID.
func writeCurrent(spaceDir, id string) error {
	if err := os.MkdirAll(spaceDir, 0755); err != nil {
		return err
	}
	return os.WriteFile(currentFile(spaceDir), []byte(id+"\n"), 0644)
}

// clearCurrent removes the current-mission pointer if present.
func clearCurrent(spaceDir string) {
	_ = os.Remove(currentFile(spaceDir))
}

// newMissionID generates a compact sortable mission ID (M + base36 ms since 2026-01-01).
func newMissionID() string {
	epoch := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	ms := time.Now().UTC().Sub(epoch).Milliseconds()
	if ms < 0 {
		ms = 0
	}
	return "M" + strings.ToUpper(strconv.FormatInt(ms, 36))
}

// normalizeID upper-cases a mission selector so lookups are case-insensitive.
func normalizeID(sel string) string {
	return strings.ToUpper(strings.TrimSpace(sel))
}
