// Package secretlog provides a structured access log for secret reads,
// recording which keys were accessed, when, and by which profile.
package secretlog

import (
	"errors"
	"sync"
	"time"
)

// Entry represents a single secret access event.
type Entry struct {
	Key       string    `json:"key"`
	VaultPath string    `json:"vault_path"`
	Profile   string    `json:"profile,omitempty"`
	AccessedAt time.Time `json:"accessed_at"`
}

// Logger records secret access entries in memory.
type Logger struct {
	mu      sync.Mutex
	entries []Entry
}

// New returns a new Logger.
func New() *Logger {
	return &Logger{}
}

// Record appends an access entry. Returns an error if key or vaultPath is empty.
func (l *Logger) Record(key, vaultPath, profile string, at time.Time) error {
	if key == "" {
		return errors.New("secretlog: key must not be empty")
	}
	if vaultPath == "" {
		return errors.New("secretlog: vaultPath must not be empty")
	}
	if at.IsZero() {
		at = time.Now().UTC()
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	l.entries = append(l.entries, Entry{
		Key:        key,
		VaultPath:  vaultPath,
		Profile:    profile,
		AccessedAt: at,
	})
	return nil
}

// All returns a copy of all recorded entries.
func (l *Logger) All() []Entry {
	l.mu.Lock()
	defer l.mu.Unlock()
	out := make([]Entry, len(l.entries))
	copy(out, l.entries)
	return out
}

// Clear removes all recorded entries.
func (l *Logger) Clear() {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.entries = nil
}

// ForKey returns all entries matching the given key.
func (l *Logger) ForKey(key string) []Entry {
	l.mu.Lock()
	defer l.mu.Unlock()
	var out []Entry
	for _, e := range l.entries {
		if e.Key == key {
			out = append(out, e)
		}
	}
	return out
}
