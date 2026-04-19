// Package secretslice provides a utility for extracting an ordered slice
// of key-value pairs from a secrets map, with optional key filtering.
package secretslice

import (
	"errors"
	"sort"
)

// Entry represents a single secret key-value pair.
type Entry struct {
	Key   string
	Value string
}

// Slicer converts a secrets map into an ordered slice of entries.
type Slicer struct {
	keys []string // if non-empty, only these keys are included (in order)
}

// Option configures a Slicer.
type Option func(*Slicer)

// WithKeys restricts the output to the given keys, preserving their order.
func WithKeys(keys []string) Option {
	return func(s *Slicer) {
		s.keys = keys
	}
}

// New creates a new Slicer with the given options.
func New(opts ...Option) *Slicer {
	s := &Slicer{}
	for _, o := range opts {
		o(s)
	}
	return s
}

// Apply converts the secrets map to an ordered slice of Entry values.
// If WithKeys was provided, only those keys are included in that order.
// Otherwise all keys are returned sorted alphabetically.
// Returns an error if secrets is empty.
func (s *Slicer) Apply(secrets map[string]string) ([]Entry, error) {
	if len(secrets) == 0 {
		return nil, errors.New("secretslice: secrets map is empty")
	}

	if len(s.keys) > 0 {
		var entries []Entry
		for _, k := range s.keys {
			v, ok := secrets[k]
			if !ok {
				return nil, errors.New("secretslice: key not found: " + k)
			}
			entries = append(entries, Entry{Key: k, Value: v})
		}
		return entries, nil
	}

	keys := make([]string, 0, len(secrets))
	for k := range secrets {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	entries := make([]Entry, 0, len(keys))
	for _, k := range keys {
		entries = append(entries, Entry{Key: k, Value: secrets[k]})
	}
	return entries, nil
}
