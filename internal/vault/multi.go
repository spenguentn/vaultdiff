package vault

import (
	"context"
	"fmt"
)

// MultiReader reads secrets from multiple Vault environments concurrently.
type MultiReader struct {
	clients map[string]*Client
}

// NewMultiReader constructs a MultiReader from a map of named Vault clients.
func NewMultiReader(clients map[string]*Client) *MultiReader {
	return &MultiReader{clients: clients}
}

// SecretResult holds the result of reading a secret from a named environment.
type SecretResult struct {
	Env     string
	Secrets map[string]string
	Err     error
}

// ReadAll reads the secret at path from every registered environment concurrently
// and returns a slice of SecretResult, one per environment.
func (m *MultiReader) ReadAll(ctx context.Context, path string) []SecretResult {
	results := make([]SecretResult, 0, len(m.clients))
	type item struct {
		env    string
		data   map[string]string
		err    error
	}
	ch := make(chan item, len(m.clients))

	for env, c := range m.clients {
		env, c := env, c
		go func() {
			data, err := c.ReadSecret(ctx, path)
			if err != nil {
				err = fmt.Errorf("env %q: %w", env, err)
			}
			ch <- item{env: env, data: data, err: err}
		}()
	}

	for i := 0; i < len(m.clients); i++ {
		it := <-ch
		results = append(results, SecretResult{
			Env:     it.env,
			Secrets: it.data,
			Err:     it.err,
		})
	}
	return results
}
