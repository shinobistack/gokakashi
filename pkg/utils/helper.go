package utils

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/shinobistack/gokakashi/notifier"

	v0 "github.com/shinobistack/gokakashi/internal/config/v0"
	"github.com/shinobistack/gokakashi/pkg/registry"
	"github.com/shinobistack/gokakashi/pkg/scanner"
)

const reportsRootDir = "reports/"

// Todo: to re-arrange and restructure

// InitializeRegistry initializes the Docker registry and performs login if necessary.
func InitializeRegistry(target v0.ScanTarget) (registry.Registry, error) {
	log.Printf("Initializing registry: %s", target.Registry)
	reg, err := registry.NewRegistry(target.Registry)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize registry: %v", err)
	}

	log.Printf("Logging in to registry: %s", target.Registry)
	if err := reg.Login(target); err != nil {
		return nil, fmt.Errorf("registry login failed: %v", err)
	}
	log.Println("Successfully logged in.")
	return reg, nil
}

// PullAndScanImage pulls the Docker image and runs the scan using Trivy.
func PullAndScanImage(reg registry.Registry, image v0.Image, tag string, severityLevels []string) (string, []notifier.Vulnerability, error) {
	imageWithTag := fmt.Sprintf("%s:%s", image.Name, tag)
	log.Printf("Pulling and scanning image: %s", imageWithTag)

	if err := reg.PullImage(imageWithTag); err != nil {
		return "", nil, fmt.Errorf("failed to pull Docker image: %v", err)
	}
	log.Printf("Successfully pulled image: %s", imageWithTag)

	// ToDo: refurb the scanner
	// Initialize the scanner (Trivy)
	trivyScanner := scanner.NewTrivyScanner()

	// Scan the Docker image using Trivy
	report, vulnerabilities, err := trivyScanner.ScanImage(imageWithTag, severityLevels)
	if err != nil {
		return "", nil, fmt.Errorf("error scanning Docker image: %v", err)
	}
	log.Println("Scan completed successfully.")

	return report, vulnerabilities, nil
}

// SaveScanReport saves the scan report to a file.
func SaveScanReport(image v0.Image, tag, report string, websites map[string]v0.Website, apiPublishTarget string) ([]string, error) {
	var savedPaths []string
	restructuredImageName := strings.ReplaceAll(image.Name, "/", "_") // Replace slashes with underscores
	// if publish != empty then go on create file path   if doesn't exists and save report'
	// else skip saving report
	// If an API-specified publish target is provided, use it first
	publishTargets := image.Publish
	if apiPublishTarget != "" {
		publishTargets = []string{apiPublishTarget}
	}
	if len(publishTargets) > 0 {
		for _, publishTarget := range publishTargets {
			websiteConfig, ok := websites[publishTarget]
			if !ok {
				log.Printf("Publish target %s not found in website config. Skipping...", publishTarget)
				continue
			}
			reportFilePath := filepath.Join(reportsRootDir, websiteConfig.ReportSubDir, fmt.Sprintf("%s_%s_report.json", restructuredImageName, tag))
			// It creates the directroy defined under cfg.Website.ReportSubDir if it doesn't exist
			if err := os.MkdirAll(filepath.Dir(reportFilePath), 0755); err != nil {
				return nil, fmt.Errorf("failed to create directory for report: %v", err)
			}

			// Save the report
			if err := os.WriteFile(reportFilePath, []byte(report), 0644); err != nil {
				return nil, fmt.Errorf("failed to write report file: %v", err)
			}
			log.Printf("Report saved successfully at: %s", reportFilePath)
			savedPaths = append(savedPaths, reportFilePath)
		}
	} else {
		log.Printf("No publish target specified for image %s:%s. Skipping saving report.", image.Name, tag)
	}

	if len(savedPaths) == 0 {
		return nil, fmt.Errorf("no reports saved for image %s:%s", image.Name, tag)
	}
	return savedPaths, nil
}

