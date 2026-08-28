# Spacecraft for Google Antigravity

Spacecraft provides a first-class mission-control harness for **Google Antigravity** with strict **Source of Trust (Hierarchy of Authority)**, specialized multi-agent subagent swarms, and deterministic **UX/UI/Frontend design gates**.

---

## 1. Overview & Architecture

Spacecraft integrates with Antigravity via the **Customization System**:
- **Global Plugin:** Installed at `~/.gemini/config/plugins/spacecraft/` (available across all workspaces).
- **Project Scaffold:** `.agents/` directory (rules, skills, hooks) + `GEMINI.md` at repository root.
- **Mission Evidence & Artifacts:** `.space/` (isolated, gitignored mission state).
- **CLI Utility:** `spacecraft` linked to `~/.local/bin/spacecraft`.

```
~/.gemini/config/plugins/spacecraft/
├── plugin.json       # Plugin manifest
├── rules/            # Consolidated rulebook (AGENTS.md)
├── hooks.json        # PreToolUse safety hooks (branch protection)
├── hooks/            # Safety script (safety-check.mjs)
├── agents/           # 10 Subagent contracts (sc-planner, sc-designer, sc-coder, ...)
└── skills/           # 25 Modular lifecycle & domain skills (sc-run, sc-discuss, sc-ship, ...)
```

---

## 2. Hierarchy of Authority (Source of Trust)

Every prompt, agent action, and code modification is strictly governed by this descending order of authority:

```
┌────────────────────────────────────────────────────────┐
│ 1. Explicit Human User Directive                      │  (Highest Authority)
├────────────────────────────────────────────────────────┤
│ 2. Visual SoT ([data-draft-surface] in Draft HTML)    │
│    + Behavioral SoT (spec.md + design-contract.md)     │
├────────────────────────────────────────────────────────┤
│ 3. Design Tokens (DESIGN.md Orbital Console)          │
│    + Process Rules (GEMINI.md / AGENTS.md)             │
├────────────────────────────────────────────────────────┤
│ 4. Automated Evidence & Oracles                        │
│    (evidence.jsonl, review.json, sc-judge, tests)      │
├────────────────────────────────────────────────────────┤
│ 5. Existing Codebase Implementation                    │  (Lowest Authority)
└────────────────────────────────────────────────────────┘
```

> [!IMPORTANT]
> **Prompt Integrity Standard (Spec Contract):**
> Every prompt and skill instruction must define **Goal**, **Output**, **Good vs Bad**, and **Verify**. Never invent verification steps.

---

## 3. UX / UI / Frontend Engineering System

### A. Design System (`DESIGN.md` — Orbital Console)
- **Token Compliance:** Strict usage of `--sc-*` custom properties (`--sc-bg: #0e1116`, `--sc-surface: #151a21`, `--sc-accent: #f6b44b`, `--sc-cyan: #62d6cf`).
- **Layout & Spacing:** Strict `4px / 8px` rhythm. Left-aligned technical operator layouts over marketing shells.
- **Anti-Slop Directives:** Forbid purple gradient blobs, multi-layer drop shadows, nested floating cards, and pill-happy interactive controls.

### B. Visual Draft HTML (Visual Source of Truth)
- For visual UI missions, `/sc-discuss` requires an approved Draft HTML under `.space/missions/<id>/design/drafts/`.
- Implementation in `/sc-run` must **port** structure, tokens, spacing, typography, and chrome from `[data-draft-surface]`.

### C. Draft Parity & Paired Screenshots
- Capture screenshots of `[data-draft-surface]` and the live running application at identical viewports:
  - Mobile: `375px`
  - Tablet: `768px`
  - Desktop: `1280px`
  - Wide: `1536px`
- Side-by-side comparison before marking UI tasks ready.

### D. Chrome DevTools MCP Live Probing
- Before visual UI can be marked `ready`, `sc-browser-probe` sweeps the running product using **Chrome DevTools MCP** tools:
  - `navigate_page` / `resize_page`
  - `take_screenshot` / `take_snapshot`
  - `list_console_messages` (zero unexpected errors)
  - `list_network_requests` (zero failed endpoints)
  - `lighthouse_audit` (Core Web Vitals & accessibility)
- Runs an AFK fix-loop until `PROBE: CLEAN`.

---

## 4. Multi-Agent Subagent Swarm

Commander coordinates the mission without directly writing product code:

| Subagent | Role Contract |
|---|---|
| `sc-planner` | Slices `spec.md` into ≤7 jigsaw tasks per phase with design contracts and approved scenarios. |
| `sc-designer` | UX lead: orchestrates Impeccable (primary craft); owns Spacecraft port gates; approved draft HTML remains visual SoT. |
| `sc-coder` | Surgical implementation of code making failing tests pass (GREEN) or executing triage skips. |
| `sc-tester` | TDD RED tests, scenario test suites, diff-coverage, and mutation checks. |
| `sc-browser-probe` | Live browser sweep using Chrome DevTools MCP with AFK fix-loop to `PROBE: CLEAN`. |
| `sc-reviewer` | Review gate against requirements; enforces zero findings in `review.json`. |
| `sc-judge` | Adversarial verify; hunts for draft drift and ungrounded claims (`VERIFIED` vs `REFUTED`). |
| `sc-writer` | Specifications, documentation, prompt refinement, and narrative context harvest. |
| `sc-adviser` | Architectural trade-offs and domain consultation. |
| `sc-firmware` | Embedded hardware and firmware specialist. |

---

## 5. Lifecycle & Lanes

```
/sc-discuss (Clarify & Draft) → /sc-run (TDD, Build, Probe, Judge) → Human Check → /sc-ship
```

- **Discuss Lane (`/sc-discuss`):** Clarify requirements, brainstorm, size missions (≤7 tasks), approve Draft HTML.
- **Mission Lane (`/sc-run`):** AFK runner executing jigsaw tasks, TDD, Chrome DevTools probe, review, and `sc-judge`.
- **Quick Lane (`/sc-quick`):** Direct surgical edits without mission ceremony (still enforces `INTENT:`, `AUTH:`, `TWINS:`).
- **Ship Lane (`/sc-ship`):** Squashes commits (≤5), merges `--no-ff`, and archives mission state.

---

## 6. Safety Hooks & Inner-Loop Gates

Three layers: **hooks** (hard) → **short AGENTS.md** (always-on soft) → **skills** (on demand).

- **Hooks (`safety-check.mjs`):** deny secret paths, force-push, catastrophic `rm`, mutate on `main`/`master`, ship without `SPACECRAFT_SHIP=1`. **`git push` denied in-agent** (human pushes after AUTH).
- **`INTENT:`** Class + intended behavior before behavior-changing edits.
- **`AUTH:`** Quoted user permission before outward actions (does not unlock hooks alone).
- **`TWINS:`** Search project-wide for identical bug construct after fixing a defect.
- **3-Cycle Stop:** Escalate after 3 failed fix-verify cycles.

---

## 7. Bilingual Standard (TH / EN)

- **English:** Technical substance, architecture, code, JSON schemas, diffs, test logs, evidence files.
- **Thai:** Human-in-the-loop (HIL) clarification questions, decision options, short progress status, and phase handoffs.

---

## 8. Installation Commands

### Global Install (Plugin for Antigravity)
```sh
make install-antigravity
```

### Project Scaffold (.agents + GEMINI.md)
```sh
./bootstrap.sh --antigravity /path/to/project
# or from repo:
make install-antigravity-project PROJECT=/path/to/project
```
