// Package vault provides a thin client layer over the HashiCorp Vault API,
// covering KV v2 operations used by vaultdiff: reading secrets at specific
// versions, listing available versions for a path, and fetching secret
// metadata such as creation timestamps and destruction status.
//
// All exported functions accept a context.Context so callers can enforce
// deadlines and cancellation across concurrent comparisons.
package vault
