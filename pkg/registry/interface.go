package registry

import "github.com/hasura/goKakashi/pkg/config"

// Registry interface defines methods for container registries
type Registry interface {
	Login(cfg *config.Config) error
	PullImage(image string) error
}
