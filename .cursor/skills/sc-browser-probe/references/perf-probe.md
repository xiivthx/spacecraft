# Perf probe (optional)

Use only when the user asks about rate/throughput/feasibility, or the product documents a target (e.g. tags/sec).

## Goal

Answer whether a stated rate is plausible under the current setup, with measured or clearly assumed data - not vibes.

## Output

```markdown
### Perf
- Target rate: <e.g. 20 tags/sec>
- Method: measured | estimated
- Per scenario (or load step):
  | step | count | wall time | rate | notes |
  |------|-------|-----------|------|-------|
  | … | … | … | …/sec | … |
- Feasible at target?: yes | no | unknown
- Bottleneck guess: UI thread | network | device | sim | unknown
- Recommended rate: <n>/sec (with margin) or n/a
- Machine / setup estimate: <inputs assumed: CPU, browser, sim vs hardware, batch size>
- Assumptions: …
```

## Procedure

1. **Define unit** - What is one "tag" / event / op? Must match user language.
2. **Measure when possible** - Time N operations in-browser or via exposed sim controls; compute rate = N / seconds.
3. **Per scenario** - If scenarios differ (short vs long payload), report rate per scenario, not one blended number.
4. **Feasibility** - Compare measured/estimated rate to target with margin (default ~20% headroom unless user sets otherwise).
5. **Recommend** - Suggest a sustainable rate and/or setup changes (batch size, fewer DOM updates, hardware vs sim, headless load tool). Label guesses as guesses.
6. **Unknown** - If UI cannot be timed, say `unknown` and list what instrumentation would unlock a real answer.

## Good / Bad

- Good: unit defined; table with wall time; assumptions listed; separate functional pass from perf
- Bad: "should be fine"; inventing 20/sec without timing; treating perf fail as functional `PROBE: ISSUES` unless user made rate a hard bar

## Verify

- Unit of work defined
- At least one timed or explicitly estimated row
- Feasible answer is `yes` | `no` | `unknown` with reason
