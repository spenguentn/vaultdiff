package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Config holds the top-level vaultdiff configuration.
type Config struct {
	Vault   VaultConfig   `yaml:"vault"`
	Audit   AuditConfig   `yaml:"audit"`
	Output  OutputConfig  `yaml:"output"`
}

// VaultConfig holds Vault connection settings.
type VaultConfig struct {
	Address   string `yaml:"address"`
	Token     string `yaml:"token"`
	Namespace string `yaml:"namespace"`
	TLSSkip   bool   `yaml:"tls_skip_verify"`
}

// AuditConfig holds audit logging settings.
type AuditConfig struct {
	Enabled bool   `yaml:"enabled"`
	Format  string `yaml:"format"` // json | text
	Path    string `yaml:"path"`
}

// OutputConfig holds output formatting settings.
type OutputConfig struct {
	Format      string `yaml:"format"`       // text | json
	MaskSecrets bool   `yaml:"mask_secrets"`
}

// Load reads a YAML config file from the given path and returns a Config.
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("config: reading file %q: %w", path, err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("config: parsing YAML: %w", err)
	}

	if err := cfg.validate(); err != nil {
		return nil, err
	}

	return &cfg, nil
}

// validate applies defaults and checks required fields.
func (c *Config) validate() error {
	if c.Vault.Address == "" {
		c.Vault.Address = "http://127.0.0.1:8200"
	}
	if c.Output.Format == "" {
		c.Output.Format = "text"
	}
	allowed := map[string]bool{"text": true, "json": true}
	if !allowed[c.Output.Format] {
		return fmt.Errorf("config: unsupported output format %q (want text|json)", c.Output.Format)
	}
	if c.Audit.Enabled {
		if c.Audit.Format == "" {
			c.Audit.Format = "json"
		}
		if !allowed[c.Audit.Format] {
			return fmt.Errorf("config: unsupported audit format %q (want text|json)", c.Audit.Format)
		}
	}
	return nil
}
