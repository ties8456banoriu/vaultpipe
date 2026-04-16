// Package scopedenv provides scoping of secrets to named environments
// (e.g. "dev", "staging", "prod"), allowing callers to namespace secret
// keys and merge only the relevant subset.
package scopedenv

import (
	"errors"
	"fmt"
	"strings"
)

// ErrNoScope is returned when no scope is provided.
var ErrNoScope = errors.New("scopedenv: scope name must not be empty")

// ErrEmptySecrets is returned when the secrets map is empty.
var ErrEmptySecrets = errors.New("scopedenv: secrets must not be empty")

// Scoper namespaces and filters secrets by environment scope.
type Scoper struct {
	scope string
}

// New creates a Scoper for the given scope name (e.g. "dev").
func New(scope string) (*Scoper, error) {
	scope = strings.TrimSpace(scope)
	if scope == "" {
		return nil, ErrNoScope
	}
	return &Scoper{scope: strings.ToLower(scope)}, nil
}

// Tag prefixes each key with "<SCOPE>_" and returns the new map.
func (s *Scoper) Tag(secrets map[string]string) (map[string]string, error) {
	if len(secrets) == 0 {
		return nil, ErrEmptySecrets
	}
	out := make(map[string]string, len(secrets))
	prefix := strings.ToUpper(s.scope) + "_"
	for k, v := range secrets {
		out[prefix+k] = v
	}
	return out, nil
}

// Filter returns only keys that are prefixed with the scope, stripping the prefix.
func (s *Scoper) Filter(secrets map[string]string) (map[string]string, error) {
	if len(secrets) == 0 {
		return nil, ErrEmptySecrets
	}
	prefix := strings.ToUpper(s.scope) + "_"
	out := make(map[string]string)
	for k, v := range secrets {
		if strings.HasPrefix(k, prefix) {
			out[strings.TrimPrefix(k, prefix)] = v
		}
	}
	if len(out) == 0 {
		return nil, fmt.Errorf("scopedenv: no keys matched scope %q", s.scope)
	}
	return out, nil
}

// Scope returns the scoper's current scope name.
func (s *Scoper) Scope() string { return s.scope }
