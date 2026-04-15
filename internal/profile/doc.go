// Package profile provides named environment profile management for vaultpipe.
//
// A profile groups a Vault secret path with an output .env file under a
// human-readable name (e.g. "dev", "staging", "prod"). Profiles are persisted
// as JSON and can be loaded, listed, switched, and deleted via the CLI.
package profile
