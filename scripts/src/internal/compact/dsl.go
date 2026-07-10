package compact

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strings"
)

// DSLConfig is the top-level structure for the filter config DSL.
type DSLConfig struct {
	Rules []Rule `json:"rules"`
}

// Rule maps a command pattern to a pipeline of filter stages.
type Rule struct {
	Exe    string  `json:"exe"`
	Arg1   string  `json:"arg1,omitempty"`
	Stages []Stage `json:"stages"`
}

// Stage is a single filter stage in a pipeline.
type Stage struct {
	Include     *IncludeStage     `json:"include,omitempty"`
	Exclude     *ExcludeStage     `json:"exclude,omitempty"`
	Dedup       *DedupStage       `json:"dedup,omitempty"`
	Truncate    *TruncateStage    `json:"truncate,omitempty"`
	StripPrefix *StripPrefixStage `json:"stripPrefix,omitempty"`
	Passthrough *PassthroughStage `json:"passthrough,omitempty"`
}

// IncludeStage keeps only lines matching the pattern.
type IncludeStage struct {
	Pattern string `json:"pattern"`
}

func (s *IncludeStage) Apply(stdout string) string {
	if stdout == "" || s.Pattern == "" {
		return stdout
	}
	re, err := regexp.Compile(s.Pattern)
	if err != nil {
		return stdout
	}
	lines := strings.Split(stdout, "\n")
	var kept []string
	for _, line := range lines {
		if re.MatchString(line) {
			kept = append(kept, line)
		}
	}
	return strings.Join(kept, "\n")
}

// ExcludeStage drops lines matching the pattern.
type ExcludeStage struct {
	Pattern string `json:"pattern"`
}

func (s *ExcludeStage) Apply(stdout string) string {
	if stdout == "" || s.Pattern == "" {
		return stdout
	}
	re, err := regexp.Compile(s.Pattern)
	if err != nil {
		return stdout
	}
	lines := strings.Split(stdout, "\n")
	var kept []string
	for _, line := range lines {
		if !re.MatchString(line) {
			kept = append(kept, line)
		}
	}
	return strings.Join(kept, "\n")
}

// DedupStage collapses consecutive identical lines.
type DedupStage struct{}

func (s *DedupStage) Apply(stdout string) string {
	if stdout == "" {
		return ""
	}
	return dedupLines(stdout)
}

// TruncateStage keeps head+tail lines with a summary separator.
type TruncateStage struct {
	Head int `json:"head"`
	Tail int `json:"tail"`
}

func (s *TruncateStage) Apply(stdout string) string {
	if stdout == "" {
		return ""
	}
	lines := strings.Split(stdout, "\n")
	head := s.Head
	tail := s.Tail
	if head <= 0 {
		head = 250
	}
	if tail <= 0 {
		tail = 250
	}
	total := len(lines)
	if total <= head+tail {
		return stdout
	}

	var result strings.Builder
	result.Grow((head + tail + 3) * 80)

	for i := 0; i < head && i < total; i++ {
		result.WriteString(lines[i])
		result.WriteByte('\n')
	}

	skipped := total - head - tail
	fmt.Fprintf(&result, "--- %d lines skipped (total: %d) ---\n", skipped, total)

	for i := total - tail; i < total; i++ {
		result.WriteString(lines[i])
		result.WriteByte('\n')
	}

	return strings.TrimRight(result.String(), "\n")
}

// StripPrefixStage removes a prefix from each line.
type StripPrefixStage struct {
	Prefix string `json:"prefix"`
}

func (s *StripPrefixStage) Apply(stdout string) string {
	if stdout == "" || s.Prefix == "" {
		return stdout
	}
	lines := strings.Split(stdout, "\n")
	for i, line := range lines {
		if strings.HasPrefix(line, s.Prefix) {
			lines[i] = line[len(s.Prefix):]
		}
	}
	return strings.Join(lines, "\n")
}

// PassthroughStage is an explicit no-op stage.
type PassthroughStage struct{}

func (s *PassthroughStage) Apply(stdout string) string { return stdout }

// LoadDSLFilter loads a DSL filter from path for the given command.
// If the config file does not exist, it returns nil, nil.
func LoadDSLFilter(path string, ci *CommandInfo) (Filter, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config DSLConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	rule := matchRule(&config, ci)
	if rule == nil {
		return nil, nil
	}

	return buildPipeline(rule)
}

// matchRule finds the first rule whose exe matches ci.Exe and,
// if rule.Arg1 is non-empty, rule.Arg1 matches ci.Arg1.
// Rules are checked in order; first match wins.
func matchRule(config *DSLConfig, ci *CommandInfo) *Rule {
	if ci == nil {
		return nil
	}
	for i := range config.Rules {
		r := &config.Rules[i]
		if r.Exe != ci.Exe {
			continue
		}
		if r.Arg1 != "" && r.Arg1 != ci.Arg1 {
			continue
		}
		return r
	}
	return nil
}

// buildPipeline creates a composite Filter from a Rule's stages.
func buildPipeline(rule *Rule) (Filter, error) {
	stages := make([]Filter, 0, len(rule.Stages))
	for _, s := range rule.Stages {
		f, err := stageToFilter(s)
		if err != nil {
			return nil, err
		}
		if f == nil {
			continue // skip empty stages
		}
		stages = append(stages, f)
	}
	if len(stages) == 0 {
		return nil, nil
	}
	return &pipelineFilter{stages: stages}, nil
}

// pipelineFilter chains multiple Filters, feeding output of one to input of next.
type pipelineFilter struct {
	stages []Filter
}

func (p *pipelineFilter) Apply(stdout string) string {
	for _, s := range p.stages {
		stdout = s.Apply(stdout)
	}
	return stdout
}

// stageToFilter converts a Stage to a Filter based on which field is set.
// Returns an error for empty stages or invalid regex patterns.
func stageToFilter(s Stage) (Filter, error) {
	switch {
	case s.Include != nil:
		if _, err := regexp.Compile(s.Include.Pattern); err != nil {
			return nil, fmt.Errorf("invalid include pattern %q: %w", s.Include.Pattern, err)
		}
		return s.Include, nil
	case s.Exclude != nil:
		if _, err := regexp.Compile(s.Exclude.Pattern); err != nil {
			return nil, fmt.Errorf("invalid exclude pattern %q: %w", s.Exclude.Pattern, err)
		}
		return s.Exclude, nil
	case s.Dedup != nil:
		return s.Dedup, nil
	case s.Truncate != nil:
		return s.Truncate, nil
	case s.StripPrefix != nil:
		return s.StripPrefix, nil
	case s.Passthrough != nil:
		return s.Passthrough, nil
	default:
		return nil, fmt.Errorf("stage has no action defined")
	}
}
