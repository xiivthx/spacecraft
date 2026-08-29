---
name: sc-automate-slack
description: "Optional Automations+Slack HIL AFK lane. Schedule/git may drive toward ready under /sc-run gates; Slack notify on handback/needs-HIL; Slack reply may cue resume. Never required for ready/ship. Slack resume ≠ AUTH / VERIFIED / ship. On Spacecraft stop or ready handoff: disarm Automations lane."
disable-model-invocation: true
---

# sc-automate-slack

## Goal

Optional Automations + Slack HIL AFK. Schedule/git may drive toward `ready` only as AFK under `/sc-run` gates; Slack notify on handback/needs-HIL; Slack reply may cue resume. Mission artifacts under `.space/missions/<id>/` remain SoT.

## Output

Greppable disposition (exactly one):

- `Automate-Slack: ran`
- `Automate-Slack: resumed`
- `Automate-Slack: stopped: <reason>`
- `Automate-Slack: skipped: <reason>`

## When to use

Overnight/AFK `/sc-run` when Automations + Slack HIL AFK helps. Else `Automate-Slack: skipped: <reason>`.

## Workflow

1. **Optional arm** - invoke Automate via `~/.cursor/skills-cursor/automate/SKILL.md`; SDK pointer only `~/.cursor/skills-cursor/sdk/SKILL.md` (no copy, no YAML/SDK bootstrap under this skill dir). Emit `Automate-Slack: ran` when armed or setup ran.
2. **Notify** - on handback or needs-HIL, notify Slack (channel/thread per Automation). Notify is cue only; artifacts stay SoT.
3. **Resume cue** - on Slack reply, cue resume under `/sc-run` gates. Emit `Automate-Slack: resumed` when a reply actually cued resume. Resume does **not** satisfy AUTH, ship form, or `set-state ready` alone.
4. **Spacecraft stop** - on `3-cycle` | `timebox` | `blocked`, when armed: **disarm**/teardown (editor disable / pause / unsubscribe) + `handback.md` → `Automate-Slack: stopped: <reason>`.
5. **Ready handoff** - when armed at ready: **disarm**/teardown → `Automate-Slack: stopped: ready-handoff`. Do not leave schedule/git Automation accepting Slack resume after ready or stop.
6. **Skip** - unused → `Automate-Slack: skipped: <reason>`.

Shared firewall: [../sc-run/references/optional-lanes.md](../sc-run/references/optional-lanes.md)

## Verify

- Disposition: `Automate-Slack: ran` | `Automate-Slack: resumed` | `Automate-Slack: stopped: <reason>` | `Automate-Slack: skipped: <reason>`
- Stop or ready handoff with prior armed lane: disarm/teardown; on stop also `handback.md`
- Upstream Automate + SDK by reference only

## Must not

- Treat Slack resume as AUTH / `VERIFIED` / ship (Slack resume ≠ AUTH)
- Treat Automation schedule/git/success chat as AUTH / `VERIFIED` / `set-state ready` / ship
- Soft-pass `set-state ready` or progress ship from Slack notify/reply alone
- Merge/push without Spacecraft ship gates + quoted `AUTH:`
- Leave schedule/git Automation accepting Slack resume after stop or ready handoff
- Add non-Slack first-party HIL on this map; copy Automate/SDK bodies; add Automation YAML/SDK bootstrap under this skill dir
