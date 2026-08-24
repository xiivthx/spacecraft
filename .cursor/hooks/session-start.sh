#!/bin/sh
# sessionStart hook: inject active mission context (always succeeds).

if [ -x ./spacecraft ]; then
  ./spacecraft context 2>&1
elif command -v spacecraft >/dev/null 2>&1; then
  spacecraft context 2>&1
else
  printf '%s\n' "No active spacecraft mission."
fi
exit 0
