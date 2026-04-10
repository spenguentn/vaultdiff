package vault

import (
	"context"
	"fmt"
	"sync"
)

// BatchResult holds the outcome of reading a single secret path.
type BatchResult struct {
	Path    string
	Secrets map[string]string
	Err     error
}

// BatchReader reads multiple secret paths concurrently.
type BatchReader struct {
	client    *Client
	concurrency int
}

// NewBatchReader creates a BatchReader with the given client and max concurrency.
// If concurrency is <= 0, it defaults to 5.
func NewBatchReader(client *Client, concurrency int) *BatchReader {
	if concurrency <= 0 {
		concurrency = 5
	}
	return &BatchReader{
		client:      client,
		concurrency: concurrency,
	}
}

// ReadAll reads all provided paths concurrently and returns a slice of BatchResult.
// Order of results matches the order of the input paths.
func (b *BatchReader) ReadAll(ctx context.Context, mount string, paths []string) []BatchResult {
	results := make([]BatchResult, len(paths))
	sem := make(chan struct{}, b.concurrency)
	var wg sync.WaitGroup

	for i, p := range paths {
		wg.Add(1)
		go func(idx int, path string) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()

			sp, err := NewSecretPath(mount, path)
			if err != nil {
				results[idx] = BatchResult{Path: path, Err: fmt.Errorf("invalid path: %w", err)}
				return
			}

			secrets, err := ReadSecret(ctx, b.client, sp)
			results[idx] = BatchResult{Path: path, Secrets: secrets, Err: err}
		}(i, p)
	}

	wg.Wait()
	return results
}

// Errors returns only the failed BatchResults.
func Errors(results []BatchResult) []BatchResult {
	var failed []BatchResult
	for _, r := range results {
		if r.Err != nil {
			failed = append(failed, r)
		}
	}
	return failed
}
