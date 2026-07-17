---
name: sc-security
description: >
  Static security review of source and manifests. Activate on "security review", "check for secrets", "OWASP check", "injection scan", "audit dependencies", or hardcoded credentials.
---

# sc-security

Static security review of source code and dependency manifests. Apply pattern-based heuristics only; never execute dynamic audit tools. Flag findings.

## When to use

Activate on these triggers:

- "security review" or "security check"
- "check for secrets", "find hardcoded credentials", "scan for tokens"
- "OWASP check" or "OWASP Top 10"
- "injection scan", "SQL injection", "command injection"
- "audit dependencies" or "check manifest for vulnerabilities"
- Before final evidence capture for any mission task with security exposure

## Workflow

### 1. Scope from mission artifacts

Read `spec.md` and the active `plan.json` task to identify security boundaries, sensitive data, and acceptance criteria. Limit the review to files listed or implied by the task scope.

### 2. Static pattern sweep

Scan all in-scope source and manifest files for the patterns in the Rules section. Do not execute code, run tools, or fetch live vulnerability databases.

### 3. Classify findings

For each match, assign severity (`critical` / `high` / `medium` / `low` / `informational`) and an OWASP category when applicable.

### 4. Report

Emit findings in the output format. Include file path, line, matched pattern, evidence snippet, and a concrete fix. Do not modify files during review.

### 5. Record in mission evidence

Add findings to the task output or `evidence.jsonl` notes.

## Rules

### General

- **Must**: Review source and manifests only. No execution, no dynamic scans.
- **Must not**: Run any runtime or dynamic vulnerability audit tool, dependency scanner, or live CVE database query.
- **Must**: Treat every hardcoded secret as critical until proven otherwise.
- **Must**: Prefer false positives over missed negatives on the first sweep; triage during classification.
- **Must**: Use `spec.md` acceptance criteria and `plan.json` scope to decide whether a finding is relevant.

### OWASP pattern detection

- **Must** scan for:
  - **Injection** — SQL, command, LDAP, XPath, template injection
  - **Broken authentication** — weak password checks, missing MFA patterns, session tokens committed to source
  - **Sensitive data exposure** — PII regexes, plaintext password storage, weak cryptography
  - **XML external entities** — `DOCTYPE`, `ENTITY`, external DTD references
  - **Broken access control** — missing authorization checks, path traversal patterns
  - **Security misconfiguration** — default credentials, debug flags, CORS wildcard, verbose error leaks
  - **Cross-site scripting** — unescaped output in HTML, inline event handlers, `innerHTML` with user input
  - **Insecure deserialization** — `eval`, `pickle.loads`, `ObjectInputStream`, YAML `load`
  - **Known vulnerable components** — manifest patterns below
  - **Insufficient logging/monitoring** — no log/audit traces around authentication, authorization, or payment flows

### Hardcoded secrets/keys/tokens

- **Must** flag these patterns in source or config:
  - Private keys: `-----BEGIN.*PRIVATE KEY-----`
  - Password assignments: `password\s*=\s*['"]`
  - API key assignments: `api[_]?key\s*=\s*['"]`
  - Secret assignments: `secret\s*=\s*['"][^'"]{8,}`
  - Token assignments: `token\s*=\s*['"][^'"]{8,}`
  - Bearer tokens: `Bearer\s+[A-Za-z0-9\-_]{20,}`
  - JWT secrets or HMAC keys in config or source
  - Database connection strings with embedded credentials
- **Must not** ignore test fixtures unless they are explicitly documented as fake/non-sensitive.

### SQL and command injection

- **Must** flag:
  - String concatenation in SQL: `"SELECT * FROM " + table`, `+ userInput`, template literals inside query strings
  - Dynamic table/column names in queries without an allowlist
  - Shell command building: `exec(`, `spawn(`, `os.system(`, `subprocess` with string arguments
  - User input flowing unsanitized into queries or commands
- **Must** note: parameterized queries help for values, but identifiers (tables/columns) still require allowlisting.

### Manifest and lock files

- **Must** read in-scope manifests statically:
  - `package.json` / `package-lock.json` — deprecated packages, very old version ranges, extraneous `resolutions`/`overrides` that hide versions
  - `go.mod` / `go.sum` — `replace` directives that may mask vulnerable versions, local replace paths
  - `requirements.txt` / `Pipfile.lock` / `Cargo.lock` — package versions matching known CVE patterns, yanked or deprecated package names
- **Must not**: treat absence of a pattern match as a clean bill of health; this is heuristic only.

## Out of scope

- Dynamic or runtime security testing
- Penetration testing, fuzzing, or network scanning
- Social engineering or physical security
- Infrastructure hardening beyond source/manifest patterns
- General code quality or architecture review
- Mission lifecycle management
- Writing fixes — report only; leave implementation to the build step

## Output format

```
Scope: <files reviewed>
Patterns: <list of pattern categories used>
Findings:
  - Severity: critical/high/medium/low/info
    Category: <OWASP category or pattern type>
    File: <path>
    Line: <number>
    Pattern: <matched pattern or heuristic>
    Evidence: <short snippet>
    Fix: <concrete remediation>
Summary:
  Total: <n>
  Critical: <n>
  By category: <counts>
Recommendation: pass / fix-before-merge / block
```

## Checklist

- [ ] Mission scope confirmed from `spec.md` and `plan.json`
- [ ] All source files in task scope reviewed for OWASP patterns
- [ ] Hardcoded secrets scan completed on source and config files
- [ ] SQL/command injection patterns checked
- [ ] Manifest and lock files reviewed for suspicious version patterns
- [ ] Findings classified by severity and OWASP category
- [ ] Each finding includes file path, line, evidence snippet, and fix guidance
- [ ] No dynamic audit tools executed
- [ ] Results recorded in task output or `evidence.jsonl`

## References

- OWASP Top 10 (2021) — https://owasp.org/Top10/
- CWE/SANS Top 25 — https://cwe.mitre.org/top25/
- OWASP Cheat Sheet Series — https://cheatsheetseries.owasp.org/
- Mission artifacts: `spec.md`, `plan.json`, `evidence.jsonl`
