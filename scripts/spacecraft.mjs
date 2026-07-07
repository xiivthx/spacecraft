#!/usr/bin/env node
import { spawn } from "node:child_process";
import { constants as fsConstants, existsSync } from "node:fs";
import fs from "node:fs/promises";
import path from "node:path";
import process from "node:process";

const ROOT = process.cwd();
const SPACE_DIR = path.join(ROOT, ".space");
const MISSIONS_DIR = path.join(SPACE_DIR, "missions");
const ARCHIVE_DIR = path.join(SPACE_DIR, "archive");
const CURRENT_FILE = path.join(SPACE_DIR, "current");
const ID_EPOCH_MS = Date.UTC(2026, 0, 1);
const SHORT_ID_WIDTH = 8;

const STATES = new Set([
  "draft",
  "specified",
  "planned",
  "implementing",
  "verifying",
  "reviewing",
  "ready",
  "shipped",
  "blocked"
]);

const CLARIFICATION_STATUSES = new Set(["open", "clear", "deferred"]);

function usage() {
  return `Spacecraft local mission helper

Usage:
  node scripts/spacecraft.mjs init
  node scripts/spacecraft.mjs new <title>
  node scripts/spacecraft.mjs current
  node scripts/spacecraft.mjs resolve [selector] [--json]
  node scripts/spacecraft.mjs missions
  node scripts/spacecraft.mjs use <number|id|title>
  node scripts/spacecraft.mjs bind-branch [selector]
  node scripts/spacecraft.mjs status
  node scripts/spacecraft.mjs flow [--json]
  node scripts/spacecraft.mjs git-info
  node scripts/spacecraft.mjs git-suggest [type] [slug]
  node scripts/spacecraft.mjs set-state <state>
  node scripts/spacecraft.mjs clarify-status <open|clear|deferred>
  node scripts/spacecraft.mjs evidence <label> -- <command...>
  node scripts/spacecraft.mjs validate
  node scripts/spacecraft.mjs closeout-check
  node scripts/spacecraft.mjs archive [selector]
`;
}

function fail(message, code = 1) {
  const error = new Error(message);
  error.exitCode = code;
  throw error;
}

function shortTimeId(prefix, date = new Date()) {
  const offset = date.getTime() - ID_EPOCH_MS;
  if (!Number.isFinite(offset) || offset < 0) {
    fail(`Cannot create ${prefix} id before 2026-01-01T00:00:00.000Z.`);
  }
  const encoded = Math.floor(offset).toString(36).toUpperCase();
  return `${prefix}${encoded.padStart(SHORT_ID_WIDTH, "0")}`;
}

function missionId(date = new Date()) {
  return shortTimeId("M", date);
}

function evidenceId(date = new Date()) {
  return shortTimeId("E", date);
}

function isoNow() {
  return new Date().toISOString();
}

async function exists(filePath) {
  try {
    await fs.access(filePath, fsConstants.F_OK);
    return true;
  } catch {
    return false;
  }
}

async function ensureCurrentFile() {
  if (!(await exists(CURRENT_FILE))) {
    await fs.writeFile(CURRENT_FILE, "", { flag: "wx" });
  }
}

async function initSpacecraft({ silent = false } = {}) {
  await fs.mkdir(MISSIONS_DIR, { recursive: true });
  await ensureCurrentFile();
  if (!silent) {
    console.log("Spacecraft initialized at .space/");
  }
}

async function readCurrentMissionId({ required = false } = {}) {
  if (!(await exists(CURRENT_FILE))) {
    if (required) {
      fail("No current Spacecraft mission. Start one with /sc-start <title>.");
    }
    return null;
  }

  const value = (await fs.readFile(CURRENT_FILE, "utf8")).trim();
  if (!value) {
    if (required) {
      fail("No current Spacecraft mission. Start one with /sc-start <title>.");
    }
    return null;
  }
  return normalizeMissionId(value) ?? value;
}

function missionDir(id) {
  return path.join(MISSIONS_DIR, id);
}

function artifactPath(id, artifact) {
  return path.join(missionDir(id), artifact);
}

function displayPath(filePath) {
  return path.relative(ROOT, filePath) || ".";
}

async function readJson(filePath) {
  const content = await fs.readFile(filePath, "utf8");
  return JSON.parse(content);
}

async function readJsonIfExists(filePath) {
  if (!(await exists(filePath))) {
    return null;
  }
  return readJson(filePath);
}

async function writeJson(filePath, data) {
  await fs.writeFile(filePath, `${JSON.stringify(data, null, 2)}\n`);
}

function normalizeMissionId(value) {
  const text = String(value ?? "");
  const legacy = text.match(/\b[Mm]-(\d{8}-\d{6})\b/);
  if (legacy) {
    return `M-${legacy[1]}`;
  }
  const compact = text.match(/(?:^|[^A-Za-z0-9])([Mm][0-9A-Za-z]{8})(?=$|[^A-Za-z0-9])/);
  return compact ? compact[1].toUpperCase() : null;
}

function missionBranchNames(mission) {
  const candidates = [
    mission?.branch,
    mission?.workBranch,
    mission?.git?.workBranch
  ];
  return candidates.filter((value) => typeof value === "string" && value.trim());
}

function missionActive(mission) {
  return mission?.state !== "shipped";
}

async function listMissionRecords() {
  if (!(await exists(MISSIONS_DIR))) {
    return [];
  }

  const entries = await fs.readdir(MISSIONS_DIR, { withFileTypes: true });
  const records = [];
  for (const entry of entries) {
    if (!entry.isDirectory()) {
      continue;
    }
    const id = entry.name;
    const missionPath = path.join(MISSIONS_DIR, id, "mission.json");
    if (!(await exists(missionPath))) {
      continue;
    }
    const mission = await readJson(missionPath);
    records.push({
      id,
      mission,
      dir: missionDir(id),
      active: missionActive(mission),
      branches: missionBranchNames(mission)
    });
  }
  return records.sort((a, b) => a.id.localeCompare(b.id));
}

function missionSummary(record, signal = null) {
  return {
    id: record.id,
    title: displayMissionTitle(record.mission?.title),
    state: record.mission?.state ?? "unknown",
    active: record.active,
    branches: record.branches,
    signal
  };
}

function displayMissionTitle(title) {
  const text = String(title ?? "(untitled)").replace(/\s+/g, " ").trim();
  if (text.length <= 88) {
    return text || "(untitled)";
  }
  return `${text.slice(0, 85)}...`;
}

function findMissionRecord(records, id) {
  return records.find((record) => record.id === id) ?? null;
}

function missionDisplayRecords(records) {
  return [...records].sort((a, b) => {
    if (a.active !== b.active) {
      return a.active ? -1 : 1;
    }
    return b.id.localeCompare(a.id);
  });
}

function findMissionBySelector(records, selector, orderedRecords = missionDisplayRecords(records)) {
  const text = String(selector ?? "").trim();
  if (!text) {
    return null;
  }

  if (/^\d+$/.test(text)) {
    const index = Number.parseInt(text, 10) - 1;
    return orderedRecords[index] ?? null;
  }

  const id = normalizeMissionId(text);
  if (id) {
    return findMissionRecord(records, id);
  }

  const exactTitle = records.filter((record) => record.mission?.title === text);
  if (exactTitle.length === 1) {
    return exactTitle[0];
  }

  const normalized = text.toLowerCase();
  const titleMatches = records.filter((record) => String(record.mission?.title ?? "").toLowerCase().includes(normalized));
  return titleMatches.length === 1 ? titleMatches[0] : null;
}

function currentSessionKey() {
  return process.env.SPACECRAFT_SESSION
    || process.env.OPENCODE_SESSION_ID
    || process.env.CODEX_SESSION_ID
    || null;
}

function sessionFilePath() {
  const key = currentSessionKey();
  if (!key) {
    return null;
  }
  const safeKey = slugify(key).slice(0, 80);
  if (!safeKey) {
    return null;
  }
  return path.join(SPACE_DIR, "sessions", `${safeKey}.current`);
}

async function readSessionMissionId() {
  const sessionFile = sessionFilePath();
  if (!sessionFile) {
    return null;
  }
  if (!(await exists(sessionFile))) {
    return null;
  }
  return normalizeMissionId(await fs.readFile(sessionFile, "utf8"));
}

