#!/bin/sh
# Deny git merge/push/tag unless SPACECRAFT_SHIP=1 and closeout-check passes.
# Exception: SPACECRAFT_QUICK=1 with SPACECRAFT_SHIP=1 skips closeout (no-mission /sc-quick lane).

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

# Classify: real ship git first (never allowlist-bypass), then strict self-test.
ship_class=$(printf '%s' "$command" | python3 -c '
import re, sys

cmd = sys.stdin.read()
# Optional bare VAR=val / export / env prefixes before git (agent ship form).
GIT_PREFIX = (
    r"(?:(?:export\s+)?(?:[A-Za-z_][A-Za-z0-9_]*=\S+\s+)*)?"
    r"(?:env(?:\s+[A-Za-z_][A-Za-z0-9_]*=\S+)*)?\s*"
    r"(?:[^\s;|&]+/)?git(?:\s+(?:-C\s+\S+|-c\s+\S+))*"
)
SEP = r"(?:^|[;&|]|&&|\|\|)\s*"
REAL_SHIP = re.compile(SEP + GIT_PREFIX + r"\s+(merge|push|tag)\b")

PRIMARY_SELF = re.compile(
    r"^\s*(bash|sh)\s+\.?/?\.cursor/hooks/(hooks_test|check-ship-commands|check-main-write)\.sh(\s|$)"
)
PIPE_SELF = re.compile(
    r"^\s*.+\|\s*\.?/?\.cursor/hooks/check-(ship-commands|main-write)\.sh\s*$"
)

if REAL_SHIP.search(cmd):
    print("ship")
    sys.exit(0)
if PRIMARY_SELF.search(cmd):
    print("selftest")
    sys.exit(0)
if PIPE_SELF.search(cmd) and not REAL_SHIP.search(cmd):
    print("selftest")
    sys.exit(0)
print("other")
')

case "$ship_class" in
  selftest) allow ;;
  other) allow ;;
  ship) ;;
  *) allow ;;
esac

# Resolve ship flags: process env (hooks_test / sessionStart) OR command-string
# assignments. Cursor beforeShellExecution does not inherit agent shell exports
# or `VAR=1 cmd` prefixes into the hook process, so the documented
# `SPACECRAFT_SHIP=1 [SPACECRAFT_QUICK=1] git merge|push|tag` form must be
# read from the command text.
ship_flags=$(printf '%s' "$command" | python3 -c '
import os, re, sys
cmd = sys.stdin.read()
ship = os.environ.get("SPACECRAFT_SHIP", "")
quick = os.environ.get("SPACECRAFT_QUICK", "")
# Assignments must prefix the ship git invocation (not echo/strings elsewhere).
pat = re.compile(
    r"(?:^|[;&|]|&&|\|\|)\s*(?:export\s+)?"
    r"((?:[A-Za-z_][A-Za-z0-9_]*=\S+\s+)*)"
    r"(?:env(?:\s+[A-Za-z_][A-Za-z0-9_]*=\S+)*)?\s*"
    r"(?:[^\s;|&]+/)?git(?:\s+(?:-C\s+\S+|-c\s+\S+))*\s+(?:merge|push|tag)\b"
)
for m in pat.finditer(cmd):
    assigns = m.group(1) or ""
    env_chunk = m.group(0)
    blob = assigns + " " + env_chunk
    if re.search(r"(?:^|[\s])SPACECRAFT_SHIP=1(?:\s|$)", blob):
        ship = "1"
    if re.search(r"(?:^|[\s])SPACECRAFT_QUICK=1(?:\s|$)", blob):
        quick = "1"
print(f"{ship}\t{quick}")
')
SHIP_FLAG=$(printf '%s' "$ship_flags" | cut -f1)
QUICK_FLAG=$(printf '%s' "$ship_flags" | cut -f2)

# Real git merge|push|tag: require SPACECRAFT_SHIP=1, then closeout (unless quick lane).
if [ "$SHIP_FLAG" != "1" ]; then
  deny \
    "Ship gate blocked this command. Run /sc-ship (or /sc-quick ship) and set SPACECRAFT_SHIP=1 only for gated git merge/push/tag, then unset it." \
    "Do not merge, push, or tag unless the user explicitly requested ship via /sc-ship or /sc-quick. Prefix gated git with SPACECRAFT_SHIP=1 (and SPACECRAFT_QUICK=1 for no-mission /sc-quick), then unset after."
fi

# No-mission /sc-quick: skip mission closeout-check.
if [ "$QUICK_FLAG" = "1" ]; then
  allow
fi

# SPACECRAFT_SHIP=1 (mission ship): run closeout before allowing.
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
    "Ship gate blocked: closeout-check failed. Fix closeout issues, then retry with SPACECRAFT_SHIP=1. For no-mission /sc-quick use SPACECRAFT_SHIP=1 SPACECRAFT_QUICK=1." \
    "closeout-check failed (exit $closeout_rc). Output:
$truncated"
fi

allow
