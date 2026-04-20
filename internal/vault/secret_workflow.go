package vault

import (
	"errors"
	"fmt"
	"time"
)

// WorkflowStage represents a stage in a secret workflow.
type WorkflowStage string

const (
	WorkflowStageDraft    WorkflowStage = "draft"
	WorkflowStagePending  WorkflowStage = "pending"
	WorkflowStageApproved WorkflowStage = "approved"
	WorkflowStageRejected WorkflowStage = "rejected"
	WorkflowStageDeployed WorkflowStage = "deployed"
	WorkflowStageRetired  WorkflowStage = "retired"
)

// IsValidWorkflowStage returns true if s is a known workflow stage.
func IsValidWorkflowStage(s WorkflowStage) bool {
	switch s {
	case WorkflowStageDraft, WorkflowStagePending, WorkflowStageApproved,
		WorkflowStageRejected, WorkflowStageDeployed, WorkflowStageRetired:
		return true
	}
	return false
}

// SecretWorkflow tracks the lifecycle workflow of a secret.
type SecretWorkflow struct {
	Mount       string        `json:"mount"`
	Path        string        `json:"path"`
	Stage       WorkflowStage `json:"stage"`
	AssignedTo  string        `json:"assigned_to"`
	InitiatedBy string        `json:"initiated_by"`
	UpdatedAt   time.Time     `json:"updated_at"`
	Comment     string        `json:"comment,omitempty"`
}

// FullPath returns the canonical path for this workflow entry.
func (w SecretWorkflow) FullPath() string {
	return fmt.Sprintf("%s/%s", w.Mount, w.Path)
}

// IsTerminal returns true if the workflow has reached a final stage.
func (w SecretWorkflow) IsTerminal() bool {
	return w.Stage == WorkflowStageDeployed || w.Stage == WorkflowStageRetired || w.Stage == WorkflowStageRejected
}

// Validate checks that the workflow entry has all required fields.
func (w SecretWorkflow) Validate() error {
	if w.Mount == "" {
		return errors.New("workflow: mount is required")
	}
	if w.Path == "" {
		return errors.New("workflow: path is required")
	}
	if w.InitiatedBy == "" {
		return errors.New("workflow: initiated_by is required")
	}
	if !IsValidWorkflowStage(w.Stage) {
		return fmt.Errorf("workflow: unknown stage %q", w.Stage)
	}
	return nil
}
