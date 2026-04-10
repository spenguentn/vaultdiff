package vault

import (
	"testing"
)

func TestNewBatchReader_DefaultConcurrency(t *testing.T) {
	client := &Client{}
	br := NewBatchReader(client, 0)
	if br.concurrency != 5 {
		t.Errorf("expected default concurrency 5, got %d", br.concurrency)
	}
}

func TestNewBatchReader_CustomConcurrency(t *testing.T) {
	client := &Client{}
	br := NewBatchReader(client, 10)
	if br.concurrency != 10 {
		t.Errorf("expected concurrency 10, got %d", br.concurrency)
	}
}

func TestNewBatchReader_NegativeConcurrency(t *testing.T) {
	client := &Client{}
	br := NewBatchReader(client, -3)
	if br.concurrency != 5 {
		t.Errorf("expected default concurrency 5 for negative input, got %d", br.concurrency)
	}
}

func TestBatchResult_Fields(t *testing.T) {
	r := BatchResult{
		Path:    "secret/app",
		Secrets: map[string]string{"key": "value"},
		Err:     nil,
	}
	if r.Path != "secret/app" {
		t.Errorf("unexpected path: %s", r.Path)
	}
	if r.Secrets["key"] != "value" {
		t.Errorf("unexpected secret value")
	}
}

func TestErrors_FiltersFailures(t *testing.T) {
	results := []BatchResult{
		{Path: "a", Err: nil},
		{Path: "b", Err: fmt.Errorf("read failed")},
		{Path: "c", Err: nil},
		{Path: "d", Err: fmt.Errorf("timeout")},
	}
	failed := Errors(results)
	if len(failed) != 2 {
		t.Errorf("expected 2 failures, got %d", len(failed))
	}
}

func TestErrors_NoneReturnsNil(t *testing.T) {
	results := []BatchResult{
		{Path: "a", Err: nil},
		{Path: "b", Err: nil},
	}
	failed := Errors(results)
	if len(failed) != 0 {
		t.Errorf("expected 0 failures, got %d", len(failed))
	}
}
