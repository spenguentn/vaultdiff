package vault

import "strings"

// SecretClassification represents the sensitivity level of a secret.
type SecretClassification string

const (
	ClassificationPublic       SecretClassification = "public"
	ClassificationInternal     SecretClassification = "internal"
	ClassificationConfidential SecretClassification = "confidential"
	ClassificationRestricted   SecretClassification = "restricted"
)

// IsValidClassification returns true if the given classification is known.
func IsValidClassification(c SecretClassification) bool {
	switch c {
	case ClassificationPublic, ClassificationInternal, ClassificationConfidential, ClassificationRestricted:
		return true
	}
	return false
}

// ClassifyRequest holds the inputs needed to classify a secret.
type ClassifyRequest struct {
	Mount          string
	Path           string
	Classification SecretClassification
	ClassifiedBy   string
}

// Validate returns an error if the request is incomplete.
func (r ClassifyRequest) Validate() error {
	if r.Mount == "" {
		return ErrMissingMount
	}
	if r.Path == "" {
		return ErrMissingPath
	}
	if !IsValidClassification(r.Classification) {
		return ErrInvalidClassification
	}
	if r.ClassifiedBy == "" {
		return ErrMissingClassifiedBy
	}
	return nil
}

// FullPath returns the canonical mount+path identifier.
func (r ClassifyRequest) FullPath() string {
	return strings.TrimRight(r.Mount, "/") + "/" + strings.TrimLeft(r.Path, "/")
}

// ClassificationRecord stores the resolved classification for a secret.
type ClassificationRecord struct {
	ClassifyRequest
	Label string
}

// NewClassificationRecord builds a record from a valid request.
func NewClassificationRecord(req ClassifyRequest) (ClassificationRecord, error) {
	if err := req.Validate(); err != nil {
		return ClassificationRecord{}, err
	}
	return ClassificationRecord{
		ClassifyRequest: req,
		Label:           string(req.Classification),
	}, nil
}
