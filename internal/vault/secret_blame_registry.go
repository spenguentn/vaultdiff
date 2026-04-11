package vault

import (
	"fmt"
	"sync"
)

// SecretBlameRegistry stores blame entries keyed by mount+path+version.
type SecretBlameRegistry struct {
	mu      sync.RWMutex
	entries map[string]BlameEntry
}

// NewSecretBlameRegistry returns an initialised SecretBlameRegistry.
func NewSecretBlameRegistry() *SecretBlameRegistry {
	return &SecretBlameRegistry{
		entries: make(map[string]BlameEntry),
	}
}

func blameKey(mount, path string, version int) string {
	return fmt.Sprintf("%s/%s@v%d", mount, path, version)
}

// Record stores a BlameEntry after validation.
func (r *SecretBlameRegistry) Record(entry BlameEntry) error {
	if err := entry.Validate(); err != nil {
		return err
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	r.entries[blameKey(entry.Mount, entry.Path, entry.Version)] = entry
	return nil
}

// Get retrieves the BlameEntry for a specific secret version.
func (r *SecretBlameRegistry) Get(mount, path string, version int) (BlameEntry, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	e, ok := r.entries[blameKey(mount, path, version)]
	return e, ok
}

// Latest returns the highest-version blame entry for a given mount+path.
func (r *SecretBlameRegistry) Latest(mount, path string) (BlameEntry, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var best BlameEntry
	found := false
	for _, e := range r.entries {
		if e.Mount == mount && e.Path == path {
			if !found || e.Version > best.Version {
				best = e
				found = true
			}
		}
	}
	return best, found
}

// Remove deletes a blame entry for a specific version.
func (r *SecretBlameRegistry) Remove(mount, path string, version int) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.entries, blameKey(mount, path, version))
}
