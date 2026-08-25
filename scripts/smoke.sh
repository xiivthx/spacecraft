#!/bin/sh
# smoke.sh - post-install validation for a spacecraft project-layer install.
#
# Checks: always-on hard-contract + session-start/safety hooks, MCP/hooks JSON,
# .space scaffold, and CLI help. Domain rules/skills: require at least one entry
# under .cursor/rules and .cursor/skills — does NOT require all catalog packs
# (selective SPACECRAFT_PACKS / spacecraft-profile.json subsets are valid).
# Agents stay User-layer only; safety hooks are dual-layer.
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

# 1. File counts - project layer: some rules + skills present (agents are User-layer).
# Min 1 each: selective profiles may install a subset of domain packs, never all-packs.
for pair in "rules:1" "skills:1"; do
  dir=${pair%%:*}; min=${pair##*:}
  n=$(count "$TARGET/.cursor/$dir")
  if [ "$n" -ge "$min" ]; then
    pass "$dir: $n"
  else
    bad "$dir: $n (expected >= $min)"
  fi
done

# 1b. If profile selects packs, assert inventory matches subset (not all catalog packs).
# Frozen map mirrors design-contract v1 (smoke stays independent of Node catalog load).
profile="$TARGET/.cursor/spacecraft-profile.json"
if [ -f "$profile" ] && python3 -m json.tool "$profile" >/dev/null 2>&1; then
  # shellcheck disable=SC2016
  packs=$(python3 -c 'import json,sys; print(" ".join(json.load(open(sys.argv[1])).get("packs") or []))' "$profile")
  # Catalog-managed skill → pack (selectable only). Coming packs never appear in profile.
  check_skill_pack() {
    skill=$1; want=$2
    path="$TARGET/.cursor/skills/$skill"
    if [ "$want" = present ]; then
      if [ -f "$path/SKILL.md" ]; then
        pass "profile skill $skill present"
      else
        bad "profile skill $skill missing (selected pack inventory)"
      fi
    else
      if [ -e "$path" ]; then
        bad "profile skill $skill present but pack not selected (must not require/install all packs)"
      else
        pass "profile skill $skill absent (unselected)"
      fi
    fi
  }
  # Defaults: all catalog domain skills absent unless pack selected.
  fe=absent be=absent db=absent emb=absent qu=absent fpga=absent
  for p in $packs; do
    case "$p" in
      frontend) fe=present ;;
      backend) be=present ;;
      database) db=present ;;
      embedded) emb=present ;;
      quality) qu=present ;;
      fpga) fpga=present ;;
    esac
  done
  check_skill_pack sc-web-frontend "$fe"
  check_skill_pack sc-ux-design "$fe"
  check_skill_pack sc-browser-probe "$fe"
  check_skill_pack sc-web-backend "$be"
  check_skill_pack sc-database "$db"
  check_skill_pack sc-firmware "$emb"
  check_skill_pack sc-rtl "$fpga"
  check_skill_pack sc-rtl-verify "$fpga"
  check_skill_pack sc-security "$qu"
  check_skill_pack sc-performance "$qu"
  check_skill_pack sc-solid "$qu"
  check_skill_pack sc-architect "$qu"
  check_skill_pack sc-diagram "$qu"
  for id in iot pcb management; do
    if [ -e "$TARGET/.cursor/skills/$id" ]; then
      bad "fictional coming-pack skill dir $id present"
    else
      pass "no fictional coming skill dir $id"
    fi
  done
fi

# 1c. session-start + safety hooks (dual-layer for cloud).
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
