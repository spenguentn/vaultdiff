package vault

import (
	"testing"
)

func TestComputeSecretScore_AllPresent(t *testing.T) {
	req := ScoreRequest{
		Mount:       "secret",
		Path:        "app/db",
		HasExpiry:   true,
		HasOwner:    true,
		HasSchema:   true,
		HasChecksum: true,
		Versions:    3,
	}
	result, err := ComputeSecretScore(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Score != 100 {
		t.Errorf("expected score 100, got %d", result.Score)
	}
	if result.Grade() != "A" {
		t.Errorf("expected grade A, got %s", result.Grade())
	}
}

func TestComputeSecretScore_NonePresent(t *testing.T) {
	req := ScoreRequest{
		Mount:    "secret",
		Path:     "app/db",
		Versions: 1,
	}
	result, err := ComputeSecretScore(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Score != 0 {
		t.Errorf("expected score 0, got %d", result.Score)
	}
	if result.Grade() != "F" {
		t.Errorf("expected grade F, got %s", result.Grade())
	}
}

func TestComputeSecretScore_MissingMount(t *testing.T) {
	req := ScoreRequest{Path: "app/db"}
	_, err := ComputeSecretScore(req)
	if err == nil {
		t.Fatal("expected error for missing mount")
	}
}

func TestComputeSecretScore_MissingPath(t *testing.T) {
	req := ScoreRequest{Mount: "secret"}
	_, err := ComputeSecretScore(req)
	if err == nil {
		t.Fatal("expected error for missing path")
	}
}

func TestSecretScore_FullPath(t *testing.T) {
	s := SecretScore{Mount: "secret", Path: "app/db"}
	if s.FullPath() != "secret/app/db" {
		t.Errorf("unexpected full path: %s", s.FullPath())
	}
}

func TestSecretScore_Grade_Boundaries(t *testing.T) {
	cases := []struct {
		score    int
		expected string
	}{
		{90, "A"},
		{75, "B"},
		{60, "C"},
		{40, "D"},
		{39, "F"},
	}
	for _, tc := range cases {
		s := SecretScore{Score: tc.score}
		if g := s.Grade(); g != tc.expected {
			t.Errorf("score %d: expected grade %s, got %s", tc.score, tc.expected, g)
		}
	}
}

func TestComputeSecretScore_ReasonsCount(t *testing.T) {
	req := ScoreRequest{
		Mount: "kv",
		Path:  "service/token",
	}
	result, err := ComputeSecretScore(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Reasons) != 5 {
		t.Errorf("expected 5 reasons, got %d", len(result.Reasons))
	}
}
