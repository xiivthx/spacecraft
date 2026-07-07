---
description: Implement the next smallest task in the resolved Spacecraft mission
agent: sc-commander
---
Use sc-mission, sc-clarify, and sc-git.
Run:
scripts/spacecraft resolve --json
If resolver safety is not `safe` or no mission is selected, stop before editing. Show the conflict/candidates and tell the user to run `scripts/spacecraft missions` then `scripts/spacecraft use <number|id|title>`, or set `SPACECRAFT_MISSION=<mission-id>` for one command.
Treat `.space/current` as fallback state, not sole authority.
Read the resolved mission's mission.json, spec.md, and plan.json.
If spec.md or plan.json is missing, stop and tell the user to run /sc-start or /sc-plan.
Read questions.md and decisions.md when present.
If blocking clarification remains open, stop.
Do not implement product code until clarification status is clear or explicitly deferred by the user.
If the user explicitly chooses to defer a decision, record the deferral in decisions.md and keep the implementation limited to unaffected tasks.
Run:
scripts/spacecraft git-info
If the workspace is not a git worktree, stop before editing product files unless the user has explicitly accepted no-git implementation risk for this mission in decisions.md.
If git exists, inspect dirty state before edits. Work with user changes; do not revert unrelated changes.
For large, risky, or multi-session slices, prefer a separate branch or git worktree.
Use sc-git branch naming when suggesting a branch:
<type>/<id>/<title>
Use:
scripts/spacecraft git-suggest <type> <slug>
to generate the suggested branch and commit examples.
Never edit product files directly on main.
If the current branch is main and the task clearly requires edits, create or switch to a feature/issue branch without asking again. If branch intent is unclear, use `scripts/spacecraft git-suggest`.
The agent may create checkpoint commits only on a valid non-main work branch.
Before staging or committing, ensure `.gitignore` is current and unsafe files are not staged.
Do not stage secrets, local env files, private data, dependency folders, build outputs, caches, logs, local databases, or machine-specific files.
Before code or dependency changes, check official current docs/registry/releases for direct dependencies and framework APIs. Use latest stable direct versions unless a deep dependency, ecosystem pin, or explicit user instruction says otherwise. Record source/version/date when it affects implementation.
Use rtk for noisy shell commands when available: installed hook, `rtk <supported command...>`, or `rtk proxy <command...>` for passthrough/tracking. Raw commands are allowed when rtk is missing or exact output is needed.
Implement only the next smallest pending task.
If implementing UI, read DESIGN.md and use sc-design.
If UI mood/theme/art direction is not recorded in decisions.md or plan.json, stop and recommend /sc-design before implementation.
Implement only the next planned UI slice.
Do not invent a new visual language when DESIGN.md exists.
Prefer CSS custom properties and local component styles.
Do not add broad styling frameworks unless explicitly requested.
Prefer test-first work where practical.
Before ending, run a bounded lightweight self-review/self-test loop over the touched diff:
- make at most two short self-review passes unless the user explicitly asks for more
- in each pass, check the change against the current task and acceptance checks
- remove accidental debug code, unrelated edits, dead code, and noisy formatting churn
- check obvious edge cases, error/loading states, accessibility labels, and mobile layout when relevant
- run the nearest cheap focused self-test when practical, such as a touched unit test, typecheck, lint, or targeted build
- if the pass finds small, low-risk issues, fix them and repeat the loop once
- stop when a pass finds no small self-found issues, or when the next fix would be broad, risky, design-level, blocked by ambiguity, or better suited to independent review
Keep this self-review lightweight. Do not invoke sc-reviewer, do not write review.md or review.json, and do not treat it as independent expert review.
If the loop stops with known remaining risk, state that risk and recommend the next concrete command or clarification.
If a useful self-test is skipped, say why and recommend the next concrete verification command.
Update plan.json task status conservatively.
Do not claim completion until /sc-verify captures evidence.
End with next action and session advice. Prefer continuing for immediate /sc-verify, or `/sc-flow` when the user wants the runner to verify and continue through remaining tasks.
