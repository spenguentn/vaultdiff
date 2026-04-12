package vault

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

func reminderKey(mount, path string) string {
	return fmt.Sprintf("%s/%s", mount, path)
}

// SecretReminderRegistry stores and retrieves SecretReminder entries.
type SecretReminderRegistry struct {
	mu      sync.RWMutex
	entries map[string]SecretReminder
}

// NewSecretReminderRegistry returns an initialised SecretReminderRegistry.
func NewSecretReminderRegistry() *SecretReminderRegistry {
	return &SecretReminderRegistry{
		entries: make(map[string]SecretReminder),
	}
}

// Set validates and stores a reminder, stamping CreatedAt when zero.
func (reg *SecretReminderRegistry) Set(r SecretReminder) error {
	if err := r.Validate(); err != nil {
		return err
	}
	if r.CreatedAt.IsZero() {
		r.CreatedAt = time.Now().UTC()
	}
	reg.mu.Lock()
	defer reg.mu.Unlock()
	reg.entries[reminderKey(r.Mount, r.Path)] = r
	return nil
}

// Get retrieves the reminder for the given mount and path.
func (reg *SecretReminderRegistry) Get(mount, path string) (SecretReminder, error) {
	reg.mu.RLock()
	defer reg.mu.RUnlock()
	r, ok := reg.entries[reminderKey(mount, path)]
	if !ok {
		return SecretReminder{}, errors.New("reminder: not found")
	}
	return r, nil
}

// Remove deletes the reminder for the given mount and path.
func (reg *SecretReminderRegistry) Remove(mount, path string) error {
	reg.mu.Lock()
	defer reg.mu.Unlock()
	k := reminderKey(mount, path)
	if _, ok := reg.entries[k]; !ok {
		return errors.New("reminder: not found")
	}
	delete(reg.entries, k)
	return nil
}

// Due returns all reminders that are due relative to now.
func (reg *SecretReminderRegistry) Due(now time.Time) []SecretReminder {
	reg.mu.RLock()
	defer reg.mu.RUnlock()
	var due []SecretReminder
	for _, r := range reg.entries {
		if r.IsDue(now) {
			due = append(due, r)
		}
	}
	return due
}
