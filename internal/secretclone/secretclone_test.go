package secretclone_test

import (
	"testing"

	"github.com/yourusername/vaultpipe/internal/secretclone"
)

func TestClone_EmptySecrets_ReturnsError(t *testing.T) {
	c := secretclone.New()
	_, err := c.Clone(map[string]string{})
	if err == nil {
		t.Fatal("expected error for empty secrets")
	}
}

func TestClone_BasicCopy(t *testing.T) {
	src := map[string]string{"foo": "bar", "baz": "qux"}
	c := secretclone.New()
	out, err := c.Clone(src)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != len(src) {
		t.Fatalf("expected %d keys, got %d", len(src), len(out))
	}
	for k, v := range src {
		if out[k] != v {
			t.Errorf("key %q: expected %q, got %q", k, v, out[k])
		}
	}
}

func TestClone_IsolatedFromSource(t *testing.T) {
	src := map[string]string{"key": "val"}
	c := secretclone.New()
	out, _ := c.Clone(src)
	out["key"] = "changed"
	if src["key"] != "val" {
		t.Error("source map was mutated")
	}
}

func TestClone_WithPrefix(t *testing.T) {
	src := map[string]string{"TOKEN": "abc"}
	c := secretclone.New(secretclone.WithPrefix("APP_"))
	out, err := c.Clone(src)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := out["APP_TOKEN"]; !ok {
		t.Error("expected key APP_TOKEN")
	}
}

func TestClone_WithUppercase(t *testing.T) {
	src := map[string]string{"db_host": "localhost"}
	c := secretclone.New(secretclone.WithUppercase())
	out, err := c.Clone(src)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := out["DB_HOST"]; !ok {
		t.Error("expected key DB_HOST")
	}
}

func TestClone_WithPrefixAndUppercase(t *testing.T) {
	src := map[string]string{"secret": "value"}
	c := secretclone.New(secretclone.WithPrefix("VP_"), secretclone.WithUppercase())
	out, err := c.Clone(src)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["VP_SECRET"] != "value" {
		t.Errorf("expected VP_SECRET=value, got %v", out)
	}
}
