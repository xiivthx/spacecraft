import { spawnSync } from "node:child_process";
import { mkdtemp, mkdir, writeFile } from "node:fs/promises";
import { tmpdir } from "node:os";
import path from "node:path";
import assert from "node:assert/strict";
import test from "node:test";
import { fileURLToPath } from "node:url";

const repoRoot = path.resolve(fileURLToPath(new URL("..", import.meta.url)));
const spacecraft = path.join(repoRoot, "scripts", "spacecraft.mjs");

async function createWorkspace() {
  return mkdtemp(path.join(tmpdir(), "spacecraft-resolver-"));
}

async function writeMission(cwd, {
  id,
  title,
  state = "implementing",
  git = {}
}) {
  const dir = path.join(cwd, ".space", "missions", id);
  await mkdir(dir, { recursive: true });
  await writeFile(path.join(dir, "mission.json"), `${JSON.stringify({
    id,
    title,
    state,
    createdAt: "2026-07-07T00:00:00.000Z",
    updatedAt: "2026-07-07T00:00:00.000Z",
    git
  }, null, 2)}\n`);
}

async function writeCurrent(cwd, id) {
  await mkdir(path.join(cwd, ".space"), { recursive: true });
  await writeFile(path.join(cwd, ".space", "current"), `${id}\n`);
}

async function writeSession(cwd, key, id) {
  await mkdir(path.join(cwd, ".space", "sessions"), { recursive: true });
  await writeFile(path.join(cwd, ".space", "sessions", `${key}.current`), `${id}\n`);
}

function cleanEnv(extra = {}) {
  const env = { ...process.env };
  delete env.SPACECRAFT_SESSION;
  delete env.OPENCODE_SESSION_ID;
  delete env.CODEX_SESSION_ID;
  delete env.SPACECRAFT_MISSION;

  for (const [key, value] of Object.entries(extra)) {
    if (value === undefined || value === null) {
      delete env[key];
    } else {
      env[key] = String(value);
    }
  }
  return env;
}

function run(command, args, { cwd, env = {}, check = true } = {}) {
  const result = spawnSync(command, args, {
    cwd,
    env: cleanEnv(env),
    encoding: "utf8"
  });

  if (check) {
    assert.equal(result.status, 0, [
      `${command} ${args.join(" ")} failed`,
      `stdout:\n${result.stdout}`,
      `stderr:\n${result.stderr}`
    ].join("\n"));
  }
  return result;
}

function runSpacecraft(cwd, args, options = {}) {
  return run(process.execPath, [spacecraft, ...args], { cwd, ...options });
}

function resolveJson(cwd, args = [], options = {}) {
  const result = runSpacecraft(cwd, ["resolve", ...args, "--json"], options);
  return JSON.parse(result.stdout);
}

function runGit(cwd, args) {
  return run("git", args, { cwd });
}

function ids(records) {
  return records.map((record) => record.id).sort();
}

async function initGitOnBranch(cwd, branch) {
  runGit(cwd, ["init", "-b", "main"]);
  runGit(cwd, ["checkout", "-b", branch]);
}

async function initGitOnMain(cwd) {
  runGit(cwd, ["init", "-b", "main"]);
}

test("explicit selectors override lower-priority signals", async () => {
  const cwd = await createWorkspace();
  const alpha = "M-20260707-000001";
  const beta = "M-20260707-000002";

  await writeMission(cwd, { id: alpha, title: "Alpha mission" });
  await writeMission(cwd, { id: beta, title: "Beta target" });
  await writeCurrent(cwd, alpha);
  await writeSession(cwd, "case-session", alpha);
  await initGitOnBranch(cwd, `feat/${alpha.toLowerCase()}-alpha`);

  const byTitle = resolveJson(cwd, ["Beta target"], {
    env: { SPACECRAFT_SESSION: "case-session" }
  });
  assert.equal(byTitle.selected.id, beta);
  assert.equal(byTitle.source, "selector");
  assert.equal(byTitle.safety, "safe");

  const byEnv = resolveJson(cwd, [], {
    env: {
      SPACECRAFT_MISSION: beta,
      SPACECRAFT_SESSION: "case-session"
    }
  });
  assert.equal(byEnv.selected.id, beta);
  assert.equal(byEnv.source, "SPACECRAFT_MISSION");
  assert.equal(byEnv.safety, "safe");
});

