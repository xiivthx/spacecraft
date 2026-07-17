#!/bin/sh
# config-smoke.sh — static validation of the spacecraft .cursor config surface.
#
# Checks (no build required):
#   - mcp.json parses as JSON (hooks.json too, when present)
#   - every skill (.cursor/skills/*/SKILL.md) has YAML frontmatter with name + description
#   - every agent (.cursor/agents/*.md) has YAML frontmatter with name + description
#   - the legacy .cursor/commands/ directory is absent
#
# Usage: config-smoke.sh [root-dir]   (default: current directory)
set -e

ROOT="${1:-.}"
CURSOR="$ROOT/.cursor"

fail=0
pass() { echo "  ok   $1"; }
bad()  { echo "  FAIL $1"; fail=1; }

echo "Config smoke: $ROOT"

# 1. JSON files parse.
if [ -f "$CURSOR/mcp.json" ]; then
  if python3 -m json.tool "$CURSOR/mcp.json" >/dev/null 2>&1; then
    pass "mcp.json parses"
  else
    bad "mcp.json is invalid JSON"
  fi
else
  bad "mcp.json missing"
fi

if [ -f "$CURSOR/hooks.json" ]; then
  if python3 -m json.tool "$CURSOR/hooks.json" >/dev/null 2>&1; then
    pass "hooks.json parses"
  else
    bad "hooks.json is invalid JSON"
  fi
fi

# frontmatter <file> <label> — verify leading --- block carries name + description.
frontmatter() {
  file="$1"; label="$2"
  # First line must open the YAML frontmatter block.
  if [ "$(sed -n '1p' "$file")" != "---" ]; then
    bad "$label: missing frontmatter opener"
    return
  fi
  # Extract the block between the first two --- markers.
  block=$(awk 'NR==1{next} /^---[[:space:]]*$/{exit} {print}' "$file")
  if ! printf '%s\n' "$block" | grep -Eq '^name:[[:space:]]*[^[:space:]]'; then
    bad "$label: frontmatter missing name"
    return
  fi
  if ! printf '%s\n' "$block" | grep -Eq '^description:[[:space:]]*[^[:space:]]'; then
    bad "$label: frontmatter missing description"
    return
  fi
  pass "$label"
}

# 2. Skills frontmatter.
skills=0
for f in "$CURSOR"/skills/*/SKILL.md; do
  [ -f "$f" ] || continue
  skills=$((skills + 1))
  frontmatter "$f" "skill $(basename "$(dirname "$f")")"
done
if [ "$skills" -eq 0 ]; then
  bad "no skills found under $CURSOR/skills"
else
  pass "skills scanned: $skills"
fi

# 3. Agents frontmatter.
agents=0
for f in "$CURSOR"/agents/*.md; do
  [ -f "$f" ] || continue
  agents=$((agents + 1))
  frontmatter "$f" "agent $(basename "$f" .md)"
done
if [ "$agents" -eq 0 ]; then
  bad "no agents found under $CURSOR/agents"
else
  pass "agents scanned: $agents"
fi

# 4. Legacy commands/ directory must be gone.
if [ -d "$CURSOR/commands" ]; then
  bad ".cursor/commands/ present (should be removed)"
else
  pass "no .cursor/commands/"
fi

if [ "$fail" -ne 0 ]; then
  echo "Config smoke FAILED"
  exit 1
fi
echo "Config smoke passed"
