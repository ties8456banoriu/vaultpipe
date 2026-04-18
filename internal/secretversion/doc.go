// Package secretversion tracks version history for secrets fetched from
// HashiCorp Vault. Each time a secret is retrieved, callers record the
// env key, vault path, and version number. The tracker provides access
// to the latest recorded version as well as the full history per key.
package secretversion
