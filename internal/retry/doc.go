// Package retry provides a configurable retry mechanism with exponential
// backoff for use when fetching secrets from Vault or performing other
// transient operations that may fail intermittently.
//
// Usage:
//
//	d := retry.New(retry.DefaultConfig())
//	err := d.Do(func() error {
//		return vault.FetchSecrets()
//	})
package retry
