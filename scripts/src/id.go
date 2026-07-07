package main

import (
	"time"

	"spacecraft/internal/id"
)

func missionId(date ...time.Time) string {
	d := time.Now()
	if len(date) > 0 {
		d = date[0]
	}
	v, err := id.MissionId(d)
	if err != nil {
		fail(err.Error())
	}
	return v
}

func evidenceId(date ...time.Time) string {
	d := time.Now()
	if len(date) > 0 {
		d = date[0]
	}
	v, err := id.EvidenceId(d)
	if err != nil {
		fail(err.Error())
	}
	return v
}
