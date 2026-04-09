// Package snapshot provides types and utilities for capturing and comparing
// point-in-time snapshots of Vault secret paths.
package snapshot

import "time"

// Snapshot represents a captured state of secrets at a given path and version.
type Snapshot struct {
	Path      string
	Version   int
	CapturedAt time.Time
	Secrets   map[string]string
	Meta      Meta
}

// Meta holds optional metadata associated with a snapshot.
type Meta struct {
	Environment string
	Operator    string
	Label       string
}

// New creates a new Snapshot with the current timestamp.
func New(path string, version int, secrets map[string]string, meta Meta) *Snapshot {
	if secrets == nil {
		secrets = make(map[string]string)
	}
	return &Snapshot{
		Path:       path,
		Version:    version,
		CapturedAt: time.Now().UTC(),
		Secrets:    secrets,
		Meta:       meta,
	}
}

// KeyCount returns the number of secret keys in the snapshot.
func (s *Snapshot) KeyCount() int {
	return len(s.Secrets)
}

// HasKey reports whether the snapshot contains the given key.
func (s *Snapshot) HasKey(key string) bool {
	_, ok := s.Secrets[key]
	return ok
}

// Equal reports whether two snapshots contain identical secret data.
func (s *Snapshot) Equal(other *Snapshot) bool {
	if len(s.Secrets) != len(other.Secrets) {
		return false
	}
	for k, v := range s.Secrets {
		if other.Secrets[k] != v {
			return false
		}
	}
	return true
}
