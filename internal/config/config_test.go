package config

import (
	"os"
	"testing"
)

func TestLoadFromEnv_Defaults(t *testing.T) {
	os.Unsetenv("VAULT_ADDR")
	os.Unsetenv("VAULT_TOKEN")
	os.Unsetenv("VAULT_NAMESPACE")
	os.Unsetenv("VAULTPIPE_OUTPUT")

	cfg := LoadFromEnv()

	if cfg.VaultAddr != "http://127.0.0.1:8200" {
		t.Errorf("expected default VaultAddr, got %q", cfg.VaultAddr)
	}
	if cfg.OutputFile != ".env" {
		t.Errorf("expected default OutputFile '.env', got %q", cfg.OutputFile)
	}
	if cfg.VaultToken != "" {
		t.Errorf("expected empty VaultToken, got %q", cfg.VaultToken)
	}
}

func TestLoadFromEnv_OverridesFromEnv(t *testing.T) {
	t.Setenv("VAULT_ADDR", "https://vault.example.com")
	t.Setenv("VAULT_TOKEN", "s.testtoken")
	t.Setenv("VAULT_NAMESPACE", "dev")
	t.Setenv("VAULTPIPE_OUTPUT", ".env.local")

	cfg := LoadFromEnv()

	if cfg.VaultAddr != "https://vault.example.com" {
		t.Errorf("unexpected VaultAddr: %q", cfg.VaultAddr)
	}
	if cfg.VaultToken != "s.testtoken" {
		t.Errorf("unexpected VaultToken: %q", cfg.VaultToken)
	}
	if cfg.VaultNamespace != "dev" {
		t.Errorf("unexpected VaultNamespace: %q", cfg.VaultNamespace)
	}
	if cfg.OutputFile != ".env.local" {
		t.Errorf("unexpected OutputFile: %q", cfg.OutputFile)
	}
}

func TestValidate_MissingToken(t *testing.T) {
	cfg := &Config{
		VaultAddr:  "http://127.0.0.1:8200",
		SecretPath: "secret/data/myapp",
	}
	if err := cfg.Validate(); err == nil {
		t.Error("expected validation error for missing VaultToken")
	}
}

func TestValidate_MissingSecretPath(t *testing.T) {
	cfg := &Config{
		VaultAddr:  "http://127.0.0.1:8200",
		VaultToken: "s.testtoken",
	}
	if err := cfg.Validate(); err == nil {
		t.Error("expected validation error for missing SecretPath")
	}
}

func TestValidate_Valid(t *testing.T) {
	cfg := &Config{
		VaultAddr:  "http://127.0.0.1:8200",
		VaultToken: "s.testtoken",
		SecretPath: "secret/data/myapp",
	}
	if err := cfg.Validate(); err != nil {
		t.Errorf("unexpected validation error: %v", err)
	}
}
