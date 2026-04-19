package secretslice_test

import (
	"testing"

	"github.com/yourusername/vaultpipe/internal/secretslice"
)

func TestApply_EmptySecrets_ReturnsError(t *testing.T) {
	s := secretslice.New()
	_, err := s.Apply(map[string]string{})
	if err == nil {
		t.Fatal("expected error for empty secrets")
	}
}

func TestApply_AllKeys_SortedAlphabetically(t *testing.T) {
	s := secretslice.New()
	secrets := map[string]string{"ZEBRA": "z", "ALPHA": "a", "MANGO": "m"}
	entries, err := s.Apply(secrets)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := []string{"ALPHA", "MANGO", "ZEBRA"}
	for i, e := range entries {
		if e.Key != expected[i] {
			t.Errorf("index %d: got key %q, want %q", i, e.Key, expected[i])
		}
	}
}

func TestApply_WithKeys_ReturnsInOrder(t *testing.T) {
	s := secretslice.New(secretslice.WithKeys([]string{"ZEBRA", "ALPHA"}))
	secrets := map[string]string{"ZEBRA": "z", "ALPHA": "a", "MANGO": "m"}
	entries, err := s.Apply(secrets)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
	if entries[0].Key != "ZEBRA" || entries[1].Key != "ALPHA" {
		t.Errorf("unexpected order: %v", entries)
	}
}

func TestApply_WithKeys_MissingKey_ReturnsError(t *testing.T) {
	s := secretslice.New(secretslice.WithKeys([]string{"MISSING"}))
	secrets := map[string]string{"ALPHA": "a"}
	_, err := s.Apply(secrets)
	if err == nil {
		t.Fatal("expected error for missing key")
	}
}

func TestApply_ValuesPreserved(t *testing.T) {
	s := secretslice.New()
	secrets := map[string]string{"FOO": "bar"}
	entries, err := s.Apply(secrets)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if entries[0].Value != "bar" {
		t.Errorf("expected value %q, got %q", "bar", entries[0].Value)
	}
}
