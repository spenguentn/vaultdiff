package vault

import (
	"testing"
	"time"
)

func baseDependencyLink() DependencyLink {
	return DependencyLink{
		SourceMount: "secret",
		SourcePath:  "app/api",
		TargetMount: "secret",
		TargetPath:  "shared/db",
		AddedBy:     "alice",
		AddedAt:     time.Now().UTC(),
		Note:        "api depends on db creds",
	}
}

func TestDependencyLink_FullSource(t *testing.T) {
	d := baseDependencyLink()
	if got := d.FullSource(); got != "secret/app/api" {
		t.Errorf("expected secret/app/api, got %s", got)
	}
}

func TestDependencyLink_FullTarget(t *testing.T) {
	d := baseDependencyLink()
	if got := d.FullTarget(); got != "secret/shared/db" {
		t.Errorf("expected secret/shared/db, got %s", got)
	}
}

func TestDependencyLink_Validate_Valid(t *testing.T) {
	if err := baseDependencyLink().Validate(); err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestDependencyLink_Validate_MissingSourceMount(t *testing.T) {
	d := baseDependencyLink()
	d.SourceMount = ""
	if err := d.Validate(); err == nil {
		t.Error("expected error for missing source mount")
	}
}

func TestDependencyLink_Validate_MissingTargetPath(t *testing.T) {
	d := baseDependencyLink()
	d.TargetPath = ""
	if err := d.Validate(); err == nil {
		t.Error("expected error for missing target path")
	}
}

func TestDependencyLink_Validate_MissingAddedBy(t *testing.T) {
	d := baseDependencyLink()
	d.AddedBy = ""
	if err := d.Validate(); err == nil {
		t.Error("expected error for missing added_by")
	}
}

func TestNewSecretDependencyRegistry_NotNil(t *testing.T) {
	if r := NewSecretDependencyRegistry(); r == nil {
		t.Error("expected non-nil registry")
	}
}

func TestDependencyRegistry_Add_And_Get(t *testing.T) {
	r := NewSecretDependencyRegistry()
	if err := r.Add(baseDependencyLink()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	links := r.GetDependencies("secret", "app/api")
	if len(links) != 1 {
		t.Errorf("expected 1 link, got %d", len(links))
	}
}

func TestDependencyRegistry_Add_SetsTimestamp(t *testing.T) {
	r := NewSecretDependencyRegistry()
	d := baseDependencyLink()
	d.AddedAt = time.Time{}
	_ = r.Add(d)
	links := r.GetDependencies("secret", "app/api")
	if links[0].AddedAt.IsZero() {
		t.Error("expected AddedAt to be set automatically")
	}
}

func TestDependencyRegistry_Add_Invalid(t *testing.T) {
	r := NewSecretDependencyRegistry()
	d := baseDependencyLink()
	d.SourceMount = ""
	if err := r.Add(d); err == nil {
		t.Error("expected validation error")
	}
}

func TestDependencyRegistry_Get_NotFound(t *testing.T) {
	r := NewSecretDependencyRegistry()
	links := r.GetDependencies("secret", "nonexistent")
	if len(links) != 0 {
		t.Errorf("expected empty slice, got %d", len(links))
	}
}

func TestDependencyRegistry_Remove(t *testing.T) {
	r := NewSecretDependencyRegistry()
	_ = r.Add(baseDependencyLink())
	r.Remove("secret", "app/api")
	if links := r.GetDependencies("secret", "app/api"); len(links) != 0 {
		t.Error("expected links to be removed")
	}
}

func TestDependencyRegistry_Count(t *testing.T) {
	r := NewSecretDependencyRegistry()
	_ = r.Add(baseDependencyLink())
	d2 := baseDependencyLink()
	d2.TargetPath = "shared/cache"
	_ = r.Add(d2)
	if got := r.Count(); got != 2 {
		t.Errorf("expected count 2, got %d", got)
	}
}
