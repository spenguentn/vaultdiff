package vault

import (
	"testing"
)

func TestNewSecretTag_Valid(t *testing.T) {
	tag, err := NewSecretTag("secret", "myapp/config", TagSet{"env": "prod"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tag.Mount != "secret" || tag.Path != "myapp/config" {
		t.Errorf("unexpected fields: %+v", tag)
	}
}

func TestNewSecretTag_NilTagsDefaults(t *testing.T) {
	tag, err := NewSecretTag("secret", "myapp/config", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tag.Tags == nil {
		t.Error("expected non-nil TagSet")
	}
}

func TestSecretTag_Validate_MissingMount(t *testing.T) {
	_, err := NewSecretTag("", "myapp/config", nil)
	if err == nil {
		t.Error("expected error for empty mount")
	}
}

func TestSecretTag_Validate_MissingPath(t *testing.T) {
	_, err := NewSecretTag("secret", "", nil)
	if err == nil {
		t.Error("expected error for empty path")
	}
}

func TestSecretTag_FullPath(t *testing.T) {
	tag, _ := NewSecretTag("secret", "myapp/config", nil)
	if tag.FullPath() != "secret/myapp/config" {
		t.Errorf("unexpected full path: %s", tag.FullPath())
	}
}

func TestTagSet_SetAndGet(t *testing.T) {
	ts := make(TagSet)
	_ = ts.Set("owner", "team-a")
	v, ok := ts.Get("owner")
	if !ok || v != "team-a" {
		t.Errorf("expected owner=team-a, got %q ok=%v", v, ok)
	}
}

func TestTagSet_Set_EmptyKey(t *testing.T) {
	ts := make(TagSet)
	if err := ts.Set("", "value"); err == nil {
		t.Error("expected error for empty key")
	}
}

func TestTagSet_Merge(t *testing.T) {
	base := TagSet{"a": "1"}
	other := TagSet{"b": "2", "a": "overwritten"}
	base.Merge(other)
	if base["a"] != "overwritten" || base["b"] != "2" {
		t.Errorf("unexpected merged tags: %v", base)
	}
}

func TestTagRegistry_RegisterAndGet(t *testing.T) {
	reg := NewTagRegistry()
	tag, _ := NewSecretTag("secret", "myapp/config", TagSet{"env": "prod"})
	if err := reg.Register(tag); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got, ok := reg.Get("secret/myapp/config")
	if !ok || got.FullPath() != "secret/myapp/config" {
		t.Errorf("expected tag to be found, got ok=%v", ok)
	}
}

func TestTagRegistry_Delete(t *testing.T) {
	reg := NewTagRegistry()
	tag, _ := NewSecretTag("secret", "myapp/config", nil)
	_ = reg.Register(tag)
	if !reg.Delete("secret/myapp/config") {
		t.Error("expected delete to return true")
	}
	if reg.Len() != 0 {
		t.Errorf("expected empty registry, got len=%d", reg.Len())
	}
}

func TestTagRegistry_All(t *testing.T) {
	reg := NewTagRegistry()
	t1, _ := NewSecretTag("secret", "app/a", nil)
	t2, _ := NewSecretTag("secret", "app/b", nil)
	_ = reg.Register(t1)
	_ = reg.Register(t2)
	if len(reg.All()) != 2 {
		t.Errorf("expected 2 tags, got %d", len(reg.All()))
	}
}
