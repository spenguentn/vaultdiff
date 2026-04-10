package vault

import (
	"errors"
	"sync"
)

// LeaseTracker maintains a registry of active leases keyed by their lease ID.
// It is safe for concurrent use.
type LeaseTracker struct {
	mu     sync.RWMutex
	leases map[string]LeaseInfo
}

// NewLeaseTracker returns an initialised LeaseTracker.
func NewLeaseTracker() *LeaseTracker {
	return &LeaseTracker{
		leases: make(map[string]LeaseInfo),
	}
}

// Track registers a lease. Returns an error if the lease fails validation.
func (t *LeaseTracker) Track(l LeaseInfo) error {
	if err := l.Validate(); err != nil {
		return err
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	t.leases[l.LeaseID] = l
	return nil
}

// Get retrieves a lease by ID. Returns false if not found.
func (t *LeaseTracker) Get(leaseID string) (LeaseInfo, bool) {
	t.mu.RLock()
	defer t.mu.RUnlock()
	l, ok := t.leases[leaseID]
	return l, ok
}

// Revoke removes a lease from the tracker. Returns an error if not found.
func (t *LeaseTracker) Revoke(leaseID string) error {
	t.mu.Lock()
	defer t.mu.Unlock()
	if _, ok := t.leases[leaseID]; !ok {
		return errors.New("lease tracker: lease not found: " + leaseID)
	}
	delete(t.leases, leaseID)
	return nil
}

// Expired returns all leases that have passed their expiry time.
func (t *LeaseTracker) Expired() []LeaseInfo {
	t.mu.RLock()
	defer t.mu.RUnlock()
	var out []LeaseInfo
	for _, l := range t.leases {
		if l.IsExpired() {
			out = append(out, l)
		}
	}
	return out
}

// Count returns the total number of tracked leases.
func (t *LeaseTracker) Count() int {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return len(t.leases)
}
