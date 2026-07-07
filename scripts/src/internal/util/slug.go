package util

import (
	"regexp"
	"strings"
	"time"
)

// Slugify converts a string into a URL-safe slug.
func Slugify(value string) string {
	text := strings.ToLower(strings.TrimSpace(value))
	var sb strings.Builder
	for _, r := range text {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') {
			sb.WriteRune(r)
		} else {
			sb.WriteRune('-')
		}
	}
	slug := regexp.MustCompile(`-+`).ReplaceAllString(sb.String(), "-")
	slug = strings.Trim(slug, "-")
	if len(slug) > 60 {
		slug = slug[:60]
		slug = strings.TrimRight(slug, "-")
	}
	if slug == "" {
		return "mission"
	}
	return slug
}

// CommandToString joins command parts into a safe shell string.
func CommandToString(parts []string) string {
	safeRegex := regexp.MustCompile(`^[A-Za-z0-9_./:=@%+-]+$`)
	var result []string
	for _, part := range parts {
		if safeRegex.MatchString(part) {
			result = append(result, part)
		} else {
			escaped := strings.ReplaceAll(part, "'", "'\\''")
			result = append(result, "'"+escaped+"'")
		}
	}
	return strings.Join(result, " ")
}

// IsoNow returns the current UTC time in ISO 8601 format with millisecond precision.
func IsoNow() string {
	return time.Now().UTC().Format("2006-01-02T15:04:05.000Z")
}

var legacyMissionRegex = regexp.MustCompile(`\b[Mm]-(\d{8}-\d{6})\b`)
var compactMissionRegex = regexp.MustCompile(`(?:^|[^A-Za-z0-9])([Mm][0-9A-Za-z]{8})(?:$|[^A-Za-z0-9])`)

// NormalizeMissionId extracts a normalized mission id from text.
// Supports both legacy (M-YYYYMMDD-HHmmss) and compact (Mxxxxxxxx) formats.
// Returns nil if no valid mission id is found.
func NormalizeMissionId(value string) *string {
	text := strings.TrimSpace(value)
	if legacy := legacyMissionRegex.FindStringSubmatch(text); legacy != nil {
		res := "M-" + legacy[1]
		return &res
	}
	if compact := compactMissionRegex.FindStringSubmatch(text); compact != nil {
		res := strings.ToUpper(compact[1])
		return &res
	}
	return nil
}

// RegexpReplace replaces all matches of the pattern with repl in src.
func RegexpReplace(pattern, repl, src string) string {
	re := regexp.MustCompile(pattern)
	return re.ReplaceAllString(src, repl)
}

// ContainsStr returns true if slice contains item.
func ContainsStr(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
