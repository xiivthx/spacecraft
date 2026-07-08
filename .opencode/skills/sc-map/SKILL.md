---
name: sc-map
description: Survey project structure before planning — map relevant files, dependencies, and risk zones to ensure comprehensive task coverage
license: MIT
compatibility: opencode
metadata:
  version: 1.0
  audience: xiivthx
---

# sc-map

Survey project structure before planning modifications. Produces `map.json` identifying touchpoints, dependency chains, risk zones, and architectural layers — so `sc-planning` can scope tasks with full awareness and zero side-effect blind spots.

## When to use

Activate when the user asks to:

- Survey a project before `/sc-plan`
- Understand which files a task touches and what depends on them
- Identify risk zones before implementation
- Map project structure for a new mission

Commander auto-triggers: before `/sc-plan` when mission has `spec.md` but no `outputs/map.json`.

## Workflow

Use this exact 3-phase sequence:

### Phase 0 — Pre-flight

1. **Resolve mission.** Read mission `spec.md` to extract keywords, intent, and scope.
2. **Determine mode.** If `outputs/map.json` exists AND `--incremental` flag is set, run delta-only (re-discover changed files via `git diff`). Otherwise full scan.
3. **Resolve PROJECT_ROOT.** Use repo root (git top-level). If mission targets a subdirectory, scope to that path.

### Phase 1 — DISCOVER (deterministic)

**Goal:** Build a file inventory, keyword-hit map, and import/dependency graph using only deterministic tools (find, grep) — no LLM yet.

#### Step 1.1 — Extract keywords from spec

Read mission `spec.md`. Extract:
- **Intent keywords:** nouns and verbs describing what to change (e.g., "git", "branch", "merge", "commit")
- **File hints:** explicit file paths mentioned in spec
- **Scope hints:** subdirectories, modules, or layers mentioned

Store as `$KEYWORDS` and `$FILE_HINTS`.

#### Step 1.2 — Build file inventory

Run deterministic file discovery:

```bash
# All tracked files (excluding node_modules, .git, dist, build)
git ls-files | grep -v -E '(node_modules|\.git|dist|build|\.DS_Store)' > /tmp/sc-map-files.txt
```

For each file, classify by category:
- `code` — source files (.ts, .js, .go, .py, .rs, .sh, etc.)
- `config` — configuration (.json, .yaml, .yml, .toml, .env.example)
- `docs` — documentation (.md, .rst, .txt)
- `skill` — Spacecraft skill files (matches `.opencode/skills/*/`)
- `script` — scripts, Makefiles, tooling
- `test` — test files (matches `*test*`, `*spec*`, `tests/`)

Store as `$FILE_INVENTORY` (path + category + line count).

```bash
# Get line counts for all files
wc -l $(cat /tmp/sc-map-files.txt) 2>/dev/null | tail -n +2 > /tmp/sc-map-linecounts.txt
```

#### Step 1.3 — Keyword grep

For each keyword in `$KEYWORDS`, grep across the file inventory:

```bash
grep -rli "$keyword" $(cat /tmp/sc-map-files.txt) 2>/dev/null
```

Build a hit map: `{keyword: [file_paths]}`. Files with multiple keyword hits get higher relevance.

Store as `$KEYWORD_HITS`.

#### Step 1.4 — Build dependency graph

For each code file, extract imports/dependencies:

```bash
# TypeScript/JavaScript: import/require
grep -rnE '(import |require\()' --include='*.ts' --include='*.js' --include='*.tsx' --include='*.jsx' .

# Go: import blocks
grep -rnE '^\s+"' --include='*.go' .

# Python: import/from
grep -rnE '(^import |^from )' --include='*.py' .

# Shell: source/.
grep -rnE '(^source |^\. )' --include='*.sh' .

# Markdown: links and references
grep -rnE '\[.*\]\(.*\)' --include='*.md' .
```

Build a dependency map: `{file: [imported_files]}`. Reverse it to get dependents map.

Store as `$DEPENDENCY_GRAPH`.

#### Step 1.5 — Identify entry points

Check for common entry-point patterns (in order):
- `opencode.json` (Spacecraft config)
- `AGENTS.md` (agent rules)
- `scripts/spacecraft` (CLI entry)
- `src/main.*`, `src/index.*`, `main.*`, `index.*`

