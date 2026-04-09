package vault

import "fmt"

// EnvPair holds two named Vault environments to be compared (left vs right).
// Typically left is the source environment and right is the target.
type EnvPair struct {
	Left  *Environment
	Right *Environment
}

// NewEnvPair constructs an EnvPair and validates both environments.
func NewEnvPair(left, right *Environment) (*EnvPair, error) {
	if left == nil {
		return nil, fmt.Errorf("left environment must not be nil")
	}
	if right == nil {
		return nil, fmt.Errorf("right environment must not be nil")
	}
	if err := left.Validate(); err != nil {
		return nil, fmt.Errorf("left environment invalid: %w", err)
	}
	if err := right.Validate(); err != nil {
		return nil, fmt.Errorf("right environment invalid: %w", err)
	}
	return &EnvPair{Left: left, Right: right}, nil
}

// Names returns a human-readable label for the pair, e.g. "staging -> production".
func (p *EnvPair) Names() string {
	return fmt.Sprintf("%s -> %s", p.Left.Name, p.Right.Name)
}

// SameMount reports whether both environments share the same mount path.
func (p *EnvPair) SameMount() bool {
	return p.Left.MountPath == p.Right.MountPath
}