async function writeSessionMissionId(id) {
  const sessionFile = sessionFilePath();
  if (!sessionFile) {
    return null;
  }
  await fs.mkdir(path.dirname(sessionFile), { recursive: true });
  await fs.writeFile(sessionFile, `${id}\n`);
  return sessionFile;
}

function resolveSafety(selected, conflicts, ambiguous) {
  if (conflicts.length > 0) {
    return "conflict";
  }
  if (ambiguous) {
    return "ambiguous";
  }
  if (!selected) {
    return "none";
  }
  return "safe";
}

const AUTHORITATIVE_SIGNAL_SOURCES = new Set([
  "session",
  "branch",
  "branch-metadata",
  ".space/current"
]);

function authoritativeSignals(signals) {
  return signals.filter((signal) => AUTHORITATIVE_SIGNAL_SOURCES.has(signal.source));
}

function signalConflicts(signals, explicitSelector, selectedMissionId = null, selectedSource = null) {
  if (explicitSelector) {
    return [];
  }

  const strongSignals = authoritativeSignals(signals);
  const conflicts = [];
  const selectedByStrongSignal = AUTHORITATIVE_SIGNAL_SOURCES.has(selectedSource);
  for (const signal of strongSignals) {
    if (signal.expectedMissionId && !signal.missionId) {
      conflicts.push({
        type: "missing-signal-mission",
        source: signal.source,
        missionId: signal.expectedMissionId,
        value: signal.value
      });
    }
    if (Array.isArray(signal.missionIds) && signal.missionIds.length > 1) {
      if (!selectedByStrongSignal || !selectedMissionId || !signal.missionIds.includes(selectedMissionId)) {
        conflicts.push({
          type: "ambiguous-signal",
          source: signal.source,
          value: signal.value,
          missionIds: signal.missionIds
        });
      }
    }
  }

  const resolvedSignals = strongSignals
    .filter((signal) => signal.missionId)
    .map((signal) => ({
      source: signal.source,
      missionId: signal.missionId,
      value: signal.value
    }));
  const distinctMissionIds = new Set(resolvedSignals.map((signal) => signal.missionId));
  if (distinctMissionIds.size > 1) {
    conflicts.push({
      type: "signal-mismatch",
      signals: resolvedSignals
    });
  }

  return conflicts;
}

function candidateRecordsForResolution(records, selected, activeRecords, signals) {
  const candidateIds = new Set();
  if (selected) {
    candidateIds.add(selected.id);
  }
  for (const signal of authoritativeSignals(signals)) {
    if (signal.missionId) {
      candidateIds.add(signal.missionId);
    }
    if (Array.isArray(signal.missionIds)) {
      for (const id of signal.missionIds) {
        candidateIds.add(id);
      }
    }
  }
  for (const record of activeRecords) {
    candidateIds.add(record.id);
  }

  return missionDisplayRecords(records).filter((record) => candidateIds.has(record.id));
}

async function resolveMission({ selector = null } = {}) {
  const records = await listMissionRecords();
  const signals = [];
  const conflicts = [];
  let selected = null;
  let source = null;
  let ambiguous = false;
  const explicitSelector = selector || process.env.SPACECRAFT_MISSION || null;
  const explicitSources = new Set(["selector", "SPACECRAFT_MISSION"]);

  function select(record, signal) {
    if (!record || selected) {
      return;
    }
    if (explicitSelector && !explicitSources.has(signal)) {
      return;
    }
    selected = record;
    source = signal;
  }

  if (explicitSelector) {
    const record = findMissionBySelector(records, explicitSelector);
    signals.push({ source: selector ? "selector" : "SPACECRAFT_MISSION", value: explicitSelector, missionId: record?.id ?? null });
    if (!record) {
      ambiguous = true;
    }
    select(record, selector ? "selector" : "SPACECRAFT_MISSION");
  }

  const sessionMissionId = await readSessionMissionId();
  const sessionRecord = sessionMissionId ? findMissionRecord(records, sessionMissionId) : null;
  if (sessionMissionId) {
    signals.push({ source: "session", value: sessionMissionId, expectedMissionId: sessionMissionId, missionId: sessionRecord?.id ?? null });
    select(sessionRecord, "session");
  }

  const git = await gitInfo();
  const branchMissionId = git.branch ? normalizeMissionId(git.branch) : null;
  const branchRecord = branchMissionId ? findMissionRecord(records, branchMissionId) : null;
  if (branchMissionId) {
    signals.push({ source: "branch", value: git.branch, expectedMissionId: branchMissionId, missionId: branchRecord?.id ?? null });
    select(branchRecord, "branch");
  }

  const branchMetadataMatches = git.branch
    ? records.filter((record) => record.branches.includes(git.branch))
    : [];
  if (branchMetadataMatches.length === 1) {
    signals.push({ source: "branch-metadata", value: git.branch, missionId: branchMetadataMatches[0].id });
    select(branchMetadataMatches[0], "branch-metadata");
  } else if (branchMetadataMatches.length > 1) {
    signals.push({ source: "branch-metadata", value: git.branch, missionId: null, missionIds: branchMetadataMatches.map((record) => record.id) });
  }

  const currentMissionId = await readCurrentMissionId();
  const currentRecord = currentMissionId ? findMissionRecord(records, currentMissionId) : null;
  if (currentMissionId) {
    signals.push({ source: ".space/current", value: currentMissionId, expectedMissionId: currentMissionId, missionId: currentRecord?.id ?? null });
    select(currentRecord, ".space/current");
  }

  const activeRecords = records.filter((record) => record.active);
  if (activeRecords.length === 1) {
    signals.push({ source: "single-active", value: activeRecords[0].id, missionId: activeRecords[0].id });
    select(activeRecords[0], "single-active");
  } else if (!selected && activeRecords.length > 1) {
    ambiguous = true;
  }

  conflicts.push(...signalConflicts(signals, explicitSelector, selected?.id ?? null, source));

  const orderedRecords = missionDisplayRecords(records);
  const displayNumberById = new Map(orderedRecords.map((record, index) => [record.id, index + 1]));
  const candidateRecords = candidateRecordsForResolution(records, selected, activeRecords, signals);

  return {
    selected: selected ? missionSummary(selected, source) : null,
    source,
    safety: resolveSafety(selected, conflicts, ambiguous),
    signals,
    conflicts,
    candidates: candidateRecords.map((record) => ({
      ...missionSummary(record),
      number: displayNumberById.get(record.id) ?? null
    })),
    currentMissionId,
    git: {
      branch: git.branch,
      sha: git.sha,
      isRepo: git.isRepo
    }
  };
}

function specTemplate() {
  return `# Mission Spec

## Goal

## User-visible behavior

## Non-goals

## Constraints

## Acceptance checks
`;
}

function questionsTemplate() {
  return `# Clarification Questions

## Open

## Answered
`;
}

function decisionsTemplate() {
  return `# Mission Decisions

## Confirmed

## Assumptions
`;
}

async function gitInfo() {
  const inside = await runCommand(["git", "rev-parse", "--is-inside-work-tree"]);
  if (inside.exitCode !== 0 || inside.stdout.trim() !== "true") {
    return {
      available: inside.exitCode !== 127,
      isRepo: false,
      root: null,
      branch: null,
      sha: null,
      dirty: null,
      dirtyFiles: 0
    };
  }

  const [root, branch, sha, status] = await Promise.all([
    runCommand(["git", "rev-parse", "--show-toplevel"]),
    runCommand(["git", "branch", "--show-current"]),
    runCommand(["git", "rev-parse", "HEAD"]),
    runCommand(["git", "status", "--short"])
  ]);

  const statusLines = status.stdout.split(/\r?\n/).filter((line) => line.trim());

  return {
    available: true,
    isRepo: true,
    root: root.exitCode === 0 ? root.stdout.trim() : null,
    branch: branch.exitCode === 0 ? branch.stdout.trim() || null : null,
    sha: sha.exitCode === 0 ? sha.stdout.trim() : null,
    dirty: status.exitCode === 0 ? statusLines.length > 0 : null,
    dirtyFiles: status.exitCode === 0 ? statusLines.length : 0
  };
}

