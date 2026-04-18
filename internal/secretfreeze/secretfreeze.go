// Package secretfreeze provides the ability to freeze a set of secrets,
// preventing any further modifications until explicitly thawed.
package secretfreeze

import (
	"errors"
	"fmt"
	"sync"
)

// ErrFrozen is returned when an operation is attempted on a frozen secret.
var ErrFrozen = errors.New("secretfreeze: secret is frozen")

// ErrNotFrozen is returned when thaw is called on a non-frozen secret.
var ErrNotFrozen = errors.New("secretfreeze: secret is not frozen")

// Freezer manages frozen state for secret keys.
type Freezer struct {
	mu     sync.RWMutex
	frozen map[string]struct{}
}

// New returns a new Freezer.
func New() *Freezer {
	return &Freezer{frozen: make(map[string]struct{})}
}

// Freeze marks the given key as frozen.
func (f *Freezer) Freeze(key string) error {
	if key == "" {
		return errors.New("secretfreeze: key must not be empty")
	}
	f.mu.Lock()
	defer f.mu.Unlock()
	f.frozen[key] = struct{}{}
	return nil
}

// Thaw removes the frozen state from the given key.
func (f *Freezer) Thaw(key string) error {
	if key == "" {
		return errors.New("secretfreeze: key must not be empty")
	}
	f.mu.Lock()
	defer f.mu.Unlock()
	if _, ok := f.frozen[key]; !ok {
		return ErrNotFrozen
	}
	delete(f.frozen, key)
	return nil
}

// IsFrozen reports whether the given key is frozen.
func (f *Freezer) IsFrozen(key string) bool {
	f.mu.RLock()
	defer f.mu.RUnlock()
	_, ok := f.frozen[key]
	return ok
}

// Apply returns a copy of secrets with frozen keys preserved from base.
// If a frozen key is missing from base, an error is returned.
func (f *Freezer) Apply(base, incoming map[string]string) (map[string]string, error) {
	if len(incoming) == 0 {
		return nil, errors.New("secretfreeze: incoming secrets must not be empty")
	}
	f.mu.RLock()
	defer f.mu.RUnlock()
	out := make(map[string]string, len(incoming))
	for k, v := range incoming {
		out[k] = v
	}
	for k := range f.frozen {
		v, ok := base[k]
		if !ok {
			return nil, fmt.Errorf("%w: key %q not found in base", ErrFrozen, k)
		}
		out[k] = v
	}
	return out, nil
}
