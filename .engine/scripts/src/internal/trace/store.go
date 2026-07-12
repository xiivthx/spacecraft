package trace

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"spacecraft/internal/config"
	"spacecraft/internal/util"
)

type TraceStore interface {
	LoadTraces(missionId string) ([]TraceEntry, error)
	HasTraces(missionId string) bool
	ListMissionsWithTraces() ([]string, error)
}

type FSTraceStore struct {
	cfg *config.Config
}

func NewFSTraceStore(cfg *config.Config) *FSTraceStore {
	return &FSTraceStore{cfg: cfg}
}

func (s *FSTraceStore) traceFilePath(missionId string) string {
	return filepath.Join(s.cfg.TraceStoreDir(), missionId+".jsonl")
}

func (s *FSTraceStore) LoadTraces(missionId string) ([]TraceEntry, error) {
	path := s.traceFilePath(missionId)
	if !util.Exists(path) {
		return nil, nil
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var entries []TraceEntry
	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		var entry TraceEntry
		if err := json.Unmarshal([]byte(line), &entry); err != nil {
			continue
		}
		entries = append(entries, entry)
	}
	return entries, nil
}

func (s *FSTraceStore) HasTraces(missionId string) bool {
	entries, err := s.LoadTraces(missionId)
	if err != nil {
		return false
	}
	return len(entries) > 0
}

func (s *FSTraceStore) ListMissionsWithTraces() ([]string, error) {
	dir := s.cfg.TraceStoreDir()
	if !util.Exists(dir) {
		return nil, nil
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	var ids []string
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		if !strings.HasSuffix(name, ".jsonl") {
			continue
		}
		id := strings.TrimSuffix(name, ".jsonl")
		ids = append(ids, id)
	}
	sort.Strings(ids)
	return ids, nil
}
