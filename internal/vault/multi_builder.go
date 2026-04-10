package vault

import "fmt"

// MultiClientBuilder constructs a set of named Vault clients from a slice of
// Environment definitions, applying per-environment token resolution.
type MultiClientBuilder struct {
	envs   []Environment
	tokens map[string]string // env name -> token override
}

// NewMultiClientBuilder returns a builder seeded with the given environments.
func NewMultiClientBuilder(envs []Environment) *MultiClientBuilder {
	return &MultiClientBuilder{
		envs:   envs,
		tokens: make(map[string]string),
	}
}

// WithToken sets an explicit Vault token for the named environment.
func (b *MultiClientBuilder) WithToken(envName, token string) *MultiClientBuilder {
	b.tokens[envName] = token
	return b
}

// Build validates each environment and constructs a Client per entry.
// Returns an error if any environment fails validation or client creation.
func (b *MultiClientBuilder) Build() (map[string]*Client, error) {
	clients := make(map[string]*Client, len(b.envs))
	for _, env := range b.envs {
		if err := env.Validate(); err != nil {
			return nil, fmt.Errorf("environment %q invalid: %w", env.Name, err)
		}
		token := b.tokens[env.Name]
		cfg := Config{
			Address:   env.Address,
			Token:     token,
			MountPath: env.MountPath,
		}
		c, err := NewClient(cfg)
		if err != nil {
			return nil, fmt.Errorf("client for %q: %w", env.Name, err)
		}
		clients[env.Name] = c
	}
	return clients, nil
}
