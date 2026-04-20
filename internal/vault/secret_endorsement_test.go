package vault

import (
	"testing"
	"time"
)

var baseEndorsement = &SecretEndorsement{
	Mount:      "secret",
	Path:       "app/db",
	EndorsedBy: "alice",
	Status:     "approved",
	Comment:    "looks good",
}

func TestIsValidEndorsementStatus_Known(t *testing.T) {
	for _, s := range []string{"pending", "approved", "rejected", "revoked"} {
		if !IsValidEndorsementStatus(s) {
			t.Errorf("expected %q to be valid", s)
		}
	}
}

func TestIsValidEndorsementStatus_Unknown(t *testing.T) {
	if IsValidEndorsementStatus("unknown") {
		t.Error("expected 'unknown' to be invalid")
	}
}

func TestSecretEndorsement_FullPath(t *testing.T) {
	got := baseEndorsement.FullPath()
	want := "secret/app/db"
	if got != want {
		t.Errorf("FullPath() = %q, want %q", got, want)
	}
}

func TestSecretEndorsement_Validate_Valid(t *testing.T) {
	e := *baseEndorsement
	if err := e.Validate(); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestSecretEndorsement_Validate_MissingMount(t *testing.T) {
	e := *baseEndorsement
	e.Mount = ""
	if err := e.Validate(); err == nil {
		t.Error("expected error for missing mount")
	}
}

func TestSecretEndorsement_Validate_MissingEndorsedBy(t *testing.T) {
	e := *baseEndorsement
	e.EndorsedBy = ""
	if err := e.Validate(); err == nil {
		t.Error("expected error for missing endorsed_by")
	}
}

func TestSecretEndorsement_IsApproved_True(t *testing.T) {
	e := *baseEndorsement
	if !e.IsApproved() {
		t.Error("expected IsApproved() = true")
	}
}

func TestSecretEndorsement_IsApproved_False(t *testing.T) {
	e := *baseEndorsement
	e.Status = "pending"
	if e.IsApproved() {
		t.Error("expected IsApproved() = false for pending")
	}
}

func TestNewSecretEndorsementRegistry_NotNil(t *testing.T) {
	if NewSecretEndorsementRegistry() == nil {
		t.Error("expected non-nil registry")
	}
}

func TestEndorsementRegistry_Submit_And_Get(t *testing.T) {
	r := NewSecretEndorsementRegistry()
	e := *baseEndorsement
	if err := r.Submit(&e); err != nil {
		t.Fatalf("Submit() error: %v", err)
	}
	got, ok := r.Get(e.Mount, e.Path, e.EndorsedBy)
	if !ok {
		t.Fatal("expected endorsement to be found")
	}
	if got.Status != "approved" {
		t.Errorf("Status = %q, want %q", got.Status, "approved")
	}
}

func TestEndorsementRegistry_Submit_SetsEndorsedAt(t *testing.T) {
	r := NewSecretEndorsementRegistry()
	e := *baseEndorsement
	e.EndorsedAt = time.Time{}
	_ = r.Submit(&e)
	if e.EndorsedAt.IsZero() {
		t.Error("expected EndorsedAt to be set")
	}
}

func TestEndorsementRegistry_Submit_Invalid(t *testing.T) {
	r := NewSecretEndorsementRegistry()
	e := *baseEndorsement
	e.Mount = ""
	if err := r.Submit(&e); err == nil {
		t.Error("expected validation error")
	}
}

func TestEndorsementRegistry_Get_NotFound(t *testing.T) {
	r := NewSecretEndorsementRegistry()
	_, ok := r.Get("secret", "app/db", "bob")
	if ok {
		t.Error("expected not found")
	}
}

func TestEndorsementRegistry_Remove(t *testing.T) {
	r := NewSecretEndorsementRegistry()
	e := *baseEndorsement
	_ = r.Submit(&e)
	r.Remove(e.Mount, e.Path, e.EndorsedBy)
	_, ok := r.Get(e.Mount, e.Path, e.EndorsedBy)
	if ok {
		t.Error("expected endorsement to be removed")
	}
}

func TestEndorsementRegistry_List(t *testing.T) {
	r := NewSecretEndorsementRegistry()
	e1 := *baseEndorsement
	e2 := *baseEndorsement
	e2.EndorsedBy = "bob"
	e2.Status = "pending"
	_ = r.Submit(&e1)
	_ = r.Submit(&e2)
	list := r.List("secret", "app/db")
	if len(list) != 2 {
		t.Errorf("List() len = %d, want 2", len(list))
	}
}
