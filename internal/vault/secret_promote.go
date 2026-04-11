package vault

import "fmt"

// PromoteRequest describes a secret promotion from one environment to another.
type PromoteRequest struct {
	SourceMount string
	SourcePath  string
	DestMount   string
	DestPath    string
	Version     int // 0 means latest
	DryRun      bool
}

// Validate checks that the PromoteRequest has all required fields.
func (r PromoteRequest) Validate() error {
	if r.SourceMount == "" {
		return fmt.Errorf("source mount is required")
	}
	if r.SourcePath == "" {
		return fmt.Errorf("source path is required")
	}
	if r.DestMount == "" {
		return fmt.Errorf("destination mount is required")
	}
	if r.DestPath == "" {
		return fmt.Errorf("destination path is required")
	}
	if r.Version < 0 {
		return fmt.Errorf("version must be >= 0")
	}
	return nil
}

// PromoteResult holds the outcome of a single promotion operation.
type PromoteResult struct {
	Request PromoteRequest
	Err     error
}

// IsSuccess returns true when the promotion completed without error.
func (r PromoteResult) IsSuccess() bool {
	return r.Err == nil
}

// String returns a human-readable summary of the result.
func (r PromoteResult) String() string {
	if r.IsSuccess() {
		if r.Request.DryRun {
			return fmt.Sprintf("[dry-run] would promote %s/%s -> %s/%s",
				r.Request.SourceMount, r.Request.SourcePath,
				r.Request.DestMount, r.Request.DestPath)
		}
		return fmt.Sprintf("promoted %s/%s -> %s/%s",
			r.Request.SourceMount, r.Request.SourcePath,
			r.Request.DestMount, r.Request.DestPath)
	}
	return fmt.Sprintf("failed to promote %s/%s: %v",
		r.Request.SourceMount, r.Request.SourcePath, r.Err)
}
