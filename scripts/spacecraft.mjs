#!/usr/bin/env node
import { spawn } from "node:child_process";
import { constants as fsConstants, existsSync } from "node:fs";
import fs from "node:fs/promises";
import path from "node:path";
import process from "node:process";

const ROOT = process.cwd();
const SPACE_DIR = path.join(ROOT, ".space");
const MISSIONS_DIR = path.join(SPACE_DIR, "missions");
const CURRENT_FILE = path.join(SPACE_DIR, "current");

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
  node scripts/spacecraft.mjs status
  node scripts/spacecraft.mjs git-info
  node scripts/spacecraft.mjs git-suggest [type] [slug]
  node scripts/spacecraft.mjs set-state <state>
  node scripts/spacecraft.mjs clarify-status <open|clear|deferred>
  node scripts/spacecraft.mjs evidence <label> -- <command...>
  node scripts/spacecraft.mjs validate
  node scripts/spacecraft.mjs closeout-check
`;
}

function fail(message, code = 1) {
  const error = new Error(message);
  error.exitCode = code;
  throw error;
}

function pad(value) {
  return String(value).padStart(2, "0");
}

function timestampForId(date = new Date()) {
  return [
    date.getFullYear(),
    pad(date.getMonth() + 1),
    pad(date.getDate())
  ].join("") + "-" + [
    pad(date.getHours()),
    pad(date.getMinutes()),
    pad(date.getSeconds())
  ].join("");
}

function missionId(date = new Date()) {
  return `M-${timestampForId(date)}`;
}

function evidenceId(date = new Date()) {
  return `E-${timestampForId(date)}`;
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
  return value;
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
  const id = await readCurrentMissionId();
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

  console.log(`Created Spacecraft mission ${id}`);
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

async function countEvidence(filePath) {
  if (!(await exists(filePath))) {
    return 0;
  }
  const content = await fs.readFile(filePath, "utf8");
  return content.split(/\r?\n/).filter((line) => line.trim()).length;
}

async function printStatus() {
  const id = await readCurrentMissionId();
  if (!id) {
    console.log("No current Spacecraft mission. Start one with /sc-start <title>.");
    return;
  }

  const dir = missionDir(id);
  const missionPath = path.join(dir, "mission.json");
  if (!(await exists(missionPath))) {
    fail(`Current mission ${id} is missing mission.json.`);
  }

  const mission = await readJson(missionPath);
  const plan = await readJsonIfExists(path.join(dir, "plan.json"));
  const review = await readJsonIfExists(path.join(dir, "review.json"));
  const taskCount = Array.isArray(plan?.tasks) ? plan.tasks.length : 0;
  const evidenceCount = await countEvidence(path.join(dir, "evidence.jsonl"));

  console.log(`Mission: ${mission.id}`);
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

  const id = await readCurrentMissionId({ required: true });
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

  const id = await readCurrentMissionId({ required: true });
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

function releaseGateSatisfied(gate) {
  if (!gate || typeof gate !== "object" || Array.isArray(gate)) {
    return false;
  }

  const status = String(gate.status ?? "").trim().toLowerCase();
  const satisfiedStatuses = new Set([
    "bumped",
    "checked",
    "complete",
    "completed",
    "deferred",
    "done",
    "passed",
    "planned",
    "present",
    "updated"
  ]);
  if (!satisfiedStatuses.has(status)) {
    return false;
  }
  if (status === "deferred" && !String(gate.rationale ?? "").trim()) {
    return false;
  }
  return true;
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
      candidate = `${base}-${index}`;
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

  const id = await readCurrentMissionId({ required: true });
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
  const id = await readCurrentMissionId({ required: true });
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

async function releaseCloseoutCheck() {
  const errors = [];
  const warnings = [];
  const id = await readCurrentMissionId({ required: true });
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

  const findings = Array.isArray(review?.findings) ? review.findings : [];
  const blockingFindings = findings.filter((finding) => finding?.blocksShip || finding?.severity === "critical");
  if (blockingFindings.length > 0) {
    errors.push(`Fix blocking review findings: ${blockingFindings.map((finding) => finding.id ?? finding.summary ?? "unnamed").join(", ")}.`);
  }

  const releaseReadiness = review?.releaseReadiness ?? null;
  const releaseGates = [
    ["version", "Record version bump or explicit deferral with rationale in review.json releaseReadiness.version."],
    ["changelog", "Record changelog update or explicit deferral with rationale in review.json releaseReadiness.changelog."],
    ["specNote", "Record short spec/release note update or explicit deferral with rationale in review.json releaseReadiness.specNote."],
    ["tagPlan", "Record the post-merge version tag plan in review.json releaseReadiness.tagPlan."],
    ["postRebaseVerification", "Record verification after latest rebase in review.json releaseReadiness.postRebaseVerification."]
  ];
  for (const [key, message] of releaseGates) {
    if (!releaseGateSatisfied(releaseReadiness?.[key])) {
      errors.push(message);
    }
  }

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
    case "status":
      await printStatus();
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
