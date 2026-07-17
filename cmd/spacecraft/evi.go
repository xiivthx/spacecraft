package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

func eviCmd(args []string, spaceDir, mid string) int {
	var label string
	var cmdArgs []string

	// Parse --mission flag
	for i := 0; i < len(args); i++ {
		if args[i] == "--mission" && i+1 < len(args) {
			mid = args[i+1]
			i++
		} else if label == "" && !strings.HasPrefix(args[i], "--") {
			label = args[i]
		} else if args[i] == "--" {
			cmdArgs = args[i+1:]
			break
		}
	}

	if mid == "" {
		fmt.Fprintln(os.Stderr, "spacecraft evi: no active mission — use --mission <id> or run from feat/<id>/ branch")
		return 1
	}
	if label == "" || len(cmdArgs) == 0 {
		fmt.Fprintln(os.Stderr, "Usage: spacecraft evi [--mission <id>] <label> -- <command...>")
		return 1
	}
	cmd := exec.Command(cmdArgs[0], cmdArgs[1:]...)
	cmd.Dir = filepath.Dir(spaceDir)
	out, err := cmd.CombinedOutput()
	output := string(out)

	entry := evidenceEntry{
		Label:   label,
		Command: strings.Join(cmdArgs, " "),
		Output:  output,
		TS:      time.Now().UTC().Format(time.RFC3339),
	}

	data, _ := json.Marshal(entry)
	evidencePath := filepath.Join(missionDir(spaceDir, mid), "evidence.jsonl")
	os.MkdirAll(filepath.Dir(evidencePath), 0755)

	f, err := os.OpenFile(evidencePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Fprintln(os.Stderr, "spacecraft evi:", err)
		return 1
	}
	defer f.Close()

	f.Write(data)
	f.Write([]byte("\n"))

	fmt.Print(output)
	return 0
}

type evidenceEntry struct {
	Label   string `json:"label"`
	Command string `json:"command"`
	Output  string `json:"output"`
	TS      string `json:"ts"`
}
