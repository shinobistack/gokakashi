package scanner

import "encoding/json"

// Scanner interface defines methods for image scanning
type Scanner interface {
	Scan(image string, severityLevels []string) (string, error)
	// ParseReport For CEL evaluation
	ParseReport(reportData []byte) (map[string]interface{}, error)
	GenerateFingerprint(image string, reportData []byte, celExpression string) (string, error)
	FormatVulnerabilityReport(image string, vulnerabilities []Vulnerability) string
	FormatReportForNotify(scanReport json.RawMessage, severities []string, scanImage string) ([]Vulnerability, error)
	FilterVulnerabilitiesBySeverity(vulnerabilities []Vulnerability, severityLevels []string) []Vulnerability
	ConvertVulnerabilities(filteredVulnerabilities []Vulnerability) []string
	GenerateDefaultHash(image string, vulnerabilities []string) string
}
