package snapshot

import (
	"errors"
	"fmt"
)

// FallbackLoader wraps a primary fetch function with snapshot-based fallback.
// If the primary fetch fails and a snapshot exists, the snapshot secrets are
// returned along with a non-nil FallbackUsed sentinel error so callers can
// log or warn accordingly.

// ErrFallbackUsed indicates that secrets were sourced from a local snapshot
// rather than a live Vault fetch.
var ErrFallbackUsed = errors.New("vault unreachable: using cached snapshot")

// FetchFunc is the signature expected for a live secret fetch operation.
type FetchFunc func() (map[string]string, error)

// WithFallback attempts fetchFn and, on failure, loads secrets from snapPath.
// If both fail, the original fetch error is returned.
func WithFallback(fetchFn FetchFunc, snapPath string) (map[string]string, error) {
	secrets, err := fetchFn()
	if err == nil {
		return secrets, nil
	}

	fetchErr := err

	snap, snapErr := Load(snapPath)
	if snapErr != nil {
		if errors.Is(snapErr, ErrNoSnapshot) {
			return nil, fmt.Errorf("fetch failed and no snapshot available: %w", fetchErr)
		}
		return nil, fmt.Errorf("fetch failed (%v) and snapshot load failed: %w", fetchErr, snapErr)
	}

	return snap.Secrets, ErrFallbackUsed
}
