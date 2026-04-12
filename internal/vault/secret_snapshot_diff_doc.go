// Package vault provides snapshot diffing utilities for comparing secret
// data between two environments or versions.
//
// # SnapshotDiff
//
// DiffSnapshots accepts two flat maps of secret key-value pairs (left and
// right) and produces a SnapshotDiffResult that classifies every key as
// one of: added, removed, modified, or unchanged.
//
// Example usage:
//
//	left := map[string]string{"DB_PASS": "old", "API_KEY": "abc"}
//	right := map[string]string{"DB_PASS": "new", "API_KEY": "abc", "NEW": "x"}
//
//	result := vault.DiffSnapshots("staging", "production", left, right)
//	fmt.Println(result.Summary())
//	for _, e := range result.ChangedOnly() {
//		fmt.Printf("%s -> %s\n", e.FullPath(), e.ChangeType)
//	}
//
// SnapshotDiffResult.ChangedOnly filters out unchanged entries so callers
// can focus on actionable differences. Summary provides a quick one-line
// overview suitable for audit logs or CLI output.
package vault
