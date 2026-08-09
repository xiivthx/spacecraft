# Call examples

## Generic full-system

```
/sc-browser-probe
target: http://localhost:5173
scope: full
examples: empty form, valid sample, max-length string
perf: none
timebox: 20m
```

## Feature-scoped after fix

```
/sc-browser-probe
target: start app
scope: feature:checkout
examples: expired card then valid card
perf: none
timebox: 15m
```

## NFC / sim-shaped (product-specific - call site only)

```
/sc-browser-probe
target: <nfc roll sim URL>
scope: feature:nfc-roll
examples: short url, long url, trc sample
perf: is 20 tags/sec feasible? per-scenario rate + recommend speed / machine estimate
timebox: 25m
```

Do not bake NFC fixtures into SKILL.md core - pass them as `examples:` / `perf:`.

## After sc-run handoff (recommend-only escape net)

When `/sc-run` reaches green and the change includes UI or multi-step workflow, end handoff with:

```
Recommend: /sc-browser-probe (UI/workflow touched)
```

Recommend-only escape net - not a ready gate. Does not replace sc-verification or sc-judge.
