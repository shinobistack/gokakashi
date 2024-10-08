package config

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

const ReportsRootDir = "reports/"

// Path to store hash JSON
const HashFilePath = "./hashes.json"

type Config struct {
	ScanTargets []ScanTarget       `yaml:"scan_targets"`
	Websites    map[string]Website `yaml:"websites"`
	// ToDo: To remove this and maybe validate for each webserver defined under Website
	APIToken string `yaml:"api_token"`
}

type ScanTarget struct {
	Registry string    `yaml:"registry"`
	Auth     Auth      `yaml:"auth"`
	Images   []Image   `yaml:"images"`
	Scanner  []Scanner `yaml:"scanner"`
}

type Auth struct {
	Type        string `yaml:"type"` // Types of authentication
	Username    string `yaml:"username"`
	Password    string `yaml:"password"`
	JSONKeyPath string `yaml:"json_key_path"` // Optional: For GCR service account
}

type Image struct {
	Name       string     `yaml:"name"`
	Tags       []string   `yaml:"tags"`
	Publish    []string   `yaml:"publish"`
	ScanPolicy ScanPolicy `yaml:"scan_policy"`
}

type ScanPolicy struct {
	Severity     []string          `yaml:"severity"` // e.g. Critical, High
	Notify       map[string]Notify `yaml:"notify"`
	CronSchedule string            `yaml:"cron_schedule"`
}

type Notify struct {
	APIKey          string `yaml:"api_key"`
	ProjectID       string `yaml:"project_id"`
	IssueTitle      string `yaml:"issue_title"`
	IssuePriority   int    `yaml:"issue_priority"`
	IssueAssigneeID string `yaml:"issue_assignee_id"`
	IssueLabel      string `yaml:"issue_label"`
	IssueDueDate    string `yaml:"issue_due_date"`
	TeamID          string `yaml:"team_id"`
	IssueStateID    string `yaml:"issue_state_id"`
}

type Scanner struct {
	Tool string `yaml:"tool"` // Example: Trivy, Synk, etc.
}

type Website struct {
	Hostname string `yaml:"hostname"`
	//FilesPath string     `yaml:"files_path"`
	Port int `yaml:"port"`
	// ToDo: How do we validated the token string for each webserver, currently this is not being used
	APIToken string `yaml:"api_token"`
	// ToDo: define to which webserver config user wants to publish to
	Publish      string `yaml:"visibility"`
	ReportSubDir string `yaml:"report_sub_dir"`
}

//type PortConfig struct {
//	Port int `yaml:"port"`
//}

func LoadConfig(configFile string) (*Config, error) {
	config := &Config{}

	yamlFile, err := ioutil.ReadFile(configFile)
	if err != nil {
		return nil, fmt.Errorf("Failed to read config file: %v", err)
	}

	err = yaml.Unmarshal(yamlFile, config)
	if err != nil {
		return nil, fmt.Errorf("Error parsing YAML file: %v", err)
	}

	return config, nil
}

// validateConfig validates the loaded configuration to ensure required fields are present
func ValidateConfig(cfg *Config) error {
	//Ensure the Websites fields are not empty
	if len(cfg.Websites) == 0 {
		return fmt.Errorf("Websites cannot be empty")
	}

	// Validate scan targets
	if len(cfg.ScanTargets) == 0 {
		return fmt.Errorf("No scan targets specified in the configuration")
	}

	for _, target := range cfg.ScanTargets {
		if target.Registry == "" {
			return fmt.Errorf("Registry for scan target cannot be empty")
		}
		if len(target.Images) == 0 {
			return fmt.Errorf("No images specified for registry: %s", target.Registry)
		}
	}

	return nil
	// ToDo: more validations to be added as per config design
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
