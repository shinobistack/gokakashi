package scanner

import (
	"encoding/json"
	"fmt"
	"github.com/ashwiniag/goKakashi/notifier" // Import the notifier package
	"log"
	"os/exec"
	"strings"
)

type TrivyScanner struct{}

// NewTrivyScanner initializes a new Trivy scanner
func NewTrivyScanner() *TrivyScanner {
	return &TrivyScanner{}
}

// ScanImage scans the given Docker image for vulnerabilities using Trivy
// It returns both the raw JSON output and the parsed vulnerabilities
func (t *TrivyScanner) ScanImage(image string, severityLevels []string) (string, []notifier.Vulnerability, error) {
	// Check if severity levels are provided
	var cmd *exec.Cmd
	if len(severityLevels) > 0 {
		// Use Trivy's --severity flag if severity levels are specified
		severity := strings.Join(severityLevels, ",")
		log.Printf("Scanning Docker image: %s with Trivy (severity: %s)", image, severity)
		cmd = exec.Command("trivy", "image", "--format", "json", "--severity", severity, image)
	} else {
		// Perform a normal scan if no severity levels are provided
		log.Printf("Scanning Docker image: %s with Trivy (no severity filter)", image)
		cmd = exec.Command("trivy", "image", "--format", "json", image)
	}

	// Run the Trivy scan
	output, err := cmd.Output()
	if err != nil {
		return "", nil, fmt.Errorf("Trivy scan failed: %v", err)
	}

	// Validate JSON output
	if !isValidJSON(output) {
		return "", nil, fmt.Errorf("invalid JSON output from Trivy")
	}

	// Parse the JSON output into a TrivyReport struct
	var report notifier.TrivyReport
	err = json.Unmarshal(output, &report)
	if err != nil {
		return "", nil, fmt.Errorf("Failed to parse Trivy output: %v", err)
	}

	// Extract vulnerabilities from the report
	var vulnerabilities []notifier.Vulnerability
	for _, result := range report.Results {
		vulnerabilities = append(vulnerabilities, result.Vulnerabilities...)
	}

	// Return both the raw JSON (for storage) and the parsed vulnerabilities (for processing)
	return string(output), vulnerabilities, nil
}

// isValidJSON checks if the provided string is valid JSON
func isValidJSON(s []byte) bool {
	var js interface{}
	return json.Unmarshal(s, &js) == nil
}
