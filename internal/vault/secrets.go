package vault

import (
	"context"
	"fmt"
	"path"
)

// SecretVersion holds the data and metadata for a specific version of a secret.
type SecretVersion struct {
	Version  int
	Data     map[string]string
	Metadata map[string]string
}

// ReadSecret reads the latest version of a KV v2 secret at the given mount and secretPath.
func (c *Client) ReadSecret(ctx context.Context, mount, secretPath string) (*SecretVersion, error) {
	fullPath := path.Join(mount, "data", secretPath)
	secret, err := c.logical.ReadWithContext(ctx, fullPath)
	if err != nil {
		return nil, fmt.Errorf("reading secret %q: %w", fullPath, err)
	}
	if secret == nil || secret.Data == nil {
		return nil, fmt.Errorf("secret %q not found", fullPath)
	}
	return parseSecretVersion(secret.Data)
}

// ReadSecretVersion reads a specific version of a KV v2 secret.
func (c *Client) ReadSecretVersion(ctx context.Context, mount, secretPath string, version int) (*SecretVersion, error) {
	fullPath := path.Join(mount, "data", secretPath)
	params := map[string][]string{
		"version": {fmt.Sprintf("%d", version)},
	}
	secret, err := c.logical.ReadWithDataWithContext(ctx, fullPath, params)
	if err != nil {
		return nil, fmt.Errorf("reading secret %q version %d: %w", fullPath, version, err)
	}
	if secret == nil || secret.Data == nil {
		return nil, fmt.Errorf("secret %q version %d not found", fullPath, version)
	}
	return parseSecretVersion(secret.Data)
}

func parseSecretVersion(raw map[string]interface{}) (*SecretVersion, error) {
	sv := &SecretVersion{
		Data:     make(map[string]string),
		Metadata: make(map[string]string),
	}

	if data, ok := raw["data"].(map[string]interface{}); ok {
		for k, v := range data {
			sv.Data[k] = fmt.Sprintf("%v", v)
		}
	}

	if meta, ok := raw["metadata"].(map[string]interface{}); ok {
		for k, v := range meta {
			sv.Metadata[k] = fmt.Sprintf("%v", v)
		}
		if ver, ok := meta["version"]; ok {
			fmt.Sscanf(fmt.Sprintf("%v", ver), "%d", &sv.Version)
		}
	}

	return sv, nil
}
