package compact

import (
	"path/filepath"
	"strings"
)

// FilterLs compacts directory listing output.
// Produces tree-like grouping by directory, strips . and ..
// Handles both simple listing (ls) and detailed listing (ls -la).
type FilterLs struct{}

func (FilterLs) Apply(stdout string) string {
	lines := strings.Split(stdout, "\n")
	if len(lines) == 0 {
		return ""
	}

	// Detect if this is detailed output (has permission strings like "drwxr-xr-x")
	// or also totals line like "total 123".
	isDetailed := false
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}
		if looksLikeDetailLine(trimmed) {
			isDetailed = true
			break
		}
	}

	return buildLsTree(lines, isDetailed)
}

// looksLikeDetailLine checks if a line starts with a file type + permissions pattern.
func looksLikeDetailLine(line string) bool {
	if len(line) < 10 {
		return false
	}
	// Permission pattern: [-dlbcps][rwx-]{9}
	first := line[0]
	if first != '-' && first != 'd' && first != 'l' && first != 'b' && first != 'c' && first != 'p' && first != 's' {
		return false
	}
	for i := 1; i < 10 && i < len(line); i++ {
		ch := line[i]
		if ch != 'r' && ch != 'w' && ch != 'x' && ch != '-' && ch != 's' && ch != 'S' && ch != 't' && ch != 'T' {
			return false
		}
	}
	return true
}

func buildLsTree(lines []string, isDetailed bool) string {
	var entries []string

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}
		// Skip "total" lines in detailed output.
		if isDetailed && strings.HasPrefix(trimmed, "total ") {
			continue
		}
		// Skip . and .. entries in any output format.
		if trimmed == "." || trimmed == ".." {
			continue
		}
		// For permission-denied or error lines, preserve them.
		if strings.Contains(trimmed, "Permission denied") ||
			strings.Contains(trimmed, "cannot access") ||
			strings.Contains(trimmed, "No such file") {
			entries = append(entries, trimmed)
			continue
		}
		// If detailed, extract the filename (last field).
		if isDetailed {
			parts := strings.Fields(trimmed)
			if len(parts) >= 9 {
				// Handle symlink "->" notation.
				name := parts[8]
				for i := 9; i < len(parts); i++ {
					name += " " + parts[i]
				}
				// Skip . and .. after extraction.
				if name == "." || name == ".." {
					continue
				}
				entries = append(entries, name)
			}
		} else {
			// Skip . and .. in simple listing.
			if trimmed == "." || trimmed == ".." {
				continue
			}
			entries = append(entries, trimmed)
		}
	}

	if len(entries) == 0 {
		return "empty"
	}
	return groupByDir(entries)
}

// groupByDir groups entries by parent directory for compact tree display.
func groupByDir(entries []string) string {
	dirs := make(map[string][]string)
	var rootFiles []string

	for _, e := range entries {
		// Strip trailing / or * from ls output.
		e = strings.TrimRight(e, "/*@=|>&")

		dir := filepath.Dir(e)
		base := filepath.Base(e)
		if dir == "." {
			rootFiles = append(rootFiles, base)
		} else {
			dirs[dir] = append(dirs[dir], base)
		}
	}

	var result strings.Builder

	if len(rootFiles) > 0 {
		for _, f := range rootFiles {
			result.WriteString(f)
			result.WriteByte('\n')
		}
	}

	for dir, names := range dirs {
		if len(names) == 1 && names[0] == filepath.Base(dir) {
			result.WriteString(dir + "/")
			result.WriteByte('\n')
		} else {
			result.WriteString(dir + "/ {")
			result.WriteString(strings.Join(names, ", "))
			result.WriteString("}\n")
		}
	}

	return strings.TrimRight(result.String(), "\n")
}

// FilterCat compacts `cat` output.
// Strips blank lines and standalone comment lines (//, #),
// preserves all code and error lines.
type FilterCat struct{}

func (FilterCat) Apply(stdout string) string {
	lines := strings.Split(stdout, "\n")
	var result strings.Builder
	blankCount := 0

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		// Track blank lines.
		if trimmed == "" {
			blankCount++
			continue
		}

		// Skip standalone comment lines.
		if strings.HasPrefix(trimmed, "//") || strings.HasPrefix(trimmed, "#") {
			continue
		}

		result.WriteString(line)
		result.WriteByte('\n')
	}

	out := strings.TrimRight(result.String(), "\n")
	if blankCount == len(lines) || (out == "" && blankCount > 0) {
		return "empty"
	}
	return out
}
