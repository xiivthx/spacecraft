#!/bin/sh
# Deny mutating git commands on main/master.
# Tests may set SPACECRAFT_HOOK_BRANCH_OVERRIDE to simulate branch (ignored in production unless set).

deny() {
  user_msg="$1"
  agent_msg="$2"
  printf '%s\n' "{\"permission\":\"deny\",\"user_message\":$(python3 -c 'import json,sys; print(json.dumps(sys.argv[1]))' "$user_msg"),\"agent_message\":$(python3 -c 'import json,sys; print(json.dumps(sys.argv[1]))' "$agent_msg")}"
  exit 0
}

allow() {
  printf '%s\n' '{"permission":"allow"}'
  exit 0
}

input=$(cat)

command=$(printf '%s' "$input" | python3 -c '
import json, sys
raw = sys.stdin.read()
try:
    data = json.loads(raw)
except Exception:
    sys.exit(2)
if not isinstance(data, dict):
    sys.exit(2)
cmd = data.get("command", "")
if not isinstance(cmd, str):
    sys.exit(2)
print(cmd)
' 2>/dev/null) || deny \
  "hook could not parse command" \
  "The main-write hook could not parse the shell command JSON from stdin. Fix the hook input or retry."

# Override may be explicitly empty (detached HEAD simulation); treat set-even-if-empty.
if [ "${SPACECRAFT_HOOK_BRANCH_OVERRIDE+set}" = "set" ]; then
  branch="$SPACECRAFT_HOOK_BRANCH_OVERRIDE"
else
  branch=$(git branch --show-current 2>/dev/null || true)
fi

# Detached HEAD / empty branch: fail closed. Optionally map HEAD to main/master tip name.
case "$branch" in
  main|master) ;;
  "")
    head=$(git rev-parse HEAD 2>/dev/null || true)
    if [ -n "$head" ]; then
      main_tip=$(git rev-parse main 2>/dev/null || true)
      master_tip=$(git rev-parse master 2>/dev/null || true)
      if [ -n "$main_tip" ] && [ "$head" = "$main_tip" ]; then
        branch=main
      elif [ -n "$master_tip" ] && [ "$head" = "$master_tip" ]; then
        branch=master
      fi
    fi
    ;;
  *) allow ;;
esac

# On main/master (or unknown empty): mutate first, then strict self-test only.
mw_class=$(printf '%s' "$command" | python3 -c '
import re, sys

cmd = sys.stdin.read()
GIT_PREFIX = (
    r"(?:(?:export\s+)?(?:[A-Za-z_][A-Za-z0-9_]*=\S+\s+)*)?"
    r"(?:env(?:\s+[A-Za-z_][A-Za-z0-9_]*=\S+)*)?\s*"
    r"(?:[^\s;|&]+/)?git(?:\s+(?:-C\s+\S+|-c\s+\S+))*"
)
SEP = r"(?:^|[;&|]|&&|\|\|)\s*"

MUTATE = re.compile(
    SEP + GIT_PREFIX + r"\s+(commit|merge|rebase|cherry-pick|reset|push|am|pull|tag)\b"
)
BRANCH_DEL = re.compile(
    SEP + GIT_PREFIX + r"\s+branch\s+(-D|--delete\s+-f|--force\s+-d)\b"
)

PRIMARY_SELF = re.compile(
    r"^\s*(bash|sh)\s+\.?/?\.cursor/hooks/(hooks_test|check-ship-commands|check-main-write)\.sh(\s|$)"
)
PIPE_SELF = re.compile(
    r"^\s*.+\|\s*\.?/?\.cursor/hooks/check-(ship-commands|main-write)\.sh\s*$"
)

has_mutate = bool(MUTATE.search(cmd) or BRANCH_DEL.search(cmd))
if has_mutate:
    print("mutate")
    sys.exit(0)
if PRIMARY_SELF.search(cmd):
    print("selftest")
    sys.exit(0)
if PIPE_SELF.search(cmd):
    print("selftest")
    sys.exit(0)
print("other")
')

case "$mw_class" in
  mutate)
    # Allow ship-gated merge/push/tag on main (documented SPACECRAFT_SHIP=1 form).
    # Cursor hooks do not inherit agent shell env; read assignment from command.
    ship_ok=$(printf '%s' "$command" | python3 -c '
import os, re, sys
cmd = sys.stdin.read()
env_ship = os.environ.get("SPACECRAFT_SHIP", "") == "1"
pat = re.compile(
    r"(?:^|[;&|]|&&|\|\|)\s*(?:export\s+)?"
    r"((?:[A-Za-z_][A-Za-z0-9_]*=\S+\s+)*)"
    r"(?:env(?:\s+[A-Za-z_][A-Za-z0-9_]*=\S+)*)?\s*"
    r"(?:[^\s;|&]+/)?git(?:\s+(?:-C\s+\S+|-c\s+\S+))*\s+(?:merge|push|tag)\b"
)
cmd_ship = False
for m in pat.finditer(cmd):
    blob = (m.group(1) or "") + " " + m.group(0)
    if re.search(r"(?:^|[\s])SPACECRAFT_SHIP=1(?:\s|$)", blob):
        cmd_ship = True
        break
# Exempt only merge/push/tag when ship gate is present (env or command prefix).
print("1" if (env_ship or cmd_ship) and pat.search(cmd) else "0")
')
    if [ "$ship_ok" = "1" ]; then
      allow
    fi
    deny \
      "Mutating git on main/master is blocked. Use a feature branch (feat/<id>/<title>). Never write on main." \
      "You are on main/master. Do not run mutating git commands here. Checkout a feature branch first: git checkout -b feat/<id>/<title>. Ship merge/push/tag must use SPACECRAFT_SHIP=1 prefix."
    ;;
  selftest|other) allow ;;
  *) allow ;;
esac
