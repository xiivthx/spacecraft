#!/usr/bin/env python3
"""Merge or unmerge Cursor hooks without clobbering unrelated user hooks.

Usage:
  hooks-merge.py merge   <target-hooks.json> <source-hooks.json>
  hooks-merge.py unmerge <target-hooks.json> <source-hooks.json>

merge:   add every hook entry defined in source into the matching event
         array in target, skipping entries already present. Unrelated
         events and entries are left intact. Schema version is preserved
         (defaults to 1).
unmerge: remove from target only the hook entries that source defines,
         leaving all other user hooks intact. Empty target files are
         removed to avoid leaving dangling config behind.
"""
import json
import os
import sys


def load(path):
    if not os.path.exists(path):
        return {}
    with open(path, "r", encoding="utf-8") as f:
        text = f.read().strip()
    if not text:
        return {}
    data = json.loads(text)
    if not isinstance(data, dict):
        raise ValueError(f"{path}: expected a JSON object")
    return data


def write(path, data):
    os.makedirs(os.path.dirname(os.path.abspath(path)), exist_ok=True)
    with open(path, "w", encoding="utf-8") as f:
        json.dump(data, f, indent=2)
        f.write("\n")


def main():
    if len(sys.argv) != 4 or sys.argv[1] not in ("merge", "unmerge"):
        sys.stderr.write(__doc__)
        return 2

    mode, target_path, source_path = sys.argv[1], sys.argv[2], sys.argv[3]
    source = load(source_path)
    target = load(target_path)
    src_hooks = source.get("hooks", {})
    tgt_hooks = target.get("hooks", {})

    if mode == "merge":
        added = 0
        for event, entries in src_hooks.items():
            existing = tgt_hooks.setdefault(event, [])
            for entry in entries:
                if entry not in existing:
                    existing.append(entry)
                    added += 1
        target["hooks"] = tgt_hooks
        target.setdefault("version", source.get("version", 1))
        write(target_path, target)
        print(f"merged {added} hook(s) into {target_path}")
        return 0

    # unmerge
    removed = 0
    for event, entries in src_hooks.items():
        existing = tgt_hooks.get(event)
        if not existing:
            continue
        kept = [e for e in existing if e not in entries]
        removed += len(existing) - len(kept)
        if kept:
            tgt_hooks[event] = kept
        else:
            del tgt_hooks[event]

    if tgt_hooks:
        target["hooks"] = tgt_hooks
        write(target_path, target)
    elif os.path.exists(target_path) and set(target.keys()) <= {"hooks", "version"}:
        os.remove(target_path)
    else:
        target["hooks"] = tgt_hooks
        write(target_path, target)
    print(f"removed {removed} hook(s) from {target_path}")
    return 0


if __name__ == "__main__":
    sys.exit(main())
