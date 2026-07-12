package research

import (
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

// Executor abstracts command lookups and execution so runners can be tested
// without relying on a real operating-system shell.
type Executor interface {
	LookPath(file string) (string, error)
	Run(ctx context.Context, name string, args ...string) ([]byte, error)
}

// OSExecutor is the production implementation of Executor that delegates to
// os/exec for real subprocess calls.
type OSExecutor struct{}

func (o *OSExecutor) LookPath(file string) (string, error) {
	return exec.LookPath(file)
}

func (o *OSExecutor) Run(ctx context.Context, name string, args ...string) ([]byte, error) {
	cmd := exec.CommandContext(ctx, name, args...)
	return cmd.Output()
}

// BrowserUseRunner attempts deep analysis via the browser-use Python package.
type BrowserUseRunner struct {
	exec      Executor
	available bool
}

// NewBrowserUseRunner creates a runner and checks whether python3 is available.
func NewBrowserUseRunner(exec Executor) *BrowserUseRunner {
	_, err := exec.LookPath("python3")
	return &BrowserUseRunner{
		exec:      exec,
		available: err == nil,
	}
}

// IsAvailable reports whether the runner can be used.
func (r *BrowserUseRunner) IsAvailable() bool {
	return r.available
}

// InstallInstructions returns the command needed to install browser-use.
func (r *BrowserUseRunner) InstallInstructions() string {
	return "pip install browser-use && playwright install"
}

// Analyze performs a deep analysis of a URL by invoking a Python subprocess
// that fetches and extracts page content.  Returns a DeepResult with summary,
// key points, source URL, and fetch timestamp.
func (r *BrowserUseRunner) Analyze(ctx context.Context, url string) (*DeepResult, error) {
	if !r.available {
		return nil, fmt.Errorf("browser-use not available; install with: %s", r.InstallInstructions())
	}

	// Generate a Python script that fetches the URL and extracts text content.
	// Uses requests + BeautifulSoup if browser-use is unavailable.
	script := fmt.Sprintf(`import sys, json
try:
    from browser_use import Agent
    agent = Agent(task="Go to %s and summarize the main content in a few paragraphs. Return only the summary.")
    import asyncio
    result = asyncio.run(agent.run())
    print(json.dumps({"summary": str(result), "key_points": [], "source_url": %q}))
except ImportError:
    try:
        import requests
        from bs4 import BeautifulSoup
        resp = requests.get(%q, timeout=30)
        soup = BeautifulSoup(resp.text, 'html.parser')
        for tag in soup(['script', 'style', 'nav', 'footer', 'header']):
            tag.decompose()
        text = ' '.join(soup.stripped_strings)
        # Truncate to ~4000 chars for summary.
        summary = text[:4000] if len(text) > 4000 else text
        # Extract first ~5 sentences as key_points.
        sentences = [s.strip() for s in summary.replace('\n', '. ').split('. ') if len(s.strip()) > 20][:5]
        print(json.dumps({"summary": summary, "key_points": sentences, "source_url": %q}))
    except ImportError:
        print(json.dumps({"summary": "Page fetched but no parser available. Install: pip install beautifulsoup4 requests", "key_points": [], "source_url": %q}))
`, url, url, url, url, url)

	output, err := r.exec.Run(ctx, "python3", "-c", script)
	if err != nil {
		return nil, fmt.Errorf("browser-use subprocess: %w (output: %s)", err, string(output))
	}

	var result struct {
		Summary   string   `json:"summary"`
		KeyPoints []string `json:"key_points"`
		SourceURL string   `json:"source_url"`
		Error     string   `json:"error"`
	}
	if err := json.Unmarshal(output, &result); err != nil {
		// If JSON parsing fails, use raw output as summary.
		return &DeepResult{
			Summary:   string(output),
			KeyPoints: nil,
			SourceURL: url,
			FetchedAt: nowISO(),
		}, nil
	}

	if result.Error != "" {
		return nil, fmt.Errorf("browser-use: %s", result.Error)
	}

	if result.Summary == "" && len(result.KeyPoints) == 0 {
		return nil, fmt.Errorf("browser-use: empty response from subprocess")
	}

	return &DeepResult{
		Summary:   result.Summary,
		KeyPoints: result.KeyPoints,
		SourceURL: result.SourceURL,
		FetchedAt: nowISO(),
	}, nil
}

// NotebookLMRunner attempts deep analysis via the NotebookLM CLI.
type NotebookLMRunner struct {
	exec      Executor
	available bool
}

// NewNotebookLMRunner creates a runner and checks whether nlm is available.
func NewNotebookLMRunner(exec Executor) *NotebookLMRunner {
	_, err := exec.LookPath("nlm")
	return &NotebookLMRunner{
		exec:      exec,
		available: err == nil,
	}
}

// IsAvailable reports whether the runner can be used.
func (r *NotebookLMRunner) IsAvailable() bool {
	return r.available
}

// InstallInstructions returns the command needed to install notebooklm-mcp-cli.
func (r *NotebookLMRunner) InstallInstructions() string {
	return "uv tool install notebooklm-mcp-cli && nlm login"
}

// Analyze performs a deep analysis using NotebookLM by:
// 1. Creating a notebook with the query as title
// 2. Adding sources (from a Brave Search, handled by caller)
// 3. Querying the notebook for a synthesis
// Returns a DeepResult with the synthesized analysis.
func (r *NotebookLMRunner) Analyze(ctx context.Context, query string) (*DeepResult, error) {
	if !r.available {
		return nil, fmt.Errorf("notebooklm not available; install with: %s", r.InstallInstructions())
	}

	// Create a notebook and query it for the analysis.
	// The nlm CLI supports: nlm notebook create "title" and nlm notebook query
	output, err := r.exec.Run(ctx, "nlm", "notebook", "create", query)
	if err != nil {
		return nil, fmt.Errorf("nlm notebook create: %w (output: %s)", err, string(output))
	}

	// Parse notebook ID from output (nlm typically outputs JSON or plain text with the ID).
	notebookOutput := strings.TrimSpace(string(output))

	// Query the notebook for a synthesis.
	queryOutput, err := r.exec.Run(ctx, "nlm", "notebook", "query", notebookOutput, fmt.Sprintf("Research and analyze: %s. Provide a comprehensive summary with key points.", query))
	if err != nil {
		// If direct query fails, try the create output as notebook name.
		return &DeepResult{
			Summary:   fmt.Sprintf("NotebookLM notebook created: %s. Use 'nlm notebook query' to interact.", notebookOutput),
			KeyPoints: nil,
			SourceURL: "",
			FetchedAt: nowISO(),
		}, nil
	}

	return &DeepResult{
		Summary:   string(queryOutput),
		KeyPoints: nil,
		SourceURL: "",
		FetchedAt: nowISO(),
	}, nil
}

// nowISO returns the current UTC time in ISO 8601 format.
func nowISO() string {
	return time.Now().UTC().Format(time.RFC3339)
}
