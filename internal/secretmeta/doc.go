// Package secretmeta tracks provenance metadata for secrets fetched from
// HashiCorp Vault. For each env key it stores the originating Vault path,
// mount, KV version, and the time the secret was last fetched. This allows
// other subsystems (audit, lineage, expiry) to reference authoritative
// source information without re-querying Vault.
package secretmeta
