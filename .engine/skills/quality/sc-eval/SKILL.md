---
name: sc-eval
description: >
  Agent-facing eval framework. Create labelled eval examples, write rubrics,
  and run eval suites for mission evidence quality assessment.
license: MIT
compatibility: opencode
metadata:
  version: 1.0
  audience: xiivthx
---

# sc-eval

Create and run structured evaluation suites for agent mission evidence. The eval framework supports deterministic checks, rubric-based 0-4 scoring across 5 dimensions, and optional LM Judge evaluation with graceful fallback.

## When to use

Activate when the user asks to:

- Create an eval dataset for a mission
- Write a rubric for scoring agent output
- Run an eval suite against mission evidence
- Check eval coverage before shipping
- Add labelled eval examples

## Workflow

### Creating an eval dataset

1. Scaffold the eval directory:
   ```
   scripts/spacecraft eval init <mission-id>
   ```
   Creates `.space/evals/<mission-id>/` with:
   - `rubric.json` — 5 dimensions on 0-4 scale
   - `dataset.json` — labelled eval examples
   - `config.json` — coverage threshold (default 0.8)

2. Add labelled examples to `dataset.json`:
   ```json
   {
     "examples": [
       {
         "id": "ex1",
         "label": "T1 compile check",
         "evidenceRefs": ["E07N4OOBX"],
         "expectedPass": true,
         "expectedScore": 4
       }
     ]
   }
   ```

3. Customise `rubric.json`:
   ```json
   {
     "dimensions": [
       {
         "name": "task_success",
         "description": "Did the agent accomplish stated task objectives?"
       }
     ]
   }
   ```
   The 5 standard dimensions are: `task_success`, `tool_use_quality`, `trajectory_compliance`, `hallucination`, `response_quality`.
   Each dimension is scored 0-4.

### Running eval suites

```
scripts/spacecraft eval <mission-id>
```

This runs the full eval pipeline:
1. **Deterministic checks** — Compile, test, lint pass/fail from evidence exit codes
2. **Trajectory analysis** — Tool-call sequence and step correctness from evidence entries
3. **Rubric scoring** — Each dimension scored 0-4 against rubric criteria
4. **LM Judge** (optional) — Secondary model evaluation of non-deterministic output quality

Results are written to `evidence.jsonl` with `"type": "eval"`.

### Configuring LM Judge

Set the `SPACECRAFT_JUDGE_MODEL` environment variable to enable LM Judge evaluation:

```
export SPACECRAFT_JUDGE_MODEL="gemini-pro"
```

The LM Judge invokes the model via `agy`, `llm`, or `opencode` CLI (tried in that order). When no judge is configured or available, the eval falls back to deterministic-only.

### Configuring coverage threshold

Edit `.space/evals/<mission-id>/config.json`:
```json
{
  "coverageThreshold": 0.8
}
```

Coverage is calculated as `covered_checks / total_checks` where covered checks are labelled examples in the dataset and total checks are plan tasks.

## Ship gating

The `sc-ship` process checks eval coverage via the `releaseReadiness.evalCoverage` gate in `review.json`. When coverage is below the configured threshold, `spacecraft closeout-check` blocks the ship.

To satisfy the gate after running eval:
```json
{
  "releaseReadiness": {
    "evalCoverage": {
      "status": "passed"
    }
  }
}
```

## Eval result structure

Results written to `evidence.jsonl`:
```json
{
  "id": "E07N4YG9Z",
  "type": "eval",
  "label": "eval-M07N361SC",
  "command": "spacecraft eval M07N361SC",
  "exitCode": 0,
  "stdout": "{...full eval result...}",
  "createdAt": "2026-07-12T..."
}
```

The `stdout` field contains the full `EvalResult` JSON with deterministic checks, rubric scores, LM judge results, and coverage metrics.

## Commands reference

| Command | Purpose |
|---------|---------|
| `spacecraft eval <id>` | Run full eval suite |
| `spacecraft eval init <id>` | Scaffold eval directory |

## References

- `spec.md` — Mission acceptance checks
- `plan.json` — Tasks for coverage calculation
- `evidence.jsonl` — Input evidence and output target for eval results
