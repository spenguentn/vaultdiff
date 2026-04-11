package vault

import (
	"testing"
	"time"
)

var baseArchiveEntry = SecretArchiveEntry{
	Mount:      "secret",
	Path:       "myapp/db",
	Version:    3,
	Reason:     ArchiveReasonRotated,
	ArchivedBy: "admin",
}

func TestSecretArchiveEntry_Validate_Valid(t *testing.T) {
	e := baseArchiveEntry
	if err := e.Validate(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestSecretArchiveEntry_Validate_MissingMount(t *testing.T) {
	e := baseArchiveEntry
	e.Mount = ""
	if err := e.Validate(); err == nil {
		t.Fatal("expected error for missing mount")
	}
}

func TestSecretArchiveEntry_Validate_MissingPath(t *testing.T) {
	e := baseArchiveEntry
	e.Path = ""
	if err := e.Validate(); err == nil {
		t.Fatal("expected error for missing path")
	}
}

func TestSecretArchiveEntry_Validate_MissingArchivedBy(t *testing.T) {
	e := baseArchiveEntry
	e.ArchivedBy = ""
	if err := e.Validate(); err == nil {
		t.Fatal("expected error for missing archived_by")
	}
}

func TestSecretArchiveEntry_FullPath(t *testing.T) {
	e := baseArchiveEntry
	if got := e.FullPath(); got != "secret/myapp/db" {
		t.Fatalf("expected secret/myapp/db, got %s", got)
	}
}

func TestIsReasonValid_Known(t *testing.T) {
	for _, r := range []ArchiveReason{ArchiveReasonDeprecated, ArchiveReasonRotated, ArchiveReasonMigrated, ArchiveReasonManual} {
		if !IsReasonValid(r) {
			t.Fatalf("expected %q to be valid", r)
		}
	}
}

func TestIsReasonValid_Unknown(t *testing.T) {
	if IsReasonValid(ArchiveReason("unknown")) {
		t.Fatal("expected unknown reason to be invalid")
	}
}

func TestArchiveRegistry_ArchiveAndGet(t *testing.T) {
	reg := NewSecretArchiveRegistry()
	e := baseArchiveEntry
	if err := reg.Archive(&e); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got, ok := reg.Get(e.Mount, e.Path)
	if !ok {
		t.Fatal("expected entry to be found")
	}
	if got.Path != e.Path {
		t.Fatalf("path mismatch: %s", got.Path)
	}
}

func TestArchiveRegistry_Archive_SetsTimestamp(t *testing.T) {
	reg := NewSecretArchiveRegistry()
	e := baseArchiveEntry
	e.ArchivedAt = time.Time{}
	_ = reg.Archive(&e)
	if e.ArchivedAt.IsZero() {
		t.Fatal("expected ArchivedAt to be set")
	}
}

func TestArchiveRegistry_Get_NotFound(t *testing.T) {
	reg := NewSecretArchiveRegistry()
	_, ok := reg.Get("secret", "missing/path")
	if ok {
		t.Fatal("expected not found")
	}
}

func TestArchiveRegistry_Remove(t *testing.T) {
	reg := NewSecretArchiveRegistry()
	e := baseArchiveEntry
	_ = reg.Archive(&e)
	if !reg.Remove(e.Mount, e.Path) {
		t.Fatal("expected remove to return true")
	}
	if _, ok := reg.Get(e.Mount, e.Path); ok {
		t.Fatal("expected entry to be gone after remove")
	}
}

func TestArchiveRegistry_Remove_NotFound(t *testing.T) {
	reg := NewSecretArchiveRegistry()
	if reg.Remove("secret", "ghost") {
		t.Fatal("expected remove to return false for missing entry")
	}
}

func TestArchiveRegistry_Count(t *testing.T) {
	reg := NewSecretArchiveRegistry()
	if reg.Count() != 0 {
		t.Fatal("expected empty registry")
	}
	e := baseArchiveEntry
	_ = reg.Archive(&e)
	if reg.Count() != 1 {
		t.Fatalf("expected count 1, got %d", reg.Count())
	}
}

func TestArchiveRegistry_All(t *testing.T) {
	reg := NewSecretArchiveRegistry()
	e := baseArchiveEntry
	_ = reg.Archive(&e)
	all := reg.All()
	if len(all) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(all))
	}
}
