// Package snapshot provides types for capturing point-in-time representations
// of Vault secret paths, including a thread-safe in-memory store for managing
// multiple named snapshots during a diff or audit session.
//
// Typical usage:
//
//	snap := snapshot.New("secret/myapp", 2, secrets, snapshot.Meta{Environment: "staging"})
//	store := snapshot.NewStore()
//	store.Save("staging-v2", snap)
package snapshot
