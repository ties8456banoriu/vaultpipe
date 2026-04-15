// Package diff compares two snapshots of secret key-value maps and
// surfaces which keys were added, removed, or had their values changed.
//
// It is used by the refresh watcher to decide whether the .env file
// needs to be rewritten and to emit structured audit events that
// describe exactly what changed between polling cycles.
package diff
