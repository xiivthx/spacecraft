package hooks

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestLoadConfig_ValidJSON(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "hooks.json")
	cfg := Config{
		Hooks: []Hook{
			{Event: "mission.created", Label: "create-hook", Command: "echo hi", Blocking: false, Timeout: 30},
			{Event: "mission.state.changed", Label: "state-hook", Command: "make lint", Blocking: true, Timeout: 60},
		},
	}
	data, _ := json.Marshal(cfg)
	os.WriteFile(p, data, 0644)

	got, err := LoadConfig(p)
	if err != nil {
		t.Fatalf("LoadConfig error: %v", err)
	}
	if len(got.Hooks) != 2 {
		t.Fatalf("expected 2 hooks, got %d", len(got.Hooks))
	}
	if got.Hooks[0].Event != "mission.created" || got.Hooks[0].Label != "create-hook" {
		t.Errorf("hook 0: %+v", got.Hooks[0])
	}
	if got.Hooks[1].Event != "mission.state.changed" || !got.Hooks[1].Blocking {
		t.Errorf("hook 1: %+v", got.Hooks[1])
	}
	if got.Hooks[1].Timeout != 60 {
		t.Errorf("hook 1 timeout = %d, want 60", got.Hooks[1].Timeout)
	}
}

func TestLoadConfig_MissingFile(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "nonexistent.json")

	got, err := LoadConfig(p)
	if err != nil {
		t.Fatalf("expected no error for missing file, got: %v", err)
	}
	if len(got.Hooks) != 0 {
		t.Errorf("expected empty hooks, got %d", len(got.Hooks))
	}
}

func TestLoadConfig_InvalidJSON(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "hooks.json")
	os.WriteFile(p, []byte("not json{{{{"), 0644)

	got, err := LoadConfig(p)
	if err != nil {
		t.Fatalf("expected no error on invalid JSON, got: %v", err)
	}
	if len(got.Hooks) != 0 {
		t.Errorf("expected empty hooks on invalid JSON, got %d", len(got.Hooks))
	}
}

func TestLoadConfig_MissingRequiredFields(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "hooks.json")
	os.WriteFile(p, []byte(`{
		"hooks": [
			{"label": "no-event", "command": "echo x"},
			{"event": "mission.created", "label": "no-command"},
			{"event": "mission.evidence.appended", "command": "echo ok", "label": "good"}
		]
	}`), 0644)

	got, err := LoadConfig(p)
	if err != nil {
		t.Fatalf("LoadConfig error: %v", err)
	}
	if len(got.Hooks) != 1 {
		t.Fatalf("expected 1 valid hook, got %d", len(got.Hooks))
	}
	if got.Hooks[0].Event != "mission.evidence.appended" || got.Hooks[0].Label != "good" {
		t.Errorf("hook: %+v", got.Hooks[0])
	}
}

func TestLoadConfig_UnknownEvent(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "hooks.json")
	os.WriteFile(p, []byte(`{
		"hooks": [
			{"event": "unknown.event", "command": "echo x", "label": "bad"},
			{"event": "mission.created", "command": "echo ok", "label": "good"}
		]
	}`), 0644)

	got, err := LoadConfig(p)
	if err != nil {
		t.Fatalf("LoadConfig error: %v", err)
	}
	if len(got.Hooks) != 1 {
		t.Fatalf("expected 1 valid hook, got %d", len(got.Hooks))
	}
	if got.Hooks[0].Event != "mission.created" {
		t.Errorf("hook: %+v", got.Hooks[0])
	}
}

func TestLoadConfig_DefaultTimeout(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "hooks.json")
	os.WriteFile(p, []byte(`{
		"hooks": [
			{"event": "mission.created", "command": "echo x"}
		]
	}`), 0644)

	got, err := LoadConfig(p)
	if err != nil {
		t.Fatalf("LoadConfig error: %v", err)
	}
	if len(got.Hooks) != 1 {
		t.Fatalf("expected 1 hook, got %d", len(got.Hooks))
	}
	if got.Hooks[0].Timeout != 30 {
		t.Errorf("default timeout = %d, want 30", got.Hooks[0].Timeout)
	}
}

