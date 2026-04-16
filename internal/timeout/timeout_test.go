package timeout_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/yourusername/vaultpipe/internal/timeout"
)

func TestDo_SucceedsWithinDeadline(t *testing.T) {
	d := timeout.New(500 * time.Millisecond)

	secrets, err := d.Do(context.Background(), func(ctx context.Context) (map[string]string, error) {
		return map[string]string{"KEY": "value"}, nil
	})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if secrets["KEY"] != "value" {
		t.Errorf("expected KEY=value, got %q", secrets["KEY"])
	}
}

func TestDo_ExceedsDeadline_ReturnsErrDeadlineExceeded(t *testing.T) {
	d := timeout.New(50 * time.Millisecond)

	_, err := d.Do(context.Background(), func(ctx context.Context) (map[string]string, error) {
		time.Sleep(200 * time.Millisecond)
		return nil, nil
	})

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, timeout.ErrDeadlineExceeded) {
		t.Errorf("expected ErrDeadlineExceeded, got %v", err)
	}
}

func TestDo_ParentContextCancelled_ReturnsError(t *testing.T) {
	d := timeout.New(5 * time.Second)

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // cancel immediately

	_, err := d.Do(ctx, func(ctx context.Context) (map[string]string, error) {
		time.Sleep(100 * time.Millisecond)
		return nil, nil
	})

	if err == nil {
		t.Fatal("expected error for cancelled parent context, got nil")
	}
}

func TestDo_FetchReturnsError_PropagatesError(t *testing.T) {
	d := timeout.New(500 * time.Millisecond)
	wantErr := errors.New("vault unavailable")

	_, err := d.Do(context.Background(), func(ctx context.Context) (map[string]string, error) {
		return nil, wantErr
	})

	if !errors.Is(err, wantErr) {
		t.Errorf("expected %v, got %v", wantErr, err)
	}
}

func TestNew_ZeroDuration_UsesDefault(t *testing.T) {
	// A zero duration should not panic and should still allow fast operations.
	d := timeout.New(0)

	secrets, err := d.Do(context.Background(), func(ctx context.Context) (map[string]string, error) {
		return map[string]string{"A": "1"}, nil
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if secrets["A"] != "1" {
		t.Errorf("expected A=1, got %q", secrets["A"])
	}
}

func TestDo_FetchReturnsErrorAfterDeadline_ReturnsErrDeadlineExceeded(t *testing.T) {
	// When both the deadline is exceeded and the fetch returns an error,
	// ErrDeadlineExceeded should take precedence.
	d := timeout.New(50 * time.Millisecond)
	fetchErr := errors.New("fetch failed")

	_, err := d.Do(context.Background(), func(ctx context.Context) (map[string]string, error) {
		time.Sleep(200 * time.Millisecond)
		return nil, fetchErr
	})

	if !errors.Is(err, timeout.ErrDeadlineExceeded) {
		t.Errorf("expected ErrDeadlineExceeded, got %v", err)
	}
}
