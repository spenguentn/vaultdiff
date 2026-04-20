package vault

import (
	"fmt"
	"sync"
	"time"
)

func dispositionKey(mount, path string) string {
	return fmt.Sprintf("%s/%s", mount, path)
}

// SecretDispositionRegistry stores and retrieves SecretDisposition records.
type SecretDispositionRegistry struct {
	mu      sync.RWMutex
	records map[string]*SecretDisposition
}

// NewSecretDispositionRegistry returns an initialised registry.
func NewSecretDispositionRegistry() *SecretDispositionRegistry {
	return &SecretDispositionRegistry{
		records: make(map[string]*SecretDisposition),
	}
}

// Set validates and stores a disposition record, stamping CreatedAt if unset.
func (r *SecretDispositionRegistry) Set(d *SecretDisposition) error {
	if err := d.Validate(); err != nil {
		return err
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if d.CreatedAt.IsZero() {
		d.CreatedAt = time.Now().UTC()
	}
	r.records[dispositionKey(d.Mount, d.Path)] = d
	return nil
}

// Get retrieves the disposition for the given mount and path.
func (r *SecretDispositionRegistry) Get(mount, path string) (*SecretDisposition, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	v, ok := r.records[dispositionKey(mount, path)]
	return v, ok
}

// Remove deletes a disposition record.
func (r *SecretDispositionRegistry) Remove(mount, path string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.records, dispositionKey(mount, path))
}

// Due returns all dispositions that are due as of the provided time.
func (r *SecretDispositionRegistry) Due(now time.Time) []*SecretDisposition {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var out []*SecretDisposition
	for _, d := range r.records {
		if d.IsDue(now) {
			out = append(out, d)
		}
	}
	return out
}