function slugify(value) {
  const slug = String(value || "")
    .normalize("NFKD")
    .replace(/[\u0300-\u036f]/g, "")
    .toLowerCase()
    .replace(/[^a-z0-9]+/g, "-")
    .replace(/^-+|-+$/g, "")
    .replace(/-{2,}/g, "-")
    .slice(0, 60)
    .replace(/-+$/g, "");
  return slug || "mission";
}

async function printGitSuggestion(args) {
  const resolution = await resolveMission();
  if (resolution.safety !== "safe" || !resolution.selected) {
    fail(formatResolutionBlock(resolution, "git-suggest"));
  }

  const id = resolution.selected?.id ?? null;
  const mission = id ? await readJsonIfExists(path.join(missionDir(id), "mission.json")) : null;
  const missionPart = id ? id.toLowerCase() : "no-mission";
  const branchTypes = new Set(["feat", "fix", "docs", "refactor", "test", "build", "ci", "chore", "perf", "style", "issue", "release"]);
  const requestedType = args[0]?.toLowerCase();
  const type = branchTypes.has(requestedType) ? requestedType : "feat";
  const slugParts = branchTypes.has(requestedType) ? args.slice(1) : args;
  const slug = slugify(slugParts.join(" ") || mission?.title || "mission");
  const version = slugParts.join(" ").trim()
    .toLowerCase()
    .replace(/^v?/, "v")
    .replace(/[^a-z0-9._-]+/g, "-")
    .replace(/^-+|-+$/g, "") || "v0.1.0";
  const branchSlug = `${missionPart}-${slug}`;
  const branch = type === "release" ? `release/${version}` : `${type}/${branchSlug}`;
  const commitType = ["issue", "release"].includes(type) ? "chore" : type;

  console.log("Spacecraft git strategy: release branching");
  if (resolution.selected) {
    console.log(`Mission: ${resolution.selected.title} (${resolution.selected.id})`);
    console.log(`Selected by: ${resolution.source}`);
  }
  console.log("Base: latest main");
  console.log("Main: protected; no direct writes");
  console.log(`Branch: ${branch}`);
  console.log(type === "release"
    ? "Rule: release branches are for version/changelog/spec preparation only."
    : "Rule: one branch per feature/issue/tightly scoped change.");
  console.log("Final branch history: target 1-3 commits; maximum 5 unless justified.");
  console.log("Before merge: rebase on latest main, verify, bump version, update changelog/spec.");
  console.log("Merge: git merge --no-ff <branch>");
  console.log("After merge: create annotated version tag.");
  console.log("Use a separate worktree for large, risky, or multi-session branches.");
  console.log("");
  console.log("Conventional Commits format:");
  console.log("<type>[optional scope]: <description>");
  console.log("");
  console.log("Common types: feat, fix, docs, refactor, test, build, ci, chore, perf, style, revert");
  console.log("Examples:");
  if (type === "release") {
    console.log(`chore(release): prepare ${version}`);
  } else {
    console.log(`${commitType}(${slug}): add focused mission change`);
  }
  console.log(`docs(release): update changelog for v0.2.0`);
  console.log(`chore(release): bump version to v0.2.0`);
}

async function createMission(title) {
  if (!title || !title.trim()) {
    fail(`Missing mission title.\n\n${usage()}`);
  }

  await initSpacecraft({ silent: true });

  const id = missionId();
  const dir = missionDir(id);
  await fs.mkdir(dir, { recursive: false });
  await fs.mkdir(path.join(dir, "outputs"));

  const now = isoNow();
  const git = await gitInfo();
  const mission = {
    id,
    title: title.trim(),
    state: "draft",
    createdAt: now,
    updatedAt: now,
    baseSha: git.sha,
    headSha: null,
    git: {
      isRepo: git.isRepo,
      root: git.root,
      branch: git.branch,
      baseSha: git.sha,
      dirtyAtStart: git.dirty,
      dirtyFilesAtStart: git.dirtyFiles
    },
    artifacts: {
      spec: "spec.md",
      plan: "plan.json",
      evidence: "evidence.jsonl",
      review: "review.md",
      reviewJson: "review.json",
      questions: "questions.md",
      decisions: "decisions.md",
      design: "design/"
    },
    clarification: {
      status: "open",
      blockingQuestions: 0,
      lastQuestion: null
    }
  };

  await writeJson(path.join(dir, "mission.json"), mission);
  await fs.mkdir(path.join(dir, "design"));
  await fs.writeFile(path.join(dir, "spec.md"), specTemplate());
  await writeJson(path.join(dir, "plan.json"), {
    missionId: id,
    tasks: []
  });
  await fs.writeFile(path.join(dir, "evidence.jsonl"), "");
  await fs.writeFile(path.join(dir, "review.md"), "# Mission Review\n");
  await fs.writeFile(path.join(dir, "questions.md"), questionsTemplate());
  await fs.writeFile(path.join(dir, "decisions.md"), decisionsTemplate());
  await writeJson(path.join(dir, "review.json"), {
    status: "not-reviewed",
    findings: []
  });
  await fs.writeFile(CURRENT_FILE, `${id}\n`);
  const sessionFile = await writeSessionMissionId(id);

  console.log(`Created Spacecraft mission ${id}`);
  if (sessionFile) {
    console.log(`Session: ${displayPath(sessionFile)}`);
  }
  if (git.isRepo) {
    console.log(`Git: ${git.branch ?? "(detached)"} ${git.sha ? git.sha.slice(0, 12) : "(no commit)"}${git.dirty ? ` dirty:${git.dirtyFiles}` : ""}`);
  } else {
    console.log("Git: not a git worktree. Use only for discovery/design/read-only work, or explicitly accept no-git implementation risk.");
  }
  console.log("Next: /sc-plan");
}

async function printCurrent() {
  const id = await readCurrentMissionId();
  if (!id) {
    console.log("No current Spacecraft mission. Start one with /sc-start <title>.");
    return;
  }
  console.log(id);
}

function printMissionCandidate(candidate, index = null) {
  const displayIndex = Number.isInteger(candidate.number) ? candidate.number : index;
  const prefix = displayIndex === null ? "-" : `${displayIndex}.`;
  const branchHint = candidate.branches?.length ? ` branch:${candidate.branches.join(",")}` : "";
  const signal = candidate.signal ? ` signal:${candidate.signal}` : "";
  console.log(`${prefix} ${candidate.title} (${candidate.id}) - state:${candidate.state}${signal}${branchHint}`);
}

function nextCommandForMission(record) {
  const status = record.mission?.clarification?.status ?? "open";
  const blockingQuestions = record.mission?.clarification?.blockingQuestions ?? 0;
  if (status === "open" || blockingQuestions > 0) {
    return "/sc-clarify";
  }

  switch (record.mission?.state) {
    case "draft":
    case "specified":
      return "/sc-plan";
    case "planned":
    case "implementing":
      return "/sc-work";
    case "verifying":
      return "/sc-verify";
    case "reviewing":
      return "/sc-review";
    case "ready":
      return "/sc-ship";
    case "shipped":
      return "(shipped)";
    case "blocked":
      return "/sc-status";
    default:
      return "/sc-status";
  }
}

function branchHintsForRecord(record, resolution) {
  const hints = new Set(record.branches);
  if (resolution.git.branch && normalizeMissionId(resolution.git.branch) === record.id) {
    hints.add(resolution.git.branch);
  }
  return Array.from(hints);
}

function printMissionRecord(record, index, resolution) {
  const selected = resolution.selected?.id === record.id;
  const marker = selected ? "*" : " ";
  const signal = selected && resolution.source ? ` signal:${resolution.source}` : "";
  const branchHints = branchHintsForRecord(record, resolution);
  const branchHint = branchHints.length > 0 ? ` branch:${branchHints.join(",")}` : "";
  console.log(`${index}. ${marker} ${displayMissionTitle(record.mission?.title)} (${record.id}) - state:${record.mission?.state ?? "unknown"}${signal}${branchHint} next:${nextCommandForMission(record)}`);
}

function resolutionCandidateLines(resolution) {
  return resolution.candidates.map((candidate, index) => {
    const displayIndex = Number.isInteger(candidate.number) ? candidate.number : index + 1;
    const branchHint = candidate.branches?.length ? ` branch:${candidate.branches.join(",")}` : "";
    const signal = candidate.signal ? ` signal:${candidate.signal}` : "";
    return `${displayIndex}. ${candidate.title} (${candidate.id}) - state:${candidate.state}${signal}${branchHint}`;
  });
}

