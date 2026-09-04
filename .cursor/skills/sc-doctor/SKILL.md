---
name: "sc-doctor"
description: "Grades local Cursor/Spacecraft conversations and `.cursor/skills` against efficiency and code-quality rubrics, then drafts skill edits and a local report. Use when the user wants their Cursor agent setup graded from local conversation history, or asks which of their `.cursor/skills` are actually working."
---
# sc-doctor

Grade recent local Cursor transcripts; never upload them. Let `SKILL_ROOT` be the directory containing this SKILL.md.

## Step 0: Start the run

Read `$SKILL_ROOT/references/supported-harnesses.md`. Identify Cursor from product identity, `AskQuestion`, or this skill path under `.cursor/skills`. If the executing harness is not Cursor, follow the reference's stop behavior. Do not create a report directory or read conversation history.

Check for a current git repository:

```bash
git rev-parse --show-toplevel
```

Use `AskQuestion`. When a current repository is available, ask **"Which conversations should I grade?"** with:

1. **Conversations in this repository** - recommended.
2. **All conversations**.
3. **Choose projects to analyze**.

When there is no current repository, ask the same question with:

1. **All conversations** - recommended.
2. **Choose projects to analyze**.

If the user chooses projects, ask for one or more project paths. Expand and validate every path as a git repository before continuing. The run produces one combined report across those projects.

Then ask **"Which skills should I evaluate?"** with:

1. **Project skills + global skills** - recommended.
2. **Project skills only**.

For an all-conversations run, "Project skills" means skills from local git repositories inferred from the conversations' working directories. After these answers, proceed immediately.

Never write artifacts into the user's repo. Create one fresh, collision-free scratch directory per run and use it as `REPORT_DIR` for every artifact:

```bash
REPORT_DIR="$(mktemp -d "${TMPDIR:-/tmp}/sc-doctor-XXXXXXXX")"
```

## Step 1: Collect

Build collector arguments: `--repo "$REPO"` for the current repository; repeat `--repo PATH` for selected projects; `--all-conversations` for all. Add `--include-global-skills` for project and global skills; omit it for project skills only.

```bash
python3 "$SKILL_ROOT/scripts/collect_sessions.py" \
  --out "$REPORT_DIR" \
  <conversation-scope arguments> \
  <skill-scope arguments>
```

Optional: `--cursor-home PATH` (default `~/.cursor`), `--days N` (default 45), `--max-sessions N` (default 12), `--include-subagents`.

Read `$REPORT_DIR/inventory.json`. If `sessions_sampled` is 0, tell the user there is nothing recent to score in the selected conversation scope (suggest raising `--days` or choosing different projects) and stop. If `skills_found` is 0, continue - the report becomes a case for creating skills, and `skill_coverage` is 0.

## Step 2: Score each sampled transcript

Process datasets of 50 transcripts or fewer in a single batch. For datasets with more than 50 transcripts, use parallel batches (20 transcripts per batch recommended). Score in this local process, or local child agents only. Rubrics:

- `$SKILL_ROOT/scorers/efficiency.md`
- `$SKILL_ROOT/scorers/code-quality.md`

For each transcript in `$REPORT_DIR/transcripts/`, score both rubrics: label, numeric score from the label table, and a 1-3 sentence reason citing the transcript. Apply the code-quality scorer only where the transcript shows code changes; otherwise record `insufficient_evidence` and exclude that result from the code-quality average and failed-conversation filter.

## Step 3: Aggregate

- `raw_efficiency` = mean of efficiency scores across all scored sessions.
- `raw_code_quality` = mean of code-quality scores, excluding `insufficient_evidence`. If no session had enough evidence, set it to 0.5 and say so in the findings.
- Curve qualitative rubric means into letter-grade report scores with `curve(score) = 0.5 + 0.5 * score`.
- `efficiency = curve(raw_efficiency)`.
- `code_quality = curve(raw_code_quality)`.
- `skill_coverage` = fraction of sampled sessions where at least one installed skill was detected. If `skills_found` is 0, coverage is 0.
- `overall = 0.5 * efficiency + 0.35 * code_quality + 0.15 * skill_coverage.`

Then, define `failed_conversations` from each conversation's raw, uncurved scorer results. A conversation fails when at least one applicable efficiency or code-quality score is below `0.5`. An `insufficient_evidence` result does not make a conversation fail. Use only `failed_conversations` as evidence for skill-improvement suggestions and draft skill edits.

Then derive:

- `top_findings`: the 3 most impactful, specific patterns across sessions. Concrete and concise.
- `suggestions`: each names a skill (existing or proposed-new) and a specific change. Trace to waste or defects in `failed_conversations` - cite the failed session, scorer, and moment. An installed skill that never triggered in a failed conversation is usually a description problem.

## Step 4: Draft skill edits

Follow `$SKILL_ROOT/references/skill-improvements.md`. Read each skill path from `inventory.json`, write the improved file to `$REPORT_DIR/proposed/<skill-name>/SKILL.md`, and put `diff -u <current> <proposed>` in the suggestion's `diff` field. For a proposed-new skill, write the complete SKILL.md and set `diff` to its full content as an addition. Do not modify the user's real skill files.

## Step 5: Write report.json and render

Write `$REPORT_DIR/report.json`. Store the curved `efficiency` and `code_quality` values, literal `skill_coverage`, and weighted `overall` in `scores`; do not store the raw rubric means there.

```json
{
  "title": "Agent Skill Report",
  "generated_at": "<ISO timestamp>",
  "harness": "<harness from inventory.json>",
  "handle": "<repo_name from inventory.json>",
  "stats": {"sessions_analyzed": 0, "sessions_scanned": 0, "skills_found": 0, "skills_used": 0, "window_days": 45},
  "scores": {"efficiency": 0.0, "code_quality": 0.0, "skill_coverage": 0.0, "overall": 0.0},
  "top_findings": ["", "", ""],
  "suggestions": [{"skill": "", "change": "<one-sentence summary of the edit>", "evidence": "<which session(s) and what happened that motivates this>", "proposed_path": "<path under proposed/, if an edit was drafted>", "diff": "<unified diff, or full content for a new skill>"}]
}
```

```bash
python3 "$SKILL_ROOT/scripts/render_report.py" "$REPORT_DIR/report.json" --open
```

This writes a self-contained `$REPORT_DIR/report.html` and opens it locally. Share-as-png downloads a 1200x675 image on the user's machine.

## Step 6: Output

Tell the user the grade and the three findings. Finish with this exact summary, substituting the absolute `REPORT_DIR` path:

- Your agent skill report: file://$REPORT_DIR/report.html

Want me to apply these suggestions to your skills?
