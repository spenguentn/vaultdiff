package vault

import "fmt"

// SecretScore represents a computed quality/risk score for a secret.
type SecretScore struct {
	Mount   string
	Path    string
	Score   int
	Reasons []string
}

// FullPath returns the canonical mount+path key.
func (s SecretScore) FullPath() string {
	return fmt.Sprintf("%s/%s", s.Mount, s.Path)
}

// Grade returns a letter grade based on the numeric score.
func (s SecretScore) Grade() string {
	switch {
	case s.Score >= 90:
		return "A"
	case s.Score >= 75:
		return "B"
	case s.Score >= 60:
		return "C"
	case s.Score >= 40:
		return "D"
	default:
		return "F"
	}
}

// ScoreRequest holds the inputs needed to compute a secret's score.
type ScoreRequest struct {
	Mount      string
	Path       string
	HasExpiry  bool
	HasOwner   bool
	HasSchema  bool
	HasChecksum bool
	Versions   int
}

// Validate ensures the request has required fields.
func (r ScoreRequest) Validate() error {
	if r.Mount == "" {
		return fmt.Errorf("secret score: mount is required")
	}
	if r.Path == "" {
		return fmt.Errorf("secret score: path is required")
	}
	return nil
}

// ComputeSecretScore calculates a quality score for a secret based on
// the presence of governance metadata.
func ComputeSecretScore(req ScoreRequest) (SecretScore, error) {
	if err := req.Validate(); err != nil {
		return SecretScore{}, err
	}

	score := 0
	var reasons []string

	if req.HasExpiry {
		score += 25
		reasons = append(reasons, "expiry policy set (+25)")
	} else {
		reasons = append(reasons, "no expiry policy (0)")
	}

	if req.HasOwner {
		score += 25
		reasons = append(reasons, "ownership record set (+25)")
	} else {
		reasons = append(reasons, "no ownership record (0)")
	}

	if req.HasSchema {
		score += 25
		reasons = append(reasons, "schema defined (+25)")
	} else {
		reasons = append(reasons, "no schema defined (0)")
	}

	if req.HasChecksum {
		score += 15
		reasons = append(reasons, "checksum registered (+15)")
	} else {
		reasons = append(reasons, "no checksum (0)")
	}

	if req.Versions > 1 {
		score += 10
		reasons = append(reasons, "multiple versions present (+10)")
	} else {
		reasons = append(reasons, "single version only (0)")
	}

	return SecretScore{
		Mount:   req.Mount,
		Path:    req.Path,
		Score:   score,
		Reasons: reasons,
	}, nil
}
