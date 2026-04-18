package secretfreeze_test

import (
	"testing"

	"github.com/yourusername/vaultpipe/internal/secretfreeze"
)

func TestFreeze_And_IsFrozen_RoundTrip(t *testing.T) {
	f := secretfreeze.New()
	if err := f.Freeze("DB_PASS"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !f.IsFrozen("DB_PASS") {
		t.Error("expected DB_PASS to be frozen")
	}
}

func TestFreeze_EmptyKey_ReturnsError(t *testing.T) {
	f := secretfreeze.New()
	if err := f.Freeze(""); err == nil {
		t.Error("expected error for empty key")
	}
}

func TestThaw_RemovesFrozen(t *testing.T) {
	f := secretfreeze.New()
	_ = f.Freeze("API_KEY")
	if err := f.Thaw("API_KEY"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f.IsFrozen("API_KEY") {
		t.Error("expected API_KEY to be thawed")
	}
}

func TestThaw_NotFrozen_ReturnsError(t *testing.T) {
	f := secretfreeze.New()
	if err := f.Thaw("MISSING"); err == nil {
		t.Error("expected error when thawing non-frozen key")
	}
}

func TestApply_EmptyIncoming_ReturnsError(t *testing.T) {
	f := secretfreeze.New()
	_, err := f.Apply(map[string]string{"A": "1"}, map[string]string{})
	if err == nil {
		t.Error("expected error for empty incoming secrets")
	}
}

func TestApply_FrozenKeyPreservedFromBase(t *testing.T) {
	f := secretfreeze.New()
	_ = f.Freeze("DB_PASS")
	base := map[string]string{"DB_PASS": "original"}
	incoming := map[string]string{"DB_PASS": "changed", "OTHER": "val"}
	out, err := f.Apply(base, incoming)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["DB_PASS"] != "original" {
		t.Errorf("expected original value, got %q", out["DB_PASS"])
	}
	if out["OTHER"] != "val" {
		t.Errorf("expected val, got %q", out["OTHER"])
	}
}

func TestApply_FrozenKeyMissingFromBase_ReturnsError(t *testing.T) {
	f := secretfreeze.New()
	_ = f.Freeze("DB_PASS")
	_, err := f.Apply(map[string]string{}, map[string]string{"DB_PASS": "x"})
	if err == nil {
		t.Error("expected error when frozen key missing from base")
	}
}

func TestIsFrozen_UnknownKey_ReturnsFalse(t *testing.T) {
	f := secretfreeze.New()
	if f.IsFrozen("NOPE") {
		t.Error("expected false for unknown key")
	}
}
