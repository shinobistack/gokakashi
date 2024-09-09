package registry

import "github.com/ashwiniag/goKakashi/pkg/config"

// Registry interface defines methods for container registries
type Registry interface {
	Login(target config.ScanTarget) error
	PullImage(image string) error
}
