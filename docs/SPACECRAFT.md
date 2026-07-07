# Spacecraft

Spacecraft is a lean, local-first OpenCode harness for mission-driven software development. It is not a web server, dashboard, background service, database, or cloud sync tool.

Persona lives in `PERSONA.md`. Operational rules live in `AGENTS.md`.

Spacecraft is operated from the OpenCode UI or TUI with slash commands:

- `/sc-start`
- `/sc-clarify`
- `/sc-design`
- `/sc-plan`
- `/sc-git`
- `/sc-work`
- `/sc-verify`
- `/sc-flow`
- `/sc-design-review`
- `/sc-polish`
- `/sc-review`
- `/sc-status`
- `/sc-ship`

The helper script at `scripts/spacecraft.mjs` exists only to create file-backed mission artifacts, update mission state, and capture verification evidence. OpenCode commands call it through bash during normal usage. You may run it directly for smoke testing or debugging.

## File Layout

- `AGENTS.md`: project rules for Spacecraft work.
- `PERSONA.md`: commander tone and communication rules.
- `SPEC.md`: English-only project-level specification.
- `DESIGN.md`: Orbital Console visual source of truth for UI work.
- `opencode.json`: local OpenCode configuration and conservative permissions.
- `.opencode/agents/`: Spacecraft commander, planner, designer, and reviewer agents.
- `.opencode/commands/`: slash command prompts.
- `.opencode/skills/`: reusable Spacecraft workflow skills.
- `.opencode/skills/sc-design/scripts/serve-html.mjs`: local preview server for design HTML artifacts.
- `scripts/spacecraft.mjs`: dependency-free Node.js helper.
- `.space/current`: fallback selected mission id.
- `.space/sessions/`: local session-to-mission bindings when the host exposes a stable session key.
- `.space/missions/<mission-id>/`: mission artifacts and command output.
- `.space/archive/<mission-id>/`: compact shipped mission artifacts.
- `.space/missions/<mission-id>/questions.md`: open and answered clarification questions.
- `.space/missions/<mission-id>/decisions.md`: confirmed choices and assumptions.
- `.space/missions/<mission-id>/design/`: local HTML art-direction and design exploration artifacts.

## Commands

- `node scripts/spacecraft.mjs init`: create `.space/`, `.space/missions/`, and `.space/current`.
- `node scripts/spacecraft.mjs new "<title>"`: create a new mission and select it.
- `node scripts/spacecraft.mjs current`: print the `.space/current` fallback mission id.
- `node scripts/spacecraft.mjs resolve [selector] [--json]`: resolve the active mission from explicit, session, branch, metadata, current, and single-active signals.
- `node scripts/spacecraft.mjs missions`: list missions, selected signal, safety, branch hints, and next command.
- `node scripts/spacecraft.mjs use <number|id|title>`: select a mission by list number, mission id, exact title, or unique title substring.
- `node scripts/spacecraft.mjs bind-branch [selector]`: record the current branch on a mission without switching branches.
- `node scripts/spacecraft.mjs status`: print resolved mission state, resolver source/safety, artifact paths, task count, evidence count, and review status.
- `node scripts/spacecraft.mjs flow [--json]`: print workflow readiness, blockers, next task, and checkpoint policy.
- `node scripts/spacecraft.mjs git-info`: print git worktree, branch, HEAD, and dirty status.
- `node scripts/spacecraft.mjs git-suggest [type] [slug]`: print recommended release-branching branch name and Conventional Commit examples.
- `node scripts/spacecraft.mjs set-state <state>`: set mission state.
- `node scripts/spacecraft.mjs clarify-status <open|clear|deferred>`: set mission clarification status.
- `node scripts/spacecraft.mjs evidence "<label>" -- <command>`: run a command and record stdout, stderr, exit code, and evidence metadata.
- `node scripts/spacecraft.mjs validate`: validate the resolved mission artifacts.
- `node scripts/spacecraft.mjs closeout-check`: block or confirm release closeout readiness.
- `node scripts/spacecraft.mjs archive [selector]`: compact a shipped mission into `.space/archive/` and remove the live mission copy.
- `npm run sc:git`: shortcut for git safety status.
- `npm run sc:missions`: shortcut for mission list and resolver safety.
- `npm run sc:use -- <number|id|title>`: shortcut for selecting a mission.
- `npm run sc:flow`: shortcut for workflow readiness.
- `npm run sc:git:suggest -- [type] [slug]`: shortcut for branch and commit suggestions.
- `npm run sc:archive -- [selector]`: shortcut for mission archive compaction.
- `npm run sc:closeout`: shortcut for release closeout readiness.
- `node .opencode/skills/sc-design/scripts/serve-html.mjs [artifact-or-dir] --open`: serve and open design HTML artifacts.
- `npm run sc:design:open -- [artifact-or-dir]`: shortcut for the same design preview server.

