#!/bin/sh
# harness-scorecard.sh - thin fail-closed wrapper over required harness fixtures.
#
# Prints one greppable line per required dimension:
#   SCORECARD <dimension-id> <pass|fail>
# Exit 0 iff every required dimension is pass.
# Env HARNESS_SCORECARD_FORCE_FAIL=<dimension-id> marks that dimension fail and
# forces aggregate exit non-zero (negative prove without mutating fixtures).
#
# Usage: harness-scorecard.sh <repo-root> [spacecraft-binary]
set -u

ROOT="${1:?usage: harness-scorecard.sh <repo-root> [spacecraft-binary]}"
BIN="${2:-}"

ROOT="$(cd "$ROOT" && pwd)"

if [ -z "$BIN" ]; then
  if [ -x "$ROOT/spacecraft" ]; then
    BIN="$ROOT/spacecraft"
  elif command -v spacecraft >/dev/null 2>&1; then
    BIN="$(command -v spacecraft)"
  fi
fi
if [ -n "$BIN" ]; then
  case "$BIN" in
    /*) ;;
    *) BIN="$(cd "$(dirname "$BIN")" && pwd)/$(basename "$BIN")" ;;
  esac
fi

FORCE_FAIL="${HARNESS_SCORECARD_FORCE_FAIL:-}"
aggregate=0

emit() {
  dim=$1
  status=$2
  if [ -n "$FORCE_FAIL" ] && [ "$FORCE_FAIL" = "$dim" ]; then
    status=fail
  fi
  printf 'SCORECARD %s %s\n' "$dim" "$status"
  if [ "$status" != pass ]; then
    aggregate=1
  fi
}

# --- install-smoke: scripts/smoke.sh <ROOT> [BIN] ---
rc=0
if [ -n "$BIN" ]; then
  sh "$ROOT/scripts/smoke.sh" "$ROOT" "$BIN" || rc=$?
else
  sh "$ROOT/scripts/smoke.sh" "$ROOT" || rc=$?
fi
if [ "$rc" -eq 0 ]; then
  emit install-smoke pass
else
  emit install-smoke fail
fi

# --- false-completion: scripts/check-judge-break.sh <ROOT> [BIN] ---
rc=0
if [ -n "$BIN" ]; then
  sh "$ROOT/scripts/check-judge-break.sh" "$ROOT" "$BIN" || rc=$?
else
  sh "$ROOT/scripts/check-judge-break.sh" "$ROOT" || rc=$?
fi
if [ "$rc" -eq 0 ]; then
  emit false-completion pass
else
  emit false-completion fail
fi

# --- judge-skill: all .cursor/skills/sc-judge/test/*smoke*.sh (no ROOT/BIN argv) ---
# Most smokes: detection/pass success = exit 0.
# judge-smoke.sh (false-completion plant): detection success = nonzero + VERDICT: REFUTED.
judge_dir="$ROOT/.cursor/skills/sc-judge/test"
judge_rc=0
judge_count=0
mkdir -p "$ROOT/.tmp"
for f in "$judge_dir"/*smoke*.sh; do
  [ -f "$f" ] || continue
  judge_count=$((judge_count + 1))
  smoke_out=$(mktemp "$ROOT/.tmp/judge-skill-smoke.XXXXXX")
  smoke_rc=0
  sh "$f" >"$smoke_out" 2>&1 || smoke_rc=$?
  cat "$smoke_out"
  if [ "$smoke_rc" -eq 0 ]; then
    :
  elif [ "$(basename "$f")" = "judge-smoke.sh" ] && grep -Fq 'VERDICT: REFUTED' "$smoke_out"; then
    :
  else
    judge_rc=1
  fi
  rm -f "$smoke_out"
done
if [ "$judge_count" -eq 0 ]; then
  judge_rc=1
fi
if [ "$judge_rc" -eq 0 ]; then
  emit judge-skill pass
else
  emit judge-skill fail
fi

# --- process-grammar: check-workflow / check-sc-planning / check-sc-planner ---
proc_rc=0
proc_count=0
for f in \
  "$ROOT"/scripts/check-workflow-*.sh \
  "$ROOT"/scripts/check-sc-planning-*.sh \
  "$ROOT"/scripts/check-sc-planner-*.sh
do
  [ -f "$f" ] || continue
  proc_count=$((proc_count + 1))
  sh "$f" "$ROOT" || proc_rc=1
done
if [ "$proc_count" -eq 0 ]; then
  proc_rc=1
fi
if [ "$proc_rc" -eq 0 ]; then
  emit process-grammar pass
else
  emit process-grammar fail
fi

exit "$aggregate"
