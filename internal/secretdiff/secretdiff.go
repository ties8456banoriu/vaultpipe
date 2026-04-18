// Package secretdiff provides a tracker for recording and retrieving
// secret diff events between successive fetches.
package secretdiff

import (
	"errors"
	"sync"
	"time"
)

// Event represents a single diff event for a secret key.
type Event struct {
	EnvKey    string
	VaultPath string
	OldValue  string
	NewValue  string
	ChangedAt time.Time
}

// Tracker records diff events keyed by env key.
type Tracker struct {
	mu     sync.RWMutex
	events map[string]Event
}

// New returns a new Tracker.
func New() *Tracker {
	return &Tracker{events: make(map[string]Event)}
}

// Record stores a diff event. ChangedAt is set to now if zero.
func (t *Tracker) Record(e Event) error {
	if e.EnvKey == "" {
		return errors.New("secretdiff: env key must not be empty")
	}
	if e.VaultPath == "" {
		return errors.New("secretdiff: vault path must not be empty")
	}
	if e.ChangedAt.IsZero() {
		e.ChangedAt = time.Now().UTC()
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	t.events[e.EnvKey] = e
	return nil
}

// Get returns the most recent diff event for the given env key.
func (t *Tracker) Get(envKey string) (Event, error) {
	t.mu.RLock()
	defer t.mu.RUnlock()
	e, ok := t.events[envKey]
	if !ok {
		return Event{}, errors.New("secretdiff: no event for key: " + envKey)
	}
	return e, nil
}

// All returns a copy of all recorded diff events.
func (t *Tracker) All() []Event {
	t.mu.RLock()
	defer t.mu.RUnlock()
	out := make([]Event, 0, len(t.events))
	for _, e := range t.events {
		out = append(out, e)
	}
	return out
}

// Clear removes all recorded events.
func (t *Tracker) Clear() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.events = make(map[string]Event)
}