func TestLoadConfig_DefaultLabel(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "hooks.json")
	os.WriteFile(p, []byte(`{
		"hooks": [
			{"event": "mission.created", "command": "echo x"}
		]
	}`), 0644)

	got, err := LoadConfig(p)
	if err != nil {
		t.Fatalf("LoadConfig error: %v", err)
	}
	if got.Hooks[0].Label != "mission.created" {
		t.Errorf("default label = %q, want 'mission.created'", got.Hooks[0].Label)
	}
}

func TestMatch_NamedEvent(t *testing.T) {
	cfg := &Config{
		Hooks: []Hook{
			{Event: "mission.created", Label: "a", Command: "echo a"},
			{Event: "mission.state.changed", Label: "b", Command: "echo b"},
		},
	}
	matched := Match(cfg, "mission.created")
	if len(matched) != 1 {
		t.Fatalf("expected 1 match, got %d", len(matched))
	}
	if matched[0].Label != "a" {
		t.Errorf("matched hook = %v", matched[0])
	}
}

func TestMatch_Wildcard(t *testing.T) {
	cfg := &Config{
		Hooks: []Hook{
			{Event: "*", Label: "wild", Command: "echo wild"},
			{Event: "mission.created", Label: "named", Command: "echo named"},
		},
	}
	matched := Match(cfg, "mission.state.changed")
	if len(matched) != 1 {
		t.Fatalf("expected 1 wildcard match, got %d", len(matched))
	}
	if matched[0].Label != "wild" {
		t.Errorf("wildcard hook = %v", matched[0])
	}
}

func TestMatch_WildcardPlusNamed(t *testing.T) {
	cfg := &Config{
		Hooks: []Hook{
			{Event: "*", Label: "wild", Command: "echo wild"},
			{Event: "mission.created", Label: "named", Command: "echo named"},
		},
	}
	matched := Match(cfg, "mission.created")
	if len(matched) != 2 {
		t.Fatalf("expected 2 matches (wildcard + named), got %d", len(matched))
	}
}

func TestMatch_NoMatch(t *testing.T) {
	cfg := &Config{
		Hooks: []Hook{
			{Event: "mission.created", Label: "a", Command: "echo a"},
		},
	}
	matched := Match(cfg, "mission.archived")
	if len(matched) != 0 {
		t.Errorf("expected 0 matches, got %d", len(matched))
	}
}

func TestMatch_NilConfig(t *testing.T) {
	matched := Match(nil, "mission.created")
	if matched != nil {
		t.Error("expected nil for nil config")
	}
}

func TestRun_Success(t *testing.T) {
	ctx := context.Background()
	hook := Hook{
		Event:   "mission.created",
		Label:   "test-ok",
		Command: "echo hello world",
		Timeout: 10,
	}

	err := Run(ctx, hook)
	if err != nil {
		t.Errorf("Run returned error: %v", err)
	}
}

func TestRun_ExitCode1(t *testing.T) {
	ctx := context.Background()
	hook := Hook{
		Event:   "mission.created",
		Label:   "test-fail",
		Command: "exit 1",
		Timeout: 10,
	}

	err := Run(ctx, hook)
	if err == nil {
		t.Error("expected error for exit code 1")
	}
	if !strings.Contains(err.Error(), "exited with code 1") {
		t.Errorf("error message: %v", err)
	}
}

func TestRun_Timeout(t *testing.T) {
	ctx := context.Background()
	hook := Hook{
		Event:   "mission.created",
		Label:   "test-slow",
		Command: "sleep 10",
		Timeout: 1,
	}

	start := time.Now()
	err := Run(ctx, hook)
	elapsed := time.Since(start)

	if err == nil {
		t.Error("expected timeout error")
	}
	if !strings.Contains(err.Error(), "timed out") {
		t.Errorf("error message: %v", err)
	}
	if elapsed > 5*time.Second {
		t.Errorf("timeout took too long: %v", elapsed)
	}
}

func TestRun_ShellPipeline(t *testing.T) {
	ctx := context.Background()
	hook := Hook{
		Event:   "mission.created",
		Label:   "pipe",
		Command: "echo one && echo two | tr 'a-z' 'A-Z'",
		Timeout: 10,
	}

	err := Run(ctx, hook)
	if err != nil {
		t.Errorf("pipeline Run error: %v", err)
	}
}