## Mission Resolution

Spacecraft can have multiple active missions. Agents should resolve the selected mission with the helper instead of trusting `.space/current` alone.

Resolver priority:

1. explicit selector argument or `SPACECRAFT_MISSION`
2. local session binding from `SPACECRAFT_SESSION`, `OPENCODE_SESSION_ID`, or `CODEX_SESSION_ID`
3. mission id embedded in the current git branch name
4. branch metadata recorded in `mission.json`
5. `.space/current` fallback
6. single active mission fallback

Strong signals are session, branch id, branch metadata, and `.space/current`. If strong signals point to different missions, point to a missing mission, or branch metadata matches multiple missions, resolver safety becomes `conflict` or `ambiguous`. Write paths such as `set-state`, `clarify-status`, `evidence`, `validate`, `closeout-check`, and unqualified `bind-branch` must block until the user selects a mission.

Use `node scripts/spacecraft.mjs missions` when selection is unclear. Users can choose without knowing mission ids:

```sh
node scripts/spacecraft.mjs use 2
node scripts/spacecraft.mjs use "Add multi-mission active detection"
node scripts/spacecraft.mjs use "multi-mission"
```

An explicit selector or `SPACECRAFT_MISSION=<mission-id>` is an advanced one-command override. `.space/current` remains useful as portable fallback state, but it is not the only source of active mission truth.

## Id Format

New mission and evidence ids are compact sortable ids with no separator:

```text
M07FYB5W5
E07FYB5W5
```

The first character is the kind (`M` for mission, `E` for evidence). The remaining 8 characters are fixed-width uppercase base36 milliseconds since `2026-01-01T00:00:00.000Z`. This preserves lexicographic time ordering and covers about 89 years. Legacy ids such as `M-20260707-141230` remain valid for old missions, branch names, and session/current files.

## Mission Lifecycle

1. `/sc-start` creates a mission, drafts a minimal `spec.md`, and identifies blocking ambiguity.
2. `/sc-clarify` resolves one blocking question at a time.
3. `/sc-design` clears the main UI direction when the mission includes screens or components.
4. `/sc-plan` creates or updates `plan.json`.
5. `/sc-git` prepares or reviews branch, commit, release, merge, and tag policy when implementation is next.
6. `/sc-work` implements the next smallest pending task and runs a lightweight self-review/self-test to catch obvious issues.
7. `/sc-verify` captures command evidence.
8. `/sc-flow` can run the safe loop of `/sc-work Txx`, `/sc-verify Txx`, checkpoint commit, and next task until a real gate blocks.
9. `/sc-design-review` checks visual quality and anti-slop issues when UI changed.
10. `/sc-polish` performs focused design cleanup before shipping UI work.
11. `/sc-review` reviews the diff and evidence.
12. `/sc-ship` or a user request to ship/release/merge/finish the mission/close the branch runs release closeout. It ships only when evidence, clarification, git, release, and review gates pass, then archives the shipped mission unless the user asks to keep it live.

Do not implement product code before `spec.md` and `plan.json` exist. Do not claim work is done, verified, ready, or shipped without evidence in `evidence.jsonl`.

