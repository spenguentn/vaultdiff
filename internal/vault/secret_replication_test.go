package vault

import (
	"testing"
)

func baseReplicationPolicy() ReplicationPolicy {
	return ReplicationPolicy{
		Mount:     "secret",
		Path:      "app/config",
		Targets:   []string{"us-east", "eu-west"},
		Mode:      ReplicationModeSync,
		CreatedBy: "ops-team",
	}
}

func TestIsValidReplicationMode_Known(t *testing.T) {
	for _, m := range []ReplicationMode{ReplicationModeSync, ReplicationModeAsync} {
		if !IsValidReplicationMode(m) {
			t.Errorf("expected %q to be valid", m)
		}
	}
}

func TestIsValidReplicationMode_Unknown(t *testing.T) {
	if IsValidReplicationMode("push") {
		t.Error("expected unknown mode to be invalid")
	}
}

func TestReplicationPolicy_FullPath(t *testing.T) {
	p := baseReplicationPolicy()
	if got := p.FullPath(); got != "secret/app/config" {
		t.Errorf("unexpected full path: %s", got)
	}
}

func TestReplicationPolicy_Validate_Valid(t *testing.T) {
	if err := baseReplicationPolicy().Validate(); err != nil {
		t.Errorf("expected valid, got: %v", err)
	}
}

func TestReplicationPolicy_Validate_MissingMount(t *testing.T) {
	p := baseReplicationPolicy()
	p.Mount = ""
	if err := p.Validate(); err == nil {
		t.Error("expected error for missing mount")
	}
}

func TestReplicationPolicy_Validate_MissingPath(t *testing.T) {
	p := baseReplicationPolicy()
	p.Path = ""
	if err := p.Validate(); err == nil {
		t.Error("expected error for missing path")
	}
}

func TestReplicationPolicy_Validate_NoTargets(t *testing.T) {
	p := baseReplicationPolicy()
	p.Targets = nil
	if err := p.Validate(); err == nil {
		t.Error("expected error for empty targets")
	}
}

func TestReplicationPolicy_Validate_InvalidMode(t *testing.T) {
	p := baseReplicationPolicy()
	p.Mode = "broadcast"
	if err := p.Validate(); err == nil {
		t.Error("expected error for invalid mode")
	}
}

func TestReplicationPolicy_Validate_MissingCreatedBy(t *testing.T) {
	p := baseReplicationPolicy()
	p.CreatedBy = ""
	if err := p.Validate(); err == nil {
		t.Error("expected error for missing created_by")
	}
}

func TestNewSecretReplicationRegistry_NotNil(t *testing.T) {
	if NewSecretReplicationRegistry() == nil {
		t.Error("expected non-nil registry")
	}
}

func TestReplicationRegistry_Set_And_Get(t *testing.T) {
	reg := NewSecretReplicationRegistry()
	p := baseReplicationPolicy()
	if err := reg.Set(p); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got, ok := reg.Get(p.Mount, p.Path)
	if !ok {
		t.Fatal("expected policy to be found")
	}
	if got.CreatedBy != p.CreatedBy {
		t.Errorf("got created_by %q, want %q", got.CreatedBy, p.CreatedBy)
	}
}

func TestReplicationRegistry_Set_Invalid(t *testing.T) {
	reg := NewSecretReplicationRegistry()
	p := baseReplicationPolicy()
	p.Mount = ""
	if err := reg.Set(p); err == nil {
		t.Error("expected validation error")
	}
}

func TestReplicationRegistry_Get_NotFound(t *testing.T) {
	reg := NewSecretReplicationRegistry()
	_, ok := reg.Get("secret", "missing")
	if ok {
		t.Error("expected not found")
	}
}

func TestReplicationRegistry_Remove(t *testing.T) {
	reg := NewSecretReplicationRegistry()
	p := baseReplicationPolicy()
	_ = reg.Set(p)
	reg.Remove(p.Mount, p.Path)
	_, ok := reg.Get(p.Mount, p.Path)
	if ok {
		t.Error("expected policy to be removed")
	}
}

func TestReplicationRegistry_All(t *testing.T) {
	reg := NewSecretReplicationRegistry()
	_ = reg.Set(baseReplicationPolicy())
	if len(reg.All()) != 1 {
		t.Errorf("expected 1 policy, got %d", len(reg.All()))
	}
}
