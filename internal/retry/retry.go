package retry

import (
	"errors"
	"time"
)

// ErrMaxAttemptsReached is returned when all retry attempts are exhausted.
var ErrMaxAttemptsReached = errors.New("retry: max attempts reached")

// Config holds the retry policy configuration.
type Config struct {
	MaxAttempts int
	Delay       time.Duration
	Multiplier  float64 // exponential backoff multiplier; 1.0 = constant delay
}

// DefaultConfig returns a sensible default retry configuration.
func DefaultConfig() Config {
	return Config{
		MaxAttempts: 3,
		Delay:       500 * time.Millisecond,
		Multiplier:  2.0,
	}
}

// Doer executes a retry loop with the given config.
type Doer struct {
	cfg   Config
	sleep func(time.Duration)
}

// New creates a new Doer with the provided config.
func New(cfg Config) *Doer {
	return &Doer{
		cfg:   cfg,
		sleep: time.Sleep,
	}
}

// Do calls fn up to MaxAttempts times. It stops early if fn returns nil.
// On each failure the delay is multiplied by Multiplier before the next attempt.
func (d *Doer) Do(fn func() error) error {
	delay := d.cfg.Delay
	var lastErr error

	for attempt := 1; attempt <= d.cfg.MaxAttempts; attempt++ {
		if err := fn(); err == nil {
			return nil
		} else {
			lastErr = err
		}

		if attempt < d.cfg.MaxAttempts {
			d.sleep(delay)
			delay = time.Duration(float64(delay) * d.cfg.Multiplier)
		}
	}

	if lastErr != nil {
		return lastErr
	}
	return ErrMaxAttemptsReached
}
