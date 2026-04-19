// Package secrettrim provides substring trimming of secret values
// using per-key start/end index rules.
//
// This is useful when a Vault secret contains a longer string and only
// a slice of it is needed for a particular environment variable.
package secrettrim
