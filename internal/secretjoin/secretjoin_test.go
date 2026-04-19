package secretjoin_test

import (
	"testing"

	"github.com/yourusername/vaultpipe/internal/secretjoin"
)

func TestNew_NoRules_ReturnsError(t *testing.T) {
	_, err := secretjoin.New(nil)
	if err == nil {
		t.Fatal("expected error for no rules")
	}
}

func TestNew_EmptyTargetKey_ReturnsError(t *testing.T) {
	_, err := secretjoin.New([]secretjoin.Rule{
		{Keys: []string{"A", "B"}, TargetKey: "", Separator: "-"},
	})
	if err == nil {
		t.Fatal("expected error for empty target key")
	}
}

func TestNew_TooFewSourceKeys_ReturnsError(t *testing.T) {
	_, err := secretjoin.New([]secretjoin.Rule{
		{Keys: []string{"A"}, TargetKey: "OUT", Separator: "-"},
	})
	if err == nil {
		t.Fatal("expected error for fewer than two source keys")
	}
}

func TestApply_EmptySecrets_ReturnsError(t *testing.T) {
	j, _ := secretjoin.New([]secretjoin.Rule{
		{Keys: []string{"A", "B"}, TargetKey: "OUT", Separator: "-"},
	})
	_, err := j.Apply(map[string]string{})
	if err == nil {
		t.Fatal("expected error for empty secrets")
	}
}

func TestApply_MissingSourceKey_ReturnsError(t *testing.T) {
	j, _ := secretjoin.New([]secretjoin.Rule{
		{Keys: []string{"A", "MISSING"}, TargetKey: "OUT", Separator: "-"},
	})
	_, err := j.Apply(map[string]string{"A": "hello"})
	if err == nil {
		t.Fatal("expected error for missing source key")
	}
}

func TestApply_JoinsValues(t *testing.T) {
	j, _ := secretjoin.New([]secretjoin.Rule{
		{Keys: []string{"FIRST", "LAST"}, TargetKey: "FULL_NAME", Separator: " "},
	})
	out, err := j.Apply(map[string]string{"FIRST": "John", "LAST": "Doe"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := out["FULL_NAME"]; got != "John Doe" {
		t.Errorf("expected 'John Doe', got %q", got)
	}
}

func TestApply_PreservesOriginalKeys(t *testing.T) {
	j, _ := secretjoin.New([]secretjoin.Rule{
		{Keys: []string{"A", "B"}, TargetKey: "AB", Separator: ":"},
	})
	out, err := j.Apply(map[string]string{"A": "foo", "B": "bar"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["A"] != "foo" || out["B"] != "bar" {
		t.Error("original keys should be preserved")
	}
	if out["AB"] != "foo:bar" {
		t.Errorf("expected 'foo:bar', got %q", out["AB"])
	}
}
