package secrethash_test

import (
	"strings"
	"testing"

	"github.com/yourusername/vaultpipe/internal/secrethash"
)

func TestNew_UnsupportedAlgorithm_ReturnsError(t *testing.T) {
	_, err := secrethash.New("bcrypt")
	if err == nil {
		t.Fatal("expected error for unsupported algorithm, got nil")
	}
}

func TestNew_ValidAlgorithm_ReturnsHasher(t *testing.T) {
	for _, algo := range []secrethash.Algorithm{
		secrethas.AlgoSHA256,
		secrethas.AlgoSHA512,
		secrethas.AlgoMD5,
	} {
		h, err := secrethash.New(algo)
		if err != nil {
			t.Fatalf("unexpected error for algo %q: %v", algo, err)
		}
		if h == nil {
			t.Fatalf("expected non-nil hasher for algo %q", algo)
		}
	}
}

func TestApply_EmptySecrets_ReturnsError(t *testing.T) {
	h, _ := secrethash.New(secrethash.AlgoSHA256)
	_, err := h.Apply(map[string]string{})
	if err == nil {
		t.Fatal("expected error for empty secrets, got nil")
	}
}

func TestApply_SHA256_ProducesHex(t *testing.T) {
	h, _ := secrethash.New(secrethash.AlgoSHA256)
	secrets := map[string]string{"DB_PASS": "hunter2"}
	out, err := h.Apply(secrets)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	val, ok := out["DB_PASS"]
	if !ok {
		t.Fatal("expected key DB_PASS in output")
	}
	if len(val) != 64 {
		t.Fatalf("expected 64-char SHA256 hex, got %d chars: %s", len(val), val)
	}
}

func TestApply_MD5_ProducesHex(t *testing.T) {
	h, _ := secrethash.New(secrethash.AlgoMD5)
	secrets := map[string]string{"TOKEN": "abc123"}
	out, err := h.Apply(secrets)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out["TOKEN"]) != 32 {
		t.Fatalf("expected 32-char MD5 hex, got %d", len(out["TOKEN"]))
	}
}

func TestApply_Deterministic(t *testing.T) {
	h, _ := secrethash.New(secrethash.AlgoSHA256)
	secrets := map[string]string{"KEY": "stable-value"}
	out1, _ := h.Apply(secrets)
	out2, _ := h.Apply(secrets)
	if out1["KEY"] != out2["KEY"] {
		t.Fatal("expected deterministic output, got different hashes")
	}
}

func TestHash_SingleValue_SHA512(t *testing.T) {
	h, _ := secrethash.New(secrethash.AlgoSHA512)
	result, err := h.Hash("mysecret")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 128 {
		t.Fatalf("expected 128-char SHA512 hex, got %d", len(result))
	}
	if !isHex(result) {
		t.Fatalf("result is not valid hex: %s", result)
	}
}

func isHex(s string) bool {
	for _, c := range strings.ToLower(s) {
		if !((c >= '0' && c <= '9') || (c >= 'a' && c <= 'f')) {
			return false
		}
	}
	return true
}
