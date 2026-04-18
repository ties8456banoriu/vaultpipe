package secretpriority_test

import (
	"testing"

	"github.com/celtechstarter/vaultpipe/internal/secretpriority"
)

func TestAdd_And_Resolve_RoundTrip(t *testing.T) {
	m := secretpriority.New()
	if err := m.Add("DB_PASS", "low", 1); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := m.Add("DB_PASS", "high", 10); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	v, err := m.Resolve("DB_PASS")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v != "high" {
		t.Errorf("expected %q, got %q", "high", v)
	}
}

func TestAdd_EmptyKey_ReturnsError(t *testing.T) {
	m := secretpriority.New()
	if err := m.Add("", "val", 1); err == nil {
		t.Fatal("expected error for empty key")
	}
}

func TestResolve_UnknownKey_ReturnsError(t *testing.T) {
	m := secretpriority.New()
	if _, err := m.Resolve("MISSING"); err == nil {
		t.Fatal("expected error for unknown key")
	}
}

func TestResolve_TiedPriority_ReturnsFirstAdded(t *testing.T) {
	m := secretpriority.New()
	_ = m.Add("KEY", "first", 5)
	_ = m.Add("KEY", "second", 5)
	v, err := m.Resolve("KEY")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v != "first" {
		t.Errorf("expected %q, got %q", "first", v)
	}
}

func TestResolveAll_ReturnsAllKeys(t *testing.T) {
	m := secretpriority.New()
	_ = m.Add("A", "a-low", 1)
	_ = m.Add("A", "a-high", 9)
	_ = m.Add("B", "b-only", 3)

	out, err := m.ResolveAll()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["A"] != "a-high" {
		t.Errorf("A: expected %q, got %q", "a-high", out["A"])
	}
	if out["B"] != "b-only" {
		t.Errorf("B: expected %q, got %q", "b-only", out["B"])
	}
}

func TestResolveAll_Empty_ReturnsError(t *testing.T) {
	m := secretpriority.New()
	if _, err := m.ResolveAll(); err == nil {
		t.Fatal("expected error for empty manager")
	}
}