When the user clearly asks for mutating work and no suitable mission or non-main branch exists, Spacecraft may create the mission and branch without asking another blocking question. This keeps flow moving while preserving the spec, plan, git, evidence, and review gates.

## Workflow Runner

`/sc-flow` reduces repeated HIL during implementation. It uses mission artifacts as workflow memory and continues through:

```text
/sc-work Txx -> /sc-verify Txx -> checkpoint commit -> next task
```

The runner stops when a real gate blocks progress: unsafe resolver state, open blocking clarification, missing spec or plan, missing design direction for UI work, dependency/API freshness work, main-branch write risk, dirty or unsafe files, failed verification, failed validation, critical review findings, release actions, or context that is too heavy for safe continuation.

Checkpoint commits are local rollback points created only on a valid non-main work branch after passing evidence exists for the task. Before `/sc-ship`, checkpoint commits should be squashed or fixed up into logical Conventional Commits.

## Subagent Invocation Model

Some slash commands include required read-only subagents as part of their normal workflow. Invoking these commands is explicit permission to use the named subagent; do not ask for separate subagent permission:

- `/sc-plan`: use `sc-planner` to draft the plan.
- `/sc-design`: use `sc-designer` to shape UI direction.
- `/sc-design-review`: use `sc-designer` to critique UI/design quality.
- `/sc-polish`: use `sc-designer` to identify focused polish items before small UI cleanup.
- `/sc-review`: use `sc-reviewer` for independent review of diff and evidence.

When `/sc-review` covers UI changes, it may also use a focused read-only `sc-designer` sidecar for design-risk triage without asking again. For a full design critique, use `/sc-design-review`.

Other commands should not spawn subagents unless the user explicitly asks for delegation or the command is updated to make that subagent part of its contract.

## Git And Release Branching

Git is the default rollback boundary for implementation work. Spacecraft should allow discovery, clarification, design exploration, planning, and read-only review without git, but `/sc-work` must not edit product files outside a git worktree unless the user explicitly accepts no-git implementation risk for that mission.

Git policy lives in `sc-git`. Spacecraft may create a non-main branch when the user's mutating intent is clear and policy permits it. Spacecraft does not auto-run `git init`, create worktrees, rebase, merge, tag, or push unless the user explicitly asks or the request is a release closeout/ship flow that authorizes release prep. Push always requires explicit user request.

Spacecraft uses release branching:

- `main` is protected and release-only.
- Do not write product changes directly on `main`.
- Use one non-main branch per feature, fix, issue, or tightly scoped change.
- Branch from the latest `main`.
- Branch names should follow `<type>/<issue-or-mission>-<slug>`, for example `feat/m07fyb5w5-workflow-runner`.
- Use `release/v<major>.<minor>.<patch>` only for release preparation work.
- Merge only after verification, git, release, and review gates pass.

For large, risky, or multi-session implementation slices, prefer a separate branch or git worktree:

```sh
git worktree add ../spacecraft-flow -b feat/m07fyb5w5-workflow-runner main
```

For small adjacent edits in an existing git worktree, the current branch is acceptable when the dirty state has been inspected and unrelated user changes are preserved.

New missions record the current git base sha when git is available. `/sc-status` reports the current git state. `/sc-ship` should include a rollback story: git base sha, branch/worktree, or explicit no-git risk acceptance recorded in `decisions.md`.

The agent may commit frequently inside a valid non-main work branch. Before merge, squash or fixup checkpoint commits into logical commits. Target 1 to 3 final commits per branch; a branch merged into `main` should not exceed 5 final commits unless explicitly justified. If more than 5 final commits seem necessary, split the feature earlier.

Keep `.gitignore` current before staging, committing, or merging. When new tools, frameworks, generated outputs, caches, logs, env files, local databases, or machine-specific files appear, update ignore rules before staging. Do not ignore source files, migrations, lockfiles, intended configs, tests, specs, changelog entries, or release notes just to make status clean.

