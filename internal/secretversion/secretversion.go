// Package secretversion tracks the version history of secrets fetched from Vault.
package secretversion

import (
	"errors"
	"sync"
	"time"
)

// Entry records a single version snapshot for a secret key.
type Entry struct {
	EnvKey    string
	VaultPath string
	Version   int
	FetchedAt time.Time
}

// Tracker maintains an ordered version history per env key.
type Tracker struct {
	mu      sync.RWMutex
	history map[string][]Entry
}

// New returns a new Tracker.
func New() *Tracker {
	return &Tracker{history: make(map[string][]Entry)}
}

// Record appends a version entry for the given env key.
func (t *Tracker) Record(envKey, vaultPath string, version int, fetchedAt time.Time) error {
	if envKey == "" {
		return errors.New("secretversion: envKey must not be empty")
	}
	if vaultPath == "" {
		return errors.New("secretversion: vaultPath must not be empty")
	}
	if version < 0 {
		return errors.New("secretversion: version must be non-negative")
	}
	if fetchedAt.IsZero() {
		fetchedAt = time.Now().UTC()
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	t.history[envKey] = append(t.history[envKey], Entry{
		EnvKey:    envKey,
		VaultPath: vaultPath,
		Version:   version,
		FetchedAt: fetchedAt,
	})
	return nil
}

// Latest returns the most recently recorded entry for the given env key.
func (t *Tracker) Latest(envKey string) (Entry, error) {
	t.mu.RLock()
	defer t.mu.RUnlock()
	entries, ok := t.history[envKey]
	if !ok || len(entries) == 0 {
		return Entry{}, errors.New("secretversion: no history for key: " + envKey)
	}
	return entries[len(entries)-1], nil
}

// All returns a copy of the full history for the given env key.
func (t *Tracker) All(envKey string) ([]Entry, error) {
	t.mu.RLock()
	defer t.mu.RUnlock()
	entries, ok := t.history[envKey]
	if !ok {
		return nil, errors.New("secretversion: no history for key: " + envKey)
	}
	copy := make([]Entry, len(entries))
	for i, e := range entries {
		copy[i] = e
	}
	return copy, nil
}
