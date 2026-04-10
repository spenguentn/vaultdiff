package vault

import (
	"fmt"
	"os"
	"strings"
)

const envVaultNamespace = "VAULT_NAMESPACE"

// NamespaceSource describes where a namespace value was resolved from.
type NamespaceSource string

const (
	NamespaceSourceDirect NamespaceSource = "direct"
	NamespaceSourceEnv    NamespaceSource = "env"
	NamespaceSourceRoot   NamespaceSource = "root"
)

// ResolvedNamespace holds a resolved namespace and its source.
type ResolvedNamespace struct {
	Namespace Namespace
	Source    NamespaceSource
}

// ResolveNamespace resolves the active Vault namespace using the following priority:
//  1. Explicit value passed directly (non-empty).
//  2. VAULT_NAMESPACE environment variable.
//  3. Root namespace (empty).
func ResolveNamespace(direct string) (ResolvedNamespace, error) {
	if direct = strings.TrimSpace(direct); direct != "" {
		n := NewNamespace(direct)
		if err := n.Validate(); err != nil {
			return ResolvedNamespace{}, fmt.Errorf("invalid namespace %q: %w", direct, err)
		}
		return ResolvedNamespace{Namespace: n, Source: NamespaceSourceDirect}, nil
	}

	if env := strings.TrimSpace(os.Getenv(envVaultNamespace)); env != "" {
		n := NewNamespace(env)
		if err := n.Validate(); err != nil {
			return ResolvedNamespace{}, fmt.Errorf("invalid namespace from env %q: %w", env, err)
		}
		return ResolvedNamespace{Namespace: n, Source: NamespaceSourceEnv}, nil
	}

	return ResolvedNamespace{Namespace: NewNamespace(""), Source: NamespaceSourceRoot}, nil
}
