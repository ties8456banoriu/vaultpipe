package secrettag_test

import (
	"testing"

	"github.com/yourusername/vaultpipe/internal/secrettag"
)

func TestTag_And_Get_RoundTrip(t *testing.T) {
	tr := secrettag.New()
	if err := tr.Tag("DB_PASS", []string{"sensitive", "db"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	tags, err := tr.Get("DB_PASS")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(tags) != 2 || tags[0] != "sensitive" || tags[1] != "db" {
		t.Fatalf("unexpected tags: %v", tags)
	}
}

func TestTag_EmptyKey_ReturnsError(t *testing.T) {
	tr := secrettag.New()
	if err := tr.Tag("", []string{"x"}); err != secrettag.ErrEmptyKey {
		t.Fatalf("expected ErrEmptyKey, got %v", err)
	}
}

func TestTag_EmptyTags_ReturnsError(t *testing.T) {
	tr := secrettag.New()
	if err := tr.Tag("KEY", nil); err != secrettag.ErrEmptyTags {
		t.Fatalf("expected ErrEmptyTags, got %v", err)
	}
}

func TestGet_UnknownKey_ReturnsError(t *testing.T) {
	tr := secrettag.New()
	_, err := tr.Get("MISSING")
	if err != secrettag.ErrUnknownKey {
		t.Fatalf("expected ErrUnknownKey, got %v", err)
	}
}

func TestHasTag_True(t *testing.T) {
	tr := secrettag.New()
	_ = tr.Tag("API_KEY", []string{"sensitive"})
	if !tr.HasTag("API_KEY", "sensitive") {
		t.Fatal("expected HasTag to return true")
	}
}

func TestHasTag_False(t *testing.T) {
	tr := secrettag.New()
	_ = tr.Tag("API_KEY", []string{"sensitive"})
	if tr.HasTag("API_KEY", "db") {
		t.Fatal("expected HasTag to return false")
	}
}

func TestAll_ReturnsCopy(t *testing.T) {
	tr := secrettag.New()
	_ = tr.Tag("A", []string{"x"})
	_ = tr.Tag("B", []string{"y", "z"})
	all := tr.All()
	if len(all) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(all))
	}
	// mutate copy, ensure original unaffected
	all["A"][0] = "mutated"
	orig, _ := tr.Get("A")
	if orig[0] == "mutated" {
		t.Fatal("All() should return a deep copy")
	}
}
