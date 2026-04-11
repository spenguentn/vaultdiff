package vault

import (
	"fmt"
	"sync"
	"time"
)

// SecretApprovalRegistry manages in-memory approval requests for secret changes.
type SecretApprovalRegistry struct {
	mu       sync.RWMutex
	requests map[string]*ApprovalRequest
}

// NewSecretApprovalRegistry returns an initialised SecretApprovalRegistry.
func NewSecretApprovalRegistry() *SecretApprovalRegistry {
	return &SecretApprovalRegistry{
		requests: make(map[string]*ApprovalRequest),
	}
}

func approvalKey(mount, path string) string {
	return fmt.Sprintf("%s::%s", mount, path)
}

// Submit adds a new approval request. Returns an error if the request is invalid
// or a pending request already exists for the same path.
func (r *SecretApprovalRegistry) Submit(req *ApprovalRequest) error {
	if err := req.Validate(); err != nil {
		return err
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	key := approvalKey(req.Mount, req.Path)
	if existing, ok := r.requests[key]; ok && existing.Status == ApprovalPending {
		return fmt.Errorf("approval: pending request already exists for %s", key)
	}
	req.Status = ApprovalPending
	req.CreatedAt = time.Now().UTC()
	r.requests[key] = req
	return nil
}

// Get returns the approval request for the given mount and path.
func (r *SecretApprovalRegistry) Get(mount, path string) (*ApprovalRequest, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	v, ok := r.requests[approvalKey(mount, path)]
	return v, ok
}

// Review updates the status of an existing request.
func (r *SecretApprovalRegistry) Review(mount, path, reviewer string, approved bool) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	key := approvalKey(mount, path)
	req, ok := r.requests[key]
	if !ok {
		return fmt.Errorf("approval: no request found for %s", key)
	}
	if req.IsTerminal() {
		return fmt.Errorf("approval: request for %s is already in terminal state %s", key, req.Status)
	}
	now := time.Now().UTC()
	req.ReviewedBy = reviewer
	req.ReviewedAt = &now
	if approved {
		req.Status = ApprovalApproved
	} else {
		req.Status = ApprovalRejected
	}
	return nil
}

// Revoke moves an approved request back to revoked state.
func (r *SecretApprovalRegistry) Revoke(mount, path string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	key := approvalKey(mount, path)
	req, ok := r.requests[key]
	if !ok {
		return fmt.Errorf("approval: no request found for %s", key)
	}
	if req.Status != ApprovalApproved {
		return fmt.Errorf("approval: only approved requests can be revoked")
	}
	req.Status = ApprovalRevoked
	return nil
}

// All returns a snapshot of all stored requests.
func (r *SecretApprovalRegistry) All() []*ApprovalRequest {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]*ApprovalRequest, 0, len(r.requests))
	for _, v := range r.requests {
		out = append(out, v)
	}
	return out
}
