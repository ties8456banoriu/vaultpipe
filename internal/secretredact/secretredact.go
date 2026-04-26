// Package secretredact provides a pipeline step that redacts secret values
// in a map based on configurable key patterns. Redacted values are replaced
// with a placeholder string, making it safe to log or display secrets without
// exposing sensitive data.
package secretredact

import (
	"errors"
	"fmt"
	"path"
	"strings"
)

// DefaultPlaceholder is used when no custom placeholder is configured.
const DefaultPlaceholder = "[REDACTED]"

// Redactor replaces secret values matching configured patterns with a placeholder.
type Redactor struct {
	patterns    []string
	placeholder string
}

// Option is a functional option for Redactor.
type Option func(*Redactor)

// WithPlaceholder sets a custom placeholder string.
func WithPlaceholder(p string) Option {
	return func(r *Redactor) {
		if p != "" {
			r.placeholder = p
		}
	}
}

// New creates a Redactor that redacts values whose keys match any of the
// provided glob patterns. At least one non-empty pattern is required.
func New(patterns []string, opts ...Option) (*Redactor, error) {
	var valid []string
	for _, p := range patterns {
		trimmed := strings.TrimSpace(p)
		if trimmed == "" {
			continue
		}
		// Validate the glob pattern is syntactically correct.
		if _, err := path.Match(trimmed, ""); err != nil {
			return nil, fmt.Errorf("secretredact: invalid pattern %q: %w", trimmed, err)
		}
		valid = append(valid, trimmed)
	}
	if len(valid) == 0 {
		return nil, errors.New("secretredact: at least one pattern is required")
	}

	r := &Redactor{
		patterns:    valid,
		placeholder: DefaultPlaceholder,
	}
	for _, o := range opts {
		o(r)
	}
	return r, nil
}

// Apply returns a copy of secrets where any value whose key matches a
// configured pattern is replaced with the placeholder. Keys that do not
// match any pattern are left unchanged.
func (r *Redactor) Apply(secrets map[string]string) (map[string]string, error) {
	if len(secrets) == 0 {
		return nil, errors.New("secretredact: secrets map is empty")
	}

	out := make(map[string]string, len(secrets))
	for k, v := range secrets {
		if r.matches(k) {
			out[k] = r.placeholder
		} else {
			out[k] = v
		}
	}
	return out, nil
}

// Patterns returns the list of active glob patterns.
func (r *Redactor) Patterns() []string {
	cp := make([]string, len(r.patterns))
	copy(cp, r.patterns)
	return cp
}

// matches reports whether key matches any of the configured patterns.
func (r *Redactor) matches(key string) bool {
	for _, p := range r.patterns {
		if ok, _ := path.Match(p, key); ok {
			return true
		}
	}
	return false
}
