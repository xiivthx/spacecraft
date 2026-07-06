---
name: sc-mission
description: Manage Spacecraft mission artifacts and lifecycle for local OpenCode development
license: MIT
compatibility: opencode
---
- Read `.space/current`.
- Read `mission.json`, `spec.md`, `questions.md`, `decisions.md`, `plan.json`, design artifacts, `evidence.jsonl`, and `review.json` when available.
- sc-mission owns lifecycle but must route ambiguity to sc-clarify.
- Enforce order: mission -> clarify -> spec -> design if needed -> plan -> work -> verify -> review -> ship.
- Do not skip clarification when user intent, scope, acceptance criteria, or design direction is materially ambiguous.
- Do not implement if spec or plan is missing.
- Use sc-git for git safety, release branching, commits, rebasing, merge, version bump, changelog/spec notes, and tagging.
- Mutating work must satisfy sc-git gates before implementation and before ship.
- Keep mission artifacts small and human-readable.
- Prefer explicit evidence over narrative claims.
- End mission responses with the next action and session advice: continue current chat for small adjacent steps, start a new session for phase changes, large next tasks, or context-heavy threads.
