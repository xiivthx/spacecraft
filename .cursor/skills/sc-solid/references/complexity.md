# Managing Complexity

> Consult when: code feels hard to understand, a small change touches many files, or you're tempted to add "just in case" abstractions.

## Essential vs Accidental

- **Essential complexity** = inherent to the problem domain. Business rules, domain logic. Cannot remove, only manage.
- **Accidental complexity** = introduced by our solutions. Poor abstractions, framework ceremony, premature optimization. Must minimize.

**Goal:** Minimize accidental complexity while clearly expressing essential complexity.

---

## Three Detection Signals

| Signal | Symptom | Common Cause |
|--------|---------|--------------|
| **Change amplification** | "To add this field, I need to update 15 files." | Scattered responsibilities, poor boundaries |
| **Cognitive load** | "I need to understand 10 classes to understand this one." | Tight coupling, hidden deps, unclear naming |
| **Unknown unknowns** | "I changed this and something unrelated broke." | Global state, implicit contracts |

---

## KISS - Keep It Simple

> The simplest solution that passes the test is the best.

- Start with the obvious solution. Add complexity only when required.
- Prefer boring, well-understood approaches.
- Question every abstraction: "Does this reduce or increase cognitive load?"

---

## YAGNI - You Aren't Gonna Need It

> Don't build what you don't need NOW. Build what the mission task asks for.

Cost of YAGNI violations:
1. Development time wasted on unused features
2. Maintenance burden - code that exists must be maintained
3. Cognitive load - more code to understand
4. Wrong abstraction - guessing future needs incorrectly

**If a mission task doesn't ask for it, don't build it.**

---

## DRY - But Wait for Three

> Every piece of knowledge should have a single representation. But don't extract prematurely.

### The Rule of Three

```
Duplication #1 → Leave it
Duplication #2 → Note it, maybe leave it
Duplication #3 → NOW extract it
```

Why? The wrong abstraction is 10x worse than duplication. Wait for the pattern to emerge.

---

## The Four Elements of Simple Design (XP)

In priority order:

1. **Runs all the tests** - if it doesn't work, nothing else matters
2. **Expresses intent** - clear names, obvious structure. Code tells the story.
3. **No duplication** - DRY (Rule of Three). Single source of truth.
4. **Minimal** - fewest classes and methods possible. Remove anything unnecessary.

If all four are true, the design is simple enough. Stop optimizing.

---

## Technical Debt

### When to pay down:
- It's in your path (you're already touching the code)
- It blocks a mission task
- It's causing bugs

### When to defer:
- Code works and won't change
- Being replaced soon
- No test coverage (add tests first)

### Boy Scout Rule:
Leave the code better than you found it. One small improvement per touch: fix a name, extract a method, add a missing test.

---

## Spacecraft integration

- When a `plan.json` task feels complex, run the 3 detection signals. If any fire, split the task.
- Record complexity trade-offs in `decisions.md` (e.g., "deferred extraction - only 2 duplications, waiting for third").
- The 4 elements of simple design are the acceptance criteria for "done" beyond test passing.
