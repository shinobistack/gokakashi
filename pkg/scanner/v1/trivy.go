package scanner

import (
	"encoding/json"
	"fmt"
	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/checker/decls"
	"log"
	"os"
	"os/exec"
	"strings"
)

type TrivyScanner struct{}

func (t *TrivyScanner) Scan(image string, severityLevels []string) (string, error) {
	// Create a temporary file for the report
	// Todo: to make use of workspace for agents
	outputFile, err := os.CreateTemp("", "trivy-report-*.json")
	if err != nil {
		return "", fmt.Errorf("failed to create temp file for Trivy report: %w", err)
	}
	defer outputFile.Close()

	// Build the Trivy command
	var cmd *exec.Cmd
	if len(severityLevels) > 0 {
		severity := strings.Join(severityLevels, ",")
		log.Printf("[INFO] Scanning Docker image: %s with Trivy (severity: %s)", image, severity)
		cmd = exec.Command("trivy", "image", "--format", "json", "--output", outputFile.Name(), "--severity", severity, image)
	} else {
		log.Printf("[INFO] Scanning Docker image: %s with Trivy (no severity filter)", image)
		cmd = exec.Command("trivy", "image", "--format", "json", "--output", outputFile.Name(), image)
	}

	// Execute the command and capture output for debugging
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("[ERROR] Trivy scan failed. Output: %s", string(output))
		return "", fmt.Errorf("Trivy scan failed: %w", err)
	}

	log.Printf("[INFO] Trivy scan completed for image: %s. Report saved to: %s", image, outputFile.Name())

	// Validate JSON file contents
	reportContents, err := os.ReadFile(outputFile.Name())
	if err != nil {
		return "", fmt.Errorf("failed to read Trivy report: %w", err)
	}
	if !isValidJSON(reportContents) {
		return "", fmt.Errorf("invalid JSON output in Trivy report")
	}

	return outputFile.Name(), nil
}

func (t *TrivyScanner) ParseReport(reportData []byte) (map[string]interface{}, error) {
	var jsonMap map[string]interface{}
	if err := json.Unmarshal(reportData, &jsonMap); err != nil {
		return nil, fmt.Errorf("failed to parse Trivy report: %w", err)
	}

	return jsonMap, nil
}

//func (t *TrivyScanner) GetExpectedFields() []string {
//	return []string{"Results", "Vulnerabilities", "Severity", "PkgName"}
//}

func (t *TrivyScanner) GenerateFingerprint(image string, reportData []byte, celExpression string) (string, error) {
	parsedReport, err := t.ParseReport(reportData)
	if err != nil {
		return "", fmt.Errorf("failed to parse report: %w", err)
	}

	fmt.Printf("[DEBUG] Evaluating CEL: %s\n", celExpression)

	// Use google.protobuf.Struct for CEL compatibility
	env, err := cel.NewEnv(
		cel.Declarations(
			decls.NewVar("report", decls.NewObjectType("google.protobuf.Struct")),
		),
	)
	if err != nil {
		return "", fmt.Errorf("failed to create CEL environment: %w", err)
	}

	// Compile CEL Expression
	ast, issues := env.Compile(celExpression)
	if issues != nil && issues.Err() != nil {
		return "", fmt.Errorf("CEL compilation error: %w", issues.Err())
	}

	prg, err := env.Program(ast)
	if err != nil {
		return "", fmt.Errorf("failed to create CEL program: %w", err)
	}

	// Execute CEL Expression
	out, _, err := prg.Eval(map[string]interface{}{
		"report": parsedReport,
	})
	if err != nil {
		return "", fmt.Errorf("CEL evaluation error: %w", err)
	}

	// Ensure output is a string
	fingerprint, ok := out.Value().(string)
	if !ok {
		return "", fmt.Errorf("unexpected fingerprint type: %T", out.Value())
	}

	return fingerprint, nil
}

// isValidJSON checks if the provided byte slice is valid JSON
func isValidJSON(s []byte) bool {
	var js interface{}
	return json.Unmarshal(s, &js) == nil
}
