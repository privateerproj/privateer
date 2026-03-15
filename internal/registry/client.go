// Package registry provides a client for the Privateer plugin registry (vetted plugins and plugin metadata).
package registry

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"
)

const (
	// DefaultBaseURL is the default base URL for the plugin registry.
	DefaultBaseURL = "https://revanite.io/privateer"
	// defaultHTTPTimeout is the timeout for registry HTTP requests.
	defaultHTTPTimeout = 15 * time.Second
	// userAgent identifies the client.
	userAgent = "pvtr/1.0"
)

// VettedListResponse is the response from the vetted plugins list endpoint.
type VettedListResponse struct {
	Message string   `json:"message"`
	Updated string   `json:"updated"`
	Plugins []string `json:"plugins"`
}

// PluginData is the per-plugin metadata from the registry.
type PluginData struct {
	Name              string   `json:"name"`
	Updated           string   `json:"updated"`
	Source            string   `json:"source"`
	Image             string   `json:"image"`
	Latest            string   `json:"latest"`
	SupportedVersions []string `json:"supported-versions"`
	BinaryPathInImage string   `json:"binaryPath,omitempty"` // optional path inside OCI image
	DownloadURL       string   `json:"download,omitempty"`   // optional direct URL to plugin binary
}

// Client fetches vetted plugin list and plugin metadata from the registry.
type Client struct {
	BaseURL    string
	HTTPClient *http.Client
}

// NewClient returns a registry client with default base URL and HTTP client.
func NewClient() *Client {
	baseURL := os.Getenv("PVTR_REGISTRY_URL")
	if baseURL == "" {
		baseURL = DefaultBaseURL
	}
	return &Client{
		BaseURL: strings.TrimSuffix(baseURL, "/"),
		HTTPClient: &http.Client{
			Timeout: defaultHTTPTimeout,
		},
	}
}

// GetVettedList fetches the list of vetted plugin names.
// Returns an error on network failure or non-200 response.
func (c *Client) GetVettedList() (*VettedListResponse, error) {
	url := c.BaseURL + "/vetted-plugins.json"
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("build request: %w", err)
	}
	req.Header.Set("User-Agent", userAgent)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetch vetted list: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("vetted list returned status %d", resp.StatusCode)
	}

	var out VettedListResponse
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, fmt.Errorf("decode vetted list: %w", err)
	}
	return &out, nil
}

// GetPluginData fetches metadata for a single plugin by name (e.g. "ossf/pvtr-github-repo-scanner").
// Returns an error if the plugin is not found (404) or on network/parse failure.
func (c *Client) GetPluginData(name string) (*PluginData, error) {
	if name == "" || strings.Contains(name, "..") {
		return nil, fmt.Errorf("invalid plugin name: %q", name)
	}
	// Name may contain slashes; path is plugin-data/owner/repo.json
	url := c.BaseURL + "/plugin-data/" + name + ".json"
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("build request: %w", err)
	}
	req.Header.Set("User-Agent", userAgent)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetch plugin data: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("plugin %q not found in vetted registry", name)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("plugin data returned status %d", resp.StatusCode)
	}

	var out PluginData
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, fmt.Errorf("decode plugin data: %w", err)
	}
	return &out, nil
}
