package secretshuffle_test

import (
	"math/rand"
	"sort"
	"testing"

	"github.com/yourusername/vaultpipe/internal/secretshuffle"
)

func baseSecrets() map[string]string {
	return map[string]string{
		"DB_HOST":     "localhost",
		"DB_PORT":     "5432",
		"API_KEY":     "abc123",
		"SECRET_TOKEN": "xyz789",
	}
}

func TestShuffle_EmptySecrets_ReturnsError(t *testing.T) {
	s := secretshuffle.New(nil)
	_, err := s.Shuffle(map[string]string{})
	if err == nil {
		t.Fatal("expected error for empty secrets")
	}
}

func TestShuffle_ReturnsAllKeys(t *testing.T) {
	s := secretshuffle.New(rand.NewSource(42))
	secrets := baseSecrets()
	keys, err := s.Shuffle(secrets)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(keys) != len(secrets) {
		t.Fatalf("expected %d keys, got %d", len(secrets), len(keys))
	}
	sorted := make([]string, len(keys))
	copy(sorted, keys)
	sort.Strings(sorted)
	for _, k := range sorted {
		if _, ok := secrets[k]; !ok {
			t.Errorf("unexpected key %q in result", k)
		}
	}
}

func TestShuffle_DifferentSeedsProduceDifferentOrders(t *testing.T) {
	secrets := baseSecrets()
	s1 := secretshuffle.New(rand.NewSource(1))
	s2 := secretshuffle.New(rand.NewSource(99))
	var diff bool
	for i := 0; i < 20; i++ {
		k1, _ := s1.Shuffle(secrets)
		k2, _ := s2.Shuffle(secrets)
		for j := range k1 {
			if k1[j] != k2[j] {
				diff = true
			}
		}
	}
	if !diff {
		t.Error("expected at least one ordering difference across seeds")
	}
}

func TestApply_EmptySecrets_ReturnsError(t *testing.T) {
	s := secretshuffle.New(nil)
	_, _, err := s.Apply(map[string]string{})
	if err == nil {
		t.Fatal("expected error for empty secrets")
	}
}

func TestApply_RetainsAllValues(t *testing.T) {
	s := secretshuffle.New(rand.NewSource(7))
	secrets := baseSecrets()
	out, keys, err := s.Apply(secrets)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != len(secrets) {
		t.Fatalf("expected %d entries, got %d", len(secrets), len(out))
	}
	for _, k := range keys {
		if out[k] != secrets[k] {
			t.Errorf("value mismatch for key %q: got %q want %q", k, out[k], secrets[k])
		}
	}
}
