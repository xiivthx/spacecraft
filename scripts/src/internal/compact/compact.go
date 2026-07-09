// Package compact filters and compresses command output for LLM context.
// It runs an arbitrary command, passes stdout through a configurable filter,
// and returns compact, token-optimized output while preserving the exit code
// and forwarding stderr unfiltered.
package compact

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
)

// Filter transforms raw command stdout into compact output.
// Implementations should preserve error/fatal lines and never alter semantics.
type Filter interface {
	Apply(stdout string) string
}

// CommandInfo stores just enough to auto-detect the right filter.
type CommandInfo struct {
	Exe  string // executable name (e.g., "git", "go", "ls")
	Arg1 string // first positional argument (e.g., "status", "test", "-la")
}

// ParseCommand extracts CommandInfo from a raw argument list.
// Returns nil if the slice is empty.
func ParseCommand(args []string) *CommandInfo {
	if len(args) == 0 {
		return nil
	}
	ci := &CommandInfo{Exe: args[0]}
	if len(args) > 1 && len(args[1]) > 0 && args[1][0] != '-' {
		ci.Arg1 = args[1]
	}
	return ci
}

// Result holds the outcome of a Runner execution.
type Result struct {
	ExitCode int
	Stdout   string // raw stdout (before filtering)
	Stderr   string // raw stderr (never filtered)
	Output   string // filtered stdout + unfiltered stderr (what to display)
}

// Runner executes a command and applies a Filter to its stdout.
type Runner struct {
	command []string
	filter  Filter
}

// NewRunner creates a Runner for the given command and filter.
// Use nil filter to pass stdout through unfiltered.
func NewRunner(command []string, filter Filter) *Runner {
	return &Runner{command: command, filter: filter}
}

// Run executes the command and returns the Result.
// Stderr is never filtered — it is forwarded as-is.
// If the filter is nil, stdout is passed through unfiltered.
func (r *Runner) Run() (*Result, error) {
	if len(r.command) == 0 {
		return nil, fmt.Errorf("compact: empty command")
	}

	cmd := exec.Command(r.command[0], r.command[1:]...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	cmd.Stdin = os.Stdin

	err := cmd.Run()

	exitCode := 0
	if err != nil {
		if ee, ok := err.(*exec.ExitError); ok {
			exitCode = ee.ExitCode()
		} else {
			exitCode = 127
			// Append execution error to stderr so it's visible.
			stderr.WriteString(err.Error() + "\n")
		}
	}

	rawStdout := stdout.String()
	rawStderr := stderr.String()

	filtered := rawStdout
	if r.filter != nil {
		filtered = r.filter.Apply(rawStdout)
	}

	output := filtered
	if rawStderr != "" {
		if output != "" && !bytes.HasSuffix([]byte(output), []byte("\n")) {
			output += "\n"
		}
		output += rawStderr
	}

	return &Result{
		ExitCode: exitCode,
		Stdout:   rawStdout,
		Stderr:   rawStderr,
		Output:   output,
	}, nil
}

// noFilter lets stdout pass through unchanged.
type noFilter struct{}

func (noFilter) Apply(s string) string { return s }
