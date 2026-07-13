---
name: sc-creator
description: >
  Create new Spacecraft skills from datasources. Activate on "create skill", "new skill", "add skill",
  or when building a sc-* skill from a template or research.
license: MIT
compatibility: opencode
metadata:
  version: 1.0
  audience: xiivthx
---

# sc-creator

End-to-end skill creation: gather content, scaffold from template, wire to agents and commands, polish for spacecraft conventions, register in docs, and verify system consistency. Output is a fully integrated `sc-*` skill ready for use.

## When to use

Activate when the user asks to:

- **"Create skill X" / "new skill" / "add skill sc-*"** — explicit skill creation
- **"Build a skill from this datasource"** — content from files, research, or internal knowledge
- When a new `sc-*` capability needs to be codified into the system

## Workflow

Six phases — execute in order unless the user specifies otherwise.

### Phase 1: Gather

1. **Identify datasource** — User provides content (`.space/temp/` files, a URL, a description) or Commander researches internally. If domain knowledge is missing, run `spacecraft research "<topic>"` before proceeding.

2. **Read template** — `references/template.md`. This is the canonical structure. Always start here — never write a skill from memory.

3. **Survey conventions** — Read 1–2 existing skills (e.g., `sc-git/SKILL.md`, `sc-verification/SKILL.md`). Absorb: naming patterns, section tone, cross-reference style, how rules are phrased.

4. **Map content** — Decide what goes in SKILL.md vs references. SKILL.md is operational (trigger, workflow, rules, checklist). References are detail (examples, deep dives, tables, code samples).

### Phase 2: Create

1. **Create directory** — `.engine/skills/sc-<name>/`. Add `references/` subdirectory if the skill has reference content.

2. **Write SKILL.md** from the template — every section is required:
   - `description` — under 200 chars, includes 2–3 trigger phrases. This is the only thing the agent sees before loading.
   - `## When to use` — concrete trigger patterns, not abstract categories (e.g., "when user says 'commit'" not "when git operations occur")
   - `## Workflow` — exact sequence with commands in backticks. Use sub-bullets for variations and edge cases.
   - `## Rules` — Must / Must not / Prefer. Non-negotiable items first. Every Must must be verifiable.
   - `## Out of scope` — what this skill does NOT handle. Keeps boundaries clean and prevents scope creep.
   - `## Output format` — exact expected output shape. Show the format, not just describe it.
   - `## Checklist` — actionable items. Every item must be verifiable (can a reviewer confirm it?).
   - `## References` — list reference files with one-line descriptions. Load-on-demand only.

3. **Write references** — One file per topic. Each reference:
   - Starts with `> Consult when:` one-liner so the agent knows when to load it
   - Is 30–150 lines — concise enough to load quickly, long enough to be useful
   - Has a `## Spacecraft integration` section at the end linking to missions, plan.json, evidence
   - Uses TypeScript examples consistent with the codebase style

### Phase 3: Wire

1. **Map to agents** — Relevant agents by role: sc-coder (implementation/code quality), sc-tester (testing/verification), sc-planner (architecture/planning), sc-reviewer (review/gates), sc-designer (UI/design). sc-commander has `sc-*` wildcard — never needs explicit wiring.
2. **Update agent files** — Add `"sc-<name>": allow` to `skill.permission` block, alphabetical order.

### Phase 4: Polish

1. **Spacecraft-native rewrite** — Reference spacecraft concepts (missions, `plan.json`, `evidence.jsonl`, `decisions.md`, `spec.md`). Operational tone, not educational.
2. **Remove agent names** — Skills are passive resources. Never mention Commander, sc-coder, sc-tester, sc-reviewer, sc-planner, or sc-designer in skill content.
3. **Remove mutual cross-references** — No loops between skills. Each skill is self-contained.
4. **Consolidate content** — Single source of truth per domain. No duplicated knowledge across skills.
5. **Verify limits** — `description` < 200 chars. SKILL.md 80–150 lines (meta-skills may exceed). References 30–150 lines each.

