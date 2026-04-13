package vault

import (
	"testing"
	"time"
)

var baseBookmark = SecretBookmark{
	Mount:     "secret",
	Path:      "myapp/db",
	Alias:     "db-prod",
	CreatedBy: "alice",
}

func TestSecretBookmark_FullPath(t *testing.T) {
	got := baseBookmark.FullPath()
	if got != "secret/myapp/db" {
		t.Errorf("expected secret/myapp/db, got %s", got)
	}
}

func TestSecretBookmark_Validate_Valid(t *testing.T) {
	if err := baseBookmark.Validate(); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestSecretBookmark_Validate_MissingMount(t *testing.T) {
	b := baseBookmark
	b.Mount = ""
	if err := b.Validate(); err == nil {
		t.Error("expected error for missing mount")
	}
}

func TestSecretBookmark_Validate_MissingAlias(t *testing.T) {
	b := baseBookmark
	b.Alias = ""
	if err := b.Validate(); err == nil {
		t.Error("expected error for missing alias")
	}
}

func TestNewSecretBookmarkRegistry_NotNil(t *testing.T) {
	if NewSecretBookmarkRegistry() == nil {
		t.Error("expected non-nil registry")
	}
}

func TestBookmarkRegistry_Add_And_GetByAlias(t *testing.T) {
	r := NewSecretBookmarkRegistry()
	if err := r.Add(baseBookmark); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got, err := r.GetByAlias("db-prod")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Path != baseBookmark.Path {
		t.Errorf("expected path %s, got %s", baseBookmark.Path, got.Path)
	}
}

func TestBookmarkRegistry_Add_SetsCreatedAt(t *testing.T) {
	r := NewSecretBookmarkRegistry()
	before := time.Now().UTC()
	_ = r.Add(baseBookmark)
	got, _ := r.GetByAlias("db-prod")
	if got.CreatedAt.Before(before) {
		t.Error("expected CreatedAt to be set on add")
	}
}

func TestBookmarkRegistry_Add_DuplicateAlias(t *testing.T) {
	r := NewSecretBookmarkRegistry()
	_ = r.Add(baseBookmark)
	if err := r.Add(baseBookmark); err == nil {
		t.Error("expected error for duplicate alias")
	}
}

func TestBookmarkRegistry_GetByPath(t *testing.T) {
	r := NewSecretBookmarkRegistry()
	_ = r.Add(baseBookmark)
	got, err := r.GetByPath("secret", "myapp/db")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Alias != "db-prod" {
		t.Errorf("expected alias db-prod, got %s", got.Alias)
	}
}

func TestBookmarkRegistry_Remove(t *testing.T) {
	r := NewSecretBookmarkRegistry()
	_ = r.Add(baseBookmark)
	if err := r.Remove("db-prod"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, err := r.GetByAlias("db-prod"); err == nil {
		t.Error("expected error after removal")
	}
}

func TestBookmarkRegistry_All(t *testing.T) {
	r := NewSecretBookmarkRegistry()
	_ = r.Add(baseBookmark)
	if len(r.All()) != 1 {
		t.Errorf("expected 1 bookmark, got %d", len(r.All()))
	}
}
