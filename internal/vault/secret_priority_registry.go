package vault

import (
	"fmt"
	"sync"
	"time"
)

func priorityKey(mount, path string) string {
	return mount + "|" + path
}

// SecretPriorityRegistry stores and retrieves SecretPriority records.
type SecretPriorityRegistry struct {
	mu      sync.RWMutex
	entries map[string]SecretPriority
}

// NewSecretPriorityRegistry returns an initialised registry.
func NewSecretPriorityRegistry() *SecretPriorityRegistry {
	return &SecretPriorityRegistry{
		entries: make(map[string]SecretPriority),
	}
}

// Set validates and stores a priority record, stamping AssignedAt when zero.
func (r *SecretPriorityRegistry) Set(p SecretPriority) error {
	if err := p.Validate(); err != nil {
		return err
	}
	if p.AssignedAt.IsZero() {
		p.AssignedAt = time.Now().UTC()
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	r.entries[priorityKey(p.Mount, p.Path)] = p
	return nil
}

// Get retrieves the priority for the given mount and path.
func (r *SecretPriorityRegistry) Get(mount, path string) (SecretPriority, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	p, ok := r.entries[priorityKey(mount, path)]
	if !ok {
		return SecretPriority{}, fmt.Errorf("secret priority: no entry for %s/%s", mount, path)
	}
	return p, nil
}

// Remove deletes the priority record for the given mount and path.
func (r *SecretPriorityRegistry) Remove(mount, path string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.entries, priorityKey(mount, path))
}

// All returns a snapshot of every stored priority record.
func (r *SecretPriorityRegistry) All() []SecretPriority {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]SecretPriority, 0, len(r.entries))
	for _, p := range r.entries {
		out = append(out, p)
	}
	return out
}
