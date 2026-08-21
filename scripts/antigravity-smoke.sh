#!/bin/sh
# antigravity-smoke.sh - smoke tests for Antigravity Spacecraft plugin & configuration.
set -e

ROOT="${1:-.}"
PLUGIN="$ROOT/plugins/spacecraft"

echo "Spacecraft Antigravity Smoke Test"
echo "================================="

# 1. Validate plugin.json
if [ ! -f "$PLUGIN/plugin.json" ]; then
  echo "FAIL: missing $PLUGIN/plugin.json" >&2
  exit 1
fi
node -e "JSON.parse(require('fs').readFileSync('$PLUGIN/plugin.json', 'utf8'))"
echo "  ok: plugin.json parses cleanly"

# 2. Validate hooks.json
if [ ! -f "$PLUGIN/hooks.json" ]; then
  echo "FAIL: missing $PLUGIN/hooks.json" >&2
  exit 1
fi
node -e "JSON.parse(require('fs').readFileSync('$PLUGIN/hooks.json', 'utf8'))"
echo "  ok: hooks.json parses cleanly"

# 3. Test safety hook execution
node -e '
const { execSync } = require("child_process");
const script = "'"$PLUGIN"'/hooks/safety-check.mjs";

function run(toolCall) {
  return JSON.parse(execSync("node " + script, {
    input: JSON.stringify({ toolCall }),
    encoding: "utf8"
  }));
}

const safe = run({ name: "run_command", args: { CommandLine: "ls -la" } });
if (safe.decision !== "allow") throw new Error("safe command not allowed");

const status = run({ name: "run_command", args: { CommandLine: "git status" } });
if (status.decision !== "allow") throw new Error("git status not allowed");

const shipPush = run({ name: "run_command", args: { CommandLine: "SPACECRAFT_SHIP=1 git push" } });
if (shipPush.decision !== "deny") throw new Error("ship push must be denied in-agent");

const force = run({ name: "run_command", args: { CommandLine: "git push --force origin main" } });
if (force.decision !== "deny") throw new Error("force push must be denied");

const merge = run({ name: "run_command", args: { CommandLine: "SPACECRAFT_SHIP=1 git merge --no-ff feat/x" } });
if (merge.decision !== "allow") throw new Error("ship merge not allowed");
'
echo "  ok: safety-check.mjs decisions pass"

# 4. Check rules/AGENTS.md and GEMINI.md (short hard contract)
for rule_file in "$PLUGIN/rules/AGENTS.md" "$ROOT/GEMINI.md"; do
  if [ ! -f "$rule_file" ]; then
    echo "FAIL: missing $rule_file" >&2
    exit 1
  fi
  grep -qi "hard contract" "$rule_file" || { echo "FAIL: missing hard contract in $rule_file" >&2; exit 1; }
  grep -q "AUTH:" "$rule_file" || { echo "FAIL: missing AUTH: in $rule_file" >&2; exit 1; }
  grep -q "safety-check" "$rule_file" || { echo "FAIL: missing safety-check pointer in $rule_file" >&2; exit 1; }
done
echo "  ok: AGENTS.md & GEMINI.md contain short hard contract"


# 5. Check skills
skill_count=0
for s in "$PLUGIN"/skills/sc-*; do
  [ -d "$s" ] || continue
  [ -f "$s/SKILL.md" ] || { echo "FAIL: missing $s/SKILL.md" >&2; exit 1; }
  # Validate YAML frontmatter
  node -e "
    const fs = require('fs');
    const content = fs.readFileSync('$s/SKILL.md', 'utf8');
    if (!content.startsWith('---')) throw new Error('missing frontmatter in $s');
  "
  skill_count=$((skill_count + 1))
done
echo "  ok: $skill_count skills verified in Antigravity plugin"

# 6. Check agents
agent_count=0
for a in "$PLUGIN"/agents/sc-*.md; do
  [ -f "$a" ] || continue
  agent_count=$((agent_count + 1))
done
echo "  ok: $agent_count subagents verified in Antigravity plugin"

echo "All Antigravity smoke checks passed!"
