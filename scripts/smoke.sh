#!/bin/sh
# smoke.sh — post-install validation for a spacecraft install.
#
# Checks: expected rules/agents/skills are present, mcp.json parses as JSON,
# hooks.json parses when present, and the spacecraft CLI reports help.
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

# 1. File counts — every surface must have at least one entry.
for pair in "rules:1" "agents:1" "skills:1"; do
  dir=${pair%%:*}; min=${pair##*:}
  n=$(count "$TARGET/.cursor/$dir")
  if [ "$n" -ge "$min" ]; then
    pass "$dir: $n"
  else
    bad "$dir: $n (expected >= $min)"
  fi
done

# 2. .space scaffold.
for d in missions archive; do
  if [ -d "$TARGET/.space/$d" ]; then
    pass ".space/$d present"
  else
    bad ".space/$d missing"
  fi
done

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
