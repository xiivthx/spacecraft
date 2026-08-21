#!/bin/sh
# Assert deny/allow paths for check-main-write.sh and check-ship-commands.sh.
set -eu

ROOT=$(CDPATH= cd -- "$(dirname "$0")/../.." && pwd)
MAIN_WRITE="$ROOT/.cursor/hooks/check-main-write.sh"
SHIP="$ROOT/.cursor/hooks/check-ship-commands.sh"
FAIL=0

assert_perm() {
  label="$1"
  expected="$2"
  actual="$3"
  if [ "$actual" = "$expected" ]; then
    printf 'PASS %s (permission=%s)\n' "$label" "$actual"
  else
    printf 'FAIL %s: expected permission=%s got=%s\n' "$label" "$expected" "$actual"
    FAIL=1
  fi
}

extract_perm() {
  printf '%s' "$1" | python3 -c 'import json,sys; print(json.load(sys.stdin).get("permission",""))'
}

run_main() {
  branch="$1"
  cmd_json="$2"
  SPACECRAFT_HOOK_BRANCH_OVERRIDE="$branch" "$MAIN_WRITE" <<EOF
$cmd_json
EOF
}

# run_ship <SPACECRAFT_SHIP|empty> <cmd_json> [SPACECRAFT_CLOSEOUT_CMD] [SPACECRAFT_QUICK]
run_ship() {
  env_ship="$1"
  cmd_json="$2"
  closeout_cmd="${3-}"
  env_quick="${4-}"
  if [ -n "$env_ship" ]; then
    if [ -n "$closeout_cmd" ]; then
      if [ -n "$env_quick" ]; then
        SPACECRAFT_SHIP="$env_ship" SPACECRAFT_QUICK="$env_quick" SPACECRAFT_CLOSEOUT_CMD="$closeout_cmd" "$SHIP" <<EOF
$cmd_json
EOF
      else
        env -u SPACECRAFT_QUICK SPACECRAFT_SHIP="$env_ship" SPACECRAFT_CLOSEOUT_CMD="$closeout_cmd" "$SHIP" <<EOF
$cmd_json
EOF
      fi
    else
      if [ -n "$env_quick" ]; then
        env -u SPACECRAFT_CLOSEOUT_CMD SPACECRAFT_SHIP="$env_ship" SPACECRAFT_QUICK="$env_quick" "$SHIP" <<EOF
$cmd_json
EOF
      else
        env -u SPACECRAFT_CLOSEOUT_CMD -u SPACECRAFT_QUICK SPACECRAFT_SHIP="$env_ship" "$SHIP" <<EOF
$cmd_json
EOF
      fi
    fi
  else
    env -u SPACECRAFT_SHIP -u SPACECRAFT_QUICK -u SPACECRAFT_CLOSEOUT_CMD "$SHIP" <<EOF
$cmd_json
EOF
  fi
}

# --- check-main-write.sh ---
out=$(run_main "feat/x" '{"command":"git commit -m x"}')
assert_perm "main-write: feat allows commit" "allow" "$(extract_perm "$out")"

out=$(run_main "main" '{"command":"git commit -m x"}')
assert_perm "main-write: main denies commit" "deny" "$(extract_perm "$out")"

out=$(run_main "master" '{"command":"git merge feat/x"}')
assert_perm "main-write: master denies merge" "deny" "$(extract_perm "$out")"

out=$(run_main "main" '{"command":"git status"}')
assert_perm "main-write: main allows status" "allow" "$(extract_perm "$out")"

out=$(run_main "main" '{"command":"git log -1"}')
assert_perm "main-write: main allows log" "allow" "$(extract_perm "$out")"

out=$(run_main "main" '{"command":"git push origin main"}')
assert_perm "main-write: main denies push" "deny" "$(extract_perm "$out")"

out=$(run_main "main" 'not-json')
assert_perm "main-write: bad json denies" "deny" "$(extract_perm "$out")"

# Without OVERRIDE: resolve branch via workspace_roots (Cursor hook cwd is often not the repo).
out=$(
  env -u SPACECRAFT_HOOK_BRANCH_OVERRIDE "$MAIN_WRITE" <<EOF
{"command":"git commit -m x","workspace_roots":["$ROOT"]}
EOF
)
# ROOT may be on main or a feature branch; only assert allow when not on main/master.
cur=$(git -C "$ROOT" branch --show-current 2>/dev/null || true)
case "$cur" in
  main|master)
    assert_perm "main-write: workspace_roots on main denies commit" "deny" "$(extract_perm "$out")"
    ;;
  *)
    assert_perm "main-write: workspace_roots feat allows commit" "allow" "$(extract_perm "$out")"
    ;;
esac

# --- check-ship-commands.sh ---
out=$(run_ship "" '{"command":"git push"}')
assert_perm "ship: without env denies push" "deny" "$(extract_perm "$out")"

out=$(run_ship "" '{"command":"git merge --no-ff feat/x"}')
assert_perm "ship: without env denies merge" "deny" "$(extract_perm "$out")"

