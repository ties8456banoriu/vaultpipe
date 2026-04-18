// Package secretarchive provides archival and retrieval of secret snapshots
// with named archive slots for point-in-time recovery.
package secretarchive

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

// ErrNotFound is returned when a named archive does not exist.
var ErrNotFound = errors.New("archive not found")

// ErrEmptyName is returned when an archive name is blank.
var ErrEmptyName = errors.New("archive name must not be empty")

// ErrEmptySecrets is returned when secrets map is empty.
var ErrEmptySecrets = errors.New("secrets must not be empty")

// Entry holds an archived snapshot of secrets.
type Entry struct {
	Name      string
	Secrets   map[string]string
	ArchivedAt time.Time
}

// Archive stores named secret snapshots.
type Archive struct {
	mu      sync.RWMutex
	entries map[string]Entry
}

// New creates a new Archive.
func New() *Archive {
	return &Archive{entries: make(map[string]Entry)}
}

// Store saves a named snapshot of secrets.
func (a *Archive) Store(name string, secrets map[string]string) error {
	if name == "" {
		return ErrEmptyName
	}
	if len(secrets) == 0 {
		return ErrEmptySecrets
	}
	copy := make(map[string]string, len(secrets))
	for k, v := range secrets {
		copy[k] = v
	}
	a.mu.Lock()
	defer a.mu.Unlock()
	a.entries[name] = Entry{Name: name, Secrets: copy, ArchivedAt: time.Now().UTC()}
	return nil
}

// Get retrieves a named archive entry.
func (a *Archive) Get(name string) (Entry, error) {
	if name == "" {
		return Entry{}, ErrEmptyName
	}
	a.mu.RLock()
	defer a.mu.RUnlock()
	e, ok := a.entries[name]
	if !ok {
		return Entry{}, fmt.Errorf("%w: %s", ErrNotFound, name)
	}
	out := Entry{Name: e.Name, ArchivedAt: e.ArchivedAt, Secrets: make(map[string]string, len(e.Secrets))}
	for k, v := range e.Secrets {
		out.Secrets[k] = v
	}
	return out, nil
}

// Delete removes a named archive entry.
func (a *Archive) Delete(name string) error {
	if name == "" {
		return ErrEmptyName
	}
	a.mu.Lock()
	defer a.mu.Unlock()
	if _, ok := a.entries[name]; !ok {
		return fmt.Errorf("%w: %s", ErrNotFound, name)
	}
	delete(a.entries, name)
	return nil
}

// All returns a copy of all archive entries.
func (a *Archive) All() []Entry {
	a.mu.RLock()
	defer a.mu.RUnlock()
	out := make([]Entry, 0, len(a.entries))
	for _, e := range a.entries {
		copy := Entry{Name: e.Name, ArchivedAt: e.ArchivedAt, Secrets: make(map[string]string, len(e.Secrets))}
		for k, v := range e.Secrets {
			copy.Secrets[k] = v
		}
		out = append(out, copy)
	}
	return out
}
