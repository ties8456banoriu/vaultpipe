// Package secretnamespace provides namespace-based isolation and grouping
// for secrets fetched from HashiCorp Vault.
//
// Namespaces allow multiple logical environments (e.g. "prod", "staging", "test")
// to coexist in memory simultaneously, and can be merged on demand.
package secretnamespace
