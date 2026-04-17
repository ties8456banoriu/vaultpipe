// Package secretlog provides an in-memory access log for secrets fetched
// from HashiCorp Vault. Each Record call appends a timestamped Entry
// containing the env key, vault path, and optional profile name.
//
// Entries can be retrieved in bulk via All, filtered by key via ForKey,
// or cleared with Clear.
package secretlog
