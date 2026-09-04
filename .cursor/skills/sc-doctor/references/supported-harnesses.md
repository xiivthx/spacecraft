# Supported harnesses

sc-doctor grades Cursor only. Collector ID `cursor`. Source: `~/.cursor/projects/<slug>/agent-transcripts`.

## Startup gate

Identify the executing harness from runtime context. Do not infer it from conversation files on disk. Identify Cursor from Cursor product identity, availability of `AskQuestion`, or this skill path under `.cursor/skills`.

If the executing harness is not Cursor, or cannot be identified confidently, stop before creating a report directory or reading conversation history. Tell the user:

> sc-doctor supports Cursor only (local agent transcripts). This run appears to be using an unsupported harness, so no conversations were read.

## Cursor source details

- Project slug: POSIX absolute workspace path, strip leading `/`, join remaining parts with `-`. Example: `/Users/me/src/spacecraft` becomes `Users-me-src-spacecraft`.
- Parent session: `~/.cursor/projects/<slug>/agent-transcripts/<uuid>/<uuid>.jsonl`
- Subagent: `.../agent-transcripts/<uuid>/subagents/<id>.jsonl` (only with `--include-subagents`)
- Repo filter uses slug match, not a `cwd` field inside the JSONL.
- User-question tool: `AskQuestion`
- `--cursor-home PATH` - nonstandard Cursor home (default `~/.cursor`). Do not scan other harness homes.

## Skill locations

- Project: `.cursor/skills`
- Global (with `--include-global-skills`): `~/.cursor/skills` (or `$CURSOR_HOME/skills`)

Never treat `~/.cursor/skills-cursor` as user or project skills (Cursor builtin).
