---
name: sc-mission
description: "Manage mission artifacts and lifecycle. Activate on /sc-run mission work, mission creation, status check, or lifecycle management."
---

# sc-mission

Manage mission artifacts and lifecycle for local Cursor development.

## When to use

Activate when the user asks to:

- Resolve, list, or select a mission
- Create mission artifacts (spec, plan, evidence)
- Check mission state or lifecycle
- Route ambiguity to clarification
- Handle session handoff or release closeout

## Workflow

Use this exact sequence unless the user specifies otherwise:

1. **Resolve mission** - Run `spacecraft resolve` (or `spacecraft status` / `spacecraft missions`); `.space/current` is fallback state, not sole authority. On conflict or ambiguity, use `spacecraft use <selector>`. Commander auto-detects development lane (Advisory, Mission, Debug, Quick) based on user intent.
2. **Read artifacts** - Read `mission.json`, `spec.md`, `questions.md`, `decisions.md`, `plan.json`, design artifacts, `evidence.jsonl`, and `review.json` when available.
3. **Index artifacts** - After creating or updating spec.md, plan.json, decisions.md, or questions.md, ctx_index them with source label `sc-memory/<mission-id>/<type>` (best-effort: warn on failure, never block). See sc-memory for conventions.
4. **Route ambiguity** - If intent, scope, or acceptance criteria is ambiguous, route to sc-clarify before proceeding.
5. **Enforce lifecycle** - Follow: mission -> clarify -> spec -> visual design if needed -> plan -> build -> verify -> review -> ship. For roadmap AFK, `/sc-run` orchestrates jigsaw plan → per-acceptance RED-GREEN (checkpoint commits) → combine/refactor → review.
6. **Release or handoff** - On ship intent, run `/sc-ship` closeout. On session end, give handoff summary.

### Edge cases

- **No mission exists and user wants mutating work** - Create mission with `spacecraft new "<title>"` and a work branch. Record in `decisions.md`: "Auto-created mission for: <reason>."
- **Multiple missions match selector** - Show candidates. Ask user to pick with `spacecraft use <number>`.
- **Resolve conflict or ambiguity** - Resolve via `spacecraft resolve`; on conflict/ambiguity use `spacecraft use <selector>`. Block mutating work until resolved.
- **Session ends mid-work** - Handoff: summarize current task, blockers, dirty git, next pickup command. Do not merge, tag, or delete branches.
- **State transition conflict** - If the current state does not permit the requested transition (e.g., trying to build without a plan), block and explain the prerequisite.
- **Artifact corruption** - If `plan.json` or `evidence.jsonl` is unparseable, flag it. Do not proceed. Ask the user whether to regenerate or restore from git.
- **Missing spec when plan exists** - Inconsistent state. Treat as corrupted - spec.md is the source of truth. Regenerate plan from spec or restore spec from git.

## Lifecycle states

Mission states enforce the development lane gates (CLI truth):

| State | Meaning | Permitted transitions |
|-------|---------|----------------------|
| `active` | Mission created / initial state | → `planned`, `blocked` |
| `planned` | Plan ready, awaiting build | → `in_progress`, `blocked` |
| `in_progress` | Implementation in progress | → `ready`, `blocked` |
| `ready` | Review passed, ready to ship | → `shipped`, `blocked` |
| `blocked` | Gated by critical finding | → `active`, `in_progress` |
| `shipped` | Merged and archived | Terminal |

Happy path: `active` → `planned` → `in_progress` → `ready` → `shipped`. `blocked` is reachable from any active (non-shipped) state.

Set state with `spacecraft set-state [mission-id] <new-state>` (mission-id optional when a mission is already resolved). The CLI enforces valid transitions.

## Rules

