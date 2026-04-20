package vault

import (
	"testing"
)

var baseWorkflow = SecretWorkflow{
	Mount:       "secret",
	Path:        "app/db-password",
	Stage:       WorkflowStagePending,
	InitiatedBy: "alice",
	AssignedTo:  "bob",
}

func TestIsValidWorkflowStage_Known(t *testing.T) {
	for _, s := range []WorkflowStage{
		WorkflowStageDraft, WorkflowStagePending, WorkflowStageApproved,
		WorkflowStageRejected, WorkflowStageDeployed, WorkflowStageRetired,
	} {
		if !IsValidWorkflowStage(s) {
			t.Errorf("expected %q to be valid", s)
		}
	}
}

func TestIsValidWorkflowStage_Unknown(t *testing.T) {
	if IsValidWorkflowStage("unknown") {
		t.Error("expected unknown stage to be invalid")
	}
}

func TestSecretWorkflow_FullPath(t *testing.T) {
	w := baseWorkflow
	if got := w.FullPath(); got != "secret/app/db-password" {
		t.Errorf("unexpected full path: %s", got)
	}
}

func TestSecretWorkflow_IsTerminal_Deployed(t *testing.T) {
	w := baseWorkflow
	w.Stage = WorkflowStageDeployed
	if !w.IsTerminal() {
		t.Error("expected deployed to be terminal")
	}
}

func TestSecretWorkflow_IsTerminal_Pending(t *testing.T) {
	w := baseWorkflow
	if w.IsTerminal() {
		t.Error("expected pending to be non-terminal")
	}
}

func TestSecretWorkflow_Validate_Valid(t *testing.T) {
	if err := baseWorkflow.Validate(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestSecretWorkflow_Validate_MissingMount(t *testing.T) {
	w := baseWorkflow
	w.Mount = ""
	if err := w.Validate(); err == nil {
		t.Error("expected error for missing mount")
	}
}

func TestSecretWorkflow_Validate_MissingInitiatedBy(t *testing.T) {
	w := baseWorkflow
	w.InitiatedBy = ""
	if err := w.Validate(); err == nil {
		t.Error("expected error for missing initiated_by")
	}
}

func TestSecretWorkflow_Validate_InvalidStage(t *testing.T) {
	w := baseWorkflow
	w.Stage = "bogus"
	if err := w.Validate(); err == nil {
		t.Error("expected error for invalid stage")
	}
}
