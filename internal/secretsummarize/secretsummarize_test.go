package secretsummarize_test

import (
	"testing"

	"github.com/wunderkind/vaultpipe/internal/secretsummarize"
)

func TestApply_EmptySecrets_ReturnsError(t *testing.T) {
	s := secretsummarize.New()
	_, err := s.Apply(map[string]string{})
	if err == nil {
		t.Fatal("expected error for empty secrets")
	}
}

func TestApply_NilSecrets_ReturnsError(t *testing.T) {
	s := secretsummarize.New()
	_, err := s.Apply(nil)
	if err == nil {
		t.Fatal("expected error for nil secrets")
	}
}

func TestApply_TotalKeys(t *testing.T) {
	s := secretsummarize.New()
	secrets := map[string]string{"A": "foo", "BB": "bar", "CCC": "baz"}
	sum, err := s.Apply(secrets)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if sum.TotalKeys != 3 {
		t.Errorf("expected TotalKeys=3, got %d", sum.TotalKeys)
	}
}

func TestApply_EmptyValueCount(t *testing.T) {
	s := secretsummarize.New()
	secrets := map[string]string{"A": "val", "B": "", "C": ""}
	sum, err := s.Apply(secrets)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if sum.EmptyValues != 2 {
		t.Errorf("expected EmptyValues=2, got %d", sum.EmptyValues)
	}
}

func TestApply_AvgValueLen(t *testing.T) {
	s := secretsummarize.New()
	// values: "ab"(2), "cdef"(4) => avg 3.0
	secrets := map[string]string{"K1": "ab", "K2": "cdef"}
	sum, err := s.Apply(secrets)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if sum.AvgValueLen != 3.0 {
		t.Errorf("expected AvgValueLen=3.0, got %f", sum.AvgValueLen)
	}
}

func TestApply_LongestAndShortestKey(t *testing.T) {
	s := secretsummarize.New()
	secrets := map[string]string{"A": "1", "BB": "2", "CCC": "3"}
	sum, err := s.Apply(secrets)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if sum.LongestKey != "CCC" {
		t.Errorf("expected LongestKey=CCC, got %s", sum.LongestKey)
	}
	if sum.ShortestKey != "A" {
		t.Errorf("expected ShortestKey=A, got %s", sum.ShortestKey)
	}
}

func TestApply_UniqueValues(t *testing.T) {
	s := secretsummarize.New()
	secrets := map[string]string{"A": "x", "B": "x", "C": "y"}
	sum, err := s.Apply(secrets)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if sum.UniqueValues != 2 {
		t.Errorf("expected UniqueValues=2, got %d", sum.UniqueValues)
	}
}
