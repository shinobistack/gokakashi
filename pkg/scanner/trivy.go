package scanner

import (
	"encoding/json"
	"fmt"
	"github.com/ashwiniag/goKakashi/notifier" // Import the notifier package
	"log"
	"os/exec"
)

type TrivyScanner struct{}

// NewTrivyScanner initializes a new Trivy scanner
func NewTrivyScanner() *TrivyScanner {
	return &TrivyScanner{}
}

// ScanImage scans the given Docker image for vulnerabilities using Trivy
// It returns both the raw JSON output and the parsed vulnerabilities
func (t *TrivyScanner) ScanImage(image string) (string, []notifier.Vulnerability, error) {
	log.Printf("Scanning Docker image: %s with Trivy", image)

	// Run Trivy scan
	cmd := exec.Command("trivy", "image", "--format", "json", image)
	output, err := cmd.Output()
	if err != nil {
		return "", nil, fmt.Errorf("Trivy scan failed: %v", err)
	}

	// Validate JSON output
	if !isValidJSON(output) {
		return "", nil, fmt.Errorf("invalid JSON output from Trivy")
	}

	// Parse the JSON output into a slice of Vulnerability structs
	var vulnerabilities []notifier.Vulnerability
	err = json.Unmarshal(output, &vulnerabilities)
	if err != nil {
		return "", nil, fmt.Errorf("Failed to parse Trivy output: %v", err)
	}

	// Return both the raw JSON (for storage) and the parsed vulnerabilities (for processing)
	return string(output), vulnerabilities, nil
}

// isValidJSON checks if the provided string is valid JSON
func isValidJSON(s []byte) bool {
	var js interface{}
	return json.Unmarshal(s, &js) == nil
}
