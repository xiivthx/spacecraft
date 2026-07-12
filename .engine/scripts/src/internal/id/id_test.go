package id

import (
	"strings"
	"testing"
	"time"
)

func TestShortTimeIdPrefix(t *testing.T) {
	date := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	id, err := ShortTimeId("M", date)
	if err != nil {
		t.Fatalf("ShortTimeId error: %v", err)
	}
	if !strings.HasPrefix(id, "M") {
		t.Errorf("id should start with M, got %s", id)
	}
	if len(id) != 9 {
		t.Errorf("id should be 9 chars, got %d: %s", len(id), id)
	}
}

func TestShortTimeIdWidth(t *testing.T) {
	date := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	id, err := ShortTimeId("E", date)
	if err != nil {
		t.Fatal(err)
	}
	encoded := id[1:]
	if len(encoded) != 8 {
		t.Errorf("encoded part should be 8 chars, got %d: %s", len(encoded), encoded)
	}
}

func TestShortTimeIdIncrement(t *testing.T) {
	d1 := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	d2 := d1.Add(time.Hour)

	id1, _ := ShortTimeId("M", d1)
	id2, _ := ShortTimeId("M", d2)

	if id1 >= id2 {
		t.Errorf("later date should produce larger id: %s >= %s", id1, id2)
	}
}

func TestShortTimeIdLexicographic(t *testing.T) {
	times := []time.Time{
		time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2026, 6, 15, 12, 30, 0, 0, time.UTC),
		time.Date(2027, 1, 1, 0, 0, 0, 0, time.UTC),
	}

	var ids []string
	for _, d := range times {
		id, _ := ShortTimeId("M", d)
		ids = append(ids, id)
	}

	for i := 1; i < len(ids); i++ {
		if ids[i-1] >= ids[i] {
			t.Errorf("ids should be increasing: %s >= %s", ids[i-1], ids[i])
		}
	}
}

func TestMissionIdPrefix(t *testing.T) {
	date := time.Date(2026, 7, 7, 0, 0, 0, 0, time.UTC)
	id, err := MissionId(date)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(id, "M") {
		t.Errorf("mission id should start with M, got %s", id)
	}
	if len(id) != 9 {
		t.Errorf("mission id should be 9 chars, got %d: %s", len(id), id)
	}
}

func TestEvidenceIdPrefix(t *testing.T) {
	date := time.Date(2026, 7, 7, 0, 0, 0, 0, time.UTC)
	id, err := EvidenceId(date)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(id, "E") {
		t.Errorf("evidence id should start with E, got %s", id)
	}
}

func TestShortTimeIdNoHyphen(t *testing.T) {
	date := time.Date(2026, 7, 7, 0, 0, 0, 0, time.UTC)
	for i := 0; i < 10; i++ {
		id, err := ShortTimeId("M", date.Add(time.Duration(i)*time.Hour))
		if err != nil {
			t.Fatal(err)
		}
		if strings.Contains(id, "-") {
			t.Errorf("id should not contain hyphen: %s", id)
		}
	}
}

func TestShortTimeIdUpperCase(t *testing.T) {
	date := time.Date(2026, 7, 7, 0, 0, 0, 0, time.UTC)
	id, err := ShortTimeId("M", date)
	if err != nil {
		t.Fatal(err)
	}
	for _, c := range id[1:] {
		if c >= 'a' && c <= 'z' {
			t.Errorf("id should be uppercase, got '%c' in %s", c, id)
		}
	}
}

func TestMissionIdDistinctEvidenceId(t *testing.T) {
	date := time.Date(2026, 7, 7, 0, 0, 0, 0, time.UTC)
	mID, _ := MissionId(date)
	eID, _ := EvidenceId(date)
	if mID[0] != 'M' || eID[0] != 'E' {
		t.Errorf("prefix mismatch: %s vs %s", mID, eID)
	}
	if mID[1:] != eID[1:] {
		t.Errorf("encoded parts should match: %s vs %s", mID[1:], eID[1:])
	}
}

func TestShortTimeIdBeforeEpoch(t *testing.T) {
	date := time.Date(2025, 12, 31, 0, 0, 0, 0, time.UTC)
	_, err := ShortTimeId("M", date)
	if err == nil {
		t.Error("should error for date before epoch")
	}
}

func TestShortTimeIdAtEpoch(t *testing.T) {
	date := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	id, err := ShortTimeId("M", date)
	if err != nil {
		t.Fatal(err)
	}
	// At epoch, offset is 0, encoded is "00000000"
	if id != "M00000000" {
		t.Errorf("epoch id should be M00000000, got %s", id)
	}
}
