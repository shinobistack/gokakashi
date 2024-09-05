package scanner

// Scanner interface defines methods for image scanning
type Scanner interface {
	ScanImage(image string) (string, error)
}
