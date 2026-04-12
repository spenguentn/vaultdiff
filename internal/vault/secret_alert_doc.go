// Package vault provides the SecretAlert and SecretAlertRegistry types for
// recording and querying alerts raised against Vault secrets.
//
// # Overview
//
// A SecretAlert captures a discrete event — such as an imminent expiry,
// policy violation, or unexpected modification — associated with a specific
// secret identified by its mount and path.
//
// # Severity levels
//
// Four severity levels are supported: low, medium, high, and critical.
// Use IsValidSeverity to check an arbitrary string before constructing an
// alert.
//
// # Registry
//
// SecretAlertRegistry provides a thread-safe, in-memory store for alerts.
// Multiple alerts may be recorded per path. Use Clear to reset alerts for a
// given path once they have been acknowledged.
package vault
