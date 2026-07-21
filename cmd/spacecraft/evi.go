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

	for i := 0; i < len(args); i++ {
		if args[i] == "--mission" && i+1 < len(args) {
			mid = args[i+1]
			i++
		} else if args[i] == "--" {
			cmdArgs = args[i+1:]
			break
		} else if label == "" && !strings.HasPrefix(args[i], "--") {
			label = args[i]
		}
	}

	if mid == "" {
		fmt.Fprintln(os.Stderr, "spacecraft evidence: no active mission - use --mission <id> or run from feat/<id>/ branch")
		return 1
	}
	if label == "" || len(cmdArgs) == 0 {
		fmt.Fprintln(os.Stderr, "Usage: spacecraft evidence [--mission <id>] <label> -- <command...>")
		return 1
	}

	cmd := exec.Command(cmdArgs[0], cmdArgs[1:]...)
	cmd.Dir = filepath.Dir(spaceDir)
	out, err := cmd.CombinedOutput()
	output := string(out)

	exitCode := 0
	if err != nil {
		if ee, ok := err.(*exec.ExitError); ok {
			exitCode = ee.ExitCode()
		} else {
			exitCode = 127
			output += err.Error() + "\n"
		}
	}

	entry := evidenceEntry{
		Label:      label,
		Command:    strings.Join(cmdArgs, " "),
		Output:     output,
		OutputHash: outputSHA256Hex(output),
		ExitCode:   exitCode,
		TS:         time.Now().UTC().Format(time.RFC3339),
	}

	data, _ := json.Marshal(entry)
	evidencePath := filepath.Join(missionDir(spaceDir, mid), "evidence.jsonl")
	os.MkdirAll(filepath.Dir(evidencePath), 0755)

	f, ferr := os.OpenFile(evidencePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if ferr != nil {
		fmt.Fprintln(os.Stderr, "spacecraft evidence:", ferr)
		return 1
	}
	f.Write(data)
	f.Write([]byte("\n"))
	f.Close()

	fmt.Print(output)
	fmt.Printf("Exit code: %d\n", exitCode)

	return exitCode
}

type evidenceEntry struct {
	Label      string `json:"label"`
	Command    string `json:"command"`
	Output     string `json:"output"`
	OutputHash string `json:"outputHash,omitempty"`
	ExitCode   int    `json:"exitCode"`
	TS         string `json:"ts"`
}
