# vaultdiff

A CLI tool to diff and audit changes between HashiCorp Vault secret versions across environments.

## Installation

```bash
go install github.com/yourusername/vaultdiff@latest
```

Or download pre-built binaries from the [releases page](https://github.com/yourusername/vaultdiff/releases).

## Usage

```bash
# Compare secrets between two Vault environments
vaultdiff --source https://vault-dev.example.com --target https://vault-prod.example.com --path secret/myapp

# Compare specific versions of a secret
vaultdiff --vault https://vault.example.com --path secret/myapp --version1 5 --version2 8

# Output as JSON for automation
vaultdiff --source https://vault-dev.example.com --target https://vault-prod.example.com --path secret/myapp --format json

# Audit mode: show all changes with metadata
vaultdiff --vault https://vault.example.com --path secret/myapp --audit
```

## Authentication

vaultdiff supports multiple authentication methods:
- `VAULT_TOKEN` environment variable
- `VAULT_ADDR` for Vault address
- Token file at `~/.vault-token`
- AppRole, Kubernetes, and other Vault auth methods

## Features

- 🔍 Compare secrets across environments or versions
- 📊 Clear diff output showing added, removed, and modified keys
- 🔐 Secure handling of sensitive data
- 📝 Audit trail with timestamps and version metadata
- 🚀 Fast parallel secret fetching
- 📦 Multiple output formats (text, JSON, YAML)

## License

MIT License - see [LICENSE](LICENSE) for details.