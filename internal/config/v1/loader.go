package v1

import (
	"fmt"
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

func LoadConfig(configFile string) (*Config, error) {
	config := &Config{}

	yamlFile, err := os.ReadFile(configFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %v", err)
	}

	err = yaml.Unmarshal(yamlFile, config)
	if err != nil {
		return nil, fmt.Errorf("error parsing YAML file: %v", err)
	}

	return config, nil
}

// validateConfig validates the loaded configuration to ensure required fields are present
func ValidateConfig(config *Config) error {
	if config.Site.APIToken == "" {
		return fmt.Errorf("API token is missing")
	}

	// ToDo: Minimum one integration must be provided and that needs to be ...?
	if len(config.Integrations) == 0 {
		return fmt.Errorf("at least one integration must be defined")
	}
	for _, integration := range config.Integrations {
		if integration.Name == "" || integration.Type == "" {
			return fmt.Errorf("Integration name and type are required")
		}
	}
	for _, policy := range config.Policies {
		if policy.Name == "" {
			return fmt.Errorf("Policy name is required")
		}
		if policy.Image.Registry == "" || policy.Image.Name == "" {
			return fmt.Errorf("Policy image registry and name are required")
		}
		if policy.Trigger.Type == "cron" && policy.Trigger.Schedule == "" {
			return fmt.Errorf("Policy with cron trigger must define a schedule")
		}
		for _, notify := range policy.Notify {
			if notify.To == "" {
				return fmt.Errorf("Notify 'to' field is required")
			}
			if notify.When == "" {
				return fmt.Errorf("Notify 'when' field is required")
			}
		}
	}
	return nil

	// ToDo: more validations to be added as per config design
	// ToDo: How do we validated the token string for each webserver, currently this is not being used
}

func LoadAndValidateConfig(configPath string) (*Config, error) {
	// Load the YAML configuration
	log.Printf("Loading configuration from YAML file: %s", configPath)
	cfg, err := LoadConfig(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	// Validate the configuration
	if err := ValidateConfig(cfg); err != nil {
		return nil, fmt.Errorf("configuration validation failed: %w", err)
	}
	log.Println("Configuration loaded and validated successfully.")
	return cfg, nil
}
