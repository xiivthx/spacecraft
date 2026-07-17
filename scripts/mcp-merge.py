#!/usr/bin/env python3
"""Merge or unmerge MCP servers without clobbering unrelated user servers.

Usage:
  mcp-merge.py merge   <target-mcp.json> <source-mcp.json>
  mcp-merge.py unmerge <target-mcp.json> <source-mcp.json>

merge:   add/update every server defined in source into target.
unmerge: remove from target only the servers that source defines,
         leaving all other user servers intact. Empty target files are
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
    src_servers = source.get("mcpServers", {})
    tgt_servers = target.get("mcpServers", {})

    if mode == "merge":
        for name, cfg in src_servers.items():
            tgt_servers[name] = cfg
        target["mcpServers"] = tgt_servers
        write(target_path, target)
        print(f"merged {len(src_servers)} server(s) into {target_path}")
        return 0

    # unmerge
    removed = 0
    for name in src_servers:
        if name in tgt_servers:
            del tgt_servers[name]
            removed += 1
    if tgt_servers:
        target["mcpServers"] = tgt_servers
        write(target_path, target)
    elif os.path.exists(target_path) and list(target.keys()) in ([], ["mcpServers"]):
        os.remove(target_path)
    else:
        target["mcpServers"] = tgt_servers
        write(target_path, target)
    print(f"removed {removed} server(s) from {target_path}")
    return 0


if __name__ == "__main__":
    sys.exit(main())
