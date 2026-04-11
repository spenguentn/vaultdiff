package vault

import (
	"errors"
	"sync"
)

// RotationScheduler manages a collection of rotation policies.
type RotationScheduler struct {
	mu       sync.RWMutex
	policies map[string]RotationPolicy
}

// NewRotationScheduler creates an empty RotationScheduler.
func NewRotationScheduler() *RotationScheduler {
	return &RotationScheduler{
		policies: make(map[string]RotationPolicy),
	}
}

func policyKey(mount, path string) string {
	return mount + "/" + path
}

// Register adds or replaces a rotation policy. Returns an error if invalid.
func (s *RotationScheduler) Register(p RotationPolicy) error {
	if err := p.Validate(); err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.policies[policyKey(p.Mount, p.Path)] = p
	return nil
}

// Remove deletes a policy by mount and path. Returns an error if not found.
func (s *RotationScheduler) Remove(mount, path string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	key := policyKey(mount, path)
	if _, ok := s.policies[key]; !ok {
		return errors.New("rotation scheduler: policy not found")
	}
	delete(s.policies, key)
	return nil
}

// DuePolicies returns all policies that are currently due for rotation.
func (s *RotationScheduler) DuePolicies() []RotationPolicy {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var due []RotationPolicy
	for _, p := range s.policies {
		if p.IsDue() {
			due = append(due, p)
		}
	}
	return due
}

// Count returns the total number of registered policies.
func (s *RotationScheduler) Count() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.policies)
}
