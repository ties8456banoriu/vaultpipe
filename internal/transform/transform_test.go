package transform

import (
	"testing"
)

var baseSecrets = map[string]string{
	"DB_PASSWORD": "  secret123  ",
	"API_KEY":     "myapikey",
	"APP_ENV":     "production",
}

func TestApply_NoRules_ReturnsUnchanged(t *testing.T) {
	tr := NewTransformer(nil)
	out, err := tr.Apply(baseSecrets)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["API_KEY"] != "myapikey" {
		t.Errorf("expected unchanged value, got %q", out["API_KEY"])
	}
}

func TestApply_EmptySecrets_ReturnsError(t *testing.T) {
	tr := NewTransformer(nil)
	_, err := tr.Apply(map[string]string{})
	if err == nil {
		t.Fatal("expected error for empty secrets")
	}
}

func TestApply_TrimRule(t *testing.T) {
	tr := NewTransformer([]Rule{{Key: "DB_PASSWORD", Type: "trim"}})
	out, err := tr.Apply(baseSecrets)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["DB_PASSWORD"] != "secret123" {
		t.Errorf("expected trimmed value, got %q", out["DB_PASSWORD"])
	}
}

func TestApply_UpperRule(t *testing.T) {
	tr := NewTransformer([]Rule{{Key: "API_KEY", Type: "upper"}})
	out, err := tr.Apply(baseSecrets)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["API_KEY"] != "MYAPIKEY" {
		t.Errorf("expected upper value, got %q", out["API_KEY"])
	}
}

func TestApply_PrefixRule(t *testing.T) {
	tr := NewTransformer([]Rule{{Key: "APP_ENV", Type: "prefix", Arg: "vp_"}})
	out, err := tr.Apply(baseSecrets)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["APP_ENV"] != "vp_production" {
		t.Errorf("expected prefixed value, got %q", out["APP_ENV"])
	}
}

func TestApply_UnknownKey_IsSkipped(t *testing.T) {
	tr := NewTransformer([]Rule{{Key: "NONEXISTENT", Type: "upper"}})
	_, err := tr.Apply(baseSecrets)
	if err != nil {
		t.Fatalf("unexpected error for unknown key: %v", err)
	}
}

func TestApply_UnknownType_ReturnsError(t *testing.T) {
	tr := NewTransformer([]Rule{{Key: "API_KEY", Type: "rot13"}})
	_, err := tr.Apply(baseSecrets)
	if err == nil {
		t.Fatal("expected error for unknown transform type")
	}
}

func TestParseRules_Valid(t *testing.T) {
	rules, err := ParseRules([]string{"DB_PASSWORD:trim", "API_KEY:prefix:Bearer "})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(rules) != 2 {
		t.Fatalf("expected 2 rules, got %d", len(rules))
	}
	if rules[1].Arg != "Bearer " {
		t.Errorf("expected arg 'Bearer ', got %q", rules[1].Arg)
	}
}

func TestParseRules_Invalid(t *testing.T) {
	_, err := ParseRules([]string{"BADFORMAT"})
	if err == nil {
		t.Fatal("expected error for invalid rule format")
	}
}
