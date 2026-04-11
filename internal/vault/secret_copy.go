package vault

import (
	"errors"
	"fmt"
)

// CopyRequest describes a request to copy a secret from one path to another.
type CopyRequest struct {
	SourceMount string
	SourcePath  string
	DestMount   string
	DestPath    string
	Overwrite   bool
}

// Validate returns an error if the CopyRequest is incomplete.
func (r CopyRequest) Validate() error {
	if r.SourceMount == "" {
		return errors.New("source mount is required")
	}
	if r.SourcePath == "" {
		return errors.New("source path is required")
	}
	if r.DestMount == "" {
		return errors.New("destination mount is required")
	}
	if r.DestPath == "" {
		return errors.New("destination path is required")
	}
	return nil
}

// CopyResult holds the outcome of a copy operation.
type CopyResult struct {
	Request CopyRequest
	Err     error
}

// IsSuccess returns true when no error occurred.
func (r CopyResult) IsSuccess() bool { return r.Err == nil }

// String returns a human-readable summary of the result.
func (r CopyResult) String() string {
	if r.IsSuccess() {
		return fmt.Sprintf("copied %s/%s -> %s/%s",
			r.Request.SourceMount, r.Request.SourcePath,
			r.Request.DestMount, r.Request.DestPath)
	}
	return fmt.Sprintf("copy failed %s/%s -> %s/%s: %v",
		r.Request.SourceMount, r.Request.SourcePath,
		r.Request.DestMount, r.Request.DestPath, r.Err)
}
