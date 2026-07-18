#!/bin/sh
# Deny git merge/push/tag unless SPACECRAFT_SHIP=1 and closeout-check passes.

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

truncate_msg() {
  python3 -c '
import sys
text = sys.stdin.read()
limit = 1500
if len(text) > limit:
    text = text[:limit] + "\n...[truncated]"
print(text)
'
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
  "The ship gate hook could not parse the shell command JSON from stdin. Fix the hook input or retry."

# Hook self-tests: allow even if the command string mentions git push/merge/tag.
# Must run before closeout so self-tests never hit live closeout-check.
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

# Only gate real ship invocations (not matcher false positives like `ls`).
printf '%s' "$command" | python3 -c '
import re, sys
cmd = sys.stdin.read()
if re.search(r"(^|[;&|]|&&|\|\|)\s*git\s+(merge|push|tag)\b", cmd):
    sys.exit(3)
sys.exit(0)
' && allow

# Real git merge|push|tag: require SPACECRAFT_SHIP=1, then closeout.
if [ "${SPACECRAFT_SHIP:-}" != "1" ]; then
  deny \
    "Ship gate blocked this command. Run /sc-ship and set SPACECRAFT_SHIP=1 only for gated git merge/push/tag, then unset it." \
    "Do not merge, push, or tag unless the user explicitly requested ship via /sc-ship. Export SPACECRAFT_SHIP=1 for those gated git commands only, then unset SPACECRAFT_SHIP."
fi

# SPACECRAFT_SHIP=1: run closeout before allowing (never allow immediately).
if [ -n "${SPACECRAFT_CLOSEOUT_CMD:-}" ]; then
  closeout_out=$(sh -c "$SPACECRAFT_CLOSEOUT_CMD" 2>&1)
  closeout_rc=$?
elif [ -x ./spacecraft ]; then
  closeout_out=$(./spacecraft closeout-check 2>&1)
  closeout_rc=$?
else
  closeout_out=$(spacecraft closeout-check 2>&1)
  closeout_rc=$?
fi

if [ "$closeout_rc" -ne 0 ]; then
  truncated=$(printf '%s' "$closeout_out" | truncate_msg)
  deny \
    "Ship gate blocked: closeout-check failed. Fix closeout issues, then retry with SPACECRAFT_SHIP=1." \
    "closeout-check failed (exit $closeout_rc). Output:
$truncated"
fi

allow
