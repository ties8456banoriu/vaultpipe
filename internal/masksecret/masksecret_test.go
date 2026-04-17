package masksecret

import (
	"testing"
)

func TestNew_EmptyPlaceholder_UsesDefault(t *testing.T) {
	m, err := New(PolicyAll, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if m.placeholder != "***" {
		t.Errorf("expected *** got %s", m.placeholder)
	}
}

func TestApply_EmptySecrets_ReturnsError(t *testing.T) {
	m, _ := New(PolicyAll, "")
	_, err := m.Apply(map[string]string{})
	if err == nil {
		t.Fatal("expected error for empty secrets")
	}
}

func TestApply_PolicyAll_MasksEverything(t *testing.T) {
	m, _ := New(PolicyAll, "REDACTED")
	out, err := m.Apply(map[string]string{"APP_NAME": "myapp", "DB_PASS": "s3cr3t"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for k, v := range out {
		if v != "REDACTED" {
			t.Errorf("key %s: expected REDACTED got %s", k, v)
		}
	}
}

func TestApply_PolicyNone_RetainsValues(t *testing.T) {
	m, _ := New(PolicyNone, "")
	secrets := map[string]string{"APP_NAME": "myapp", "API_TOKEN": "tok123"}
	out, err := m.Apply(secrets)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for k, want := range secrets {
		if out[k] != want {
			t.Errorf("key %s: expected %s got %s", k, want, out[k])
		}
	}
}

func TestApply_PolicySensitive_OnlyMasksSensitiveKeys(t *testing.T) {
	m, _ := New(PolicySensitive, "***")
	secrets := map[string]string{
		"APP_NAME":    "myapp",
		"DB_PASSWORD": "hunter2",
		"API_TOKEN":   "tok123",
		"LOG_LEVEL":   "debug",
	}
	out, err := m.Apply(secrets)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["APP_NAME"] != "myapp" {
		t.Errorf("APP_NAME should not be masked")
	}
	if out["LOG_LEVEL"] != "debug" {
		t.Errorf("LOG_LEVEL should not be masked")
	}
	if out["DB_PASSWORD"] != "***" {
		t.Errorf("DB_PASSWORD should be masked")
	}
	if out["API_TOKEN"] != "***" {
		t.Errorf("API_TOKEN should be masked")
	}
}

func TestApply_ReturnsCopy_DoesNotMutateInput(t *testing.T) {
	m, _ := New(PolicyAll, "X")
	input := map[string]string{"KEY": "value"}
	_, _ = m.Apply(input)
	if input["KEY"] != "value" {
		t.Error("Apply must not mutate the input map")
	}
}
