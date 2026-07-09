> Consult when: writing or validating `map.json` output — schema reference and field descriptions.

# map.json Schema

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

## Field reference

| Field | Type | Description |
|-------|------|-------------|
| `version` | string | Schema version (semver) |
| `missionId` | string | Mission ID this map was generated for |
| `generatedAt` | string | ISO 8601 timestamp |
| `project.root` | string | Absolute path to repo root |
| `project.name` | string | Project name |
| `project.languages` | string[] | Detected languages |
| `project.totalFiles` | int | Total tracked files |
| `project.scannedFiles` | int | Files included in analysis |
| `spec.keywords` | string[] | Extracted from mission spec.md |
| `spec.fileHints` | string[] | Explicit file paths from spec |
| `spec.intent` | string | Mission intent summary |
| `files` | object[] | All analyzed files with path, category, layer, lines, relevance, riskZone |
| `touchpoints` | object[] | Files most relevant to the mission — priority-sorted with dependencies and dependents |
| `dependencies.graph` | map[string][]string | Per-file import/dependency list |
| `dependencies.shared` | object[] | Files imported by multiple consumers — high-risk flag |
| `riskZones` | object | Files grouped by risk category: red/yellow/green |
| `layers` | object | Files grouped by architectural layer |
| `coverageGap` | string[] | Files or areas not mapped, with reasons |

## Spacecraft integration

- Written to `.space/missions/<id>/outputs/map.json`
- Consumed by `sc-planning` for task scoping and risk warnings
- Consumed by `/sc-build` for file lists and red-zone alerts
- Auto-triggered by Commander before `/sc-plan` when `map.json` is missing
