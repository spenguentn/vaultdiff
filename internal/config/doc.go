// Package config provides loading and validation of vaultdiff configuration
// files. Configuration is expressed as YAML and covers Vault connection
// settings, audit logging options, and output formatting preferences.
//
// Configuration file locations are searched in the following order:
//
//  1. Path explicitly passed to Load
//  2. $VAULTDIFF_CONFIG environment variable
//  3. $HOME/.vaultdiff/config.yaml
//  4. /etc/vaultdiff/config.yaml
//
// Example usage:
//
//	cfg, err := config.Load("/etc/vaultdiff/config.yaml")
//	if err != nil {
//		log.Fatal(err)
//	}
//
// Environment variables prefixed with VAULTDIFF_ can override individual
// configuration fields. For example, VAULTDIFF_VAULT_ADDR overrides the
// vault_addr field in the configuration file.
package config
