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

# Hook self-tests: allow before branch/mutation checks (JSON may mention git verbs).
printf '%s' "$command" | python3 -c '
import sys
cmd = sys.stdin.read()
for p in (
    ".cursor/hooks/hooks_test.sh",
    ".cursor/hooks/check-ship-commands.sh",
    ".cursor/hooks/check-main-write.sh",
):
    if p in cmd:
        sys.exit(0)
sys.exit(1)
' && allow

if [ -n "${SPACECRAFT_HOOK_BRANCH_OVERRIDE:-}" ]; then
  branch="$SPACECRAFT_HOOK_BRANCH_OVERRIDE"
else
  branch=$(git branch --show-current 2>/dev/null || true)
fi

case "$branch" in
  main|master) ;;
  *) allow ;;
esac

# Mutating git on main/master (real invocations, not mere mention in strings)
printf '%s' "$command" | python3 -c '
import re, sys
cmd = sys.stdin.read()
patterns = [
    r"(^|[;&|]|&&|\|\|)\s*git\s+commit\b",
    r"(^|[;&|]|&&|\|\|)\s*git\s+merge\b",
    r"(^|[;&|]|&&|\|\|)\s*git\s+rebase\b",
    r"(^|[;&|]|&&|\|\|)\s*git\s+cherry-pick\b",
    r"(^|[;&|]|&&|\|\|)\s*git\s+reset\b",
    r"(^|[;&|]|&&|\|\|)\s*git\s+push\b",
    r"(^|[;&|]|&&|\|\|)\s*git\s+am\b",
    r"(^|[;&|]|&&|\|\|)\s*git\s+pull\b",
    r"(^|[;&|]|&&|\|\|)\s*git\s+tag\b",
    r"(^|[;&|]|&&|\|\|)\s*git\s+branch\s+(-D|--delete\s+-f|--force\s+-d)\b",
]
for p in patterns:
    if re.search(p, cmd):
        sys.exit(3)
sys.exit(0)
' && allow || deny \
  "Mutating git on main/master is blocked. Use a feature branch (feat/<id>/<title>). Never write on main." \
  "You are on main/master. Do not run mutating git commands here. Checkout a feature branch first: git checkout -b feat/<id>/<title>"
