package vault

import "errors"

// CopyPlan accumulates a set of CopyRequests to be executed together.
type CopyPlan struct {
	requests []CopyRequest
}

// NewCopyPlan returns an empty CopyPlan.
func NewCopyPlan() *CopyPlan {
	return &CopyPlan{}
}

// Add appends a validated CopyRequest to the plan.
func (p *CopyPlan) Add(req CopyRequest) error {
	if err := req.Validate(); err != nil {
		return err
	}
	p.requests = append(p.requests, req)
	return nil
}

// Requests returns a copy of the accumulated requests.
func (p *CopyPlan) Requests() []CopyRequest {
	out := make([]CopyRequest, len(p.requests))
	copy(out, p.requests)
	return out
}

// Len returns the number of requests in the plan.
func (p *CopyPlan) Len() int { return len(p.requests) }

// Validate returns an error when the plan contains no requests.
func (p *CopyPlan) Validate() error {
	if len(p.requests) == 0 {
		return errors.New("copy plan has no requests")
	}
	return nil
}
