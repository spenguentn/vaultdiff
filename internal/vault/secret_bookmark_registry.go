package vault

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

func bookmarkKey(mount, path string) string {
	return fmt.Sprintf("%s/%s", mount, path)
}

// SecretBookmarkRegistry stores and retrieves secret bookmarks by alias.
type SecretBookmarkRegistry struct {
	mu        sync.RWMutex
	byAlias   map[string]SecretBookmark
	byPath    map[string]SecretBookmark
}

// NewSecretBookmarkRegistry returns an initialized SecretBookmarkRegistry.
func NewSecretBookmarkRegistry() *SecretBookmarkRegistry {
	return &SecretBookmarkRegistry{
		byAlias: make(map[string]SecretBookmark),
		byPath:  make(map[string]SecretBookmark),
	}
}

// Add registers a new bookmark. Returns an error if the alias is already in use or the bookmark is invalid.
func (r *SecretBookmarkRegistry) Add(b SecretBookmark) error {
	if err := b.Validate(); err != nil {
		return err
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.byAlias[b.Alias]; exists {
		return fmt.Errorf("bookmark: alias %q already registered", b.Alias)
	}
	if b.CreatedAt.IsZero() {
		b.CreatedAt = time.Now().UTC()
	}
	r.byAlias[b.Alias] = b
	r.byPath[bookmarkKey(b.Mount, b.Path)] = b
	return nil
}

// GetByAlias retrieves a bookmark by its alias.
func (r *SecretBookmarkRegistry) GetByAlias(alias string) (SecretBookmark, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	b, ok := r.byAlias[alias]
	if !ok {
		return SecretBookmark{}, errors.New("bookmark: alias not found")
	}
	return b, nil
}

// GetByPath retrieves a bookmark by mount and path.
func (r *SecretBookmarkRegistry) GetByPath(mount, path string) (SecretBookmark, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	b, ok := r.byPath[bookmarkKey(mount, path)]
	if !ok {
		return SecretBookmark{}, errors.New("bookmark: path not found")
	}
	return b, nil
}

// Remove deletes a bookmark by alias.
func (r *SecretBookmarkRegistry) Remove(alias string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	b, ok := r.byAlias[alias]
	if !ok {
		return errors.New("bookmark: alias not found")
	}
	delete(r.byAlias, alias)
	delete(r.byPath, bookmarkKey(b.Mount, b.Path))
	return nil
}

// All returns a slice of all registered bookmarks.
func (r *SecretBookmarkRegistry) All() []SecretBookmark {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]SecretBookmark, 0, len(r.byAlias))
	for _, b := range r.byAlias {
		out = append(out, b)
	}
	return out
}
