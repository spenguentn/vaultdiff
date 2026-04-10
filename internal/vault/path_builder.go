package vault

import (
	"fmt"
	"strings"
)

// PathBuilder constructs KV v2 secret paths for a given environment.
type PathBuilder struct {
	env Environment
}

// NewPathBuilder returns a PathBuilder scoped to the provided environment.
func NewPathBuilder(env Environment) *PathBuilder {
	return &PathBuilder{env: env}
}

// Secret returns a SecretPath for the given logical secret path.
// The path is trimmed of leading/trailing slashes before construction.
func (b *PathBuilder) Secret(path string) (SecretPath, error) {
	path = strings.Trim(path, "/")
	if path == "" {
		return SecretPath{}, fmt.Errorf("path_builder: secret path must not be empty")
	}
	return NewSecretPath(b.env.MountPath, path)
}

// MustSecret is like Secret but panics on error. Useful in tests.
func (b *PathBuilder) MustSecret(path string) SecretPath {
	sp, err := b.Secret(path)
	if err != nil {
		panic(err)
	}
	return sp
}

// Batch returns SecretPaths for each of the provided logical paths.
// It returns the first error encountered.
func (b *PathBuilder) Batch(paths []string) ([]SecretPath, error) {
	result := make([]SecretPath, 0, len(paths))
	for _, p := range paths {
		sp, err := b.Secret(p)
		if err != nil {
			return nil, fmt.Errorf("path_builder: batch error on %q: %w", p, err)
		}
		result = append(result, sp)
	}
	return result, nil
}

// EnvPrefix returns the environment name for use as a log/display prefix.
func (b *PathBuilder) EnvPrefix() string {
	return fmt.Sprintf("[%s]", b.env.Name)
}
