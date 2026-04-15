package refresh_test

import (
	"context"
	"errors"
	"log"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/yourusername/vaultpipe/internal/refresh"
)

// --- fakes ---

type fakeFetcher struct {
	mu      sync.Mutex
	calls   int
	secrets map[string]string
	err     error
}

func (f *fakeFetcher) GetSecret(_ string) (map[string]string, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.calls++
	return f.secrets, f.err
}

type fakeWriter struct {
	mu     sync.Mutex
	writes []map[string]string
	err    error
}

func (fw *fakeWriter) Write(secrets map[string]string) error {
	fw.mu.Lock()
	defer fw.mu.Unlock()
	fw.writes = append(fw.writes, secrets)
	return fw.err
}

func testLogger() *log.Logger {
	return log.New(os.Stderr, "", 0)
}

// --- tests ---

func TestWatcher_PerformsInitialRefresh(t *testing.T) {
	fetcher := &fakeFetcher{secrets: map[string]string{"KEY": "val"}}
	writer := &fakeWriter{}

	w := refresh.NewWatcher(fetcher, writer, "secret/app", 500*time.Millisecond, testLogger())

	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	_ = w.Start(ctx)

	writer.mu.Lock()
	defer writer.mu.Unlock()
	if len(writer.writes) == 0 {
		t.Fatal("expected at least one write from initial refresh")
	}
}

func TestWatcher_TickerTriggersRefresh(t *testing.T) {
	fetcher := &fakeFetcher{secrets: map[string]string{"TOKEN": "abc"}}
	writer := &fakeWriter{}

	w := refresh.NewWatcher(fetcher, writer, "secret/app", 50*time.Millisecond, testLogger())

	ctx, cancel := context.WithTimeout(context.Background(), 180*time.Millisecond)
	defer cancel()

	_ = w.Start(ctx)

	writer.mu.Lock()
	defer writer.mu.Unlock()
	if len(writer.writes) < 2 {
		t.Fatalf("expected multiple writes, got %d", len(writer.writes))
	}
}

func TestWatcher_FetchErrorDoesNotStop(t *testing.T) {
	fetcher := &fakeFetcher{err: errors.New("vault unavailable")}
	writer := &fakeWriter{}

	w := refresh.NewWatcher(fetcher, writer, "secret/app", 40*time.Millisecond, testLogger())

	ctx, cancel := context.WithTimeout(context.Background(), 150*time.Millisecond)
	defer cancel()

	// Should return ctx.Err(), not the fetch error
	err := w.Start(ctx)
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("expected DeadlineExceeded, got %v", err)
	}

	fetcher.mu.Lock()
	defer fetcher.mu.Unlock()
	if fetcher.calls == 0 {
		t.Fatal("expected fetcher to be called despite errors")
	}
}
