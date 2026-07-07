package main

import (
	"strings"
	"testing"
	"time"
)

func TestShortTimeIdPrefix(t *testing.T) {
	date := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC) // epoch
	id := shortTimeId("M", date)
	if !strings.HasPrefix(id, "M") {
		t.Errorf("id should start with prefix M, got %s", id)
	}
	if len(id) != 9 { // M + 8 chars
		t.Errorf("id should be 9 chars (1 prefix + 8 encoded), got %d: %s", len(id), id)
	}
}

func TestShortTimeIdWidth(t *testing.T) {
	// At epoch, offset is 0, encoded is "00000000"
	date := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	id := shortTimeId("E", date)
	if len(id) != 9 {
		t.Errorf("epoch id should be 9 chars, got %d: %s", len(id), id)
	}
	// Should have 8 characters after prefix
	encoded := id[1:]
	if len(encoded) != 8 {
		t.Errorf("encoded part should be 8 chars, got %d: %s", len(encoded), encoded)
	}
}

func TestShortTimeIdIncrement(t *testing.T) {
	d1 := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	d2 := d1.Add(time.Hour)

	id1 := shortTimeId("M", d1)
	id2 := shortTimeId("M", d2)

	if id1 >= id2 {
		t.Errorf("later date should produce lexicographically larger id: %s >= %s", id1, id2)
	}
}

func TestShortTimeIdLexicographic(t *testing.T) {
	// Verify that base36 encoding preserves lexicographic order for increasing values
	times := []time.Time{
		time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2026, 6, 15, 12, 30, 0, 0, time.UTC),
		time.Date(2027, 1, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2028, 1, 1, 0, 0, 0, 0, time.UTC),
	}

	var ids []string
	for _, d := range times {
		ids = append(ids, shortTimeId("M", d))
	}

	for i := 1; i < len(ids); i++ {
		if ids[i-1] >= ids[i] {
			t.Errorf("ids should be strictly increasing: %s >= %s at index %d", ids[i-1], ids[i], i)
		}
	}
}

func TestMissionId(t *testing.T) {
	date := time.Date(2026, 7, 7, 0, 0, 0, 0, time.UTC)
	id := missionId(date)
	if !strings.HasPrefix(id, "M") {
		t.Errorf("mission id should start with M, got %s", id)
	}
	if len(id) != 9 {
		t.Errorf("mission id should be 9 chars, got %d: %s", len(id), id)
	}
}

func TestEvidenceId(t *testing.T) {
	date := time.Date(2026, 7, 7, 0, 0, 0, 0, time.UTC)
	id := evidenceId(date)
	if !strings.HasPrefix(id, "E") {
		t.Errorf("evidence id should start with E, got %s", id)
	}
	if len(id) != 9 {
		t.Errorf("evidence id should be 9 chars, got %d: %s", len(id), id)
	}
}

func TestShortTimeIdNoHyphen(t *testing.T) {
	// Verify that ids never contain hyphens (compact format)
	date := time.Date(2026, 7, 7, 0, 0, 0, 0, time.UTC)
	for i := 0; i < 10; i++ {
		id := shortTimeId("M", date.Add(time.Duration(i)*time.Hour))
		if strings.Contains(id, "-") {
			t.Errorf("id should not contain hyphen: %s", id)
		}
	}
}

func TestShortTimeIdUpperCase(t *testing.T) {
	date := time.Date(2026, 7, 7, 0, 0, 0, 0, time.UTC)
	id := shortTimeId("M", date)
	// base36 encoding produces uppercase letters only
	for _, c := range id[1:] {
		if c >= 'a' && c <= 'z' {
			t.Errorf("id should be uppercase, got lowercase '%c' in %s", c, id)
		}
	}
}

func TestMissionIdDistinctEvidenceId(t *testing.T) {
	date := time.Date(2026, 7, 7, 0, 0, 0, 0, time.UTC)
	mID := missionId(date)
	eID := evidenceId(date)
	if mID[0] != 'M' || eID[0] != 'E' {
		t.Errorf("mission id should start with M and evidence with E, got %s and %s", mID, eID)
	}
	// After prefix removal, the encoded parts should be identical (same timestamp)
	if mID[1:] != eID[1:] {
		t.Errorf("encoded parts should match for same timestamp: %s vs %s", mID[1:], eID[1:])
	}
}
