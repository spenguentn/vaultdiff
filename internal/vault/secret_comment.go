package vault

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

// SecretComment represents a human-readable annotation attached to a secret.
type SecretComment struct {
	Mount     string    `json:"mount"`
	Path      string    `json:"path"`
	Author    string    `json:"author"`
	Body      string    `json:"body"`
	CreatedAt time.Time `json:"created_at"`
}

// FullPath returns the canonical mount+path identifier.
func (c SecretComment) FullPath() string {
	return fmt.Sprintf("%s/%s", c.Mount, c.Path)
}

// Validate returns an error if the comment is missing required fields.
func (c SecretComment) Validate() error {
	if c.Mount == "" {
		return errors.New("comment: mount is required")
	}
	if c.Path == "" {
		return errors.New("comment: path is required")
	}
	if c.Author == "" {
		return errors.New("comment: author is required")
	}
	if c.Body == "" {
		return errors.New("comment: body is required")
	}
	return nil
}

// SecretCommentRegistry stores comments keyed by mount+path.
type SecretCommentRegistry struct {
	mu       sync.RWMutex
	comments map[string][]SecretComment
}

// NewSecretCommentRegistry creates an empty registry.
func NewSecretCommentRegistry() *SecretCommentRegistry {
	return &SecretCommentRegistry{
		comments: make(map[string][]SecretComment),
	}
}

func commentKey(mount, path string) string {
	return fmt.Sprintf("%s/%s", mount, path)
}

// Add validates and appends a comment to the registry.
func (r *SecretCommentRegistry) Add(c SecretComment) error {
	if err := c.Validate(); err != nil {
		return err
	}
	if c.CreatedAt.IsZero() {
		c.CreatedAt = time.Now().UTC()
	}
	key := commentKey(c.Mount, c.Path)
	r.mu.Lock()
	defer r.mu.Unlock()
	r.comments[key] = append(r.comments[key], c)
	return nil
}

// Get returns all comments for a given mount and path.
func (r *SecretCommentRegistry) Get(mount, path string) ([]SecretComment, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	list, ok := r.comments[commentKey(mount, path)]
	return list, ok
}

// Remove deletes all comments for a given mount and path.
func (r *SecretCommentRegistry) Remove(mount, path string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.comments, commentKey(mount, path))
}

// Count returns the total number of comments across all secrets.
func (r *SecretCommentRegistry) Count() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	total := 0
	for _, list := range r.comments {
		total += len(list)
	}
	return total
}
