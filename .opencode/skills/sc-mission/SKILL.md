---
name: sc-mission
description: Manage Spacecraft mission artifacts and lifecycle for local OpenCode development
license: MIT
compatibility: opencode
---
- Resolve the active mission with `scripts/spacecraft resolve [selector] [--json]`, `status`, or `missions`; `.space/current` is fallback state, not sole authority.
- Resolver priority is explicit selector or `SPACECRAFT_MISSION`, session binding, branch mission id, branch metadata, `.space/current`, then single active mission.
- Strong signal conflicts or ambiguous active missions block mission writes until the user selects with `scripts/spacecraft use <number|id|title>` or an explicit selector.
- Users may choose by list number, mission id, exact title, or unique title substring; do not expect the user to know a mission id.
- New mission and evidence ids are compact sortable ids with no hyphen, such as `M07FYB5W5`; legacy `M-YYYYMMDD-HHmmss` ids remain valid.
- Read the resolved mission's `mission.json`, `spec.md`, `questions.md`, `decisions.md`, `plan.json`, design artifacts, `evidence.jsonl`, and `review.json` when available.
- sc-mission owns lifecycle but must route ambiguity to sc-clarify.
- Enforce order: mission -> clarify -> spec -> design if needed -> plan -> work -> verify -> review -> ship. `/sc-flow` may repeat work -> verify -> checkpoint commit for successive tasks until a gate blocks.
- Do not skip clarification when user intent, scope, acceptance criteria, or design direction is materially ambiguous.
- If clear mutating work is requested and no suitable mission or branch exists, create the mission and non-main branch without another blocking question when policy permits it.
- Do not implement if spec or plan is missing.
- Use sc-git for git safety, release branching, commits, rebasing, merge, version bump, changelog/spec notes, and tagging.
- Mutating work must satisfy sc-git gates before implementation and before ship.
- Treat stop-chat/close-session/end-session/new-session requests as session handoff unless release intent is explicit.
- If "close session" is ambiguous and work appears ready, recommend `/sc-ship`; do not merge automatically.
- Treat ship/release/merge/finish-mission/close-branch requests as release closeout prep. Block closeout when gates are incomplete.
- After successful release closeout, archive shipped mission artifacts under `.space/archive/` unless the user asks to keep the full live mission folder.
- Keep mission artifacts small and human-readable.
- Prefer explicit evidence over narrative claims.
- End each Spacecraft session with a recommended next action and session advice: continue this chat for small adjacent steps, or start a new session when the phase changed, the thread is context-heavy, or mission artifacts are sufficient for handoff.
