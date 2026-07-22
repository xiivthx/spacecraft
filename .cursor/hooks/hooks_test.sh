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
assert_perm "ship: SPACECRAFT_SHIP=1 closeout ok allows push" "allow" "$(extract_perm "$out")"

out=$(run_ship "1" '{"command":"git tag v1.0.0"}' "exit 0")
assert_perm "ship: SPACECRAFT_SHIP=1 closeout ok allows tag" "allow" "$(extract_perm "$out")"

out=$(run_ship "1" '{"command":"git push"}' "false")
assert_perm "ship: SPACECRAFT_SHIP=1 closeout fail denies push" "deny" "$(extract_perm "$out")"

out=$(run_ship "1" '{"command":"git merge --no-ff feat/x"}' "exit 1")
assert_perm "ship: SPACECRAFT_SHIP=1 closeout fail denies merge" "deny" "$(extract_perm "$out")"

# /sc-quick: SPACECRAFT_QUICK=1 skips closeout even when closeout would fail.
out=$(run_ship "1" '{"command":"git push"}' "false" "1")
assert_perm "ship: SPACECRAFT_QUICK=1 skips closeout allows push" "allow" "$(extract_perm "$out")"

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

if [ "$FAIL" -ne 0 ]; then
  printf 'hooks_test.sh: FAILED\n'
  exit 1
fi
printf 'hooks_test.sh: all assertions passed\n'
exit 0
