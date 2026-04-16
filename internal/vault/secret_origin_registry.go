package vault

import (
	"fmt"
	"strings"
	"sync"
	"time"
)

func originKey(mount, path string) string {
	return strings.Trim(mount, "/") + "/" + strings.Trim(path, "/")
}

// SecretOriginRegistry stores origin records keyed by mount+path.
type SecretOriginRegistry struct {
	mu      sync.RWMutex
	entries map[string]SecretOrigin
}

// NewSecretOriginRegistry returns an initialised registry.
func NewSecretOriginRegistry() *SecretOriginRegistry {
	return &SecretOriginRegistry{
		entries: make(map[string]SecretOrigin),
	}
}

// Record stores an origin entry, stamping CreatedAt if zero.
func (r *SecretOriginRegistry) Record(o SecretOrigin) error {
	if err := o.Validate(); err != nil {
		return err
	}
	if o.CreatedAt.IsZero() {
		o.CreatedAt = time.Now().UTC()
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	r.entries[originKey(o.Mount, o.Path)] = o
	return nil
}

// Get retrieves the origin for a mount+path combination.
func (r *SecretOriginRegistry) Get(mount, path string) (SecretOrigin, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	o, ok := r.entries[originKey(mount, path)]
	if !ok {
		return SecretOrigin{}, fmt.Errorf("origin: no record for %s", originKey(mount, path))
	}
	return o, nil
}

// Remove deletes an origin record.
func (r *SecretOriginRegistry) Remove(mount, path string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.entries, originKey(mount, path))
}

// All returns a snapshot of all recorded origins.
func (r *SecretOriginRegistry) All() []SecretOrigin {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]SecretOrigin, 0, len(r.entries))
	for _, v := range r.entries {
		out = append(out, v)
	}
	return out
}