func TestRun_ContextCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	hook := Hook{
		Event:   "mission.created",
		Label:   "test-cancel",
		Command: "sleep 10",
		Timeout: 30,
	}

	go func() {
		time.Sleep(200 * time.Millisecond)
		cancel()
	}()

	start := time.Now()
	err := Run(ctx, hook)
	elapsed := time.Since(start)

	if err != context.Canceled {
		t.Errorf("expected context.Canceled, got: %v", err)
	}
	if elapsed > 5*time.Second {
		t.Errorf("cancellation took too long: %v", elapsed)
	}
}

func TestFire_NoHooks(t *testing.T) {
	cfg := &Config{}
	err := Fire(context.Background(), cfg, "mission.created")
	if err != nil {
		t.Errorf("Fire with empty config returned error: %v", err)
	}
}

func TestFire_NilConfig(t *testing.T) {
	err := Fire(context.Background(), nil, "mission.created")
	if err != nil {
		t.Errorf("Fire with nil config returned error: %v", err)
	}
}

func TestFire_NonBlockingFailure(t *testing.T) {
	cfg := &Config{
		Hooks: []Hook{
			{Event: "mission.created", Label: "fail-nonblock", Command: "exit 1", Blocking: false, Timeout: 10},
			{Event: "mission.created", Label: "succeed-after", Command: "echo ok", Blocking: false, Timeout: 10},
		},
	}
	err := Fire(context.Background(), cfg, "mission.created")
	if err != nil {
		t.Errorf("non-blocking failure should not return error: %v", err)
	}
}

func TestFire_BlockingFailure(t *testing.T) {
	cfg := &Config{
		Hooks: []Hook{
			{Event: "mission.created", Label: "fail-block", Command: "exit 1", Blocking: true, Timeout: 10},
			{Event: "mission.created", Label: "never-runs", Command: "echo nope", Blocking: false, Timeout: 10},
		},
	}
	err := Fire(context.Background(), cfg, "mission.created")
	if err == nil {
		t.Error("expected error for blocking failure")
	}
	if !strings.Contains(err.Error(), "blocking") {
		t.Errorf("error should mention 'blocking': %v", err)
	}
	if !strings.Contains(err.Error(), "fail-block") {
		t.Errorf("error should mention hook label: %v", err)
	}
}

func TestFire_ConfigOrder(t *testing.T) {
	cfg := &Config{
		Hooks: []Hook{
			{Event: "mission.created", Label: "first", Command: "echo 1", Timeout: 10},
			{Event: "mission.created", Label: "second", Command: "echo 2", Timeout: 10},
			{Event: "mission.created", Label: "third", Command: "echo 3", Timeout: 10},
		},
	}
	err := Fire(context.Background(), cfg, "mission.created")
	if err != nil {
		t.Errorf("Fire error: %v", err)
	}
}

func TestFire_OnlyMatchingEvents(t *testing.T) {
	cfg := &Config{
		Hooks: []Hook{
			{Event: "mission.created", Label: "should-run", Command: "echo ok", Timeout: 10},
			{Event: "mission.archived", Label: "should-not-run", Command: "echo nope", Timeout: 10},
		},
	}
	err := Fire(context.Background(), cfg, "mission.created")
	if err != nil {
		t.Errorf("Fire error: %v", err)
	}
}

func TestRun_OutputPrefix(t *testing.T) {
	// Capture stdout to verify prefix format.
	origStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	ctx := context.Background()
	hook := Hook{
		Event:   "mission.created",
		Label:   "output-test",
		Command: "echo line1 && echo line2",
		Timeout: 10,
	}

	err := Run(ctx, hook)
	w.Close()
	os.Stdout = origStdout

	if err != nil {
		t.Fatalf("Run error: %v", err)
	}

	var output strings.Builder
	buf := make([]byte, 4096)
	for {
		n, _ := r.Read(buf)
		if n == 0 {
			break
		}
		output.Write(buf[:n])
	}

	outStr := output.String()
	fmt.Fprintf(origStdout, "CAPTURED OUTPUT:\n%s\n", outStr)

	if !strings.Contains(outStr, "[hook] output-test") {
		t.Error("missing [hook] lifecycle prefix")
	}
	if !strings.Contains(outStr, "[output-test] line1") {
		t.Error("missing child stdout prefix [output-test] line1")
	}
	if !strings.Contains(outStr, "[output-test] line2") {
		t.Error("missing child stdout prefix [output-test] line2")
	}
	if !strings.Contains(outStr, "[hook] output-test ✓") {
		t.Error("missing success ✓ marker")
	}
}

