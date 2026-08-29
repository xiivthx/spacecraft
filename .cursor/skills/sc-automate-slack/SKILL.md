---
name: sc-automate-slack
description: "Optional Automations+Slack HIL AFK lane. Schedule/git may drive toward ready under /sc-run gates; Slack notify on handback/needs-HIL; Slack reply may cue resume. Never required for ready/ship. Slack resume ≠ AUTH / VERIFIED / ship. On Spacecraft stop or ready handoff: disarm Automations lane."
disable-model-invocation: true
---

# sc-automate-slack

## Goal

Optional Automations + Slack HIL AFK lane. Schedule or git triggers may drive work toward `ready` only as AFK under existing `/sc-run` gates; Slack notify on handback / needs-HIL; a Slack reply may cue resume toward `ready`. Humans or agents may arm it; it is **never required** for ready or ship. Spacecraft owns ready/ship. Mission artifacts under `.space/missions/<id>/` remain SoT. Slack resume is **not** AUTH. Automation schedule/git/success chat is never AUTH / `VERIFIED` / `set-state ready` / ship authority.

## Output

Greppable disposition (exactly one):

- `Automate-Slack: ran`
- `Automate-Slack: resumed`
- `Automate-Slack: stopped: <reason>`
- `Automate-Slack: skipped: <reason>`

## Good / Bad

- Good: optional Automations+Slack during overnight/AFK `/sc-run`; Slack notify on handback/needs-HIL; Slack reply cues resume under existing `/sc-run` gates only; Spacecraft stop or ready handoff → disarm/teardown Automations lane; mission artifacts SoT; upstream Automate/SDK by reference; disposition greppable; never required for ready/ship
- Bad: treating Slack reply, Automation schedule/git/success chat, or notify delivery as AUTH / `VERIFIED` / `set-state ready` / ship authority; leaving schedule/git Automation accepting Slack resume after ready or stop; merge/push without Spacecraft ship gates + quoted `AUTH:`; requiring Automations+Slack for ready/ship; non-Slack first-party HIL on this map; copying upstream Automate/SDK skill bodies; Automation YAML or SDK bootstrap under this skill dir; soft-passing `set-state ready` from Slack alone

## When to use

During overnight/AFK `/sc-run` when Automations + Slack HIL AFK helps. Skip when the human runs `/sc-run` without this lane (`Automate-Slack: skipped: <reason>`).

## Workflow

1. **Optional arm** - when Automations + Slack HIL AFK is wanted, invoke Cursor Automate by reading `~/.cursor/skills-cursor/automate/SKILL.md` (do not copy its body into this repo). SDK details only via `~/.cursor/skills-cursor/sdk/SKILL.md` (pointer only; do not copy; no SDK bootstrap under this skill dir). Emit `Automate-Slack: ran` when the lane was armed or Automations setup ran for this mission.
2. **Notify** - on handback or needs-HIL, notify Slack (channel/thread per the Automation). Mission artifacts under `.space/missions/<id>/` stay SoT; notify is cue only.
3. **Resume cue** - on Slack reply, cue resume AFK under existing `/sc-run` gates. Emit `Automate-Slack: resumed` when a Slack reply actually cued resume. Resume does **not** satisfy AUTH, ship form, or `set-state ready` alone.
4. **Spacecraft stop** - on `3-cycle` | `timebox` | `blocked`, when the lane was armed, Must both:
   1. **Disarm** / teardown the Automations lane: editor disable / pause / unsubscribe (do not leave schedule/git Automation accepting Slack resume).
   2. **Write** `.space/missions/<id>/handback.md` with stop reason + remaining work cue (per `/sc-run`).
   Then emit `Automate-Slack: stopped: <reason>` (e.g. stop reason).
5. **Ready handoff** - when overnight/AFK `/sc-run` reaches **ready** handoff and Automations+Slack was armed, Must **disarm**/teardown (editor disable / pause / unsubscribe) and emit `Automate-Slack: stopped: ready-handoff`. Do not leave schedule/git Automation accepting Slack resume after ready or stop.
6. **Skip** - if the lane is not used → `Automate-Slack: skipped: <reason>`.

## Boundaries

- This lane is **never required** for ready or ship.
- Must not treat Slack as AUTH / `VERIFIED` / ship authority. Slack resume ≠ AUTH.
- Must not treat Automation schedule/git/success chat as AUTH, `VERIFIED`, `set-state ready`, or ship authority. Schedule/git may drive toward ready only as AFK under existing `/sc-run` gates.
- Must not merge or push without Spacecraft ship gates + quoted `AUTH:`.
- Mission artifacts under `.space/missions/<id>/` remain SoT.
- Slack only for first-party HIL on this map; other messengers out of scope. Must not add non-Slack first-party HIL on this map.
- Must not leave schedule/git Automation accepting Slack resume after Spacecraft stop or ready handoff; Must disarm/teardown when the lane was armed.
- Must not copy upstream Automate or SDK skill bodies into this repo; reference only.
- Must not add Automation YAML or SDK bootstrap under this skill directory.
- Must not soft-pass `set-state ready` or progress ship from Slack notify/reply alone.

## Verify

- Disposition line present: `Automate-Slack: ran` | `Automate-Slack: resumed` | `Automate-Slack: stopped: <reason>` | `Automate-Slack: skipped: <reason>`
- Firewall greppable: Slack resume ≠ AUTH; Must not AUTH/`VERIFIED`/ship from Slack; Must not treat Automation schedule/git/success as AUTH/`VERIFIED`/`set-state ready`/ship; drive-toward-ready only under `/sc-run` gates; never required for ready/ship; artifacts SoT
- On Spacecraft stop or ready handoff with a prior armed lane: disarm/teardown (editor disable / pause / unsubscribe); on stop also `handback.md`
- Upstream invoked only via `~/.cursor/skills-cursor/automate/SKILL.md`; SDK pointer only `~/.cursor/skills-cursor/sdk/SKILL.md`
