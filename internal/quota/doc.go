// Package quota provides per-key fetch quota enforcement for vaultpipe.
//
// An Enforcer tracks how many times each secret key has been fetched within a
// sliding time window and rejects requests that exceed the configured limit.
//
// Usage:
//
//	e, _ := quota.New(time.Minute)
//	e.SetLimit("DB_PASSWORD", 5)
//	if err := e.Check("DB_PASSWORD"); err != nil {
//		// handle quota exceeded
//	}
package quota
