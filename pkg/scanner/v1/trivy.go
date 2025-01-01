package scanner

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

type TrivyScanner struct{}

func (t *TrivyScanner) Scan(image string, severityLevels []string) (string, error) {
	// Create a temporary file for the report
	outputFile, err := os.CreateTemp("", "trivy-report-*.json")
	if err != nil {
		return "", fmt.Errorf("failed to create temp file for Trivy report: %w", err)
	}
	defer outputFile.Close()

	// Build the command
	var cmd *exec.Cmd
	if len(severityLevels) > 0 {
		severity := strings.Join(severityLevels, ",")
		log.Printf("[INFO] Scanning Docker image: %s with Trivy (severity: %s)", image, severity)
		cmd = exec.Command("trivy", "image", "--format", "json", "--output", outputFile.Name(), "--severity", severity, image)
	} else {
		log.Printf("[INFO] Scanning Docker image: %s with Trivy (no severity filter)", image)
		cmd = exec.Command("trivy", "image", "--format", "json", "--output", outputFile.Name(), image)
	}

	// Capture standard output and error for debugging
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("[ERROR] Trivy output: %s", string(output))
		return "", fmt.Errorf("Trivy scan failed: %w", err)
	}

	log.Printf("[INFO] Trivy scan completed for image: %s. Report saved to: %s", image, outputFile.Name())

	// Validate JSON output
	if !isValidJSON(output) {
		return "", fmt.Errorf("invalid JSON output from Trivy")
	}

	return outputFile.Name(), nil
}

// isValidJSON checks if the provided string is valid JSON
func isValidJSON(s []byte) bool {
	var js interface{}
	return json.Unmarshal(s, &js) == nil
}