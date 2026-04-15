// Package timeout provides a configurable deadline wrapper for Vault fetch
// operations. It wraps any fetch function and enforces a maximum execution
// duration, returning a descriptive error if the deadline is exceeded.
package timeout

import (
	"context"
	"errors"
	"fmt"
	"time"
)

// ErrDeadlineExceeded is returned when the wrapped operation exceeds its
// allowed duration.
var ErrDeadlineExceeded = errors.New("operation exceeded deadline")

// FetchFunc is any function that accepts a context and returns a map of
// secrets or an error — matching the signature used by the Vault client.
type FetchFunc func(ctx context.Context) (map[string]string, error)

// Doer wraps a FetchFunc with a configurable timeout.
type Doer struct {
	timeout time.Duration
}

// New creates a Doer that will cancel any fetch operation after d.
// If d is zero or negative, a default of 10 seconds is used.
func New(d time.Duration) *Doer {
	if d <= 0 {
		d = 10 * time.Second
	}
	return &Doer{timeout: d}
}

// Do executes fn within the configured deadline. If the context supplied by
// the caller is already cancelled, that error is returned immediately. If the
// operation exceeds the deadline, ErrDeadlineExceeded is returned wrapping the
// underlying context error.
func (d *Doer) Do(parent context.Context, fn FetchFunc) (map[string]string, error) {
	ctx, cancel := context.WithTimeout(parent, d.timeout)
	defer cancel()

	type result struct {
		secrets map[string]string
		err     error
	}

	ch := make(chan result, 1)
	go func() {
		s, err := fn(ctx)
		ch <- result{s, err}
	}()

	select {
	case res := <-ch:
		return res.secrets, res.err
	case <-ctx.Done():
		return nil, fmt.Errorf("%w: %w", ErrDeadlineExceeded, ctx.Err())
	}
}
