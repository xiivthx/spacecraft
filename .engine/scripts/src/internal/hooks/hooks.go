// Package hooks provides a lifecycle hook system for running user-defined shell
// commands in response to Spacecraft mission events.
package hooks

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"sync"
	"time"
)

var knownEvents = map[string]bool{
	"mission.created":           true,
	"mission.state.changed":     true,
	"mission.evidence.appended": true,
	"mission.plan.saved":        true,
	"mission.review.completed":  true,
	"mission.validated":         true,
	"mission.shipped":           true,
	"mission.archived":          true,
	"*":                         true,
}

type Hook struct {
	Event    string `json:"event"`
	Label    string `json:"label"`
	Command  string `json:"command"`
	Blocking bool   `json:"blocking"`
	Timeout  int    `json:"timeout"`
}

type Config struct {
	Hooks []Hook `json:"hooks"`
}

func LoadConfig(configPath string) (*Config, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return &Config{}, nil
		}
		return nil, err
	}

	var raw struct {
		Hooks []map[string]interface{} `json:"hooks"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		fmt.Fprintf(os.Stderr, "hooks: invalid JSON in .space/hooks.json: %v\n", err)
		return &Config{}, nil
	}

	var hooks []Hook
	for i, h := range raw.Hooks {
		event, _ := h["event"].(string)
		command, _ := h["command"].(string)

		if event == "" {
			fmt.Fprintf(os.Stderr, "hooks: hook #%d missing required field 'event'\n", i+1)
			continue
		}
		if command == "" {
			fmt.Fprintf(os.Stderr, "hooks: hook #%d missing required field 'command'\n", i+1)
			continue
		}

		label := event
		if l, ok := h["label"].(string); ok && l != "" {
			label = l
		}

		blocking := false
		if b, ok := h["blocking"].(bool); ok {
			blocking = b
		}

		timeout := 30
		if t, ok := h["timeout"].(float64); ok && t > 0 {
			timeout = int(t)
		}

		if event != "*" && !knownEvents[event] {
			fmt.Fprintf(os.Stderr, "hooks: hook \"%s\" has unknown event \"%s\"\n", label, event)
			continue
		}

		hooks = append(hooks, Hook{
			Event:    event,
			Label:    label,
			Command:  command,
			Blocking: blocking,
			Timeout:  timeout,
		})
	}

	return &Config{Hooks: hooks}, nil
}

func Match(cfg *Config, event string) []Hook {
	if cfg == nil {
		return nil
	}
	var matched []Hook
	for _, h := range cfg.Hooks {
		if h.Event == event || h.Event == "*" {
			matched = append(matched, h)
		}
	}
	return matched
}

func Run(ctx context.Context, hook Hook) error {
	timeout := hook.Timeout
	if timeout <= 0 {
		timeout = 30
	}

	runCtx, cancel := context.WithTimeout(ctx, time.Duration(timeout)*time.Second)
	defer cancel()

	fmt.Printf("[hook] %s\n", hook.Label)

	cmd := exec.CommandContext(runCtx, "sh", "-c", hook.Command)
	cmd.Env = os.Environ()

	prefix := fmt.Sprintf("[%s]", hook.Label)

	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("hook %q: stdout pipe: %w", hook.Label, err)
	}
	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("hook %q: stderr pipe: %w", hook.Label, err)
	}

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		scanLines(os.Stdout, stdoutPipe, prefix)
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		scanLines(os.Stderr, stderrPipe, prefix)
	}()

	err = cmd.Start()
	if err != nil {
		fmt.Printf("[hook] %s ✗ (exec: %v)\n", hook.Label, err)
		wg.Wait()
		return fmt.Errorf("hook %q: start: %w", hook.Label, err)
	}

	waitErr := cmd.Wait()
	wg.Wait()

	if waitErr != nil {
		if runCtx.Err() == context.DeadlineExceeded {
			fmt.Printf("[hook] %s ⏱ timed out\n", hook.Label)
			return fmt.Errorf("hook %q timed out after %ds", hook.Label, timeout)
		}
		if ctx.Err() == context.Canceled {
			fmt.Printf("[hook] %s ✗ (canceled)\n", hook.Label)
			return context.Canceled
		}
		if exitErr, ok := waitErr.(*exec.ExitError); ok {
			code := exitErr.ExitCode()
			fmt.Printf("[hook] %s ✗ (exit: %d)\n", hook.Label, code)
			return fmt.Errorf("hook %q exited with code %d", hook.Label, code)
		}
		fmt.Printf("[hook] %s ✗ (error: %v)\n", hook.Label, waitErr)
		return fmt.Errorf("hook %q: %w", hook.Label, waitErr)
	}

	fmt.Printf("[hook] %s ✓\n", hook.Label)
	return nil
}

func Fire(ctx context.Context, cfg *Config, event string) error {
	if cfg == nil || len(cfg.Hooks) == 0 {
		return nil
	}
	hooks := Match(cfg, event)
	for _, h := range hooks {
		if err := Run(ctx, h); err != nil {
			if h.Blocking {
				return fmt.Errorf("blocking hook %q failed: %w", h.Label, err)
			}
		}
	}
	return nil
}

func scanLines(w io.Writer, r io.Reader, prefix string) {
	scanner := bufio.NewScanner(r)
	scanner.Buffer(make([]byte, 1024*1024), 1024*1024)
	for scanner.Scan() {
		fmt.Fprintf(w, "%s %s\n", prefix, scanner.Text())
	}
}
