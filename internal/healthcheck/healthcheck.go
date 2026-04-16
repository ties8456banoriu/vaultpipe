// Package healthcheck provides a simple Vault connectivity check.
package healthcheck

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

// Result holds the outcome of a health check.
type Result struct {
	Reachable  bool
	StatusCode int
	Latency    time.Duration
	Err        error
}

// Checker performs health checks against a Vault instance.
type Checker struct {
	baseURL    string
	httpClient *http.Client
}

// New creates a Checker for the given Vault base URL.
func New(baseURL string, timeout time.Duration) *Checker {
	if timeout == 0 {
		timeout = 5 * time.Second
	}
	return &Checker{
		baseURL: baseURL,
		httpClient: &http.Client{Timeout: timeout},
	}
}

// Check calls the Vault sys/health endpoint and returns a Result.
func (c *Checker) Check(ctx context.Context) Result {
	url := fmt.Sprintf("%s/v1/sys/health", c.baseURL)
	start := time.Now()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return Result{Err: fmt.Errorf("build request: %w", err)}
	}

	resp, err := c.httpClient.Do(req)
	latency := time.Since(start)
	if err != nil {
		return Result{Latency: latency, Err: fmt.Errorf("vault unreachable: %w", err)}
	}
	defer resp.Body.Close()

	// Vault returns 200 (active), 429 (standby), 472/473 (DR/perf).
	// All indicate the server is reachable.
	reachable := resp.StatusCode < 500
	return Result{
		Reachable:  reachable,
		StatusCode: resp.StatusCode,
		Latency:    latency,
	}
}
