package vault

import (
	"errors"
	"fmt"
	"time"
)

// RenewResult holds the outcome of a lease renewal attempt.
type RenewResult struct {
	LeaseID    string
	RenewedAt  time.Time
	NewTTL     time.Duration
	Err        error
}

// IsSuccess returns true if the renewal completed without error.
func (r RenewResult) IsSuccess() bool {
	return r.Err == nil
}

// RenewRequest describes a single lease renewal request.
type RenewRequest struct {
	LeaseID   string
	Increment time.Duration // hint; Vault may ignore it
}

// Validate returns an error if the request is not usable.
func (r RenewRequest) Validate() error {
	if r.LeaseID == "" {
		return errors.New("renew: lease ID must not be empty")
	}
	if r.Increment < 0 {
		return fmt.Errorf("renew: increment must not be negative, got %s", r.Increment)
	}
	return nil
}

// Renewer issues renewal requests against a Vault logical backend.
type Renewer struct {
	client LogicalWriter
}

// LogicalWriter is the minimal interface required to renew leases.
type LogicalWriter interface {
	Write(path string, data map[string]interface{}) (*SecretResponse, error)
}

// NewRenewer constructs a Renewer. Panics if client is nil.
func NewRenewer(client LogicalWriter) *Renewer {
	if client == nil {
		panic("renewer: client must not be nil")
	}
	return &Renewer{client: client}
}

// Renew attempts to renew the given lease and returns the result.
func (r *Renewer) Renew(req RenewRequest) RenewResult {
	if err := req.Validate(); err != nil {
		return RenewResult{LeaseID: req.LeaseID, Err: err}
	}

	data := map[string]interface{}{
		"lease_id":  req.LeaseID,
		"increment": int(req.Increment.Seconds()),
	}

	resp, err := r.client.Write("sys/leases/renew", data)
	if err != nil {
		return RenewResult{LeaseID: req.LeaseID, Err: err}
	}

	lease, err := ParseLease(resp)
	if err != nil {
		return RenewResult{LeaseID: req.LeaseID, Err: err}
	}

	return RenewResult{
		LeaseID:   lease.ID,
		RenewedAt: time.Now().UTC(),
		NewTTL:    lease.TTL,
	}
}
