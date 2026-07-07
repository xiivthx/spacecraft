# Spacecraft Project Specification

## Purpose

Spacecraft is a local-first mission-control harness for OpenCode-driven software development. It coordinates specification, planning, implementation, verification, review, release closeout, and session handoff through file-backed mission artifacts and slash command prompts.

## Core Contract

- Mission state is stored under `.space/missions/<mission-id>/`.
- Active mission resolution must use `node scripts/spacecraft.mjs resolve [selector] [--json]`, `status`, or `missions`; `.space/current` is fallback state, not sole authority.
- Product implementation must not begin before `spec.md` and `plan.json` exist for the resolved mission.
- Work is not done, passed, verified, ready, shipped, or releasable without evidence in `evidence.jsonl`.
- Git is the default rollback boundary for implementation work.
- Product changes must not be written directly on `main`.
- Release closeout is distinct from session handoff.

## Mission Ids

New mission ids and evidence ids use compact sortable ids with no hyphen:

```text
M07FYB5W5
E07FYB5W5
```

The prefix identifies the artifact kind. The 8-character suffix is fixed-width uppercase base36 milliseconds since `2026-01-01T00:00:00.000Z`, preserving lexicographic time order for about 89 years. Legacy mission ids such as `M-20260707-141230` remain supported for existing missions and branches.

## Workflow

The normal lifecycle is:

```text
/sc-start -> /sc-clarify -> /sc-design when needed -> /sc-plan -> /sc-git -> /sc-work -> /sc-verify -> /sc-review -> /sc-ship
```

`/sc-flow` reduces repetitive HIL during implementation by continuing this safe loop inside the same chat:

```text
/sc-work Txx -> /sc-verify Txx -> checkpoint commit -> next task
```

The workflow runner stops on resolver conflicts, open blocking clarification, missing mission artifacts, main-branch write risk, unsafe dirty files, missing design direction, failed verification, failed validation, critical review findings, release actions, or context that is too heavy for safe continuation.

## Verification

Verification evidence is captured with:

```sh
node scripts/spacecraft.mjs evidence "<label>" -- <command>
```

Each evidence entry records the evidence id, label, command, exit code, output file paths, and creation timestamp. Failures are evidence too.

## Git And Release

Spacecraft uses release branching:

- one non-main branch per feature, fix, issue, or tightly scoped change
- branch names follow `<type>/<issue-or-mission>-<slug>`
- checkpoint commits are allowed only on valid non-main work branches
- final branch history should be 1 to 3 logical Conventional Commits and must not exceed 5 without justification
- work branches are rebased on latest `main` before merge
- releases merge into `main` with `git merge --no-ff <branch>`
- version tags are annotated and created after the merge commit
- pushes require explicit user request

## Archive

After successful release closeout, shipped missions are compacted from `.space/missions/` into `.space/archive/` unless the user asks to keep the full live mission folder. Archive compaction requires completed plan tasks, evidence, ready review, no critical findings, and recorded release-readiness gates. Archives keep durable summary artifacts and omit bulky command output files.
