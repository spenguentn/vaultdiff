package vault

import (
	"fmt"
	"sync"
	"time"
)

func impactKey(mount, path string) string {
	return fmt.Sprintf("%s/%s", mount, path)
}

// SecretImpactRegistry stores impact assessments keyed by mount+path.
type SecretImpactRegistry struct {
	mu      sync.RWMutex
	entries map[string]SecretImpact
}

// NewSecretImpactRegistry returns an initialised registry.
func NewSecretImpactRegistry() *SecretImpactRegistry {
	return &SecretImpactRegistry{
		entries: make(map[string]SecretImpact),
	}
}

// Set validates and stores an impact assessment, stamping AssessedBy timestamp
// via the existing field (caller supplies assessed_by string).
func (r *SecretImpactRegistry) Set(entry SecretImpact) error {
	if err := entry.Validate(); err != nil {
		return err
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	// Annotate with current time in justification if empty.
	if entry.Justification == "" {
		entry.Justification = fmt.Sprintf("assessed at %s", time.Now().UTC().Format(time.RFC3339))
	}
	r.entries[impactKey(entry.Mount, entry.Path)] = entry
	return nil
}

// Get retrieves an impact assessment by mount and path.
func (r *SecretImpactRegistry) Get(mount, path string) (SecretImpact, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	e, ok := r.entries[impactKey(mount, path)]
	return e, ok
}

// Remove deletes an impact entry.
func (r *SecretImpactRegistry) Remove(mount, path string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.entries, impactKey(mount, path))
}

// All returns a snapshot of every registered impact entry.
func (r *SecretImpactRegistry) All() []SecretImpact {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]SecretImpact, 0, len(r.entries))
	for _, e := range r.entries {
		out = append(out, e)
	}
	return out
}
