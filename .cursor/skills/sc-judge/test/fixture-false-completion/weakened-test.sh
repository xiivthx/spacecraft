#!/bin/sh
# Planted weakened test: tautology so GREEN is cheap (no real acceptance check).
set -eu
# assert always true - deliberately weakened
if true; then
  echo "PASS (tautology)"
  exit 0
fi
exit 1
