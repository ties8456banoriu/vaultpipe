// Package secretstrip removes secrets from a map by key pattern.
package secretstrip

import (
	"errors"
	"fmt"
	"path"
)

// Stripper removes secrets whose keys match one or more glob patterns.
type Stripper struct {
	patterns []string
}

// New returns a Stripper that will remove keys matching any of the given glob
// patterns. At least one pattern must be provided.
func New(patterns []string) (*Stripper, error) {
	if len(patterns) == 0 {
		return nil, errors.New("secretstrip: at least one pattern is required")
	}
	for _, p := range patterns {
		if p == "" {
			return nil, errors.New("secretstrip: pattern must not be empty")
		}
		// validate pattern syntax early
		if _, err := path.Match(p, ""); err != nil {
			return nil, fmt.Errorf("secretstrip: invalid pattern %q: %w", p, err)
		}
	}
	return &Stripper{patterns: patterns}, nil
}

// Apply returns a copy of secrets with all keys matching any configured
// pattern removed. Returns an error if secrets is empty.
func (s *Stripper) Apply(secrets map[string]string) (map[string]string, error) {
	if len(secrets) == 0 {
		return nil, errors.New("secretstrip: secrets must not be empty")
	}
	out := make(map[string]string, len(secrets))
	for k, v := range secrets {
		if s.matches(k) {
			continue
		}
		out[k] = v
	}
	return out, nil
}

func (s *Stripper) matches(key string) bool {
	for _, p := range s.patterns {
		ok, _ := path.Match(p, key)
		if ok {
			return true
		}
	}
	return false
}
