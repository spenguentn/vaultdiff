package vault

import (
	"errors"
	"fmt"
	"strings"
)

// TagSet represents a collection of key-value metadata tags attached to a secret.
type TagSet map[string]string

// SecretTag associates a Vault secret path with a set of metadata tags.
type SecretTag struct {
	Mount string
	Path  string
	Tags  TagSet
}

// Validate checks that the SecretTag has required fields set.
func (t SecretTag) Validate() error {
	if strings.TrimSpace(t.Mount) == "" {
		return errors.New("secret tag: mount must not be empty")
	}
	if strings.TrimSpace(t.Path) == "" {
		return errors.New("secret tag: path must not be empty")
	}
	return nil
}

// FullPath returns the combined mount and path for display purposes.
func (t SecretTag) FullPath() string {
	return fmt.Sprintf("%s/%s", strings.Trim(t.Mount, "/"), strings.Trim(t.Path, "/"))
}

// Get returns the value for a tag key and whether it was found.
func (ts TagSet) Get(key string) (string, bool) {
	v, ok := ts[key]
	return v, ok
}

// Set adds or updates a tag.
func (ts TagSet) Set(key, value string) error {
	if strings.TrimSpace(key) == "" {
		return errors.New("tag key must not be empty")
	}
	ts[key] = value
	return nil
}

// Keys returns all tag keys in the set.
func (ts TagSet) Keys() []string {
	keys := make([]string, 0, len(ts))
	for k := range ts {
		keys = append(keys, k)
	}
	return keys
}

// Merge combines another TagSet into ts, overwriting existing keys.
func (ts TagSet) Merge(other TagSet) {
	for k, v := range other {
		ts[k] = v
	}
}

// NewSecretTag constructs a SecretTag and validates it.
func NewSecretTag(mount, path string, tags TagSet) (SecretTag, error) {
	if tags == nil {
		tags = make(TagSet)
	}
	t := SecretTag{Mount: mount, Path: path, Tags: tags}
	if err := t.Validate(); err != nil {
		return SecretTag{}, err
	}
	return t, nil
}
