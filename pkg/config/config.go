package config

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type Config struct {
	ScanTargets []ScanTarget `yaml:"scan_targets"`
	Website     Website      `yaml:"website"`
}

type ScanTarget struct {
	Registry string    `yaml:"registry"`
	Auth     Auth      `yaml:"auth"`
	Images   []Image   `yaml:"images"`
	Scanner  []Scanner `yaml:"scanner"`
}

type Auth struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

type Image struct {
	Name       string     `yaml:"name"`
	Tags       []string   `yaml:"tags"`
	ScanPolicy ScanPolicy `yaml:"scan_policy"`
}

type ScanPolicy struct {
	Vulnerabilities []string          `yaml:"vulnerabilities"` // e.g. Critical, High
	Notify          map[string]Notify `yaml:"notify"`
	CronSchedule    string            `yaml:"cron_schedule"`
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
	Hostname  string     `yaml:"hostname"`
	FilesPath string     `yaml:"files_path"`
	Public    PortConfig `yaml:"public"`
	Private   PortConfig `yaml:"private"`
}

type PortConfig struct {
	Port int `yaml:"port"`
}

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