function resolutionConflictLines(resolution) {
  return resolution.conflicts.map((conflict) => {
    if (conflict.type === "branch-current") {
      return `branch mission ${conflict.branchMissionId} differs from .space/current ${conflict.currentMissionId}`;
    }
    if (conflict.type === "missing-current") {
      return `.space/current points to missing mission ${conflict.currentMissionId}`;
    }
    if (conflict.type === "missing-signal-mission") {
      return `${conflict.source} points to missing mission ${conflict.missionId}`;
    }
    if (conflict.type === "ambiguous-signal") {
      return `${conflict.source} matches multiple missions for ${conflict.value}: ${conflict.missionIds.join(", ")}`;
    }
    if (conflict.type === "signal-mismatch") {
      const signals = conflict.signals
        .map((signal) => `${signal.source} -> ${signal.missionId}`)
        .join("; ");
      return `mission signals disagree: ${signals}`;
    }
    return conflict.type;
  });
}

function formatResolutionBlock(resolution, commandName) {
  const lines = [
    `Spacecraft cannot safely choose a mission for ${commandName}.`,
    "No product files or mission artifacts changed.",
    `Safety: ${resolution.safety}`
  ];
  if (resolution.git.isRepo) {
    lines.push(`Branch: ${resolution.git.branch ?? "(detached)"}`);
  }
  if (resolution.currentMissionId) {
    lines.push(`Current: ${resolution.currentMissionId}`);
  }
  const conflicts = resolutionConflictLines(resolution);
  if (conflicts.length > 0) {
    lines.push("Conflicts:");
    lines.push(...conflicts.map((line) => `- ${line}`));
  }
  const candidates = resolutionCandidateLines(resolution);
  if (candidates.length > 0) {
    lines.push("Candidates:");
    lines.push(...candidates);
  }
  lines.push("Select a mission with: node scripts/spacecraft.mjs use <number|title>");
  lines.push("Advanced fallback: set SPACECRAFT_MISSION=<mission-id> for one command.");
  return lines.join("\n");
}

async function requireResolvedMission(commandName) {
  const resolution = await resolveMission();
  if (resolution.safety !== "safe" || !resolution.selected) {
    fail(formatResolutionBlock(resolution, commandName));
  }
  return resolution;
}

async function printMissions() {
  const records = await listMissionRecords();
  const resolution = await resolveMission();
  const orderedRecords = missionDisplayRecords(records);

  if (orderedRecords.length === 0) {
    console.log("No Spacecraft missions. Start one with /sc-start <title>.");
    return;
  }

  if (resolution.selected) {
    console.log(`Selected: ${resolution.selected.title} (${resolution.selected.id})`);
    console.log(`Selected by: ${resolution.source}`);
  } else {
    console.log("Selected: unresolved");
  }
  console.log(`Safety: ${resolution.safety}`);
  if (resolution.currentMissionId) {
    console.log(`Current: ${resolution.currentMissionId}`);
  }
  if (resolution.git.isRepo) {
    console.log(`Branch: ${resolution.git.branch ?? "(detached)"}`);
  }
  if (resolution.conflicts.length > 0) {
    console.log("Conflicts:");
    for (const line of resolutionConflictLines(resolution)) {
      console.log(`- ${line}`);
    }
  }

  console.log("Missions:");
  orderedRecords.forEach((record, index) => printMissionRecord(record, index + 1, resolution));
  console.log("Use: node scripts/spacecraft.mjs use <number|id|title>");
}

async function useMission(args) {
  const selector = args.join(" ").trim();
  if (!selector) {
    await printMissions();
    fail("Missing mission selector.");
  }

  const records = await listMissionRecords();
  const resolution = await resolveMission();
  const orderedRecords = missionDisplayRecords(records);
  const record = findMissionBySelector(records, selector, orderedRecords);
  if (!record) {
    console.log(`No unique mission matches "${selector}".`);
    console.log("Candidates:");
    orderedRecords.forEach((candidate, index) => printMissionRecord(candidate, index + 1, resolution));
    fail("Choose a mission by number, mission id, exact title, or unique title substring.");
  }

  await fs.writeFile(CURRENT_FILE, `${record.id}\n`);
  const sessionFile = await writeSessionMissionId(record.id);
  console.log(`Selected mission: ${displayMissionTitle(record.mission?.title)} (${record.id})`);
  if (sessionFile) {
    console.log(`Session: ${displayPath(sessionFile)}`);
  }
  console.log(`Next: ${nextCommandForMission(record)}`);
}

async function bindBranch(args) {
  const selector = args.join(" ").trim() || null;
  const resolution = selector ? await resolveMission({ selector }) : await requireResolvedMission("bind-branch");
  if (!resolution.selected || resolution.safety !== "safe") {
    fail(formatResolutionBlock(resolution, "bind-branch"));
  }

  const git = await gitInfo();
  if (!git.isRepo) {
    fail("Cannot bind branch outside a git worktree.");
  }
  if (!git.branch) {
    fail("Cannot bind branch while HEAD is detached.");
  }

  const missionPath = artifactPath(resolution.selected.id, "mission.json");
  const mission = await readJson(missionPath);
  mission.git = {
    ...(mission.git ?? {}),
    workBranch: git.branch,
    workBranchBoundAt: isoNow()
  };
  mission.updatedAt = isoNow();
  await writeJson(missionPath, mission);

  console.log(`Bound branch: ${git.branch}`);
  console.log(`Mission: ${displayMissionTitle(mission.title)} (${mission.id})`);
}

async function printResolvedMission(args) {
  const json = args.includes("--json");
  const selector = args.filter((arg) => arg !== "--json").join(" ").trim() || null;
  const resolution = await resolveMission({ selector });

  if (json) {
    console.log(JSON.stringify(resolution, null, 2));
    return;
  }

  if (resolution.selected) {
    console.log(`Mission: ${resolution.selected.title} (${resolution.selected.id})`);
    console.log(`Source: ${resolution.source}`);
  } else {
    console.log("Mission: unresolved");
  }
  console.log(`Safety: ${resolution.safety}`);
  if (resolution.git.isRepo) {
    console.log(`Branch: ${resolution.git.branch ?? "(detached)"}`);
  }
  if (resolution.currentMissionId) {
    console.log(`Current: ${resolution.currentMissionId}`);
  }

  if (resolution.conflicts.length > 0) {
    console.log("Conflicts:");
    for (const line of resolutionConflictLines(resolution)) {
      console.log(`- ${line}`);
    }
  }

  if (resolution.candidates.length > 0) {
    console.log("Candidates:");
    resolution.candidates.forEach((candidate, index) => printMissionCandidate(candidate, index + 1));
  }
}

async function countEvidence(filePath) {
  if (!(await exists(filePath))) {
    return 0;
  }
  const content = await fs.readFile(filePath, "utf8");
  return content.split(/\r?\n/).filter((line) => line.trim()).length;
}

async function readTextIfExists(filePath) {
  if (!(await exists(filePath))) {
    return null;
  }
  return fs.readFile(filePath, "utf8");
}

async function readEvidenceEntries(filePath) {
  const content = await readTextIfExists(filePath);
  if (!content) {
    return [];
  }

  const entries = [];
  for (const [index, line] of content.split(/\r?\n/).entries()) {
    if (!line.trim()) {
      continue;
    }
    try {
      entries.push(JSON.parse(line));
    } catch (error) {
      entries.push({
        id: `invalid-line-${index + 1}`,
        label: "Invalid evidence entry",
        command: "",
        exitCode: 1,
        createdAt: null,
        parseError: error.message
      });
    }
  }
  return entries;
}

function taskDisplay(task) {
  const id = task?.id ? String(task.id) : "task";
  const title = task?.title ? ` ${task.title}` : "";
  return `${id}${title}`;
}

function nextOpenTask(tasks) {
  return tasks.find((task) => task?.status !== "completed") ?? null;
}