func TestRun_StderrOutput(t *testing.T) {
	origStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	ctx := context.Background()
	hook := Hook{
		Event:   "mission.created",
		Label:   "stderr-test",
		Command: "echo error >&2",
		Timeout: 10,
	}

	err := Run(ctx, hook)
	w.Close()
	os.Stderr = origStderr

	if err != nil {
		t.Fatalf("Run error: %v", err)
	}

	var output strings.Builder
	buf := make([]byte, 4096)
	for {
		n, _ := r.Read(buf)
		if n == 0 {
			break
		}
		output.Write(buf[:n])
	}

	outStr := output.String()
	if !strings.Contains(outStr, "[stderr-test] error") {
		t.Errorf("missing stderr prefix: %s", outStr)
	}
}

func TestRun_FailureOutput(t *testing.T) {
	origStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	ctx := context.Background()
	hook := Hook{
		Event:   "mission.created",
		Label:   "fail-output",
		Command: "exit 42",
		Timeout: 10,
	}

	_ = Run(ctx, hook)
	w.Close()
	os.Stdout = origStdout

	var output strings.Builder
	buf := make([]byte, 4096)
	for {
		n, _ := r.Read(buf)
		if n == 0 {
			break
		}
		output.Write(buf[:n])
	}

	outStr := output.String()
	if !strings.Contains(outStr, "[hook] fail-output ✗ (exit: 42)") {
		t.Errorf("missing failure marker with exit code: %s", outStr)
	}
}

func TestRun_NoTrailingNewline(t *testing.T) {
	origStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	ctx := context.Background()
	hook := Hook{
		Event:   "mission.created",
		Label:   "no-nl",
		Command: "printf noline",
		Timeout: 10,
	}

	_ = Run(ctx, hook)
	w.Close()
	os.Stdout = origStdout

	var output strings.Builder
	buf := make([]byte, 4096)
	for {
		n, _ := r.Read(buf)
		if n == 0 {
			break
		}
		output.Write(buf[:n])
	}

	outStr := output.String()
	if !strings.Contains(outStr, "[no-nl] noline") {
		t.Errorf("output should contain prefixed line even without trailing newline: %s", outStr)
	}
}

func TestDeployHooks(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "hooks.json")
	cfg := Config{
		Hooks: []Hook{
			{Event: "deploy.before", Label: "pre-deploy", Command: "echo deploying", Blocking: true, Timeout: 30},
			{Event: "deploy.after", Label: "post-deploy", Command: "echo deployed", Blocking: false, Timeout: 30},
		},
	}
	data, _ := json.Marshal(cfg)
	os.WriteFile(p, data, 0644)

	got, err := LoadConfig(p)
	if err != nil {
		t.Fatalf("LoadConfig error: %v", err)
	}
	if len(got.Hooks) != 2 {
		t.Fatalf("expected 2 hooks, got %d", len(got.Hooks))
	}

	beforeHooks := Match(got, "deploy.before")
	if len(beforeHooks) != 1 {
		t.Errorf("expected 1 deploy.before hook, got %d", len(beforeHooks))
	}
	if beforeHooks[0].Label != "pre-deploy" {
		t.Errorf("expected pre-deploy, got %s", beforeHooks[0].Label)
	}

	afterHooks := Match(got, "deploy.after")
	if len(afterHooks) != 1 {
		t.Errorf("expected 1 deploy.after hook, got %d", len(afterHooks))
	}
	if afterHooks[0].Label != "post-deploy" {
		t.Errorf("expected post-deploy, got %s", afterHooks[0].Label)
	}
}

func TestDeployHookExecution(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "hooks.json")
	cfg := Config{
		Hooks: []Hook{
			{Event: "deploy.before", Label: "pre-deploy", Command: "echo before-deploy", Blocking: true, Timeout: 30},
		},
	}
	data, _ := json.Marshal(cfg)
	os.WriteFile(p, data, 0644)

	loaded, err := LoadConfig(p)
	if err != nil {
		t.Fatalf("LoadConfig error: %v", err)
	}

	ctx := context.Background()
	err = Fire(ctx, loaded, "deploy.before")
	if err != nil {
		t.Errorf("Fire deploy.before failed: %v", err)
	}
}
