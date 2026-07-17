# Mission Templates

## spec.md Template

```markdown
# Mission Spec

## Goal
What this mission accomplishes and why it matters.

## User-visible behavior
- Observable changes from the user's perspective
- API changes, CLI flags, UI updates

## Non-goals
- What this mission explicitly does NOT include
- Deferred features or out-of-scope items

## Constraints
- Backward compatibility requirements
- Performance budgets
- Security considerations
- Dependencies on other systems

## Edge Cases
- Empty input: how should the system behave?
- Concurrent access: what if multiple users/agents act simultaneously?
- Partial failures: what happens when a subtask fails mid-execution?
- Boundary conditions: max/min values, zero-length inputs, null states

## Error Handling
- Expected failure modes and their user-facing messages
- Retry strategies for transient failures
- Fallback behavior when primary path fails
- Error logging and debugging information

## Integration Points
- External APIs or services this mission interacts with
- Database schemas or migrations required
- File system paths or configurations affected
- Dependencies on other missions or features

## Test Plan
- Unit tests: which functions/modules need isolated testing?
- Integration tests: which cross-component flows need verification?
- E2E tests: what user scenarios must pass?
- Edge case tests: which boundary conditions need coverage?
- Performance tests: what load/latency requirements must be met?

## Acceptance checks
- [ ] Check 1: specific, verifiable condition
- [ ] Check 2: specific, verifiable condition
- [ ] Check 3: specific, verifiable condition
```

## AI Review Checklist

When performing AI-powered code review, check for:

### Hallucination Detection
- [ ] Verify all imported packages exist in official registries (pkg.go.dev, npmjs.com, pypi.org)
- [ ] Check that referenced functions/methods actually exist in the imported packages
- [ ] Validate that configuration keys match the actual library API
- [ ] Confirm that error codes and status codes are real, not invented

### Security Anti-Patterns
- [ ] SQL queries use parameterized statements, not string formatting
- [ ] User input is validated and sanitized before use
- [ ] Secrets are not hardcoded or logged
- [ ] HTTP handlers include authentication/authorization checks
- [ ] File operations validate paths to prevent directory traversal

### Error Handling
- [ ] All error-returning functions have their errors checked
- [ ] Errors are propagated or handled appropriately (not silently ignored)
- [ ] Error messages provide actionable information for debugging
- [ ] Resource cleanup happens in error paths (defer, finally, etc.)

### Code Quality
- [ ] Functions have single responsibility
- [ ] Magic numbers are replaced with named constants
- [ ] Complex logic has explanatory comments
- [ ] Tests cover happy path, error paths, and edge cases

## plan.json Template

```json
{
  "missionId": "M07ABC123",
  "tasks": [
    {
      "id": "T1",
      "title": "Task title",
      "status": "pending",
      "dependsOn": [],
      "files": ["path/to/file.go"],
      "acceptance": [
        "Specific, verifiable acceptance criterion"
      ],
      "verify": "command to verify task completion",
      "evidence": "evidence-label"
    }
  ]
}
```

## evidence.jsonl Format

Each line is a JSON object:
```json
{
  "id": "E07ABC123",
  "label": "task-verification",
  "command": "go test ./... -v",
  "exitCode": 0,
  "stdout": "test output...",
  "stderr": "",
  "createdAt": "2026-07-14T12:00:00Z"
}
```
