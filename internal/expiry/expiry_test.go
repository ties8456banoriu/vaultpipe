package expiry

import (
	"errors"
	"testing"
	"time"
)

func TestNew_InvalidWarnBefore(t *testing.T) {
	_, err := New(0)
	if err == nil {
		t.Fatal("expected error for zero warnBefore")
	}
}

func TestTrack_EmptyKey_ReturnsError(t *testing.T) {
	tr, _ := New(5 * time.Minute)
	if err := tr.Track("", time.Hour); err == nil {
		t.Fatal("expected error for empty key")
	}
}

func TestTrack_NegativeTTL_ReturnsError(t *testing.T) {
	tr, _ := New(5 * time.Minute)
	if err := tr.Track("MY_SECRET", -time.Second); err == nil {
		t.Fatal("expected error for negative ttl")
	}
}

func TestCheck_UnknownKey_ReturnsErrNoExpiry(t *testing.T) {
	tr, _ := New(5 * time.Minute)
	_, _, err := tr.Check("MISSING")
	if !errors.Is(err, ErrNoExpiry) {
		t.Fatalf("expected ErrNoExpiry, got %v", err)
	}
}

func TestCheck_StatusOK(t *testing.T) {
	tr, _ := New(5 * time.Minute)
	_ = tr.Track("KEY", time.Hour)
	status, remaining, err := tr.Check("KEY")
	if err != nil {
		t.Fatal(err)
	}
	if status != StatusOK {
		t.Fatalf("expected OK, got %s", status)
	}
	if remaining <= 0 {
		t.Fatal("expected positive remaining")
	}
}

func TestCheck_StatusWarning(t *testing.T) {
	tr, _ := New(10 * time.Minute)
	_ = tr.Track("KEY", 5*time.Minute)
	status, _, err := tr.Check("KEY")
	if err != nil {
		t.Fatal(err)
	}
	if status != StatusWarning {
		t.Fatalf("expected Warning, got %s", status)
	}
}

func TestCheck_StatusExpired(t *testing.T) {
	tr, _ := New(5 * time.Minute)
	fixed := time.Now()
	tr.now = func() time.Time { return fixed }
	_ = tr.Track("KEY", time.Second)
	tr.now = func() time.Time { return fixed.Add(2 * time.Second) }
	status, remaining, err := tr.Check("KEY")
	if err != nil {
		t.Fatal(err)
	}
	if status != StatusExpired {
		t.Fatalf("expected Expired, got %s", status)
	}
	if remaining != 0 {
		t.Fatalf("expected zero remaining for expired, got %v", remaining)
	}
}

func TestRemove_DeletesEntry(t *testing.T) {
	tr, _ := New(5 * time.Minute)
	_ = tr.Track("KEY", time.Hour)
	tr.Remove("KEY")
	_, _, err := tr.Check("KEY")
	if !errors.Is(err, ErrNoExpiry) {
		t.Fatal("expected ErrNoExpiry after remove")
	}
}
