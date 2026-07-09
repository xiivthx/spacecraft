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

2. **Read template** — `docs/templates/skill.md`. This is the canonical structure. Always start here — never write a skill from memory.

3. **Survey conventions** — Read 1–2 existing skills (e.g., `sc-git/SKILL.md`, `sc-verification/SKILL.md`). Absorb: naming patterns, section tone, cross-reference style, how rules are phrased.

4. **Map content** — Decide what goes in SKILL.md vs references. SKILL.md is operational (trigger, workflow, rules, checklist). References are detail (examples, deep dives, tables, code samples).

### Phase 2: Create

1. **Create directory** — `.opencode/skills/sc-<name>/`. Add `references/` subdirectory if the skill has reference content.

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

1. **Map to agents** — Determine which agents need this skill based on their role:
   - `sc-coder` — for implementation and code quality skills
   - `sc-tester` — for testing and verification skills
   - `sc-planner` — for architecture and planning skills
   - `sc-reviewer` — for review and quality gate skills
   - `sc-designer` — for UI and design skills
   - `sc-commander` — already has `sc-*` wildcard, never needs explicit wiring

2. **Update agent files** — For each relevant agent, add `"sc-<name>": allow` to the `skill.permission` block in `.opencode/agents/<agent>.md`. Insert in alphabetical order among existing entries.

### Phase 4: Polish

1. **Spacecraft-native rewrite** — Content must reference spacecraft concepts: missions, `plan.json` tasks, `evidence.jsonl`, `decisions.md`, `spec.md`. Remove generic textbook tone. Every sentence should be operational, not educational.

2. **Remove agent names from skill content** — Skills are passive resources. They describe what to do, never which agent does it. Direction is strictly Commander → skills. Never mention Commander, sc-coder, sc-tester, sc-reviewer, sc-planner, or sc-designer in skill content.

3. **Remove mutual cross-references** — If skill A points to skill B (e.g., Out of scope: "use sc-B"), then skill B must NOT point back to skill A. This creates a loop. Each skill is self-contained; the Commander decides which skills to load together. Concrete example: if sc-solid says "use sc-tdd" and sc-tdd says "use sc-solid", that's a loop — remove both.

4. **Consolidate content** — Single source of truth per domain. If two skills overlap (e.g., both have testing strategy content), move all content into one skill and remove from the other. No duplicated knowledge.

5. **Verify limits** — `description` under 200 chars. SKILL.md typically 80–150 lines (longer for meta-skills is acceptable). References 30–150 lines each. If a reference would exceed 150 lines, split into two reference files.

### Phase 5: Register

1. **Update command files** — Add `sc-<name>` to the `Use:` line in `.opencode/commands/` for each relevant command:
   - `/sc-build` — for implementation, testing, and code quality skills
   - `/sc-review` — for review and quality gate skills
   - `/sc-plan` — for planning and architecture skills
   - `/sc-design` — for UI and design skills
   - Other commands only if the skill is directly relevant to that command's phase

2. **Update SPACECRAFT.md** — Four spots to check and update:
   - **§Slash commands** (line ~19) — add command name if creating a new command
   - **§Routing table** (lines ~97-107) — add skill to the command's skill and permission columns
   - **§Subagent table** (lines ~112-119) — verify each subagent's skill list matches their agent file
   - **§Skill references table** (lines ~138-152) — add new row: skill name, file path, used by

3. **Go scripts** — No changes needed. Skills are an OpenCode agent-layer concept. The Go CLI handles missions and workflow states, not skill loading.

### Phase 6: Verify

1. **Cross-check system consistency** — Audit these four sources against each other:
   - Skills on disk (`.opencode/skills/sc-*/SKILL.md` frontmatter `name` fields)
   - Agent `skill.permission` blocks (`.opencode/agents/*.md`)
   - Command `Use:` lines (`.opencode/commands/*.md`)
   - SPACECRAFT.md tables (§Routing, §Subagent, §Skill references)

   Confirm: zero skills referenced that don't exist on disk. Zero skills on disk missing from docs.

2. **Fix inconsistencies** — Common issues to check:
   - Stale references to skills that don't exist on disk (phantoms)
   - Skills on disk not listed in SPACECRAFT.md skill references table
   - Agent files referencing skills not on disk
   - Command `Use:` lines referencing skills not on disk
   - `description` fields exceeding 200 chars
   - Missing `/sc-quick` from slash commands or routing table

3. **Self-review** — Run the new skill's own checklist against itself. Verify all 7 template sections present. Check for accidental agent names or cross-references.

4. **Commit** — Stage only new and intentionally modified files (never bulk-add unrelated dirty files). Use Conventional Commits:
   ```
   feat: add sc-<name> skill
   ```
   Body: list what was created (SKILL.md lines, reference count), which agents were wired, which commands registered.

## Rules

- **Must**: Start from `docs/templates/skill.md`. Survey 1–2 existing skills for conventions before writing.
- **Must**: All 7 template sections present. Description under 200 chars with trigger phrases.
- **Must**: Rewrite datasource content for spacecraft context. Never copy-paste verbatim.
- **Must**: Delegate detail to references. SKILL.md is operational; references are deep dives.
- **Must**: Wire to agents, register in commands and SPACECRAFT.md before claiming done.
- **Must**: Cross-check all four system layers (skills, agents, commands, docs) before commit.
- **Must not**: Include agent names in skill content. Skills are passive resources.
- **Must not**: Create mutual cross-references between skills. No loops — single direction only.
- **Must not**: Write on `main`. Create a work branch before any mutating work.

## Out of scope

- Mission creation — use sc-mission, /sc-start
- Code implementation — use /sc-build
- Testing — use sc-tdd
- Git operations — use sc-git
- UI design — use sc-design

## Output format

```
Phase 1: Gather
  Datasource: [path / url / research topic]
  Template: docs/templates/skill.md ✓
  Conventions surveyed: [skill names]

Phase 2: Create
  Directory: .opencode/skills/sc-<name>/
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

- [ ] Template (`docs/templates/skill.md`) consulted
- [ ] 1–2 existing skills surveyed for conventions
- [ ] Directory created at `.opencode/skills/sc-<name>/`
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

- `docs/templates/skill.md` — canonical skill template with field annotations
- `.opencode/skills/sc-git/SKILL.md` — reference: complex rules, workflow, checklist
- `.opencode/skills/sc-verification/SKILL.md` — reference: concise, evidence-driven
- `.opencode/agents/` — agent config: `skill.permission` blocks
- `.opencode/commands/` — command config: `Use:` frontmatter
- `docs/SPACECRAFT.md` — master registry: routing, subagent, skill reference tables
