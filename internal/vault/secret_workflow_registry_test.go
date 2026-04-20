package vault

import (
	"testing"
)

func TestNewSecretWorkflowRegistry_NotNil(t *testing.T) {
	if NewSecretWorkflowRegistry() == nil {
		t.Error("expected non-nil registry")
	}
}

func TestWorkflowRegistry_Set_And_Get(t *testing.T) {
	r := NewSecretWorkflowRegistry()
	if err := r.Set(baseWorkflow); err != nil {
		t.Fatalf("Set failed: %v", err)
	}
	got, ok := r.Get(baseWorkflow.Mount, baseWorkflow.Path)
	if !ok {
		t.Fatal("expected entry to exist")
	}
	if got.Stage != baseWorkflow.Stage {
		t.Errorf("stage mismatch: got %s", got.Stage)
	}
}

func TestWorkflowRegistry_Set_SetsUpdatedAt(t *testing.T) {
	r := NewSecretWorkflowRegistry()
	_ = r.Set(baseWorkflow)
	got, _ := r.Get(baseWorkflow.Mount, baseWorkflow.Path)
	if got.UpdatedAt.IsZero() {
		t.Error("expected UpdatedAt to be set")
	}
}

func TestWorkflowRegistry_Set_Invalid(t *testing.T) {
	r := NewSecretWorkflowRegistry()
	w := baseWorkflow
	w.Mount = ""
	if err := r.Set(w); err == nil {
		t.Error("expected error for invalid entry")
	}
}

func TestWorkflowRegistry_Get_NotFound(t *testing.T) {
	r := NewSecretWorkflowRegistry()
	_, ok := r.Get("missing", "path")
	if ok {
		t.Error("expected not found")
	}
}

func TestWorkflowRegistry_Advance_Valid(t *testing.T) {
	r := NewSecretWorkflowRegistry()
	_ = r.Set(baseWorkflow)
	err := r.Advance(baseWorkflow.Mount, baseWorkflow.Path, WorkflowStageApproved, "carol", "looks good")
	if err != nil {
		t.Fatalf("Advance failed: %v", err)
	}
	got, _ := r.Get(baseWorkflow.Mount, baseWorkflow.Path)
	if got.Stage != WorkflowStageApproved {
		t.Errorf("expected approved, got %s", got.Stage)
	}
}

func TestWorkflowRegistry_Advance_NotFound(t *testing.T) {
	r := NewSecretWorkflowRegistry()
	if err := r.Advance("x", "y", WorkflowStageApproved, "actor", ""); err == nil {
		t.Error("expected error for missing entry")
	}
}

func TestWorkflowRegistry_Remove(t *testing.T) {
	r := NewSecretWorkflowRegistry()
	_ = r.Set(baseWorkflow)
	if !r.Remove(baseWorkflow.Mount, baseWorkflow.Path) {
		t.Error("expected Remove to return true")
	}
	_, ok := r.Get(baseWorkflow.Mount, baseWorkflow.Path)
	if ok {
		t.Error("expected entry to be gone")
	}
}

func TestWorkflowRegistry_All(t *testing.T) {
	r := NewSecretWorkflowRegistry()
	_ = r.Set(baseWorkflow)
	if len(r.All()) != 1 {
		t.Errorf("expected 1 entry, got %d", len(r.All()))
	}
}
