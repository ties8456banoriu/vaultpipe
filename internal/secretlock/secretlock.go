// Package secretlock provides a mechanism to lock specific secret keys,
// preventing them from being overwritten during a merge or refresh cycle.
package secretlock

import (
	"errors"
	"sync"
)

// ErrKeyLocked is returned when an attempt is made to modify a locked key.
var ErrKeyLocked = errors.New("secretlock: key is locked")

// ErrEmptyKey is returned when an empty key is provided.
var ErrEmptyKey = errors.New("secretlock: key must not be empty")

// Locker manages a set of locked secret keys.
type Locker struct {
	mu      sync.RWMutex
	locked  map[string]struct{}
}

// New returns a new Locker.
func New() *Locker {
	return &Locker{locked: make(map[string]struct{})}
}

// Lock marks the given key as locked.
func (l *Locker) Lock(key string) error {
	if key == "" {
		return ErrEmptyKey
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	l.locked[key] = struct{}{}
	return nil
}

// Unlock removes the lock from the given key.
func (l *Locker) Unlock(key string) error {
	if key == "" {
		return ErrEmptyKey
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	delete(l.locked, key)
	return nil
}

// IsLocked reports whether the given key is locked.
func (l *Locker) IsLocked(key string) bool {
	l.mu.RLock()
	defer l.mu.RUnlock()
	_, ok := l.locked[key]
	return ok
}

// Apply filters out locked keys from the incoming secrets map, returning
// only the secrets that are not locked. Returns an error if secrets is empty.
func (l *Locker) Apply(secrets map[string]string) (map[string]string, error) {
	if len(secrets) == 0 {
		return nil, errors.New("secretlock: secrets must not be empty")
	}
	l.mu.RLock()
	defer l.mu.RUnlock()
	out := make(map[string]string, len(secrets))
	for k, v := range secrets {
		if _, locked := l.locked[k]; !locked {
			out[k] = v
		}
	}
	return out, nil
}

// All returns a copy of all currently locked keys.
func (l *Locker) All() []string {
	l.mu.RLock()
	defer l.mu.RUnlock()
	keys := make([]string, 0, len(l.locked))
	for k := range l.locked {
		keys = append(keys, k)
	}
	return keys
}
