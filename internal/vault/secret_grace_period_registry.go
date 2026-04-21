package vault

import (
	"fmt"
	"sync"
	"time"
)

func gracePeriodKey(mount, path string) string {
	return fmt.Sprintf("%s/%s", mount, path)
}

// SecretGracePeriodRegistry manages grace-period records for secrets.
type SecretGracePeriodRegistry struct {
	mu      sync.RWMutex
	records map[string]SecretGracePeriod
}

// NewSecretGracePeriodRegistry returns an initialised registry.
func NewSecretGracePeriodRegistry() *SecretGracePeriodRegistry {
	return &SecretGracePeriodRegistry{
		records: make(map[string]SecretGracePeriod),
	}
}

// Set stores a grace-period record, stamping CreatedAt if zero.
func (r *SecretGracePeriodRegistry) Set(g SecretGracePeriod) error {
	if err := g.Validate(); err != nil {
		return err
	}
	if g.CreatedAt.IsZero() {
		g.CreatedAt = time.Now().UTC()
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	r.records[gracePeriodKey(g.Mount, g.Path)] = g
	return nil
}

// Get retrieves the grace-period record for a given mount+path.
func (r *SecretGracePeriodRegistry) Get(mount, path string) (SecretGracePeriod, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	g, ok := r.records[gracePeriodKey(mount, path)]
	return g, ok
}

// Remove deletes the grace-period record for a given mount+path.
func (r *SecretGracePeriodRegistry) Remove(mount, path string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.records, gracePeriodKey(mount, path))
}

// Expired returns all records whose grace period has elapsed.
func (r *SecretGracePeriodRegistry) Expired() []SecretGracePeriod {
	r.mu.RLock()
	defer r.mu.RUnlock()
	now := time.Now().UTC()
	var out []SecretGracePeriod
	for _, g := range r.records {
		if !g.ExpiresAt.IsZero() && now.After(g.ExpiresAt) {
			out = append(out, g)
		}
	}
	return out
}
