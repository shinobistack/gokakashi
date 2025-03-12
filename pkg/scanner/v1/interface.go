package scanner

// Scanner interface defines methods for image scanning
type Scanner interface {
	Scan(image string, severityLevels []string) (string, error)
	// ParseReport For CEL evaluation
	ParseReport(reportData []byte) (map[string]interface{}, error)
	GenerateFingerprint(image string, reportData []byte, celExpression string) (string, error)
}
