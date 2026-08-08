#!/bin/sh
# check-judge-break.sh - prove closeout/ready gates reject known-bad fixtures.
#
# Loads each pack under .cursor/skills/sc-judge/references/judge-break/*/
# into a temp .space, runs spacecraft closeout-check, and asserts non-zero
# exit plus the expected failure substring. No LLM calls.
#
# Usage: check-judge-break.sh [repo-root] [spacecraft-binary]
set -e

ROOT="${1:-.}"
BIN="${2:-}"

# Resolve before run_fixture cds into temp (relative BIN would break).
ROOT="$(cd "$ROOT" && pwd)"
FIXDIR="$ROOT/.cursor/skills/sc-judge/references/judge-break"

if [ ! -d "$FIXDIR" ]; then
  echo "FAIL: missing $FIXDIR"
  exit 1
fi

if [ -z "$BIN" ]; then
  if [ -x "$ROOT/spacecraft" ]; then
    BIN="$ROOT/spacecraft"
  elif command -v spacecraft >/dev/null 2>&1; then
    BIN="$(command -v spacecraft)"
  else
    echo "FAIL: spacecraft binary required (pass as arg2 or build ./spacecraft)"
    exit 1
  fi
fi

case "$BIN" in
  /*) ;;
  *) BIN="$(cd "$(dirname "$BIN")" && pwd)/$(basename "$BIN")" ;;
esac

if [ ! -x "$BIN" ]; then
  echo "FAIL: not executable: $BIN"
  exit 1
fi

read_expect() {
  file="$1"
  field="$2"
  python3 -c "
import json,sys
d=json.load(open(sys.argv[1]))
print(d.get(sys.argv[2],''))
" "$file" "$field"
}

run_fixture() {
  name="$1"
  dir="$FIXDIR/$name"
  expect="$dir/expect.json"
  if [ ! -f "$expect" ]; then
    echo "FAIL: $name missing expect.json"
    return 1
  fi

  mid="$(read_expect "$expect" id)"
  must="$(read_expect "$expect" mustContain)"
  if [ -z "$mid" ] || [ -z "$must" ]; then
    echo "FAIL: $name expect.json missing id or mustContain"
    return 1
  fi

  mkdir -p "$ROOT/.tmp"
  tmp=$(mktemp -d "$ROOT/.tmp/judge-break.XXXXXX")

  space="$tmp/.space"
  mdir="$space/missions/$mid"
  mkdir -p "$mdir"
  for f in mission.json spec.md plan.json evidence.jsonl review.json; do
    if [ -f "$dir/$f" ]; then
      cp "$dir/$f" "$mdir/$f"
    fi
  done
  printf '%s\n' "$mid" > "$space/current"

  if [ "$name" = "false-completion" ]; then
    if ! grep -Eq '"status"[[:space:]]*:[[:space:]]*"done"' "$mdir/plan.json"; then
      echo "FAIL: $name plan.json must mark a task done"
      rm -rf "$tmp"
      return 1
    fi
    if [ -s "$mdir/evidence.jsonl" ]; then
      echo "FAIL: $name evidence.jsonl must be empty"
      rm -rf "$tmp"
      return 1
    fi
  fi

  set +e
  out=$(
    cd "$tmp" &&
      SPACECRAFT_CLOSEOUT_SKIP_CHANGELOG=1 "$BIN" closeout-check 2>&1
  )
  code=$?
  set -e

  rm -rf "$tmp"

  if [ "$code" -eq 0 ]; then
    echo "FAIL: $name closeout-check unexpectedly passed"
    echo "$out"
    return 1
  fi
  if ! printf '%s\n' "$out" | grep -qi "$must"; then
    echo "FAIL: $name expected failure containing: $must"
    echo "$out"
    return 1
  fi
  echo "ok: $name rejected ($must)"
}

count=0
for expect in "$FIXDIR"/*/expect.json; do
  [ -f "$expect" ] || continue
  name=$(basename "$(dirname "$expect")")
  run_fixture "$name"
  count=$((count + 1))
done

if [ "$count" -eq 0 ]; then
  echo "FAIL: no fixtures under $FIXDIR"
  exit 1
fi

echo "ok: judge-break $count fixture(s) rejected as expected"
exit 0
