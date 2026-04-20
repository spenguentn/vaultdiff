package vault

import (
	"fmt"
	"sync"
	"time"
)

func workflowKey(mount, path string) string {
	return fmt.Sprintf("%s/%s", mount, path)
}

// SecretWorkflowRegistry stores and manages secret workflow states.
type SecretWorkflowRegistry struct {
	mu      sync.RWMutex
	entries map[string]SecretWorkflow
}

// NewSecretWorkflowRegistry returns an initialised SecretWorkflowRegistry.
func NewSecretWorkflowRegistry() *SecretWorkflowRegistry {
	return &SecretWorkflowRegistry{
		entries: make(map[string]SecretWorkflow),
	}
}

// Set stores or updates a workflow entry after validation.
func (r *SecretWorkflowRegistry) Set(w SecretWorkflow) error {
	if err := w.Validate(); err != nil {
		return err
	}
	w.UpdatedAt = time.Now().UTC()
	r.mu.Lock()
	defer r.mu.Unlock()
	r.entries[workflowKey(w.Mount, w.Path)] = w
	return nil
}

// Get retrieves the workflow entry for the given mount and path.
func (r *SecretWorkflowRegistry) Get(mount, path string) (SecretWorkflow, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	w, ok := r.entries[workflowKey(mount, path)]
	return w, ok
}

// Advance transitions the workflow to the next stage.
func (r *SecretWorkflowRegistry) Advance(mount, path string, next WorkflowStage, actor, comment string) error {
	if !IsValidWorkflowStage(next) {
		return fmt.Errorf("workflow: unknown stage %q", next)
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	key := workflowKey(mount, path)
	w, ok := r.entries[key]
	if !ok {
		return fmt.Errorf("workflow: no entry for %s/%s", mount, path)
	}
	w.Stage = next
	w.AssignedTo = actor
	w.Comment = comment
	w.UpdatedAt = time.Now().UTC()
	r.entries[key] = w
	return nil
}

// Remove deletes the workflow entry for the given mount and path.
func (r *SecretWorkflowRegistry) Remove(mount, path string) bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	key := workflowKey(mount, path)
	_, ok := r.entries[key]
	delete(r.entries, key)
	return ok
}

// All returns a snapshot of all workflow entries.
func (r *SecretWorkflowRegistry) All() []SecretWorkflow {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]SecretWorkflow, 0, len(r.entries))
	for _, w := range r.entries {
		out = append(out, w)
	}
	return out
}
