package snapshot

import (
	"fmt"
	"sync"
)

// Store is an in-memory registry of named snapshots.
type Store struct {
	mu    sync.RWMutex
	items map[string]*Snapshot
}

// NewStore creates an empty Store.
func NewStore() *Store {
	return &Store{items: make(map[string]*Snapshot)}
}

// Save stores a snapshot under the given name, overwriting any existing entry.
func (s *Store) Save(name string, snap *Snapshot) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.items[name] = snap
}

// Get retrieves a snapshot by name. Returns an error if not found.
func (s *Store) Get(name string) (*Snapshot, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	snap, ok := s.items[name]
	if !ok {
		return nil, fmt.Errorf("snapshot %q not found", name)
	}
	return snap, nil
}

// Delete removes a snapshot by name. Returns an error if it does not exist.
func (s *Store) Delete(name string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.items[name]; !ok {
		return fmt.Errorf("snapshot %q not found", name)
	}
	delete(s.items, name)
	return nil
}

// List returns all stored snapshot names.
func (s *Store) List() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	names := make([]string, 0, len(s.items))
	for k := range s.items {
		names = append(names, k)
	}
	return names
}

// Len returns the number of stored snapshots.
func (s *Store) Len() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.items)
}
