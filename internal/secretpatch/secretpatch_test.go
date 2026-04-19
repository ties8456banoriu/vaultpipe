package secretpatch_test

import (
	"testing"

	"github.com/yourusername/vaultpipe/internal/secretpatch"
)

func base() map[string]string {
	return map[string]string{"A": "1", "B": "2"}
}

func patch() map[string]string {
	return map[string]string{"B": "99", "C": "3"}
}

func TestNew_InvalidPolicy(t *testing.T) {
	_, err := secretpatch.New("invalid")
	if err == nil {
		t.Fatal("expected error for invalid policy")
	}
}

func TestNew_ValidPolicy(t *testing.T) {
	for _, p := range []secretpatch.Policy{secretpatch.PolicyOverwrite, secretpatch.PolicyKeep} {
		_, err := secretpatch.New(p)
		if err != nil {
			t.Fatalf("unexpected error for policy %q: %v", p, err)
		}
	}
}

func TestApply_EmptyBase_ReturnsError(t *testing.T) {
	p, _ := secretpatch.New(secretpatch.PolicyOverwrite)
	_, err := p.Apply(nil, patch(), nil)
	if err != secretpatch.ErrEmptyBase {
		t.Fatalf("expected ErrEmptyBase, got %v", err)
	}
}

func TestApply_EmptyPatch_ReturnsError(t *testing.T) {
	p, _ := secretpatch.New(secretpatch.PolicyOverwrite)
	_, err := p.Apply(base(), nil, nil)
	if err != secretpatch.ErrEmptyPatch {
		t.Fatalf("expected ErrEmptyPatch, got %v", err)
	}
}

func TestApply_PolicyOverwrite_ReplacesConflict(t *testing.T) {
	p, _ := secretpatch.New(secretpatch.PolicyOverwrite)
	out, err := p.Apply(base(), patch(), nil)
	if err != nil {
		t.Fatal(err)
	}
	if out["B"] != "99" {
		t.Errorf("expected B=99, got %s", out["B"])
	}
	if out["C"] != "3" {
		t.Errorf("expected C=3, got %s", out["C"])
	}
	if out["A"] != "1" {
		t.Errorf("expected A=1, got %s", out["A"])
	}
}

func TestApply_PolicyKeep_PreservesBase(t *testing.T) {
	p, _ := secretpatch.New(secretpatch.PolicyKeep)
	out, err := p.Apply(base(), patch(), nil)
	if err != nil {
		t.Fatal(err)
	}
	if out["B"] != "2" {
		t.Errorf("expected B=2 (kept), got %s", out["B"])
	}
	if out["C"] != "3" {
		t.Errorf("expected C=3 (new), got %s", out["C"])
	}
}

func TestApply_WithKeys_OnlyAppliesListed(t *testing.T) {
	p, _ := secretpatch.New(secretpatch.PolicyOverwrite)
	out, err := p.Apply(base(), patch(), []string{"C"})
	if err != nil {
		t.Fatal(err)
	}
	if out["B"] != "2" {
		t.Errorf("B should remain 2, got %s", out["B"])
	}
	if out["C"] != "3" {
		t.Errorf("C should be 3, got %s", out["C"])
	}
}

func TestApply_EmptyKeyInList_ReturnsError(t *testing.T) {
	p, _ := secretpatch.New(secretpatch.PolicyOverwrite)
	_, err := p.Apply(base(), patch(), []string{""})
	if err != secretpatch.ErrEmptyKey {
		t.Fatalf("expected ErrEmptyKey, got %v", err)
	}
}
