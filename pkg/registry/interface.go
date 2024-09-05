package registry

import "github.com/ashwiniag/goKakashi/pkg/config"

// Registry interface defines methods for container registries
type Registry interface {
	Login(cfg *config.Config) error
	PullImage(image string) error
}
