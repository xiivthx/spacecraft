package research

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// BraveClient calls the Brave Search API.
type BraveClient struct {
	client  *http.Client
	apiKey  string
	baseURL string
}

// NewBraveClient creates a Brave Search client with a 10-second default timeout.
func NewBraveClient(apiKey, baseURL string) *BraveClient {
	return &BraveClient{
		client:  &http.Client{Timeout: 10 * time.Second},
		apiKey:  apiKey,
		baseURL: strings.TrimSuffix(baseURL, "/"),
	}
}

// Search performs a Brave web search scoped to the given domains.
// count limits the number of results requested (0 means use Brave's default).
func (c *BraveClient) Search(ctx context.Context, query string, domains []string, count int) ([]ResearchResult, error) {
	q := query
	if len(domains) > 0 {
		scopes := make([]string, 0, len(domains)+1)
		for _, d := range domains {
			scopes = append(scopes, "site:"+d)
		}
		scopes = append(scopes, query)
		q = strings.Join(scopes, " ")
	}

	u, err := url.Parse(c.baseURL + "/res/v1/web/search")
	if err != nil {
		return nil, fmt.Errorf("brave search: invalid url: %w", err)
	}
	values := u.Query()
	values.Set("q", q)
	if count > 0 {
		values.Set("count", fmt.Sprintf("%d", count))
	}
	u.RawQuery = values.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("brave search: create request: %w", err)
	}
	req.Header.Set("X-Subscription-Token", c.apiKey)
	req.Header.Set("Accept", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("brave search: do request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("brave search returned status %d", resp.StatusCode)
	}

	var payload struct {
		Web struct {
			Results []struct {
				Title       string `json:"title"`
				URL         string `json:"url"`
				Description string `json:"description"`
			} `json:"results"`
		} `json:"web"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return nil, fmt.Errorf("brave search: decode response: %w", err)
	}

	results := make([]ResearchResult, 0, len(payload.Web.Results))
	for _, r := range payload.Web.Results {
		results = append(results, ResearchResult{
			Title:     r.Title,
			URL:       r.URL,
			Snippet:   r.Description,
			Relevance: 0.0,
		})
	}
	return results, nil
}
