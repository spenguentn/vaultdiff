// Package vault provides utilities for interacting with HashiCorp Vault.
//
// # Secret Search
//
// SearchSecrets scans an in-memory map of secret data and returns all keys
// whose names match a given query string.  It is designed to work with
// snapshot data already loaded into memory, keeping it free of network I/O.
//
// Usage:
//
//	data := map[string]map[string]string{
//	    "secret/app/config": {"db_password": "x", "api_key": "y"},
//	}
//	results := vault.SearchSecrets(data, vault.SearchOptions{
//	    Query:      "db_",
//	    Mount:      "secret",
//	    MaxResults: 10,
//	})
//	for _, r := range results {
//	    fmt.Println(r.FullPath())
//	}
//
// SearchOptions fields:
//   - Query:      case-insensitive substring matched against key names.
//   - Mount:      when non-empty, restricts results to that mount prefix.
//   - MaxResults: caps total results; 0 means unlimited.
package vault
