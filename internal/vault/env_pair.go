package vault

import "fmt"

// EnvPair holds two Vault environments to be compared.
type EnvPair struct {
	Left  Environment
	Right Environment
}

// NewEnvPair constructs an EnvPair from two environments and validates both.
func NewEnvPair(left, right Environment) (EnvPair, error) {
	if err := left.Validate(); err != nil {
		return EnvPair{}, fmt.Errorf("left environment invalid: %w", err)
	}
	if err := right.Validate(); err != nil {
		return EnvPair{}, fmt.Errorf("right environment invalid: %w", err)
	}
	return EnvPair{Left: left, Right: right}, nil
}

// Names returns a human-readable label for the pair.
func (p EnvPair) Names() string {
	return fmt.Sprintf("%s → %s", p.Left.Name, p.Right.Name)
}

// SameMount reports whether both environments share the same mount path.
func (p EnvPair) SameMount() bool {
	return p.Left.MountPath == p.Right.MountPath
}
