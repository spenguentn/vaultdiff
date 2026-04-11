package vault

import (
	"testing"
	"time"
)

var baseOwnership = OwnershipRecord{
	Mount:   "secret",
	Path:    "app/db-password",
	Owner:   "alice",
	Team:    "platform",
	Contact: "alice@example.com",
}

func TestOwnershipRecord_FullPath(t *testing.T) {
	rec := baseOwnership
	if got := rec.FullPath(); got != "secret/app/db-password" {
		t.Errorf("expected 'secret/app/db-password', got %q", got)
	}
}

func TestOwnershipRecord_Validate_Valid(t *testing.T) {
	if err := baseOwnership.Validate(); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestOwnershipRecord_Validate_MissingMount(t *testing.T) {
	rec := baseOwnership
	rec.Mount = ""
	if err := rec.Validate(); err == nil {
		t.Error("expected error for missing mount")
	}
}

func TestOwnershipRecord_Validate_MissingPath(t *testing.T) {
	rec := baseOwnership
	rec.Path = ""
	if err := rec.Validate(); err == nil {
		t.Error("expected error for missing path")
	}
}

func TestOwnershipRecord_Validate_MissingOwner(t *testing.T) {
	rec := baseOwnership
	rec.Owner = ""
	if err := rec.Validate(); err == nil {
		t.Error("expected error for missing owner")
	}
}

func TestNewOwnershipRegistry_NotNil(t *testing.T) {
	if r := NewOwnershipRegistry(); r == nil {
		t.Error("expected non-nil registry")
	}
}

func TestOwnershipRegistry_Assign_And_Get(t *testing.T) {
	r := NewOwnershipRegistry()
	if err := r.Assign(baseOwnership); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	rec, ok := r.Get(baseOwnership.Mount, baseOwnership.Path)
	if !ok {
		t.Fatal("expected record to be found")
	}
	if rec.Owner != "alice" {
		t.Errorf("expected owner 'alice', got %q", rec.Owner)
	}
}

func TestOwnershipRegistry_Assign_SetsTimestamp(t *testing.T) {
	r := NewOwnershipRegistry()
	rec := baseOwnership
	rec.AssignedAt = time.Time{}
	_ = r.Assign(rec)
	got, _ := r.Get(rec.Mount, rec.Path)
	if got.AssignedAt.IsZero() {
		t.Error("expected AssignedAt to be set automatically")
	}
}

func TestOwnershipRegistry_Assign_Invalid(t *testing.T) {
	r := NewOwnershipRegistry()
	rec := baseOwnership
	rec.Owner = ""
	if err := r.Assign(rec); err == nil {
		t.Error("expected validation error")
	}
}

func TestOwnershipRegistry_Get_NotFound(t *testing.T) {
	r := NewOwnershipRegistry()
	_, ok := r.Get("secret", "nonexistent")
	if ok {
		t.Error("expected not found")
	}
}

func TestOwnershipRegistry_Remove(t *testing.T) {
	r := NewOwnershipRegistry()
	_ = r.Assign(baseOwnership)
	if !r.Remove(baseOwnership.Mount, baseOwnership.Path) {
		t.Error("expected Remove to return true")
	}
	_, ok := r.Get(baseOwnership.Mount, baseOwnership.Path)
	if ok {
		t.Error("expected record to be gone after removal")
	}
}

func TestOwnershipRegistry_Remove_NotFound(t *testing.T) {
	r := NewOwnershipRegistry()
	if r.Remove("secret", "missing") {
		t.Error("expected Remove to return false for unknown key")
	}
}

func TestOwnershipRegistry_All(t *testing.T) {
	r := NewOwnershipRegistry()
	_ = r.Assign(baseOwnership)
	second := baseOwnership
	second.Path = "app/api-key"
	_ = r.Assign(second)
	if len(r.All()) != 2 {
		t.Errorf("expected 2 records, got %d", len(r.All()))
	}
}
