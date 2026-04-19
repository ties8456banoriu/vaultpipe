package secretflatten_test

import (
	"testing"

	"github.com/yourusername/vaultpipe/internal/secretflatten"
)

func TestNew_EmptyDelimiter_ReturnsError(t *testing.T) {
	_, err := secretflatten.New("", "_")
	if err == nil {
		t.Fatal("expected error for empty delimiter")
	}
}

func TestNew_EmptySeparator_ReturnsError(t *testing.T) {
	_, err := secretflatten.New(".", "")
	if err == nil {
		t.Fatal("expected error for empty separator")
	}
}

func TestNew_Valid_ReturnsFlattener(t *testing.T) {
	f, err := secretflatten.New(".", "_")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f == nil {
		t.Fatal("expected non-nil flattener")
	}
}

func TestApply_EmptySecrets_ReturnsError(t *testing.T) {
	f, _ := secretflatten.New(".", "_")
	_, err := f.Apply(map[string]string{})
	if err == nil {
		t.Fatal("expected error for empty secrets")
	}
}

func TestApply_ReplacesDelimiter(t *testing.T) {
	f, _ := secretflatten.New(".", "_")
	secrets := map[string]string{
		"db.host": "localhost",
		"db.port": "5432",
		"plain":   "value",
	}
	out, err := f.Apply(secrets)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["db_host"] != "localhost" {
		t.Errorf("expected db_host=localhost, got %q", out["db_host"])
	}
	if out["db_port"] != "5432" {
		t.Errorf("expected db_port=5432, got %q", out["db_port"])
	}
	if out["plain"] != "value" {
		t.Errorf("expected plain=value, got %q", out["plain"])
	}
}

func TestApply_WithUppercase(t *testing.T) {
	f, _ := secretflatten.New(".", "_", secretflatten.WithUppercase())
	secrets := map[string]string{
		"db.host": "localhost",
	}
	out, err := f.Apply(secrets)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["DB_HOST"] != "localhost" {
		t.Errorf("expected DB_HOST=localhost, got %v", out)
	}
}

func TestApply_NoDelimiterInKey_Unchanged(t *testing.T) {
	f, _ := secretflatten.New(".", "_")
	secrets := map[string]string{"MYKEY": "val"}
	out, err := f.Apply(secrets)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["MYKEY"] != "val" {
		t.Errorf("expected MYKEY=val, got %q", out["MYKEY"])
	}
}
