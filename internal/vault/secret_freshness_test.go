package vault

import (
	"testing"
	"time"
)

func TestIsValidFreshnessStatus_Known(t *testing.T) {
	for _, s := range []FreshnessStatus{FreshnessFresh, FreshnessStale, FreshnessExpired, FreshnessUnknown} {
		if !IsValidFreshnessStatus(s) {
			t.Errorf("expected %q to be valid", s)
		}
	}
}

func TestIsValidFreshnessStatus_Unknown(t *testing.T) {
	if IsValidFreshnessStatus("bogus") {
		t.Error("expected 'bogus' to be invalid")
	}
}

func TestSecretFreshness_FullPath(t *testing.T) {
	f := &SecretFreshness{Mount: "secret", Path: "app/db"}
	if got := f.FullPath(); got != "secret/app/db" {
		t.Errorf("unexpected FullPath: %q", got)
	}
}

func TestSecretFreshness_Validate_Valid(t *testing.T) {
	f := &SecretFreshness{
		Mount:       "secret",
		Path:        "app/key",
		LastUpdated: time.Now().Add(-time.Hour),
		MaxAge:      24 * time.Hour,
		Status:      FreshnessFresh,
	}
	if err := f.Validate(); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestSecretFreshness_Validate_MissingMount(t *testing.T) {
	f := &SecretFreshness{Path: "app/key", MaxAge: time.Hour, LastUpdated: time.Now(), Status: FreshnessFresh}
	if err := f.Validate(); err == nil {
		t.Error("expected error for missing mount")
	}
}

func TestSecretFreshness_Validate_MissingPath(t *testing.T) {
	f := &SecretFreshness{Mount: "secret", MaxAge: time.Hour, LastUpdated: time.Now(), Status: FreshnessFresh}
	if err := f.Validate(); err == nil {
		t.Error("expected error for missing path")
	}
}

func TestSecretFreshness_Validate_ZeroMaxAge(t *testing.T) {
	f := &SecretFreshness{Mount: "secret", Path: "app/key", LastUpdated: time.Now(), Status: FreshnessFresh}
	if err := f.Validate(); err == nil {
		t.Error("expected error for zero max_age")
	}
}

func TestComputeFreshness_Fresh(t *testing.T) {
	lastUpdated := time.Now().Add(-1 * time.Hour)
	maxAge := 24 * time.Hour
	f, err := ComputeFreshness("secret", "app/key", lastUpdated, maxAge)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f.Status != FreshnessFresh {
		t.Errorf("expected fresh, got %q", f.Status)
	}
}

func TestComputeFreshness_Stale(t *testing.T) {
	lastUpdated := time.Now().Add(-20 * time.Hour)
	maxAge := 24 * time.Hour
	f, err := ComputeFreshness("secret", "app/key", lastUpdated, maxAge)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f.Status != FreshnessStale {
		t.Errorf("expected stale, got %q", f.Status)
	}
}

func TestComputeFreshness_Expired(t *testing.T) {
	lastUpdated := time.Now().Add(-48 * time.Hour)
	maxAge := 24 * time.Hour
	f, err := ComputeFreshness("secret", "app/key", lastUpdated, maxAge)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f.Status != FreshnessExpired {
		t.Errorf("expected expired, got %q", f.Status)
	}
}

func TestComputeFreshness_MissingMount(t *testing.T) {
	_, err := ComputeFreshness("", "app/key", time.Now(), time.Hour)
	if err == nil {
		t.Error("expected error for missing mount")
	}
}

func TestComputeFreshness_ZeroMaxAge(t *testing.T) {
	_, err := ComputeFreshness("secret", "app/key", time.Now(), 0)
	if err == nil {
		t.Error("expected error for zero max_age")
	}
}
