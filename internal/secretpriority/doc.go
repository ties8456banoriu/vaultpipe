// Package secretpriority provides priority-based resolution of secret values.
//
// When the same env key is supplied by multiple sources (e.g. a Vault path,
// a local override, and a profile default) the Manager picks the value
// associated with the highest numeric priority level, making precedence
// rules explicit and auditable.
package secretpriority