// FilterVulnerabilitiesBySeverity filters vulnerabilities based on the provided severity levels
func FilterVulnerabilitiesBySeverity(vulnerabilities []notifier.Vulnerability, severityLevels []string) []notifier.Vulnerability {
	var filtered []notifier.Vulnerability
	for _, v := range vulnerabilities {
		for _, level := range severityLevels {
			if v.Severity == level {
				filtered = append(filtered, v)
			}
		}
	}
	return filtered
}

// CheckAndSaveHash checks if the hash already exists and saves it if not.
func CheckAndSaveHash(image v0.Image, tag string, vulnerabilities []notifier.Vulnerability) (bool, error) {
	vulnerabilityData, vulnerabilityEntries := ConvertVulnerabilities(vulnerabilities)
	hash := GenerateHash(image.Name, tag, vulnerabilityEntries)

	exists, err := HashExists(hashFilePath, hash)
	if err != nil {
		return false, fmt.Errorf("error checking hash: %v", err)
	}

	if exists {
		log.Printf("Hash already exists for %s:%s, skipping ticket creation", image.Name, tag)
		return true, nil
	}

	entry := HashEntry{
		Image:           image.Name,
		Tag:             tag,
		Vulnerabilities: vulnerabilityData,
		Hash:            hash,
	}

	if err := SaveHashToFile(hashFilePath, entry); err != nil {
		return false, fmt.Errorf("error saving hash: %v", err)
	}

	return false, nil
}

func RunImageScan(target v0.ScanTarget, image v0.Image, cfg *v0.Config) error {
	log.Printf("Processing registry: %s", target.Registry)

	// Initialize the registry
	reg, err := InitializeRegistry(target)
	if err != nil {
		log.Fatalf("Error: %v", err)
		return nil
	}

	for _, tag := range image.Tags {
		// Pull and scan the image
		report, vulnerabilities, err := PullAndScanImage(reg, image, tag, image.ScanPolicy.Severity)
		if err != nil {
			log.Printf("Error: %v", err)
			continue
		}

		// Save the report
		_, err = SaveScanReport(image, tag, report, cfg.Websites, "")
		if err != nil {
			log.Fatalf("Error saving report: %v", err)
			return nil
		}

		// Filter vulnerabilities based on severity levels
		filteredVulnerabilities := FilterVulnerabilitiesBySeverity(vulnerabilities, image.ScanPolicy.Severity)

		// Skip ticket creation if no matching vulnerabilities found
		if len(filteredVulnerabilities) == 0 {
			log.Printf("No vulnerabilities found for image: %s:%s", image.Name, tag)
			continue
		}

		// Check for hash and save if new
		hashExists, err := CheckAndSaveHash(image, tag, filteredVulnerabilities)
		if err != nil {
			log.Fatalf("Error: %v", err)
			return nil
		}
		if hashExists {
			continue
		}

		// Notify via configured notification channels
		Notify(image, tag, filteredVulnerabilities)
	}
	return nil
}

func Notify(image v0.Image, tag string, vulnerabilities []notifier.Vulnerability) {
	for toolName, notifyConfig := range image.ScanPolicy.Notify {
		if toolName == "Linear" {
			linearNotifier := notifier.NewLinearNotifier()
			finalTitle := fmt.Sprintf("%s %s", fmt.Sprintf("%s:%s", image.Name, tag), notifyConfig.IssueTitle)
			err := linearNotifier.SendNotification(notifier.TrivyReport{
				ArtifactName: fmt.Sprintf("%s:%s", image.Name, tag),
				Results:      []notifier.Result{},
			}, vulnerabilities, notifier.NotifyConfig{
				APIKey:    notifyConfig.APIKey,
				TeamID:    notifyConfig.TeamID,
				ProjectID: notifyConfig.ProjectID,
				Title:     finalTitle,
				Priority:  notifyConfig.IssuePriority,
				Assignee:  notifyConfig.IssueAssigneeID,
				StateID:   notifyConfig.IssueStateID,
				DueDate:   notifyConfig.IssueDueDate,
			})
			if err != nil {
				log.Printf("Failed to send notification: %v", err)
			}
		}
		// Add other notifiers (e.g., Jira) here
	}
}
