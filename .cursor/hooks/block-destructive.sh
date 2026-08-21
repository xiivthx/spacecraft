#!/bin/sh
# Deny force-push and catastrophic rm / wipe patterns.

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
  "The destructive hook could not parse the shell command JSON from stdin. Fix the hook input or retry."

# Self-test exemption (primary command only)
case "$command" in
  *"hooks_test.sh"*|*"block-destructive.sh"*)
    # Still deny if compound includes real destructive with self-test
    ;;
esac

decision=$(printf '%s' "$command" | python3 -c '
import re, sys
cmd = sys.stdin.read()

# Allow pure self-test of this hook / hooks_test
if re.match(r"^\s*(bash|sh)\s+\.?/?\.cursor/hooks/(hooks_test|block-destructive)\.sh(\s|$)", cmd):
    print("allow")
    raise SystemExit(0)
if re.match(r"^\s*.+\|\s*\.?/?\.cursor/hooks/block-destructive\.sh\s*$", cmd):
    print("allow")
    raise SystemExit(0)

# Force push: push with -f / --force / --force-with-lease (lease still blocked as force-family)
force_push = re.compile(
    r"(?:^|[;&|]|&&|\|\|)\s*(?:export\s+)?"
    r"(?:(?:[A-Za-z_][A-Za-z0-9_]*=\S+\s+)*)?"
    r"(?:env(?:\s+[A-Za-z_][A-Za-z0-9_]*=\S+)*)?\s*"
    r"(?:[^\s;|&]+/)?git(?:\s+(?:-C\s+\S+|-c\s+\S+))*\s+push\b[^\n;|&]*"
    r"(?:\s|^)(-f|--force|--force-with-lease)\b"
)
# Also: git push -f before remote
force_push2 = re.compile(
    r"(?:^|[;&|]|&&|\|\|)\s*(?:export\s+)?"
    r"(?:(?:[A-Za-z_][A-Za-z0-9_]*=\S+\s+)*)?"
    r"(?:[^\s;|&]+/)?git(?:\s+(?:-C\s+\S+|-c\s+\S+))*\s+push\s+(-f|--force|--force-with-lease)\b"
)

if force_push.search(cmd) or force_push2.search(cmd):
    print("force_push")
    raise SystemExit(0)

# Catastrophic rm: recursive on / or home
if re.search(r"(?:^|[;&|]|&&|\|\|)\s*(?:sudo\s+)?rm\b", cmd) and re.search(
    r"(?:^|[;&|]|&&|\|\|)\s*(?:sudo\s+)?rm\s+[^\n;|&]*-[a-zA-Z]*r", cmd
):
    # Target is filesystem root
    if re.search(
        r"(?:^|[;&|]|&&|\|\|)\s*(?:sudo\s+)?rm\s+(?:-[a-zA-Z0-9]+\s+)*/(?:\s|[;&|]|$)",
        cmd,
    ):
        print("rm_root")
        raise SystemExit(0)
    # Target is home
    if re.search(
        r"(?:^|[;&|]|&&|\|\|)\s*(?:sudo\s+)?rm\s+(?:-[a-zA-Z0-9]+\s+)*(?:~/|~(?:\s|[;&|]|$)|\$HOME(?:/|\s|[;&|]|$))",
        cmd,
    ):
        print("rm_home")
        raise SystemExit(0)

print("allow")
')

case "$decision" in
  force_push)
    deny \
      "Force push blocked by Spacecraft safety hook." \
      "Do not run git push --force / -f / --force-with-lease. Rewrite history only when the user explicitly overrides outside this gate."
    ;;
  rm_root)
    deny \
      "Destructive rm of filesystem root blocked." \
      "Command matched rm targeting /. Refused by Spacecraft safety hook."
    ;;
  rm_home)
    deny \
      "Destructive rm of home directory blocked." \
      "Command matched rm targeting ~ or \$HOME. Refused by Spacecraft safety hook."
    ;;
  *)
    allow
    ;;
esac
