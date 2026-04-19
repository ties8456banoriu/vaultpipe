package secretmergetag

import (
	"testing"
)

func TestNew_InvalidPolicy(t *testing.T) {
	_, err := New("invalid")
	if err == nil {
		t.Fatal("expected error for invalid policy")
	}
}

func TestNew_ValidPolicy(t *testing.T) {
	for _, p := range []ConflictPolicy{PolicySkip, PolicyOverwrite, PolicyError} {
		_, err := New(p)
		if err != nil {
			t.Fatalf("unexpected error for policy %q: %v", p, err)
		}
	}
}

func TestMerge_NoSources_ReturnsError(t *testing.T) {
	m, _ := New(PolicySkip)
	_, err := m.Merge()
	if err == nil {
		t.Fatal("expected error for no sources")
	}
}

func TestMerge_NoConflict_CombinesAll(t *testing.T) {
	m, _ := New(PolicySkip)
	src1 := map[string]map[string]string{"DB_URL": {"env": "prod"}}
	src2 := map[string]map[string]string{"API_KEY": {"team": "backend"}}
	out, err := m.Merge(src1, src2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["DB_URL"]["env"] != "prod" {
		t.Errorf("expected prod, got %q", out["DB_URL"]["env"])
	}
	if out["API_KEY"]["team"] != "backend" {
		t.Errorf("expected backend, got %q", out["API_KEY"]["team"])
	}
}

func TestMerge_PolicySkip_KeepsFirst(t *testing.T) {
	m, _ := New(PolicySkip)
	src1 := map[string]map[string]string{"KEY": {"env": "prod"}}
	src2 := map[string]map[string]string{"KEY": {"env": "staging"}}
	out, err := m.Merge(src1, src2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["KEY"]["env"] != "prod" {
		t.Errorf("expected prod, got %q", out["KEY"]["env"])
	}
}

func TestMerge_PolicyOverwrite_KeepsLast(t *testing.T) {
	m, _ := New(PolicyOverwrite)
	src1 := map[string]map[string]string{"KEY": {"env": "prod"}}
	src2 := map[string]map[string]string{"KEY": {"env": "staging"}}
	out, err := m.Merge(src1, src2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["KEY"]["env"] != "staging" {
		t.Errorf("expected staging, got %q", out["KEY"]["env"])
	}
}

func TestMerge_PolicyError_ReturnsErrorOnConflict(t *testing.T) {
	m, _ := New(PolicyError)
	src1 := map[string]map[string]string{"KEY": {"env": "prod"}}
	src2 := map[string]map[string]string{"KEY": {"env": "staging"}}
	_, err := m.Merge(src1, src2)
	if err == nil {
		t.Fatal("expected error on conflict")
	}
}
