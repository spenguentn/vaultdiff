package vault

import (
	"testing"
	"time"
)

var sampleSecrets = map[string]string{
	"username": "admin",
	"password": "s3cr3t",
}

func TestSecretCache_SetAndGet(t *testing.T) {
	c := NewSecretCache(0)
	c.Set("secret/foo", sampleSecrets)

	got, ok := c.Get("secret/foo")
	if !ok {
		t.Fatal("expected cache hit")
	}
	if got["username"] != "admin" {
		t.Errorf("expected admin, got %s", got["username"])
	}
}

func TestSecretCache_MissOnUnknownPath(t *testing.T) {
	c := NewSecretCache(0)
	_, ok := c.Get("secret/missing")
	if ok {
		t.Fatal("expected cache miss")
	}
}

func TestSecretCache_Invalidate(t *testing.T) {
	c := NewSecretCache(0)
	c.Set("secret/bar", sampleSecrets)
	c.Invalidate("secret/bar")

	_, ok := c.Get("secret/bar")
	if ok {
		t.Fatal("expected cache miss after invalidation")
	}
}

func TestSecretCache_Flush(t *testing.T) {
	c := NewSecretCache(0)
	c.Set("secret/a", sampleSecrets)
	c.Set("secret/b", sampleSecrets)
	c.Flush()

	if c.Size() != 0 {
		t.Errorf("expected size 0 after flush, got %d", c.Size())
	}
}

func TestSecretCache_TTL_Expired(t *testing.T) {
	c := NewSecretCache(1 * time.Millisecond)
	c.Set("secret/ttl", sampleSecrets)

	time.Sleep(5 * time.Millisecond)

	_, ok := c.Get("secret/ttl")
	if ok {
		t.Fatal("expected cache miss after TTL expiry")
	}
}

func TestSecretCache_TTL_NotExpired(t *testing.T) {
	c := NewSecretCache(10 * time.Second)
	c.Set("secret/fresh", sampleSecrets)

	_, ok := c.Get("secret/fresh")
	if !ok {
		t.Fatal("expected cache hit before TTL expiry")
	}
}

func TestCacheEntry_Expired_ZeroTTL(t *testing.T) {
	entry := &CacheEntry{
		FetchedAt: time.Now().Add(-24 * time.Hour),
		TTL:       0,
	}
	if entry.Expired() {
		t.Error("zero TTL entry should never expire")
	}
}

func TestSecretCache_Size(t *testing.T) {
	c := NewSecretCache(0)
	if c.Size() != 0 {
		t.Errorf("expected 0, got %d", c.Size())
	}
	c.Set("secret/x", sampleSecrets)
	c.Set("secret/y", sampleSecrets)
	if c.Size() != 2 {
		t.Errorf("expected 2, got %d", c.Size())
	}
}
