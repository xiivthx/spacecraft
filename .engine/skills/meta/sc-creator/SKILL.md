---
name: sc-creator
description: >
  Create spacecraft skills, agents, and commands from templates. Activate on "create skill", "new agent",
  "add command", or when building sc-* artifacts from templates or research.
license: MIT
compatibility: opencode
metadata:
  version: 1.0
  audience: xiivthx
---

# sc-creator

Create spacecraft artifacts (skills, agents, commands) from templates. End-to-end workflow: gather content, scaffold from template, wire to system, polish for conventions, register in docs, verify consistency.

## When to use

Activate when the user asks to:

- **"Create skill X" / "new skill" / "add skill sc-*"** — skill creation
- **"Create agent X" / "new agent" / "add agent sc-*"** — agent creation
- **"Create command X" / "new command" / "add command sc-*"** — command creation
- **"Build from this datasource"** — content from files, research, or internal knowledge
- When a new sc-* capability needs to be codified into the system

## Workflow

Three creation modes. Execute in order unless the user specifies otherwise.

### Mode selection

1. **Identify artifact type** — skill, agent, or command
2. **Read template** — `templates/skill.md`, `templates/agent.md`, or `templates/command.md`
3. **Survey conventions** — Read 1–2 existing artifacts of the same type
4. **Map content** — Decide structure based on artifact type

### Phase 1: Gather

1. **Identify datasource** — User provides content (`.space/temp/` files, a URL, a description) or research internally. If domain knowledge is missing, run `spacecraft research "<topic>"` before proceeding.

2. **Read template** — Based on artifact type:
   - **Skill**: `templates/skill.md` — canonical skill structure
   - **Agent**: `templates/agent.md` — canonical agent structure
   - **Command**: `templates/command.md` — canonical command structure

3. **Survey conventions** — Read 1–2 existing artifacts of the same type. Absorb: naming patterns, section tone, cross-reference style, how rules are phrased.

4. **Map content** — Decide structure:
   - **Skill**: SKILL.md (operational) + references/ (detail)
   - **Agent**: Single markdown file with role, constraints, tools
   - **Command**: Single markdown file with frontmatter, workflow, gates

### Phase 2: Create

#### For skills:

1. **Create directory** — `.engine/skills/<category>/sc-<name>/`. Add `references/` subdirectory if needed.

2. **Write SKILL.md** from template — every section required:
   - `description` — under 200 chars, includes 2–3 trigger phrases
   - `## When to use` — concrete trigger patterns
   - `## Workflow` — exact sequence with commands in backticks
   - `## Rules` — Must / Must not / Prefer. Every Must verifiable.
   - `## Out of scope` — what this skill does NOT handle
   - `## Output format` — exact expected output shape
   - `## Checklist` — actionable, verifiable items
   - `## References` — list reference files with one-line descriptions

3. **Write references** — One file per topic. Each reference:
   - Starts with `> Consult when:` one-liner
   - Is 30–150 lines
   - Has `## Spacecraft integration` section linking to missions, plan.json, evidence

#### For agents:

1. **Create file** — `.engine/agents/sc-<name>.md`

2. **Write agent** from template — every section required:
   - `name` — agent identifier
   - `description` — under 200 chars, role summary
   - `## Role` — what this agent does
   - `## Constraints` — what this agent must/must not do
   - `## Tools` — which tools this agent can use
   - `## Output format` — expected response shape
   - `## Checklist` — verification items

#### For commands:

1. **Create file** — `.engine/commands/sc-<name>.md`

2. **Write command** from template — every section required:
   - Frontmatter: `name`, `description`, `subtask` (if applicable)
   - `## Pre-flight Checks` — validation before execution
   - `## Workflow` — exact execution sequence
   - `## Hard Stop Gates` — conditions that block execution
   - `## Error Handling` — how to handle failures
   - `## Output format` — expected output shape

### Phase 3: Wire

1. **Map to system** — Based on artifact type:
   - **Skill**: Map to agents by role. Add `"sc-<name>": allow` to agent `skill.permission` blocks.
   - **Agent**: Add to `opencode.json` agent registry. Wire to skills via `skill.permission`.
   - **Command**: Add to `opencode.json` command registry. Wire to agents via `agent` field.

2. **Update registries** — Add to appropriate config files in alphabetical order.

### Phase 4: Polish

1. **Spacecraft-native rewrite** — Reference spacecraft concepts (missions, `plan.json`, `evidence.jsonl`, `decisions.md`, `spec.md`). Operational tone, not educational.

