// Package cache implements a lightweight in-memory TTL cache for Vault secret
// maps. It is intended to reduce repeated Vault API calls when vaultpipe
// operates in watch/refresh mode with a short tick interval.
//
// Usage:
//
//	c := cache.New(30 * time.Second)
//	c.Set("/secret/data/myapp", secrets)
//	if vals, err := c.Get("/secret/data/myapp"); err == nil {
//	    // use cached vals
//	}
package cache
