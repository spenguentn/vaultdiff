package vault

import (
	"fmt"
	"sync"
)

// TagRegistry stores and retrieves SecretTag entries keyed by their full path.
type TagRegistry struct {
	mu   sync.RWMutex
	data map[string]SecretTag
}

// NewTagRegistry returns an initialised, empty TagRegistry.
func NewTagRegistry() *TagRegistry {
	return &TagRegistry{
		data: make(map[string]SecretTag),
	}
}

// Register adds or replaces a SecretTag in the registry.
func (r *TagRegistry) Register(tag SecretTag) error {
	if err := tag.Validate(); err != nil {
		return fmt.Errorf("tag_registry: register: %w", err)
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	r.data[tag.FullPath()] = tag
	return nil
}

// Get retrieves a SecretTag by its full path (mount/path).
func (r *TagRegistry) Get(fullPath string) (SecretTag, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	t, ok := r.data[fullPath]
	return t, ok
}

// Delete removes a SecretTag from the registry.
func (r *TagRegistry) Delete(fullPath string) bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	_, ok := r.data[fullPath]
	if ok {
		delete(r.data, fullPath)
	}
	return ok
}

// All returns a snapshot of all registered SecretTags.
func (r *TagRegistry) All() []SecretTag {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]SecretTag, 0, len(r.data))
	for _, t := range r.data {
		out = append(out, t)
	}
	return out
}

// Len returns the number of registered tags.
func (r *TagRegistry) Len() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.data)
}
