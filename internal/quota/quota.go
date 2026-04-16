// Package quota enforces per-key secret fetch quotas within a time window.
package quota

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

// ErrQuotaExceeded is returned when a key has exceeded its allowed fetches.
var ErrQuotaExceeded = errors.New("quota: fetch limit exceeded for key")

// ErrNoQuota is returned when no quota is configured for a key.
var ErrNoQuota = errors.New("quota: no quota configured for key")

type entry struct {
	count     int
	windowEnd time.Time
}

// Enforcer tracks and enforces fetch quotas per secret key.
type Enforcer struct {
	mu      sync.Mutex
	limits  map[string]int
	window  time.Duration
	entries map[string]*entry
}

// New creates an Enforcer with the given window duration.
func New(window time.Duration) (*Enforcer, error) {
	if window <= 0 {
		return nil, errors.New("quota: window must be positive")
	}
	return &Enforcer{
		limits:  make(map[string]int),
		window:  window,
		entries: make(map[string]*entry),
	}, nil
}

// SetLimit configures a maximum fetch count for the given key within the window.
func (e *Enforcer) SetLimit(key string, max int) error {
	if key == "" {
		return errors.New("quota: key must not be empty")
	}
	if max <= 0 {
		return errors.New("quota: max must be positive")
	}
	e.mu.Lock()
	defer e.mu.Unlock()
	e.limits[key] = max
	return nil
}

// Check records a fetch attempt for key and returns ErrQuotaExceeded if the
// limit has been reached within the current window.
func (e *Enforcer) Check(key string) error {
	if key == "" {
		return errors.New("quota: key must not be empty")
	}
	e.mu.Lock()
	defer e.mu.Unlock()
	max, ok := e.limits[key]
	if !ok {
		return fmt.Errorf("%w: %s", ErrNoQuota, key)
	}
	now := time.Now()
	ent, exists := e.entries[key]
	if !exists || now.After(ent.windowEnd) {
		e.entries[key] = &entry{count: 1, windowEnd: now.Add(e.window)}
		return nil
	}
	if ent.count >= max {
		return fmt.Errorf("%w: %s", ErrQuotaExceeded, key)
	}
	ent.count++
	return nil
}

// Remaining returns how many fetches are left for key in the current window.
func (e *Enforcer) Remaining(key string) (int, error) {
	e.mu.Lock()
	defer e.mu.Unlock()
	max, ok := e.limits[key]
	if !ok {
		return 0, fmt.Errorf("%w: %s", ErrNoQuota, key)
	}
	ent, exists := e.entries[key]
	if !exists || time.Now().After(ent.windowEnd) {
		return max, nil
	}
	return max - ent.count, nil
}
