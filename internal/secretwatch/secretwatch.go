// Package secretwatch tracks which secrets have been accessed during a session.
package secretwatch

import (
	"errors"
	"sync"
	"time"
)

// Access records a single secret access event.
type Access struct {
	EnvKey    string
	VaultPath string
	AccessedAt time.Time
}

// Watcher tracks secret access events.
type Watcher struct {
	mu      sync.Mutex
	events  []Access
	seen    map[string]int // envKey -> index in events
}

// New returns a new Watcher.
func New() *Watcher {
	return &Watcher{
		seen: make(map[string]int),
	}
}

// Record registers an access event for the given env key and vault path.
func (w *Watcher) Record(envKey, vaultPath string) error {
	if envKey == "" {
		return errors.New("secretwatch: envKey must not be empty")
	}
	if vaultPath == "" {
		return errors.New("secretwatch: vaultPath must not be empty")
	}
	w.mu.Lock()
	defer w.mu.Unlock()
	a := Access{
		EnvKey:     envKey,
		VaultPath:  vaultPath,
		AccessedAt: time.Now().UTC(),
	}
	if idx, ok := w.seen[envKey]; ok {
		w.events[idx] = a
		return nil
	}
	w.seen[envKey] = len(w.events)
	w.events = append(w.events, a)
	return nil
}

// All returns a copy of all recorded access events.
func (w *Watcher) All() []Access {
	w.mu.Lock()
	defer w.mu.Unlock()
	out := make([]Access, len(w.events))
	copy(out, w.events)
	return out
}

// Has returns true if the given envKey has been recorded.
func (w *Watcher) Has(envKey string) bool {
	w.mu.Lock()
	defer w.mu.Unlock()
	_, ok := w.seen[envKey]
	return ok
}

// Clear removes all recorded events.
func (w *Watcher) Clear() {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.events = nil
	w.seen = make(map[string]int)
}
