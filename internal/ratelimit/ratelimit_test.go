package ratelimit_test

import (
	"testing"
	"time"

	"github.com/yourusername/vaultpipe/internal/ratelimit"
)

func TestAllow_WithinLimit_ReturnsNil(t *testing.T) {
	l := ratelimit.New(3, time.Second)
	for i := 0; i < 3; i++ {
		if err := l.Allow(); err != nil {
			t.Fatalf("call %d: unexpected error: %v", i+1, err)
		}
	}
}

func TestAllow_ExceedsLimit_ReturnsError(t *testing.T) {
	l := ratelimit.New(2, time.Second)
	_ = l.Allow()
	_ = l.Allow()
	if err := l.Allow(); err == nil {
		t.Fatal("expected ErrRateLimited, got nil")
	}
}

func TestAllow_WindowExpiry_AllowsAgain(t *testing.T) {
	now := time.Now()
	calls := 0

	l := ratelimit.New(2, 500*time.Millisecond)
	// Inject controlled clock via unexported field workaround: use real time + sleep.
	_ = l

	// Simpler: use a real short window and sleep.
	l2 := ratelimit.New(1, 100*time.Millisecond)
	if err := l2.Allow(); err != nil {
		t.Fatalf("first call should succeed: %v", err)
	}
	if err := l2.Allow(); err == nil {
		t.Fatal("second call should be rate limited")
	}
	time.Sleep(110 * time.Millisecond)
	if err := l2.Allow(); err != nil {
		t.Fatalf("call after window expiry should succeed: %v", err)
	}
	_ = calls
	_ = now
}

func TestReset_ClearsTimestamps(t *testing.T) {
	l := ratelimit.New(2, time.Second)
	_ = l.Allow()
	_ = l.Allow()
	if err := l.Allow(); err == nil {
		t.Fatal("expected rate limit before reset")
	}
	l.Reset()
	if err := l.Allow(); err != nil {
		t.Fatalf("expected success after reset: %v", err)
	}
}

func TestRemaining_DecreasesWithCalls(t *testing.T) {
	l := ratelimit.New(3, time.Second)
	if got := l.Remaining(); got != 3 {
		t.Fatalf("expected 3 remaining, got %d", got)
	}
	_ = l.Allow()
	if got := l.Remaining(); got != 2 {
		t.Fatalf("expected 2 remaining, got %d", got)
	}
	_ = l.Allow()
	_ = l.Allow()
	if got := l.Remaining(); got != 0 {
		t.Fatalf("expected 0 remaining, got %d", got)
	}
}

func TestNew_PanicsOnInvalidArgs(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for maxRequests < 1")
		}
	}()
	ratelimit.New(0, time.Second)
}
