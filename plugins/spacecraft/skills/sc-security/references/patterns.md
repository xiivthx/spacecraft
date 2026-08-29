# Security pattern detection (heuristic)

On-demand static patterns for `sc-security`. Domain checklist / severity mapping: `.cursor/rules/300-security.mdc`. Finding craft: `.cursor/skills/sc-run/references/defect-finding.md`. Load this file when running a heuristic sweep.

## OWASP pattern detection

Scan for:

- **Injection** - SQL, command, LDAP, XPath, template injection
- **Broken authentication** - weak password checks, missing MFA patterns, session tokens committed to source
- **Sensitive data exposure** - PII regexes, plaintext password storage, weak cryptography
- **XML external entities** - `DOCTYPE`, `ENTITY`, external DTD references
- **Broken access control** - missing authorization checks, path traversal patterns
- **Security misconfiguration** - default credentials, debug flags, CORS wildcard, verbose error leaks
- **Cross-site scripting** - unescaped output in HTML, inline event handlers, `innerHTML` with user input
- **Insecure deserialization** - `eval`, `pickle.loads`, `ObjectInputStream`, YAML `load`
- **Known vulnerable components** - manifest patterns below
- **Insufficient logging/monitoring** - no log/audit traces around authentication, authorization, or payment flows

## Hardcoded secrets/keys/tokens

Flag these patterns in source or config:

- Private keys: `-----BEGIN.*PRIVATE KEY-----`
- Password assignments: `password\s*=\s*['"]`
- API key assignments: `api[_]?key\s*=\s*['"]`
- Secret assignments: `secret\s*=\s*['"][^'"]{8,}`
- Token assignments: `token\s*=\s*['"][^'"]{8,}`
- Bearer tokens: `Bearer\s+[A-Za-z0-9\-_]{20,}`
- JWT secrets or HMAC keys in config or source
- Database connection strings with embedded credentials

Do not ignore test fixtures unless they are explicitly documented as fake/non-sensitive.

## SQL and command injection

Flag:

- String concatenation in SQL: `"SELECT * FROM " + table`, `+ userInput`, template literals inside query strings
- Dynamic table/column names in queries without an allowlist
- Shell command building: `exec(`, `spawn(`, `os.system(`, `subprocess` with string arguments
- User input flowing unsanitized into queries or commands

Note: parameterized queries help for values, but identifiers (tables/columns) still require allowlisting.

## Manifest and lock files

Read in-scope manifests statically:

- `package.json` / `package-lock.json` - deprecated packages, very old version ranges, extraneous `resolutions`/`overrides` that hide versions
- `go.mod` / `go.sum` - `replace` directives that may mask vulnerable versions, local replace paths
- `requirements.txt` / `Pipfile.lock` / `Cargo.lock` - package versions matching known CVE patterns, yanked or deprecated package names

Absence of a pattern match is not a clean bill of health; this is heuristic only.

## Related

- Rule: `.cursor/rules/300-security.mdc`
- Finding craft: `.cursor/skills/sc-run/references/defect-finding.md`
- Skill: `../SKILL.md`
