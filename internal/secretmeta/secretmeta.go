// Package secretmeta tracks metadata about secrets fetched from Vault,
// including version, mount path, and fetch timestamp.
package secretmeta

import (
	"errors"
	"sync"
	"time"
)

// Meta holds metadata for a single secret key.
type Meta struct {
	EnvKey    string
	VaultPath string
	Mount     string
	Version   int
	FetchedAt time.Time
}

// Store holds metadata records keyed by env variable name.
type Store struct {
	mu      sync.RWMutex
	records map[string]Meta
}

// New returns an initialised Store.
func New() *Store {
	return &Store{records: make(map[string]Meta)}
}

// Record stores metadata for the given env key.
func (s *Store) Record(m Meta) error {
	if m.EnvKey == "" {
		return errors.New("secretmeta: env key must not be empty")
	}
	if m.VaultPath == "" {
		return errors.New("secretmeta: vault path must not be empty")
	}
	if m.FetchedAt.IsZero() {
		m.FetchedAt = time.Now().UTC()
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.records[m.EnvKey] = m
	return nil
}

// Get returns the metadata for the given env key.
func (s *Store) Get(envKey string) (Meta, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	m, ok := s.records[envKey]
	if !ok {
		return Meta{}, errors.New("secretmeta: no metadata for key: " + envKey)
	}
	return m, nil
}

// All returns a copy of all stored metadata records.
func (s *Store) All() []Meta {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]Meta, 0, len(s.records))
	for _, m := range s.records {
		out = append(out, m)
	}
	return out
}

// Delete removes the metadata record for the given env key.
func (s *Store) Delete(envKey string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.records, envKey)
}
