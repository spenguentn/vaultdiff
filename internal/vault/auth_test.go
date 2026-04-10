package vault

import (
	"testing"
)

func TestAuthConfig_Validate_Token(t *testing.T) {
	cfg := AuthConfig{Method: AuthToken, Token: "s.abc123"}
	if err := cfg.Validate(); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestAuthConfig_Validate_TokenEmpty(t *testing.T) {
	cfg := AuthConfig{Method: AuthToken}
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error for empty token")
	}
}

func TestAuthConfig_Validate_AppRole(t *testing.T) {
	cfg := AuthConfig{Method: AuthAppRole, RoleID: "role", SecretID: "secret"}
	if err := cfg.Validate(); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestAuthConfig_Validate_AppRole_MissingRoleID(t *testing.T) {
	cfg := AuthConfig{Method: AuthAppRole, SecretID: "secret"}
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error for missing role_id")
	}
}

func TestAuthConfig_Validate_AppRole_MissingSecretID(t *testing.T) {
	cfg := AuthConfig{Method: AuthAppRole, RoleID: "role"}
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error for missing secret_id")
	}
}

func TestAuthConfig_Validate_Kubernetes(t *testing.T) {
	cfg := AuthConfig{Method: AuthKubernetes, Role: "my-role", JWT: "jwt-token"}
	if err := cfg.Validate(); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestAuthConfig_Validate_UnsupportedMethod(t *testing.T) {
	cfg := AuthConfig{Method: "ldap"}
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error for unsupported method")
	}
}

func TestAuthConfig_EffectiveMountPath_Default(t *testing.T) {
	cfg := AuthConfig{Method: AuthAppRole}
	if got := cfg.EffectiveMountPath(); got != "approle" {
		t.Fatalf("expected approle, got %q", got)
	}
}

func TestAuthConfig_EffectiveMountPath_Custom(t *testing.T) {
	cfg := AuthConfig{Method: AuthAppRole, MountPath: "/custom-approle/"}
	if got := cfg.EffectiveMountPath(); got != "custom-approle" {
		t.Fatalf("expected custom-approle, got %q", got)
	}
}

func TestAuthConfig_LoginPath(t *testing.T) {
	cfg := AuthConfig{Method: AuthKubernetes, Role: "r", JWT: "j"}
	if got := cfg.LoginPath(); got != "auth/kubernetes/login" {
		t.Fatalf("expected auth/kubernetes/login, got %q", got)
	}
}
