package registry

import (
	"fmt"
)

// NewRegistry initializes the registry based on the provider
func NewRegistry(registry string) (Registry, error) {
	switch registry {
	case "dockerhub":
		return NewDockerHub(), nil
	//case "ecr":
	//	return NewECR(), nil
	case "gcr":
		return NewGCR(), nil
	//case "acr":
	//	return NewACR(), nil
	default:
		return nil, fmt.Errorf("unsupported registry provider: %s", registry)
	}
}
