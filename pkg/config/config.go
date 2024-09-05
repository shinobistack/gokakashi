package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	DockerUsername   string
	DockerPassword   string
	DockerImage      string
	RegistryProvider string
	PublicPort       string
	PrivatePort      string
	SkipDockerLogin  bool
	//AWSRegion        string
	//GCRProjectID     string
	//ACRResourceGroup string
	//ACRRegistry      string
}

func LoadConfig() (*Config, error) {
	// Load environment variables from .env file (if it exists)
	err := godotenv.Load()
	if err != nil {
		fmt.Println("No .env file found. Proceeding with environment variables.")
	}

	skipDockerLogin, _ := strconv.ParseBool(os.Getenv("SKIP_DOCKER_LOGIN"))

	cfg := &Config{
		DockerUsername:   os.Getenv("DOCKER_USERNAME"),
		DockerPassword:   os.Getenv("DOCKER_PASSWORD"),
		DockerImage:      os.Getenv("DOCKER_IMAGE"),
		RegistryProvider: os.Getenv("REGISTRY_PROVIDER"),
		PublicPort:       getEnv("PUBLIC_PORT", "8080"),
		PrivatePort:      getEnv("PRIVATE_PORT", "9090"),
		SkipDockerLogin:  skipDockerLogin,
		//AWSRegion:        getEnv("AWS_REGION", "us-east-1"),
		//GCRProjectID:     os.Getenv("GCR_PROJECT_ID"),
		//ACRResourceGroup: os.Getenv("ACR_RESOURCE_GROUP"),
		//ACRRegistry:      os.Getenv("ACR_REGISTRY"),
	}

	// Validate required fields
	if !cfg.SkipDockerLogin && (cfg.DockerUsername == "" || cfg.DockerPassword == "") {
		return nil, fmt.Errorf("DOCKER_USERNAME and DOCKER_PASSWORD are required unless login is skipped")
	}

	if cfg.DockerImage == "" {
		return nil, fmt.Errorf("DOCKER_IMAGE is required")
	}

	if cfg.RegistryProvider == "" {
		return nil, fmt.Errorf("REGISTRY_PROVIDER is required")
	}

	return cfg, nil
}

// getEnv fetches an environment variable or returns the fallback value
func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
