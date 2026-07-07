package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
)

func resolveSafety(selected *MissionRecord, conflicts []ConflictInfo, ambiguous bool) string {
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

func authoritativeSignals(signals []SignalInfo) []SignalInfo {
	var res []SignalInfo
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

func signalConflicts(signals []SignalInfo, explicitSelector bool, selectedMissionId *string, selectedSource *string) []ConflictInfo {
	if explicitSelector {
		return []ConflictInfo{}
	}
	strongSignals := authoritativeSignals(signals)
	var conflicts []ConflictInfo
	selectedByStrongSignal := false
	if selectedSource != nil {
		selectedByStrongSignal = isAuthoritativeSignalSource(*selectedSource)
	}

	for _, sig := range strongSignals {
		if sig.ExpectedMissionId != "" && sig.MissionId == nil {
			conflicts = append(conflicts, ConflictInfo{
				Type:      "missing-signal-mission",
				Source:    sig.Source,
				MissionId: sig.ExpectedMissionId,
				Value:     sig.Value,
			})
		}
		if len(sig.MissionIds) > 1 {
			if !selectedByStrongSignal || selectedMissionId == nil || !containsStr(sig.MissionIds, *selectedMissionId) {
				conflicts = append(conflicts, ConflictInfo{
					Type:       "ambiguous-signal",
					Source:     sig.Source,
					Value:      sig.Value,
					MissionIds: sig.MissionIds,
				})
			}
		}
	}

	var resolvedSignals []SignalInfo
	for _, sig := range strongSignals {
		if sig.MissionId != nil {
			resolvedSignals = append(resolvedSignals, SignalInfo{
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
		conflicts = append(conflicts, ConflictInfo{
			Type:    "signal-mismatch",
			Signals: resolvedSignals,
		})
	}
	if conflicts == nil {
		return []ConflictInfo{}
	}
	return conflicts
}

func stringPtr(s string) *string {
	return &s
}
func intPtr(i int) *int {
	return &i
}

func missionSummary(r *MissionRecord, source *string) MissionInfo {
	return MissionInfo{
		ID:       r.ID,
		Title:    r.Mission.Title,
		State:    r.Mission.State,
		Active:   r.Active,
		Branches: r.Branches,
		Signal:   source,
	}
}

func candidateRecordsForResolution(records []MissionRecord, selected *MissionRecord, activeRecords []MissionRecord, signals []SignalInfo) []MissionRecord {
	candidateIds := make(map[string]bool)
	var result []MissionRecord

	add := func(r MissionRecord) {
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
			if rec := findMissionRecord(records, *sig.MissionId); rec != nil {
				add(*rec)
			}
		}
		for _, mid := range sig.MissionIds {
			if rec := findMissionRecord(records, mid); rec != nil {
				add(*rec)
			}
		}
	}
	for _, r := range activeRecords {
		add(r)
	}
	if len(result) == 0 && len(records) > 0 {
		ordered := missionDisplayRecords(records)
		for i := 0; i < len(ordered) && i < 10; i++ {
			add(ordered[i])
		}
	}
	return result
}

func resolveMission(selector string) ResolveOutput {
	records := listMissionRecords()
	orderedRecords := missionDisplayRecords(records)
	var signals []SignalInfo
	var conflicts []ConflictInfo
	var selected *MissionRecord
	var source *string
	ambiguous := false

	explicitSelector := selector
	if explicitSelector == "" {
		explicitSelector = os.Getenv("SPACECRAFT_MISSION")
	}

	selectFunc := func(record *MissionRecord, signal string) {
		if record == nil || selected != nil {
			return
		}
		if explicitSelector != "" && signal != "selector" && signal != "SPACECRAFT_MISSION" {
			return
		}
		selected = record
		source = stringPtr(signal)
	}

	if explicitSelector != "" {
		record := findMissionBySelector(records, explicitSelector, orderedRecords)
		src := "SPACECRAFT_MISSION"
		if selector != "" {
			src = "selector"
		}
		var mid *string
		if record != nil {
			mid = stringPtr(record.ID)
		}
		signals = append(signals, SignalInfo{
			Source:    src,
			Value:     explicitSelector,
			MissionId: mid,
		})
		if record == nil {
			ambiguous = true
		}
		selectFunc(record, src)
	}

	sessionMissionId := readSessionMissionId()
	if sessionMissionId != nil {
		record := findMissionRecord(records, *sessionMissionId)
		var mid *string
		if record != nil {
			mid = stringPtr(record.ID)
		}
		signals = append(signals, SignalInfo{
			Source:            "session",
			Value:             *sessionMissionId,
			ExpectedMissionId: *sessionMissionId,
			MissionId:         mid,
		})
		selectFunc(record, "session")
	}

	git := gitInfo()
	var branchMissionId *string
	if git.Branch != "" {
		branchMissionId = normalizeMissionId(git.Branch)
	}
	if branchMissionId != nil {
		record := findMissionRecord(records, *branchMissionId)
		var mid *string
		if record != nil {
			mid = stringPtr(record.ID)
		}
		signals = append(signals, SignalInfo{
			Source:            "branch",
			Value:             git.Branch,
			ExpectedMissionId: *branchMissionId,
			MissionId:         mid,
		})
		selectFunc(record, "branch")
	}

	var branchMetadataMatches []MissionRecord
	if git.Branch != "" {
		for _, r := range records {
			if containsStr(r.Branches, git.Branch) {
				branchMetadataMatches = append(branchMetadataMatches, r)
			}
		}
	}
	if len(branchMetadataMatches) == 1 {
		signals = append(signals, SignalInfo{
			Source:    "branch-metadata",
			Value:     git.Branch,
			MissionId: stringPtr(branchMetadataMatches[0].ID),
		})
		selectFunc(&branchMetadataMatches[0], "branch-metadata")
	} else if len(branchMetadataMatches) > 1 {
		var ids []string
		for _, r := range branchMetadataMatches {
			ids = append(ids, r.ID)
		}
		signals = append(signals, SignalInfo{
			Source:     "branch-metadata",
			Value:      git.Branch,
			MissionIds: ids,
		})
	}

	currentMissionId := readCurrentMissionId(false)
	if currentMissionId != nil {
		record := findMissionRecord(records, *currentMissionId)
		var mid *string
		if record != nil {
			mid = stringPtr(record.ID)
		}
		signals = append(signals, SignalInfo{
			Source:            ".space/current",
			Value:             *currentMissionId,
			ExpectedMissionId: *currentMissionId,
			MissionId:         mid,
		})
		selectFunc(record, ".space/current")
	}

	var activeRecords []MissionRecord
	for _, r := range records {
		if r.Active {
			activeRecords = append(activeRecords, r)
		}
	}
	if len(activeRecords) == 1 {
		signals = append(signals, SignalInfo{
			Source:    "single-active",
			Value:     activeRecords[0].ID,
			MissionId: stringPtr(activeRecords[0].ID),
		})
		selectFunc(&activeRecords[0], "single-active")
	} else if selected == nil && len(activeRecords) > 1 {
		ambiguous = true
	}

	var selId *string
	if selected != nil {
		selId = stringPtr(selected.ID)
	}
	cfs := signalConflicts(signals, explicitSelector != "", selId, source)
	conflicts = append(conflicts, cfs...)

	displayNumberById := make(map[string]int)
	for i, r := range orderedRecords {
		displayNumberById[r.ID] = i + 1
	}

	candidateRecords := candidateRecordsForResolution(records, selected, activeRecords, signals)

	var selInfo *MissionInfo
	if selected != nil {
		info := missionSummary(selected, source)
		selInfo = &info
	}

	var candidates []CandidateInfo
	for _, c := range candidateRecords {
		ci := CandidateInfo{
			MissionInfo: missionSummary(&c, nil),
		}
		if num, ok := displayNumberById[c.ID]; ok {
			ci.Number = intPtr(num)
		}
		candidates = append(candidates, ci)
	}

	if signals == nil {
		signals = []SignalInfo{}
	}
	if conflicts == nil {
		conflicts = []ConflictInfo{}
	}
	if candidates == nil {
		candidates = []CandidateInfo{}
	}

	return ResolveOutput{
		Selected:         selInfo,
		Source:           source,
		Safety:           resolveSafety(selected, conflicts, ambiguous),
		Signals:          signals,
		Conflicts:        conflicts,
		Candidates:       candidates,
		CurrentMissionId: currentMissionId,
		Git: GitInfo{
			Branch: git.Branch,
			Sha:    git.Sha,
			IsRepo: git.IsRepo,
		},
	}
}

func printResolvedMission(args []string) {
	selector := ""
	if len(args) > 0 && args[0] != "--json" {
		selector = args[0]
	}
	out := resolveMission(selector)
	isJson := false
	for _, arg := range args {
		if arg == "--json" {
			isJson = true
		}
	}
	if isJson {
		var buf bytes.Buffer
		enc := json.NewEncoder(&buf)
		enc.SetEscapeHTML(false)
		enc.SetIndent("", "  ")
		enc.Encode(out)
		fmt.Print(buf.String())
	} else {
		if out.Selected != nil {
			fmt.Printf("Mission: %s\n", out.Selected.ID)
		} else {
			fmt.Println("No mission resolved.")
			os.Exit(1)
		}
	}
}

func requireResolvedMission(commandName string) ResolveOutput {
	res := resolveMission("")
	if res.Safety != "safe" || res.Selected == nil {
		fail(formatResolutionBlock(res, commandName))
	}
	return res
}

func formatResolutionBlock(out ResolveOutput, context string) string {
	return "Resolution failed or blocked."
}
