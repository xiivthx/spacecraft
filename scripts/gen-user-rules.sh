#!/bin/sh
# gen-user-rules.sh - emit a short Spacecraft CORE User Rules file from
# 010-hard-contract.mdc (+ optional HIL one-liner from 027), stripping YAML
# frontmatter. Keep under ~40 lines so User Rules stay readable.
#
# Usage: sh scripts/gen-user-rules.sh <rules-dir> <out-file>
set -e

RULES_DIR="${1:?usage: gen-user-rules.sh <rules-dir> <out-file>}"
OUT="${2:?usage: gen-user-rules.sh <rules-dir> <out-file>}"

SOURCES="010-hard-contract.mdc"

mkdir -p "$(dirname "$OUT")"
: > "$OUT"

strip_frontmatter() {
  awk '
    NR == 1 && $0 == "---" { infm = 1; next }
    infm && $0 == "---"    { infm = 0; next }
    infm                   { next }
    { print }
  ' "$1"
}

for name in $SOURCES; do
  src="$RULES_DIR/$name"
  if [ ! -f "$src" ]; then
    echo "FAIL: missing rule source: $src" >&2
    exit 1
  fi
  strip_frontmatter "$src" >> "$OUT"
  echo "" >> "$OUT"
done

# Tiny HIL reminder from 027 without pulling the full language rule.
if [ -f "$RULES_DIR/027-th-en-hil.mdc" ]; then
  printf '%s\n' "## HIL language (CORE)" >> "$OUT"
  printf '%s\n' "English for technical substance; Thai for HIL questions, short status, and handoffs. No dual language blocks." >> "$OUT"
  echo "" >> "$OUT"
fi
