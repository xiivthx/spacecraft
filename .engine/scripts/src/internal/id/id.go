// Package id generates compact sortable mission and evidence IDs.
//
// IDs are 9-character strings: a single-letter prefix (M for mission, E for evidence)
// followed by an 8-character base36 encoding of milliseconds since ID_EPOCH_MS.
// This produces lexicographically sortable IDs with no hyphens.
package id

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

const ID_EPOCH_MS = 1767225600000 // Date.UTC(2026, 0, 1)

// ShortTimeId generates a prefixed compact time-based id.
// prefix is a single uppercase letter (e.g. "M" for mission, "E" for evidence).
// date is the timestamp to encode; defaults to time.Now().
func ShortTimeId(prefix string, date time.Time) (string, error) {
	offset := date.UnixMilli() - ID_EPOCH_MS
	if offset < 0 {
		return "", fmt.Errorf("cannot create %s id before 2026-01-01T00:00:00.000Z", prefix)
	}
	encoded := strings.ToUpper(strconv.FormatInt(offset, 36))
	if len(encoded) < 8 {
		encoded = strings.Repeat("0", 8-len(encoded)) + encoded
	}
	return prefix + encoded, nil
}

// MissionId generates a compact mission id (e.g. "M07ABCDEF").
// If date is provided, uses it; otherwise uses time.Now().
func MissionId(date ...time.Time) (string, error) {
	d := time.Now()
	if len(date) > 0 {
		d = date[0]
	}
	return ShortTimeId("M", d)
}

// EvidenceId generates a compact evidence id (e.g. "E07ABCDEF").
// If date is provided, uses it; otherwise uses time.Now().
func EvidenceId(date ...time.Time) (string, error) {
	d := time.Now()
	if len(date) > 0 {
		d = date[0]
	}
	return ShortTimeId("E", d)
}
