package roadmap

import (
	"encoding/json"
	"time"
)

type Issue struct {
	Number int      `json:"number"`
	Title  string   `json:"title"`
	URL    string   `json:"url"`
	State  string   `json:"state"`
	Labels []string `json:"labels,omitempty"`
	Phase  string   `json:"phase,omitempty"`
}

type MissionEntry struct {
	ID          string `json:"id"`
	Description string `json:"description,omitempty"`
}

type Roadmap struct {
	ID          string         `json:"id"`
	Title       string         `json:"title"`
	Description string         `json:"description"`
	Missions    []MissionEntry `json:"missions"`
	Issues      []Issue        `json:"issues,omitempty"`
	CreatedAt   time.Time      `json:"createdAt"`
	UpdatedAt   time.Time      `json:"updatedAt"`
}

func (r *Roadmap) UnmarshalJSON(data []byte) error {
	type Alias Roadmap
	aux := &struct {
		Missions []json.RawMessage `json:"missions"`
		*Alias
	}{
		Alias: (*Alias)(r),
	}
	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}
	r.Missions = make([]MissionEntry, 0, len(aux.Missions))
	for _, raw := range aux.Missions {
		var entry MissionEntry
		if err := json.Unmarshal(raw, &entry); err == nil {
			r.Missions = append(r.Missions, entry)
		} else {
			var id string
			if err := json.Unmarshal(raw, &id); err != nil {
				return err
			}
			r.Missions = append(r.Missions, MissionEntry{ID: id})
		}
	}
	return nil
}