out=$(run_ship "" '{"command":"git tag v1.0.0"}')
assert_perm "ship: without env denies tag" "deny" "$(extract_perm "$out")"

# AUTH (quoted user authorization) is a separate always-on gate; it does not set
# SPACECRAFT_SHIP or allow merge/push/tag. Env gate still required.
out=$(run_ship "" '{"command":"git push origin main"}')
assert_perm "ship: AUTH does not bypass SPACECRAFT_SHIP - denies push" "deny" "$(extract_perm "$out")"
if ! printf '%s' "$out" | grep -Eq 'deny|Ship gate'; then
  printf 'FAIL ship: AUTH no-bypass message missing deny|Ship gate\n'
  FAIL=1
else
  printf 'PASS ship: AUTH no-bypass message contains deny|Ship gate\n'
fi

out=$(run_ship "1" '{"command":"git push"}' "true")
assert_perm "ship: SPACECRAFT_SHIP=1 closeout ok asks push" "ask" "$(extract_perm "$out")"

out=$(run_ship "1" '{"command":"git tag v1.0.0"}' "exit 0")
assert_perm "ship: SPACECRAFT_SHIP=1 closeout ok allows tag" "allow" "$(extract_perm "$out")"

out=$(run_ship "1" '{"command":"git push"}' "false")
assert_perm "ship: SPACECRAFT_SHIP=1 closeout fail denies push" "deny" "$(extract_perm "$out")"

out=$(run_ship "1" '{"command":"git merge --no-ff feat/x"}' "exit 1")
assert_perm "ship: SPACECRAFT_SHIP=1 closeout fail denies merge" "deny" "$(extract_perm "$out")"

# /sc-quick: SPACECRAFT_QUICK=1 skips closeout; push still asks; merge/tag allow.
out=$(run_ship "1" '{"command":"git push"}' "false" "1")
assert_perm "ship: SPACECRAFT_QUICK=1 skips closeout asks push" "ask" "$(extract_perm "$out")"

out=$(run_ship "1" '{"command":"git merge --no-ff docs/x"}' "exit 1" "1")
assert_perm "ship: SPACECRAFT_QUICK=1 skips closeout allows merge" "allow" "$(extract_perm "$out")"

out=$(run_ship "1" '{"command":"git tag v1.0.0"}' "false" "1")
assert_perm "ship: SPACECRAFT_QUICK=1 skips closeout allows tag" "allow" "$(extract_perm "$out")"

# QUICK alone is not enough — still need SPACECRAFT_SHIP=1.
out=$(env -u SPACECRAFT_SHIP SPACECRAFT_QUICK=1 env -u SPACECRAFT_CLOSEOUT_CMD "$SHIP" <<'EOF'
{"command":"git push"}
EOF
)
assert_perm "ship: SPACECRAFT_QUICK=1 without SHIP denies push" "deny" "$(extract_perm "$out")"

# Command-string assignments (Cursor hook process has no agent shell env).
out=$(env -u SPACECRAFT_SHIP -u SPACECRAFT_QUICK -u SPACECRAFT_CLOSEOUT_CMD "$SHIP" <<'EOF'
{"command":"SPACECRAFT_SHIP=1 SPACECRAFT_QUICK=1 git merge --no-ff docs/x"}
EOF
)
assert_perm "ship: command-prefix SHIP+QUICK allows merge" "allow" "$(extract_perm "$out")"

out=$(env -u SPACECRAFT_SHIP -u SPACECRAFT_QUICK SPACECRAFT_CLOSEOUT_CMD="true" "$SHIP" <<'EOF'
{"command":"SPACECRAFT_SHIP=1 git tag -a v1.0.0 -m v1.0.0"}
EOF
)
assert_perm "ship: command-prefix SHIP allows tag with closeout" "allow" "$(extract_perm "$out")"

out=$(env -u SPACECRAFT_SHIP -u SPACECRAFT_QUICK -u SPACECRAFT_CLOSEOUT_CMD "$SHIP" <<'EOF'
{"command":"echo SPACECRAFT_SHIP=1; git push"}
EOF
)
assert_perm "ship: echo of SHIP=1 does not unlock push" "deny" "$(extract_perm "$out")"

out=$(run_main "main" '{"command":"SPACECRAFT_SHIP=1 git merge --no-ff feat/x"}')
assert_perm "main-write: ship-prefixed merge on main allows" "allow" "$(extract_perm "$out")"

out=$(run_main "main" '{"command":"SPACECRAFT_SHIP=1 git tag -a v1.0.0 -m v1.0.0"}')
assert_perm "main-write: ship-prefixed tag on main allows" "allow" "$(extract_perm "$out")"

out=$(run_main "main" '{"command":"SPACECRAFT_SHIP=1 git commit -m x"}')
assert_perm "main-write: ship-prefixed commit on main still denies" "deny" "$(extract_perm "$out")"

