// Package resolver provides mission resolution logic based on signals
// from the environment, git branch, session bindings, and current file.
package resolver

import (
	"fmt"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"spacecraft/internal/gitutil"
	"spacecraft/internal/mission"
)

// Resolver resolves the current active mission from multiple signal sources.
type Resolver struct {
	store  mission.MissionStore
	runner gitutil.CommandRunner
	getenv func(string) string
}

// New creates a new Resolver with the given store, command runner, and optional env func.
// If getenv is nil, os.Getenv is used.
func New(store mission.MissionStore, runner gitutil.CommandRunner, getenv func(string) string) *Resolver {
	if getenv == nil {
		getenv = os.Getenv
	}
	return &Resolver{
		store:  store,
		runner: runner,
		getenv: getenv,
	}
}

// RequireResolved returns the resolved mission or an error if resolution fails.
func (r *Resolver) RequireResolved(commandName string) (mission.ResolveOutput, error) {
	res := r.Resolve("")
	if res.Safety != "safe" || res.Selected == nil {
		return res, fmt.Errorf("resolution %s: safety=%s, selected=%v", commandName, res.Safety, res.Selected)
	}
	return res, nil
}

// Resolve resolves the current mission using the given selector string.
// An empty selector triggers automatic signal-based resolution.
func (r *Resolver) Resolve(selector string) mission.ResolveOutput {
	records := r.listMissionRecords()
	orderedRecords := displayRecords(records)
	var signals []mission.SignalInfo
	var conflicts []mission.ConflictInfo
	var selected *mission.MissionRecord
	var source *string
	ambiguous := false

	explicitSelector := selector
	if explicitSelector == "" {
		explicitSelector = r.getenv("SPACECRAFT_MISSION")
	}

	selectFunc := func(record *mission.MissionRecord, signal string) {
		if record == nil || selected != nil {
			return
		}
		if explicitSelector != "" && signal != "selector" && signal != "SPACECRAFT_MISSION" {
			return
		}
		selected = record
		s := signal
		source = &s
	}

	// 1. Explicit selector or SPACECRAFT_MISSION env
	if explicitSelector != "" {
		record := findMissionBySelector(records, explicitSelector, orderedRecords)
		src := "SPACECRAFT_MISSION"
		if selector != "" {
			src = "selector"
		}
		var mid *string
		if record != nil {
			mid = &record.ID
		}
		signals = append(signals, mission.SignalInfo{
			Source:    src,
			Value:     explicitSelector,
			MissionId: mid,
		})
		if record == nil {
			ambiguous = true
		}
		selectFunc(record, src)
	}

	// 2. Session binding
	if sessionKey := currentSessionKey(r.getenv); sessionKey != nil {
		sessionId, _ := r.store.ReadSession(*sessionKey)
		if sessionId != nil {
			record := findMissionRecordFn(records, *sessionId)
			var mid *string
			if record != nil {
				mid = &record.ID
			}
			signals = append(signals, mission.SignalInfo{
				Source:            "session",
				Value:             *sessionId,
				ExpectedMissionId: *sessionId,
				MissionId:         mid,
			})
			selectFunc(record, "session")
		}
	}

	// 3. Branch mission id (extracted from branch name)
	git := r.gitInfo()
	var branchMissionId *string
	if git.Branch != "" {
		branchMissionId = normalizeMissionId(git.Branch)
	}
	if branchMissionId != nil {
		record := findMissionRecordFn(records, *branchMissionId)
		var mid *string
		if record != nil {
			mid = &record.ID
		}
		signals = append(signals, mission.SignalInfo{
			Source:            "branch",
			Value:             git.Branch,
			ExpectedMissionId: *branchMissionId,
			MissionId:         mid,
		})
		selectFunc(record, "branch")
	}

	// 4. Branch metadata (branch stored in mission.json)
	var branchMetadataMatches []mission.MissionRecord
	if git.Branch != "" {
		for _, r := range records {
			if containsStr(r.Branches, git.Branch) {
				branchMetadataMatches = append(branchMetadataMatches, r)
			}
		}
	}
	if len(branchMetadataMatches) == 1 {
		signals = append(signals, mission.SignalInfo{
			Source:    "branch-metadata",
			Value:     git.Branch,
			MissionId: &branchMetadataMatches[0].ID,
		})
		selectFunc(&branchMetadataMatches[0], "branch-metadata")
	} else if len(branchMetadataMatches) > 1 {
		var ids []string
		for _, r := range branchMetadataMatches {
			ids = append(ids, r.ID)
		}
		signals = append(signals, mission.SignalInfo{
			Source:     "branch-metadata",
			Value:      git.Branch,
			MissionIds: ids,
		})
	}

	// 5. .space/current
	currentMissionId, _ := r.store.ReadCurrent()
	if currentMissionId != nil {
		record := findMissionRecordFn(records, *currentMissionId)
		var mid *string
		if record != nil {
			mid = &record.ID
		}
		signals = append(signals, mission.SignalInfo{
			Source:            ".space/current",
			Value:             *currentMissionId,
			ExpectedMissionId: *currentMissionId,
			MissionId:         mid,
		})
		selectFunc(record, ".space/current")
	}

	// 6. Single active mission fallback
	var activeRecords []mission.MissionRecord
	for _, r := range records {
		if r.Active {
			activeRecords = append(activeRecords, r)
		}
	}
	if len(activeRecords) == 1 {
		signals = append(signals, mission.SignalInfo{
			Source:    "single-active",
			Value:     activeRecords[0].ID,
			MissionId: &activeRecords[0].ID,
		})
		selectFunc(&activeRecords[0], "single-active")
	} else if selected == nil && len(activeRecords) > 1 {
		ambiguous = true
	}

	var selId *string
	if selected != nil {
		selId = &selected.ID
	}
	cfs := signalConflicts(signals, explicitSelector != "", selId, source)
	conflicts = append(conflicts, cfs...)

	displayNumberById := make(map[string]int)
	for i, r := range orderedRecords {
		displayNumberById[r.ID] = i + 1
	}

	candidateRecords := candidateRecordsForResolution(records, selected, activeRecords, signals)

	var selInfo *mission.MissionInfo
	if selected != nil {
		info := missionSummary(selected, source)
		selInfo = &info
	}

	var candidates []mission.CandidateInfo
	for _, c := range candidateRecords {
		ci := mission.CandidateInfo{
			MissionInfo: missionSummary(&c, nil),
		}
		if num, ok := displayNumberById[c.ID]; ok {
			ci.Number = &num
		}
		candidates = append(candidates, ci)
	}

	if signals == nil {
		signals = []mission.SignalInfo{}
	}
	if conflicts == nil {
		conflicts = []mission.ConflictInfo{}
	}
	if candidates == nil {
		candidates = []mission.CandidateInfo{}
	}

	return mission.ResolveOutput{
		Selected:         selInfo,
		Source:           source,
		Safety:           resolveSafety(selected, conflicts, ambiguous),
		Signals:          signals,
		Conflicts:        conflicts,
		Candidates:       candidates,
		CurrentMissionId: currentMissionId,
		Git: git,
	}
}

