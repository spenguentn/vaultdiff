package vault

import (
	"fmt"
	"sync"
	"time"
)

func riskKey(mount, path string) string {
	return fmt.Sprintf("%s/%s", mount, path)
}

// SecretRiskRegistry stores and retrieves risk assessments for secrets.
type SecretRiskRegistry struct {
	mu      sync.RWMutex
	entries map[string]SecretRisk
}

// NewSecretRiskRegistry returns an initialised SecretRiskRegistry.
func NewSecretRiskRegistry() *SecretRiskRegistry {
	return &SecretRiskRegistry{
		entries: make(map[string]SecretRisk),
	}
}

// Set stores a risk assessment, validating it first and stamping AssessedAt.
func (r *SecretRiskRegistry) Set(risk SecretRisk) error {
	if err := risk.Validate(); err != nil {
		return err
	}
	if risk.AssessedAt.IsZero() {
		risk.AssessedAt = time.Now().UTC()
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	r.entries[riskKey(risk.Mount, risk.Path)] = risk
	return nil
}

// Get retrieves the risk assessment for the given mount and path.
func (r *SecretRiskRegistry) Get(mount, path string) (SecretRisk, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	v, ok := r.entries[riskKey(mount, path)]
	return v, ok
}

// Remove deletes the risk assessment for the given mount and path.
func (r *SecretRiskRegistry) Remove(mount, path string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.entries, riskKey(mount, path))
}

// All returns a copy of all stored risk assessments.
func (r *SecretRiskRegistry) All() []SecretRisk {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]SecretRisk, 0, len(r.entries))
	for _, v := range r.entries {
		out = append(out, v)
	}
	return out
}
