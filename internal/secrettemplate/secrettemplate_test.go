package secrettemplate_test

import (
	"strings"
	"testing"

	"github.com/yourusername/vaultpipe/internal/secrettemplate"
)

func TestRender_EmptySecrets_ReturnsError(t *testing.T) {
	r := secrettemplate.New()
	_, err := r.Render("hello", nil)
	if err == nil {
		t.Fatal("expected error for empty secrets")
	}
}

func TestRender_EmptyTemplate_ReturnsError(t *testing.T) {
	r := secrettemplate.New()
	_, err := r.Render("", map[string]string{"KEY": "val"})
	if err == nil {
		t.Fatal("expected error for empty template")
	}
}

func TestRender_SimpleSubstitution(t *testing.T) {
	r := secrettemplate.New()
	secrets := map[string]string{"DB_HOST": "localhost"}
	out, err := r.Render(`host={{ index . "DB_HOST" }}`, secrets)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != "host=localhost" {
		t.Errorf("got %q, want %q", out, "host=localhost")
	}
}

func TestRender_MissingKey_ReturnsError(t *testing.T) {
	r := secrettemplate.New()
	secrets := map[string]string{"OTHER": "val"}
	_, err := r.Render(`{{ index . "MISSING" }}`, secrets)
	if err == nil {
		t.Fatal("expected error for missing key")
	}
}

func TestRender_InvalidTemplate_ReturnsError(t *testing.T) {
	r := secrettemplate.New()
	secrets := map[string]string{"K": "v"}
	_, err := r.Render(`{{ .Unclosed`, secrets)
	if err == nil {
		t.Fatal("expected parse error")
	}
}

func TestRenderAll_EmptyTemplates_ReturnsError(t *testing.T) {
	r := secrettemplate.New()
	_, err := r.RenderAll(nil, map[string]string{"K": "v"})
	if err == nil {
		t.Fatal("expected error for empty templates")
	}
}

func TestRenderAll_MultipleTemplates(t *testing.T) {
	r := secrettemplate.New()
	secrets := map[string]string{"USER": "admin", "PASS": "secret"}
	templates := map[string]string{
		"dsn":  `user={{ index . "USER" }} pass={{ index . "PASS" }}`,
		"user": `{{ index . "USER" }}`,
	}
	out, err := r.RenderAll(templates, secrets)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out["dsn"], "admin") {
		t.Errorf("dsn missing user, got %q", out["dsn"])
	}
	if out["user"] != "admin" {
		t.Errorf("user got %q, want admin", out["user"])
	}
}
