package vault

import (
	"testing"
)

func TestSecretMigration_FullSource(t *testing.T) {
	m := SecretMigration{SourceMount: "secret", SourcePath: "app/db"}
	if got := m.FullSource(); got != "secret/app/db" {
		t.Errorf("expected secret/app/db, got %s", got)
	}
}

func TestSecretMigration_FullDest(t *testing.T) {
	m := SecretMigration{DestMount: "kv", DestPath: "prod/db"}
	if got := m.FullDest(); got != "kv/prod/db" {
		t.Errorf("expected kv/prod/db, got %s", got)
	}
}

func TestSecretMigration_IsTerminal_Completed(t *testing.T) {
	m := SecretMigration{Status: MigrationCompleted}
	if !m.IsTerminal() {
		t.Error("expected completed to be terminal")
	}
}

func TestSecretMigration_IsTerminal_Pending(t *testing.T) {
	m := SecretMigration{Status: MigrationPending}
	if m.IsTerminal() {
		t.Error("expected pending to not be terminal")
	}
}

func TestSecretMigration_Validate_Valid(t *testing.T) {
	m := SecretMigration{
		SourceMount: "secret", SourcePath: "app/key",
		DestMount: "kv", DestPath: "prod/key",
		InitiatedBy: "alice", Status: MigrationPending,
	}
	if err := m.Validate(); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestSecretMigration_Validate_MissingSourceMount(t *testing.T) {
	m := SecretMigration{
		SourcePath: "app/key", DestMount: "kv", DestPath: "prod/key",
		InitiatedBy: "alice", Status: MigrationPending,
	}
	if err := m.Validate(); err == nil {
		t.Error("expected error for missing source_mount")
	}
}

func TestSecretMigration_Validate_InvalidStatus(t *testing.T) {
	m := SecretMigration{
		SourceMount: "secret", SourcePath: "app/key",
		DestMount: "kv", DestPath: "prod/key",
		InitiatedBy: "alice", Status: "unknown",
	}
	if err := m.Validate(); err == nil {
		t.Error("expected error for invalid status")
	}
}

func TestIsValidMigrationStatus_Known(t *testing.T) {
	for _, s := range []MigrationStatus{MigrationPending, MigrationRunning, MigrationCompleted, MigrationFailed} {
		if !IsValidMigrationStatus(s) {
			t.Errorf("expected %q to be valid", s)
		}
	}
}

func TestIsValidMigrationStatus_Unknown(t *testing.T) {
	if IsValidMigrationStatus("archived") {
		t.Error("expected 'archived' to be invalid")
	}
}

func TestMigrationRegistry_Submit_And_Get(t *testing.T) {
	r := NewSecretMigrationRegistry()
	m := SecretMigration{
		SourceMount: "secret", SourcePath: "svc/token",
		DestMount: "kv", DestPath: "prod/token",
		InitiatedBy: "bob", Status: MigrationPending,
	}
	result, err := r.Submit(m)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.ID == "" {
		t.Error("expected ID to be set")
	}
	got, ok := r.Get("secret", "svc/token")
	if !ok {
		t.Fatal("expected migration to be found")
	}
	if got.InitiatedBy != "bob" {
		t.Errorf("expected bob, got %s", got.InitiatedBy)
	}
}

func TestMigrationRegistry_Submit_DuplicateActive(t *testing.T) {
	r := NewSecretMigrationRegistry()
	m := SecretMigration{
		SourceMount: "secret", SourcePath: "svc/token",
		DestMount: "kv", DestPath: "prod/token",
		InitiatedBy: "bob", Status: MigrationPending,
	}
	if _, err := r.Submit(m); err != nil {
		t.Fatalf("first submit failed: %v", err)
	}
	if _, err := r.Submit(m); err == nil {
		t.Error("expected error for duplicate active migration")
	}
}

func TestMigrationRegistry_UpdateStatus(t *testing.T) {
	r := NewSecretMigrationRegistry()
	m := SecretMigration{
		SourceMount: "secret", SourcePath: "svc/key",
		DestMount: "kv", DestPath: "prod/key",
		InitiatedBy: "carol", Status: MigrationPending,
	}
	result, _ := r.Submit(m)
	if err := r.UpdateStatus(result.ID, MigrationCompleted); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	got, _ := r.Get("secret", "svc/key")
	if got.Status != MigrationCompleted {
		t.Errorf("expected completed, got %s", got.Status)
	}
}

func TestMigrationRegistry_Get_NotFound(t *testing.T) {
	r := NewSecretMigrationRegistry()
	if _, ok := r.Get("secret", "missing"); ok {
		t.Error("expected not found")
	}
}
