// Package secretflatten provides utilities for flattening nested secret keys
// by replacing a delimiter character with a separator string.
//
// This is useful when secrets sourced from Vault contain hierarchical key names
// (e.g. "db.host") that need to be normalised into env-friendly keys
// (e.g. "DB_HOST") before being written to a .env file.
package secretflatten
