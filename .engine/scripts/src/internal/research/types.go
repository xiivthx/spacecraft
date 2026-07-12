package research

import "context"

// ResearchResult represents a single external research hit.
type ResearchResult struct {
	Title     string
	URL       string
	Snippet   string
	Relevance float64
}

// PackageInfo holds metadata about a package in a software registry.
type PackageInfo struct {
	Name          string
	LatestVersion string
	License       string
	Published     string
	Homepage      string
	Source        string
}

// SearchScope restricts research to a set of domains with a description.
type SearchScope struct {
	Domains     []string
	Description string
}

// RegistryQuery identifies a package lookup within a specific ecosystem.
type RegistryQuery struct {
	PackageName string
	Ecosystem   string
}

// SearchOptions controls optional behavior for research searches.
type SearchOptions struct{}

// DeepOptions controls optional behavior for deep URL analysis.
type DeepOptions struct{}

// DeepResult contains a synthesized deep analysis of a web page.
type DeepResult struct {
	Summary   string
	KeyPoints []string
	SourceURL string
	FetchedAt string
}

// ResearchProvider abstracts external research capabilities.
type ResearchProvider interface {
	Search(ctx context.Context, query string, opts SearchOptions) ([]ResearchResult, error)
	LookupPackage(ctx context.Context, q RegistryQuery) (*PackageInfo, error)
	DeepAnalyze(ctx context.Context, url string, opts DeepOptions) (*DeepResult, error)
}

// ResearchStore abstracts persistence for external research results.
type ResearchStore interface {
	Save(ctx context.Context, result any) error
	Load(ctx context.Context, path string) ([]byte, error)
}
