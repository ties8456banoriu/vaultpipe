// Package secretpin provides version-pinning for individual secrets fetched
// from HashiCorp Vault. Pinned secrets are skipped during auto-refresh,
// ensuring a specific version remains active in the local environment until
// the pin is explicitly removed.
package secretpin
