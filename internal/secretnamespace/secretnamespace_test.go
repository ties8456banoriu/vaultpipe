package secretnamespace_test

import (
	"testing"

	"github.com/yourusername/vaultpipe/internal/secretnamespace"
)

func baseSecrets() map[string]string {
	return map[string]string{"KEY_A": "val_a", "KEY_B": "val_b"}
}

func TestSet_And_Get_RoundTrip(t *testing.T) {
	m := secretnamespace.New()
	if err := m.Set("prod", baseSecrets()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got, err := m.Get("prod")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["KEY_A"] != "val_a" {
		t.Errorf("expected val_a, got %s", got["KEY_A"])
	}
}

func TestSet_EmptyNamespace_ReturnsError(t *testing.T) {
	m := secretnamespace.New()
	if err := m.Set("", baseSecrets()); err == nil {
		t.Fatal("expected error for empty namespace")
	}
}

func TestSet_EmptySecrets_ReturnsError(t *testing.T) {
	m := secretnamespace.New()
	if err := m.Set("prod", map[string]string{}); err == nil {
		t.Fatal("expected error for empty secrets")
	}
}

func TestGet_NotFound_ReturnsError(t *testing.T) {
	m := secretnamespace.New()
	_, err := m.Get("missing")
	if err == nil {
		t.Fatal("expected error for missing namespace")
	}
}

func TestGet_ReturnsCopy(t *testing.T) {
	m := secretnamespace.New()
	_ = m.Set("prod", baseSecrets())
	got, _ := m.Get("prod")
	got["KEY_A"] = "mutated"
	again, _ := m.Get("prod")
	if again["KEY_A"] == "mutated" {
		t.Error("Get should return a copy, not a reference")
	}
}

func TestDelete_RemovesNamespace(t *testing.T) {
	m := secretnamespace.New()
	_ = m.Set("prod", baseSecrets())
	if err := m.Delete("prod"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_, err := m.Get("prod")
	if err == nil {
		t.Fatal("expected namespace to be deleted")
	}
}

func TestDelete_NotFound_ReturnsError(t *testing.T) {
	m := secretnamespace.New()
	if err := m.Delete("ghost"); err == nil {
		t.Fatal("expected error deleting missing namespace")
	}
}

func TestList_ReturnsAllNamespaces(t *testing.T) {
	m := secretnamespace.New()
	_ = m.Set("prod", baseSecrets())
	_ = m.Set("staging", baseSecrets())
	list := m.List()
	if len(list) != 2 {
		t.Errorf("expected 2 namespaces, got %d", len(list))
	}
}

func TestMerge_CombinesNamespaces(t *testing.T) {
	m := secretnamespace.New()
	_ = m.Set("ns1", map[string]string{"A": "1"})
	_ = m.Set("ns2", map[string]string{"B": "2"})
	result, err := m.Merge("ns1", "ns2")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result["A"] != "1" || result["B"] != "2" {
		t.Errorf("unexpected merge result: %v", result)
	}
}

func TestMerge_LaterNamespaceWins(t *testing.T) {
	m := secretnamespace.New()
	_ = m.Set("ns1", map[string]string{"KEY": "first"})
	_ = m.Set("ns2", map[string]string{"KEY": "second"})
	result, _ := m.Merge("ns1", "ns2")
	if result["KEY"] != "second" {
		t.Errorf("expected second, got %s", result["KEY"])
	}
}

func TestMerge_MissingNamespace_ReturnsError(t *testing.T) {
	m := secretnamespace.New()
	_ = m.Set("ns1", baseSecrets())
	_, err := m.Merge("ns1", "missing")
	if err == nil {
		t.Fatal("expected error for missing namespace in merge")
	}
}
