package cache_test

import (
	"testing"
	"time"

	"github.com/yourusername/vaultpipe/internal/cache"
)

func TestGet_MissOnEmpty(t *testing.T) {
	c := cache.New(time.Minute)
	_, err := c.Get("/secret/data/app")
	if err != cache.ErrCacheMiss {
		t.Fatalf("expected ErrCacheMiss, got %v", err)
	}
}

func TestSetAndGet_ReturnsSecrets(t *testing.T) {
	c := cache.New(time.Minute)
	secrets := map[string]string{"DB_PASS": "hunter2", "API_KEY": "abc123"}
	c.Set("/secret/data/app", secrets)

	got, err := c.Get("/secret/data/app")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["DB_PASS"] != "hunter2" || got["API_KEY"] != "abc123" {
		t.Errorf("unexpected values: %v", got)
	}
}

func TestGet_ReturnsCopy(t *testing.T) {
	c := cache.New(time.Minute)
	secrets := map[string]string{"KEY": "val"}
	c.Set("/secret/data/app", secrets)

	got, _ := c.Get("/secret/data/app")
	got["KEY"] = "mutated"

	again, _ := c.Get("/secret/data/app")
	if again["KEY"] != "val" {
		t.Error("cache returned reference instead of copy")
	}
}

func TestGet_MissAfterExpiry(t *testing.T) {
	c := cache.New(10 * time.Millisecond)
	c.Set("/secret/data/app", map[string]string{"K": "V"})

	time.Sleep(20 * time.Millisecond)

	_, err := c.Get("/secret/data/app")
	if err != cache.ErrCacheMiss {
		t.Fatalf("expected ErrCacheMiss after expiry, got %v", err)
	}
}

func TestInvalidate_RemovesEntry(t *testing.T) {
	c := cache.New(time.Minute)
	c.Set("/secret/data/app", map[string]string{"K": "V"})
	c.Invalidate("/secret/data/app")

	_, err := c.Get("/secret/data/app")
	if err != cache.ErrCacheMiss {
		t.Fatalf("expected ErrCacheMiss after invalidation, got %v", err)
	}
}

func TestFlush_RemovesAllEntries(t *testing.T) {
	c := cache.New(time.Minute)
	c.Set("/secret/data/a", map[string]string{"K": "V"})
	c.Set("/secret/data/b", map[string]string{"X": "Y"})
	c.Flush()

	for _, path := range []string{"/secret/data/a", "/secret/data/b"} {
		if _, err := c.Get(path); err != cache.ErrCacheMiss {
			t.Errorf("expected miss for %s after flush", path)
		}
	}
}

func TestZeroTTL_AlwaysMisses(t *testing.T) {
	c := cache.New(0)
	c.Set("/secret/data/app", map[string]string{"K": "V"})

	_, err := c.Get("/secret/data/app")
	if err != cache.ErrCacheMiss {
		t.Fatalf("expected ErrCacheMiss with zero TTL, got %v", err)
	}
}
