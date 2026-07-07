package main

import (
	"strings"
	"encoding/json"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

type EvidenceEntry struct {
	ID        string `json:"id"`
	Label     string `json:"label"`
	Command   string `json:"command"`
	ExitCode  int    `json:"exitCode"`
	Stdout    string `json:"stdout"`
	Stderr    string `json:"stderr"`
	CreatedAt string `json:"createdAt"`
}

type CommandResult struct {
	exitCode int
	stdout   string
	stderr   string
}

func runCommand(commandParts []string) CommandResult {
	if len(commandParts) == 0 {
		return CommandResult{exitCode: 1}
	}
	cmd := exec.Command(commandParts[0], commandParts[1:]...)
	cmd.Dir = ROOT
	cmd.Env = os.Environ()

	var stdoutBuf bytes.Buffer
	var stderrBuf bytes.Buffer
	cmd.Stdout = &stdoutBuf
	cmd.Stderr = &stderrBuf

	err := cmd.Run()
	exitCode := 0
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			exitCode = exitError.ExitCode()
		} else {
			exitCode = 127
			stderrBuf.WriteString(err.Error() + "\n")
		}
	}

	return CommandResult{
		exitCode: exitCode,
		stdout:   stdoutBuf.String(),
		stderr:   stderrBuf.String(),
	}
}

func reserveEvidencePaths(dir string) (string, string, string) {
	base := evidenceId()
	candidate := base
	index := 2
	outputsDir := filepath.Join(dir, "outputs")

	for {
		stdoutPath := filepath.Join(outputsDir, candidate+".stdout.txt")
		stderrPath := filepath.Join(outputsDir, candidate+".stderr.txt")

		f, err := os.OpenFile(stdoutPath, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0644)
		if err == nil {
			f.Close()
			return candidate, stdoutPath, stderrPath
		}
		if !os.IsExist(err) {
			fail("Failed to reserve evidence path: " + err.Error())
		}
		candidate = evidenceId(time.Now().Add(time.Duration(index) * time.Millisecond))
		index++
	}
}

func recordEvidence(args []string) {
	separator := -1
	for i, a := range args {
		if a == "--" {
			separator = i
			break
		}
	}
	if separator == -1 {
		fail("Missing -- before evidence command.\n\n" + usage())
	}

	label := strings.TrimSpace(stringsJoin(args[:separator], " "))
	commandParts := args[separator+1:]

	if label == "" {
		fail("Missing evidence label.\n\n" + usage())
	}
	if len(commandParts) == 0 {
		fail("Missing command after --.\n\n" + usage())
	}

	res := resolveMission("")
	if res.Safety != "safe" || res.Selected == nil {
		fail(formatResolutionBlock(res, "evidence"))
	}

	id := res.Selected.ID
	dir := missionDir(id)
	outputsDir := filepath.Join(dir, "outputs")
	os.MkdirAll(outputsDir, 0755)

	evidence, stdoutPath, stderrPath := reserveEvidencePaths(dir)
	result := runCommand(commandParts)

	os.WriteFile(stdoutPath, []byte(result.stdout), 0644)
	os.WriteFile(stderrPath, []byte(result.stderr), 0644)

	entry := EvidenceEntry{
		ID:        evidence,
		Label:     label,
		Command:   commandToString(commandParts),
		ExitCode:  result.exitCode,
		Stdout:    displayPath(stdoutPath),
		Stderr:    displayPath(stderrPath),
		CreatedAt: isoNow(),
	}

	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(false)
	enc.Encode(entry)

	f, _ := os.OpenFile(filepath.Join(dir, "evidence.jsonl"), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	f.Write(buf.Bytes())
	f.Close()

	fmt.Printf("Evidence: %s\n", evidence)
	fmt.Printf("Exit code: %d\n", result.exitCode)
	os.Exit(result.exitCode)
}
