package secretwatch_test

import (
	"testing"

	"github.com/yourusername/vaultpipe/internal/secretwatch"
)

func TestRecord_And_All_RoundTrip(t *testing.T) {
	w := secretwatch.New()
	if err := w.Record("DB_PASS", "secret/data/db"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	events := w.All()
	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}
	if events[0].EnvKey != "DB_PASS" || events[0].VaultPath != "secret/data/db" {
		t.Errorf("unexpected event: %+v", events[0])
	}
}

func TestRecord_EmptyKey_ReturnsError(t *testing.T) {
	w := secretwatch.New()
	if err := w.Record("", "secret/data/db"); err == nil {
		t.Fatal("expected error for empty envKey")
	}
}

func TestRecord_EmptyVaultPath_ReturnsError(t *testing.T) {
	w := secretwatch.New()
	if err := w.Record("DB_PASS", ""); err == nil {
		t.Fatal("expected error for empty vaultPath")
	}
}

func TestRecord_DuplicateKey_UpdatesEvent(t *testing.T) {
	w := secretwatch.New()
	_ = w.Record("DB_PASS", "secret/data/db")
	_ = w.Record("DB_PASS", "secret/data/db-v2")
	events := w.All()
	if len(events) != 1 {
		t.Fatalf("expected 1 event after duplicate, got %d", len(events))
	}
	if events[0].VaultPath != "secret/data/db-v2" {
		t.Errorf("expected updated vault path, got %s", events[0].VaultPath)
	}
}

func TestHas_ReturnsTrueAfterRecord(t *testing.T) {
	w := secretwatch.New()
	_ = w.Record("API_KEY", "secret/data/api")
	if !w.Has("API_KEY") {
		t.Error("expected Has to return true")
	}
	if w.Has("MISSING") {
		t.Error("expected Has to return false for unknown key")
	}
}

func TestClear_RemovesAllEvents(t *testing.T) {
	w := secretwatch.New()
	_ = w.Record("API_KEY", "secret/data/api")
	w.Clear()
	if len(w.All()) != 0 {
		t.Error("expected empty events after Clear")
	}
	if w.Has("API_KEY") {
		t.Error("expected Has to return false after Clear")
	}
}

func TestAll_ReturnsCopy(t *testing.T) {
	w := secretwatch.New()
	_ = w.Record("X", "secret/data/x")
	a := w.All()
	a[0].EnvKey = "MUTATED"
	if w.All()[0].EnvKey == "MUTATED" {
		t.Error("All should return a copy, not a reference")
	}
}
