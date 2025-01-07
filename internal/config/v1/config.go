package v1

import (
	"fmt"
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

const ReportsRootDir = "reports/"

// Path to store hash JSON
const HashFilePath = "./hashes.json"

// Integration defines the configuration for external services
type Integration struct {
	Name   string                 `yaml:"name"`
	Type   string                 `yaml:"type"`
	Config map[string]interface{} `yaml:"config"`
}

// SiteConfig defines the API server configuration
type SiteConfig struct {
	APIToken string `yaml:"api_token"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
}

// Trigger specifies the action schedule (cron or CI-based)
type Trigger struct {
	Type     string `yaml:"type"`
	Schedule string `yaml:"schedule,omitempty"`
	// Notify
}

// ImagePolicy defines details about the image to scan
type ImagePolicy struct {
	Registry string   `yaml:"registry"`
	Name     string   `yaml:"name"`
	Tags     []string `yaml:"tags"`
}

// CheckCondition specifies conditions and notification settings
type CheckCondition struct {
	Condition string   `yaml:"condition"`
	Notify    []string `yaml:"notify"`
}

//type Notify struct {
//	OnSuccess []string          `yaml:"on_success,omitempty"`
//	Severity  []string          `yaml:"severity,omitempty"` // To test
//	Labels    map[string]string `yaml:"labels,omitempty"`   // To test - to restrict image level labels at policies
//}

// Policy defines the scanning policies for specific images
type Policy struct {
	Name    string            `yaml:"name"`
	Image   ImagePolicy       `yaml:"image"`
	Trigger Trigger           `yaml:"trigger"`
	Labels  map[string]string `yaml:"labels,omitempty"`
	Check   CheckCondition    `yaml:"check"`
	// ToDo: to update the scanner field to takein tools and tool's argument
	Scanner string `yaml:"scanner"`
}

type DbConnection struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Name     string `yaml:"name"`
}

// Config represents the complete configuration for GoKakashi
type Config struct {
	Integrations []Integration `yaml:"integrations"`
	Site         SiteConfig    `yaml:"site"`
	Policies     []Policy      `yaml:"policies"`
	Database     DbConnection  `yaml:"database"`
}

//type Website struct {
//	Hostname         string `yaml:"hostname"`
//	Port             int    `yaml:"port"`
//	APIToken         string `yaml:"api_token"`
//	Publish          string `yaml:"visibility"`
//	ReportSubDir     string `yaml:"report_sub_dir"`
//	ConfiguredDomain string `yaml:"configured_domain"`
//}

//type Notify struct {
//	APIKey          string `yaml:"api_key"`
//	ProjectID       string `yaml:"project_id"`
//	IssueTitle      string `yaml:"issue_title"`
//	IssuePriority   int    `yaml:"issue_priority"`
//	IssueAssigneeID string `yaml:"issue_assignee_id"`
//	IssueLabel      string `yaml:"issue_label"`
//	IssueDueDate    string `yaml:"issue_due_date"`
//	TeamID          string `yaml:"team_id"`
//	IssueStateID    string `yaml:"issue_state_id"`
//}

type Scanner struct {
	Tool string `yaml:"tool"` // Example: Trivy, Synk, etc.
}

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
		return fmt.Errorf("At least one integration must be defined")
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
