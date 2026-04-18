package secretpin_test

import (
	"testing"

	"github.com/elizabethwanjiku703/vaultpipe/internal/secretpin"
)

func TestPin_And_Get_RoundTrip(t *testing.T) {
	p := secretpin.New()
	if err := p.Pin("DB_PASS", 3); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	pin, err := p.Get("DB_PASS")
	if err != nil {
		t.Fatalf("expected pin, got error: %v", err)
	}
	if pin.Version != 3 {
		t.Errorf("expected version 3, got %d", pin.Version)
	}
	if pin.PinnedAt.IsZero() {
		t.Error("expected PinnedAt to be set")
	}
}

func TestPin_EmptyKey_ReturnsError(t *testing.T) {
	p := secretpin.New()
	if err := p.Pin("", 1); err == nil {
		t.Fatal("expected error for empty key")
	}
}

func TestPin_ZeroVersion_ReturnsError(t *testing.T) {
	p := secretpin.New()
	if err := p.Pin("API_KEY", 0); err == nil {
		t.Fatal("expected error for version < 1")
	}
}

func TestGet_NotPinned_ReturnsErrNotPinned(t *testing.T) {
	p := secretpin.New()
	_, err := p.Get("UNKNOWN")
	if err != secretpin.ErrNotPinned {
		t.Errorf("expected ErrNotPinned, got %v", err)
	}
}

func TestUnpin_RemovesPin(t *testing.T) {
	p := secretpin.New()
	_ = p.Pin("TOKEN", 2)
	if err := p.Unpin("TOKEN"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.IsPinned("TOKEN") {
		t.Error("expected key to be unpinned")
	}
}

func TestUnpin_EmptyKey_ReturnsError(t *testing.T) {
	p := secretpin.New()
	if err := p.Unpin(""); err == nil {
		t.Fatal("expected error for empty key")
	}
}

func TestIsPinned_TrueAndFalse(t *testing.T) {
	p := secretpin.New()
	_ = p.Pin("SECRET", 1)
	if !p.IsPinned("SECRET") {
		t.Error("expected IsPinned true")
	}
	if p.IsPinned("OTHER") {
		t.Error("expected IsPinned false for unknown key")
	}
}

func TestAll_ReturnsCopy(t *testing.T) {
	p := secretpin.New()
	_ = p.Pin("A", 1)
	_ = p.Pin("B", 2)
	all := p.All()
	if len(all) != 2 {
		t.Errorf("expected 2 pins, got %d", len(all))
	}
}
