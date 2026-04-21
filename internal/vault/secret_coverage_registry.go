package vault

import (
	"fmt"
	"sync"
	"time"
)

func coverageKey(mount, path string) string {
	return fmt.Sprintf("%s/%s", mount, path)
}

// SecretCoverageRegistry stores and retrieves SecretCoverage records.
type SecretCoverageRegistry struct {
	mu      sync.RWMutex
	records map[string]SecretCoverage
}

// NewSecretCoverageRegistry returns an initialised SecretCoverageRegistry.
func NewSecretCoverageRegistry() *SecretCoverageRegistry {
	return &SecretCoverageRegistry{
		records: make(map[string]SecretCoverage),
	}
}

// Set stores a coverage record after validation, stamping AssessedAt if zero.
func (r *SecretCoverageRegistry) Set(c SecretCoverage) error {
	if c.AssessedAt.IsZero() {
		c.AssessedAt = time.Now().UTC()
	}
	if err := c.Validate(); err != nil {
		return err
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	r.records[coverageKey(c.Mount, c.Path)] = c
	return nil
}

// Get retrieves the coverage record for the given mount and path.
func (r *SecretCoverageRegistry) Get(mount, path string) (SecretCoverage, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	c, ok := r.records[coverageKey(mount, path)]
	return c, ok
}

// Remove deletes the coverage record for the given mount and path.
func (r *SecretCoverageRegistry) Remove(mount, path string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.records, coverageKey(mount, path))
}

// All returns a snapshot of every stored coverage record.
func (r *SecretCoverageRegistry) All() []SecretCoverage {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]SecretCoverage, 0, len(r.records))
	for _, v := range r.records {
		out = append(out, v)
	}
	return out
}
