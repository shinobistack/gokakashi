package scanner

import (
	"fmt"
)

// NewScanner initializes the scanner based on the provided tool name.
func NewScanner(tool string) (Scanner, error) {
	switch tool {
	case "trivy":
		return &TrivyScanner{}, nil
	// New scanners cases here like snyk
	default:
		return nil, fmt.Errorf("unsupported scanner: %s", tool)
	}
}
