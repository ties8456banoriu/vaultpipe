// Package mapping provides utilities for transforming Vault secret keys into
// environment variable names suitable for use in .env files.
//
// By default, keys are converted to UPPER_SNAKE_CASE. Custom mappings can be
// provided via Rule structs or parsed from "vault_key=ENV_KEY" strings, which
// can be supplied through CLI flags or a configuration file.
package mapping
