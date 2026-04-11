package vault

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

// DependencyLink represents a directional dependency between two secrets.
type DependencyLink struct {
	SourceMount string
	SourcePath  string
	TargetMount string
	TargetPath  string
	AddedBy     string
	AddedAt     time.Time
	Note        string
}

// FullSource returns the full source path.
func (d DependencyLink) FullSource() string {
	return fmt.Sprintf("%s/%s", d.SourceMount, d.SourcePath)
}

// FullTarget returns the full target path.
func (d DependencyLink) FullTarget() string {
	return fmt.Sprintf("%s/%s", d.TargetMount, d.TargetPath)
}

// Validate checks that the dependency link is well-formed.
func (d DependencyLink) Validate() error {
	if d.SourceMount == "" {
		return errors.New("source mount is required")
	}
	if d.SourcePath == "" {
		return errors.New("source path is required")
	}
	if d.TargetMount == "" {
		return errors.New("target mount is required")
	}
	if d.TargetPath == "" {
		return errors.New("target path is required")
	}
	if d.AddedBy == "" {
		return errors.New("added_by is required")
	}
	return nil
}

// SecretDependencyRegistry tracks dependencies between secrets.
type SecretDependencyRegistry struct {
	mu    sync.RWMutex
	links map[string][]DependencyLink
}

func dependencyKey(mount, path string) string {
	return fmt.Sprintf("%s/%s", mount, path)
}

// NewSecretDependencyRegistry creates an empty registry.
func NewSecretDependencyRegistry() *SecretDependencyRegistry {
	return &SecretDependencyRegistry{
		links: make(map[string][]DependencyLink),
	}
}

// Add records a dependency link after validation.
func (r *SecretDependencyRegistry) Add(link DependencyLink) error {
	if err := link.Validate(); err != nil {
		return fmt.Errorf("invalid dependency link: %w", err)
	}
	if link.AddedAt.IsZero() {
		link.AddedAt = time.Now().UTC()
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	key := dependencyKey(link.SourceMount, link.SourcePath)
	r.links[key] = append(r.links[key], link)
	return nil
}

// GetDependencies returns all targets that the given source depends on.
func (r *SecretDependencyRegistry) GetDependencies(mount, path string) []DependencyLink {
	r.mu.RLock()
	defer r.mu.RUnlock()
	key := dependencyKey(mount, path)
	out := make([]DependencyLink, len(r.links[key]))
	copy(out, r.links[key])
	return out
}

// Remove deletes all dependency links for the given source.
func (r *SecretDependencyRegistry) Remove(mount, path string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.links, dependencyKey(mount, path))
}

// Count returns the total number of dependency links registered.
func (r *SecretDependencyRegistry) Count() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	total := 0
	for _, v := range r.links {
		total += len(v)
	}
	return total
}
