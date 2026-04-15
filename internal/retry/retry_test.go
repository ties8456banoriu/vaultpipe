package retry

import (
	"errors"
	"testing"
	"time"
)

var errTransient = errors.New("transient error")

func newTestDoer(cfg Config) *Doer {
	d := New(cfg)
	d.sleep = func(time.Duration) {} // no-op sleep for tests
	return d
}

func TestDo_SucceedsOnFirstAttempt(t *testing.T) {
	d := newTestDoer(DefaultConfig())
	calls := 0
	err := d.Do(func() error {
		calls++
		return nil
	})
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if calls != 1 {
		t.Fatalf("expected 1 call, got %d", calls)
	}
}

func TestDo_RetriesAndSucceeds(t *testing.T) {
	d := newTestDoer(DefaultConfig())
	calls := 0
	err := d.Do(func() error {
		calls++
		if calls < 3 {
			return errTransient
		}
		return nil
	})
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if calls != 3 {
		t.Fatalf("expected 3 calls, got %d", calls)
	}
}

func TestDo_ExhaustsAttempts(t *testing.T) {
	cfg := Config{MaxAttempts: 3, Delay: time.Millisecond, Multiplier: 1.0}
	d := newTestDoer(cfg)
	calls := 0
	err := d.Do(func() error {
		calls++
		return errTransient
	})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, errTransient) {
		t.Fatalf("expected errTransient, got %v", err)
	}
	if calls != 3 {
		t.Fatalf("expected 3 calls, got %d", calls)
	}
}

func TestDo_AppliesDelay(t *testing.T) {
	cfg := Config{MaxAttempts: 3, Delay: 100 * time.Millisecond, Multiplier: 2.0}
	d := New(cfg)

	var delays []time.Duration
	d.sleep = func(dur time.Duration) { delays = append(delays, dur) }

	_ = d.Do(func() error { return errTransient })

	if len(delays) != 2 {
		t.Fatalf("expected 2 sleep calls, got %d", len(delays))
	}
	if delays[0] != 100*time.Millisecond {
		t.Errorf("expected first delay 100ms, got %v", delays[0])
	}
	if delays[1] != 200*time.Millisecond {
		t.Errorf("expected second delay 200ms, got %v", delays[1])
	}
}

func TestDo_SingleAttemptNoDelay(t *testing.T) {
	cfg := Config{MaxAttempts: 1, Delay: time.Second, Multiplier: 1.0}
	d := New(cfg)
	slept := false
	d.sleep = func(time.Duration) { slept = true }

	_ = d.Do(func() error { return errTransient })

	if slept {
		t.Error("expected no sleep on single attempt")
	}
}