2. **Remove agent names** — Skills and commands are passive resources. Never mention Commander, sc-coder, sc-tester, sc-reviewer, sc-planner, or sc-designer in content.

3. **Remove mutual cross-references** — No loops between artifacts. Each is self-contained.

4. **Consolidate content** — Single source of truth per domain. No duplicated knowledge.

5. **Verify limits**:
   - **Skill**: `description` < 200 chars. SKILL.md 80–150 lines (meta-skills may exceed). References 30–150 lines each.
   - **Agent**: `description` < 200 chars. File 50–200 lines.
   - **Command**: `description` < 200 chars. File 30–150 lines.

### Phase 5: Register

1. **Update SPACECRAFT.md** — Based on artifact type:
   - **Skill**: §Slash commands, §Routing table, §Subagent table, §Skill references table
   - **Agent**: §Subagent table, §Agent references table
   - **Command**: §Slash commands, §Routing table

2. **Update AGENTS.md** — Add to available skills/agents/commands tables.

3. **Go scripts** — No changes needed. Artifacts are OpenCode agent-layer; Go CLI handles missions and workflow.

### Phase 6: Verify

1. **Cross-check system consistency** — Audit artifacts on disk, agent permissions, command registries, and SPACECRAFT.md tables against each other. Confirm: zero phantom refs, zero artifacts missing from docs.

2. **Fix common issues** — Stale references, missing SPACECRAFT.md entries, description overflow, missing registry entries.

3. **Self-review** — Run new artifact's own checklist. Verify all required sections, no agent names, no cross-reference loops.

4. **Commit** — Conventional Commits, non-main branch only. Body lists created files, wiring, registrations.

## Rules

- **Must**: Start from appropriate template (`templates/skill.md`, `templates/agent.md`, or `templates/command.md`). Survey 1–2 existing artifacts of the same type.
- **Must**: All required sections present. Description under 200 chars with trigger phrases.
- **Must**: Rewrite datasource content for spacecraft context. Never copy-paste verbatim.
- **Must**: Delegate detail to references (skills only). Main file is operational; references are deep dives.
- **Must**: Wire to system, register in SPACECRAFT.md and AGENTS.md before claiming done.
- **Must**: Cross-check all system layers (artifacts, agents, commands, docs) before commit.
- **Must not**: Include agent names in skill/command content. Skills and commands are passive resources.
- **Must not**: Create mutual cross-references between artifacts. No loops — single direction only.
- **Must not**: Write on `main`. Create a work branch before any mutating work.

## Out of scope

- Mission creation — use sc-mission, the start command
- Code implementation — use the build command
- Testing — use sc-tdd
- Git operations — use sc-git
- UI design — use sc-design

## Output format

```
Mode: [skill | agent | command]

Phase 1: Gather
  Datasource: [path / url / research topic]
  Template: templates/<type>.md ✓
  Conventions surveyed: [artifact names]

Phase 2: Create
  Directory/File: .engine/<type>/<path>
  Main file: [N] lines, description [N] chars
  References: [count] files ([total] lines) — skills only

Phase 3: Wire
  System: [agents / commands / registries]

Phase 4: Polish
  Agent names removed ✓ | Loops checked ✓ | Consolidated ✓ | Limits ok ✓

Phase 5: Register
  SPACECRAFT.md: [tables updated]
  AGENTS.md: [tables updated]

Phase 6: Verify
  Cross-check: [N] phantom refs, [N] missing docs, [N] overflow descs
  Commit: feat: add sc-<name>
```

## Checklist

Before claiming an artifact is created:

- [ ] Template consulted (`templates/skill.md`, `templates/agent.md`, or `templates/command.md`)
- [ ] 1–2 existing artifacts of same type surveyed for conventions
- [ ] Directory/file created at correct path
- [ ] All required sections present
- [ ] `description` under 200 chars with trigger phrases
- [ ] Content rewritten for spacecraft context (not copy-paste)
- [ ] No agent names in skill/command content
- [ ] No mutual cross-references with other artifacts (loops)
- [ ] System wiring complete (agent permissions, command registries)
- [ ] SPACECRAFT.md updated: relevant tables
- [ ] AGENTS.md updated: relevant tables
- [ ] Cross-check: zero phantom refs, zero missing doc entries
- [ ] Committed with Conventional Commit on a non-main branch

---

## References

- `templates/skill.md` — canonical skill template with field annotations
- `templates/agent.md` — canonical agent template with field annotations
- `templates/command.md` — canonical command template with field annotations
