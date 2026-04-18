package secretarchive_test

import (
	"errors"
	"testing"

	"github.com/yourusername/vaultpipe/internal/secretarchive"
)

var base = map[string]string{"DB_PASS": "secret", "API_KEY": "abc123"}

func TestStore_And_Get_RoundTrip(t *testing.T) {
	a := secretarchive.New()
	if err := a.Store("v1", base); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	e, err := a.Get("v1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if e.Secrets["DB_PASS"] != "secret" {
		t.Errorf("expected secret, got %s", e.Secrets["DB_PASS"])
	}
}

func TestStore_EmptyName_ReturnsError(t *testing.T) {
	a := secretarchive.New()
	if err := a.Store("", base); !errors.Is(err, secretarchive.ErrEmptyName) {
		t.Errorf("expected ErrEmptyName, got %v", err)
	}
}

func TestStore_EmptySecrets_ReturnsError(t *testing.T) {
	a := secretarchive.New()
	if err := a.Store("v1", map[string]string{}); !errors.Is(err, secretarchive.ErrEmptySecrets) {
		t.Errorf("expected ErrEmptySecrets, got %v", err)
	}
}

func TestGet_NotFound_ReturnsError(t *testing.T) {
	a := secretarchive.New()
	_, err := a.Get("missing")
	if !errors.Is(err, secretarchive.ErrNotFound) {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

func TestGet_ReturnsCopy(t *testing.T) {
	a := secretarchive.New()
	_ = a.Store("v1", base)
	e, _ := a.Get("v1")
	e.Secrets["DB_PASS"] = "mutated"
	e2, _ := a.Get("v1")
	if e2.Secrets["DB_PASS"] == "mutated" {
		t.Error("expected isolation from mutation")
	}
}

func TestDelete_RemovesEntry(t *testing.T) {
	a := secretarchive.New()
	_ = a.Store("v1", base)
	if err := a.Delete("v1"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, err := a.Get("v1"); !errors.Is(err, secretarchive.ErrNotFound) {
		t.Error("expected ErrNotFound after delete")
	}
}

func TestDelete_NotFound_ReturnsError(t *testing.T) {
	a := secretarchive.New()
	if err := a.Delete("ghost"); !errors.Is(err, secretarchive.ErrNotFound) {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

func TestAll_ReturnsAllEntries(t *testing.T) {
	a := secretarchive.New()
	_ = a.Store("v1", base)
	_ = a.Store("v2", map[string]string{"X": "y"})
	if got := len(a.All()); got != 2 {
		t.Errorf("expected 2 entries, got %d", got)
	}
}

func TestStore_SetsArchivedAt(t *testing.T) {
	a := secretarchive.New()
	_ = a.Store("v1", base)
	e, _ := a.Get("v1")
	if e.ArchivedAt.IsZero() {
		t.Error("expected ArchivedAt to be set")
	}
}
