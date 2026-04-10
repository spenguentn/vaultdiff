package vault

import (
	"testing"
)

func TestNewNamespace_TrimsSlashes(t *testing.T) {
	n := NewNamespace("/team/platform/")
	if n.Path != "team/platform" {
		t.Errorf("expected 'team/platform', got %q", n.Path)
	}
}

func TestNamespace_IsRoot_Empty(t *testing.T) {
	n := NewNamespace("")
	if !n.IsRoot() {
		t.Error("expected IsRoot() == true for empty path")
	}
}

func TestNamespace_IsRoot_NonEmpty(t *testing.T) {
	n := NewNamespace("admin")
	if n.IsRoot() {
		t.Error("expected IsRoot() == false for non-empty path")
	}
}

func TestNamespace_String_Root(t *testing.T) {
	n := NewNamespace("")
	if n.String() != "root" {
		t.Errorf("expected 'root', got %q", n.String())
	}
}

func TestNamespace_String_NonRoot(t *testing.T) {
	n := NewNamespace("team/platform")
	if n.String() != "team/platform" {
		t.Errorf("expected 'team/platform', got %q", n.String())
	}
}

func TestNamespace_Child_FromRoot(t *testing.T) {
	n := NewNamespace("")
	child := n.Child("admin")
	if child.Path != "admin" {
		t.Errorf("expected 'admin', got %q", child.Path)
	}
}

func TestNamespace_Child_FromNonRoot(t *testing.T) {
	n := NewNamespace("team")
	child := n.Child("platform")
	if child.Path != "team/platform" {
		t.Errorf("expected 'team/platform', got %q", child.Path)
	}
}

func TestNamespace_Child_EmptySegment(t *testing.T) {
	n := NewNamespace("team")
	child := n.Child("")
	if child.Path != "team" {
		t.Errorf("expected unchanged 'team', got %q", child.Path)
	}
}

func TestNamespace_Validate_Valid(t *testing.T) {
	n := NewNamespace("team/platform")
	if err := n.Validate(); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestNamespace_Validate_WithSpaces(t *testing.T) {
	n := NewNamespace("team platform")
	if err := n.Validate(); err == nil {
		t.Error("expected error for namespace with spaces")
	}
}

func TestNamespace_HeaderValue_Root(t *testing.T) {
	n := NewNamespace("")
	if n.HeaderValue() != "" {
		t.Errorf("expected empty header value for root, got %q", n.HeaderValue())
	}
}

func TestNamespace_HeaderValue_NonRoot(t *testing.T) {
	n := NewNamespace("team/platform")
	if n.HeaderValue() != "team/platform" {
		t.Errorf("expected 'team/platform', got %q", n.HeaderValue())
	}
}
