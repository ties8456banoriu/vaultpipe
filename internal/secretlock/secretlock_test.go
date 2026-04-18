package secretlock_test

import (
	"testing"

	"github.com/yourusername/vaultpipe/internal/secretlock"
)

func TestLock_And_IsLocked_RoundTrip(t *testing.T) {
	l := secretlock.New()
	if err := l.Lock("DB_PASS"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !l.IsLocked("DB_PASS") {
		t.Error("expected DB_PASS to be locked")
	}
}

func TestLock_EmptyKey_ReturnsError(t *testing.T) {
	l := secretlock.New()
	if err := l.Lock(""); err == nil {
		t.Error("expected error for empty key")
	}
}

func TestUnlock_RemovesLock(t *testing.T) {
	l := secretlock.New()
	_ = l.Lock("API_KEY")
	if err := l.Unlock("API_KEY"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if l.IsLocked("API_KEY") {
		t.Error("expected API_KEY to be unlocked")
	}
}

func TestUnlock_EmptyKey_ReturnsError(t *testing.T) {
	l := secretlock.New()
	if err := l.Unlock(""); err == nil {
		t.Error("expected error for empty key")
	}
}

func TestApply_EmptySecrets_ReturnsError(t *testing.T) {
	l := secretlock.New()
	_, err := l.Apply(map[string]string{})
	if err == nil {
		t.Error("expected error for empty secrets")
	}
}

func TestApply_FiltersLockedKeys(t *testing.T) {
	l := secretlock.New()
	_ = l.Lock("DB_PASS")
	secrets := map[string]string{"DB_PASS": "secret", "API_KEY": "abc123"}
	out, err := l.Apply(secrets)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := out["DB_PASS"]; ok {
		t.Error("expected DB_PASS to be filtered out")
	}
	if out["API_KEY"] != "abc123" {
		t.Error("expected API_KEY to be retained")
	}
}

func TestApply_NoLockedKeys_ReturnsAll(t *testing.T) {
	l := secretlock.New()
	secrets := map[string]string{"A": "1", "B": "2"}
	out, err := l.Apply(secrets)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 2 {
		t.Errorf("expected 2 secrets, got %d", len(out))
	}
}

func TestAll_ReturnsCopy(t *testing.T) {
	l := secretlock.New()
	_ = l.Lock("X")
	_ = l.Lock("Y")
	keys := l.All()
	if len(keys) != 2 {
		t.Errorf("expected 2 locked keys, got %d", len(keys))
	}
}
