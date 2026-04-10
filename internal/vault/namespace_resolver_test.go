package vault

import (
	"os"
	"testing"
)

func TestResolveNamespace_Direct(t *testing.T) {
	res, err := ResolveNamespace("team/platform")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Source != NamespaceSourceDirect {
		t.Errorf("expected source 'direct', got %q", res.Source)
	}
	if res.Namespace.Path != "team/platform" {
		t.Errorf("expected path 'team/platform', got %q", res.Namespace.Path)
	}
}

func TestResolveNamespace_DirectTrimsSlashes(t *testing.T) {
	res, err := ResolveNamespace("/admin/")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Namespace.Path != "admin" {
		t.Errorf("expected 'admin', got %q", res.Namespace.Path)
	}
}

func TestResolveNamespace_EnvVar(t *testing.T) {
	t.Setenv(envVaultNamespace, "ops/staging")
	res, err := ResolveNamespace("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Source != NamespaceSourceEnv {
		t.Errorf("expected source 'env', got %q", res.Source)
	}
	if res.Namespace.Path != "ops/staging" {
		t.Errorf("expected 'ops/staging', got %q", res.Namespace.Path)
	}
}

func TestResolveNamespace_DirectTakesPrecedenceOverEnv(t *testing.T) {
	t.Setenv(envVaultNamespace, "ops/staging")
	res, err := ResolveNamespace("team/platform")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Source != NamespaceSourceDirect {
		t.Errorf("expected source 'direct', got %q", res.Source)
	}
}

func TestResolveNamespace_Root(t *testing.T) {
	os.Unsetenv(envVaultNamespace)
	res, err := ResolveNamespace("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Source != NamespaceSourceRoot {
		t.Errorf("expected source 'root', got %q", res.Source)
	}
	if !res.Namespace.IsRoot() {
		t.Error("expected root namespace")
	}
}

func TestResolveNamespace_InvalidDirect(t *testing.T) {
	_, err := ResolveNamespace("bad namespace")
	if err == nil {
		t.Error("expected error for namespace with spaces")
	}
}
