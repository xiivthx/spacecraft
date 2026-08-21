---
name: sc-writer
description: "On-demand prose craft for docs, prompts, and narrative: prompt-refine diagnose→rewrite, rhythm rewrite, and context harvest before draft. Activate when refining agent/skill/rule prompts, rewriting user-facing prose, drafting narrative, or Task(sc-writer)."
---

# sc-writer

## Goal

On-demand prose craft for docs, prompts, and user-facing narrative - without changing gates, policy behavior, or product code.

## When to use

- Diagnosing or rewriting agent/skill/rule prompt fidelity → `references/prompt-refine.md`
- Rewriting narrative or user-facing prose for engagement → `references/prose-rhythm.md`
- High-stakes narrative draft when context is thin → `references/narrative-context.md` before drafting
- `Task(sc-writer)` for docs, prompt text, messages, handoff copy

## Workflow

1. **Classify** - prompt fidelity diagnose/rewrite vs narrative/user-facing prose vs machine-checkable gates (Verify bars, JSON, CLI flags, Never/Always lists, code).
2. **Prompt refine** - if agent/skill/rule prompt fidelity → follow `prompt-refine.md`.
3. **Narrative** - follow `prose-rhythm.md` for rewrite; if context is thin for high-stakes draft, run `narrative-context.md` first.
4. **Gates / Verify / schemas** - stay short, precise, compressed; do not apply rhythm craft.
5. **Edit** - wording/structure only; policy and gate behavior stay as-is.

## Good / Bad

- Good: prompt-refine diagnose→rewrite for agent/skill/rule fidelity; rhythm mix for narrative; context harvest when thin; short/precise for gates; US English; ASCII hyphen-minus only
- Bad: lyrical rhythm on Verify/gates/JSON; questionnaire dumps; expertise cosplay; behavior or policy changes via wording

## Verify

Commander re-runs task `verify` (e.g. `rg` for expected phrase) or reads diff for wording-only changes. On-demand craft runs emit their reference Verify bars (`## Prompt refine`, `### essence ###`, `## Narrative context`, or skip lines) - none block discuss clear or ship.

## Must / Must not

- **Must**: Match existing file structure and section names when editing
- **Must**: Use reference procedures for on-demand craft jobs
- **Must not**: Edit product code or tests
- **Must not**: Change what a gate, rule, or check *does* while editing wording
- **Must not**: Add always-on communication compression or chat rules
- **Must not**: Unicode em dash

## Out of scope

- Product code and tests
- Gates, policy, and runtime behavior changes
- Always-on chat compression or communication rules
- Visual UI critique (sc-designer)

## Related

- Agent: `.cursor/agents/sc-writer.md`
- Prompt refine: `references/prompt-refine.md`
- Rhythm rewrite: `references/prose-rhythm.md`
- Narrative context: `references/narrative-context.md`
