#!/bin/sh
# Deny reading secret files (.env, keys, credentials). Allow .env.example / *.env.sample.

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

file_path=$(printf '%s' "$input" | python3 -c '
import json, sys
raw = sys.stdin.read()
try:
    data = json.loads(raw)
except Exception:
    sys.exit(2)
if not isinstance(data, dict):
    sys.exit(2)
# Cursor beforeReadFile: file_path; some payloads use path / filePath
path = data.get("file_path") or data.get("filePath") or data.get("path") or ""
if not isinstance(path, str):
    sys.exit(2)
print(path)
' 2>/dev/null) || deny \
  "hook could not parse file path" \
  "The secrets-read hook could not parse file_path from stdin JSON. Fix the hook input or retry."

if [ -z "$file_path" ]; then
  allow
fi

decision=$(printf '%s' "$file_path" | python3 -c '
import os, re, sys
path = sys.stdin.read().strip()
base = os.path.basename(path)
# Allow documented samples
if re.search(r"(?i)^\.env\.example$", base):
    print("allow")
    raise SystemExit(0)
if re.search(r"(?i)\.env\.sample$", base):
    print("allow")
    raise SystemExit(0)
if re.search(r"(?i)^\.env$", base):
    print("deny")
    raise SystemExit(0)
if re.search(r"(?i)^\.env\.", base):
    print("deny")
    raise SystemExit(0)
# credentials / keys
if re.search(r"(?i)(^|/)(credentials|credentials\.json|id_rsa|id_ed25519)(\.|$)", path):
    print("deny")
    raise SystemExit(0)
if re.search(r"(?i)\.(pem|p12|pfx|key)$", base) and not re.search(r"(?i)\.(pub)$", base):
    # allow *.pub; deny private key-like extensions
    if re.search(r"(?i)\.pub$", base):
        print("allow")
    else:
        print("deny")
    raise SystemExit(0)
print("allow")
')

case "$decision" in
  deny)
    deny \
      "Secret file read blocked. Do not read .env, credentials, or private keys." \
      "Hook denied reading a secret path ($file_path). Use .env.example or ask the user; never load secrets into context."
    ;;
  *)
    allow
    ;;
esac