Never allow secrets, credentials, local env files, private data, dependency folders, build outputs, caches, logs, local databases, or machine-specific files into git or public release artifacts. If unsure whether a file is safe to track, ask before staging it.

Commit messages should follow Conventional Commits:

```text
<type>[optional scope]: <description>
```

Common types are `feat`, `fix`, `docs`, `refactor`, `test`, `build`, `ci`, `chore`, `perf`, `style`, and `revert`.

Use `!` after the type or scope, or a `BREAKING CHANGE:` footer, for breaking changes. Keep the subject imperative, lowercase, and short. Put mission ids or evidence ids in the body/footer when useful rather than forcing them into the subject.

Before merge:

- rebase the work branch on latest `main`
- re-run tests, verification, and validation after the latest rebase
- confirm `.gitignore` is current and no unsafe files are staged or tracked accidentally
- bump version unless explicitly deferred with rationale
- update changelog
- update a short spec/release note when product behavior changed

Merge into `main` only with:

```sh
git merge --no-ff <branch>
```

After merge, create an annotated version tag:

```sh
git tag -a v<major>.<minor>.<patch> -m "v<major>.<minor>.<patch>"
```

Do not push commits or tags unless the user explicitly asks.

After a successful merge into `main`, delete the merged local branch unless the user asks to keep it.

## Session Handoff

Session handoff means the chat can stop while the branch or mission remains unfinished.

Use session handoff when:

- the user asks to stop this chat, end the session, or continue in a new session
- work is mid-task, unreviewed, dirty, blocked, or not release-ready
- Spacecraft itself recommends a new session because context is heavy or the phase changed

Handoff must include:

- resolved mission and branch
- state, blockers, and dirty git status
- exact pickup command, usually `/sc-status` followed by the next intended command

Handoff must not merge, tag, delete branches, or claim release readiness.

If "close session" is ambiguous, default to handoff. If the mission appears release-ready, recommend `/sc-ship` as the next command instead of merging automatically.

## Release Closeout

Release closeout means the work is intended to become releasable on `main`.

Use release closeout for `/sc-ship` or explicit requests to ship, release, merge, finish the mission, or close the branch.

Release closeout must check:

- resolved mission, spec, plan, questions, decisions, evidence, and review
- git branch, dirty state, final commit count, Conventional Commit subjects, and rebase status
- version bump, changelog/spec release notes, no-ff merge plan, and tag plan
- verification evidence after the latest rebase
- no critical review or design findings

`review.json.releaseReadiness` records release gates for `closeout-check`:

- `version`
- `changelog`
- `specNote`
- `tagPlan`
- `postRebaseVerification`

Each gate must be an object with a `status`, such as `{ "status": "completed" }`. Accepted statuses include complete, completed, checked, updated, passed, bumped, present, done, and deferred. The `planned` status is only accepted for the `tagPlan` gate. Deferred gates need a non-empty `rationale`. Strings and booleans are not valid gates.

If any gate is missing, Spacecraft blocks release closeout and names the exact missing action. It does not call work done.

If gates pass, Spacecraft prepares merge into `main` using `git merge --no-ff <branch>`, creates the version tag after merge when needed, cleans up the merged branch, and compacts the shipped mission into `.space/archive/` unless the user asks to keep the full live mission folder. Push still requires explicit user request.

## Mission Archive

`node scripts/spacecraft.mjs archive [selector]` moves a shipped mission from `.space/missions/` to `.space/archive/` after release-readiness gates are recorded. It writes compact durable artifacts:

- `SUMMARY.md`
- compact `mission.json`
- compact `plan.json`
- compact `evidence.jsonl` without stdout/stderr output files
- review, spec, decisions, and questions files when present

The archive command refuses non-shipped missions, missing or incomplete plans, missing evidence, non-ready reviews, critical findings, and incomplete release-readiness gates. It also clears `.space/current` and session bindings that still point to the archived mission.

## Common Scenarios

