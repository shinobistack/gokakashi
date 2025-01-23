package v1

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/shinobistack/gokakashi/ent/schema"

	"gopkg.in/yaml.v2"
)

// Integration defines the configuration for external services
type Integration struct {
	Name   string                 `yaml:"name"`
	Type   string                 `yaml:"type"`
	Config map[string]interface{} `yaml:"config"`
}

// SiteConfig defines the API server configuration
type SiteConfig struct {
	APIToken             string `yaml:"api_token"`
	LogAPITokenOnStartup bool   `yaml:"log_api_token_on_startup"`
	Host                 string `yaml:"host"`
	Port                 int    `yaml:"port"`
}

type WebServerConfig struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
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

type Notify = schema.Notify

// Policy defines the scanning policies for specific images
type Policy struct {
	Name    string            `yaml:"name"`
	Image   ImagePolicy       `yaml:"image"`
	Trigger Trigger           `yaml:"trigger"`
	Labels  map[string]string `yaml:"labels,omitempty"`
	Notify  []Notify          `yaml:"notify,omitempty"`
	Scanner string            `yaml:"scanner"`
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
	Integrations []Integration   `yaml:"integrations"`
	Site         SiteConfig      `yaml:"site"`
	Policies     []Policy        `yaml:"policies"`
	Database     DbConnection    `yaml:"database"`
	WebServer    WebServerConfig `yaml:"web_server"`
}

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

func DefaultConfig() (*Config, error) {
	apiToken := os.Getenv("GOKAKASHI_API_SERVER_TOKEN")
	logAPITokenOnStartup := false
	if apiToken == "" {
		generatedToken, err := generateToken(32)
		if err != nil {
			return nil, fmt.Errorf("error generating an api token: %w", err)
		}
		apiToken = generatedToken
		logAPITokenOnStartup = true
	}

	apiHost := os.Getenv("GOKAKASHI_API_SERVER_HOST")
	apiPort, _ := strconv.Atoi(os.Getenv("GOKAKASHI_API_SERVER_PORT"))
	if apiPort == 0 {
		apiPort = 5556
	}

	webHost := os.Getenv("GOKAKASHI_WEB_SERVER_HOST")
	webPort, _ := strconv.Atoi(os.Getenv("GOKAKASHI_WEB_SERVER_PORT"))
	if webPort == 0 {
		webPort = 5555
	}

	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" {
		dbHost = "localhost"
	}
	dbPort, _ := strconv.Atoi(os.Getenv("DB_PORT"))
	if dbPort == 0 {
		dbPort = 5432
	}
	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		dbName = "postgres"
	}
	dbUser := os.Getenv("DB_USER")
	if dbUser == "" {
		dbUser = "postgres"
	}
	dbPassword := os.Getenv("DB_PASSWORD")
	if dbPassword == "" {
		dbPassword = "postgres"
	}

	return &Config{
		Site: SiteConfig{
			APIToken:             apiToken,
			LogAPITokenOnStartup: logAPITokenOnStartup,
			Host:                 apiHost,
			Port:                 apiPort,
		},
		WebServer: WebServerConfig{
			Host: webHost,
			Port: webPort,
		},
		Database: DbConnection{
			Host:     dbHost,
			Port:     dbPort,
			User:     dbUser,
			Password: dbPassword,
			Name:     dbName,
		},
	}, nil
}

func (c *Config) WebServerURL() string {
	host := c.WebServer.Host
	if host == "" {
		host = "localhost"
	}
	return fmt.Sprintf("http://%s:%d", host, c.WebServer.Port)
}

func (c *Config) APIServerURL() string {
	host := c.Site.Host
	if host == "" {
		host = "localhost"
	}
	return fmt.Sprintf("http://%s:%d", host, c.Site.Port)
}

func generateToken(length int) (string, error) {
	token := make([]byte, length)

	_, err := rand.Read(token)
	if err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(token), nil
}
