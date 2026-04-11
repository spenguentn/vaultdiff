// Package vault provides a SecretLock and SecretLockRegistry for coordinating
// exclusive access to secret paths during sensitive operations such as
// promotions, rollbacks, or bulk writes.
//
// # Overview
//
// A SecretLock associates an owner identity with a mount/path pair and an
// optional expiry time. The SecretLockRegistry is an in-memory, goroutine-safe
// store of active locks.
//
// # Usage
//
//	reg := vault.NewSecretLockRegistry()
//
//	lock := vault.SecretLock{
//		Mount:     "secret",
//		Path:      "app/db",
//		LockedBy:  "deploy-bot",
//		ExpiresAt: time.Now().Add(5 * time.Minute),
//	}
//
//	if err := reg.Lock(lock); err != nil {
//		log.Fatal(err)
//	}
//	defer reg.Unlock(lock.Mount, lock.Path, lock.LockedBy)
//
// Locks with a past ExpiresAt are considered expired and may be overwritten by
// a new owner without an explicit unlock.
package vault
