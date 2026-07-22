# Learned

Knowledge migrated from missions for internal research reuse.

## Solved

| Mission | Date | Problem | Solution | Evidence |
|---------|------|---------|----------|----------|
| M07SMYDMZ | 2026-07-21 | Lean CLI rewrite inverted #38: branch mid outranked `.space/current` from `spacecraft use`; evidence ignored current | Prefer `.space/current` in `resolveActive`; route evi without `--mission` through it; update sc-mission priority docs | cmd-spacecraft-full-suite, validate --strict |
| M07SML0VX | 2026-07-21 | Lean CLI rewrite dropped SHA-256 evidence integrity (#37); capture had no hash and validate did not check | Restore optional `outputHash` (hex SHA-256 of `output`) on write; validate mismatches when present; legacy omit still passes | ship-reverify, validate --strict |

## Lessons

| Mission | Date | Lesson | Why it matters |
|---------|------|--------|----------------|
| M7ZMDXPQ | 2026-07-21 | Forced decision artifacts (INTENT/AUTH/TWINS/verdicts) beat mid-list prose for agent compliance | Harness gates that need rare high-stakes compliance should require greppable artifacts |
| M7ZMDXPQ | 2026-07-21 | Adversarial prove before ready must re-observe evidence; never soft-ship past REFUTED | Release gates in agent workflows should re-run claims and hunt weakened tests / false completion / unauthorized action |
| M7ZR8E5V | 2026-07-21 | Docs/prose "must contain phrase X" acceptances are TDD skip - do not invent phrase-echo RED harnesses | AFK otherwise burns checkpoints proving nothing beyond the sentence about to be pasted |
| M7ZR8E5V | 2026-07-21 | ≤7 jigsaw tasks per phase stays a hard Must; escape via phases or roadmaps only | Soft prefer / 8-9 bands remove the pressure that keeps plans honest |
| M7YB6X29 | 2026-07-21 | If a work branch encodes a mission id, create matching mission artifacts before `/sc-ship` even when the change is already complete | Closeout cannot invent Verify; retroactive scaffold + real evidence unblocks docs-only ships |
| M07SMYDMZ | 2026-07-21 | After rewriting a signal-priority resolver, keep regression tests that `use`/explicit override still beats heuristics even if CHANGELOG claims the old fix | Priority bugs return silently when only the previous implementation was tested |
| M07SML0VX | 2026-07-21 | When rewriting a subsystem, re-port integrity checks that lived on the old schema rather than assuming changelog history means the behavior still exists | Silent loss of security/integrity features after migrations ships false confidence |
| M81KH2YS | 2026-07-22 | Draft HTML quality improves when agents pick from a named locked layout pool plus shared directives, rather than inventing page architecture each run | Separating shared fidelity rules from optional art-direction packs keeps anti-slop authoritative and reduces one-off structure drift |