- Mid-task, user says "open a new session": handoff only, no release.
- Plan done, implementation next, context heavy: recommend new session with `/sc-status` then `/sc-work`; no release.
- Verification failed or review blocked: handoff or fix loop; release closeout blocks.
- Review ready, evidence fresh, branch clean, user says "ship": release closeout, rebase, verify, no-ff merge, tag, cleanup branch.
- User says "close branch" while tasks are incomplete: release closeout check blocks and lists missing tasks.
- User says only "close session" after work seems ready: handoff plus recommend `/sc-ship`; no merge without release intent.

## Dependency Freshness

Before code work that chooses direct dependencies, framework versions, generated scaffolds, or current APIs, Spacecraft checks official docs, registries, or releases. It uses latest stable direct versions unless a deep dependency, ecosystem pin, lockfile constraint, security advisory, or explicit user instruction says otherwise.

Record source, version, and date in `decisions.md` or `evidence.jsonl` when the choice affects implementation.

## Shell Output

Use rtk as the shell-output proxy when available:

- prefer installed rtk hooks for automatic command rewrite
- use `rtk <supported command...>` for noisy commands such as git, grep, tests, build, lint, docker, and cargo
- use `rtk proxy <command...>` for unsupported commands when passthrough plus tracking is useful
- use raw commands when exact byte output is required or rtk is unavailable

Do not use rtk to bypass denied git, push, destructive, or secret-touching operations.

## Clarification Workflow

`/sc-start` may stop after creating a mission if the request has blocking ambiguity. Spacecraft should ask one blocking question at a time, include why it matters, provide a recommended answer, and explain what happens if the recommendation is accepted.

`/sc-clarify` resolves one question at a time. It reads the resolved mission, inspects repo files when they can answer the question, records open and answered questions in `questions.md`, and records confirmed choices or low-risk assumptions in `decisions.md`.

`/sc-design` must clarify visual and product direction before locking design. `/sc-plan` should not finalize a plan while blocking questions remain. `/sc-work` should not implement while blocking clarification remains open unless the user explicitly defers the decision and the work is limited to unaffected tasks.

`questions.md` stores open and answered questions. `decisions.md` stores confirmed choices and assumptions. Older missions may not have these files; commands should create them when needed.

## Design Workflow

`DESIGN.md` is the visual source of truth for Spacecraft-created web interfaces. The default design language is Orbital Console: local-first, precise, calm, technical, sparse, and usable before impressive.

Use `/sc-design` before implementing a new screen or major component. It must clear the main design image: mood, tone, theme, palette, composition, UX feel, art treatment, 3D, transition, animation, accessibility constraints, and the first UI surface to build.

`/sc-design` should default to one simple text question with a recommendation. When seeing options materially helps the user choose, it creates a static HTML artifact under `.space/missions/<mission-id>/design/`.

Do not create HTML for every config. Text-first decisions usually include user journey, navigation model, interaction rules, density intent, states, accessibility constraints, risks, and priorities. Visual artifacts are reserved for decisions where side-by-side comparison helps, such as layout, palette, typography feel, art treatment, 3D depth, and motion.

When the design feels weak, generic, or hard to imagine, `/sc-design` should run a reference scout before deeper config decisions. It should gather a small set of current public references, group them into 2 or 3 directions, and separate reference purpose: layout/template, mood/art, and interaction/motion. Bootstrap examples and template galleries are useful for structure; Pinterest and moodboards are useful for taste, palette, and art direction.

Selected references should be recorded in `decisions.md` or `.space/missions/<mission-id>/design/references.md` with source URLs, useful parts, what to borrow as a pattern, and what not to copy. References are calibration, not source material to clone.

After UI implementation, `/sc-design-review` can use selected references as a quality check: hierarchy, layout rhythm, density, palette discipline, art direction, and interaction feel. Review should compare against references without forcing exact imitation.

Design should proceed one config at a time where possible: concept, user journey, layout/IA, navigation, interaction, palette, typography, art/3D, motion, density, and states. This lets the user mix choices, such as layout A with palette B. Spacecraft should record each chosen config separately in `decisions.md`, then synthesize the final design brief into `spec.md` or `plan.json`.

