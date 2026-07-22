# Comprehension quiz

Hard gate at `/sc-discuss` exit: after blocking Q&A and (when visual) draft approval, **before** `clarify-status clear`.

Treat the human as the **customer / stakeholder**. Probe whether they understand the **mission spec** - requirement, process, and result - not the harness paperwork.

## Rules (agent-only - do not paste into chat)

- **Ask only useful questions.** Each question must measure real understanding of Goal / Output / Good-Bad / Verify / a recorded product decision / out-of-scope. If nothing material is unclear or worth probing, **do not invent filler** - record `Comprehension quiz: skipped - no material gaps to probe` and proceed.
- **1-5 questions** (prefer fewer). Zero is allowed via that skip.
- **Questions:** 5W1H open-ended. One W/H per question. Scope clear (name the product behavior or requirement, not a grab-bag).
- **Explanations (keys + fail coaching):** Feynman - plain words, teach a friend, one analogy max.
- Chat shows **only the questions**. No brand labels, no process narration.
- Human answers **before** model answers.
- On pass: `Comprehension quiz: passed`. On human skip: `Comprehension quiz: skipped - <reason>`.

### Ask about (customer lens)

| Topic | Examples |
|-------|----------|
| Spec / requirement | What are we building? Who is it for? What must be true when done? |
| Process | How does a user get from A to B? What happens on failure? |
| Result / Verify | How would you notice success without reading the code? |
| Decision / tradeoff | What did we choose, why not the alternative? |
| Out of scope | What are we explicitly not doing this round? |
| Visual (if UI) | What must the approved look keep? |

### Never ask (harness trivia)

- Where to write `decisions.md` / which exact line
- `clarify-status`, slash skill names, quiz mechanics, branch names, evidence labels
- Anything the customer would not need to know to approve the work

Exception: the **mission itself** is about that harness behavior - then ask as product behavior ("What new step blocks clear?"), still not "which markdown line".

## Chat template

```
### Quiz

1. …
2. …
3. …
```

## Example (≤7-task plan cap mission)

```
### Quiz

1. What hard limit does a plan phase get after this ships?
2. How may we split work when that limit is not enough?
3. Why reject a soft "prefer ≤7" or an 8-9 exception band?
```

Model answers (after human answers) - Feynman tone:

```
1. At most 7 tasks per phase - mandatory, not a suggestion.
2. Split into more phases in the same mission, or more missions on a roadmap.
3. Soft caps and exception bands are how the limit quietly dies; keep one clear rule.
```

## Grading rubric

- **Pass:** answer matches intended product/spec behavior in the human's own words.
- **Fail:** wrong scope/result, or cannot explain a must-know requirement.
- **Overall pass:** all asked items pass → human confirms → `Comprehension quiz: passed`.
- **Overall fail:** critical miss → do not clear; tighten spec or re-quiz.
- **Skip (human):** explicit request → `Comprehension quiz: skipped - <reason>`.
- **Skip (agent):** nothing material to probe → `Comprehension quiz: skipped - no material gaps to probe`.

If a fail came from a vague or harness-trivia question, rewrite toward spec/result or drop the question - do not punish the human for a bad question.

## Must not

- Mid-clarify quiz, or any quiz under `/sc-run`
- Filler / harness-meta questions to "have a quiz"
- Chat fluff; labeling "(Feynman)" / "(5W1H)"
- Showing the answer key before the human answers
- Clearing while a posed quiz is unanswered and no skip is recorded
