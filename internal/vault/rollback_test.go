package vault

import (
	"context"
	"errors"
	"testing"
	"time"
)

// TestRollbackResult_IsSuccess verifies the success predicate.
func TestRollbackResult_IsSuccess(t *testing.T) {
	r := RollbackResult{RolledBack: true, RolledAt: time.Now()}
	if !r.IsSuccess() {
		t.Fatal("expected IsSuccess true")
	}
}

// TestRollbackResult_IsSuccess_WithErr returns false when Err is set.
func TestRollbackResult_IsSuccess_WithErr(t *testing.T) {
	r := RollbackResult{RolledBack: true, Err: errors.New("oops")}
	if r.IsSuccess() {
		t.Fatal("expected IsSuccess false when Err is set")
	}
}

// TestRollbackResult_String_OK checks the success string format.
func TestRollbackResult_String_OK(t *testing.T) {
	r := RollbackResult{
		Mount:      "secret",
		Path:       "app/db",
		Version:    3,
		RolledBack: true,
		RolledAt:   time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
	}
	got := r.String()
	if got != "rollback secret/app/db@v3 OK at 2024-01-15T10:00:00Z" {
		t.Fatalf("unexpected string: %s", got)
	}
}

// TestRollbackResult_String_Err checks the failure string format.
func TestRollbackResult_String_Err(t *testing.T) {
	r := RollbackResult{
		Mount:   "secret",
		Path:    "app/db",
		Version: 2,
		Err:     errors.New("permission denied"),
	}
	got := r.String()
	expected := "rollback secret/app/db@v2 FAILED: permission denied"
	if got != expected {
		t.Fatalf("got %q, want %q", got, expected)
	}
}

// TestNewRollbacker_NilPanics ensures nil client panics.
func TestNewRollbacker_NilPanics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for nil client")
		}
	}()
	NewRollbacker(nil)
}

// TestRollback_InvalidVersion checks that version < 1 returns an error.
func TestRollback_InvalidVersion(t *testing.T) {
	// We cannot create a real Client easily in unit tests, so we test the
	// validation path by inspecting the returned error from a zero-version req.
	// Build a minimal client config to avoid nil-pointer in NewClient.
	cfg := Config{Address: "http://127.0.0.1:8200", Token: "test"}
	c, err := NewClient(cfg)
	if err != nil {
		t.Skipf("skipping: cannot create client: %v", err)
	}
	rb := NewRollbacker(c)
	res := rb.Rollback(context.Background(), RollbackRequest{
		Mount: "secret", Path: "app/db", Version: 0,
	})
	if res.Err == nil {
		t.Fatal("expected error for version 0")
	}
	if res.IsSuccess() {
		t.Fatal("expected IsSuccess false")
	}
}
