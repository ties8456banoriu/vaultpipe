package deporder_test

import (
	"errors"
	"testing"

	"github.com/eliziario/vaultpipe/internal/deporder"
)

var baseSecrets = map[string]string{
	"DB_URL":       "postgres://localhost/db",
	"DB_USER":      "admin",
	"APP_SECRET":   "s3cr3t",
}

func TestResolve_EmptySecrets_ReturnsError(t *testing.T) {
	r := deporder.New(nil)
	_, err := r.Resolve(map[string]string{})
	if err == nil {
		t.Fatal("expected error for empty secrets")
	}
}

func TestResolve_NoDeps_ReturnsAllKeys(t *testing.T) {
	r := deporder.New(nil)
	keys, err := r.Resolve(baseSecrets)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(keys) != len(baseSecrets) {
		t.Fatalf("expected %d keys, got %d", len(baseSecrets), len(keys))
	}
}

func TestResolve_WithDeps_OrderRespected(t *testing.T) {
	secrets := map[string]string{
		"A": "1",
		"B": "2",
		"C": "3",
	}
	deps := map[string][]string{
		"C": {"B"},
		"B": {"A"},
	}
	r := deporder.New(deps)
	keys, err := r.Resolve(secrets)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	pos := make(map[string]int, len(keys))
	for i, k := range keys {
		pos[k] = i
	}
	if pos["A"] >= pos["B"] {
		t.Errorf("expected A before B, got order %v", keys)
	}
	if pos["B"] >= pos["C"] {
		t.Errorf("expected B before C, got order %v", keys)
	}
}

func TestResolve_CycleDetected_ReturnsErrCycle(t *testing.T) {
	secrets := map[string]string{"X": "1", "Y": "2"}
	deps := map[string][]string{
		"X": {"Y"},
		"Y": {"X"},
	}
	r := deporder.New(deps)
	_, err := r.Resolve(secrets)
	if !errors.Is(err, deporder.ErrCycle) {
		t.Fatalf("expected ErrCycle, got %v", err)
	}
}

func TestResolve_UnknownDep_ReturnsErrUnknownKey(t *testing.T) {
	secrets := map[string]string{"A": "1"}
	deps := map[string][]string{
		"A": {"MISSING"},
	}
	r := deporder.New(deps)
	_, err := r.Resolve(secrets)
	if !errors.Is(err, deporder.ErrUnknownKey) {
		t.Fatalf("expected ErrUnknownKey, got %v", err)
	}
}
