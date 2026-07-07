package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

const ID_EPOCH_MS = 1767225600000 // Date.UTC(2026, 0, 1)
const SHORT_ID_WIDTH = 8

func shortTimeId(prefix string, date time.Time) string {
	offset := date.UnixMilli() - ID_EPOCH_MS
	if offset < 0 {
		fail(fmt.Sprintf("Cannot create %s id before 2026-01-01T00:00:00.000Z.", prefix))
	}
	encoded := strings.ToUpper(strconv.FormatInt(offset, 36))
	if len(encoded) < SHORT_ID_WIDTH {
		encoded = strings.Repeat("0", SHORT_ID_WIDTH-len(encoded)) + encoded
	}
	return prefix + encoded
}

func missionId(date ...time.Time) string {
	d := time.Now()
	if len(date) > 0 {
		d = date[0]
	}
	return shortTimeId("M", d)
}

func evidenceId(date ...time.Time) string {
	d := time.Now()
	if len(date) > 0 {
		d = date[0]
	}
	return shortTimeId("E", d)
}
