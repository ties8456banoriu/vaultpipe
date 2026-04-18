package secretdiff_test

import (
	"testing"
	"time"

	"github.com/yourusername/vaultpipe/internal/secretdiff"
)

func TestRecord_And_Get_RoundTrip(t *testing.T) {
	tr := secretdiff.New()
	e := secretdiff.Event{
		EnvKey:    "DB_PASS",
		VaultPath: "secret/db",
		OldValue:  "old",
		NewValue:  "new",
	}
	if err := tr.Record(e); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got, err := tr.Get("DB_PASS")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.NewValue != "new" || got.OldValue != "old" {
		t.Errorf("unexpected event: %+v", got)
	}
}

func TestRecord_SetsTimestampIfZero(t *testing.T) {
	tr := secretdiff.New()
	before := time.Now().UTC()
	_ = tr.Record(secretdiff.Event{EnvKey: "K", VaultPath: "p"})
	got, _ := tr.Get("K")
	if got.ChangedAt.Before(before) {
		t.Errorf("expected timestamp >= %v, got %v", before, got.ChangedAt)
	}
}

func TestRecord_EmptyEnvKey_ReturnsError(t *testing.T) {
	tr := secretdiff.New()
	err := tr.Record(secretdiff.Event{VaultPath: "p"})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestRecord_EmptyVaultPath_ReturnsError(t *testing.T) {
	tr := secretdiff.New()
	err := tr.Record(secretdiff.Event{EnvKey: "K"})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestGet_UnknownKey_ReturnsError(t *testing.T) {
	tr := secretdiff.New()
	_, err := tr.Get("MISSING")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestAll_ReturnsCopy(t *testing.T) {
	tr := secretdiff.New()
	_ = tr.Record(secretdiff.Event{EnvKey: "A", VaultPath: "p"})
	_ = tr.Record(secretdiff.Event{EnvKey: "B", VaultPath: "p"})
	all := tr.All()
	if len(all) != 2 {
		t.Fatalf("expected 2 events, got %d", len(all))
	}
}

func TestClear_RemovesAllEvents(t *testing.T) {
	tr := secretdiff.New()
	_ = tr.Record(secretdiff.Event{EnvKey: "A", VaultPath: "p"})
	tr.Clear()
	if len(tr.All()) != 0 {
		t.Fatal("expected empty after clear")
	}
}
