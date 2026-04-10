package vault

import (
	"sync"
	"time"
)

// CacheEntry holds a cached secret value with an expiration timestamp.
type CacheEntry struct {
	Secrets   map[string]string
	FetchedAt time.Time
	TTL       time.Duration
}

// Expired reports whether the cache entry has passed its TTL.
func (e *CacheEntry) Expired() bool {
	if e.TTL <= 0 {
		return false
	}
	return time.Since(e.FetchedAt) > e.TTL
}

// SecretCache is a thread-safe in-memory cache for secret data keyed by path.
type SecretCache struct {
	mu      sync.RWMutex
	entries map[string]*CacheEntry
	ttl     time.Duration
}

// NewSecretCache creates a SecretCache with the given TTL.
// A zero TTL means entries never expire.
func NewSecretCache(ttl time.Duration) *SecretCache {
	return &SecretCache{
		entries: make(map[string]*CacheEntry),
		ttl:     ttl,
	}
}

// Set stores secret data for the given path.
func (c *SecretCache) Set(path string, secrets map[string]string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.entries[path] = &CacheEntry{
		Secrets:   secrets,
		FetchedAt: time.Now(),
		TTL:       c.ttl,
	}
}

// Get retrieves cached secrets for the given path.
// Returns nil, false if the entry is missing or expired.
func (c *SecretCache) Get(path string) (map[string]string, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	entry, ok := c.entries[path]
	if !ok || entry.Expired() {
		return nil, false
	}
	return entry.Secrets, true
}

// Invalidate removes the cache entry for the given path.
func (c *SecretCache) Invalidate(path string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.entries, path)
}

// Flush removes all entries from the cache.
func (c *SecretCache) Flush() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.entries = make(map[string]*CacheEntry)
}

// Size returns the number of entries currently in the cache.
func (c *SecretCache) Size() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.entries)
}
