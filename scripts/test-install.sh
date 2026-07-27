#!/bin/sh
# test-install.sh - install/bootstrap smoke test in a throwaway temp dir.
#
# Seeds a pre-existing unrelated hook so the hooks merge must preserve user
# hooks alongside the spacecraft ship-gate hook, then runs install-global
# against a fake HOME and checks the core skills land. Never writes
# ~/.cursorrules.
#
# Usage: test-install.sh <repo-root> <spacecraft-binary>
set -e

ROOT="${1:?usage: test-install.sh <repo-root> <spacecraft-binary>}"
BIN="${2:?usage: test-install.sh <repo-root> <spacecraft-binary>}"

mkdir -p "$ROOT/.tmp"
tmp=$(mktemp -d "$ROOT/.tmp/install-smoke.XXXXXX")
fake_home="$tmp/home"
# Always succeed so EXIT trap cannot flip a green run to status 1.
cleanup() { rm -rf "$tmp"; }
trap cleanup EXIT INT TERM

echo "Install smoke in $tmp"
mkdir -p "$tmp/.cursor" "$fake_home"
printf '%s\n' '{"version":1,"hooks":{"beforeShellExecution":[{"command":".cursor/hooks/user-unrelated.sh","matcher":"echo"}]}}' \
  > "$tmp/.cursor/hooks.json"

HOME="$fake_home" sh "$ROOT/scripts/install-cursor.sh" "$tmp" "$ROOT"
sh "$ROOT/scripts/smoke.sh" "$tmp" "$BIN"

hooks="$tmp/.cursor/hooks.json"
if ! grep -q 'user-unrelated.sh' "$hooks"; then
  echo "FAIL: install clobbered pre-existing hooks (user-unrelated.sh missing)"
  exit 1
fi
if ! grep -q 'check-ship-commands.sh' "$hooks"; then
  echo "FAIL: install missing spacecraft ship-gate hook (check-ship-commands.sh)"
  exit 1
fi
echo "  ok   hooks merge preserves user + ship-gate"

if [ -f "$fake_home/.cursorrules" ]; then
  echo "FAIL: install wrote legacy ~/.cursorrules"
  exit 1
fi
echo "  ok   no legacy ~/.cursorrules"

# T4 acceptance 1: install-project places domain rules 300-620 + the
# session-start hook into a project target distinct from the source repo.
for rule in 300-security 400-performance 500-database 600-firmware \
  610-firmware-peripherals 620-firmware-testing; do
  test -f "$tmp/.cursor/rules/$rule.mdc" \
    || { echo "FAIL: install-project missing domain rule $rule.mdc in $tmp/.cursor/rules"; exit 1; }
done
grep -q 'session-start.sh' "$hooks" \
  || { echo "FAIL: install-project did not wire session-start hook into $hooks"; exit 1; }
test -f "$tmp/.cursor/hooks/session-start.sh" \
  || { echo "FAIL: install-project did not copy session-start.sh into $tmp/.cursor/hooks/"; exit 1; }
echo "  ok   install-project places domain rules 300-620 + session-start hook"

# T4 acceptance 2: install-project does NOT copy alwaysApply rules (000/025/
# 050/100/200) into a project target distinct from the source repo.
for rule in 000-spacecraft 025-english-coach 050-style 100-conventions 200-workflow; do
  if [ -f "$tmp/.cursor/rules/$rule.mdc" ]; then
    echo "FAIL: install-project copied alwaysApply rule $rule.mdc into project target $tmp"
    exit 1
  fi
done
echo "  ok   install-project excludes alwaysApply rules from project target"

# STORM Tier 0/3: project install must land sc-storm + discuss lens-pass reference.
test -f "$tmp/.cursor/skills/sc-storm/SKILL.md" \
  || { echo "FAIL: install-project missing sc-storm skill"; exit 1; }
