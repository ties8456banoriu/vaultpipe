package secretencode_test

import (
	"testing"

	"github.com/yourusername/vaultpipe/internal/secretencode"
)

func TestNew_UnsupportedEncoding_ReturnsError(t *testing.T) {
	_, err := secretencode.New("rot13", false)
	if err == nil {
		t.Fatal("expected error for unsupported encoding")
	}
}

func TestNew_ValidEncoding_ReturnsEncoder(t *testing.T) {
	e, err := secretencode.New(secretencode.EncodingBase64, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if e == nil {
		t.Fatal("expected non-nil encoder")
	}
}

func TestApply_EmptySecrets_ReturnsError(t *testing.T) {
	e, _ := secretencode.New(secretencode.EncodingBase64, false)
	_, err := e.Apply(map[string]string{})
	if err == nil {
		t.Fatal("expected error for empty secrets")
	}
}

func TestApply_Base64Encode(t *testing.T) {
	e, _ := secretencode.New(secretencode.EncodingBase64, false)
	out, err := e.Apply(map[string]string{"KEY": "hello"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["KEY"] != "aGVsbG8=" {
		t.Errorf("expected aGVsbG8=, got %q", out["KEY"])
	}
}

func TestApply_Base64Decode(t *testing.T) {
	e, _ := secretencode.New(secretencode.EncodingBase64, true)
	out, err := e.Apply(map[string]string{"KEY": "aGVsbG8="})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["KEY"] != "hello" {
		t.Errorf("expected hello, got %q", out["KEY"])
	}
}

func TestApply_Base64URLEncode(t *testing.T) {
	e, _ := secretencode.New(secretencode.EncodingBase64URL, false)
	out, err := e.Apply(map[string]string{"K": "hello world"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["K"] == "" {
		t.Error("expected non-empty encoded value")
	}
}

func TestApply_InvalidBase64_ReturnsError(t *testing.T) {
	e, _ := secretencode.New(secretencode.EncodingBase64, true)
	_, err := e.Apply(map[string]string{"KEY": "not-valid-base64!!!"})
	if err == nil {
		t.Fatal("expected error for invalid base64")
	}
}

func TestApply_RoundTrip(t *testing.T) {
	original := map[string]string{"SECRET": "my-super-secret-value"}
	enc, _ := secretencode.New(secretencode.EncodingBase64, false)
	dec, _ := secretencode.New(secretencode.EncodingBase64, true)

	encoded, err := enc.Apply(original)
	if err != nil {
		t.Fatalf("encode error: %v", err)
	}
	decoded, err := dec.Apply(encoded)
	if err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if decoded["SECRET"] != original["SECRET"] {
		t.Errorf("round-trip mismatch: got %q", decoded["SECRET"])
	}
}
