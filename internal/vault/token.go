package vault

import (
	"errors"
	"os"
	"strings"
)

// TokenSource defines how a Vault token is resolved.
type TokenSource int

const (
	// TokenSourceEnv resolves the token from an environment variable.
	TokenSourceEnv TokenSource = iota
	// TokenSourceFile resolves the token from a file path.
	TokenSourceFile
	// TokenSourceDirect holds the token value directly.
	TokenSourceDirect
)

// Token holds a resolved Vault authentication token.
type Token struct {
	Value  string
	Source TokenSource
}

// ResolveToken returns a Token from the first successful source.
// Priority: direct value > environment variable > file path.
func ResolveToken(direct, envVar, filePath string) (*Token, error) {
	if direct = strings.TrimSpace(direct); direct != "" {
		return &Token{Value: direct, Source: TokenSourceDirect}, nil
	}

	if envVar != "" {
		if val := os.Getenv(envVar); val != "" {
			return &Token{Value: val, Source: TokenSourceEnv}, nil
		}
	}

	if filePath != "" {
		data, err := os.ReadFile(filePath)
		if err != nil {
			return nil, fmt.Errorf("reading token file %q: %w", filePath, err)
		}
		val := strings.TrimSpace(string(data))
		if val == "" {
			return nil, errors.New("token file is empty")
		}
		return &Token{Value: val, Source: TokenSourceFile}, nil
	}

	return nil, errors.New("no vault token found: provide a direct value, environment variable, or file path")
}