test -f "$tmp/.cursor/skills/sc-discuss/references/lens-pass.md" \
  || { echo "FAIL: install-project missing sc-discuss/references/lens-pass.md"; exit 1; }
echo "  ok   install-project installs sc-storm and discuss lens-pass reference"

mkdir -p "$fake_home/.cursor"
printf '%s\n' '{"version":1,"hooks":{"beforeShellExecution":[{"command":"~/.cursor/hooks/user-unrelated-global.sh","matcher":"echo"}]}}' \
  > "$fake_home/.cursor/hooks.json"

HOME="$fake_home" make -C "$ROOT" install-global \
  GLOBAL="$fake_home/.cursor" LOCAL_BIN="$fake_home/.local/bin" BIN="$BIN"
for skill in sc-run sc-ship sc-quick sc-storm; do
  test -f "$fake_home/.cursor/skills/$skill/SKILL.md" \
    || { echo "FAIL: install-global missing $skill skill"; exit 1; }
done
test -f "$fake_home/.cursor/skills/sc-discuss/references/lens-pass.md" \
  || { echo "FAIL: install-global missing sc-discuss/references/lens-pass.md"; exit 1; }
echo "  ok   install-global installs sc-run, sc-ship, sc-quick, sc-storm, and lens-pass"

user_rules="$fake_home/.cursor/spacecraft/USER-RULES.txt"
test -f "$user_rules" \
  || { echo "FAIL: install-global did not write $user_rules"; exit 1; }
for marker in 'Spacecraft' 'English prompt coach' 'Coding Standards' 'Project Structure' 'Lane Detection' 'lens pass'; do
  grep -q "$marker" "$user_rules" \
    || { echo "FAIL: USER-RULES.txt missing marker: $marker"; exit 1; }
done
echo "  ok   install-global generates USER-RULES.txt with five-source markers (+ lens pass)"

if [ -f "$fake_home/.cursorrules" ]; then
  echo "FAIL: install-global wrote legacy ~/.cursorrules"
  exit 1
fi
echo "  ok   install-global still writes no legacy ~/.cursorrules"

# T3 acceptance 1: global hooks.json gets the safety hooks merged in, and the
# pre-existing unrelated global hook (seeded above) survives the merge.
global_hooks="$fake_home/.cursor/hooks.json"
if ! grep -q 'user-unrelated-global.sh' "$global_hooks"; then
  echo "FAIL: install-global clobbered pre-existing global hook (user-unrelated-global.sh missing)"
  exit 1
fi
if ! grep -q 'check-main-write.sh' "$global_hooks"; then
  echo "FAIL: install-global did not merge check-main-write.sh into $global_hooks"
  exit 1
fi
if ! grep -q 'check-ship-commands.sh' "$global_hooks"; then
  echo "FAIL: install-global did not merge check-ship-commands.sh into $global_hooks"
  exit 1
fi
echo "  ok   install-global hooks merge preserves unrelated hook + adds safety hooks"

# T3 acceptance 2: installed hook commands use absolute or ~ paths (never the
# project-layer's repo-relative .cursor/hooks/...), and the scripts themselves
# are copied into ~/.cursor/hooks/ so they work outside this repo.
if ! grep -qE '"command": "(~|/)[^"]*check-main-write\.sh"' "$global_hooks"; then
  echo "FAIL: check-main-write.sh command in $global_hooks is not an absolute or ~ path"
  exit 1
fi
if ! grep -qE '"command": "(~|/)[^"]*check-ship-commands\.sh"' "$global_hooks"; then
  echo "FAIL: check-ship-commands.sh command in $global_hooks is not an absolute or ~ path"
  exit 1
fi
for hook_script in check-main-write.sh check-ship-commands.sh; do
  test -f "$fake_home/.cursor/hooks/$hook_script" \
    || { echo "FAIL: install-global did not copy $hook_script into $fake_home/.cursor/hooks/"; exit 1; }
done
echo "  ok   install-global hook paths are absolute/~ and scripts land in ~/.cursor/hooks/"
