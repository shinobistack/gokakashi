package scanner

import (
	"encoding/json"
	"fmt"
	"os/exec"
)

// TrivyScanner struct
type TrivyScanner struct{}

// NewTrivyScanner creates a new TrivyScanner instance
func NewTrivyScanner() *TrivyScanner {
	return &TrivyScanner{}
}

// ScanImage scans the Docker image using Trivy and returns the report in JSON format
func (t *TrivyScanner) ScanImage(image string) (string, error) {
	// Execute Trivy scan
	cmd := exec.Command("trivy", "image", "--quiet", "--format", "json", image)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("trivy scan failed: %v, %s", err, string(output))
	}

	// Validate JSON output
	if !isValidJSON(string(output)) {
		return "", fmt.Errorf("invalid JSON output from Trivy")
	}

	return string(output), nil
}

// isValidJSON checks if the provided string is valid JSON
func isValidJSON(s string) bool {
	var js interface{}
	return json.Unmarshal([]byte(s), &js) == nil
}
