// Package config provides loading and validation of vaultdiff configuration
// files. Configuration is expressed as YAML and covers Vault connection
// settings, audit logging options, and output formatting preferences.
//
// Example usage:
//
//	cfg, err := config.Load("/etc/vaultdiff/config.yaml")
//	if err != nil {
//		log.Fatal(err)
//	}
package config
