// Package validate provides pre-write validation of secret key-value pairs.
//
// It checks that every key conforms to the POSIX environment variable naming
// convention ([A-Za-z_][A-Za-z0-9_]*) and optionally warns when a value is
// an empty string. Validation runs before secrets are written to the .env
// file so that invalid data is caught early and surfaced to the operator.
package validate
