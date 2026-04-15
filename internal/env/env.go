// Package env provides utilities for loading and resolving environment
// variable overrides that take precedence over Vault-fetched secrets.
package env

import (
	"errors"
	"os"
	"strings"
)

// ErrNoOverrides is returned when no environment overrides are found.
var ErrNoOverrides = errors.New("env: no overrides found")

// Resolver reads environment variables and returns a map of key→value
// pairs that should override secrets fetched from Vault.
type Resolver struct {
	prefix string
	env    func(string) (string, bool)
}

// NewResolver creates a Resolver that matches environment variables
// beginning with prefix (e.g. "VAULTPIPE_OVERRIDE_").
func NewResolver(prefix string) *Resolver {
	return &Resolver{
		prefix: strings.ToUpper(prefix),
		env:    os.LookupEnv,
	}
}

// Resolve scans os.Environ for variables matching the configured prefix
// and returns them as a map with the prefix stripped from each key.
// Returns ErrNoOverrides if no matching variables exist.
func (r *Resolver) Resolve() (map[string]string, error) {
	result := make(map[string]string)

	for _, entry := range os.Environ() {
		k, v, ok := strings.Cut(entry, "=")
		if !ok {
			continue
		}
		upper := strings.ToUpper(k)
		if strings.HasPrefix(upper, r.prefix) {
			stripped := upper[len(r.prefix):]
			if stripped != "" {
				result[stripped] = v
			}
		}
	}

	if len(result) == 0 {
		return nil, ErrNoOverrides
	}
	return result, nil
}

// Apply merges overrides on top of base, returning a new map.
// Keys present in overrides replace those in base; all other base
// keys are preserved unchanged.
func Apply(base, overrides map[string]string) map[string]string {
	merged := make(map[string]string, len(base))
	for k, v := range base {
		merged[k] = v
	}
	for k, v := range overrides {
		merged[k] = v
	}
	return merged
}
