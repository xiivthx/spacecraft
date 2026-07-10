package compact

import "strings"

// FilterDockerPs compacts `docker ps` output.
// Strips header row, abbreviates long image tags (>20 chars), keeps data rows.
type FilterDockerPs struct{}

func (FilterDockerPs) Apply(stdout string) string {
	if stdout == "" {
		return "no containers"
	}
	lines := strings.Split(stdout, "\n")
	var dataLines []string
	headerSeen := false

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}
		// Detect and skip header row (CONTAINER ID header).
		if !headerSeen && strings.HasPrefix(trimmed, "CONTAINER ID") {
			headerSeen = true
			continue
		}
		// Abbreviate long image tag fields (column 2 in docker ps output).
		// docker ps output is column-aligned with 2+ spaces between columns.
		line = abbreviateImageTag(line)
		dataLines = append(dataLines, line)
	}

	if len(dataLines) == 0 {
		return "no containers"
	}
	return strings.Join(dataLines, "\n")
}

// abbreviateImageTag truncates long image tags in the second column.
func abbreviateImageTag(line string) string {
	const maxImageLen = 21
	// Find the image column: starts after CONTAINER ID, separated by 2+ spaces.
	// docker ps format: CONTAINER ID   IMAGE   COMMAND   CREATED   STATUS   PORTS   NAMES
	// We look for the image portion — between the first and second set of spaces.
	fields := splitColumns(line)
	if len(fields) < 2 {
		return line
	}
	image := fields[1]
	if len(image) > maxImageLen {
		image = image[:maxImageLen-3] + "..."
	}
	// Rebuild line with abbreviated image — preserve original spacing.
	return strings.Replace(line, fields[1], image, 1)
}

// splitColumns splits a docker ps row into columns by 2+ space separators.
func splitColumns(line string) []string {
	var cols []string
	start := 0
	inSpace := false
	spaceStart := 0
	for i, ch := range line {
		if ch == ' ' {
			if !inSpace {
				inSpace = true
				spaceStart = i
			}
		} else {
			if inSpace && i-spaceStart >= 2 {
				cols = append(cols, line[start:spaceStart])
				start = i
			}
			inSpace = false
		}
	}
	cols = append(cols, line[start:])
	return cols
}
