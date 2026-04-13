package vault

import (
	"fmt"
	"sync"
	"time"
)

func pinKey(mount, path string) string {
	return fmt.Sprintf("%s::%s", mount, path)
}

// SecretPinRegistry stores and manages pinned secret versions.
type SecretPinRegistry struct {
	mu   sync.RWMutex
	pins map[string]SecretPin
}

// NewSecretPinRegistry returns an initialised SecretPinRegistry.
func NewSecretPinRegistry() *SecretPinRegistry {
	return &SecretPinRegistry{
		pins: make(map[string]SecretPin),
	}
}

// Pin stores a pin after validation, setting PinnedAt if not already set.
func (r *SecretPinRegistry) Pin(p SecretPin) error {
	if err := p.Validate(); err != nil {
		return err
	}
	if p.PinnedAt.IsZero() {
		p.PinnedAt = time.Now().UTC()
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	r.pins[pinKey(p.Mount, p.Path)] = p
	return nil
}

// Get retrieves a pin by mount and path.
func (r *SecretPinRegistry) Get(mount, path string) (SecretPin, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	p, ok := r.pins[pinKey(mount, path)]
	return p, ok
}

// Unpin removes a pin for the given mount and path.
func (r *SecretPinRegistry) Unpin(mount, path string) bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	key := pinKey(mount, path)
	_, ok := r.pins[key]
	if ok {
		delete(r.pins, key)
	}
	return ok
}

// IsPinned reports whether a non-expired pin exists for the path.
func (r *SecretPinRegistry) IsPinned(mount, path string) bool {
	p, ok := r.Get(mount, path)
	if !ok {
		return false
	}
	return !p.IsExpired()
}

// All returns a snapshot of all current pins.
func (r *SecretPinRegistry) All() []SecretPin {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]SecretPin, 0, len(r.pins))
	for _, p := range r.pins {
		out = append(out, p)
	}
	return out
}
