package scanner

// Scanner interface defines methods for image scanning
type Scanner interface {
	Scan(image string, severityLevels []string) (string, error)
	// ParseReport For CEL evaluation
	ParseReport(reportPath []byte) (map[string]interface{}, error)
	// GetExpectedFields for valid fields
}
