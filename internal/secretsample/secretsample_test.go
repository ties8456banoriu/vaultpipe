package secretsample_test

import (
	"testing"

	"github.com/yourusername/vaultpipe/internal/secretsample"
)

var baseSecrets = map[string]string{
	"DB_HOST":     "localhost",
	"DB_PORT":     "5432",
	"DB_USER":     "admin",
	"DB_PASSWORD": "s3cr3t",
	"API_KEY":     "abc123",
}

func TestNew_InvalidN_ReturnsError(t *testing.T) {
	_, err := secretsample.New(0)
	if err == nil {
		t.Fatal("expected error for n=0, got nil")
	}
}

func TestNew_ValidN_ReturnsSampler(t *testing.T) {
	s, err := secretsample.New(2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s == nil {
		t.Fatal("expected non-nil sampler")
	}
}

func TestApply_EmptySecrets_ReturnsError(t *testing.T) {
	s, _ := secretsample.New(2)
	_, err := s.Apply(map[string]string{})
	if err == nil {
		t.Fatal("expected error for empty secrets, got nil")
	}
}

func TestApply_ReturnsSampleOfCorrectSize(t *testing.T) {
	s, _ := secretsample.New(3, secretsample.WithSeed(42))
	out, err := s.Apply(baseSecrets)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 3 {
		t.Errorf("expected 3 keys, got %d", len(out))
	}
}

func TestApply_NLargerThanSecrets_ReturnsAll(t *testing.T) {
	s, _ := secretsample.New(100, secretsample.WithSeed(1))
	out, err := s.Apply(baseSecrets)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != len(baseSecrets) {
		t.Errorf("expected %d keys, got %d", len(baseSecrets), len(out))
	}
}

func TestApply_SameSeedProducesSameResult(t *testing.T) {
	s1, _ := secretsample.New(2, secretsample.WithSeed(99))
	s2, _ := secretsample.New(2, secretsample.WithSeed(99))

	out1, _ := s1.Apply(baseSecrets)
	out2, _ := s2.Apply(baseSecrets)

	for k := range out1 {
		if _, ok := out2[k]; !ok {
			t.Errorf("key %q present in first sample but not second", k)
		}
	}
}

func TestApply_ReturnsNewMap_OriginalUnchanged(t *testing.T) {
	src := map[string]string{"FOO": "bar", "BAZ": "qux"}
	s, _ := secretsample.New(1, secretsample.WithSeed(7))
	out, err := s.Apply(src)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Mutate output; source must be unchanged.
	for k := range out {
		out[k] = "MODIFIED"
	}
	if src["FOO"] == "MODIFIED" || src["BAZ"] == "MODIFIED" {
		t.Error("Apply mutated the original secrets map")
	}
}
