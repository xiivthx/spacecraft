package compact

import "strings"

// FilterCurl compacts `curl -v` verbose output.
// Strips TLS handshake/connection verbose lines (* prefix),
// keeps request headers (> prefix), response headers (< prefix), and body.
type FilterCurl struct{}

func (FilterCurl) Apply(stdout string) string {
	if stdout == "" {
		return ""
	}
	lines := strings.Split(stdout, "\n")
	var kept []string
	for _, line := range lines {
		// Strip verbose connection/TLS lines (start with "* ").
		if strings.HasPrefix(line, "* ") {
			continue
		}
		kept = append(kept, line)
	}
	if len(kept) == 0 || strings.TrimSpace(strings.Join(kept, "")) == "" {
		return "connection failed"
	}
	return strings.Join(kept, "\n")
}
