package vault

import (
	"fmt"
	"sync"
	"time"
)

func deprecationKey(mount, path string) string {
	return fmt.Sprintf("%s/%s", mount, path)
}

// SecretDeprecationRegistry tracks deprecation records for secrets.
type SecretDeprecationRegistry struct {
	mu      sync.RWMutex
	records map[string]SecretDeprecation
}

// NewSecretDeprecationRegistry returns an initialised registry.
func NewSecretDeprecationRegistry() *SecretDeprecationRegistry {
	return &SecretDeprecationRegistry{
		records: make(map[string]SecretDeprecation),
	}
}

// Set stores a deprecation record, stamping DeprecatedAt if zero.
func (r *SecretDeprecationRegistry) Set(d SecretDeprecation) error {
	if err := d.Validate(); err != nil {
		return err
	}
	if d.DeprecatedAt.IsZero() {
		d.DeprecatedAt = time.Now().UTC()
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	r.records[deprecationKey(d.Mount, d.Path)] = d
	return nil
}

// Get retrieves the deprecation record for a given mount+path.
func (r *SecretDeprecationRegistry) Get(mount, path string) (SecretDeprecation, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	d, ok := r.records[deprecationKey(mount, path)]
	return d, ok
}

// Remove deletes the deprecation record for a given mount+path.
func (r *SecretDeprecationRegistry) Remove(mount, path string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.records, deprecationKey(mount, path))
}

// Sunsetted returns all records whose sunset date has passed.
func (r *SecretDeprecationRegistry) Sunsetted() []SecretDeprecation {
	r.mu.RLock()
	defer r.mu.RUnlock()
	now := time.Now().UTC()
	var out []SecretDeprecation
	for _, d := range r.records {
		if d.IsSunset(now) {
			out = append(out, d)
		}
	}
	return out
}
