package secretversion_test

import (
	"testing"
	"time"

	"github.com/damon/vaultpipe/internal/secretversion"
)

func TestRecord_And_Latest_RoundTrip(t *testing.T) {
	tr := secretversion.New()
	now := time.Now().UTC()
	if err := tr.Record("DB_PASS", "secret/db", 3, now); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	e, err := tr.Latest("DB_PASS")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if e.Version != 3 || e.VaultPath != "secret/db" {
		t.Errorf("unexpected entry: %+v", e)
	}
}

func TestRecord_SetsTimestampIfZero(t *testing.T) {
	tr := secretversion.New()
	if err := tr.Record("API_KEY", "secret/api", 1, time.Time{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	e, _ := tr.Latest("API_KEY")
	if e.FetchedAt.IsZero() {
		t.Error("expected FetchedAt to be set")
	}
}

func TestRecord_EmptyEnvKey_ReturnsError(t *testing.T) {
	tr := secretversion.New()
	if err := tr.Record("", "secret/x", 1, time.Now()); err == nil {
		t.Error("expected error for empty envKey")
	}
}

func TestRecord_EmptyVaultPath_ReturnsError(t *testing.T) {
	tr := secretversion.New()
	if err := tr.Record("KEY", "", 1, time.Now()); err == nil {
		t.Error("expected error for empty vaultPath")
	}
}

func TestRecord_NegativeVersion_ReturnsError(t *testing.T) {
	tr := secretversion.New()
	if err := tr.Record("KEY", "secret/x", -1, time.Now()); err == nil {
		t.Error("expected error for negative version")
	}
}

func TestLatest_UnknownKey_ReturnsError(t *testing.T) {
	tr := secretversion.New()
	if _, err := tr.Latest("MISSING"); err == nil {
		t.Error("expected error for unknown key")
	}
}

func TestAll_ReturnsCopy(t *testing.T) {
	tr := secretversion.New()
	now := time.Now().UTC()
	_ = tr.Record("TOKEN", "secret/token", 1, now)
	_ = tr.Record("TOKEN", "secret/token", 2, now)
	entries, err := tr.All("TOKEN")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 2 {
		t.Errorf("expected 2 entries, got %d", len(entries))
	}
	entries[0].Version = 99
	latest, _ := tr.Latest("TOKEN")
	if latest.Version == 99 {
		t.Error("mutation of copy affected internal state")
	}
}

func TestAll_UnknownKey_ReturnsError(t *testing.T) {
	tr := secretversion.New()
	if _, err := tr.All("NOPE"); err == nil {
		t.Error("expected error for unknown key")
	}
}
