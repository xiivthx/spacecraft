---
name: sc-creator
description: "Creates Cursor-native spacecraft skills and agents when users request a reusable capability or specialized subagent."
---

# sc-creator

Create spacecraft skills and agents from Cursor-native conventions.

## When to use

Activate when the user asks to:

- **"Create skill X" / "new skill" / "add skill sc-*"** - skill creation
- **"Create agent X" / "new agent" / "add agent sc-*"** - agent creation
- **"Build from this datasource"** - content from files, research, or internal knowledge
- When a new sc-* capability needs to be codified into the system

## Workflow

Execute the phases in order unless the user specifies otherwise.

### Mode selection

1. **Identify artifact type** - skill or agent
2. **Read template or examples** - use `templates/skill.md` for skills or existing `.cursor/agents/*.md` files for agents
3. **Survey conventions** - Read 1–2 existing artifacts of the same type
4. **Map content** - Decide structure based on artifact type

### Phase 1: Gather

1. **Identify datasource** - User provides content (`.space/temp/` files, a URL, a description) or research internally. If domain knowledge is missing, use sc-search (WebSearch/WebFetch) for `"<topic>"` before proceeding.

2. **Read template** - Based on artifact type:
   - **Skill**: `templates/skill.md` - canonical skill structure
   - **Agent**: existing `.cursor/agents/*.md` files - canonical agent examples

3. **Survey conventions** - Read 1–2 existing artifacts of the same type. Absorb: naming patterns, section tone, cross-reference style, how rules are phrased.

4. **Map content** - Decide structure:
   - **Skill**: SKILL.md (operational) + references/ (detail)
   - **Agent**: Single markdown file with role, constraints, tools

### Phase 2: Create

#### For skills:

1. **Create directory** - `.cursor/skills/sc-<name>/`. Add `references/` subdirectory if needed.

2. **Write SKILL.md** from template - every section required:
   - `description` - under 200 chars, includes 2–3 trigger phrases
   - `## When to use` - concrete trigger patterns
   - `## Workflow` - exact sequence with commands in backticks
   - `## Rules` - Must / Must not / Prefer. Every Must verifiable.
   - `## Out of scope` - what this skill does NOT handle
   - `## Output format` - exact expected output shape
   - `## Checklist` - actionable, verifiable items
   - `## References` - list reference files with one-line descriptions

3. **Write references** - One file per topic. Each reference:
   - Starts with `> Consult when:` one-liner
   - Is 30–150 lines
   - Has `## Spacecraft integration` section linking to missions, plan.json, evidence

#### For agents:

1. **Create file** - `.cursor/agents/sc-<name>.md`

2. **Write agent** from template - every section required:
   - `name` - agent identifier
   - `description` - under 200 chars, role summary
   - `## Role` - what this agent does
   - `## Constraints` - what this agent must/must not do
   - `## Tools` - which tools this agent can use
   - `## Output format` - expected response shape
   - `## Checklist` - verification items


### Phase 3: Wire

Confirm the skill or agent is discoverable from its Cursor-native location.

### Phase 4: Polish

1. **Spacecraft-native rewrite** - Reference spacecraft concepts (missions, `plan.json`, `evidence.jsonl`, `decisions.md`, `spec.md`). Operational tone, not educational.

2. **Keep routing explicit** - Mention a subagent only when the skill delegates work to it.

3. **Remove mutual cross-references** - No loops between artifacts. Each is self-contained.

4. **Consolidate content** - Single source of truth per domain. No duplicated knowledge.

5. **Verify limits**:
   - **Skill**: `description` < 200 chars. SKILL.md 80–150 lines (meta-skills may exceed). References 30–150 lines each.
   - **Agent**: `description` < 200 chars. File 50–200 lines.

### Phase 5: Register

1. **Update Cursor rules only when needed** - Keep discovery and routing references accurate.

2. **Check agent definitions** - Confirm delegated subagents exist under `.cursor/agents/`.

3. **Go scripts** - No changes needed. Artifacts are Cursor agent-layer; Go CLI handles missions and workflow.

### Phase 6: Verify


2. **Fix common issues** - Stale references, missing Cursor rules entries, description overflow, missing registry entries.

3. **Self-review** - Run new artifact's own checklist. Verify all required sections, no agent names, no cross-reference loops.

4. **Commit** - Conventional Commits, non-main branch only. Body lists created files, wiring, registrations.

## Rules

- **Must**: Start from `templates/skill.md` or existing `.cursor/agents/*.md` examples. Survey 1–2 existing artifacts of the same type.
- **Must**: All required sections present. Description under 200 chars with trigger phrases.
- **Must**: Rewrite datasource content for spacecraft context. Never copy-paste verbatim.
- **Must**: Delegate detail to references (skills only). Main file is operational; references are deep dives.
- **Must**: Wire to system, register in the relevant Cursor rules and agent files before claiming done.
- **Must**: Cross-check affected artifacts, agents, rules, and docs before commit.
- **Must not**: Create `.cursor/commands/`; explicit workflows are Cursor skills with `disable-model-invocation: true`.
- **Must not**: Create mutual cross-references between artifacts. No loops - single direction only.
- **Must not**: Write on `main`. Create a work branch before any mutating work.

## Out of scope

- Mission creation - use sc-mission, the start command
- Code implementation - use the build command
- Testing - use sc-tdd
- Git operations - use sc-git
- UI design - use sc-design

## Output format

```
Mode: [skill | agent]

Phase 1: Gather
  Datasource: [path / url / research topic]
  Template: templates/<type>.md ✓
  Conventions surveyed: [artifact names]

Phase 2: Create
  Directory/File: .cursor/<type>/<path>
  Main file: [N] lines, description [N] chars
  References: [count] files ([total] lines) - skills only

Phase 3: Wire
  System: [skills / agents / rules]

Phase 4: Polish
  Agent names removed ✓ | Loops checked ✓ | Consolidated ✓ | Limits ok ✓

Phase 5: Register
  Cursor rules: [tables updated]
  Cursor agent definitions: [tables updated]

Phase 6: Verify
  Cross-check: [N] phantom refs, [N] missing docs, [N] overflow descs
  Commit: feat: add sc-<name>
```

## Checklist

Before claiming an artifact is created:

- [ ] Template consulted (`templates/skill.md`, `.cursor/agents/ examples`, or `.cursor/skills/ workflow examples`)
- [ ] 1–2 existing artifacts of same type surveyed for conventions
- [ ] Directory/file created at correct path
- [ ] All required sections present
- [ ] `description` under 200 chars with trigger phrases
- [ ] Content rewritten for spacecraft context (not copy-paste)
- [ ] Delegated agent names and paths are valid
- [ ] No mutual cross-references with other artifacts (loops)
- [ ] Cursor rules updated: relevant tables
- [ ] Cursor agent definitions updated: relevant tables
- [ ] Cross-check: zero phantom refs, zero missing doc entries
- [ ] Committed with Conventional Commit on a non-main branch

---

## References

- `templates/skill.md` - canonical skill template with field annotations
- `.cursor/agents/*.md` - canonical agent examples
- `.cursor/skills/*/SKILL.md` - canonical workflow and domain skill examples
