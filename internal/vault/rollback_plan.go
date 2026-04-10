package vault

import "fmt"

// RollbackPlan holds a set of rollback requests to be executed together.
type RollbackPlan struct {
	requests []RollbackRequest
}

// NewRollbackPlan creates an empty plan.
func NewRollbackPlan() *RollbackPlan {
	return &RollbackPlan{}
}

// Add appends a rollback request to the plan.
func (p *RollbackPlan) Add(req RollbackRequest) error {
	if req.Mount == "" {
		return fmt.Errorf("rollback plan: mount must not be empty")
	}
	if req.Path == "" {
		return fmt.Errorf("rollback plan: path must not be empty")
	}
	if req.Version < 1 {
		return fmt.Errorf("rollback plan: version must be >= 1")
	}
	p.requests = append(p.requests, req)
	return nil
}

// Requests returns a copy of the planned rollback requests.
func (p *RollbackPlan) Requests() []RollbackRequest {
	out := make([]RollbackRequest, len(p.requests))
	copy(out, p.requests)
	return out
}

// Len returns the number of requests in the plan.
func (p *RollbackPlan) Len() int { return len(p.requests) }

// IsEmpty returns true when no requests have been added.
func (p *RollbackPlan) IsEmpty() bool { return len(p.requests) == 0 }

// Describe returns a human-readable summary of the plan.
func (p *RollbackPlan) Describe() string {
	if p.IsEmpty() {
		return "rollback plan: (empty)"
	}
	out := fmt.Sprintf("rollback plan: %d operation(s)\n", p.Len())
	for i, r := range p.requests {
		out += fmt.Sprintf("  [%d] %s/%s -> version %d\n", i+1, r.Mount, r.Path, r.Version)
	}
	return out
}
