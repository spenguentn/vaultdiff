package vault

import (
	"fmt"
	"sync"
	"time"
)

// SecretLabelRegistry stores and retrieves labels for secrets.
type SecretLabelRegistry struct {
	mu     sync.RWMutex
	labels map[string]SecretLabel
}

// NewSecretLabelRegistry creates an empty label registry.
func NewSecretLabelRegistry() *SecretLabelRegistry {
	return &SecretLabelRegistry{
		labels: make(map[string]SecretLabel),
	}
}

// Set adds or replaces a label on a secret. CreatedAt is set automatically.
func (r *SecretLabelRegistry) Set(l SecretLabel) error {
	if err := l.Validate(); err != nil {
		return err
	}
	if l.CreatedAt.IsZero() {
		l.CreatedAt = time.Now().UTC()
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	r.labels[labelKey(l.Mount, l.Path, l.Key)] = l
	return nil
}

// Get retrieves a label by mount, path and key.
func (r *SecretLabelRegistry) Get(mount, path, key string) (SecretLabel, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	l, ok := r.labels[labelKey(mount, path, key)]
	if !ok {
		return SecretLabel{}, fmt.Errorf("label %q not found for %s/%s", key, mount, path)
	}
	return l, nil
}

// List returns all labels for a given mount and path.
func (r *SecretLabelRegistry) List(mount, path string) []SecretLabel {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var out []SecretLabel
	for _, l := range r.labels {
		if l.Mount == mount && l.Path == path {
			out = append(out, l)
		}
	}
	return out
}

// Remove deletes a label from the registry.
func (r *SecretLabelRegistry) Remove(mount, path, key string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	k := labelKey(mount, path, key)
	if _, ok := r.labels[k]; !ok {
		return fmt.Errorf("label %q not found for %s/%s", key, mount, path)
	}
	delete(r.labels, k)
	return nil
}
