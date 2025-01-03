package registry

import (
	"encoding/json"
	"fmt"
)

func NewRegistry(registry string, config map[string]interface{}) (Integration, error) {
	if len(config) == 0 {
		return nil, fmt.Errorf("configuration is empty for registry provider: %s", registry)
	}

	switch registry {
	case "docker-hub":
		var dockerHubConfig DockerHubIntegration
		configBytes, _ := json.Marshal(config) // Convert map to JSON
		if err := json.Unmarshal(configBytes, &dockerHubConfig); err != nil {
			return nil, fmt.Errorf("failed to parse docker-hub configuration: %w", err)
		}
		return &dockerHubConfig, nil
	case "google-cloud-artifact-registry":
		var gcrConfig GCRIntegration
		configBytes, _ := json.Marshal(config) // Convert map to JSON
		if err := json.Unmarshal(configBytes, &gcrConfig); err != nil {
			return nil, fmt.Errorf("failed to parse google-cloud-artifact-registry configuration: %w", err)
		}
		return &gcrConfig, nil

	default:
		return nil, fmt.Errorf("unsupported registry provider: %s", registry)
	}
}
