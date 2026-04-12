package vault

import (
	"fmt"
	"sync"
	"time"
)

func alertKey(mount, path string) string {
	return fmt.Sprintf("%s/%s", mount, path)
}

// SecretAlertRegistry stores and retrieves secret alerts keyed by mount+path.
type SecretAlertRegistry struct {
	mu     sync.RWMutex
	alerts map[string][]SecretAlert
}

// NewSecretAlertRegistry returns an initialised SecretAlertRegistry.
func NewSecretAlertRegistry() *SecretAlertRegistry {
	return &SecretAlertRegistry{
		alerts: make(map[string][]SecretAlert),
	}
}

// Record validates and stores a new alert, stamping Triggered if zero.
func (r *SecretAlertRegistry) Record(a SecretAlert) error {
	if a.Triggered.IsZero() {
		a.Triggered = time.Now().UTC()
	}
	if err := a.Validate(); err != nil {
		return err
	}
	k := alertKey(a.Mount, a.Path)
	r.mu.Lock()
	defer r.mu.Unlock()
	r.alerts[k] = append(r.alerts[k], a)
	return nil
}

// Get returns all alerts recorded for the given mount and path.
func (r *SecretAlertRegistry) Get(mount, path string) ([]SecretAlert, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	v, ok := r.alerts[alertKey(mount, path)]
	return v, ok
}

// Clear removes all alerts for the given mount and path.
func (r *SecretAlertRegistry) Clear(mount, path string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.alerts, alertKey(mount, path))
}

// All returns a flat slice of every recorded alert across all paths.
func (r *SecretAlertRegistry) All() []SecretAlert {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var out []SecretAlert
	for _, list := range r.alerts {
		out = append(out, list...)
	}
	return out
}
