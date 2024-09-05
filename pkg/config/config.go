package config

import (
	"fmt"
	"os"
)

type Config struct {
	DockerUsername   string
	DockerPassword   string
	DockerImage      string
	RegistryProvider string
	PublicPort       string
	PrivatePort      string
}

func LoadConfig() (*Config, error) {
	cfg := &Config{
		DockerUsername:   os.Getenv("DOCKER_USERNAME"),
		DockerPassword:   os.Getenv("DOCKER_PASSWORD"),
		DockerImage:      os.Getenv("DOCKER_IMAGE"),
		RegistryProvider: os.Getenv("REGISTRY_PROVIDER"), // dockerhub, ecr, gcr
	}

	// Validate required fields
	if cfg.DockerUsername == "" && cfg.RegistryProvider == "dockerhub" {
		return nil, fmt.Errorf("DOCKER_USERNAME is required for Docker Hub")
	}
	if cfg.DockerPassword == "" && cfg.RegistryProvider == "dockerhub" {
		return nil, fmt.Errorf("DOCKER_PASSWORD is required for Docker Hub")
	}
	if cfg.DockerImage == "" {
		return nil, fmt.Errorf("DOCKER_IMAGE is required")
	}
	if cfg.RegistryProvider == "" {
		return nil, fmt.Errorf("REGISTRY_PROVIDER is required")
	}

	return cfg, nil
}

// getEnv fetches the environment variable or returns the fallback value
func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
