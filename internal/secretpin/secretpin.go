// Package secretpin allows pinning secrets to specific Vault versions,
// preventing unintended updates during refresh cycles.
package secretpin

import (
	"errors"
	"sync"
	"time"
)

// ErrNotPinned is returned when a key has no pin recorded.
var ErrNotPinned = errors.New("secretpin: key is not pinned")

// Pin holds pinning metadata for a single secret.
type Pin struct {
	EnvKey  string
	Version int
	PinnedAt time.Time
}

// Pinner stores version pins for secret keys.
type Pinner struct {
	mu   sync.RWMutex
	pins map[string]Pin
}

// New returns an initialised Pinner.
func New() *Pinner {
	return &Pinner{pins: make(map[string]Pin)}
}

// Pin records a version pin for the given env key.
func (p *Pinner) Pin(envKey string, version int) error {
	if envKey == "" {
		return errors.New("secretpin: envKey must not be empty")
	}
	if version < 1 {
		return errors.New("secretpin: version must be >= 1")
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	p.pins[envKey] = Pin{EnvKey: envKey, Version: version, PinnedAt: time.Now().UTC()}
	return nil
}

// Unpin removes the pin for the given env key.
func (p *Pinner) Unpin(envKey string) error {
	if envKey == "" {
		return errors.New("secretpin: envKey must not be empty")
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	delete(p.pins, envKey)
	return nil
}

// Get returns the Pin for the given key or ErrNotPinned.
func (p *Pinner) Get(envKey string) (Pin, error) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	pin, ok := p.pins[envKey]
	if !ok {
		return Pin{}, ErrNotPinned
	}
	return pin, nil
}

// IsPinned reports whether a key is currently pinned.
func (p *Pinner) IsPinned(envKey string) bool {
	p.mu.RLock()
	defer p.mu.RUnlock()
	_, ok := p.pins[envKey]
	return ok
}

// All returns a copy of all current pins.
func (p *Pinner) All() []Pin {
	p.mu.RLock()
	defer p.mu.RUnlock()
	out := make([]Pin, 0, len(p.pins))
	for _, pin := range p.pins {
		out = append(out, pin)
	}
	return out
}
