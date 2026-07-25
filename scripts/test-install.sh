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

mkdir -p "$fake_home/.cursor"
HOME="$fake_home" make -C "$ROOT" install-global \
  GLOBAL="$fake_home/.cursor" LOCAL_BIN="$fake_home/.local/bin" BIN="$BIN"
for skill in sc-run sc-ship sc-quick; do
  test -f "$fake_home/.cursor/skills/$skill/SKILL.md" \
    || { echo "FAIL: install-global missing $skill skill"; exit 1; }
done
echo "  ok   install-global installs sc-run, sc-ship, and sc-quick"
