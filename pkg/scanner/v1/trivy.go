package scanner

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/checker/decls"
	"github.com/google/cel-go/ext"
	"log"
	"os"
	"os/exec"
	"strings"
)

type HashEntry struct {
	Image           string              `json:"image"`           // Image name
	Vulnerabilities []VulnerabilityData `json:"vulnerabilities"` // Detailed vulnerability data
	Hash            string              `json:"hash"`            // Generated hash for the entry
}

type VulnerabilityData struct {
	VulnerabilityID  string `json:"vulnerability_id"`  // The CVE or vulnerability ID
	Severity         string `json:"severity"`          // Severity level (e.g., Critical, High)
	InstalledVersion string `json:"installed_version"` // Version of the package installed
	FixedVersion     string `json:"fixed_version"`     // Version where the vulnerability is fixed (if available)
}

// TrivyReport represents the overall Trivy scan report
type Report struct {
	ArtifactName string   `json:"ArtifactName"`
	Results      []Result `json:"Results"`
}

// Result represents the result field in Trivy output
type Result struct {
	Target          string          `json:"Target"`
	Type            string          `json:"Type"`
	Vulnerabilities []Vulnerability `json:"Vulnerabilities"`
}

// Vulnerability represents a vulnerability in the scan results
type Vulnerability struct {
	VulnerabilityID  string `json:"VulnerabilityID"`
	PkgName          string `json:"PkgName"`
	Severity         string `json:"Severity"`
	InstalledVersion string `json:"InstalledVersion"`
	FixedVersion     string `json:"FixedVersion"`
	Title            string `json:"Title"`
	Description      string `json:"Description"`
	PrimaryURL       string `json:"PrimaryURL"`
	Status           string `json:"Status"`
}

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
		ext.Strings(), // Enables string functions (split, join)
		ext.Lists(),   // Enables list functions (map, filter, flatMap)
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

func (t *TrivyScanner) FormatVulnerabilityReport(image string, vulnerabilities []Vulnerability) string {
	var buffer bytes.Buffer

	// Add image information
	buffer.WriteString(fmt.Sprintf("Image: %s\n\n", image))

	// Iterate over vulnerabilities and format them in the simplified format
	for _, vuln := range vulnerabilities {
		buffer.WriteString(fmt.Sprintf("Library: %s\n", vuln.PkgName))
		buffer.WriteString(fmt.Sprintf("Vulnerability: %s\n", vuln.VulnerabilityID))
		buffer.WriteString(fmt.Sprintf("Severity: %s\n", vuln.Severity))
		buffer.WriteString(fmt.Sprintf("Status: %s\n", vuln.Status))
		buffer.WriteString(fmt.Sprintf("Installed Version: %s\n", vuln.InstalledVersion))
		buffer.WriteString(fmt.Sprintf("Fixed Version: %s\n", vuln.FixedVersion))
		buffer.WriteString(fmt.Sprintf("Title: %s\n", vuln.Title))
		if vuln.PrimaryURL != "" {
			buffer.WriteString(fmt.Sprintf("More details: %s\n", vuln.PrimaryURL))
		}
		buffer.WriteString("\n") // Add a line break between vulnerabilities
	}

	return buffer.String()
}

func (t *TrivyScanner) FormatReportForNotify(scanReport json.RawMessage, severities []string, scanImage string) ([]Vulnerability, error) {
	var report Report
	err := json.Unmarshal(scanReport, &report)
	if err != nil {
		log.Printf("Error failed to parse scan report: %v", err)
	}

	// Load the vulnerabilities from scans.report
	var vulnerabilities []Vulnerability
	for _, result := range report.Results {
		vulnerabilities = append(vulnerabilities, result.Vulnerabilities...)
	}

	filteredVulnerabilities := t.FilterVulnerabilitiesBySeverity(vulnerabilities, severities)

	return filteredVulnerabilities, nil
}

func (t *TrivyScanner) FilterVulnerabilitiesBySeverity(vulnerabilities []Vulnerability, severityLevels []string) []Vulnerability {
	var filtered []Vulnerability
	for _, v := range vulnerabilities {
		for _, level := range severityLevels {
			if v.Severity == level {
				filtered = append(filtered, v)
			}
		}
	}
	return filtered
}

func (t *TrivyScanner) ConvertVulnerabilities(filteredVulnerabilities []Vulnerability) []string {
	var vulnerabilityEntries []string
	for _, v := range filteredVulnerabilities {
		data := VulnerabilityData{
			VulnerabilityID:  v.VulnerabilityID,
			Severity:         v.Severity,
			InstalledVersion: v.InstalledVersion,
			FixedVersion:     v.FixedVersion,
		}
		entry := fmt.Sprintf("%s_%s_%s_%s", data.VulnerabilityID, data.Severity, data.InstalledVersion, data.FixedVersion)
		vulnerabilityEntries = append(vulnerabilityEntries, entry)
	}
	return vulnerabilityEntries
}

func (t *TrivyScanner) GenerateDefaultHash(image string, vulnerabilities []string) string {
	data := fmt.Sprintf("%s_%s", image, strings.Join(vulnerabilities, "_"))
	hash := sha256.New()
	hash.Write([]byte(data))
	return hex.EncodeToString(hash.Sum(nil))
}

// isValidJSON checks if the provided byte slice is valid JSON
func isValidJSON(s []byte) bool {
	var js interface{}
	return json.Unmarshal(s, &js) == nil
}
