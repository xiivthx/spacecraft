#!/usr/bin/env python3
"""Tests for sc-doctor session collection (Cursor-only)."""

import io
import json
import os
import tempfile
import unittest
from contextlib import redirect_stderr, redirect_stdout
from datetime import datetime, timedelta, timezone
from pathlib import Path

from collect_sessions import (
    cursor_project_slug,
    detect_skills_from_entries,
    discover_skills,
    find_cursor_session_files,
    parse_args,
    parse_cursor_session,
)


def write_jsonl(path, records):
    path.parent.mkdir(parents=True, exist_ok=True)
    path.write_text("\n".join(json.dumps(record) for record in records) + "\n")


class CursorSessionTests(unittest.TestCase):
    def test_cursor_project_slug_encodes_posix_path(self):
        self.assertEqual(
            cursor_project_slug(Path("/Users/me/src/spacecraft")),
            "Users-me-src-spacecraft",
        )

    def test_discovers_cursor_project_and_global_skills_not_builtin(self):
        with tempfile.TemporaryDirectory() as tmp:
            root = Path(tmp)
            repo = root / "repo"
            cursor_home = root / "cursor-home"
            project_skill = repo / ".cursor" / "skills" / "sc-doctor" / "SKILL.md"
            global_skill = cursor_home / "skills" / "sc-git" / "SKILL.md"
            builtin_skill = (
                cursor_home / "skills-cursor" / "create-skill" / "SKILL.md"
            )
            for skill in (project_skill, global_skill, builtin_skill):
                skill.parent.mkdir(parents=True)
            project_skill.write_text("---\ndescription: Doctor\n---\n")
            global_skill.write_text("---\ndescription: Git\n---\n")
            builtin_skill.write_text("---\ndescription: Create Skill\n---\n")

            with_global = discover_skills(
                [repo],
                [],
                include_global=True,
                cursor_home=cursor_home,
            )
            project_only = discover_skills(
                [repo],
                [],
                include_global=False,
                cursor_home=cursor_home,
            )

            self.assertIn("sc-doctor", with_global)
            self.assertIn("sc-git", with_global)
            self.assertNotIn("create-skill", with_global)
            self.assertEqual(set(project_only), {"sc-doctor"})

    def test_discovers_skills_and_matches_sessions_across_projects(self):
        with tempfile.TemporaryDirectory() as tmp:
            root = Path(tmp)
            first = root / "first"
            second = root / "second"
            first_skill = first / ".cursor" / "skills" / "alpha" / "SKILL.md"
            second_skill = second / ".cursor" / "skills" / "beta" / "SKILL.md"
            decoys = (
                first / ".agents" / "skills" / "decoy-agents" / "SKILL.md",
                first / ".claude" / "skills" / "decoy-claude" / "SKILL.md",
                second / ".codex" / "skills" / "decoy-codex" / "SKILL.md",
            )
            for skill in (first_skill, second_skill, *decoys):
                skill.parent.mkdir(parents=True)
                skill.write_text("---\ndescription: Skill\n---\n")

            skills = discover_skills(
                [first, second],
                [],
                include_global=False,
            )

            self.assertEqual(set(skills), {"alpha", "beta"})

    def test_discovers_global_skills_without_repositories(self):
        with tempfile.TemporaryDirectory() as tmp:
            root = Path(tmp)
            cursor_home = root / "cursor-home"
            skill = cursor_home / "skills" / "foo" / "SKILL.md"
            skill.parent.mkdir(parents=True)
            skill.write_text("---\ndescription: Foo\n---\n")

            found = discover_skills(
                [],
                [],
                include_global=True,
                cursor_home=cursor_home,
            )
            hidden = discover_skills(
                [],
                [],
                include_global=False,
                cursor_home=cursor_home,
            )

            self.assertIn("foo", found)
            self.assertNotIn("foo", hidden)

    def test_detects_skills_from_deferred_tool_entries(self):
        entries = [
            ("tool:Skill", '{"skill": "alpha"}'),
            ("tool:read", '{"path": "/repo/.agents/skills/beta/SKILL.md"}'),
            ("assistant", "Mentioning gamma here does not count."),
        ]

        self.assertEqual(
            detect_skills_from_entries(entries, {"alpha", "beta", "gamma"}),
            {"alpha", "beta"},
        )

    def test_finds_parent_cursor_sessions_and_optional_subagents(self):
        with tempfile.TemporaryDirectory() as tmp:
            cursor_home = Path(tmp)
            slug = "Users-me-src-spacecraft"
            transcripts = (
                cursor_home / "projects" / slug / "agent-transcripts"
            )
            parent = transcripts / "parent-id" / "parent-id.jsonl"
            child = (
                transcripts / "parent-id" / "subagents" / "child-id.jsonl"
            )
            old = transcripts / "old-id" / "old-id.jsonl"
            other = (
                cursor_home
                / "projects"
                / "Users-other-workspaces-elsewhere"
                / "agent-transcripts"
                / "other-id"
                / "other-id.jsonl"
            )
            for path in (parent, child, old, other):
                write_jsonl(path, [{"role": "user"}])
            old_time = (
                datetime.now(timezone.utc) - timedelta(days=10)
            ).timestamp()
            os.utime(old, (old_time, old_time))
            cutoff = datetime.now(timezone.utc) - timedelta(days=1)
            repos = [Path("/Users/me/src/spacecraft")]

            parents = find_cursor_session_files(
                cursor_home, cutoff, False, repos=repos
            )
            with_subagents = find_cursor_session_files(
                cursor_home, cutoff, True, repos=repos
            )

            self.assertEqual([path for _, path in parents], [parent])
            self.assertEqual(
                {path for _, path in with_subagents},
                {parent, child},
            )
            self.assertNotIn(other, [path for _, path in parents])
            self.assertNotIn(other, [path for _, path in with_subagents])

    def test_parses_cursor_user_query_tools_skills_edits_and_errors(self):
        with tempfile.TemporaryDirectory() as tmp:
            path = Path(tmp) / "session-1.jsonl"
            user_text = (
                "<timestamp>Friday, Sep 4, 2026, 11:15 PM (UTC+7)</timestamp>\n"
                "<manually_attached_skills>\n"
                "Skill Name: sc-doctor\n"
                "Path: /tmp/repo/.cursor/skills/sc-doctor/SKILL.md\n"
                "---\n"
                "description: Doctor\n"
                "---\n"
                "# sc-doctor\n"
                "INLINED_SKILL_BODY_SENTINEL\n"
                "</manually_attached_skills>\n"
                "<user_query>\n"
                "Grade my skills\n"
                "</user_query>"
            )
            write_jsonl(path, [
                {
                    "role": "user",
                    "message": {
                        "content": [{"type": "text", "text": user_text}],
                    },
                },
                {
                    "role": "assistant",
                    "message": {
                        "content": [
                            {"type": "text", "text": "I will inspect it."},
                            {
                                "type": "tool_use",
                                "name": "Read",
                                "input": {
                                    "path": (
                                        "/tmp/repo/.cursor/skills"
                                        "/sc-git/SKILL.md"
                                    ),
                                },
                            },
                            {
                                "type": "tool_use",
                                "name": "StrReplace",
                                "input": {
                                    "path": "/tmp/repo/file.py",
                                    "old_string": "a",
                                    "new_string": "b",
                                },
                            },
                        ],
                    },
                },
                {
                    "type": "turn_ended",
                    "status": "error",
                    "error": "User aborted request",
                },
            ])

            meta, stats, entries, skills = parse_cursor_session(
                path,
                {"sc-doctor", "sc-git"},
                False,
            )

            self.assertEqual(meta["id"], "session-1")
            self.assertEqual(meta["originator"], "cursor")
            self.assertIsNone(meta["cwd"])
            self.assertEqual(stats["user_turns"], 1)
            self.assertEqual(stats["assistant_turns"], 1)
            self.assertEqual(stats["tool_calls"], 2)
            self.assertGreaterEqual(stats["error_outputs"], 1)
            self.assertTrue(stats["has_code_edits"])
            self.assertEqual(set(skills), {"sc-doctor", "sc-git"})
            user_texts = [text for role, text in entries if role == "user"]
            self.assertTrue(
                any("Grade my skills" in text for text in user_texts)
            )
            joined_user = "\n".join(user_texts)
            self.assertNotIn("<user_query>", joined_user)
            self.assertNotIn("</user_query>", joined_user)
            self.assertNotIn("INLINED_SKILL_BODY_SENTINEL", joined_user)
            self.assertIn(("assistant", "I will inspect it."), entries)
            tool_roles = {
                role for role, _ in entries if str(role).startswith("tool:")
            }
            self.assertIn("tool:Read", tool_roles)
            self.assertIn("tool:StrReplace", tool_roles)

    def test_excludes_cursor_subagents_by_default(self):
        with tempfile.TemporaryDirectory() as tmp:
            path = (
                Path(tmp) / "parent-id" / "subagents" / "child.jsonl"
            )
            write_jsonl(path, [
                {
                    "role": "user",
                    "message": {
                        "content": [{
                            "type": "text",
                            "text": "<user_query>\nInvestigate\n</user_query>",
                        }],
                    },
                },
                {
                    "role": "assistant",
                    "message": {
                        "content": [{
                            "type": "text",
                            "text": "Looking now.",
                        }],
                    },
                },
            ])

            self.assertIsNone(parse_cursor_session(path, set(), False))
            parsed = parse_cursor_session(path, set(), True)
            self.assertEqual(parsed[0]["thread_source"], "subagent")

    def test_collect_args_are_cursor_only(self):
        args = parse_args(["--out", "/tmp/x"])
        self.assertEqual(args.out, "/tmp/x")
        self.assertTrue(hasattr(args, "cursor_home"))
        self.assertTrue(hasattr(args, "repo"))
        self.assertTrue(hasattr(args, "all_conversations"))
        self.assertTrue(hasattr(args, "include_global_skills"))
        for dest in (
            "harness",
            "claude_home",
            "codex_home",
            "warp_db",
            "pi_home",
            "grok_home",
            "zcode_home",
        ):
            self.assertFalse(hasattr(args, dest), dest)

        buf = io.StringIO()
        err = io.StringIO()
        with redirect_stdout(buf), redirect_stderr(err):
            with self.assertRaises(SystemExit):
                parse_args(["--help"])
        help_text = buf.getvalue() + err.getvalue()
        for flag in (
            "--harness",
            "--claude-home",
            "--codex-home",
            "--warp-db",
            "--pi-home",
            "--grok-home",
            "--zcode-home",
        ):
            self.assertNotIn(flag, help_text)
        for flag in (
            "--cursor-home",
            "--repo",
            "--all-conversations",
            "--include-global-skills",
        ):
            self.assertIn(flag, help_text)


if __name__ == "__main__":
    unittest.main()
