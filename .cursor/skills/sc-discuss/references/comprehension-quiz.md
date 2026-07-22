# Comprehension quiz

`/sc-discuss` exit gate: after blocking Q&A and (when visual) draft approval, **before** `clarify-status clear`.

Human = customer. Probe **spec understanding** (requirement, process, result) - not harness paperwork.

## Procedure

1. **Probe** - From `spec.md` + product lines in `decisions.md`, list risks of skim-and-run (wrong Goal, missed Verify, silent scope creep). If none material → record skip and stop:
   `Comprehension quiz: skipped - no material gaps to probe`
2. **Ask** - 1-5 open 5W1H questions (one W/H each). Prefer fewer. Chat shows **only**:
   ```
   ### Quiz

   1. …
   2. …
   ```
3. **Key** - After answers: short Feynman key (plain, teach-a-friend) + pass/fail per item. Stop. Bad question → rewrite or drop; do not fail the human for agent ambiguity.
4. **Record**
   - Pass (all asked items + human confirms): `Comprehension quiz: passed`
   - Human skip: `Comprehension quiz: skipped - <reason>`
   - Fail: do not clear; tighten spec or re-quiz

## Question craft

| Do | Don't |
|----|-------|
| Spec / requirement / process / result / tradeoff / out-of-scope / approved look | `decisions.md` lines, `clarify-status`, slash skills, quiz mechanics, branches, evidence labels |
| One clear W/H + named product behavior | Filler to "have a quiz"; grab-bag "name anything" |
| Own-words answers | Yes/no trivia; brand labels (`Feynman`, `5W1H`); key before answers |

Harness missions: ask as **product behavior** ("What new step blocks clear?"), never "which markdown line".

## Example

Mission: hard ≤7 tasks per plan phase.

```
### Quiz

1. What hard limit does a plan phase get after this ships?
2. How may we split work when that limit is not enough?
3. Why reject a soft "prefer ≤7" or an 8-9 exception band?
```

Key (after answers):

```
1. At most 7 tasks per phase - mandatory, not a suggestion.
2. More phases in the same mission, or more missions on a roadmap.
3. Soft caps and exception bands quietly kill the limit.
```

## Must not

- Mid-clarify or `/sc-run`
- Invent questions when the skip above applies
- Clear while a posed quiz is unanswered (unless skip recorded)