// --- internal helpers ---

func (r *Resolver) gitInfo() mission.GitInfoData {
	return gitutil.GitInfo(r.runner)
}

func (r *Resolver) listMissionRecords() []mission.MissionRecord {
	records, err := r.store.List()
	if err != nil {
		return nil
	}
	return records
}

// currentSessionKey returns the first non-empty env var for session identification.
func currentSessionKey(getenv func(string) string) *string {
	for _, env := range []string{"SPACECRAFT_SESSION", "OPENCODE_SESSION_ID", "CODEX_SESSION_ID"} {
		if val := getenv(env); val != "" {
			return &val
		}
	}
	return nil
}

// --- resolution helpers (extracted from original resolve.go) ---

func resolveSafety(selected *mission.MissionRecord, conflicts []mission.ConflictInfo, ambiguous bool) string {
	if len(conflicts) > 0 {
		return "conflict"
	}
	if ambiguous {
		return "ambiguous"
	}
	if selected == nil {
		return "none"
	}
	return "safe"
}

func isAuthoritativeSignalSource(source string) bool {
	return source == "session" || source == "branch" || source == "branch-metadata" || source == ".space/current"
}

func authoritativeSignals(signals []mission.SignalInfo) []mission.SignalInfo {
	var res []mission.SignalInfo
	for _, sig := range signals {
		if isAuthoritativeSignalSource(sig.Source) {
			res = append(res, sig)
		}
	}
	return res
}

func containsStr(list []string, item string) bool {
	for _, val := range list {
		if val == item {
			return true
		}
	}
	return false
}

