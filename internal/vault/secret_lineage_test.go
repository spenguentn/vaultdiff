package vault

import (
	"testing"
	"time"
)

func baseLineageEntry() LineageEntry {
	return LineageEntry{
		Mount:     "secret",
		Path:      "myapp/db",
		Version:   1,
		CreatedAt: time.Now(),
		CreatedBy: "admin",
	}
}

func TestLineageEntry_FullPath(t *testing.T) {
	e := baseLineageEntry()
	if got := e.FullPath(); got != "secret/myapp/db" {
		t.Errorf("expected 'secret/myapp/db', got %q", got)
	}
}

func TestLineageEntry_IsDeleted_False(t *testing.T) {
	e := baseLineageEntry()
	if e.IsDeleted() {
		t.Error("expected entry to not be deleted")
	}
}

func TestLineageEntry_IsDeleted_True(t *testing.T) {
	e := baseLineageEntry()
	now := time.Now()
	e.DeletedAt = &now
	if !e.IsDeleted() {
		t.Error("expected entry to be deleted")
	}
}

func TestLineageEntry_Validate_Valid(t *testing.T) {
	if err := baseLineageEntry().Validate(); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestLineageEntry_Validate_MissingMount(t *testing.T) {
	e := baseLineageEntry()
	e.Mount = ""
	if err := e.Validate(); err == nil {
		t.Error("expected error for missing mount")
	}
}

func TestLineageEntry_Validate_ZeroVersion(t *testing.T) {
	e := baseLineageEntry()
	e.Version = 0
	if err := e.Validate(); err == nil {
		t.Error("expected error for zero version")
	}
}

func TestNewSecretLineage_Empty(t *testing.T) {
	l := NewSecretLineage("secret", "myapp/db")
	if l.Len() != 0 {
		t.Errorf("expected 0 entries, got %d", l.Len())
	}
}

func TestSecretLineage_Add_Valid(t *testing.T) {
	l := NewSecretLineage("secret", "myapp/db")
	if err := l.Add(baseLineageEntry()); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if l.Len() != 1 {
		t.Errorf("expected 1 entry, got %d", l.Len())
	}
}

func TestSecretLineage_Add_Invalid(t *testing.T) {
	l := NewSecretLineage("secret", "myapp/db")
	e := baseLineageEntry()
	e.Path = ""
	if err := l.Add(e); err == nil {
		t.Error("expected error for invalid entry")
	}
}

func TestSecretLineage_Latest_Found(t *testing.T) {
	l := NewSecretLineage("secret", "myapp/db")
	_ = l.Add(baseLineageEntry())
	e2 := baseLineageEntry()
	e2.Version = 2
	_ = l.Add(e2)
	latest, ok := l.Latest()
	if !ok {
		t.Fatal("expected entry to be found")
	}
	if latest.Version != 2 {
		t.Errorf("expected version 2, got %d", latest.Version)
	}
}

func TestSecretLineage_Latest_Empty(t *testing.T) {
	l := NewSecretLineage("secret", "myapp/db")
	_, ok := l.Latest()
	if ok {
		t.Error("expected no entry on empty lineage")
	}
}

func TestSecretLineage_Entries_ReturnsCopy(t *testing.T) {
	l := NewSecretLineage("secret", "myapp/db")
	_ = l.Add(baseLineageEntry())
	entries := l.Entries()
	entries[0].Version = 99
	latest, _ := l.Latest()
	if latest.Version == 99 {
		t.Error("Entries should return a copy, not a reference")
	}
}
