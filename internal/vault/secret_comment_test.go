package vault

import (
	"testing"
)

var baseComment = SecretComment{
	Mount:  "secret",
	Path:   "app/config",
	Author: "alice",
	Body:   "rotated credentials after incident",
}

func TestSecretComment_FullPath(t *testing.T) {
	c := baseComment
	if got := c.FullPath(); got != "secret/app/config" {
		t.Errorf("expected secret/app/config, got %s", got)
	}
}

func TestSecretComment_Validate_Valid(t *testing.T) {
	if err := baseComment.Validate(); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestSecretComment_Validate_MissingMount(t *testing.T) {
	c := baseComment
	c.Mount = ""
	if err := c.Validate(); err == nil {
		t.Error("expected error for missing mount")
	}
}

func TestSecretComment_Validate_MissingAuthor(t *testing.T) {
	c := baseComment
	c.Author = ""
	if err := c.Validate(); err == nil {
		t.Error("expected error for missing author")
	}
}

func TestSecretComment_Validate_MissingBody(t *testing.T) {
	c := baseComment
	c.Body = ""
	if err := c.Validate(); err == nil {
		t.Error("expected error for missing body")
	}
}

func TestNewSecretCommentRegistry_NotNil(t *testing.T) {
	r := NewSecretCommentRegistry()
	if r == nil {
		t.Fatal("expected non-nil registry")
	}
}

func TestCommentRegistry_Add_And_Get(t *testing.T) {
	r := NewSecretCommentRegistry()
	if err := r.Add(baseComment); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	list, ok := r.Get(baseComment.Mount, baseComment.Path)
	if !ok || len(list) != 1 {
		t.Errorf("expected 1 comment, got %d", len(list))
	}
}

func TestCommentRegistry_Add_Invalid(t *testing.T) {
	r := NewSecretCommentRegistry()
	c := baseComment
	c.Author = ""
	if err := r.Add(c); err == nil {
		t.Error("expected validation error")
	}
}

func TestCommentRegistry_Get_NotFound(t *testing.T) {
	r := NewSecretCommentRegistry()
	_, ok := r.Get("secret", "missing/path")
	if ok {
		t.Error("expected not found")
	}
}

func TestCommentRegistry_Remove(t *testing.T) {
	r := NewSecretCommentRegistry()
	_ = r.Add(baseComment)
	r.Remove(baseComment.Mount, baseComment.Path)
	_, ok := r.Get(baseComment.Mount, baseComment.Path)
	if ok {
		t.Error("expected comment to be removed")
	}
}

func TestCommentRegistry_Count(t *testing.T) {
	r := NewSecretCommentRegistry()
	_ = r.Add(baseComment)
	second := baseComment
	second.Path = "app/other"
	_ = r.Add(second)
	if r.Count() != 2 {
		t.Errorf("expected count 2, got %d", r.Count())
	}
}
