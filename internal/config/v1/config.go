package v1

import (
	"fmt"
	"strings"

	"github.com/shinobistack/gokakashi/ent/schema"
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

func (cfg *Config) String() string {
	var s strings.Builder

	s.WriteString("\n")
	s.WriteString("- - - - Configuration - - - - -")
	s.WriteString("\n")

	s.WriteString(fmt.Sprintf("  API Server URL: %s\n", cfg.APIServerURL()))
	if cfg.Site.LogAPITokenOnStartup {
		s.WriteString(fmt.Sprintf("  API Token: %s\n", cfg.Site.APIToken))
	}
	s.WriteString("\n")
	s.WriteString(fmt.Sprintf("  Web Server URL: %s\n", cfg.WebServerURL()))
	s.WriteString("\n")
	s.WriteString(fmt.Sprintf("  Database Host: %s\n", cfg.Database.Host))
	s.WriteString(fmt.Sprintf("  Database Port: %d\n", cfg.Database.Port))
	s.WriteString(fmt.Sprintf("  Database User: %s\n", cfg.Database.User))
	s.WriteString(fmt.Sprintf("  Database Name: %s\n", cfg.Database.Name))
	s.WriteString(fmt.Sprintf("  Database Password: %s\n", strings.Repeat("*", len(cfg.Database.Password))))
	s.WriteString("- - - - - - - - - - - - - - - -")
	s.WriteString("\n")

	return s.String()
}