### Phase 5: Register

1. **Update command files** — Add `sc-<name>` to `Use:` lines in relevant commands: build (impl/test/quality), review (review/gates), plan (planning/architecture), design (UI/design).
2. **Update SPACECRAFT.md** — Four spots: §Slash commands, §Routing table, §Subagent table, §Skill references table. Add row with skill name, file path, used by.
3. **Go scripts** — No changes needed. Skills are OpenCode agent-layer; Go CLI handles missions and workflow.

### Phase 6: Verify

1. **Cross-check system consistency** — Audit skills on disk, agent permissions, command `Use:` lines, and SPACECRAFT.md tables against each other. Confirm: zero phantom refs, zero skills missing from docs.
2. **Fix common issues** — Stale references, missing SPACECRAFT.md entries, description overflow, missing quick command from command registry.
3. **Self-review** — Run new skill's own checklist. Verify 7 template sections, no agent names, no cross-reference loops.
4. **Commit** — Conventional Commits, non-main branch only. Body lists created files, agents wired, commands registered.

## Rules

- **Must**: Start from `references/template.md`. Survey 1–2 existing skills for conventions before writing.
- **Must**: All 7 template sections present. Description under 200 chars with trigger phrases.
- **Must**: Rewrite datasource content for spacecraft context. Never copy-paste verbatim.
- **Must**: Delegate detail to references. SKILL.md is operational; references are deep dives.
- **Must**: Wire to agents, register in commands and SPACECRAFT.md before claiming done.
- **Must**: Cross-check all four system layers (skills, agents, commands, docs) before commit.
- **Must not**: Include agent names in skill content. Skills are passive resources.
- **Must not**: Create mutual cross-references between skills. No loops — single direction only.
- **Must not**: Write on `main`. Create a work branch before any mutating work.

## Out of scope

- Mission creation — use sc-mission, the start command
- Code implementation — use the build command
- Testing — use sc-tdd
- Git operations — use sc-git
- UI design — use sc-design

## Output format

```
Phase 1: Gather
  Datasource: [path / url / research topic]
  Template: references/template.md ✓
  Conventions surveyed: [skill names]

Phase 2: Create
  Directory: .engine/skills/sc-<name>/
  SKILL.md: [N] lines, description [N] chars
  References: [count] files ([total] lines)

Phase 3: Wire
  Agents: [list with ✓ per agent]

Phase 4: Polish
  Agent names removed ✓ | Loops checked ✓ | Consolidated ✓ | Limits ok ✓

Phase 5: Register
  Commands: [list]
  SPACECRAFT.md: routing ✓ | subagent ✓ | skill refs ✓

Phase 6: Verify
  Cross-check: [N] phantom refs, [N] missing docs, [N] overflow descs
  Commit: feat: add sc-<name>
```

## Checklist

Before claiming a skill is created:

- [ ] Template (`references/template.md`) consulted
- [ ] 1–2 existing skills surveyed for conventions
- [ ] Directory created at `.engine/skills/sc-<name>/`
- [ ] SKILL.md has all 7 template sections
- [ ] `description` under 200 chars with trigger phrases
- [ ] Content rewritten for spacecraft context (not copy-paste)
- [ ] No agent names anywhere in skill content
- [ ] No mutual cross-references with other skills (loops)
- [ ] Agent files updated for all relevant agents (alphabetical order)
- [ ] Command `Use:` lines updated for all relevant commands
- [ ] SPACECRAFT.md updated: routing table, subagent table, skill refs table
- [ ] Cross-check: zero phantom refs, zero missing doc entries
- [ ] Committed with Conventional Commit on a non-main branch

---

## References

- `references/template.md` — canonical skill template with field annotations and section requirements
- `docs/SPACECRAFT.md` — master registry: routing table, subagent table, skill references table
