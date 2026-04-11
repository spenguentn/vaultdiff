package vault

import (
	"context"
	"errors"
	"testing"
)

// --- fakes ---

type fakeReader struct {
	data map[string]interface{}
	err  error
}

func (f *fakeReader) ReadSecret(_ context.Context, _, _ string) (map[string]interface{}, error) {
	return f.data, f.err
}

type fakeWriter struct{ err error }

func (f *fakeWriter) WriteSecret(_ context.Context, _, _ string, _ map[string]interface{}) error {
	return f.err
}

// --- CopyRequest ---

func TestCopyRequest_Validate_Valid(t *testing.T) {
	req := CopyRequest{SourceMount: "kv", SourcePath: "a", DestMount: "kv", DestPath: "b"}
	if err := req.Validate(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCopyRequest_Validate_MissingSourceMount(t *testing.T) {
	req := CopyRequest{SourcePath: "a", DestMount: "kv", DestPath: "b"}
	if err := req.Validate(); err == nil {
		t.Fatal("expected error for missing source mount")
	}
}

func TestCopyRequest_Validate_MissingDestPath(t *testing.T) {
	req := CopyRequest{SourceMount: "kv", SourcePath: "a", DestMount: "kv"}
	if err := req.Validate(); err == nil {
		t.Fatal("expected error for missing dest path")
	}
}

// --- CopyResult ---

func TestCopyResult_IsSuccess_True(t *testing.T) {
	r := CopyResult{Request: CopyRequest{SourceMount: "kv", SourcePath: "a", DestMount: "kv", DestPath: "b"}}
	if !r.IsSuccess() {
		t.Fatal("expected success")
	}
}

func TestCopyResult_IsSuccess_WithErr(t *testing.T) {
	r := CopyResult{Err: errors.New("boom")}
	if r.IsSuccess() {
		t.Fatal("expected failure")
	}
}

func TestCopyResult_String_OK(t *testing.T) {
	r := CopyResult{Request: CopyRequest{SourceMount: "kv", SourcePath: "a", DestMount: "kv", DestPath: "b"}}
	s := r.String()
	if s == "" {
		t.Fatal("expected non-empty string")
	}
}

// --- SecretCopier ---

func TestNewSecretCopier_NilReaderPanics(t *testing.T) {
	defer func() { recover() }()
	NewSecretCopier(nil, &fakeWriter{})
	t.Fatal("expected panic")
}

func TestNewSecretCopier_NilWriterPanics(t *testing.T) {
	defer func() { recover() }()
	NewSecretCopier(&fakeReader{}, nil)
	t.Fatal("expected panic")
}

func TestSecretCopier_Copy_Success(t *testing.T) {
	r := &fakeReader{data: map[string]interface{}{"k": "v"}}
	w := &fakeWriter{}
	c := NewSecretCopier(r, w)
	req := CopyRequest{SourceMount: "kv", SourcePath: "a", DestMount: "kv", DestPath: "b"}
	res := c.Copy(context.Background(), req)
	if !res.IsSuccess() {
		t.Fatalf("unexpected error: %v", res.Err)
	}
}

func TestSecretCopier_Copy_ReadError(t *testing.T) {
	r := &fakeReader{err: errors.New("not found")}
	w := &fakeWriter{}
	c := NewSecretCopier(r, w)
	req := CopyRequest{SourceMount: "kv", SourcePath: "a", DestMount: "kv", DestPath: "b"}
	res := c.Copy(context.Background(), req)
	if res.IsSuccess() {
		t.Fatal("expected failure on read error")
	}
}

func TestSecretCopier_CopyPlanAll_ReturnsAllResults(t *testing.T) {
	r := &fakeReader{data: map[string]interface{}{"x": "1"}}
	w := &fakeWriter{}
	c := NewSecretCopier(r, w)
	plan := NewCopyPlan()
	_ = plan.Add(CopyRequest{SourceMount: "kv", SourcePath: "a", DestMount: "kv", DestPath: "b"})
	_ = plan.Add(CopyRequest{SourceMount: "kv", SourcePath: "c", DestMount: "kv", DestPath: "d"})
	results := c.CopyPlanAll(context.Background(), plan)
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
}
