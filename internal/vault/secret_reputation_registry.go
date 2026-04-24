package vault

import (
	"fmt"
	"sync"
	"time"
)

func reputationKey(mount, path string) string {
	return fmt.Sprintf("%s/%s", mount, path)
}

// SecretReputationRegistry stores reputation records keyed by mount+path.
type SecretReputationRegistry struct {
	mu      sync.RWMutex
	records map[string]SecretReputation
}

// NewSecretReputationRegistry creates and returns an empty registry.
func NewSecretReputationRegistry() *SecretReputationRegistry {
	return &SecretReputationRegistry{
		records: make(map[string]SecretReputation),
	}
}

// Set validates and stores a reputation record, stamping AssessedAt if zero.
func (reg *SecretReputationRegistry) Set(r SecretReputation) error {
	if err := r.Validate(); err != nil {
		return err
	}
	if r.AssessedAt.IsZero() {
		r.AssessedAt = time.Now().UTC()
	}
	reg.mu.Lock()
	defer reg.mu.Unlock()
	reg.records[reputationKey(r.Mount, r.Path)] = r
	return nil
}

// Get retrieves a reputation record by mount and path.
func (reg *SecretReputationRegistry) Get(mount, path string) (SecretReputation, bool) {
	reg.mu.RLock()
	defer reg.mu.RUnlock()
	r, ok := reg.records[reputationKey(mount, path)]
	return r, ok
}

// Remove deletes a reputation record by mount and path.
func (reg *SecretReputationRegistry) Remove(mount, path string) {
	reg.mu.Lock()
	defer reg.mu.Unlock()
	delete(reg.records, reputationKey(mount, path))
}

// All returns a snapshot of all stored reputation records.
func (reg *SecretReputationRegistry) All() []SecretReputation {
	reg.mu.RLock()
	defer reg.mu.RUnlock()
	out := make([]SecretReputation, 0, len(reg.records))
	for _, r := range reg.records {
		out = append(out, r)
	}
	return out
}
