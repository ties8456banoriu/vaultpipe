// Package cache provides an in-memory TTL cache for Vault secrets,
// reducing redundant requests during short-lived refresh cycles.
package cache

import (
	"errors"
	"sync"
	"time"
)

// ErrCacheMiss is returned when a requested key is not present or has expired.
var ErrCacheMiss = errors.New("cache: miss")

// entry holds a cached secret map alongside its expiry time.
type entry struct {
	secrets   map[string]string
	expiresAt time.Time
}

// Cache is a thread-safe in-memory store for secret maps keyed by Vault path.
type Cache struct {
	mu  sync.RWMutex
	ttl time.Duration
	data map[string]entry
}

// New creates a Cache with the given TTL. A TTL of zero disables caching
// (every Get returns ErrCacheMiss).
func New(ttl time.Duration) *Cache {
	return &Cache{
		ttl:  ttl,
		data: make(map[string]entry),
	}
}

// Set stores a copy of secrets for the given path, expiring after the cache TTL.
func (c *Cache) Set(path string, secrets map[string]string) {
	if c.ttl == 0 {
		return
	}
	copy := make(map[string]string, len(secrets))
	for k, v := range secrets {
		copy[k] = v
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data[path] = entry{
		secrets:   copy,
		expiresAt: time.Now().Add(c.ttl),
	}
}

// Get retrieves secrets for the given path. Returns ErrCacheMiss if the entry
// is absent or expired.
func (c *Cache) Get(path string) (map[string]string, error) {
	if c.ttl == 0 {
		return nil, ErrCacheMiss
	}
	c.mu.RLock()
	defer c.mu.RUnlock()
	e, ok := c.data[path]
	if !ok || time.Now().After(e.expiresAt) {
		return nil, ErrCacheMiss
	}
	copy := make(map[string]string, len(e.secrets))
	for k, v := range e.secrets {
		copy[k] = v
	}
	return copy, nil
}

// Invalidate removes the cached entry for the given path.
func (c *Cache) Invalidate(path string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.data, path)
}

// Flush removes all cached entries.
func (c *Cache) Flush() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data = make(map[string]entry)
}
