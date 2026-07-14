package roadmap

import "time"

type Issue struct {
	Number int      `json:"number"`
	Title  string   `json:"title"`
	URL    string   `json:"url"`
	State  string   `json:"state"`
	Labels []string `json:"labels,omitempty"`
	Phase  string   `json:"phase,omitempty"`
}

type Roadmap struct {
	ID          string          `json:"id"`
	Title       string          `json:"title"`
	Description string          `json:"description"`
	Missions    []string        `json:"missions"`
	Issues      []Issue         `json:"issues,omitempty"`
	CreatedAt   time.Time       `json:"createdAt"`
	UpdatedAt   time.Time       `json:"updatedAt"`
}
