// Package gitutil provides git information fetching for the Spacecraft CLI.
package gitutil

import (
	"os/exec"
	"strings"

	"spacecraft/internal/mission"
)

// CommandRunner is an interface for running commands, allowing tests to mock git.
type CommandRunner interface {
	Run(name string, args ...string) (exitCode int, stdout, stderr string)
}

// OSCommandRunner runs commands using os/exec.
type OSCommandRunner struct{}

func (OSCommandRunner) Run(name string, args ...string) (int, string, string) {
	cmd := exec.Command(name, args...)
	stdout, err := cmd.Output()
	exitCode := 0
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			exitCode = exitErr.ExitCode()
		} else {
			exitCode = 127
		}
	}
	return exitCode, string(stdout), ""
}

// GitInfo fetches git repository information using the provided runner.
// Returns a GitInfoData with Available, IsRepo, Root, Branch, Sha, Dirty, and DirtyFiles.
func GitInfo(runner CommandRunner) mission.GitInfoData {
	code, out, _ := runner.Run("git", "rev-parse", "--is-inside-work-tree")
	if code != 0 || strings.TrimSpace(out) != "true" {
		available := code != 127
		return mission.GitInfoData{Available: available, IsRepo: false}
	}

	_, rootOut, _ := runner.Run("git", "rev-parse", "--show-toplevel")
	_, branchOut, _ := runner.Run("git", "branch", "--show-current")
	_, shaOut, _ := runner.Run("git", "rev-parse", "HEAD")
	statusCode, statusOut, _ := runner.Run("git", "status", "--short")

	statusText := strings.TrimSpace(statusOut)
	var statusLines []string
	if statusText != "" {
		for _, line := range strings.Split(statusText, "\n") {
			if strings.TrimSpace(line) != "" {
				statusLines = append(statusLines, line)
			}
		}
	}

	return mission.GitInfoData{
		Available:  true,
		IsRepo:     true,
		Root:       strings.TrimSpace(rootOut),
		Branch:     strings.TrimSpace(branchOut),
		Sha:        strings.TrimSpace(shaOut),
		Dirty:      statusCode == 0 && len(statusLines) > 0,
		DirtyFiles: len(statusLines),
	}
}

// NoopRunner is a CommandRunner that returns empty/false results.
// Useful as a default or for tests that don't need git.
var NoopRunner CommandRunner = OSCommandRunner{}
