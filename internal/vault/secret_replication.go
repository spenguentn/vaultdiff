package vault

import (
	"errors"
	"fmt"
	"time"
)

// ReplicationMode describes how a secret should be replicated.
type ReplicationMode string

const (
	ReplicationModeSync  ReplicationMode = "sync"
	ReplicationModeAsync ReplicationMode = "async"
)

// IsValidReplicationMode returns true if m is a known replication mode.
func IsValidReplicationMode(m ReplicationMode) bool {
	return m == ReplicationModeSync || m == ReplicationModeAsync
}

// ReplicationPolicy defines how a secret is replicated across environments.
type ReplicationPolicy struct {
	Mount       string          `json:"mount"`
	Path        string          `json:"path"`
	Targets     []string        `json:"targets"`
	Mode        ReplicationMode `json:"mode"`
	LastSyncAt  time.Time       `json:"last_sync_at,omitempty"`
	CreatedBy   string          `json:"created_by"`
}

// FullPath returns the combined mount and path.
func (r ReplicationPolicy) FullPath() string {
	return fmt.Sprintf("%s/%s", r.Mount, r.Path)
}

// Validate checks that the policy has all required fields.
func (r ReplicationPolicy) Validate() error {
	if r.Mount == "" {
		return errors.New("replication policy: mount is required")
	}
	if r.Path == "" {
		return errors.New("replication policy: path is required")
	}
	if len(r.Targets) == 0 {
		return errors.New("replication policy: at least one target is required")
	}
	if !IsValidReplicationMode(r.Mode) {
		return fmt.Errorf("replication policy: unknown mode %q", r.Mode)
	}
	if r.CreatedBy == "" {
		return errors.New("replication policy: created_by is required")
	}
	return nil
}

// replicationKey builds the registry key for a policy.
func replicationKey(mount, path string) string {
	return fmt.Sprintf("%s/%s", mount, path)
}

// SecretReplicationRegistry stores replication policies keyed by mount+path.
type SecretReplicationRegistry struct {
	policies map[string]ReplicationPolicy
}

// NewSecretReplicationRegistry returns an empty registry.
func NewSecretReplicationRegistry() *SecretReplicationRegistry {
	return &SecretReplicationRegistry{
		policies: make(map[string]ReplicationPolicy),
	}
}

// Set validates and stores a replication policy.
func (r *SecretReplicationRegistry) Set(p ReplicationPolicy) error {
	if err := p.Validate(); err != nil {
		return err
	}
	r.policies[replicationKey(p.Mount, p.Path)] = p
	return nil
}

// Get retrieves a policy by mount and path.
func (r *SecretReplicationRegistry) Get(mount, path string) (ReplicationPolicy, bool) {
	p, ok := r.policies[replicationKey(mount, path)]
	return p, ok
}

// Remove deletes a policy from the registry.
func (r *SecretReplicationRegistry) Remove(mount, path string) {
	delete(r.policies, replicationKey(mount, path))
}

// All returns a copy of all stored policies.
func (r *SecretReplicationRegistry) All() []ReplicationPolicy {
	out := make([]ReplicationPolicy, 0, len(r.policies))
	for _, p := range r.policies {
		out = append(out, p)
	}
	return out
}