Design options must be materially different. A valid option set should differ in product metaphor, information architecture, first screen layout, navigation model, interaction model, visual language, art treatment, motion model, or density. The same screen skeleton with swapped colors, labels, or copy does not count as design exploration.

HTML design artifacts should include a similarity audit explaining what is genuinely different, what is intentionally shared, and why the options are not just theme variations. If the options are too similar, Spacecraft should discard them and create a new set before asking the user to choose.

For Thai or multilingual missions, design artifacts should be Thai-first with simple English support. Headings and explanations should be easy Thai, with short English labels only where useful, such as `แผนที่ก่อน (Map-first)`.

Design artifacts should use Feynman-style clarity. The artifact should explain each option in plain language, use a familiar analogy when it helps, label the visual so the user knows what to look at, and state the gain/tradeoff of choosing that option. It should not require design theory knowledge.

Keep design artifacts compact. One artifact should answer one design config question where possible. Visible lists should have 3 bullets or fewer, paragraphs should be short, and extra theory should be omitted or placed in a small "เหตุผลสั้น ๆ (Why)" block.

Design artifacts are not product UI. They are local-first selection aids. They should be dependency-free, directly openable in a browser, and should not install Lavish, design packages, animation libraries, CSS frameworks, or external assets.

To preview design artifacts through a local server, run:

```sh
npm run sc:design:open -- .space/missions/<mission-id>/design/<artifact>.html
```

Without an artifact path, the preview server opens the resolved mission's `design/` folder:

```sh
npm run sc:design:open
```

Use `/sc-work` to implement the next planned UI slice and catch small issues with a lightweight self-review/self-test. Use `/sc-verify` to capture technical evidence. Use `/sc-design-review` to check hierarchy, layout, typography, spacing, color, interaction states, accessibility, responsiveness, Feynman clarity, and anti-slop issues. Use `/sc-polish` for final cleanup before `/sc-ship`.

Design review is not a replacement for browser or manual visual review unless browser tooling already exists. Spacecraft does not add screenshot automation by default.

## Evidence Model

Evidence is append-only JSON Lines in `evidence.jsonl`. Each entry records:

- evidence id
- label
- command
- exit code
- stdout file
- stderr file
- creation timestamp

Command output is stored under `outputs/` inside the mission directory.

## Review Model

`/sc-work` self-review is a local cleanup pass only. It should not write `review.md` or `review.json`, and it is not a replacement for independent review.

`/sc-review` must invoke the read-only `sc-reviewer` subagent. See the Subagent Invocation Model above for permission rules.

`review.md` stores the human-readable review. `review.json` stores structured review status and findings. Critical findings block `/sc-ship`.

## Response Closeout

Spacecraft command responses should end with a concrete next action and session advice.

Continue the current chat when the next step is small, adjacent, and depends on the current conversational context, such as answering a clarification question or immediately running `/sc-verify`.

Start a new session when the mission or major phase ended, the next step is a large implementation slice, the thread is context-heavy, or mission artifacts already contain enough state for a clean handoff.

When recommending a new session, Spacecraft should include a compact pickup instruction, usually `/sc-status` followed by the next intended command.

Response closeout and session handoff are not release closeout.

## Future Web Service Mission

To begin a later web service mission from OpenCode:

```text
/sc-start Build a local TypeScript Fastify web service from scratch with health, version, and a lean Orbital Console status page
```

If the user does not specify a stack, Spacecraft defaults future web service work to Node.js + TypeScript + Fastify + Vitest. Product dependencies may be installed only with user approval.

## Smoke Test Commands

```sh
node scripts/spacecraft.mjs init
node scripts/spacecraft.mjs new "Smoke test Spacecraft harness"
node scripts/spacecraft.mjs resolve --json
node scripts/spacecraft.mjs missions
node scripts/spacecraft.mjs status
node scripts/spacecraft.mjs flow
node scripts/spacecraft.mjs validate
```