test("session binding has priority over branch and current signals", async () => {
  const cwd = await createWorkspace();
  const alpha = "M-20260707-000011";
  const beta = "M-20260707-000012";

  await writeMission(cwd, { id: alpha, title: "Branch mission" });
  await writeMission(cwd, { id: beta, title: "Session mission" });
  await writeCurrent(cwd, alpha);
  await writeSession(cwd, "case-session", beta);
  await initGitOnBranch(cwd, `feat/${alpha.toLowerCase()}-branch`);

  const resolution = resolveJson(cwd, [], {
    env: { SPACECRAFT_SESSION: "case-session" }
  });
  assert.equal(resolution.selected.id, beta);
  assert.equal(resolution.source, "session");
  assert.equal(resolution.safety, "conflict");
  assert.ok(resolution.conflicts.some((conflict) => conflict.type === "signal-mismatch"));
});

test("branch mission id overrides current fallback and conflict blocks write paths", async () => {
  const cwd = await createWorkspace();
  const alpha = "M-20260707-000021";
  const beta = "M-20260707-000022";

  await writeMission(cwd, { id: alpha, title: "Current mission" });
  await writeMission(cwd, { id: beta, title: "Branch mission" });
  await writeCurrent(cwd, alpha);
  await initGitOnBranch(cwd, `feat/${beta.toLowerCase()}-branch`);

  const resolution = resolveJson(cwd);
  assert.equal(resolution.selected.id, beta);
  assert.equal(resolution.source, "branch");
  assert.equal(resolution.safety, "conflict");

  const blocked = runSpacecraft(cwd, ["validate"], { check: false });
  assert.notEqual(blocked.status, 0);
  assert.match(blocked.stderr, /cannot safely choose a mission for validate/);
  assert.match(blocked.stderr, /mission signals disagree/);
});

test("branch metadata resolves branches that do not include a mission id", async () => {
  const cwd = await createWorkspace();
  const alpha = "M-20260707-000031";
  const beta = "M-20260707-000032";
  const branch = "feature/shared-resolver-branch";

  await writeMission(cwd, { id: alpha, title: "Other active mission" });
  await writeMission(cwd, {
    id: beta,
    title: "Metadata mission",
    git: { workBranch: branch }
  });
  await initGitOnBranch(cwd, branch);

  const resolution = resolveJson(cwd);
  assert.equal(resolution.selected.id, beta);
  assert.equal(resolution.source, "branch-metadata");
  assert.equal(resolution.safety, "safe");
});

test("explicit and current selection resolve duplicate branch metadata", async () => {
  const cwd = await createWorkspace();
  const alpha = "M-20260707-000033";
  const beta = "M-20260707-000034";
  const branch = "feature/shared-work-branch";

  await writeMission(cwd, {
    id: alpha,
    title: "Alpha duplicate metadata",
    git: { workBranch: branch }
  });
  await writeMission(cwd, {
    id: beta,
    title: "Beta duplicate metadata",
    git: { workBranch: branch }
  });
  await initGitOnBranch(cwd, branch);

  const unresolved = resolveJson(cwd);
  assert.notEqual(unresolved.safety, "safe");
  assert.ok(unresolved.conflicts.some((conflict) => conflict.type === "ambiguous-signal"));

  const byEnv = resolveJson(cwd, [], {
    env: { SPACECRAFT_MISSION: beta }
  });
  assert.equal(byEnv.selected.id, beta);
  assert.equal(byEnv.source, "SPACECRAFT_MISSION");
  assert.equal(byEnv.safety, "safe");

  const suggested = runSpacecraft(cwd, ["git-suggest"], {
    env: { SPACECRAFT_MISSION: beta }
  });
  assert.match(suggested.stdout, /^Branch:/m);

  runSpacecraft(cwd, ["use", "Beta duplicate metadata"]);
  const byCurrent = resolveJson(cwd);
  assert.equal(byCurrent.selected.id, beta);
  assert.equal(byCurrent.source, ".space/current");
  assert.equal(byCurrent.safety, "safe");
});

test("single-active fallback does not resolve duplicate branch metadata", async () => {
  const cwd = await createWorkspace();
  const active = "M-20260707-000035";
  const shipped = "M-20260707-000030";
  const branch = "feature/shared-work-branch";

  await writeMission(cwd, {
    id: active,
    title: "Only active duplicate metadata",
    git: { workBranch: branch }
  });
  await writeMission(cwd, {
    id: shipped,
    title: "Shipped duplicate metadata",
    state: "shipped",
    git: { workBranch: branch }
  });
  await initGitOnBranch(cwd, branch);

  const resolution = resolveJson(cwd);
  assert.equal(resolution.selected.id, active);
  assert.equal(resolution.source, "single-active");
  assert.equal(resolution.safety, "conflict");
  assert.ok(resolution.conflicts.some((conflict) => conflict.type === "ambiguous-signal"));
});

