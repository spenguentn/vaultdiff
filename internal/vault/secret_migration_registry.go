package vault

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
)

func migrationKey(mount, path string) string {
	return fmt.Sprintf("%s/%s", mount, path)
}

// SecretMigrationRegistry tracks in-progress and completed secret migrations.
type SecretMigrationRegistry struct {
	mu      sync.RWMutex
	entries map[string]*SecretMigration
}

// NewSecretMigrationRegistry creates an empty migration registry.
func NewSecretMigrationRegistry() *SecretMigrationRegistry {
	return &SecretMigrationRegistry{
		entries: make(map[string]*SecretMigration),
	}
}

// Submit registers a new migration. Status is set to pending and timestamps are populated.
func (r *SecretMigrationRegistry) Submit(m SecretMigration) (*SecretMigration, error) {
	if err := m.Validate(); err != nil {
		return nil, err
	}
	now := time.Now().UTC()
	m.ID = uuid.NewString()
	m.Status = MigrationPending
	m.CreatedAt = now
	m.UpdatedAt = now

	r.mu.Lock()
	defer r.mu.Unlock()
	key := migrationKey(m.SourceMount, m.SourcePath)
	if existing, ok := r.entries[key]; ok && !existing.IsTerminal() {
		return nil, fmt.Errorf("active migration already exists for %s", key)
	}
	copy := m
	r.entries[key] = &copy
	return &copy, nil
}

// UpdateStatus transitions a migration to a new status.
func (r *SecretMigrationRegistry) UpdateStatus(id string, status MigrationStatus) error {
	if !IsValidMigrationStatus(status) {
		return fmt.Errorf("invalid status: %q", status)
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, m := range r.entries {
		if m.ID == id {
			m.Status = status
			m.UpdatedAt = time.Now().UTC()
			return nil
		}
	}
	return errors.New("migration not found")
}

// Get returns the migration for the given source mount and path.
func (r *SecretMigrationRegistry) Get(mount, path string) (*SecretMigration, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	m, ok := r.entries[migrationKey(mount, path)]
	if !ok {
		return nil, false
	}
	copy := *m
	return &copy, true
}

// All returns a snapshot of all registered migrations.
func (r *SecretMigrationRegistry) All() []*SecretMigration {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]*SecretMigration, 0, len(r.entries))
	for _, m := range r.entries {
		copy := *m
		out = append(out, &copy)
	}
	return out
}
