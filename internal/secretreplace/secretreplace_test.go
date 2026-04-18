package secretreplace_test

import (
	"testing"

	"github.com/yourusername/vaultpipe/internal/secretreplace"
)

func TestNew_EmptyFind_ReturnsError(t *testing.T) {
	_, err := secretreplace.New([]secretreplace.Rule{{Find: "", Replace: "x"}})
	if err == nil {
		t.Fatal("expected error for empty Find")
	}
}

func TestParseRules_Valid(t *testing.T) {
	rules, err := secretreplace.ParseRules([]string{"foo=bar", "baz="})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(rules) != 2 {
		t.Fatalf("expected 2 rules, got %d", len(rules))
	}
	if rules[0].Find != "foo" || rules[0].Replace != "bar" {
		t.Errorf("unexpected rule[0]: %+v", rules[0])
	}
	if rules[1].Find != "baz" || rules[1].Replace != "" {
		t.Errorf("unexpected rule[1]: %+v", rules[1])
	}
}

func TestParseRules_Invalid(t *testing.T) {
	_, err := secretreplace.ParseRules([]string{"noequalssign"})
	if err == nil {
		t.Fatal("expected error for missing '='")
	}
}

func TestApply_EmptySecrets_ReturnsError(t *testing.T) {
	r, _ := secretreplace.New(nil)
	_, err := r.Apply(map[string]string{})
	if err == nil {
		t.Fatal("expected error for empty secrets")
	}
}

func TestApply_NoRules_ReturnsUnchanged(t *testing.T) {
	r, _ := secretreplace.New(nil)
	secrets := map[string]string{"KEY": "hello world"}
	out, err := r.Apply(secrets)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["KEY"] != "hello world" {
		t.Errorf("expected unchanged value, got %q", out["KEY"])
	}
}

func TestApply_SingleRule_ReplacesValue(t *testing.T) {
	r, _ := secretreplace.New([]secretreplace.Rule{{Find: "prod", Replace: "dev"}})
	secrets := map[string]string{"DB_URL": "postgres://prod-host/db"}
	out, err := r.Apply(secrets)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["DB_URL"] != "postgres://dev-host/db" {
		t.Errorf("unexpected value: %q", out["DB_URL"])
	}
}

func TestApply_MultipleRules_AppliedInOrder(t *testing.T) {
	r, _ := secretreplace.New([]secretreplace.Rule{
		{Find: "aaa", Replace: "bbb"},
		{Find: "bbb", Replace: "ccc"},
	})
	secrets := map[string]string{"VAL": "aaa"}
	out, err := r.Apply(secrets)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["VAL"] != "ccc" {
		t.Errorf("expected ccc, got %q", out["VAL"])
	}
}

func TestApply_IsolatedFromSource(t *testing.T) {
	r, _ := secretreplace.New([]secretreplace.Rule{{Find: "x", Replace: "y"}})
	src := map[string]string{"K": "x"}
	out, _ := r.Apply(src)
	if src["K"] != "x" {
		t.Error("source map was mutated")
	}
	_ = out
}
