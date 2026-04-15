// Package snapshot handles persisting and restoring secret snapshots to disk.
//
// A snapshot captures the full set of secrets at a point in time and writes
// them to a JSON file with restricted permissions (0600). Snapshots can be
// used for offline fallback when Vault is unreachable, or to detect drift
// between successive fetches via the diff package.
//
// Usage:
//
//	err := snapshot.Store("/tmp/vaultpipe.snap", secrets)
//	snap, err := snapshot.Load("/tmp/vaultpipe.snap")
package snapshot