Store as `$ENTRY_POINTS`.

#### Step 1.6 — Compile discovery output

Write intermediate discovery data:

```
.space/missions/<id>/outputs/map-discovery.json
```

Contents:
- `keywords`: extracted from spec
- `fileHints`: explicit paths from spec
- `inventory`: all files with category + line count
- `keywordHits`: per-keyword file matches
- `dependencyGraph`: imports/dependencies per file
- `dependentsGraph`: reverse dependency (who imports this file)
- `entryPoints`: identified entry files

### Phase 2 — ANALYZE (LLM)

**Goal:** Feed Phase 1 discovery data to the LLM for semantic analysis. The commander reads `map-discovery.json` and performs this analysis using its own reasoning — no subagent dispatch needed for a tactical survey.

#### Step 2.1 — Touchpoint identification

From the discovery data, identify:

1. **Direct touchpoints** — files explicitly hinted in spec + files with highest keyword-hit density
2. **Indirect touchpoints** — files that share imports/dependencies with direct touchpoints
3. **Config touchpoints** — config files that control behavior of touched code

Sort by relevance (keyword hit count × dependency weight).

#### Step 2.2 — Dependency chain analysis

From the dependency graph, trace:

1. **Upstream chain** — what does each touchpoint import? (files it depends on)
2. **Downstream chain** — what imports each touchpoint? (files that depend on it — potential side effects)
3. **Shared dependencies** — files imported by multiple touchpoints (high-risk, change once affects many)

Mark shared dependencies with >3 dependents as **high-risk**.

#### Step 2.3 — Risk zone classification

Classify each file:

| Zone | Criteria | Action |
|------|----------|--------|
| **Red** | Shared utility/types imported by 5+ files | Must review before changes |
| **Yellow** | Imported by 2-4 files, or config that gates behavior | Review recommended |
| **Green** | Leaf file, single consumer, or isolated module | Safe to change |
| **Gray** | Doc, asset, unused, or test-only | Low impact |

#### Step 2.4 — Layer classification

Assign each touchpoint and its dependencies to an architectural layer:

| Layer | Spacecraft example |
|-------|-------------------|
| `skills` | `.opencode/skills/*/SKILL.md` |
| `agents` | `AGENTS.md`, `PERSONA.md` |
| `scripts` | `scripts/`, `scripts/src/` |
| `config` | `opencode.json`, `Makefile`, `package.json` |
| `docs` | `DESIGN.md`, `CHANGELOG.md`, `PERSONA.md`, `docs/` |
| `missions` | `.space/missions/`, `.space/archive/` |
| `tests` | `tests/` |

Use directory structure + content signals (frontmatter, imports, purpose) to classify.

### Phase 3 — OUTPUT

**Goal:** Assemble `map.json` and validate completeness.

#### map.json Schema

```json
{
  "version": "1.0.0",
  "missionId": "M07H3CM5S",
  "generatedAt": "<ISO 8601>",
  "project": {
    "root": "/path/to/repo",
    "name": "spacecraft",
    "languages": ["typescript", "markdown", "shell", "go"],
    "totalFiles": 150,
    "scannedFiles": 150
  },
  "spec": {
    "keywords": ["git", "branch", "merge", "commit", "tag"],
    "fileHints": [".opencode/skills/sc-git/", "scripts/spacecraft"],
    "intent": "create sc-map skill for project structure survey"
  },
  "files": [
    {
      "path": ".opencode/skills/sc-git/SKILL.md",
      "category": "skill",
      "layer": "skills",
      "lines": 180,
      "relevance": "direct",
      "riskZone": "green"
    }
  ],
  "touchpoints": [
    {
      "path": ".opencode/skills/sc-map/SKILL.md",
      "reason": "new skill being created — primary artifact",
      "priority": 1,
      "keywords": ["map", "survey", "structure"],
      "dependencies": [],
      "dependents": ["AGENTS.md", "opencode.json"],
      "riskZone": "green"
    }
  ],
  "dependencies": {
    "graph": {
      "fileA": ["fileB", "fileC"]
    },
    "shared": [
      {
        "path": "scripts/spacecraft",
        "importedBy": 8,
        "riskZone": "red"
      }
    ]
  },
  "riskZones": {
    "red": ["scripts/spacecraft"],
    "yellow": ["opencode.json", "AGENTS.md"],
    "green": [".opencode/skills/sc-map/SKILL.md"]
  },
  "layers": {
    "skills": [".opencode/skills/sc-git/SKILL.md", ".opencode/skills/sc-map/SKILL.md"],
    "agents": ["AGENTS.md", "PERSONA.md"],
    "scripts": ["scripts/spacecraft", "scripts/src/main.go"],
    "config": ["opencode.json", "Makefile"],
    "docs": ["DESIGN.md", "CHANGELOG.md", "PERSONA.md"],
    "missions": [".space/missions/M07H3CM5S/"],
    "tests": ["tests/"]
  },
  "coverageGap": []
}
```

