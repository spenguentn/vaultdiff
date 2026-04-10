package vault

import (
	"errors"
	"testing"
	"time"
)

func TestDefaultRetryConfig(t *testing.T) {
	cfg := DefaultRetryConfig()
	if cfg.MaxAttempts != 3 {
		t.Errorf("expected MaxAttempts=3, got %d", cfg.MaxAttempts)
	}
	if cfg.InitialDelay != 200*time.Millisecond {
		t.Errorf("unexpected InitialDelay: %v", cfg.InitialDelay)
	}
}

func TestRetryConfig_Validate_Valid(t *testing.T) {
	cfg := DefaultRetryConfig()
	if err := cfg.Validate(); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestRetryConfig_Validate_ZeroAttempts(t *testing.T) {
	cfg := DefaultRetryConfig()
	cfg.MaxAttempts = 0
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error for MaxAttempts=0")
	}
}

func TestRetryConfig_Validate_NegativeDelay(t *testing.T) {
	cfg := DefaultRetryConfig()
	cfg.InitialDelay = -1
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error for negative InitialDelay")
	}
}

func TestRetryConfig_Validate_MaxDelayLessThanInitial(t *testing.T) {
	cfg := RetryConfig{
		MaxAttempts:  2,
		InitialDelay: time.Second,
		MaxDelay:     time.Millisecond,
	}
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error when MaxDelay < InitialDelay")
	}
}

func TestDo_SucceedsFirstAttempt(t *testing.T) {
	cfg := RetryConfig{MaxAttempts: 3, InitialDelay: 0, MaxDelay: 0}
	calls := 0
	err := Do(cfg, func() error {
		calls++
		return nil
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if calls != 1 {
		t.Errorf("expected 1 call, got %d", calls)
	}
}

func TestDo_RetriesOnError(t *testing.T) {
	cfg := RetryConfig{MaxAttempts: 3, InitialDelay: 0, MaxDelay: 0}
	calls := 0
	sentinel := errors.New("transient")
	err := Do(cfg, func() error {
		calls++
		if calls < 3 {
			return sentinel
		}
		return nil
	})
	if err != nil {
		t.Fatalf("expected success after retries, got %v", err)
	}
	if calls != 3 {
		t.Errorf("expected 3 calls, got %d", calls)
	}
}

func TestDo_ExhaustsAttempts(t *testing.T) {
	cfg := RetryConfig{MaxAttempts: 2, InitialDelay: 0, MaxDelay: 0}
	sentinel := errors.New("permanent")
	calls := 0
	err := Do(cfg, func() error {
		calls++
		return sentinel
	})
	if !errors.Is(err, sentinel) {
		t.Fatalf("expected sentinel error, got %v", err)
	}
	if calls != 2 {
		t.Errorf("expected 2 calls, got %d", calls)
	}
}
