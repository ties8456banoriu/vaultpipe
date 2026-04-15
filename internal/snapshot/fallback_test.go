package snapshot_test

import (
	"errors"
	"testing"

	"github.com/yourusername/vaultpipe/internal/snapshot"
)

func TestWithFallback_SuccessfulFetch(t *testing.T) {
	expected := map[string]string{"KEY": "live"}
	fetchFn := func() (map[string]string, error) { return expected, nil }

	got, err := snapshot.WithFallback(fetchFn, "/unused/path")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["KEY"] != "live" {
		t.Errorf("expected live value, got %s", got["KEY"])
	}
}

func TestWithFallback_UsesSnapshotOnFetchError(t *testing.T) {
	path := tempPath(t)
	if err := snapshot.Store(path, map[string]string{"KEY": "cached"}); err != nil {
		t.Fatal(err)
	}

	fetchFn := func() (map[string]string, error) {
		return nil, errors.New("vault unreachable")
	}

	got, err := snapshot.WithFallback(fetchFn, path)
	if !errors.Is(err, snapshot.ErrFallbackUsed) {
		t.Errorf("expected ErrFallbackUsed, got %v", err)
	}
	if got["KEY"] != "cached" {
		t.Errorf("expected cached value, got %s", got["KEY"])
	}
}

func TestWithFallback_BothFail_ReturnsOriginalError(t *testing.T) {
	fetchFn := func() (map[string]string, error) {
		return nil, errors.New("vault down")
	}

	_, err := snapshot.WithFallback(fetchFn, "/no/snapshot/here")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if errors.Is(err, snapshot.ErrFallbackUsed) {
		t.Error("should not return ErrFallbackUsed when snapshot is also missing")
	}
}