func signalConflicts(signals []mission.SignalInfo, explicitSelector bool, selectedMissionId *string, selectedSource *string) []mission.ConflictInfo {
	if explicitSelector {
		return []mission.ConflictInfo{}
	}
	strongSignals := authoritativeSignals(signals)
	var conflicts []mission.ConflictInfo
	selectedByStrongSignal := false
	if selectedSource != nil {
		selectedByStrongSignal = isAuthoritativeSignalSource(*selectedSource)
	}

	for _, sig := range strongSignals {
		if sig.ExpectedMissionId != "" && sig.MissionId == nil {
			conflicts = append(conflicts, mission.ConflictInfo{
				Type:      "missing-signal-mission",
				Source:    sig.Source,
				MissionId: sig.ExpectedMissionId,
				Value:     sig.Value,
			})
		}
		if len(sig.MissionIds) > 1 {
			if !selectedByStrongSignal || selectedMissionId == nil || !containsStr(sig.MissionIds, *selectedMissionId) {
				conflicts = append(conflicts, mission.ConflictInfo{
					Type:       "ambiguous-signal",
					Source:     sig.Source,
					Value:      sig.Value,
					MissionIds: sig.MissionIds,
				})
			}
		}
	}

	var resolvedSignals []mission.SignalInfo
	for _, sig := range strongSignals {
		if sig.MissionId != nil {
			resolvedSignals = append(resolvedSignals, mission.SignalInfo{
				Source:    sig.Source,
				MissionId: sig.MissionId,
				Value:     sig.Value,
			})
		}
	}

	distinctMap := make(map[string]bool)
	for _, sig := range resolvedSignals {
		distinctMap[*sig.MissionId] = true
	}
	if len(distinctMap) > 1 {
		conflicts = append(conflicts, mission.ConflictInfo{
			Type:    "signal-mismatch",
			Signals: resolvedSignals,
		})
	}
	if conflicts == nil {
		return []mission.ConflictInfo{}
	}
	return conflicts
}

func missionSummary(r *mission.MissionRecord, source *string) mission.MissionInfo {
	return mission.MissionInfo{
		ID:       r.ID,
		Title:    r.Mission.Title,
		State:    r.Mission.State,
		Active:   r.Active,
		Branches: r.Branches,
		Signal:   source,
	}
}

func displayRecords(records []mission.MissionRecord) []mission.MissionRecord {
	res := make([]mission.MissionRecord, len(records))
	copy(res, records)
	sort.Slice(res, func(i, j int) bool {
		if res[i].Active != res[j].Active {
			return res[i].Active
		}
		return res[i].ID > res[j].ID
	})
	return res
}

func candidateRecordsForResolution(records []mission.MissionRecord, selected *mission.MissionRecord, activeRecords []mission.MissionRecord, signals []mission.SignalInfo) []mission.MissionRecord {
	candidateIds := make(map[string]bool)
	var result []mission.MissionRecord

	add := func(r mission.MissionRecord) {
		if !candidateIds[r.ID] {
			candidateIds[r.ID] = true
			result = append(result, r)
		}
	}
	if selected != nil {
		add(*selected)
	}
	for _, sig := range signals {
		if sig.MissionId != nil {
			if rec := findMissionRecordFn(records, *sig.MissionId); rec != nil {
				add(*rec)
			}
		}
		for _, mid := range sig.MissionIds {
			if rec := findMissionRecordFn(records, mid); rec != nil {
				add(*rec)
			}
		}
	}
	for _, r := range activeRecords {
		add(r)
	}
	if len(result) == 0 && len(records) > 0 {
		ordered := displayRecords(records)
		for i := 0; i < len(ordered) && i < 10; i++ {
			add(ordered[i])
		}
	}
	return result
}

// findMissionRecordFn is a standalone version for use in closures/packages.
func findMissionRecordFn(records []mission.MissionRecord, id string) *mission.MissionRecord {
	for i := range records {
		if records[i].ID == id {
			return &records[i]
		}
	}
	return nil
}

// findMissionBySelector finds a mission by number, id, exact title, or substring match.
func findMissionBySelector(records []mission.MissionRecord, selector string, orderedRecords []mission.MissionRecord) *mission.MissionRecord {
	text := strings.TrimSpace(selector)
	if text == "" {
		return nil
	}
	if num, err := strconv.Atoi(text); err == nil && num > 0 && num <= len(orderedRecords) {
		return &orderedRecords[num-1]
	}
	if id := normalizeMissionId(text); id != nil {
		return findMissionRecordFn(records, *id)
	}
	var exact []mission.MissionRecord
	for _, r := range records {
		if r.Mission != nil && r.Mission.Title == text {
			exact = append(exact, r)
		}
	}
	if len(exact) == 1 {
		return &exact[0]
	}
	var matches []mission.MissionRecord
	norm := strings.ToLower(text)
	for _, r := range records {
		if r.Mission != nil && strings.Contains(strings.ToLower(r.Mission.Title), norm) {
			matches = append(matches, r)
		}
	}
	if len(matches) == 1 {
		return &matches[0]
	}
	return nil
}

// --- ID normalization (moved from util.go to avoid circular dep) ---

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

// FormatResolutionBlock creates a brief error message for failed resolution.
func FormatResolutionBlock(out mission.ResolveOutput, context string) string {
	return "Resolution failed or blocked."
}
