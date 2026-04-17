package secretlog_test

import (
	"testing"
	"time"

	"github.com/vincentfree/vaultpipe/internal/secretlog"
)

func TestRecord_EmptyKey_ReturnsError(t *testing.T) {
	l := secretlog.New()
	err := l.Record("", "secret/data/app", "", time.Time{})
	if err == nil {
		t.Fatal("expected error for empty key")
	}
}

func TestRecord_EmptyVaultPath_ReturnsError(t *testing.T) {
	l := secretlog.New()
	err := l.Record("DB_PASS", "", "", time.Time{})
	if err == nil {
		t.Fatal("expected error for empty vaultPath")
	}
}

func TestRecord_SetsTimestampIfZero(t *testing.T) {
	l := secretlog.New()
	before := time.Now().UTC()
	if err := l.Record("DB_PASS", "secret/data/app", "", time.Time{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	entries := l.All()
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	if entries[0].AccessedAt.Before(before) {
		t.Error("timestamp should be set to at least before-time")
	}
}

func TestAll_ReturnsCopy(t *testing.T) {
	l := secretlog.New()
	_ = l.Record("KEY", "secret/data/app", "dev", time.Now())
	a := l.All()
	a[0].Key = "MUTATED"
	b := l.All()
	if b[0].Key == "MUTATED" {
		t.Error("All should return a copy, not a reference")
	}
}

func TestForKey_FiltersCorrectly(t *testing.T) {
	l := secretlog.New()
	_ = l.Record("API_KEY", "secret/data/app", "", time.Now())
	_ = l.Record("DB_PASS", "secret/data/db", "", time.Now())
	_ = l.Record("API_KEY", "secret/data/app", "prod", time.Now())

	results := l.ForKey("API_KEY")
	if len(results) != 2 {
		t.Fatalf("expected 2 entries for API_KEY, got %d", len(results))
	}
	for _, e := range results {
		if e.Key != "API_KEY" {
			t.Errorf("unexpected key %q in ForKey result", e.Key)
		}
	}
}

func TestClear_RemovesEntries(t *testing.T) {
	l := secretlog.New()
	_ = l.Record("KEY", "secret/data/app", "", time.Now())
	l.Clear()
	if len(l.All()) != 0 {
		t.Error("expected no entries after Clear")
	}
}

func TestRecord_StoresProfile(t *testing.T) {
	l := secretlog.New()
	_ = l.Record("TOKEN", "secret/data/svc", "staging", time.Now())
	entries := l.All()
	if entries[0].Profile != "staging" {
		t.Errorf("expected profile 'staging', got %q", entries[0].Profile)
	}
}
