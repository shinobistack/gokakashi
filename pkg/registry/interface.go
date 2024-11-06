package registry

import "github.com/ashwiniag/goKakashi/internal/config/v0"

// Registry interface defines methods for container registries
type Registry interface {
	Login(target config.ScanTarget) error
	PullImage(image string) error
}