out=$(run_ship "" '{"command":"ls"}')
assert_perm "ship: allows unrelated ls" "allow" "$(extract_perm "$out")"

out=$(run_ship "" '{"command":"git status"}')
assert_perm "ship: allows git status" "allow" "$(extract_perm "$out")"

out=$(run_ship "" '{"command":"bash .cursor/hooks/hooks_test.sh"}')
assert_perm "ship: allows hooks_test.sh" "allow" "$(extract_perm "$out")"

out=$(run_ship "" '{"command":"printf %s \"{\\\"command\\\":\\\"git push\\\"}\" | .cursor/hooks/check-ship-commands.sh"}')
assert_perm "ship: allows check-ship-commands.sh self-test" "allow" "$(extract_perm "$out")"

out=$(run_ship "" '{"command":"git push origin main; bash .cursor/hooks/hooks_test.sh"}')
assert_perm "ship: compound push+hooks_test without env denies" "deny" "$(extract_perm "$out")"

out=$(run_ship "" 'not-json')
assert_perm "ship: bad json denies" "deny" "$(extract_perm "$out")"

# --- check-main-write.sh hook-test allowlist ---
out=$(run_main "main" '{"command":"bash .cursor/hooks/hooks_test.sh"}')
assert_perm "main-write: allows hooks_test.sh on main" "allow" "$(extract_perm "$out")"

out=$(run_main "main" '{"command":"printf %s \"{\\\"command\\\":\\\"git commit -m x\\\"}\" | .cursor/hooks/check-main-write.sh"}')
assert_perm "main-write: allows check-main-write.sh self-test on main" "allow" "$(extract_perm "$out")"

out=$(run_main "main" '{"command":"git push; bash .cursor/hooks/hooks_test.sh"}')
assert_perm "main-write: compound push+hooks_test denies" "deny" "$(extract_perm "$out")"

out=$(run_main "" '{"command":"git commit -m x"}')
assert_perm "main-write: empty branch denies commit" "deny" "$(extract_perm "$out")"

out=$(run_main "main" '{"command":"env git commit -m x"}')
assert_perm "main-write: env git commit denies" "deny" "$(extract_perm "$out")"

out=$(run_main "main" '{"command":"git -C . commit -m x"}')
assert_perm "main-write: git -C commit denies" "deny" "$(extract_perm "$out")"

# --- block-secrets-read.sh ---
SECRETS="$ROOT/.cursor/hooks/block-secrets-read.sh"
run_secrets() {
  printf '%s' "$1" | "$SECRETS"
}

out=$(run_secrets '{"file_path":"/proj/.env"}')
assert_perm "secrets: denies .env" "deny" "$(extract_perm "$out")"

out=$(run_secrets '{"file_path":"/proj/.env.local"}')
assert_perm "secrets: denies .env.local" "deny" "$(extract_perm "$out")"

out=$(run_secrets '{"file_path":"/proj/.env.example"}')
assert_perm "secrets: allows .env.example" "allow" "$(extract_perm "$out")"

out=$(run_secrets '{"file_path":"/proj/config.env.sample"}')
assert_perm "secrets: allows .env.sample" "allow" "$(extract_perm "$out")"

out=$(run_secrets '{"file_path":"/proj/id_rsa"}')
assert_perm "secrets: denies id_rsa" "deny" "$(extract_perm "$out")"

out=$(run_secrets '{"file_path":"/proj/cert.pem"}')
assert_perm "secrets: denies .pem" "deny" "$(extract_perm "$out")"

out=$(run_secrets '{"file_path":"/proj/src/app.ts"}')
assert_perm "secrets: allows source file" "allow" "$(extract_perm "$out")"

# --- block-destructive.sh ---
DEST="$ROOT/.cursor/hooks/block-destructive.sh"
run_dest() {
  printf '%s' "$1" | "$DEST"
}

out=$(run_dest '{"command":"git push --force origin main"}')
assert_perm "destructive: denies force push" "deny" "$(extract_perm "$out")"

out=$(run_dest '{"command":"git push -f origin main"}')
assert_perm "destructive: denies push -f" "deny" "$(extract_perm "$out")"

out=$(run_dest '{"command":"SPACECRAFT_SHIP=1 git push origin main"}')
assert_perm "destructive: allows non-force push" "allow" "$(extract_perm "$out")"

out=$(run_dest '{"command":"rm -rf /"}')
assert_perm "destructive: denies rm -rf /" "deny" "$(extract_perm "$out")"

out=$(run_dest '{"command":"rm -rf ~"}')
assert_perm "destructive: denies rm -rf home" "deny" "$(extract_perm "$out")"

out=$(run_dest '{"command":"rm -rf ./build"}')
assert_perm "destructive: allows rm project dir" "allow" "$(extract_perm "$out")"

if [ "$FAIL" -ne 0 ]; then
  printf 'hooks_test.sh: FAILED\n'
  exit 1
fi
printf 'hooks_test.sh: all assertions passed\n'
exit 0
