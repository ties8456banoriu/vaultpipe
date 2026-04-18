// Package secrettag provides tagging of secrets with arbitrary string tags
// for grouping, filtering, and annotation purposes.
package secrettag

import (
	"errors"
	"sync"
)

// ErrEmptyKey is returned when an empty env key is provided.
var ErrEmptyKey = errors.New("secrettag: env key must not be empty")

// ErrEmptyTags is returned when no tags are provided.
var ErrEmptyTags = errors.New("secrettag: tags must not be empty")

// ErrUnknownKey is returned when a key has no tags recorded.
var ErrUnknownKey = errors.New("secrettag: no tags for key")

// Tagger stores string tags associated with env keys.
type Tagger struct {
	mu   sync.RWMutex
	data map[string][]string
}

// New returns a new Tagger.
func New() *Tagger {
	return &Tagger{data: make(map[string][]string)}
}

// Tag associates tags with the given env key.
func (t *Tagger) Tag(envKey string, tags []string) error {
	if envKey == "" {
		return ErrEmptyKey
	}
	if len(tags) == 0 {
		return ErrEmptyTags
	}
	copy := make([]string, len(tags))
	for i, v := range tags {
		copy[i] = v
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	t.data[envKey] = copy
	return nil
}

// Get returns the tags for the given env key.
func (t *Tagger) Get(envKey string) ([]string, error) {
	if envKey == "" {
		return nil, ErrEmptyKey
	}
	t.mu.RLock()
	defer t.mu.RUnlock()
	tags, ok := t.data[envKey]
	if !ok {
		return nil, ErrUnknownKey
	}
	out := make([]string, len(tags))
	copy(out, tags)
	return out, nil
}

// HasTag reports whether the given env key has the specified tag.
func (t *Tagger) HasTag(envKey, tag string) bool {
	tags, err := t.Get(envKey)
	if err != nil {
		return false
	}
	for _, v := range tags {
		if v == tag {
			return true
		}
	}
	return false
}

// All returns a copy of all recorded tags keyed by env key.
func (t *Tagger) All() map[string][]string {
	t.mu.RLock()
	defer t.mu.RUnlock()
	out := make(map[string][]string, len(t.data))
	for k, v := range t.data {
		cp := make([]string, len(v))
		copy(cp, v)
		out[k] = cp
	}
	return out
}
