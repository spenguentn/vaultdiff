package vault

import (
	"testing"
	"time"
)

func baseLock() SecretLock {
	return SecretLock{
		Mount:    "secret",
		Path:     "app/config",
		LockedBy: "alice",
	}
}

func TestSecretLock_Validate_Valid(t *testing.T) {
	if err := baseLock().Validate(); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestSecretLock_Validate_MissingMount(t *testing.T) {
	l := baseLock()
	l.Mount = ""
	if err := l.Validate(); err == nil {
		t.Fatal("expected error for missing mount")
	}
}

func TestSecretLock_Validate_MissingPath(t *testing.T) {
	l := baseLock()
	l.Path = ""
	if err := l.Validate(); err == nil {
		t.Fatal("expected error for missing path")
	}
}

func TestSecretLock_Validate_MissingOwner(t *testing.T) {
	l := baseLock()
	l.LockedBy = ""
	if err := l.Validate(); err == nil {
		t.Fatal("expected error for missing locked_by")
	}
}

func TestSecretLock_FullPath(t *testing.T) {
	l := baseLock()
	if got := l.FullPath(); got != "secret/app/config" {
		t.Fatalf("expected secret/app/config, got %s", got)
	}
}

func TestSecretLock_IsExpired_NoExpiry(t *testing.T) {
	l := baseLock()
	if l.IsExpired() {
		t.Fatal("expected not expired when ExpiresAt is zero")
	}
}

func TestSecretLock_IsExpired_Future(t *testing.T) {
	l := baseLock()
	l.ExpiresAt = time.Now().Add(time.Hour)
	if l.IsExpired() {
		t.Fatal("expected not expired for future time")
	}
}

func TestSecretLock_IsExpired_Past(t *testing.T) {
	l := baseLock()
	l.ExpiresAt = time.Now().Add(-time.Second)
	if !l.IsExpired() {
		t.Fatal("expected expired for past time")
	}
}

func TestSecretLockRegistry_LockAndGet(t *testing.T) {
	reg := NewSecretLockRegistry()
	l := baseLock()
	if err := reg.Lock(l); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got, ok := reg.Get(l.Mount, l.Path)
	if !ok {
		t.Fatal("expected lock to be present")
	}
	if got.LockedBy != "alice" {
		t.Fatalf("expected alice, got %s", got.LockedBy)
	}
}

func TestSecretLockRegistry_Lock_AlreadyLocked(t *testing.T) {
	reg := NewSecretLockRegistry()
	_ = reg.Lock(baseLock())
	second := baseLock()
	second.LockedBy = "bob"
	if err := reg.Lock(second); err == nil {
		t.Fatal("expected error when path already locked")
	}
}

func TestSecretLockRegistry_Lock_ExpiredAllowsRelock(t *testing.T) {
	reg := NewSecretLockRegistry()
	l := baseLock()
	l.ExpiresAt = time.Now().Add(-time.Second)
	_ = reg.Lock(l)
	second := baseLock()
	second.LockedBy = "bob"
	if err := reg.Lock(second); err != nil {
		t.Fatalf("expected relocking after expiry, got: %v", err)
	}
}

func TestSecretLockRegistry_Unlock_Valid(t *testing.T) {
	reg := NewSecretLockRegistry()
	l := baseLock()
	_ = reg.Lock(l)
	if err := reg.Unlock(l.Mount, l.Path, "alice"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if reg.IsLocked(l.Mount, l.Path) {
		t.Fatal("expected path to be unlocked")
	}
}

func TestSecretLockRegistry_Unlock_WrongOwner(t *testing.T) {
	reg := NewSecretLockRegistry()
	_ = reg.Lock(baseLock())
	if err := reg.Unlock("secret", "app/config", "bob"); err == nil {
		t.Fatal("expected error for wrong owner")
	}
}

func TestSecretLockRegistry_Unlock_NotFound(t *testing.T) {
	reg := NewSecretLockRegistry()
	if err := reg.Unlock("secret", "app/config", "alice"); err == nil {
		t.Fatal("expected error for missing lock")
	}
}

func TestSecretLockRegistry_IsLocked_False(t *testing.T) {
	reg := NewSecretLockRegistry()
	if reg.IsLocked("secret", "app/config") {
		t.Fatal("expected not locked on empty registry")
	}
}
