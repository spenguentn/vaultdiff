package vault

import (
	"fmt"
	"sync"
	"time"
)

func endorsementKey(mount, path, endorsedBy string) string {
	return fmt.Sprintf("%s/%s::%s", mount, path, endorsedBy)
}

// SecretEndorsementRegistry stores endorsement records keyed by mount, path, and endorser.
type SecretEndorsementRegistry struct {
	mu      sync.RWMutex
	entries map[string]*SecretEndorsement
}

// NewSecretEndorsementRegistry creates an empty endorsement registry.
func NewSecretEndorsementRegistry() *SecretEndorsementRegistry {
	return &SecretEndorsementRegistry{
		entries: make(map[string]*SecretEndorsement),
	}
}

// Submit records an endorsement after validation, setting EndorsedAt if zero.
func (r *SecretEndorsementRegistry) Submit(e *SecretEndorsement) error {
	if err := e.Validate(); err != nil {
		return err
	}
	if e.EndorsedAt.IsZero() {
		e.EndorsedAt = time.Now().UTC()
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	r.entries[endorsementKey(e.Mount, e.Path, e.EndorsedBy)] = e
	return nil
}

// Get retrieves an endorsement by mount, path, and endorser.
func (r *SecretEndorsementRegistry) Get(mount, path, endorsedBy string) (*SecretEndorsement, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	e, ok := r.entries[endorsementKey(mount, path, endorsedBy)]
	return e, ok
}

// Remove deletes an endorsement record.
func (r *SecretEndorsementRegistry) Remove(mount, path, endorsedBy string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.entries, endorsementKey(mount, path, endorsedBy))
}

// List returns all endorsements for a given mount and path.
func (r *SecretEndorsementRegistry) List(mount, path string) []*SecretEndorsement {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var out []*SecretEndorsement
	for _, e := range r.entries {
		if e.Mount == mount && e.Path == path {
			out = append(out, e)
		}
	}
	return out
}
