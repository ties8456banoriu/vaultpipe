// Package env resolves local environment variable overrides for vaultpipe.
//
// During development it is common to override specific secrets fetched from
// Vault with local values (e.g. pointing DB_HOST to localhost instead of a
// remote instance). This package reads variables that match a configurable
// prefix from the process environment and merges them on top of the secrets
// map produced by the Vault client, giving local overrides the highest
// precedence without modifying the Vault data itself.
package env
