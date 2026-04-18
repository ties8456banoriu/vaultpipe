// Package secretexpand implements variable interpolation for secret values.
//
// Secrets may reference other secrets using ${KEY} syntax. References are
// resolved recursively. Circular references and unknown keys are reported
// as errors.
//
// Example:
//
//	HOST=localhost
//	DSN=postgres://${HOST}/mydb  -> postgres://localhost/mydb
package secretexpand
