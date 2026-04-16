// Package expiry provides TTL-based expiry tracking for Vault secrets,
// warning when secrets are approaching or past their lease duration.
package expiry

import (
	"errors"
	"fmt"
	"time"
)

// ErrNoExpiry is returned when no expiry has been set for a key.
var ErrNoExpiry = errors.New("expiry: no expiry set for key")

// Status represents the expiry state of a secret.
type Status string

const (
	StatusOK      Status = "ok"
	StatusWarning Status = "warning"
	StatusExpired Status = "expired"
)

// Entry holds the expiry metadata for a single secret key.
type Entry struct {
	Key       string
	ExpiresAt time.Time
	SetAt     time.Time
}

// Tracker manages expiry entries for secret keys.
type Tracker struct {
	entries     map[string]Entry
	warnBefore  time.Duration
	now         func() time.Time
}

// New creates a Tracker. warnBefore is the duration before expiry at which
// Status transitions from OK to Warning.
func New(warnBefore time.Duration) (*Tracker, error) {
	if warnBefore <= 0 {
		return nil, errors.New("expiry: warnBefore must be positive")
	}
	return &Tracker{
		entries:    make(map[string]Entry),
		warnBefore: warnBefore,
		now:        time.Now,
	}, nil
}

// Track registers a key with the given TTL starting from now.
func (t *Tracker) Track(key string, ttl time.Duration) error {
	if key == "" {
		return errors.New("expiry: key must not be empty")
	}
	if ttl <= 0 {
		return errors.New("expiry: ttl must be positive")
	}
	now := t.now()
	t.entries[key] = Entry{
		Key:       key,
		ExpiresAt: now.Add(ttl),
		SetAt:     now,
	}
	return nil
}

// Check returns the Status and remaining duration for a key.
func (t *Tracker) Check(key string) (Status, time.Duration, error) {
	e, ok := t.entries[key]
	if !ok {
		return "", 0, fmt.Errorf("%w: %s", ErrNoExpiry, key)
	}
	remaining := e.ExpiresAt.Sub(t.now())
	if remaining <= 0 {
		return StatusExpired, 0, nil
	}
	if remaining <= t.warnBefore {
		return StatusWarning, remaining, nil
	}
	return StatusOK, remaining, nil
}

// Remove deletes the expiry entry for a key.
func (t *Tracker) Remove(key string) {
	delete(t.entries, key)
}