function workflowSnapshot({ resolution, mission, specExists, planExists, plan, evidenceCount, git }) {
  const tasks = Array.isArray(plan?.tasks) ? plan.tasks : [];
  const nextTask = nextOpenTask(tasks);
  const blockers = [];
  const hasBlockingClarification = mission?.clarification?.status === "open" || (mission?.clarification?.blockingQuestions ?? 0) > 0;
  const hasTaskPlan = planExists && tasks.length > 0;
  const artifactGateClear = specExists && hasTaskPlan;
  let next = nextCommandForMission({
    mission,
    active: missionActive(mission),
    branches: missionBranchNames(mission),
    id: mission.id
  });

  if (hasBlockingClarification) {
    blockers.push("blocking clarification remains open");
    next = "/sc-clarify";
  }
  if (!specExists) {
    blockers.push("spec.md is missing");
    if (!hasBlockingClarification) {
      next = "/sc-status";
    }
  }
  if (!planExists) {
    blockers.push("plan.json is missing");
    if (!hasBlockingClarification && specExists) {
      next = "/sc-plan";
    }
  } else if (!hasTaskPlan) {
    blockers.push("plan.json has no tasks");
    if (!hasBlockingClarification && specExists) {
      next = "/sc-plan";
    }
  }
  if (["planned", "implementing", "verifying"].includes(mission?.state) && git.isRepo && (!git.branch || git.branch === "main")) {
    blockers.push("implementation workflow requires a non-main work branch");
  }
  if (["planned", "implementing", "verifying"].includes(mission?.state) && git.isRepo && git.dirty) {
    blockers.push(`worktree is dirty (${git.dirtyFiles} files); inspect before automated workflow`);
  }
  if (!hasBlockingClarification && artifactGateClear && (mission?.state === "planned" || mission?.state === "implementing")) {
    next = nextTask ? `/sc-work ${nextTask.id ?? ""}`.trim() : "/sc-review";
  }
  if (!hasBlockingClarification && artifactGateClear && mission?.state === "verifying") {
    next = nextTask ? `/sc-verify ${nextTask.id ?? ""}`.trim() : "/sc-review";
  }
  if (!hasBlockingClarification && artifactGateClear && !nextTask) {
    next = mission?.state === "ready" ? "/sc-ship" : "/sc-review";
  }

  return {
    missionId: mission.id,
    title: mission.title,
    state: mission.state,
    safety: resolution.safety,
    source: resolution.source,
    next,
    nextTask: nextTask ? {
      id: nextTask.id ?? null,
      title: nextTask.title ?? null,
      status: nextTask.status ?? null
    } : null,
    tasks: {
      total: tasks.length,
      completed: tasks.filter((task) => task?.status === "completed").length
    },
    evidenceCount,
    blockers,
    checkpointPolicy: "After passing verification, checkpoint commit on a clean non-main work branch before the next task."
  };
}

async function printWorkflow(args) {
  const json = args.includes("--json");
  const resolution = await resolveMission();
  if (resolution.safety !== "safe" || !resolution.selected) {
    fail(formatResolutionBlock(resolution, "flow"));
  }

  const id = resolution.selected.id;
  const dir = missionDir(id);
  const mission = await readJson(artifactPath(id, "mission.json"));
  const specExists = await exists(path.join(dir, "spec.md"));
  const planPath = path.join(dir, "plan.json");
  const planExists = await exists(planPath);
  const plan = planExists ? await readJson(planPath) : null;
  const evidenceCount = await countEvidence(path.join(dir, "evidence.jsonl"));
  const git = await gitInfo();
  const snapshot = workflowSnapshot({ resolution, mission, specExists, planExists, plan, evidenceCount, git });

  if (json) {
    console.log(JSON.stringify(snapshot, null, 2));
    return;
  }

  console.log(`Workflow: ${snapshot.blockers.length > 0 ? "blocked" : "ready"}`);
  console.log(`Mission: ${snapshot.title} (${snapshot.missionId})`);
  console.log(`State: ${snapshot.state}`);
  console.log(`Tasks: ${snapshot.tasks.completed}/${snapshot.tasks.total} completed`);
  console.log(`Evidence: ${snapshot.evidenceCount}`);
  if (snapshot.nextTask) {
    console.log(`Next task: ${taskDisplay(snapshot.nextTask)}`);
  }
  console.log(`Next: ${snapshot.next}`);
  console.log("Loop: /sc-work Txx -> /sc-verify Txx -> checkpoint commit -> next task, until a gate blocks.");
  console.log(`Checkpoint: ${snapshot.checkpointPolicy}`);
  if (snapshot.blockers.length > 0) {
    console.log("Blockers:");
    for (const blocker of snapshot.blockers) {
      console.log(`- ${blocker}`);
    }
  }
}

async function printStatus() {
  const resolution = await resolveMission();
  if (!resolution.selected) {
    console.log("No selected Spacecraft mission. Start one with /sc-start <title>.");
    if (resolution.candidates.length > 0) {
      console.log("Candidates:");
      resolution.candidates.forEach((candidate, index) => printMissionCandidate(candidate, index + 1));
    }
    return;
  }

  const id = resolution.selected.id;
  const dir = missionDir(id);
  const missionPath = path.join(dir, "mission.json");
  if (!(await exists(missionPath))) {
    fail(`Selected mission ${id} is missing mission.json.`);
  }

  const mission = await readJson(missionPath);
  const plan = await readJsonIfExists(path.join(dir, "plan.json"));
  const review = await readJsonIfExists(path.join(dir, "review.json"));
  const taskCount = Array.isArray(plan?.tasks) ? plan.tasks.length : 0;
  const evidenceCount = await countEvidence(path.join(dir, "evidence.jsonl"));

  console.log(`Mission: ${mission.id}`);
  console.log(`Selected by: ${resolution.source}`);
  if (resolution.safety !== "safe") {
    console.log(`Mission safety: ${resolution.safety}`);
  }
  if (resolution.currentMissionId) {
    console.log(`Current: ${resolution.currentMissionId}`);
  }
  if (resolution.conflicts.length > 0) {
    console.log("Conflicts:");
    for (const line of resolutionConflictLines(resolution)) {
      console.log(`- ${line}`);
    }
  }
  if (resolution.safety !== "safe" && resolution.candidates.length > 0) {
    console.log("Candidates:");
    resolution.candidates.forEach((candidate, index) => printMissionCandidate(candidate, index + 1));
  }
  console.log(`Title: ${mission.title}`);
  console.log(`State: ${mission.state}`);
  if (mission.clarification?.status) {
    console.log(`Clarification: ${mission.clarification.status}`);
    console.log(`Blocking questions: ${mission.clarification.blockingQuestions ?? 0}`);
  }
  const git = await gitInfo();
  if (git.isRepo) {
    console.log(`Git: ${git.branch ?? "(detached)"} ${git.sha ? git.sha.slice(0, 12) : "(no commit)"}${git.dirty ? ` dirty:${git.dirtyFiles}` : " clean"}`);
    if (mission.baseSha && git.sha && mission.baseSha !== git.sha) {
      console.log(`Mission base: ${mission.baseSha.slice(0, 12)}`);
    }
  } else {
    console.log("Git: not a git worktree");
  }
  console.log("Artifacts:");
  console.log(`  spec: ${displayPath(path.join(dir, "spec.md"))}`);
  console.log(`  plan: ${displayPath(path.join(dir, "plan.json"))}`);
  console.log(`  evidence: ${displayPath(path.join(dir, "evidence.jsonl"))}`);
  console.log(`  review: ${displayPath(path.join(dir, "review.md"))}`);
  console.log(`  reviewJson: ${displayPath(path.join(dir, "review.json"))}`);
  if (await exists(path.join(dir, "questions.md"))) {
    console.log(`  questions: ${displayPath(path.join(dir, "questions.md"))}`);
  }
  if (await exists(path.join(dir, "decisions.md"))) {
    console.log(`  decisions: ${displayPath(path.join(dir, "decisions.md"))}`);
  }
  if (await exists(path.join(dir, "design"))) {
    console.log(`  design: ${displayPath(path.join(dir, "design"))}`);
  }
  console.log(`Tasks: ${taskCount}`);
  console.log(`Evidence: ${evidenceCount}`);
  console.log(`Review: ${review?.status ?? "not-reviewed"}`);
}

