package main

import "spacecraft/internal/mission"

// Type aliases for backward compatibility.
// All types are now defined in internal/mission/model.go.
// These aliases allow the existing package main code to compile without changes.

type GitBlock = mission.GitBlock
type ArtifactsBlock = mission.ArtifactsBlock
type ClarificationBlock = mission.ClarificationBlock
type Mission = mission.Mission
type MissionRecord = mission.MissionRecord
type Task = mission.Task
type Plan = mission.Plan
type EvidenceEntry = mission.EvidenceEntry
type GitInfoData = mission.GitInfoData
type MissionInfo = mission.MissionInfo
type SignalInfo = mission.SignalInfo
type ConflictInfo = mission.ConflictInfo
type CandidateInfo = mission.CandidateInfo
type GitInfo = mission.GitInfo
type ResolveOutput = mission.ResolveOutput
type TasksSummary = mission.TasksSummary
type WorkflowSnapshot = mission.WorkflowSnapshot
type ReleaseGate = mission.ReleaseGate
type ReleaseReadiness = mission.ReleaseReadiness
type Finding = mission.Finding
type Review = mission.Review
type CompactEvidenceEntry = mission.CompactEvidenceEntry
type CompactMission = mission.CompactMission
type CompactTask = mission.CompactTask
type CompactPlan = mission.CompactPlan
