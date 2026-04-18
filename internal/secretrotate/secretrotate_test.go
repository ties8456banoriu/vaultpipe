package secretrotate_test

import (
	"testing"
	"time"

	"github.com/nicholasgasior/vaultpipe/internal/secretrotate"
)

func TestRecord_And_Latest_RoundTrip(t *testing.T) {
	r := secretrotate.New()
	err := r.Record("DB_PASS", "secret/db", 1, secretrotate.PolicyManual)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	rec, err := r.Latest("DB_PASS")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if rec.EnvKey != "DB_PASS" || rec.VaultPath != "secret/db" || rec.Version != 1 {
		t.Errorf("unexpected record: %+v", rec)
	}
}

func TestRecord_EmptyKey_ReturnsError(t *testing.T) {
	r := secretrotate.New()
	if err := r.Record("", "secret/db", 1, secretrotate.PolicyManual); err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestRecord_EmptyVaultPath_ReturnsError(t *testing.T) {
	r := secretrotate.New()
	if err := r.Record("DB_PASS", "", 1, secretrotate.PolicyManual); err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestRecord_ZeroVersion_ReturnsError(t *testing.T) {
	r := secretrotate.New()
	if err := r.Record("DB_PASS", "secret/db", 0, secretrotate.PolicyManual); err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestLatest_NotFound_ReturnsError(t *testing.T) {
	r := secretrotate.New()
	if _, err := r.Latest("MISSING"); err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestAll_ReturnsCopy(t *testing.T) {
	r := secretrotate.New()
	_ = r.Record("A", "secret/a", 1, secretrotate.PolicyScheduled)
	all := r.All()
	all["A"] = nil
	all2 := r.All()
	if len(all2["A"]) == 0 {
		t.Error("expected original records to be intact")
	}
}

func TestNeedsRotation_OldRecord_ReturnsTrue(t *testing.T) {
	r := secretrotate.New()
	_ = r.Record("TOKEN", "secret/token", 2, secretrotate.PolicyScheduled)
	// Manually age the record by checking with a tiny maxAge
	needs, err := r.NeedsRotation("TOKEN", 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !needs {
		t.Error("expected NeedsRotation to be true")
	}
}

func TestNeedsRotation_FreshRecord_ReturnsFalse(t *testing.T) {
	r := secretrotate.New()
	_ = r.Record("TOKEN", "secret/token", 2, secretrotate.PolicyScheduled)
	needs, err := r.NeedsRotation("TOKEN", 24*time.Hour)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if needs {
		t.Error("expected NeedsRotation to be false")
	}
}
