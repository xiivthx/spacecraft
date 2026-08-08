#!/bin/sh
# gen-user-rules.sh - concatenate the bodies of the alwaysApply rules
# into a single output file, stripping YAML frontmatter from each source.
#
# Usage: sh scripts/gen-user-rules.sh <rules-dir> <out-file>
set -e

RULES_DIR="${1:?usage: gen-user-rules.sh <rules-dir> <out-file>}"
OUT="${2:?usage: gen-user-rules.sh <rules-dir> <out-file>}"

SOURCES="000-spacecraft.mdc 025-english-coach.mdc 026-intent-coach.mdc 027-th-en-hil.mdc 050-style.mdc 100-conventions.mdc 200-workflow.mdc"

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
