// Package secretjoin combines multiple secret values into a single key
// using a configurable separator.
package secretjoin

import (
	"errors"
	"fmt"
	"strings"
)

// Rule defines how to join multiple secret keys into one.
type Rule struct {
	Keys      []string
	TargetKey string
	Separator string
}

// Joiner combines secret values according to join rules.
type Joiner struct {
	rules []Rule
}

// New creates a Joiner with the given rules.
// Returns an error if no rules are provided or any rule is invalid.
func New(rules []Rule) (*Joiner, error) {
	if len(rules) == 0 {
		return nil, errors.New("secretjoin: at least one rule is required")
	}
	for i, r := range rules {
		if r.TargetKey == "" {
			return nil, fmt.Errorf("secretjoin: rule %d: target key must not be empty", i)
		}
		if len(r.Keys) < 2 {
			return nil, fmt.Errorf("secretjoin: rule %d: at least two source keys required", i)
		}
	}
	return &Joiner{rules: rules}, nil
}

// Apply produces a new secrets map with joined values appended.
// Original keys are preserved. Returns error if secrets is empty.
func (j *Joiner) Apply(secrets map[string]string) (map[string]string, error) {
	if len(secrets) == 0 {
		return nil, errors.New("secretjoin: secrets must not be empty")
	}
	out := make(map[string]string, len(secrets))
	for k, v := range secrets {
		out[k] = v
	}
	for _, r := range j.rules {
		parts := make([]string, 0, len(r.Keys))
		for _, k := range r.Keys {
			v, ok := secrets[k]
			if !ok {
				return nil, fmt.Errorf("secretjoin: source key %q not found", k)
			}
			parts = append(parts, v)
		}
		out[r.TargetKey] = strings.Join(parts, r.Separator)
	}
	return out, nil
}
