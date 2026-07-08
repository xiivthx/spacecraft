---
description: Review resolved mission diff and evidence
agent: sc-commander
---
Use sc-mission, sc-verification, and sc-git.
Resolve the mission. Block if unsafe.

## Pre-flight checks

Read the resolved mission's spec.md, plan.json, evidence.jsonl, review.json, and git diff when git is available.

Run:
```
scripts/spacecraft git-info
```

## Workflow

### 1. Check git readiness

Use sc-git for git readiness checks: branch hygiene, commit style, rebase readiness, merge plan, version/changelog/tag plan, and branch cleanup policy.

### 2. Check release readiness

If the diff is intended to ship, record release readiness in review.json:
- releaseReadiness.version
- releaseReadiness.changelog
- releaseReadiness.specNote
- releaseReadiness.tagPlan
- releaseReadiness.postRebaseVerification

Use structured objects with `status`, plus `rationale` for any deferred gate. Do not use string or boolean releaseReadiness values.

### 3. UI design review (if UI files changed)

If UI files changed, invoke sc-designer as a read-only sidecar for design-risk triage. A user invocation of /sc-review is sufficient permission for this; do not ask for separate subagent permission.

Ask the subagent to review:
- hierarchy, layout, typography, spacing, color use, interaction states
- accessibility, responsiveness, anti-slop checklist
- Feynman clarity: plain-language explanation, labeled visuals, obvious gain/tradeoff, no unnecessary jargon
- visual economy: HTML was used only when visual comparison materially helped the decision
- reference fit: implementation learns from selected references without cloning them
- art direction match if UI direction was recorded
- Thai-first copy for multilingual missions
- option set diversity: flag same-y choices differing only by palette/labels

When visual review needs a browser and no app-specific dev server is running:
```
node .opencode/skills/sc-design/scripts/serve-html.mjs <artifact-or-dir> --open
```

After the subagent responds, record findings in review.md. If review.json exists, add design findings with severity: critical, important, minor. Critical design findings block /sc-ship.

### 4. Run code review

Invoke sc-reviewer as a read-only subagent. A user invocation of /sc-review is explicit permission to use the read-only sc-reviewer subagent; do not ask for separate subagent permission. The reviewer must not edit files.

### 5. Record findings

After the subagent returns findings, record the review in:
- review.md
- review.json

Set state to reviewing while working. If review status is ready, set state to ready. If review status is blocked, set state to blocked.

## Error handling

- Do not implement fixes in the same command unless the user explicitly asks.
- Critical design findings block shipping the same way critical code findings do.

End with next action and session advice. Suggest /sc-ship only when ready.
