package trace

import "encoding/json"

type EventType string

const (
	EventToolCall       EventType = "tool_call"
	EventModelInvoke    EventType = "model_invoke"
	EventSubagentSpawn  EventType = "subagent_spawn"
	EventSubagentResult EventType = "subagent_result"
	EventCheckpoint     EventType = "checkpoint"
)

type TraceEntry struct {
	ID           string          `json:"id"`
	Seq          int             `json:"seq"`
	TS           string          `json:"ts"`
	Type         EventType       `json:"type"`
	Tool         *string         `json:"tool,omitempty"`
	Args         json.RawMessage `json:"args,omitempty"`
	LatencyMs    int             `json:"latencyMs"`
	InputTokens  int             `json:"inputTokens"`
	OutputTokens int             `json:"outputTokens"`
	Model        *string         `json:"model,omitempty"`
	Subagent     *string         `json:"subagent,omitempty"`
	StepLabel    *string         `json:"stepLabel,omitempty"`
	ExitCode     *int            `json:"exitCode,omitempty"`
	TraceID      *string         `json:"traceId,omitempty"`
}
