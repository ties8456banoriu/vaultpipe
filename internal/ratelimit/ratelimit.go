// Package ratelimit provides a simple token-bucket rate limiter for
// controlling how frequently secrets are fetched from Vault, preventing
// accidental API abuse during rapid refresh cycles or misconfigured watchers.
package ratelimit

import (
	"errors"
	"sync"
	"time"
)

// ErrRateLimited is returned when a request exceeds the allowed rate.
var ErrRateLimited = errors.New("rate limit exceeded: too many requests in window")

// Limiter enforces a maximum number of allowed calls within a sliding window.
type Limiter struct {
	mu       sync.Mutex
	max      int
	window   time.Duration
	timestamps []time.Time
	now      func() time.Time
}

// New creates a Limiter that allows at most maxRequests calls per window duration.
// maxRequests must be >= 1 and window must be > 0, otherwise New panics.
func New(maxRequests int, window time.Duration) *Limiter {
	if maxRequests < 1 {
		panic("ratelimit: maxRequests must be >= 1")
	}
	if window <= 0 {
		panic("ratelimit: window must be > 0")
	}
	return &Limiter{
		max:    maxRequests,
		window: window,
		now:    time.Now,
	}
}

// Allow reports whether the current call is within the rate limit.
// It returns ErrRateLimited if the limit has been exceeded.
func (l *Limiter) Allow() error {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := l.now()
	cutoff := now.Add(-l.window)

	// Evict timestamps outside the current window.
	valid := l.timestamps[:0]
	for _, t := range l.timestamps {
		if t.After(cutoff) {
			valid = append(valid, t)
		}
	}
	l.timestamps = valid

	if len(l.timestamps) >= l.max {
		return ErrRateLimited
	}

	l.timestamps = append(l.timestamps, now)
	return nil
}

// Reset clears all recorded timestamps, fully restoring the limiter's budget.
func (l *Limiter) Reset() {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.timestamps = l.timestamps[:0]
}

// Remaining returns the number of calls still allowed in the current window.
func (l *Limiter) Remaining() int {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := l.now()
	cutoff := now.Add(-l.window)
	count := 0
	for _, t := range l.timestamps {
		if t.After(cutoff) {
			count++
		}
	}
	rem := l.max - count
	if rem < 0 {
		return 0
	}
	return rem
}
