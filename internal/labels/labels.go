// Package labels provides key-value label tagging for secrets,
// allowing secrets to be annotated with arbitrary metadata (e.g. environment, team, tier).
package labels

import (
	"errors"
	"fmt"
	"strings"
)

// Set holds labels attached to a secret key.
type Set map[string]string

// Store manages labels per secret env key.
type Store struct {
	entries map[string]Set
}

// New returns an initialised Store.
func New() *Store {
	return &Store{entries: make(map[string]Set)}
}

// Tag attaches labels to the given env key. Existing labels are merged;
// new values overwrite duplicates.
func (s *Store) Tag(envKey string, labels Set) error {
	if envKey == "" {
		return errors.New("labels: envKey must not be empty")
	}
	if len(labels) == 0 {
		return errors.New("labels: labels must not be empty")
	}
	if _, ok := s.entries[envKey]; !ok {
		s.entries[envKey] = make(Set)
	}
	for k, v := range labels {
		s.entries[envKey][k] = v
	}
	return nil
}

// Get returns the label set for an env key.
func (s *Store) Get(envKey string) (Set, error) {
	set, ok := s.entries[envKey]
	if !ok {
		return nil, fmt.Errorf("labels: no labels for key %q", envKey)
	}
	copy := make(Set, len(set))
	for k, v := range set {
		copy[k] = v
	}
	return copy, nil
}

// Filter returns all env keys whose labels match every selector in the given set.
func (s *Store) Filter(selector Set) []string {
	var result []string
	for envKey, labels := range s.entries {
		if matchesAll(labels, selector) {
			result = append(result, envKey)
		}
	}
	return result
}

// ParseRules parses label strings of the form "KEY=k1=v1,k2=v2".
func ParseRules(rules []string) (map[string]Set, error) {
	out := make(map[string]Set)
	for _, rule := range rules {
		parts := strings.SplitN(rule, "=", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("labels: invalid rule %q, expected KEY=k1=v1,...", rule)
		}
		envKey := parts[0]
		set := make(Set)
		for _, pair := range strings.Split(parts[1], ",") {
			kv := strings.SplitN(pair, "=", 2)
			if len(kv) != 2 {
				return nil, fmt.Errorf("labels: invalid label pair %q", pair)
			}
			set[kv[0]] = kv[1]
		}
		out[envKey] = set
	}
	return out, nil
}

func matchesAll(labels, selector Set) bool {
	for k, v := range selector {
		if labels[k] != v {
			return false
		}
	}
	return true
}
