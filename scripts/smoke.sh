#!/bin/sh
# smoke.sh - post-install validation for a spacecraft project-layer install.
#
# Checks: domain rules/skills present, session-start + safety hooks wired, mcp.json
# parses as JSON, hooks.json parses when present, and the spacecraft CLI
# reports help. Agents stay User-layer only; safety hooks are dual-layer.
#
# Usage: smoke.sh <target-project-dir> [spacecraft-binary]
set -e

TARGET="${1:?usage: smoke.sh <target-project-dir> [spacecraft-binary]}"
BIN="${2:-}"

fail=0
pass() { echo "  ok   $1"; }
warn() { echo "  warn $1"; }
bad()  { echo "  FAIL $1"; fail=1; }

count() { ls -1 "$1" 2>/dev/null | wc -l | tr -d ' '; }

echo "Smoke check: $TARGET"

# 1. File counts - project layer: domain rules + domain skills (agents are User-layer).
for pair in "rules:1" "skills:1"; do
  dir=${pair%%:*}; min=${pair##*:}
  n=$(count "$TARGET/.cursor/$dir")
  if [ "$n" -ge "$min" ]; then
    pass "$dir: $n"
  else
    bad "$dir: $n (expected >= $min)"
  fi
done

# 1b. session-start + safety hooks (dual-layer for cloud).
if [ -f "$TARGET/.cursor/hooks/session-start.sh" ]; then
  pass "session-start.sh present"
else
  bad "session-start.sh missing under .cursor/hooks/"
fi
if [ -f "$TARGET/.cursor/hooks.json" ] && grep -q 'session-start.sh' "$TARGET/.cursor/hooks.json"; then
  pass "hooks.json wires session-start"
else
  bad "hooks.json missing session-start.sh"
fi
for safety_hook in check-main-write.sh check-ship-commands.sh block-secrets-read.sh block-destructive.sh; do
  if [ -f "$TARGET/.cursor/hooks/$safety_hook" ]; then
    pass "safety hook $safety_hook present"
  else
    bad "safety hook $safety_hook missing under project .cursor/hooks/"
  fi
done
if [ -f "$TARGET/.cursor/rules/010-hard-contract.mdc" ]; then
  pass "010-hard-contract.mdc present"
else
  bad "010-hard-contract.mdc missing under project .cursor/rules/"
fi
leftover_agents=$(find "$TARGET/.cursor/agents" -maxdepth 1 -type f -name 'sc-*.md' 2>/dev/null | head -n 1)
if [ -n "$leftover_agents" ]; then
  bad "spacecraft agent present under project .cursor/agents (User-layer only): $leftover_agents"
else
  pass "no project-local spacecraft agents"
fi

# 2. .space scaffold.
for d in missions archive roadmaps; do
  if [ -d "$TARGET/.space/$d" ]; then
    pass ".space/$d present"
  else
    bad ".space/$d missing"
  fi
done

# 2b. Project-git ensure (install / first .space create).
if [ -d "$TARGET/.git" ] || git -C "$TARGET" rev-parse --git-dir >/dev/null 2>&1; then
  pass "git repo present"
else
  bad "git repo missing (.git / rev-parse)"
fi
if [ -f "$TARGET/.gitignore" ] && grep -Eq '^[[:space:]]*\.space/?[[:space:]]*$' "$TARGET/.gitignore"; then
  pass ".gitignore ignores .space/"
else
  bad ".gitignore missing .space/ entry"
fi


# 3. JSON parses.
if [ -f "$TARGET/.cursor/mcp.json" ]; then
  if python3 -m json.tool "$TARGET/.cursor/mcp.json" >/dev/null 2>&1; then
    pass "mcp.json parses"
  else
    bad "mcp.json is invalid JSON"
  fi
else
  warn "mcp.json not present"
fi

if [ -f "$TARGET/.cursor/hooks.json" ]; then
  if python3 -m json.tool "$TARGET/.cursor/hooks.json" >/dev/null 2>&1; then
    pass "hooks.json parses"
  else
    bad "hooks.json is invalid JSON"
  fi
fi

# 4. CLI help.
if [ -n "$BIN" ] && [ -x "$BIN" ]; then
  if "$BIN" help >/dev/null 2>&1; then
    pass "spacecraft help"
  else
    bad "spacecraft help failed"
  fi
elif command -v spacecraft >/dev/null 2>&1; then
  if spacecraft help >/dev/null 2>&1; then
    pass "spacecraft help (PATH)"
  else
    bad "spacecraft help failed"
  fi
else
  warn "spacecraft binary not found; skipping CLI check"
fi

if [ "$fail" -ne 0 ]; then
  echo "Smoke check FAILED"
  exit 1
fi
echo "Smoke check passed"