#### Step 3.1 — Validate

Before writing final output, check:

- [ ] Every touchpoint exists in file inventory
- [ ] Every dependency reference resolves to an existing file
- [ ] Shared dependencies have correct `importedBy` count
- [ ] No file appears in multiple risk zones
- [ ] Every file in inventory has a category
- [ ] Layers cover all non-trivial files

#### Step 3.2 — Write output

```bash
# Write to mission outputs directory
# Path: .space/missions/<mission-id>/outputs/map.json
```

Report summary to user:
- Total files scanned
- Touchpoints identified (with priority)
- Risk zones: red / yellow / green counts
- Shared dependencies flagged
- Coverage gaps (if any)

## Rules

- **Must**: Run Phase 1 deterministically before Phase 2 analysis — no LLM guessing about file existence.
- **Must**: Use `git ls-files` as authoritative file list (respects `.gitignore`).
- **Must**: Write `map-discovery.json` intermediate for traceability.
- **Must**: Validate `map.json` before claiming done.
- **Must**: Mark shared deps with >3 consumers as red risk zone.
- **Must**: Include coverage gaps explicitly — what was NOT mapped and why.
- **Must not**: Skip Phase 1 even for small projects — keyword extraction from spec is critical.
- **Must not**: Require external tools beyond find, grep, git, and the LLM.
- **Must not**: Produce a map larger than the project itself — keep touchpoints focused on spec scope.
- **May**: Run incrementally with `--incremental` flag (git diff since last map).

## Integration points

### Commander (AGENTS.md)

Commander auto-triggers sc-map when:
- Mission has `spec.md` but `outputs/map.json` is missing
- Before running `/sc-plan` for the first time on a mission

### sc-planning

`sc-planning` reads `map.json` as optional input:
- Uses `touchpoints` to scope task files
- Uses `dependencies.shared` to flag cross-cutting concerns
- Uses `riskZones` to warn about high-risk tasks

### /sc-build

`/sc-build` references `map.json` for:
- Listing files to modify per task
- Warning when touching red-zone files
- Suggesting test files for changed code (from dependency graph)

## Out of scope

This skill does NOT handle:

- Full knowledge graph construction — use Understand-Anything for interactive exploration
- Temporal/git-history analysis — use git log/blame directly
- Test generation — use sc-tester
- Implementation — use sc-coder or /sc-build
- UI design — use sc-design

## Output format

```
map.json written to .space/missions/<mission-id>/outputs/map.json
Summary:
  Files scanned: 150
  Touchpoints: 5 (priority 1: 2, priority 2: 3)
  Risk zones: red 1, yellow 4, green 5
  Layers: skills(12) agents(2) scripts(8) config(3) docs(5) missions(2) tests(1)
  Coverage gaps: none
```

## Checklist

Before claiming sc-map is done:

- [ ] Phase 1 discovery data written to `map-discovery.json`
- [ ] Phase 2 analysis covers touchpoints, dependencies, risks, layers
- [ ] Phase 3 `map.json` validates (no dangling refs, all fields present)
- [ ] Touchpoints are priority-sorted
- [ ] Shared dependencies >3 consumers flagged as red
- [ ] Coverage gaps documented
- [ ] Summary reported to user

---

## References

- Understand-Anything: multi-agent codebase knowledge graph (scan → batch → analyze → review → save)
- Graphiti: temporal context graphs with provenance tracking
- `.space/skill-template.md` — Spacecraft skill template reference
- `scripts/spacecraft` — CLI for mission resolution
