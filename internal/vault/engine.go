package vault

import "fmt"

// EngineType represents the type of a Vault secrets engine.
type EngineType string

const (
	EngineKV1    EngineType = "kv-v1"
	EngineKV2    EngineType = "kv-v2"
	EngineGeneric EngineType = "generic"
	EngineUnknown EngineType = "unknown"
)

// SecretsEngine describes a mounted secrets engine in Vault.
type SecretsEngine struct {
	Mount       string
	Type        EngineType
	Description string
	Local       bool
	Sealed      bool
}

// Validate returns an error if the SecretsEngine is misconfigured.
func (e SecretsEngine) Validate() error {
	if e.Mount == "" {
		return fmt.Errorf("secrets engine mount path must not be empty")
	}
	if e.Type == "" {
		return fmt.Errorf("secrets engine type must not be empty")
	}
	return nil
}

// IsVersioned reports whether the engine supports versioned secrets (KV v2).
func (e SecretsEngine) IsVersioned() bool {
	return e.Type == EngineKV2
}

// String returns a human-readable representation of the engine.
func (e SecretsEngine) String() string {
	return fmt.Sprintf("%s (%s)", e.Mount, e.Type)
}

// ParseEngineType converts a raw string to an EngineType.
func ParseEngineType(raw string) EngineType {
	switch raw {
	case "kv-v1", "kv":
		return EngineKV1
	case "kv-v2":
		return EngineKV2
	case "generic":
		return EngineGeneric
	default:
		return EngineUnknown
	}
}
