package vault

import "strings"

// SearchOptions controls how secret keys are matched during a search.
type SearchOptions struct {
	// Query is a substring matched against each secret key (case-insensitive).
	Query string
	// Mount restricts results to a specific KV mount path.
	Mount string
	// MaxResults caps the number of results returned. 0 means no limit.
	MaxResults int
}

// SearchResult holds a single matched secret key with its location.
type SearchResult struct {
	Mount string
	Path  string
	Key   string
}

// FullPath returns the canonical mount+path+key string.
func (r SearchResult) FullPath() string {
	return r.Mount + "/" + r.Path + "#" + r.Key
}

// SearchSecrets searches a map of secret data for keys matching opts.
// data is keyed by "mount/path" and values are the secret key-value pairs.
func SearchSecrets(data map[string]map[string]string, opts SearchOptions) []SearchResult {
	var results []SearchResult
	q := strings.ToLower(opts.Query)

	for location, secrets := range data {
		mount, path, _ := splitLocation(location)
		if opts.Mount != "" && mount != opts.Mount {
			continue
		}
		for key := range secrets {
			if q == "" || strings.Contains(strings.ToLower(key), q) {
				results = append(results, SearchResult{
					Mount: mount,
					Path:  path,
					Key:   key,
				})
				if opts.MaxResults > 0 && len(results) >= opts.MaxResults {
					return results
				}
			}
		}
	}
	return results
}

// splitLocation splits "mount/path" into (mount, path, ok).
func splitLocation(location string) (mount, path string, ok bool) {
	idx := strings.Index(location, "/")
	if idx < 0 {
		return location, "", false
	}
	return location[:idx], location[idx+1:], true
}
