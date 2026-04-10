package vault

import (
	"errors"
	"fmt"
	"strings"
)

// AuthMethod represents a Vault authentication method.
type AuthMethod string

const (
	AuthToken      AuthMethod = "token"
	AuthAppRole    AuthMethod = "approle"
	AuthKubernetes AuthMethod = "kubernetes"
)

// AuthConfig holds configuration for authenticating with Vault.
type AuthConfig struct {
	Method    AuthMethod
	MountPath string
	Token     string
	RoleID    string
	SecretID  string
	Role      string // used for kubernetes
	JWT       string // used for kubernetes
}

// Validate checks that the AuthConfig is well-formed for its method.
func (a AuthConfig) Validate() error {
	switch a.Method {
	case AuthToken:
		if strings.TrimSpace(a.Token) == "" {
			return errors.New("auth: token method requires a non-empty token")
		}
	case AuthAppRole:
		if strings.TrimSpace(a.RoleID) == "" {
			return errors.New("auth: approle method requires a role_id")
		}
		if strings.TrimSpace(a.SecretID) == "" {
			return errors.New("auth: approle method requires a secret_id")
		}
	case AuthKubernetes:
		if strings.TrimSpace(a.Role) == "" {
			return errors.New("auth: kubernetes method requires a role")
		}
		if strings.TrimSpace(a.JWT) == "" {
			return errors.New("auth: kubernetes method requires a jwt")
		}
	default:
		return fmt.Errorf("auth: unsupported method %q", a.Method)
	}
	return nil
}

// MountPath returns the effective mount path for the auth method.
// If MountPath is explicitly set it is returned; otherwise the default
// path for the method is used.
func (a AuthConfig) EffectiveMountPath() string {
	if strings.TrimSpace(a.MountPath) != "" {
		return strings.Trim(a.MountPath, "/")
	}
	switch a.Method {
	case AuthAppRole:
		return "approle"
	case AuthKubernetes:
		return "kubernetes"
	default:
		return string(a.Method)
	}
}

// LoginPath returns the Vault API path used to authenticate.
func (a AuthConfig) LoginPath() string {
	return fmt.Sprintf("auth/%s/login", a.EffectiveMountPath())
}
