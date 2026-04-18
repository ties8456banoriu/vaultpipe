package secretsort_test

import (
	"testing"

	"github.com/yourusername/vaultpipe/internal/secretsort"
)

var base = map[string]string{
	"ZEBRA": "last",
	"ALPHA": "first",
	"MANGO": "middle",
}

func TestNew_InvalidOrder(t *testing.T) {
	_, err := secretsort.New("random", false)
	if err != secretsort.ErrInvalidOrder {
		t.Fatalf("expected ErrInvalidOrder, got %v", err)
	}
}

func TestNew_ValidOrder(t *testing.T) {
	_, err := secretsort.New(secretsort.OrderAsc, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestApply_EmptySecrets_ReturnsError(t *testing.T) {
	s, _ := secretsort.New(secretsort.OrderAsc, false)
	_, err := s.Apply(map[string]string{})
	if err != secretsort.ErrEmptySecrets {
		t.Fatalf("expected ErrEmptySecrets, got %v", err)
	}
}

func TestApply_SortByKeyAsc(t *testing.T) {
	s, _ := secretsort.New(secretsort.OrderAsc, false)
	entries, err := s.Apply(base)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if entries[0].Key != "ALPHA" || entries[1].Key != "MANGO" || entries[2].Key != "ZEBRA" {
		t.Errorf("unexpected order: %v", entries)
	}
}

func TestApply_SortByKeyDesc(t *testing.T) {
	s, _ := secretsort.New(secretsort.OrderDesc, false)
	entries, err := s.Apply(base)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if entries[0].Key != "ZEBRA" || entries[2].Key != "ALPHA" {
		t.Errorf("unexpected order: %v", entries)
	}
}

func TestApply_SortByValueAsc(t *testing.T) {
	s, _ := secretsort.New(secretsort.OrderAsc, true)
	entries, err := s.Apply(base)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// values: first < last < middle
	if entries[0].Value != "first" || entries[1].Value != "last" || entries[2].Value != "middle" {
		t.Errorf("unexpected value order: %v", entries)
	}
}

func TestApply_ReturnsCopyNotReference(t *testing.T) {
	s, _ := secretsort.New(secretsort.OrderAsc, false)
	entries, _ := s.Apply(base)
	entries[0].Key = "MUTATED"
	// original map should be unaffected
	if _, ok := base["MUTATED"]; ok {
		t.Error("original map was mutated")
	}
}
