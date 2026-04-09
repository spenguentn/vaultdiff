package compare

import (
	"testing"

	"github.com/yourusername/vaultdiff/internal/diff"
)

// stubResults is a helper to build a minimal diff.Result slice for assertions.
func stubResults(changes ...diff.ChangeType) []diff.Result {
	results := make([]diff.Result, len(changes))
	for i, c := range changes {
		results[i] = diff.Result{Key: fmt.Sprintf("KEY_%d", i), Change: c}
	}
	return results
}

func TestSource_Defaults(t *testing.T) {
	s := Source{
		Environment: "staging",
		Mount:       "secret",
		SecretPath:  "app/config",
	}
	if s.Version != 0 {
		t.Errorf("expected default version 0 (latest), got %d", s.Version)
	}
}

func TestNewEngine_NotNil(t *testing.T) {
	// NewEngine should never return nil even with nil clients
	// (nil clients will fail only at Run time).
	e := NewEngine(nil, nil)
	if e == nil {
		t.Fatal("expected non-nil engine")
	}
}

func TestEngine_LeftRightAssignment(t *testing.T) {
	e := NewEngine(nil, nil)
	if e.Left != nil || e.Right != nil {
		t.Errorf("expected both clients nil in this test setup")
	}
}
