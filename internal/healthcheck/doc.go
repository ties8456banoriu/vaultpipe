// Package healthcheck provides a lightweight liveness probe for HashiCorp Vault.
//
// It calls the /v1/sys/health endpoint and reports whether the server is
// reachable, the HTTP status code returned, and the round-trip latency.
// Status codes below 500 are considered reachable (active, standby, DR, etc.).
package healthcheck
