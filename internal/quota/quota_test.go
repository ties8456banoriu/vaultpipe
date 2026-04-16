package quota_test

import (
	"errors"
	"testing"
	"time"

	"github.com/yourusername/vaultpipe/internal/quota"
)

func TestNew_InvalidWindow(t *testing.T) {
	_, err := quota.New(0)
	if err == nil {
		t.Fatal("expected error for zero window")
	}
}

func TestSetLimit_EmptyKey_ReturnsError(t *testing.T) {
	e, _ := quota.New(time.Minute)
	if err := e.SetLimit("", 5); err == nil {
		t.Fatal("expected error for empty key")
	}
}

func TestSetLimit_ZeroMax_ReturnsError(t *testing.T) {
	e, _ := quota.New(time.Minute)
	if err := e.SetLimit("KEY", 0); err == nil {
		t.Fatal("expected error for zero max")
	}
}

func TestCheck_NoQuota_ReturnsErrNoQuota(t *testing.T) {
	e, _ := quota.New(time.Minute)
	err := e.Check("UNKNOWN")
	if !errors.Is(err, quota.ErrNoQuota) {
		t.Fatalf("expected ErrNoQuota, got %v", err)
	}
}

func TestCheck_WithinLimit_ReturnsNil(t *testing.T) {
	e, _ := quota.New(time.Minute)
	_ = e.SetLimit("DB_PASS", 3)
	for i := 0; i < 3; i++ {
		if err := e.Check("DB_PASS"); err != nil {
			t.Fatalf("unexpected error on attempt %d: %v", i+1, err)
		}
	}
}

func TestCheck_ExceedsLimit_ReturnsErrQuotaExceeded(t *testing.T) {
	e, _ := quota.New(time.Minute)
	_ = e.SetLimit("DB_PASS", 2)
	_ = e.Check("DB_PASS")
	_ = e.Check("DB_PASS")
	err := e.Check("DB_PASS")
	if !errors.Is(err, quota.ErrQuotaExceeded) {
		t.Fatalf("expected ErrQuotaExceeded, got %v", err)
	}
}

func TestCheck_WindowExpiry_ResetsCount(t *testing.T) {
	e, _ := quota.New(20 * time.Millisecond)
	_ = e.SetLimit("KEY", 1)
	_ = e.Check("KEY") // uses up quota
	time.Sleep(30 * time.Millisecond)
	if err := e.Check("KEY"); err != nil {
		t.Fatalf("expected nil after window reset, got %v", err)
	}
}

func TestRemaining_DecreasesWithChecks(t *testing.T) {
	e, _ := quota.New(time.Minute)
	_ = e.SetLimit("API_KEY", 5)
	_ = e.Check("API_KEY")
	_ = e.Check("API_KEY")
	n, err := e.Remaining("API_KEY")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n != 3 {
		t.Fatalf("expected 3 remaining, got %d", n)
	}
}

func TestRemaining_NoQuota_ReturnsErrNoQuota(t *testing.T) {
	e, _ := quota.New(time.Minute)
	_, err := e.Remaining("MISSING")
	if !errors.Is(err, quota.ErrNoQuota) {
		t.Fatalf("expected ErrNoQuota, got %v", err)
	}
}
