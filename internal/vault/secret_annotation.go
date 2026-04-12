package vault

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

// SecretAnnotation holds a key-value annotation attached to a secret path.
type SecretAnnotation struct {
	Mount     string
	Path      string
	Key       string
	Value     string
	CreatedBy string
	CreatedAt time.Time
}

// FullPath returns the canonical mount+path string.
func (a SecretAnnotation) FullPath() string {
	return fmt.Sprintf("%s/%s", a.Mount, a.Path)
}

// Validate checks that required fields are present.
func (a SecretAnnotation) Validate() error {
	if a.Mount == "" {
		return errors.New("annotation: mount is required")
	}
	if a.Path == "" {
		return errors.New("annotation: path is required")
	}
	if a.Key == "" {
		return errors.New("annotation: key is required")
	}
	if a.CreatedBy == "" {
		return errors.New("annotation: created_by is required")
	}
	return nil
}

func annotationKey(mount, path, key string) string {
	return fmt.Sprintf("%s/%s#%s", mount, path, key)
}

// SecretAnnotationRegistry stores annotations keyed by mount/path/key.
type SecretAnnotationRegistry struct {
	mu          sync.RWMutex
	annotations map[string]SecretAnnotation
}

// NewSecretAnnotationRegistry returns an initialised registry.
func NewSecretAnnotationRegistry() *SecretAnnotationRegistry {
	return &SecretAnnotationRegistry{
		annotations: make(map[string]SecretAnnotation),
	}
}

// Set validates and stores an annotation, stamping CreatedAt if zero.
func (r *SecretAnnotationRegistry) Set(a SecretAnnotation) error {
	if err := a.Validate(); err != nil {
		return err
	}
	if a.CreatedAt.IsZero() {
		a.CreatedAt = time.Now().UTC()
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	r.annotations[annotationKey(a.Mount, a.Path, a.Key)] = a
	return nil
}

// Get retrieves an annotation by mount, path and key.
func (r *SecretAnnotationRegistry) Get(mount, path, key string) (SecretAnnotation, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	a, ok := r.annotations[annotationKey(mount, path, key)]
	return a, ok
}

// Remove deletes an annotation entry.
func (r *SecretAnnotationRegistry) Remove(mount, path, key string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.annotations, annotationKey(mount, path, key))
}

// ListForPath returns all annotations stored under a given mount/path.
func (r *SecretAnnotationRegistry) ListForPath(mount, path string) []SecretAnnotation {
	prefix := fmt.Sprintf("%s/%s#", mount, path)
	r.mu.RLock()
	defer r.mu.RUnlock()
	var out []SecretAnnotation
	for k, v := range r.annotations {
		if len(k) >= len(prefix) && k[:len(prefix)] == prefix {
			out = append(out, v)
		}
	}
	return out
}
