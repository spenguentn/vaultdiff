package vault

import "fmt"

// PromotePlan holds an ordered list of PromoteRequests to execute.
type PromotePlan struct {
	requests []PromoteRequest
}

// NewPromotePlan returns an empty PromotePlan.
func NewPromotePlan() *PromotePlan {
	return &PromotePlan{}
}

// Add validates and appends a PromoteRequest to the plan.
func (p *PromotePlan) Add(req PromoteRequest) error {
	if err := req.Validate(); err != nil {
		return fmt.Errorf("invalid promote request: %w", err)
	}
	p.requests = append(p.requests, req)
	return nil
}

// Requests returns a copy of the queued requests.
func (p *PromotePlan) Requests() []PromoteRequest {
	out := make([]PromoteRequest, len(p.requests))
	copy(out, p.requests)
	return out
}

// Len returns the number of requests in the plan.
func (p *PromotePlan) Len() int {
	return len(p.requests)
}

// Validate returns an error when the plan contains no requests.
func (p *PromotePlan) Validate() error {
	if len(p.requests) == 0 {
		return fmt.Errorf("promote plan is empty")
	}
	return nil
}
