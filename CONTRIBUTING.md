# Contributing to Spacecraft

Thanks for your interest in contributing! This guide covers development setup, testing, commit conventions, and the PR process.

## Development Setup

### Prerequisites

- Go 1.21+
- Node.js 18+ (for integration tests)
- Git
- macOS or Linux

### Initial Setup

```sh
# Clone the repository
git clone <repo-url>
cd spacecraft

# Build the Go binary
make build

# Verify everything works
make test
scripts/spacecraft help
```

### Project Structure

```
.engine/            # OpenCode engine config, skills, and conventions
  agents/           # sc-* agent definitions
  skills/           # sc-* skills by category
  opencode.json     # Agent config, permissions, models
scripts/            # Go CLI source and binary
  src/              # main.go, types.go, internal/
.space/             # Mission state (created by spacecraft init)
tests/              # Node integration tests
```

### Development Workflow

```sh
# Build after code changes
make build

# Run Go tests
make test

# Run linter
make lint

# Clean build artifacts
make clean
```

## Testing

### Running Tests

```sh
# All tests (Go unit + Node integration)
make test

# Go unit tests only
cd .engine/scripts/src && go test ./...

# Node integration tests only
cd tests && npm test
```

### Writing Tests

- Go unit tests: place `*_test.go` files alongside source in `.engine/scripts/src/`
- Integration tests: add to `tests/` directory
- Test names should describe the behavior being tested
- Include both positive and negative test cases

### Test Coverage

```sh
# Generate coverage report
cd .engine/scripts/src && go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html
```

## Commit Conventions

Spacecraft uses [Conventional Commits](https://www.conventionalcommits.org/).

### Format

```
<type>: <description>

[optional body]

[optional footer]
```

### Types

- `feat` - New features
- `fix` - Bug fixes
- `docs` - Documentation changes
- `style` - Code style (formatting, semicolons, etc.)
- `refactor` - Code refactoring (no feature or bug change)
- `test` - Adding or updating tests
- `chore` - Maintenance tasks (deps, CI, build)

### Examples

```
feat: add mission archive command
fix: resolve branch detection in worktrees
docs: update installation guide for Homebrew
test: add integration tests for evidence capture
refactor: simplify mission state transitions
```

### Scope

Keep commits focused. Target 1-3 commits per feature branch, max 5. Squash WIP before merge.

## Pull Request Process

### Before Submitting

1. **Rebase on main** - Ensure your branch is up to date:
   ```sh
   git fetch origin
   git rebase origin/main
   ```

2. **Run tests** - All tests must pass:
   ```sh
   make test
   ```

3. **Check lint** - No lint errors:
   ```sh
   make lint
   ```

4. **Verify changes** - Test your changes manually if needed

### Submitting

1. Push your branch:
   ```sh
   git push origin feat/<id>/<title>
   ```

2. Open a pull request against `main`

3. Fill out the PR template:
   - What does this change?
   - Why is it needed?
   - How was it tested?
   - Any breaking changes?

### Review Process

- All PRs require review before merge
- Reviewer checks: code quality, test coverage, documentation, commit hygiene
- Address review feedback with new commits (don't force-push during review)
- After approval, maintainer merges with `--no-ff`

### Merge Strategy

- Branches merge to `main` with `git merge --no-ff`
- Merge commits are annotated with mission ID when applicable
- After merge: version tag created, branch deleted

## Code Style

### Go

- Follow standard Go conventions (`gofmt`, `go vet`)
- Use meaningful variable names
- Keep functions focused and small
- Document exported functions and types

### Shell Scripts

- Use `#!/bin/sh` for portability
- Include `set -e` for error handling
- Quote variables: `"$var"` not `$var`
- Test on both macOS and Linux when possible

### Documentation

- Use clear, concise language
- Include code examples where helpful
- Keep examples tested and verified
- Link to related docs when relevant

## Reporting Issues

- Check existing issues before creating a new one
- Include steps to reproduce for bugs
- Include expected vs actual behavior
- Mention your OS and Go version

## Questions?

- Check the [README](./README.md) for overview and usage
- See the [installation guide](./docs/installation.md) for setup help
- Review existing missions in `.space/missions/` for examples
