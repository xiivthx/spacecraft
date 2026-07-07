# Spacecraft Agent Rules

Spacecraft is a mission-control harness.

Persona: read [PERSONA.md](PERSONA.md).

- Use the commands: `/sc-start`, `/sc-clarify`, `/sc-design`, `/sc-plan`, `/sc-git`, `/sc-work`, `/sc-verify`, `/sc-flow`, `/sc-design-review`, `/sc-polish`, `/sc-review`, `/sc-status`, and `/sc-ship`.
- Commands that define required read-only subagents are explicit permission to use them without asking again: `/sc-plan` uses `sc-planner`; `/sc-design`, `/sc-design-review`, and `/sc-polish` use `sc-designer`; `/sc-review` uses `sc-reviewer`.
- `/sc-review` may also use a focused read-only `sc-designer` sidecar for UI design-risk triage without asking again.
- Other commands should not spawn subagents unless the user explicitly asks for delegation or the command is updated to make that subagent part of its contract.
- Always resolve the active mission before mission work. Use `scripts/spacecraft resolve [selector] [--json]`, `status`, or `missions`; `.space/current` is fallback state, not sole authority.
- Resolver priority is explicit selector or `SPACECRAFT_MISSION`, session binding, branch mission id, branch metadata, `.space/current`, then single active mission.
- Strong signal conflicts or ambiguous active missions block write, verify, review, ship, and git safety work until the mission is selected with `scripts/spacecraft use <number|id|title>` or an explicit selector.
- Users may choose missions by list number, mission id, exact title, or unique title substring; do not expect them to know a mission id.
- New mission and evidence ids are compact sortable ids with no hyphen, such as `M07FYB5W5`; legacy `M-YYYYMMDD-HHmmss` ids remain valid.
- Do not implement product code before `spec.md` and `plan.json` exist.
- Do not claim done, pass, verified, or ready without evidence in `evidence.jsonl`.
- Critical review findings block shipping.
- Prefer small tasks and focused verification.
- Use `/sc-flow` to continue `work -> verify -> checkpoint commit -> next task` loops in the same chat until a real gate blocks.
- Keep mission artifacts small and human-readable.
- Keep root `SPEC.md` as the English-only project-level specification.
- Keep prompts lean: only necessary words, commands, and gates.
- Use caveman-style brevity for nonessential communication; keep technical content exact.
- Use git as the default rollback boundary for implementation work.
- Discovery, clarification, design exploration, planning, and read-only review may run without git.
- `/sc-work` must not edit product files outside a git worktree unless the user explicitly accepts no-git implementation risk for the mission.
- Use `sc-git` for branch, commit, rebase, merge, version bump, changelog/spec, and tag policy.
- If clear mutating work is requested and no suitable mission or branch exists, create the mission and non-main branch without an extra blocking question when policy already permits it.
- Do not auto-run `git init`, create worktrees, rebase, merge, tag, or push unless the user explicitly asks.
- Never write product changes directly on `main`.
- Use Spacecraft release branching: one non-main branch per feature, fix, issue, or tightly scoped change.
- Branch names follow `<type>/<id>/<title>`, for example `feat/m07fp1l7z/go-rewrite`.
- The agent may commit frequently only inside a valid non-main work branch.
- After a task has passing verification evidence, `/sc-flow` may create a local checkpoint commit before starting the next task.
- Before merge, squash/fixup checkpoint commits into logical Conventional Commits.
- A branch merged to `main` should have 1 to 3 final commits and should not exceed 5 unless justified.
- Rebase the work branch on latest `main` before merge, then merge with `git merge --no-ff <branch>`.
- Release closeout is only for `/sc-ship`, ship/release/merge requests, finished-mission closeout, or closing a branch.
- Session handoff is not release: if the user wants to stop this chat or continue in a new session while work is unfinished, summarize state and pickup command without merging.
- After a successful merge to `main`, clean up the merged branch unless the user asks to keep it.
- After successful release closeout, archive shipped mission artifacts under `.space/archive/` unless the user asks to keep the full live mission folder.
- Keep `.gitignore` current before staging, committing, or merging.
- Never let secrets, local env files, private data, caches, logs, dependency folders, build outputs, or machine-specific files enter git/public artifacts.
- Test, verify, and validate after the latest rebase and before any merge into `main`.
- Before merge, bump version and update changelog/spec notes when behavior changed; after merge, create the version tag.
- Before code or dependency work, check official current docs/registry/releases for direct dependencies and framework APIs. Use latest stable direct versions unless a deep dependency, ecosystem pin, or explicit user instruction requires otherwise. Record source/version/date in decisions or evidence when it affects implementation.
- Use rtk as the shell-output proxy when available: prefer installed hooks or `rtk <supported command...>` for noisy commands, and `rtk proxy <command...>` for unsupported passthrough/tracking. Do not use rtk to bypass denied git or destructive operations. If rtk is missing or exact raw output is required, use raw commands and note the exception when relevant.
- Future product missions may install dependencies with user approval.
- End each Spacecraft session with a recommended next action and session advice: continue this chat for small adjacent steps, or start a new session when the phase changed, the thread is context-heavy, or mission artifacts are sufficient for handoff.

