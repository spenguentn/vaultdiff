package vault

import (
	"fmt"
	"sync"
	"time"
)

// SecretAccessLogRegistry stores access log entries keyed by mount+path.
type SecretAccessLogRegistry struct {
	mu      sync.RWMutex
	entries map[string][]SecretAccessEntry
}

// NewSecretAccessLogRegistry creates an empty registry.
func NewSecretAccessLogRegistry() *SecretAccessLogRegistry {
	return &SecretAccessLogRegistry{
		entries: make(map[string][]SecretAccessEntry),
	}
}

func accessLogKey(mount, path string) string {
	return fmt.Sprintf("%s/%s", mount, path)
}

// Record appends a validated access entry to the registry.
func (r *SecretAccessLogRegistry) Record(entry SecretAccessEntry) error {
	if entry.Timestamp.IsZero() {
		entry.Timestamp = time.Now().UTC()
	}
	if err := entry.Validate(); err != nil {
		return err
	}
	key := accessLogKey(entry.Mount, entry.Path)
	r.mu.Lock()
	defer r.mu.Unlock()
	r.entries[key] = append(r.entries[key], entry)
	return nil
}

// Get returns all access entries for a given mount and path.
func (r *SecretAccessLogRegistry) Get(mount, path string) ([]SecretAccessEntry, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	entries, ok := r.entries[accessLogKey(mount, path)]
	return entries, ok
}

// All returns a flat slice of every recorded entry.
func (r *SecretAccessLogRegistry) All() []SecretAccessEntry {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var out []SecretAccessEntry
	for _, list := range r.entries {
		out = append(out, list...)
	}
	return out
}

// Clear removes all entries from the registry.
func (r *SecretAccessLogRegistry) Clear() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.entries = make(map[string][]SecretAccessEntry)
}
