package scanner

// Scanner interface defines methods for image scanning
type Scanner interface {
	Scan(image string, severityLevels []string) (string, error)
}
