package secretinspect_test

import (
	"errors"
	"testing"

	"github.com/yourusername/vaultpipe/internal/secretinspect"
)

func TestInspect_EmptySecrets_ReturnsError(t *testing.T) {
	insp := secretinspect.New()
	_, err := insp.Inspect(map[string]string{}, "KEY")
	if !errors.Is(err, secretinspect.ErrEmptySecrets) {
		t.Fatalf("expected ErrEmptySecrets, got %v", err)
	}
}

func TestInspect_UnknownKey_ReturnsError(t *testing.T) {
	insp := secretinspect.New()
	_, err := insp.Inspect(map[string]string{"FOO": "bar"}, "MISSING")
	if !errors.Is(err, secretinspect.ErrUnknownKey) {
		t.Fatalf("expected ErrUnknownKey, got %v", err)
	}
}

func TestInspect_BasicInfo(t *testing.T) {
	insp := secretinspect.New()
	secrets := map[string]string{"DB_PASS": "s3cr3t"}
	info, err := insp.Inspect(secrets, "DB_PASS")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if info.Key != "DB_PASS" {
		t.Errorf("expected key DB_PASS, got %q", info.Key)
	}
	if info.Length != 6 {
		t.Errorf("expected length 6, got %d", info.Length)
	}
	if info.IsEmpty {
		t.Error("expected IsEmpty false")
	}
}

func TestInspect_EmptyValue(t *testing.T) {
	insp := secretinspect.New()
	info, err := insp.Inspect(map[string]string{"EMPTY": ""}, "EMPTY")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !info.IsEmpty {
		t.Error("expected IsEmpty true")
	}
	if info.Uppercase {
		t.Error("expected Uppercase false for empty value")
	}
}

func TestInspect_HasSpaces(t *testing.T) {
	insp := secretinspect.New()
	info, err := insp.Inspect(map[string]string{"MSG": "hello world"}, "MSG")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !info.HasSpaces {
		t.Error("expected HasSpaces true")
	}
}

func TestInspect_Uppercase(t *testing.T) {
	insp := secretinspect.New()
	info, err := insp.Inspect(map[string]string{"ENV": "PRODUCTION"}, "ENV")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !info.Uppercase {
		t.Error("expected Uppercase true")
	}
}

func TestAll_EmptySecrets_ReturnsError(t *testing.T) {
	insp := secretinspect.New()
	_, err := insp.All(map[string]string{})
	if !errors.Is(err, secretinspect.ErrEmptySecrets) {
		t.Fatalf("expected ErrEmptySecrets, got %v", err)
	}
}

func TestAll_ReturnsSortedKeys(t *testing.T) {
	insp := secretinspect.New()
	secrets := map[string]string{"ZEBRA": "z", "ALPHA": "a", "MANGO": "m"}
	infos, err := insp.All(secrets)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(infos) != 3 {
		t.Fatalf("expected 3 infos, got %d", len(infos))
	}
	expected := []string{"ALPHA", "MANGO", "ZEBRA"}
	for i, info := range infos {
		if info.Key != expected[i] {
			t.Errorf("position %d: expected %q, got %q", i, expected[i], info.Key)
		}
	}
}