async function printGitInfo() {
  const git = await gitInfo();
  if (!git.available) {
    console.log("Git: command unavailable");
    return;
  }
  if (!git.isRepo) {
    console.log("Git: not a git worktree");
    return;
  }

  console.log("Git: worktree detected");
  console.log(`Root: ${git.root ?? "(unknown)"}`);
  console.log(`Branch: ${git.branch ?? "(detached)"}`);
  console.log(`HEAD: ${git.sha ?? "(no commit)"}`);
  console.log(`Status: ${git.dirty ? `dirty (${git.dirtyFiles} files)` : "clean"}`);
}

async function setState(state) {
  if (!state) {
    fail(`Missing state.\nAllowed states: ${Array.from(STATES).join(", ")}`);
  }
  if (!STATES.has(state)) {
    fail(`Invalid state "${state}". Allowed states: ${Array.from(STATES).join(", ")}`);
  }

  const resolution = await requireResolvedMission("set-state");
  const id = resolution.selected.id;
  const missionPath = artifactPath(id, "mission.json");
  const mission = await readJson(missionPath);
  mission.state = state;
  mission.updatedAt = isoNow();
  await writeJson(missionPath, mission);
  console.log(`Spacecraft mission ${id} state: ${state}`);
}

async function setClarificationStatus(status) {
  if (!status) {
    fail(`Missing clarification status.\nAllowed statuses: ${Array.from(CLARIFICATION_STATUSES).join(", ")}`);
  }
  if (!CLARIFICATION_STATUSES.has(status)) {
    fail(`Invalid clarification status "${status}". Allowed statuses: ${Array.from(CLARIFICATION_STATUSES).join(", ")}`);
  }

  const resolution = await requireResolvedMission("clarify-status");
  const id = resolution.selected.id;
  const missionPath = artifactPath(id, "mission.json");
  const mission = await readJson(missionPath);
  mission.clarification = {
    status,
    blockingQuestions: status === "clear" ? 0 : mission.clarification?.blockingQuestions ?? 0,
    lastQuestion: status === "clear" ? null : mission.clarification?.lastQuestion ?? null
  };
  mission.updatedAt = isoNow();
  await writeJson(missionPath, mission);
  console.log(`Spacecraft mission ${id} clarification: ${status}`);
}

function commandToString(parts) {
  return parts
    .map((part) => {
      if (/^[A-Za-z0-9_./:=@%+-]+$/.test(part)) {
        return part;
      }
      return `'${part.replaceAll("'", "'\\''")}'`;
    })
    .join(" ");
}

function evidenceFilePath(entryPath) {
  if (!entryPath || typeof entryPath !== "string") {
    return null;
  }
  return path.isAbsolute(entryPath) ? entryPath : path.join(ROOT, entryPath);
}

const DEFAULT_RELEASE_GATE_STATUSES = [
  "bumped",
  "checked",
  "complete",
  "completed",
  "deferred",
  "done",
  "passed",
  "present",
  "updated"
];

const RELEASE_GATE_DEFINITIONS = [
  ["version", "Record version bump or explicit deferral with rationale in review.json releaseReadiness.version.", DEFAULT_RELEASE_GATE_STATUSES],
  ["changelog", "Record changelog update or explicit deferral with rationale in review.json releaseReadiness.changelog.", DEFAULT_RELEASE_GATE_STATUSES],
  ["specNote", "Record short spec/release note update or explicit deferral with rationale in review.json releaseReadiness.specNote.", DEFAULT_RELEASE_GATE_STATUSES],
  ["tagPlan", "Record the post-merge version tag plan in review.json releaseReadiness.tagPlan.", [...DEFAULT_RELEASE_GATE_STATUSES, "planned"]],
  ["postRebaseVerification", "Record verification after latest rebase in review.json releaseReadiness.postRebaseVerification.", DEFAULT_RELEASE_GATE_STATUSES]
];

function releaseGateSatisfied(gate, allowedStatuses = DEFAULT_RELEASE_GATE_STATUSES) {
  if (!gate || typeof gate !== "object" || Array.isArray(gate)) {
    return false;
  }

  const status = String(gate.status ?? "").trim().toLowerCase();
  const satisfiedStatuses = new Set(allowedStatuses);
  if (!satisfiedStatuses.has(status)) {
    return false;
  }
  if (status === "deferred" && !String(gate.rationale ?? "").trim()) {
    return false;
  }
  return true;
}

function releaseReadinessErrors(releaseReadiness) {
  const errors = [];
  for (const [key, message, allowedStatuses] of RELEASE_GATE_DEFINITIONS) {
    if (!releaseGateSatisfied(releaseReadiness?.[key], allowedStatuses)) {
      errors.push(message);
    }
  }
  return errors;
}

function conventionalCommitSubject(subject) {
  return /^(feat|fix|docs|refactor|test|build|ci|chore|perf|style|revert)(\([a-z0-9._/-]+\))?!?: .+/.test(subject);
}

async function reserveEvidencePaths(dir) {
  const base = evidenceId();
  let candidate = base;
  let index = 2;
  while (true) {
    const stdoutPath = path.join(dir, "outputs", `${candidate}.stdout.txt`);
    const stderrPath = path.join(dir, "outputs", `${candidate}.stderr.txt`);
    try {
      const handle = await fs.open(stdoutPath, "wx");
      await handle.close();
      return { evidence: candidate, stdoutPath, stderrPath };
    } catch (error) {
      if (error?.code !== "EEXIST") {
        throw error;
      }
      candidate = evidenceId(new Date(Date.now() + index));
      index += 1;
    }
  }
}

function runCommand(commandParts) {
  return new Promise((resolve) => {
    const child = spawn(commandParts[0], commandParts.slice(1), {
      cwd: ROOT,
      env: process.env,
      stdio: ["ignore", "pipe", "pipe"]
    });

    let stdout = "";
    let stderr = "";
    let spawnError = null;

    child.stdout.setEncoding("utf8");
    child.stderr.setEncoding("utf8");

    child.stdout.on("data", (chunk) => {
      stdout += chunk;
    });
    child.stderr.on("data", (chunk) => {
      stderr += chunk;
    });
    child.on("error", (error) => {
      spawnError = error;
    });
    child.on("close", (code, signal) => {
      if (spawnError) {
        stderr += `${spawnError.message}\n`;
        resolve({ exitCode: 127, stdout, stderr });
        return;
      }
      if (signal) {
        stderr += `Command terminated by signal ${signal}\n`;
      }
      resolve({ exitCode: code ?? 1, stdout, stderr });
    });
  });
}

async function recordEvidence(args) {
  const separator = args.indexOf("--");
  if (separator === -1) {
    fail(`Missing -- before evidence command.\n\n${usage()}`);
  }

  const label = args.slice(0, separator).join(" ").trim();
  const commandParts = args.slice(separator + 1);

  if (!label) {
    fail(`Missing evidence label.\n\n${usage()}`);
  }
  if (commandParts.length === 0) {
    fail(`Missing command after --.\n\n${usage()}`);
  }

  const resolution = await requireResolvedMission("evidence");
  const id = resolution.selected.id;
  const dir = missionDir(id);
  const outputsDir = path.join(dir, "outputs");
  await fs.mkdir(outputsDir, { recursive: true });

  const { evidence, stdoutPath, stderrPath } = await reserveEvidencePaths(dir);
  const result = await runCommand(commandParts);

  await fs.writeFile(stdoutPath, result.stdout);
  await fs.writeFile(stderrPath, result.stderr);

  const entry = {
    id: evidence,
    label,
    command: commandToString(commandParts),
    exitCode: result.exitCode,
    stdout: displayPath(stdoutPath),
    stderr: displayPath(stderrPath),
    createdAt: isoNow()
  };

  await fs.appendFile(path.join(dir, "evidence.jsonl"), `${JSON.stringify(entry)}\n`);

  console.log(`Evidence: ${evidence}`);
  console.log(`Exit code: ${result.exitCode}`);
  process.exitCode = result.exitCode;
}

