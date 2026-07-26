#!/usr/bin/env python3
"""Filter a project hooks.json down to a safety-hooks subset with rewritten
absolute command paths, for global (non-project) installs.

Usage:
  rewrite-global-hooks.py <source-hooks.json> <global-hooks-dir> <out.json> \
      <script-basename> [<script-basename> ...]

Keeps only hook entries whose command basename is in the given safety list
(e.g. never session-start.sh, which is project-layer only), and rewrites
their "command" to an absolute path under <global-hooks-dir> so the hook
still works outside this repo.
"""
import json
import os
import sys


def main():
    if len(sys.argv) < 5:
        sys.stderr.write(__doc__)
        return 2

    src_path, hooks_dir, out_path = sys.argv[1], sys.argv[2], sys.argv[3]
    safety = set(sys.argv[4:])

    with open(src_path, "r", encoding="utf-8") as f:
        source = json.load(f)

    filtered = {"version": source.get("version", 1), "hooks": {}}
    for event, entries in source.get("hooks", {}).items():
        kept = []
        for entry in entries:
            base = os.path.basename(entry.get("command", ""))
            if base not in safety:
                continue
            rewritten = dict(entry)
            rewritten["command"] = os.path.join(hooks_dir, base)
            kept.append(rewritten)
        if kept:
            filtered["hooks"][event] = kept

    with open(out_path, "w", encoding="utf-8") as f:
        json.dump(filtered, f, indent=2)
        f.write("\n")
    return 0


if __name__ == "__main__":
    sys.exit(main())
