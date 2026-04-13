package vault

import "fmt"

func visibilityKey(mount, path string) string {
	return fmt.Sprintf("%s/%s", mount, path)
}

// SecretVisibilityRegistry stores visibility levels keyed by mount+path.
type SecretVisibilityRegistry struct {
	entries map[string]SecretVisibility
}

// NewSecretVisibilityRegistry creates an empty registry.
func NewSecretVisibilityRegistry() *SecretVisibilityRegistry {
	return &SecretVisibilityRegistry{
		entries: make(map[string]SecretVisibility),
	}
}

// Set stores a visibility record after validation.
func (r *SecretVisibilityRegistry) Set(v SecretVisibility) error {
	if err := v.Validate(); err != nil {
		return err
	}
	r.entries[visibilityKey(v.Mount, v.Path)] = v
	return nil
}

// Get retrieves a visibility record by mount and path.
func (r *SecretVisibilityRegistry) Get(mount, path string) (SecretVisibility, bool) {
	v, ok := r.entries[visibilityKey(mount, path)]
	return v, ok
}

// Remove deletes a visibility record.
func (r *SecretVisibilityRegistry) Remove(mount, path string) {
	delete(r.entries, visibilityKey(mount, path))
}

// All returns every stored visibility record.
func (r *SecretVisibilityRegistry) All() []SecretVisibility {
	out := make([]SecretVisibility, 0, len(r.entries))
	for _, v := range r.entries {
		out = append(out, v)
	}
	return out
}
