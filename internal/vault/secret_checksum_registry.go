package vault

import (
	"fmt"
	"strings"
	"sync"
)

func checksumKey(mount, path string) string {
	return strings.Trim(mount, "/") + "/" + strings.Trim(path, "/")
}

// SecretChecksumRegistry stores computed checksums for secrets.
type SecretChecksumRegistry struct {
	mu      sync.RWMutex
	entries map[string]SecretChecksum
}

// NewSecretChecksumRegistry returns an initialised registry.
func NewSecretChecksumRegistry() *SecretChecksumRegistry {
	return &SecretChecksumRegistry{
		entries: make(map[string]SecretChecksum),
	}
}

// Store saves a checksum entry, replacing any existing record for the same path.
func (r *SecretChecksumRegistry) Store(c SecretChecksum) error {
	if err := c.Validate(); err != nil {
		return err
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	r.entries[checksumKey(c.Mount, c.Path)] = c
	return nil
}

// Get retrieves the stored checksum for the given mount and path.
func (r *SecretChecksumRegistry) Get(mount, path string) (SecretChecksum, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	c, ok := r.entries[checksumKey(mount, path)]
	if !ok {
		return SecretChecksum{}, fmt.Errorf("checksum registry: no entry for %s/%s", mount, path)
	}
	return c, nil
}

// Remove deletes the checksum entry for the given mount and path.
func (r *SecretChecksumRegistry) Remove(mount, path string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.entries, checksumKey(mount, path))
}

// All returns a snapshot of all stored checksums.
func (r *SecretChecksumRegistry) All() []SecretChecksum {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]SecretChecksum, 0, len(r.entries))
	for _, c := range r.entries {
		out = append(out, c)
	}
	return out
}

// Len returns the number of stored entries.
func (r *SecretChecksumRegistry) Len() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.entries)
}