async function validateMission() {
  const errors = [];
  const resolution = await requireResolvedMission("validate");
  const id = resolution.selected.id;
  const dir = missionDir(id);

  async function requireFile(relativePath) {
    const filePath = path.join(dir, relativePath);
    if (!(await exists(filePath))) {
      errors.push(`Missing ${relativePath}`);
      return null;
    }
    return filePath;
  }

  async function parseRequiredJson(relativePath) {
    const filePath = await requireFile(relativePath);
    if (!filePath) {
      return null;
    }
    try {
      return await readJson(filePath);
    } catch (error) {
      errors.push(`Invalid JSON in ${relativePath}: ${error.message}`);
      return null;
    }
  }

  const mission = await parseRequiredJson("mission.json");
  if (mission?.clarification) {
    if (!CLARIFICATION_STATUSES.has(mission.clarification.status)) {
      errors.push(`mission.json clarification.status must be one of: ${Array.from(CLARIFICATION_STATUSES).join(", ")}`);
    }
    if (
      mission.clarification.blockingQuestions !== undefined &&
      !Number.isInteger(mission.clarification.blockingQuestions)
    ) {
      errors.push("mission.json clarification.blockingQuestions must be an integer when present");
    }
  }
  await requireFile("spec.md");
  const plan = await parseRequiredJson("plan.json");
  if (plan && !Array.isArray(plan.tasks)) {
    errors.push("plan.json must contain a tasks array");
  }

  const evidencePath = await requireFile("evidence.jsonl");
  if (evidencePath) {
    const lines = (await fs.readFile(evidencePath, "utf8")).split(/\r?\n/);
    const evidenceIds = new Map();
    const evidenceOutputs = new Map();
    lines.forEach((line, index) => {
      if (!line.trim()) {
        return;
      }
      try {
        const entry = JSON.parse(line);
        if (!entry.id || typeof entry.id !== "string") {
          errors.push(`evidence.jsonl line ${index + 1} must have string id`);
        } else if (evidenceIds.has(entry.id)) {
          errors.push(`Duplicate evidence id ${entry.id} on lines ${evidenceIds.get(entry.id)} and ${index + 1}`);
        } else {
          evidenceIds.set(entry.id, index + 1);
        }

        for (const field of ["stdout", "stderr"]) {
          const filePath = evidenceFilePath(entry[field]);
          if (!filePath) {
            errors.push(`evidence.jsonl line ${index + 1} must have ${field} path`);
            continue;
          }
          const key = path.resolve(filePath);
          if (evidenceOutputs.has(key)) {
            errors.push(`Duplicate evidence ${field} path ${entry[field]} on lines ${evidenceOutputs.get(key)} and ${index + 1}`);
          } else {
            evidenceOutputs.set(key, index + 1);
          }
          if (!existsSync(filePath)) {
            errors.push(`Missing evidence ${field} file for line ${index + 1}: ${entry[field]}`);
          }
        }
      } catch (error) {
        errors.push(`Invalid JSON in evidence.jsonl line ${index + 1}: ${error.message}`);
      }
    });
  }

  const reviewPath = path.join(dir, "review.json");
  if (await exists(reviewPath)) {
    try {
      await readJson(reviewPath);
    } catch (error) {
      errors.push(`Invalid JSON in review.json: ${error.message}`);
    }
  }

  if (errors.length > 0) {
    console.error(`Spacecraft mission ${id} is invalid:`);
    for (const error of errors) {
      console.error(`- ${error}`);
    }
    process.exitCode = 1;
    return;
  }

  console.log(`Spacecraft mission ${id} is valid.`);
}

function compactEvidenceEntry(entry) {
  return {
    id: entry.id ?? null,
    label: entry.label ?? null,
    command: entry.command ?? null,
    exitCode: entry.exitCode ?? null,
    createdAt: entry.createdAt ?? null
  };
}

function compactMissionRecord(mission, { archivedAt }) {
  return {
    id: mission.id,
    title: mission.title,
    state: mission.state,
    createdAt: mission.createdAt ?? null,
    updatedAt: mission.updatedAt ?? null,
    archivedAt,
    baseSha: mission.baseSha ?? mission.git?.baseSha ?? null,
    headSha: mission.headSha ?? null,
    git: {
      root: mission.git?.root ?? null,
      branch: mission.git?.workBranch ?? mission.git?.branch ?? null,
      baseSha: mission.git?.baseSha ?? mission.baseSha ?? null
    }
  };
}

function compactPlanRecord(plan) {
  const tasks = Array.isArray(plan?.tasks) ? plan.tasks : [];
  return {
    missionId: plan?.missionId ?? null,
    tasks: tasks.map((task) => ({
      id: task.id ?? null,
      title: task.title ?? null,
      status: task.status ?? null,
      evidence: Array.isArray(task.evidence) ? task.evidence : []
    }))
  };
}

function blockingReviewFindings(review) {
  const findings = Array.isArray(review?.findings) ? review.findings : [];
  return findings.filter((finding) => finding?.blocksShip || finding?.severity === "critical");
}

function archiveReadinessErrors({ plan, review, evidenceEntries }) {
  const errors = [];
  if (!plan) {
    errors.push("missing plan.json");
  }
  if (!review) {
    errors.push("missing review.json");
  } else if (review.status !== "ready") {
    errors.push(`review status is ${review.status ?? "missing"}`);
  }

  const tasks = Array.isArray(plan?.tasks) ? plan.tasks : [];
  if (tasks.length === 0) {
    errors.push("plan.json has no tasks");
  }
  const incompleteTasks = tasks.filter((task) => task?.status !== "completed");
  if (incompleteTasks.length > 0) {
    errors.push(`incomplete tasks: ${incompleteTasks.map((task) => task.id ?? task.title ?? "unnamed").join(", ")}`);
  }
  if (evidenceEntries.length === 0) {
    errors.push("evidence.jsonl has no evidence");
  }

  const blockingFindings = blockingReviewFindings(review);
  if (blockingFindings.length > 0) {
    errors.push(`blocking review findings: ${blockingFindings.map((finding) => finding.id ?? finding.summary ?? "unnamed").join(", ")}`);
  }
  errors.push(...releaseReadinessErrors(review?.releaseReadiness));
  return errors;
}

async function copyArchiveText(source, destination) {
  const content = await readTextIfExists(source);
  if (content === null) {
    return false;
  }
  await fs.writeFile(destination, content);
  return true;
}

async function clearArchivedMissionSelection(id) {
  const currentMissionId = await readCurrentMissionId();
  if (currentMissionId === id) {
    await fs.writeFile(CURRENT_FILE, "");
  }

  const sessionsDir = path.join(SPACE_DIR, "sessions");
  if (!(await exists(sessionsDir))) {
    return;
  }
  const entries = await fs.readdir(sessionsDir, { withFileTypes: true });
  for (const entry of entries) {
    if (!entry.isFile()) {
      continue;
    }
    const sessionPath = path.join(sessionsDir, entry.name);
    const sessionMissionId = normalizeMissionId(await fs.readFile(sessionPath, "utf8"));
    if (sessionMissionId === id) {
      await fs.writeFile(sessionPath, "");
    }
  }
}

