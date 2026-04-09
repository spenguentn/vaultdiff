package vault

import (
	"fmt"

	vaultapi "github.com/hashicorp/vault/api"
)

// Client wraps the Vault API client with additional functionality
type Client struct {
	api *vaultapi.Client
}

// Config holds the configuration for connecting to Vault
type Config struct {
	Address string
	Token   string
}

// NewClient creates a new Vault client with the provided configuration
func NewClient(cfg Config) (*Client, error) {
	if cfg.Address == "" {
		return nil, fmt.Errorf("vault address is required")
	}

	if cfg.Token == "" {
		return nil, fmt.Errorf("vault token is required")
	}

	config := vaultapi.DefaultConfig()
	config.Address = cfg.Address

	client, err := vaultapi.NewClient(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create vault client: %w", err)
	}

	client.SetToken(cfg.Token)

	return &Client{
		api: client,
	}, nil
}

// GetSecret retrieves a secret from the specified path
func (c *Client) GetSecret(path string) (map[string]interface{}, error) {
	secret, err := c.api.Logical().Read(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read secret at %s: %w", path, err)
	}

	if secret == nil {
		return nil, fmt.Errorf("secret not found at path: %s", path)
	}

	if secret.Data == nil {
		return nil, fmt.Errorf("secret data is nil at path: %s", path)
	}

	return secret.Data, nil
}

// GetSecretVersion retrieves a specific version of a secret (KV v2)
func (c *Client) GetSecretVersion(path string, version int) (map[string]interface{}, error) {
	params := map[string][]string{
		"version": {fmt.Sprintf("%d", version)},
	}

	secret, err := c.api.Logical().ReadWithData(path, params)
	if err != nil {
		return nil, fmt.Errorf("failed to read secret version %d at %s: %w", version, path, err)
	}

	if secret == nil || secret.Data == nil {
		return nil, fmt.Errorf("secret version %d not found at path: %s", version, path)
	}

	return secret.Data, nil
}

// Health checks if the Vault server is healthy and accessible
func (c *Client) Health() error {
	health, err := c.api.Sys().Health()
	if err != nil {
		return fmt.Errorf("vault health check failed: %w", err)
	}

	if health.Sealed {
		return fmt.Errorf("vault is sealed")
	}

	return nil
}