- **Must**: Resolve the active mission with `spacecraft resolve`, `status`, or `missions`; `.space/current` is fallback state, not sole authority. On conflict/ambiguity use `spacecraft use <selector>`.
- **Must**: Resolver priority is explicit selector or `SPACECRAFT_MISSION`, session binding, branch mission id, branch metadata, `.space/current`, then single active mission.
- **Must**: Strong signal conflicts or ambiguous active missions block mission writes until the user selects with `spacecraft use <number|id|title>` or an explicit selector.
- **Must**: Users may choose by list number, mission id, exact title, or unique title substring; do not expect the user to know a mission id.
- **Must**: New mission and evidence ids are compact sortable ids with no hyphen, such as `M07FYB5W5`; legacy `M-YYYYMMDD-HHmmss` ids remain valid.
- **Must**: Read the resolved mission's `mission.json`, `spec.md`, `questions.md`, `decisions.md`, `plan.json`, design artifacts, `evidence.jsonl`, and `review.json` when available.
- **Must**: sc-mission owns lifecycle but must route ambiguity to sc-clarify.
- **Must**: Enforce order: mission -> clarify -> spec -> visual design if needed -> plan -> build -> verify -> review -> ship. AFK build repeats RED → checkpoint → GREEN → evidence → checkpoint per acceptance, then combine/refactor checkpoint, until a gate blocks.
- **Must not**: Skip clarification when user intent, scope, acceptance criteria, or visual design direction is materially ambiguous.
- **Must**: If clear mutating work is requested and no suitable mission or branch exists, create the mission and non-main branch without another blocking question when policy permits it.
- **Must not**: Implement if spec or plan is missing.
- **Must**: Use sc-git for git safety, release branching, commits, rebasing, merge, version bump, changelog/spec notes, and tagging.
- **Must**: Mutating work must satisfy sc-git gates before implementation and before ship.
- **Must**: Treat stop-chat/close-session/end-session/new-session requests as session handoff unless release intent is explicit.
- **Must**: If "close session" is ambiguous and work appears ready, recommend ship; do not merge automatically.
- **Must**: Treat ship/release/merge/finish-mission/close-branch requests as release closeout prep. Block closeout when gates are incomplete.
- **Must**: After successful release closeout, archive shipped mission artifacts under `.space/archive/` unless the user asks to keep the full live mission folder.
- **Must**: After merge, run `spacecraft set-state shipped` to trigger archive and close GitHub issues. Capture evidence of the command output showing issue closing results.
- **Must**: Never claim GitHub issues are closed without running the actual command and capturing evidence. Verify by checking command output for "Issues: X closed" message.
- **Must**: Keep mission artifacts small and human-readable.
- **Must**: Prefer explicit evidence over narrative claims.
- **Must**: After creating or updating spec.md, plan.json, decisions.md, or questions.md, ctx_index them with source label `sc-memory/<mission-id>/<type>` (best-effort, non-blocking -- warn on failure). See sc-memory for label format and conventions.
- **Must**: End each session with a recommended next action and session advice: continue this chat for small adjacent steps, or start a new session when the phase changed, the thread is context-heavy, or mission artifacts are sufficient for handoff.

## Out of scope

This skill does NOT handle:

- Git operations, branching, or commits - use sc-git
- Planning tasks - use sc-planning
- Design or UI direction - use sc-design
- Evidence capture or verification - use sc-verification
- Clarification questions - use sc-clarify (sc-mission routes to it)
- Debugging or bug diagnosis - use sc-debug
- Code review - handled by reviewer subagent
- Knowledge capture and migration - use sc-learn
- Cross-mission memory, artifact indexing conventions, ctx_search/ctx_index wrapping - use sc-memory

## Output format

```
Mission artifacts follow the standard layout:
.space/missions/<id>/
  mission.json     # core metadata
  spec.md          # spec
  questions.md     # open/answered questions
  decisions.md     # confirmed choices and assumptions
  plan.json        # task plan
  evidence.jsonl   # evidence entries
  review.md        # review narrative
  review.json      # review findings and readiness
  design/          # design artifacts
  outputs/         # task outputs
```

## Checklist

Before claiming the mission lifecycle is handled:

- [ ] Development lane detected correctly from user intent
- [ ] Mission resolved with `spacecraft resolve` (on conflict/ambiguity: `spacecraft use <selector>`)
- [ ] All relevant artifacts read before any decision or mutation
- [ ] Ambiguity routed to sc-clarify when blocking (not bypassed)
- [ ] Lifecycle order enforced: clarify → spec → design → plan → build → verify → review → ship
- [ ] Lifecycle state is one of: active, planned, in_progress, ready, blocked, shipped
- [ ] Session handoff: state summary, blockers, dirty git, pickup command provided
- [ ] Release closeout: evidence checked, review gates passed, no-ff merge, tag, archive
- [ ] After merge: ran `set-state shipped`, captured evidence of issue closing, verified issues actually closed
- [ ] Never wrote product changes on `main`; always used a work branch

---

## References

- `spacecraft --help` - spacecraft CLI reference
- `spacecraft resolve --help` - resolver subcommand
- `spacecraft use --help` - mission selection
- `spacecraft evidence --help` - evidence subcommand
- `spacecraft closeout-check --help` - closeout verification