## Design discipline

- Read `DESIGN.md` before UI work.
- Use `/sc-design` before implementing a new screen or major component.
- `/sc-design` must clear the main design image: mood, theme, art direction, color, layout feel, 3D treatment, transition, and animation when relevant.
- `/sc-design` may create local HTML examples when visual comparison is clearer than text.
- `/sc-design` should default to normal chat questions; create HTML only when visual comparison materially helps.
- Preview design HTML with `node .opencode/skills/sc-design/scripts/serve-html.mjs [artifact-or-dir] --open`.
- `/sc-design` should ask one design config at a time so the user can mix choices such as layout A with palette B.
- If design direction is weak or generic, `/sc-design` should scout external references before deeper design config.
- References must be separated by purpose: layout/template, mood/art, and interaction/motion.
- References calibrate taste and structure; do not copy third-party designs, screenshots, copy, brand identity, or exact compositions.
- Design HTML must use Feynman-style clarity: plain Thai first, simple English labels only when useful, familiar analogies when helpful, labeled visuals, and explicit gain/tradeoff.
- Design HTML is a decision aid, not an essay. Keep copy short, avoid theory dumps, and make the visual explain what the user is choosing.
- Design options must differ in concept, IA, layout, interaction model, and art direction; color-swap options do not count.
- User-facing design HTML should be Thai-first with simple English labels when the user works in Thai or mixed Thai/English.
- Use `/sc-polish` before `/sc-ship` when UI changed.
- Use `/sc-design-review` for read-only design critique.
- UI work must avoid generic AI SaaS patterns.
- Design decisions must be captured in the mission spec or plan.
- Art-direction decisions must be captured in `decisions.md` before UI implementation.
- Do not introduce CSS frameworks unless explicitly requested.
- Use accessible semantic HTML and visible focus states.

## Clarification policy

Before planning, designing, or implementing, resolve meaningful ambiguity.

Use `sc-clarify` when the mission has unclear product behavior, scope, user intent, design direction, constraints, or acceptance criteria.

Do not ask questions that can be answered by reading the repo. Inspect files first.

Ask exactly one blocking question at a time. Include:
- the question
- why it matters
- the recommended answer
- what happens if the recommendation is accepted

Do not implement or finalize a plan/design while a blocking clarification question is open.

Low-risk ambiguity may be recorded as an assumption in `decisions.md`.

Confirmed user choices belong in `decisions.md`.
Open and answered questions belong in `questions.md`.

## Session handoff and release closeout

At the end of a Spacecraft command or work session, include:
- the recommended next action
- the exact next slash command when useful
- whether to continue in the current chat or start a new session

If work is unfinished and the user asks to end/close the session or start a new session, treat it as session handoff:
- summarize state, blockers, dirty git status, and exact pickup command
- do not merge, tag, delete branches, or claim release readiness
- if "close session" is ambiguous, default to handoff and recommend `/sc-ship` when release-ready

If the user asks to ship, release, merge, finish the mission, or close the branch, treat it as release closeout:
- check evidence, review, git branch, version/changelog/spec notes, rebase status, merge plan, and tag plan
- prepare merge to `main` only when all gates pass
- block and list exact missing actions when not ready
- after successful merge, delete the merged branch unless asked not to

Recommend continuing the current chat when:
- the next step is small and directly follows the current context
- the user is answering the current clarification question
- the mission is mid-task and context is still fresh

Recommend starting a new session when:
- the mission or major phase just ended
- the next task is a different phase or large implementation slice
- the chat is long, context-heavy, or token usage is high
- enough state is captured in mission artifacts for a clean handoff

When recommending a new session, provide a compact pickup instruction such as `/sc-status` followed by the next intended command.
