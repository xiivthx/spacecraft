# opencode.json Template

OpenCode configuration for Spacecraft missions. Copy to project root as `opencode.json`.

```jsonc
{
  "$schema": "https://opencode.ai/config.json",
  "default_agent": "sc-commander",
  "share": "disabled",
  "lsp": true,
  "instructions": ["PERSONA.md"],
  "agent": {
    "sc-commander": {
      "mode": "primary",
      "description": "Mission commander — orchestrates mission-driven implementation using lean prompts",
      "model": "deepseek/deepseek-v4-pro",
      // "model": "opencode-go/deepseek-v4-pro",
      // "model": "opencode-go/qwen3.7-max",
      "temperature": 0.2,
      "variant": "high",
      "color": "#4fc3f7"
    },
    "sc-planner": {
      "mode": "subagent",
      "description": "Read-only planner that turns a mission spec into a small executable plan",
      "model": "opencode-go/qwen3.7-plus",
      // "model": "opencode/deepseek-v4-flash-free",
      "temperature": 0.1,
      "variant": "low",
      "color": "#80cbc4"
    },
    "sc-designer": {
      "mode": "subagent",
      "description": "Read-only design agent for UI direction, critique, and anti-slop review",
      "model": "opencode-go/glm-5.2",
      // "model": "opencode-go/qwen3.7-plus",
      "temperature": 0.3,
      "variant": "medium",
      "color": "#ce93d8"
    },
    "sc-reviewer": {
      "mode": "subagent",
      "description": "Read-only reviewer for diff, evidence, and release readiness",
      "model": "deepseek/deepseek-v4-pro",
      // "model": "opencode-go/deepseek-v4-pro",
      // "model": "opencode-go/kimi-k2.7-code",
      "temperature": 0.1,
      "variant": "medium",
      "color": "#ef9a9a"
    },
    "sc-coder": {
      "mode": "subagent",
      "description": "Write-capable coder that implements production code to satisfy tasks and tests",
      "model": "opencode-go/kimi-k2.7-code",
      // "model": "opencode-go/deepseek-v4-pro",
      "temperature": 0.2,
      "variant": "low",
      "color": "#a5d6a7"
    },
    "sc-tester": {
      "mode": "subagent",
      "description": "Write-capable tester that writes tests and captures verification evidence (TDD)",
      "model": "opencode-go/kimi-k2.7-code",
      // "model": "opencode-go/deepseek-v4-flash",
      "temperature": 0.1,
      "color": "#ffcc80"
    }
  }
}
```

## Fields

| Field | Description |
|-------|-------------|
| `default_agent` | Primary agent for new sessions |
| `share` | Session sharing (`disabled`, `manual`, `auto`) |
| `lsp` | Enable language server protocol |
| `instructions` | Markdown files loaded as system context |
| `agent.*.mode` | `primary` (user-facing) or `subagent` (delegated) |
| `agent.*.model` | Model ID in `provider/model` format |
| `agent.*.temperature` | Lower = more deterministic |
| `agent.*.variant` | Model variant hint (`low`, `medium`, `high`) |
| `agent.*.color` | UI accent for agent messages |

## Global Install

For cross-workspace use, run `make install` from the Spacecraft repo. This writes to `~/.config/opencode/opencode.jsonc` with absolute paths and adds `skills.paths` + `plugin` for engine discovery.
