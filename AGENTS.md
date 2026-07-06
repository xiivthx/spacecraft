# Spacecraft Agent Rules

Spacecraft is a mission-control harness.

- Use the commands: `/sc-start`, `/sc-clarify`, `/sc-design`, `/sc-plan`, `/sc-git`, `/sc-work`, `/sc-verify`, `/sc-design-review`, `/sc-polish`, `/sc-review`, `/sc-status`, and `/sc-ship`.
- Always check `.space/current` when working inside a mission.
- Do not implement product code before `spec.md` and `plan.json` exist.
- Do not claim done, pass, verified, or ready without evidence in `evidence.jsonl`.
- Critical review findings block shipping.
- Prefer small tasks and focused verification.
- Keep mission artifacts small and human-readable.
- Use git as the default rollback boundary for implementation work.
- Discovery, clarification, design exploration, planning, and read-only review may run without git.
- `/sc-work` must not edit product files outside a git worktree unless the user explicitly accepts no-git implementation risk for the mission.
- Use `sc-git` for branch, commit, rebase, merge, version bump, changelog/spec, and tag policy.
- Do not auto-run `git init`, create branches/worktrees, rebase, merge, tag, or push unless the user explicitly asks.
- Never write product changes directly on `main`.
- Use Spacecraft release branching: one non-main branch per feature, fix, issue, or tightly scoped change.
- Branch names follow `<type>/<issue-or-mission>-<slug>`, for example `feat/m-20260706-120409-okinawa-planner-ui`.
- The agent may commit frequently only inside a valid non-main work branch.
- Before merge, squash/fixup checkpoint commits into logical Conventional Commits.
- A branch merged to `main` should have 1 to 3 final commits and should not exceed 5 unless justified.
- Rebase the work branch on latest `main` before merge, then merge with `git merge --no-ff <branch>`.
- Keep `.gitignore` current before staging, committing, or merging.
- Never let secrets, local env files, private data, caches, logs, dependency folders, build outputs, or machine-specific files enter git/public artifacts.
- Test, verify, and validate after the latest rebase and before any merge into `main`.
- Before merge, bump version and update changelog/spec notes when behavior changed; after merge, create the version tag.
- Future product missions may install dependencies with user approval.
- End each Spacecraft work session with a concrete next action and session advice.

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

## Session closeout

At the end of a Spacecraft command or work session, include:
- the recommended next action
- the exact next slash command when useful
- whether to continue in the current chat or start a new session

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
