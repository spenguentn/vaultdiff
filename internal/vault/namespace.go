package vault

import (
	"errors"
	"strings"
)

// Namespace represents a Vault Enterprise namespace path.
type Namespace struct {
	Path string
}

// NewNamespace creates a Namespace from the given path, trimming leading/trailing slashes.
func NewNamespace(path string) Namespace {
	return Namespace{Path: strings.Trim(path, "/")}
}

// IsRoot reports whether the namespace is the root (empty) namespace.
func (n Namespace) IsRoot() bool {
	return n.Path == ""
}

// String returns the canonical namespace path.
func (n Namespace) String() string {
	if n.IsRoot() {
		return "root"
	}
	return n.Path
}

// Child returns a new Namespace that is a child of n with the given segment appended.
func (n Namespace) Child(segment string) Namespace {
	segment = strings.Trim(segment, "/")
	if segment == "" {
		return n
	}
	if n.IsRoot() {
		return Namespace{Path: segment}
	}
	return Namespace{Path: n.Path + "/" + segment}
}

// Validate returns an error if the namespace path contains invalid characters.
func (n Namespace) Validate() error {
	if n.IsRoot() {
		return nil
	}
	if strings.Contains(n.Path, " ") {
		return errors.New("namespace path must not contain spaces")
	}
	for _, seg := range strings.Split(n.Path, "/") {
		if seg == "" {
			return errors.New("namespace path must not contain empty segments")
		}
	}
	return nil
}

// HeaderValue returns the value suitable for the X-Vault-Namespace header.
// Returns an empty string for the root namespace.
func (n Namespace) HeaderValue() string {
	return n.Path
}
