package vault

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Client is a minimal Vault HTTP client.
type Client struct {
	baseURL    string
	token      string
	httpClient *http.Client
}

// SecretData holds the key/value pairs returned from a Vault KV secret.
type SecretData map[string]string

// NewClient creates a new Vault client with the given address and token.
func NewClient(address, token string) *Client {
	return &Client{
		baseURL: address,
		token:   token,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// GetSecret fetches a KV v2 secret at the given path.
// path should be in the form "secret/data/myapp".
func (c *Client) GetSecret(path string) (SecretData, error) {
	url := fmt.Sprintf("%s/v1/%s", c.baseURL, path)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}
	req.Header.Set("X-Vault-Token", c.token)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("executing request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("vault returned status %d: %s", resp.StatusCode, string(body))
	}

	return parseKVv2Response(body)
}

// parseKVv2Response extracts the data map from a KV v2 JSON response.
func parseKVv2Response(body []byte) (SecretData, error) {
	var envelope struct {
		Data struct {
			Data map[string]interface{} `json:"data"`
		} `json:"data"`
	}
	if err := json.Unmarshal(body, &envelope); err != nil {
		return nil, fmt.Errorf("parsing vault response: %w", err)
	}

	result := make(SecretData, len(envelope.Data.Data))
	for k, v := range envelope.Data.Data {
		result[k] = fmt.Sprintf("%v", v)
	}
	return result, nil
}
