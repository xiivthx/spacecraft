package eval

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"spacecraft/internal/mission"
)

// RunLMJudge evaluates non-deterministic output quality using a secondary model.
// Falls back gracefully to deterministic-only when no judge is configured or available.
func RunLMJudge(entries []mission.EvidenceEntry) LMJudgeResult {
	judgeModel := os.Getenv("SPACECRAFT_JUDGE_MODEL")
	if judgeModel == "" {
		return LMJudgeResult{
			Score:          0,
			MaxScore:       4,
			Fallback:       true,
			FallbackReason: "SPACECRAFT_JUDGE_MODEL not configured — using deterministic-only evaluation",
		}
	}

	if !isLMJudgeAvailable() {
		return LMJudgeResult{
			Score:          0,
			MaxScore:       4,
			Fallback:       true,
			FallbackReason: "no LM judge CLI available (checked: agy, llm, opencode) — using deterministic-only evaluation",
		}
	}

	prompt := buildJudgePrompt(entries)
	reasoning, score, err := invokeJudge(judgeModel, prompt)
	if err != nil {
		return LMJudgeResult{
			Score:          0,
			MaxScore:       4,
			Fallback:       true,
			FallbackReason: fmt.Sprintf("LM judge invocation failed: %v — using deterministic-only evaluation", err),
		}
	}

	return LMJudgeResult{
		Score:     score,
		MaxScore:  4,
		Reasoning: reasoning,
		Model:     judgeModel,
		Fallback:  false,
	}
}

func isLMJudgeAvailable() bool {
	if _, err := exec.LookPath("agy"); err == nil {
		return true
	}
	if _, err := exec.LookPath("llm"); err == nil {
		return true
	}
	_, err := exec.LookPath("opencode")
	return err == nil
}

func buildJudgePrompt(entries []mission.EvidenceEntry) string {
	var sb strings.Builder
	sb.WriteString("You are an evaluation judge for an AI coding agent's output quality.\n\n")
	sb.WriteString("Review the following evidence entries and score the agent's response quality on a scale of 0-4:\n")
	sb.WriteString("0 = unusable, incorrect, or incomprehensible\n")
	sb.WriteString("1 = major issues, barely functional\n")
	sb.WriteString("2 = adequate, some issues but acceptable\n")
	sb.WriteString("3 = good quality, minor issues only\n")
	sb.WriteString("4 = excellent, clear, correct, complete\n\n")

	if len(entries) > 0 {
		sb.WriteString("Evidence entries:\n")
		for _, e := range entries {
			sb.WriteString(fmt.Sprintf("- [%s] %s: %s (exit=%d, stdout=%d bytes, stderr=%d bytes)\n",
				e.Label, e.Command, e.ID, e.ExitCode, len(e.Stdout), len(e.Stderr)))
		}
	} else {
		sb.WriteString("(no evidence entries to evaluate)\n")
	}

	sb.WriteString("\nRespond with a JSON object: {\"score\": <int 0-4>, \"reasoning\": \"<brief explanation>\"}\n")
	return sb.String()
}

func invokeJudge(model, prompt string) (reasoning string, score int, err error) {
	var cmd *exec.Cmd

	if _, e := exec.LookPath("agy"); e == nil {
		cmd = exec.Command("agy", "--model", model, "--prompt", prompt)
	} else if _, e := exec.LookPath("llm"); e == nil {
		cmd = exec.Command("llm", "-m", model, prompt)
	} else if _, e := exec.LookPath("opencode"); e == nil {
		cmd = exec.Command("opencode", "eval", "--model", model, prompt)
	} else {
		return "", 0, fmt.Errorf("no LM judge CLI available (checked: agy, llm, opencode)")
	}

	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", 0, fmt.Errorf("judge invocation failed: %w (output: %s)", err, string(out))
	}

	reasoning = strings.TrimSpace(string(out))
	score = extractScore(reasoning)

	return reasoning, score, nil
}

func extractScore(output string) int {
	score := 0
	output = strings.ToLower(output)

	if strings.Contains(output, `"score": 4`) || strings.Contains(output, `"score":4`) {
		score = 4
	} else if strings.Contains(output, `"score": 3`) || strings.Contains(output, `"score":3`) {
		score = 3
	} else if strings.Contains(output, `"score": 2`) || strings.Contains(output, `"score":2`) {
		score = 2
	} else if strings.Contains(output, `"score": 1`) || strings.Contains(output, `"score":1`) {
		score = 1
	}

	if strings.Contains(output, "excellent") || strings.Contains(output, "perfect") {
		score = maxInt(score, 4)
	} else if strings.Contains(output, "good") || strings.Contains(output, "well") {
		score = maxInt(score, 3)
	} else if strings.Contains(output, "adequate") || strings.Contains(output, "acceptable") {
		score = maxInt(score, 2)
	} else if strings.Contains(output, "poor") || strings.Contains(output, "bad") || strings.Contains(output, "failed") {
		score = minInt(score, 1)
	}

	return score
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}
