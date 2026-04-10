package vault

import (
	"testing"
)

func TestNewSecretPath_TrimsSlashes(t *testing.T) {
	p := NewSecretPath("/secret/", "/myapp/config/")
	if p.Mount != "secret" {
		t.Errorf("expected mount %q, got %q", "secret", p.Mount)
	}
	if p.SubPath != "myapp/config" {
		t.Errorf("expected sub-path %q, got %q", "myapp/config", p.SubPath)
	}
}

func TestSecretPath_DataPath(t *testing.T) {
	p := NewSecretPath("secret", "myapp/config")
	want := "secret/data/myapp/config"
	if got := p.DataPath(); got != want {
		t.Errorf("DataPath() = %q, want %q", got, want)
	}
}

func TestSecretPath_MetadataPath(t *testing.T) {
	p := NewSecretPath("secret", "myapp/config")
	want := "secret/metadata/myapp/config"
	if got := p.MetadataPath(); got != want {
		t.Errorf("MetadataPath() = %q, want %q", got, want)
	}
}

func TestSecretPath_String(t *testing.T) {
	p := NewSecretPath("secret", "myapp/config")
	want := "secret/myapp/config"
	if got := p.String(); got != want {
		t.Errorf("String() = %q, want %q", got, want)
	}
}

func TestSecretPath_Validate_Valid(t *testing.T) {
	p := NewSecretPath("secret", "myapp/config")
	if err := p.Validate(); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestSecretPath_Validate_EmptyMount(t *testing.T) {
	p := NewSecretPath("", "myapp/config")
	if err := p.Validate(); err == nil {
		t.Error("expected error for empty mount, got nil")
	}
}

func TestSecretPath_Validate_EmptySubPath(t *testing.T) {
	p := NewSecretPath("secret", "")
	if err := p.Validate(); err == nil {
		t.Error("expected error for empty sub-path, got nil")
	}
}

func TestParseSecretPath_Valid(t *testing.T) {
	p, err := ParseSecretPath("secret/myapp/config")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.Mount != "secret" {
		t.Errorf("mount = %q, want %q", p.Mount, "secret")
	}
	if p.SubPath != "myapp/config" {
		t.Errorf("sub-path = %q, want %q", p.SubPath, "myapp/config")
	}
}

func TestParseSecretPath_NoSlash(t *testing.T) {
	_, err := ParseSecretPath("secretonly")
	if err == nil {
		t.Error("expected error for path without slash, got nil")
	}
}

func TestParseSecretPath_LeadingSlashTrimmed(t *testing.T) {
	p, err := ParseSecretPath("/secret/myapp")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.Mount != "secret" {
		t.Errorf("mount = %q, want %q", p.Mount, "secret")
	}
	if p.SubPath != "myapp" {
		t.Errorf("sub-path = %q, want %q", p.SubPath, "myapp")
	}
}
