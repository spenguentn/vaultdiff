package vault

import (
	"errors"
	"sync"
	"time"
)

// LockState represents whether a secret path is locked or unlocked.
type LockState int

const (
	LockStateUnlocked LockState = iota
	LockStateLocked
)

// SecretLock represents a lock held on a specific secret path.
type SecretLock struct {
	Mount     string
	Path      string
	LockedAt  time.Time
	LockedBy  string
	ExpiresAt time.Time
}

// IsExpired returns true if the lock TTL has elapsed.
func (l SecretLock) IsExpired() bool {
	if l.ExpiresAt.IsZero() {
		return false
	}
	return time.Now().After(l.ExpiresAt)
}

// FullPath returns the mount-qualified path for the lock.
func (l SecretLock) FullPath() string {
	if l.Mount == "" {
		return l.Path
	}
	return l.Mount + "/" + l.Path
}

// Validate checks that the lock has the required fields.
func (l SecretLock) Validate() error {
	if l.Mount == "" {
		return errors.New("secret lock: mount is required")
	}
	if l.Path == "" {
		return errors.New("secret lock: path is required")
	}
	if l.LockedBy == "" {
		return errors.New("secret lock: locked_by is required")
	}
	return nil
}

// SecretLockRegistry manages in-memory locks on secret paths.
type SecretLockRegistry struct {
	mu    sync.RWMutex
	locks map[string]SecretLock
}

// NewSecretLockRegistry returns an initialised SecretLockRegistry.
func NewSecretLockRegistry() *SecretLockRegistry {
	return &SecretLockRegistry{
		locks: make(map[string]SecretLock),
	}
}

// Lock acquires a lock for the given secret. Returns an error if the path is
// already locked by someone else and the lock has not expired.
func (r *SecretLockRegistry) Lock(l SecretLock) error {
	if err := l.Validate(); err != nil {
		return err
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	existing, ok := r.locks[l.FullPath()]
	if ok && !existing.IsExpired() {
		return errors.New("secret lock: path already locked by " + existing.LockedBy)
	}
	l.LockedAt = time.Now()
	r.locks[l.FullPath()] = l
	return nil
}

// Unlock releases the lock on a path. Only the original owner may unlock.
func (r *SecretLockRegistry) Unlock(mount, path, owner string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	key := mount + "/" + path
	existing, ok := r.locks[key]
	if !ok {
		return errors.New("secret lock: no lock found for path")
	}
	if existing.LockedBy != owner {
		return errors.New("secret lock: unlock denied, owner mismatch")
	}
	delete(r.locks, key)
	return nil
}

// Get returns the current lock for a path, if any.
func (r *SecretLockRegistry) Get(mount, path string) (SecretLock, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	l, ok := r.locks[mount+"/"+path]
	return l, ok
}

// IsLocked returns true when the path has an active, non-expired lock.
func (r *SecretLockRegistry) IsLocked(mount, path string) bool {
	l, ok := r.Get(mount, path)
	return ok && !l.IsExpired()
}
