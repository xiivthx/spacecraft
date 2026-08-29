---
name: sc-security
description: "Static security review of source and manifests. Activate on \"security review\", \"check for secrets\", \"OWASP check\", \"injection scan\", \"audit dependencies\", or hardcoded credentials."
disable-model-invocation: true
---

# sc-security

Fallback and on-demand static security review of source code and dependency manifests. Cursor `security-review` is the primary ready-path security surface; this skill runs only when that subagent fails or the human/Commander explicitly requests a heuristic scan. Apply pattern-based heuristics only; never execute dynamic audit tools. Flag findings.

## When to use

**Role:** fallback + on-demand only. Ready-path security uses Cursor `security-review` first.

Activate on these triggers:

- Cursor `security-review` subagent failure on a security-in-scope ready path (fallback)
- Explicit on-demand heuristic: "security review" or "security check"
- "check for secrets", "find hardcoded credentials", "scan for tokens"
- "OWASP check" or "OWASP Top 10"
- "injection scan", "SQL injection", "command injection"
- "audit dependencies" or "check manifest for vulnerabilities"

## Workflow

### 1. Scope from mission artifacts

Read `spec.md` and the active `plan.json` task to identify security boundaries, sensitive data, and acceptance criteria. Limit the review to files listed or implied by the task scope.

### 2. Static pattern sweep

Scan all in-scope source and manifest files using `.cursor/rules/300-security.mdc` (domain checklist) and `references/patterns.md` (heuristic patterns). Do not execute code, run tools, or fetch live vulnerability databases.

### 3. Classify findings

For each match, assign checklist priority (Critical / High / Medium / Low / informational), map to house severity per rule 300 (`critical` / `important` / `minor`), and an OWASP category when applicable. Informational stays scan-note only unless promoted. Finding craft: `.cursor/skills/sc-run/references/defect-finding.md`.

### 4. Report

Emit findings in the output format. Include file path, line, matched pattern, evidence snippet, and a concrete fix. Do not modify files during review.

### 5. Record in mission evidence

Add findings to the task output or `evidence.jsonl` notes. Optional evidence label: `spacecraft evidence "sc-security-…" -- …`.

When this skill runs as ready-path **fallback** after Cursor `security-review` failure or unavailable skip, Commander Must write greppable disposition in `decisions.md`:

- `Sc-security fallback: pass` - fallback cleared security-when-in-scope (no open critical/important security findings blocking ready)
- Optionally also `Sc-security fallback: findings drained` - when fallback findings were fixed/drained before pass

Judge and reviewer accept only these greppable lines (plus SEC machine-evidence pass) - not free-text "sc-security ran" prose.

## Rules

### General

- **Must**: Review source and manifests only. No execution, no dynamic scans.
- **Must not**: Run any runtime or dynamic vulnerability audit tool, dependency scanner, or live CVE database query.
- **Must**: Treat every hardcoded secret as critical until proven otherwise.
- **Must**: Prefer false positives over missed negatives on the first sweep; triage during classification.
- **Must**: Use `spec.md` acceptance criteria and `plan.json` scope to decide whether a finding is relevant.
- **Must**: Load `references/patterns.md` for the heuristic pattern list (OWASP, secrets, injection, manifests).

## Out of scope

- Dynamic or runtime security testing
- Penetration testing, fuzzing, or network scanning
- Social engineering or physical security
- Infrastructure hardening beyond source/manifest patterns
- General code quality or architecture review
- Mission lifecycle management
- Writing fixes - report only; leave implementation to the build step

## Output format

```
Scope: <files reviewed>
Patterns: <list of pattern categories used>
Findings:
  - Severity: critical | important | minor
    Title: <impact-first, ≤12 words - required for critical/important>
    Category: <OWASP category or pattern type>
    File: <path>
    Line: <number>
    Pattern: <matched pattern or heuristic>
    Evidence: <short snippet>
    Issue: <attack vector in plain English>
    Impact: <user/system effect - required for critical/important>
    Fix: <concrete remediation = requiredFix>
    Retest: <2-3 local verify ideas - required for critical/important>
Summary:
  Total: <n>
  Critical: <n>
  Important: <n>
  Minor: <n>
  By category: <counts>
Recommendation: pass / fix-before-merge / block
```

## Checklist

- [ ] Mission scope confirmed from `spec.md` and `plan.json`
- [ ] All source files in task scope reviewed (rule 300 + `references/patterns.md`)
- [ ] Hardcoded secrets scan completed on source and config files
- [ ] SQL/command injection patterns checked
- [ ] Manifest and lock files reviewed for suspicious version patterns
- [ ] Findings classified by house severity (critical / important / minor) and OWASP category
- [ ] Each finding includes file path, line, evidence snippet, and fix guidance
- [ ] No dynamic audit tools executed
- [ ] Results recorded in task output or `evidence.jsonl`
- [ ] Ready-path fallback: greppable `Sc-security fallback: pass` (and optionally `Sc-security fallback: findings drained`) written in `decisions.md` when fallback clears the gate

## References

- `.cursor/rules/300-security.mdc` - domain checklist and severity mapping
- `references/patterns.md` - OWASP / secrets / injection / manifest heuristics (on-demand)
- `.cursor/skills/sc-run/references/defect-finding.md` - house finding craft for review.json / run summary
- OWASP Top 10 (2021) - https://owasp.org/Top10/
- CWE/SANS Top 25 - https://cwe.mitre.org/top25/
- OWASP Cheat Sheet Series - https://cheatsheetseries.owasp.org/
- Mission artifacts: `spec.md`, `plan.json`, `evidence.jsonl`
