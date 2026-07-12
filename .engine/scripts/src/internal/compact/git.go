package compact

import (
	"path/filepath"
	"strings"
)

// FilterGitStatus compacts `git status` output.
// Strips untracked hints, groups modified/staged files by directory,
// removes boilerplate help text.
type FilterGitStatus struct{}

func (FilterGitStatus) Apply(stdout string) string {
	var result strings.Builder
	lines := strings.Split(stdout, "\n")

	var staged, unstaged []string
	inStaged := false
	inUnstaged := false
	inUntracked := false

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		// Detect sections.
		if strings.Contains(trimmed, "Changes to be committed:") ||
			strings.Contains(trimmed, "Staged changes:") {
			inStaged = true
			inUnstaged = false
			inUntracked = false
			continue
		}
		if strings.Contains(trimmed, "Changes not staged for commit:") ||
			strings.Contains(trimmed, "Unstaged changes:") {
			inStaged = false
			inUnstaged = true
			inUntracked = false
			continue
		}
		if strings.Contains(trimmed, "Untracked files:") {
			inStaged = false
			inUnstaged = false
			inUntracked = true
			continue
		}

		// Skip section headers, hints, empty lines.
		if strings.HasPrefix(trimmed, "(") && strings.HasSuffix(trimmed, ")") {
			continue
		}
		if strings.HasPrefix(line, "#") {
			continue
		}
		if trimmed == "no changes added to commit" ||
			strings.Contains(trimmed, "use \"git add\"") ||
			strings.Contains(trimmed, "use \"git restore\"") ||
			strings.Contains(trimmed, "use \"git checkout\"") ||
			strings.Contains(trimmed, "use \"git rm\"") ||
			strings.Contains(trimmed, "nothing added to commit") ||
			strings.Contains(trimmed, "nothing to commit") ||
			strings.Contains(trimmed, "no changes") {
			continue
		}
		// Skip "modified:" prefix to keep just the file path.
		cleaned := trimmed
		for _, prefix := range []string{
			"modified:   ", "modified:", "new file:   ", "new file:",
			"deleted:    ", "deleted:", "renamed:    ", "renamed:",
		} {
			if strings.HasPrefix(cleaned, prefix) {
				cleaned = strings.TrimSpace(strings.TrimPrefix(cleaned, prefix))
				break
			}
		}
		if cleaned == "" {
			continue
		}

		if inStaged {
			staged = append(staged, cleaned)
		} else if inUnstaged {
			unstaged = append(unstaged, cleaned)
		} else if inUntracked {
			continue // skip untracked — high noise, low signal for LLM
		}
	}

	// Output grouped by directory.
	if len(staged) > 0 {
		result.WriteString("staged:\n")
		writeGrouped(&result, staged)
	}
	if len(unstaged) > 0 {
		if result.Len() > 0 {
			result.WriteString("\n")
		}
		result.WriteString("unstaged:\n")
		writeGrouped(&result, unstaged)
	}

	out := result.String()
	if out == "" {
		return "clean"
	}
	return strings.TrimRight(out, "\n")
}

// FilterGitDiff compacts `git diff` output.
// Keeps only +/- content lines, strips meta headers (index, @@, ---, +++).
type FilterGitDiff struct{}

func (FilterGitDiff) Apply(stdout string) string {
	var result strings.Builder
	lines := strings.Split(stdout, "\n")
	kept := 0

	for _, line := range lines {
		if line == "" {
			continue
		}
		// Keep content lines only.
		switch {
		case strings.HasPrefix(line, "--- ") || strings.HasPrefix(line, "+++ "):
			// Skip file path meta lines (before +/- check).
			continue
		case strings.HasPrefix(line, "+") || strings.HasPrefix(line, "-"):
			result.WriteString(line)
			result.WriteByte('\n')
			kept++
		case strings.HasPrefix(line, "diff --git"):
			if kept > 0 {
				result.WriteString("---\n")
			}
			// Extract just the file path.
			parts := strings.Fields(line)
			if len(parts) >= 4 {
				result.WriteString(strings.TrimPrefix(parts[3], "b/"))
				result.WriteByte('\n')
			}
		}
		// Skip index, @@, ---, +++, and other meta lines.
	}

	out := result.String()
	if out == "" {
		return "no changes"
	}
	return strings.TrimRight(out, "\n")
}

// FilterGitLog compacts `git log` output to one-line SHA+msg per commit.
type FilterGitLog struct{}

func (FilterGitLog) Apply(stdout string) string {
	var result strings.Builder
	lines := strings.Split(stdout, "\n")

	// Check if the output is already in oneline format.
	// Oneline: every non-empty line has a short hex SHA prefix followed by space.
	nonEmpty := 0
	onelineCount := 0
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}
		nonEmpty++
		// Check for SHA-like prefix: 7-12 hex chars followed by space.
		if len(trimmed) >= 8 {
			parts := strings.SplitN(trimmed, " ", 2)
			if isHexString(parts[0]) && len(parts[0]) >= 7 {
				onelineCount++
			}
		}
	}
	allShort := nonEmpty > 0 && onelineCount == nonEmpty

	if allShort {
		// Oneline format (e.g., from `git log --oneline`) — passthrough.
		for _, line := range lines {
			trimmed := strings.TrimSpace(line)
			if trimmed == "" {
				continue
			}
			result.WriteString(trimmed)
			result.WriteByte('\n')
		}
		return strings.TrimRight(result.String(), "\n")
	}

	// Standard format: strip Author/Date, keep commit SHA + message.
	var sha string
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || strings.HasPrefix(trimmed, "Author:") || strings.HasPrefix(trimmed, "Date:") {
			continue
		}
		if strings.HasPrefix(trimmed, "commit ") {
			fields := strings.Fields(trimmed)
			if len(fields) >= 2 {
				sha = fields[1]
				if len(sha) > 12 {
					sha = sha[:12]
				}
			}
		} else if sha != "" {
			result.WriteString(sha)
			result.WriteString(" ")
			result.WriteString(trimmed)
			result.WriteByte('\n')
			sha = ""
		}
	}
	return strings.TrimRight(result.String(), "\n")
}

// isHexString returns true if s contains only hex digits (0-9, a-f, A-F).
func isHexString(s string) bool {
	for _, r := range s {
		if !((r >= '0' && r <= '9') || (r >= 'a' && r <= 'f') || (r >= 'A' && r <= 'F')) {
			return false
		}
	}
	return len(s) > 0
}

// writeGrouped groups file paths by their parent directory for compact display.
func writeGrouped(b *strings.Builder, files []string) {
	dirs := make(map[string][]string)
	var rootFiles []string

	for _, f := range files {
		dir := filepath.Dir(f)
		if dir == "." {
			rootFiles = append(rootFiles, f)
		} else {
			dirs[dir] = append(dirs[dir], filepath.Base(f))
		}
	}

	for _, f := range rootFiles {
		b.WriteString("  ")
		b.WriteString(f)
		b.WriteByte('\n')
	}
	for dir, names := range dirs {
		if len(names) == 1 && dir == names[0] {
			b.WriteString("  ")
			b.WriteString(dir)
			b.WriteByte('\n')
		} else {
			b.WriteString("  ")
			b.WriteString(dir)
			b.WriteString("/ {")
			b.WriteString(strings.Join(names, ", "))
			b.WriteString("}\n")
		}
	}
}
