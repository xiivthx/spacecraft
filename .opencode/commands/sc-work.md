---
description: Implement the next smallest task in the current Spacecraft mission
agent: sc-commander
---
Use sc-mission, sc-clarify, and sc-git.
Read current mission, spec.md, and plan.json.
If spec.md or plan.json is missing, stop and tell the user to run /sc-start or /sc-plan.
Read questions.md and decisions.md when present.
If blocking clarification remains open, stop.
Do not implement product code until clarification status is clear or explicitly deferred by the user.
If the user explicitly chooses to defer a decision, record the deferral in decisions.md and keep the implementation limited to unaffected tasks.
Run:
node scripts/spacecraft.mjs git-info
If the workspace is not a git worktree, stop before editing product files unless the user has explicitly accepted no-git implementation risk for this mission in decisions.md.
If git exists, inspect dirty state before edits. Work with user changes; do not revert unrelated changes.
For large, risky, or multi-session implementation slices, recommend a separate branch or git worktree before editing.
Use sc-git branch naming when suggesting a branch:
<type>/<issue-or-mission>-<slug>
Use:
node scripts/spacecraft.mjs git-suggest <type> <slug>
to generate the suggested branch and commit examples.
Never edit product files directly on main.
If the current branch is main, stop and ask to create or switch to a feature/issue branch.
The agent may create checkpoint commits only on a valid non-main work branch.
Before staging or committing, ensure `.gitignore` is current and unsafe files are not staged.
Do not stage secrets, local env files, private data, dependency folders, build outputs, caches, logs, local databases, or machine-specific files.
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
End with the recommended next action and session advice. Prefer continuing this chat for immediate /sc-verify; recommend a new session when the next task is a separate implementation slice or context is heavy.