test("creation branch metadata does not make common branches ambiguous", async () => {
  const cwd = await createWorkspace();
  const alpha = "M-20260707-000036";
  const beta = "M-20260707-000037";

  await writeMission(cwd, {
    id: alpha,
    title: "First mission from main",
    git: { branch: "main" }
  });
  await writeMission(cwd, {
    id: beta,
    title: "Current mission from main",
    git: { branch: "main" }
  });
  await writeCurrent(cwd, beta);
  await initGitOnMain(cwd);

  const resolution = resolveJson(cwd);
  assert.equal(resolution.selected.id, beta);
  assert.equal(resolution.source, ".space/current");
  assert.equal(resolution.safety, "safe");
  assert.ok(!resolution.signals.some((signal) => signal.source === "branch-metadata"));
});

test(".space/current is used as a fallback when no stronger signal exists", async () => {
  const cwd = await createWorkspace();
  const alpha = "M-20260707-000041";
  const beta = "M-20260707-000042";

  await writeMission(cwd, { id: alpha, title: "Alpha mission" });
  await writeMission(cwd, { id: beta, title: "Current mission" });
  await writeCurrent(cwd, beta);

  const resolution = resolveJson(cwd);
  assert.equal(resolution.selected.id, beta);
  assert.equal(resolution.source, ".space/current");
  assert.equal(resolution.safety, "safe");
});

test("single active mission is used as the final fallback", async () => {
  const cwd = await createWorkspace();
  const shipped = "M-20260707-000051";
  const active = "M-20260707-000052";

  await writeMission(cwd, { id: shipped, title: "Old mission", state: "shipped" });
  await writeMission(cwd, { id: active, title: "Only active mission" });

  const resolution = resolveJson(cwd);
  assert.equal(resolution.selected.id, active);
  assert.equal(resolution.source, "single-active");
  assert.equal(resolution.safety, "safe");
});

test("multiple active missions without a signal are ambiguous candidates", async () => {
  const cwd = await createWorkspace();
  const alpha = "M-20260707-000061";
  const beta = "M-20260707-000062";

  await writeMission(cwd, { id: alpha, title: "Alpha mission" });
  await writeMission(cwd, { id: beta, title: "Beta mission" });

  const resolution = resolveJson(cwd);
  assert.equal(resolution.selected, null);
  assert.equal(resolution.safety, "ambiguous");
  assert.deepEqual(ids(resolution.candidates), [alpha, beta]);
});

test("git-suggest blocks unsafe mission resolution", async () => {
  const cwd = await createWorkspace();
  const alpha = "M-20260707-000071";
  const beta = "M-20260707-000072";

  await writeMission(cwd, { id: alpha, title: "Current mission" });
  await writeMission(cwd, { id: beta, title: "Branch mission" });
  await writeCurrent(cwd, alpha);
  await initGitOnBranch(cwd, `feat/${beta.toLowerCase()}-branch`);

  const blocked = runSpacecraft(cwd, ["git-suggest"], { check: false });
  assert.notEqual(blocked.status, 0);
  assert.match(blocked.stderr, /cannot safely choose a mission for git-suggest/);
  assert.match(blocked.stderr, /mission signals disagree/);
  assert.doesNotMatch(blocked.stdout, /^Branch:/m);
});

test("conflict candidate numbers match use selector ordering with shipped missions", async () => {
  const cwd = await createWorkspace();
  const currentShipped = "M-20260707-000080";
  const branchActive = "M-20260707-000081";
  const unrelatedShipped = "M-20260707-000083";

  await writeMission(cwd, { id: currentShipped, title: "Current shipped", state: "shipped" });
  await writeMission(cwd, { id: branchActive, title: "Branch active" });
  await writeMission(cwd, { id: unrelatedShipped, title: "Unrelated shipped", state: "shipped" });
  await writeCurrent(cwd, currentShipped);
  await initGitOnBranch(cwd, `feat/${branchActive.toLowerCase()}-branch`);

  const blocked = runSpacecraft(cwd, ["git-suggest"], { check: false });
  assert.notEqual(blocked.status, 0);
  assert.match(blocked.stderr, /cannot safely choose a mission for git-suggest/);
  assert.match(blocked.stderr, new RegExp(`^1\\. Branch active \\(${branchActive}\\) - state:implementing`, "m"));
  assert.match(blocked.stderr, new RegExp(`^3\\. Current shipped \\(${currentShipped}\\) - state:shipped`, "m"));
  assert.doesNotMatch(blocked.stderr, /^2\. Current shipped/m);

  const selected = runSpacecraft(cwd, ["use", "3"]);
  assert.match(selected.stdout, new RegExp(`Selected mission: Current shipped \\(${currentShipped}\\)`));
});
