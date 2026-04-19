// Package secretpatch applies partial updates (patches) to a secrets map,
// merging only the specified keys from a patch set into the base secrets.
package secretpatch

import (
	"errors"
	"fmt"
)

// ErrEmptySecrets is returned when the base or patch secrets map is empty.
var (
	ErrEmptyBase  = errors.New("secretpatch: base secrets must not be empty")
	ErrEmptyPatch = errors.New("secretpatch: patch secrets must not be empty")
	ErrEmptyKey   = errors.New("secretpatch: patch key must not be empty")
)

// Policy controls how conflicts between base and patch are resolved.
type Policy string

const (
	PolicyOverwrite Policy = "overwrite" // patch value wins
	PolicyKeep      Policy = "keep"      // base value wins on conflict
)

// Patcher applies a patch map to a base secrets map.
type Patcher struct {
	policy Policy
}

// New creates a new Patcher with the given conflict resolution policy.
func New(policy Policy) (*Patcher, error) {
	if policy != PolicyOverwrite && policy != PolicyKeep {
		return nil, fmt.Errorf("secretpatch: unknown policy %q", policy)
	}
	return &Patcher{policy: policy}, nil
}

// Apply merges patch into base according to the patcher's policy.
// Only keys present in keys are applied; if keys is nil all patch keys are applied.
func (p *Patcher) Apply(base, patch map[string]string, keys []string) (map[string]string, error) {
	if len(base) == 0 {
		return nil, ErrEmptyBase
	}
	if len(patch) == 0 {
		return nil, ErrEmptyPatch
	}

	result := make(map[string]string, len(base))
	for k, v := range base {
		result[k] = v
	}

	apply := patch
	if len(keys) > 0 {
		for _, k := range keys {
			if k == "" {
				return nil, ErrEmptyKey
			}
		}
		apply = make(map[string]string, len(keys))
		for _, k := range keys {
			if v, ok := patch[k]; ok {
				apply[k] = v
			}
		}
	}

	for k, v := range apply {
		if _, exists := result[k]; exists && p.policy == PolicyKeep {
			continue
		}
		result[k] = v
	}
	return result, nil
}