async function archiveMission(args) {
  const selector = args.join(" ").trim() || null;
  const resolution = selector ? await resolveMission({ selector }) : await requireResolvedMission("archive");
  if (!resolution.selected || resolution.safety !== "safe") {
    fail(formatResolutionBlock(resolution, "archive"));
  }

  const id = resolution.selected.id;
  const sourceDir = missionDir(id);
  const mission = await readJson(artifactPath(id, "mission.json"));
  if (mission.state !== "shipped") {
    fail(`Archive blocked: mission ${id} state is ${mission.state}. Archive only shipped missions.`);
  }

  const plan = await readJsonIfExists(path.join(sourceDir, "plan.json"));
  const review = await readJsonIfExists(path.join(sourceDir, "review.json"));
  const evidenceEntries = await readEvidenceEntries(path.join(sourceDir, "evidence.jsonl"));
  const readinessErrors = archiveReadinessErrors({ plan, review, evidenceEntries });
  if (readinessErrors.length > 0) {
    fail(`Archive blocked for ${id}:\n- ${readinessErrors.join("\n- ")}`);
  }

  await fs.mkdir(ARCHIVE_DIR, { recursive: true });
  const archiveDir = path.join(ARCHIVE_DIR, id);
  if (await exists(archiveDir)) {
    fail(`Archive already exists: ${displayPath(archiveDir)}`);
  }
  await fs.mkdir(archiveDir, { recursive: false });

  const archivedAt = isoNow();
  const compactEvidence = evidenceEntries.map(compactEvidenceEntry);
  const tasks = Array.isArray(plan?.tasks) ? plan.tasks : [];
  const completedTasks = tasks.filter((task) => task?.status === "completed").length;

  const summary = [
    `# Archived Mission ${id}`,
    "",
    `Title: ${mission.title ?? "(untitled)"}`,
    `State: ${mission.state}`,
    `Created: ${mission.createdAt ?? "(unknown)"}`,
    `Archived: ${archivedAt}`,
    `Branch: ${mission.git?.workBranch ?? mission.git?.branch ?? "(unknown)"}`,
    `Tasks: ${completedTasks}/${tasks.length} completed`,
    `Evidence: ${compactEvidence.length}`,
    `Review: ${review?.status ?? "missing"}`,
    "",
    "## Evidence",
    ...compactEvidence.map((entry) => `- ${entry.id}: ${entry.label ?? "(unlabeled)"} [exit ${entry.exitCode ?? "?"}] ${entry.command ?? ""}`),
    "",
    "## Kept Artifacts",
    "- SUMMARY.md",
    "- mission.json",
    "- plan.json",
    "- evidence.jsonl",
    "- review.json / review.md when present",
    "- spec.md, decisions.md, and questions.md when present",
    ""
  ].join("\n");

  await fs.writeFile(path.join(archiveDir, "SUMMARY.md"), summary);
  await writeJson(path.join(archiveDir, "mission.json"), compactMissionRecord(mission, { archivedAt }));
  await writeJson(path.join(archiveDir, "plan.json"), compactPlanRecord(plan));
  await fs.writeFile(
    path.join(archiveDir, "evidence.jsonl"),
    compactEvidence.map((entry) => JSON.stringify(entry)).join("\n") + (compactEvidence.length ? "\n" : "")
  );
  if (review) {
    await writeJson(path.join(archiveDir, "review.json"), review);
  }
  await copyArchiveText(path.join(sourceDir, "review.md"), path.join(archiveDir, "review.md"));
  await copyArchiveText(path.join(sourceDir, "spec.md"), path.join(archiveDir, "spec.md"));
  await copyArchiveText(path.join(sourceDir, "decisions.md"), path.join(archiveDir, "decisions.md"));
  await copyArchiveText(path.join(sourceDir, "questions.md"), path.join(archiveDir, "questions.md"));

  await fs.rm(sourceDir, { recursive: true, force: false });
  await clearArchivedMissionSelection(id);

  console.log(`Archived mission ${id}`);
  console.log(`Archive: ${displayPath(archiveDir)}`);
}

async function releaseCloseoutCheck() {
  const errors = [];
  const warnings = [];
  const resolution = await requireResolvedMission("closeout-check");
  const id = resolution.selected.id;
  const dir = missionDir(id);

  const mission = await readJsonIfExists(path.join(dir, "mission.json"));
  const plan = await readJsonIfExists(path.join(dir, "plan.json"));
  const review = await readJsonIfExists(path.join(dir, "review.json"));
  const evidenceCount = await countEvidence(path.join(dir, "evidence.jsonl"));

  if (!mission) {
    errors.push("Missing mission.json");
  }
  if (!plan) {
    errors.push("Missing plan.json");
  }
  if (!review) {
    errors.push("Missing review.json");
  }

  if (mission?.clarification?.status === "open" || (mission?.clarification?.blockingQuestions ?? 0) > 0) {
    errors.push("Resolve blocking clarification questions.");
  }

  const tasks = Array.isArray(plan?.tasks) ? plan.tasks : [];
  const incompleteTasks = tasks.filter((task) => task.status !== "completed");
  if (incompleteTasks.length > 0) {
    errors.push(`Complete plan tasks: ${incompleteTasks.map((task) => task.id ?? task.title).join(", ")}.`);
  }

  if (evidenceCount === 0) {
    errors.push("Capture verification evidence in evidence.jsonl.");
  }

  if (review && review.status !== "ready") {
    errors.push(`Review status must be ready; current status is ${review.status ?? "missing"}.`);
  }

  const blockingFindings = blockingReviewFindings(review);
  if (blockingFindings.length > 0) {
    errors.push(`Fix blocking review findings: ${blockingFindings.map((finding) => finding.id ?? finding.summary ?? "unnamed").join(", ")}.`);
  }

  errors.push(...releaseReadinessErrors(review?.releaseReadiness));

  const git = await gitInfo();
  if (!git.isRepo) {
    errors.push("Run closeout inside a git worktree.");
  } else {
    if (!git.branch || git.branch === "main") {
      errors.push("Closeout must run from a non-main work branch.");
    }
    if (git.dirty) {
      errors.push(`Commit, stash, or remove dirty worktree changes (${git.dirtyFiles} files).`);
    }

    const mainAncestor = await runCommand(["git", "merge-base", "--is-ancestor", "main", "HEAD"]);
    if (mainAncestor.exitCode !== 0) {
      errors.push("Rebase the work branch on latest main, then rerun verification.");
    }

    const commitCount = await runCommand(["git", "rev-list", "--count", "main..HEAD"]);
    let count = null;
    if (commitCount.exitCode === 0) {
      count = Number.parseInt(commitCount.stdout.trim(), 10);
      if (Number.isFinite(count) && count > 5) {
        errors.push(`Squash/fixup branch history to 5 or fewer final commits; current count is ${count}.`);
      }
    } else {
      warnings.push("Could not count commits from main..HEAD.");
    }

    const commitSubjects = await runCommand(["git", "log", "--format=%s", "main..HEAD"]);
    if (commitSubjects.exitCode === 0) {
      const subjects = commitSubjects.stdout.split(/\r?\n/).filter((line) => line.trim());
      if (count === 0 || subjects.length === 0) {
        errors.push("Create at least one final Conventional Commit before closeout.");
      }
      const invalidSubjects = subjects.filter((subject) => !conventionalCommitSubject(subject));
      if (invalidSubjects.length > 0) {
        errors.push(`Fix non-Conventional Commit subjects: ${invalidSubjects.join("; ")}`);
      }
    } else {
      errors.push("Check final commit subjects before closeout.");
    }
  }

  if (errors.length > 0) {
    console.log(`Spacecraft closeout blocked for ${id}:`);
    for (const error of errors) {
      console.log(`- ${error}`);
    }
    if (warnings.length > 0) {
      console.log("Warnings:");
      for (const warning of warnings) {
        console.log(`- ${warning}`);
      }
    }
    process.exitCode = 1;
    return;
  }

  console.log(`Spacecraft closeout ready for ${id}.`);
  console.log("Next: rebase already satisfied, run final verification if needed, merge with git merge --no-ff <branch>, tag release, then delete merged branch unless kept.");
  if (warnings.length > 0) {
    console.log("Warnings:");
    for (const warning of warnings) {
      console.log(`- ${warning}`);
    }
  }
}

async function main() {
  const [command, ...args] = process.argv.slice(2);

  switch (command) {
    case "init":
      await initSpacecraft();
      break;
    case "new":
      await createMission(args.join(" "));
      break;
    case "current":
      await printCurrent();
      break;
    case "resolve":
      await printResolvedMission(args);
      break;
    case "missions":
      await printMissions();
      break;
    case "use":
      await useMission(args);
      break;
    case "bind-branch":
      await bindBranch(args);
      break;
    case "status":
      await printStatus();
      break;
    case "flow":
      await printWorkflow(args);
      break;
    case "git-info":
      await printGitInfo();
      break;
    case "git-suggest":
      await printGitSuggestion(args);
      break;
    case "set-state":
      await setState(args[0]);
      break;
    case "clarify-status":
      await setClarificationStatus(args[0]);
      break;
    case "evidence":
      await recordEvidence(args);
      break;
    case "validate":
      await validateMission();
      break;
    case "closeout-check":
      await releaseCloseoutCheck();
      break;
    case "archive":
      await archiveMission(args);
      break;
    case undefined:
    case "-h":
    case "--help":
    case "help":
      console.log(usage());
      break;
    default:
      fail(`Unknown command "${command}".\n\n${usage()}`);
  }
}

main().catch((error) => {
  console.error(`Spacecraft error: ${error.message}`);
  process.exit(error.exitCode ?? 1);
});
