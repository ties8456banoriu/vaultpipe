package config

import (
	"errors"
	"os"
)

// Config holds the configuration for vaultpipe.
type Config struct {
	// VaultAddr is the address of the Vault server.
	VaultAddr string

	// VaultToken is the token used to authenticate with Vault.
	VaultToken string

	// VaultNamespace is the optional Vault namespace (Enterprise only).
	VaultNamespace string

	// SecretPath is the path in Vault to read secrets from.
	SecretPath string

	// OutputFile is the path to the .env file to write secrets to.
	OutputFile string

	// Overwrite controls whether existing .env entries should be overwritten.
	Overwrite bool
}

// LoadFromEnv populates Config fields from environment variables,
// using sensible defaults where applicable.
func LoadFromEnv() *Config {
	cfg := &Config{
		VaultAddr:      getEnvOrDefault("VAULT_ADDR", "http://127.0.0.1:8200"),
		VaultToken:     os.Getenv("VAULT_TOKEN"),
		VaultNamespace: os.Getenv("VAULT_NAMESPACE"),
		OutputFile:     getEnvOrDefault("VAULTPIPE_OUTPUT", ".env"),
		Overwrite:      false,
	}
	return cfg
}

// Validate checks that required configuration fields are set.
func (c *Config) Validate() error {
	if c.VaultAddr == "" {
		return errors.New("vault address (VAULT_ADDR) must not be empty")
	}
	if c.VaultToken == "" {
		return errors.New("vault token (VAULT_TOKEN) must not be empty")
	}
	if c.SecretPath == "" {
		return errors.New("secret path must not be empty")
	}
	return nil
}

func getEnvOrDefault(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}
