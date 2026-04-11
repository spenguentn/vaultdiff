package vault

import (
	"errors"
	"fmt"
	"time"
)

// ApprovalStatus represents the current state of a change approval request.
type ApprovalStatus string

const (
	ApprovalPending  ApprovalStatus = "pending"
	ApprovalApproved ApprovalStatus = "approved"
	ApprovalRejected ApprovalStatus = "rejected"
	ApprovalRevoked  ApprovalStatus = "revoked"
)

// ApprovalRequest represents a request to approve a secret change before it is applied.
type ApprovalRequest struct {
	ID          string
	Mount       string
	Path        string
	RequestedBy string
	Reason      string
	Status      ApprovalStatus
	CreatedAt   time.Time
	ReviewedBy  string
	ReviewedAt  *time.Time
}

// FullPath returns the mount-qualified secret path.
func (a *ApprovalRequest) FullPath() string {
	return fmt.Sprintf("%s/%s", a.Mount, a.Path)
}

// IsTerminal returns true if the request is in a final state.
func (a *ApprovalRequest) IsTerminal() bool {
	return a.Status == ApprovalApproved || a.Status == ApprovalRejected || a.Status == ApprovalRevoked
}

// Validate checks that the approval request has all required fields.
func (a *ApprovalRequest) Validate() error {
	if a.Mount == "" {
		return errors.New("approval: mount is required")
	}
	if a.Path == "" {
		return errors.New("approval: path is required")
	}
	if a.RequestedBy == "" {
		return errors.New("approval: requested_by is required")
	}
	if a.Reason == "" {
		return errors.New("approval: reason is required")
	}
	return nil
}
