package research

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// registryClient provides shared HTTP behavior for package registry lookups.
type registryClient struct {
	httpClient *http.Client
	baseURL    string
}

func newRegistryClient(baseURL string) *registryClient {
	return &registryClient{
		httpClient: &http.Client{Timeout: 10 * time.Second},
		baseURL:    strings.TrimRight(baseURL, "/"),
	}
}

// newRegistryClientWithTimeout creates a registry client with a configurable
// HTTP timeout.  Use this when the caller wants to override the default 10 s.
func newRegistryClientWithTimeout(baseURL string, timeout time.Duration) *registryClient {
	return &registryClient{
		httpClient: &http.Client{Timeout: timeout},
		baseURL:    strings.TrimRight(baseURL, "/"),
	}
}

func (c *registryClient) get(ctx context.Context, path string, headers map[string]string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+path, nil)
	if err != nil {
		return nil, err
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	return c.httpClient.Do(req)
}

func (c *registryClient) readJSON(resp *http.Response, v any) error {
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	return json.Unmarshal(body, v)
}

// GoProxyClient looks up Go modules via a Go module proxy.
type GoProxyClient struct {
	client *registryClient
}

// NewGoProxyClient creates a Go proxy client with the given base URL.
func NewGoProxyClient(baseURL string) *GoProxyClient {
	return &GoProxyClient{client: newRegistryClient(baseURL)}
}

// NewGoProxyClientWithTimeout creates a Go proxy client with a configurable timeout.
func NewGoProxyClientWithTimeout(baseURL string, timeout time.Duration) *GoProxyClient {
	return &GoProxyClient{client: newRegistryClientWithTimeout(baseURL, timeout)}
}

// Lookup fetches metadata for a Go module from the proxy.
func (c *GoProxyClient) Lookup(ctx context.Context, pkg string) (*PackageInfo, error) {
	type latest struct {
		Version string `json:"Version"`
		Time    string `json:"Time"`
	}
	resp, err := c.client.get(ctx, "/"+pkg+"/@latest", nil)
	if err != nil {
		return nil, err
	}
	var l latest
	if err := c.client.readJSON(resp, &l); err != nil {
		return nil, err
	}
	return &PackageInfo{
		Name:          pkg,
		LatestVersion: l.Version,
		Published:     l.Time,
		Source:        "proxy.golang.org",
	}, nil
}

// NpmClient looks up packages on the npm registry.
type NpmClient struct {
	client *registryClient
}

// NewNpmClient creates an npm registry client with the given base URL.
func NewNpmClient(baseURL string) *NpmClient {
	return &NpmClient{client: newRegistryClient(baseURL)}
}

// NewNpmClientWithTimeout creates an npm registry client with a configurable timeout.
func NewNpmClientWithTimeout(baseURL string, timeout time.Duration) *NpmClient {
	return &NpmClient{client: newRegistryClientWithTimeout(baseURL, timeout)}
}

// Lookup fetches metadata for an npm package.
func (c *NpmClient) Lookup(ctx context.Context, pkg string) (*PackageInfo, error) {
	type npmPkg struct {
		Name     string `json:"name"`
		Version  string `json:"version"`
		License  string `json:"license"`
		Homepage string `json:"homepage"`
	}
	resp, err := c.client.get(ctx, "/"+pkg, map[string]string{
		"Accept": "application/vnd.npm.install-v1+json",
	})
	if err != nil {
		return nil, err
	}
	var p npmPkg
	if err := c.client.readJSON(resp, &p); err != nil {
		return nil, err
	}
	return &PackageInfo{
		Name:          p.Name,
		LatestVersion: p.Version,
		License:       p.License,
		Homepage:      p.Homepage,
		Source:        "registry.npmjs.org",
	}, nil
}

// PypiClient looks up packages on the Python Package Index.
type PypiClient struct {
	client *registryClient
}

// NewPypiClient creates a PyPI client with the given base URL.
func NewPypiClient(baseURL string) *PypiClient {
	return &PypiClient{client: newRegistryClient(baseURL)}
}

// NewPypiClientWithTimeout creates a PyPI client with a configurable timeout.
func NewPypiClientWithTimeout(baseURL string, timeout time.Duration) *PypiClient {
	return &PypiClient{client: newRegistryClientWithTimeout(baseURL, timeout)}
}

// Lookup fetches metadata for a PyPI package.
func (c *PypiClient) Lookup(ctx context.Context, pkg string) (*PackageInfo, error) {
	type info struct {
		Name     string `json:"name"`
		Version  string `json:"version"`
		License  string `json:"license"`
		HomePage string `json:"home_page"`
	}
	type response struct {
		Info info `json:"info"`
	}
	resp, err := c.client.get(ctx, "/pypi/"+pkg+"/json", nil)
	if err != nil {
		return nil, err
	}
	var r response
	if err := c.client.readJSON(resp, &r); err != nil {
		return nil, err
	}
	return &PackageInfo{
		Name:          r.Info.Name,
		LatestVersion: r.Info.Version,
		License:       r.Info.License,
		Homepage:      r.Info.HomePage,
		Source:        "pypi.org",
	}, nil
}

// CargoClient looks up crates on crates.io.
type CargoClient struct {
	client *registryClient
}

// NewCargoClient creates a crates.io client with the given base URL.
func NewCargoClient(baseURL string) *CargoClient {
	return &CargoClient{client: newRegistryClient(baseURL)}
}

// NewCargoClientWithTimeout creates a crates.io client with a configurable timeout.
func NewCargoClientWithTimeout(baseURL string, timeout time.Duration) *CargoClient {
	return &CargoClient{client: newRegistryClientWithTimeout(baseURL, timeout)}
}

// Lookup fetches metadata for a crate.
func (c *CargoClient) Lookup(ctx context.Context, pkg string) (*PackageInfo, error) {
	type crate struct {
		Name             string `json:"name"`
		MaxStableVersion string `json:"max_stable_version"`
		Homepage         string `json:"homepage"`
		License          string `json:"license"`
		UpdatedAt        string `json:"updated_at"`
	}
	type response struct {
		Crate crate `json:"crate"`
	}
	resp, err := c.client.get(ctx, "/api/v1/crates/"+pkg, nil)
	if err != nil {
		return nil, err
	}
	var r response
	if err := c.client.readJSON(resp, &r); err != nil {
		return nil, err
	}
	return &PackageInfo{
		Name:          r.Crate.Name,
		LatestVersion: r.Crate.MaxStableVersion,
		License:       r.Crate.License,
		Homepage:      r.Crate.Homepage,
		Published:     r.Crate.UpdatedAt,
		Source:        "crates.io",
	}, nil
}
