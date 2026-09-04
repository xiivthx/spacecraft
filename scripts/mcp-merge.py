#!/usr/bin/env python3
"""Merge or unmerge MCP servers without clobbering unrelated user servers.

Usage:
  mcp-merge.py merge   <target-mcp.json> <source-mcp.json> [--strip-pack-mcp <cursor-dir>]
  mcp-merge.py unmerge <target-mcp.json> <source-mcp.json>

merge:   add/update every server defined in source into target.
         With --strip-pack-mcp: drop source servers whose names appear in any
         selectable pack MCP fragment under <cursor-dir> (catalog `mcp` fields,
         else all *.json under mcp-packs/). Pack MCP is applied separately by
         project install when those packs are selected.
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


def pack_managed_server_names(cursor_dir):
    """Server names owned by selectable pack MCP fragments under cursor_dir."""
    names = set()
    catalog_path = os.path.join(cursor_dir, "spacecraft-packs.json")
    if os.path.isfile(catalog_path):
        catalog = load(catalog_path)
        for entry in catalog.get("packs", []):
            if not isinstance(entry, dict):
                continue
            if entry.get("status") != "selectable":
                continue
            mcp_rel = entry.get("mcp")
            if not mcp_rel or not isinstance(mcp_rel, str):
                continue
            frag_path = os.path.join(cursor_dir, mcp_rel)
            frag = load(frag_path)
            servers = frag.get("mcpServers") or {}
            if isinstance(servers, dict):
                names.update(servers.keys())
        return names

    packs_dir = os.path.join(cursor_dir, "mcp-packs")
    if not os.path.isdir(packs_dir):
        return names
    for fname in os.listdir(packs_dir):
        if not fname.endswith(".json"):
            continue
        frag = load(os.path.join(packs_dir, fname))
        servers = frag.get("mcpServers") or {}
        if isinstance(servers, dict):
            names.update(servers.keys())
    return names


def parse_args(argv):
    if len(argv) < 4:
        return None
    mode = argv[1]
    if mode not in ("merge", "unmerge"):
        return None
    target_path, source_path = argv[2], argv[3]
    strip_cursor = None
    rest = argv[4:]
    i = 0
    while i < len(rest):
        if rest[i] == "--strip-pack-mcp":
            if i + 1 >= len(rest):
                return None
            strip_cursor = rest[i + 1]
            i += 2
            continue
        return None
    if mode == "unmerge" and strip_cursor is not None:
        return None
    return mode, target_path, source_path, strip_cursor


def main():
    parsed = parse_args(sys.argv)
    if parsed is None:
        sys.stderr.write(__doc__)
        return 2

    mode, target_path, source_path, strip_cursor = parsed
    source = load(source_path)
    target = load(target_path)
    src_servers = dict(source.get("mcpServers") or {})
    tgt_servers = target.get("mcpServers", {})
    if not isinstance(tgt_servers, dict):
        tgt_servers = {}

    if mode == "merge":
        stripped = 0
        if strip_cursor:
            managed = pack_managed_server_names(strip_cursor)
            for name in list(src_servers.keys()):
                if name in managed:
                    del src_servers[name]
                    stripped += 1
        for name, cfg in src_servers.items():
            tgt_servers[name] = cfg
        target["mcpServers"] = tgt_servers
        write(target_path, target)
        msg = f"merged {len(src_servers)} server(s) into {target_path}"
        if stripped:
            msg += f" (stripped {stripped} pack-managed)"
        print(msg)
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
