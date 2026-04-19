// Package secretcopy copies secrets from one key to another within a secret
// map. Rules are expressed as FROM:TO pairs. The original keys are preserved
// unless explicitly overwritten by a rule destination.
package secretcopy
