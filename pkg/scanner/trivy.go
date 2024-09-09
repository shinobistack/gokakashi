package scanner

import (
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
)

type TrivyScanner struct{}

// NewTrivyScanner initializes a new Trivy scanner
func NewTrivyScanner() *TrivyScanner {
	return &TrivyScanner{}
}

// ScanImage scans the given Docker image for vulnerabilities using Trivy
func (t *TrivyScanner) ScanImage(image string) (string, error) {
	log.Printf("Scanning Docker image: %s with Trivy", image)

	cmd := exec.Command("trivy", "image", "--format", "json", image)
	output, err := cmd.Output()

	if err != nil {
		return "", fmt.Errorf("Trivy scan failed: %v", err)
	}

	// Validate JSON output
	if !isValidJSON(output) {
		return "", fmt.Errorf("invalid JSON output from Trivy")
	}

	return string(output), nil
}

// isValidJSON checks if the provided string is valid JSON
func isValidJSON(s []byte) bool {
	var js interface{}
	return json.Unmarshal(s, &js) == nil
}
