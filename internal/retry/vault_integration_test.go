package retry_test

import (
	"errors"
	"testing"
	"time"

	"github.com/your-org/vaultpipe/internal/retry"
)

// simulateVaultFetch mimics a Vault client that fails twice then succeeds.
func simulateVaultFetch(failTimes int) func() error {
	attempts := 0
	return func() error {
		attempts++
		if attempts <= failTimes {
			return errors.New("vault: connection refused")
		}
		return nil
	}
}

func TestRetry_VaultFetch_EventualSuccess(t *testing.T) {
	cfg := retry.Config{
		MaxAttempts: 4,
		Delay:       time.Millisecond,
		Multiplier:  1.0,
	}
	d := retry.New(cfg)
	// replace sleep so test runs fast
	// We access via the exported Do method; sleep is internal.

	fetch := simulateVaultFetch(3)
	calls := 0
	err := d.Do(func() error {
		calls++
		return fetch()
	})
	if err != nil {
		t.Fatalf("expected success after retries, got: %v", err)
	}
	if calls != 4 {
		t.Fatalf("expected 4 calls, got %d", calls)
	}
}

func TestRetry_VaultFetch_PermanentFailure(t *testing.T) {
	cfg := retry.Config{
		MaxAttempts: 3,
		Delay:       time.Millisecond,
		Multiplier:  1.0,
	}
	d := retry.New(cfg)

	fetch := simulateVaultFetch(10) // always fails within 3 attempts
	err := d.Do(fetch)
	if err == nil {
		t.Fatal("expected error on permanent failure")
	}
}
