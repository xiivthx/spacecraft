#!/bin/sh
# sessionStart hook: inject active mission context (always succeeds).

if [ -x ./spacecraft ]; then
  ./spacecraft status 2>&1 || printf '%s\n' "No active spacecraft mission."
elif command -v spacecraft >/dev/null 2>&1; then
  spacecraft status 2>&1 || printf '%s\n' "No active spacecraft mission."
else
  printf '%s\n' "No active spacecraft mission."
fi
exit 0
