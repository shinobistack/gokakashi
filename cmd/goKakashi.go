package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/ashwiniag/goKakashi/notifier"
	"github.com/ashwiniag/goKakashi/pkg/config"
	"github.com/ashwiniag/goKakashi/pkg/registry"
	"github.com/ashwiniag/goKakashi/pkg/scanner"
	"github.com/ashwiniag/goKakashi/pkg/web"
)

func main() {
	log.Println("=== Starting goKakashi Tool ===")

	// Get config file from command-line argument
	configFile := flag.String("config", "", "Path to the config YAML file")
	flag.Parse()

	if *configFile == "" {
		log.Fatal("Please provide the path to the config YAML file using --config")
	}
	log.Printf("Using configuration file: %s", *configFile)

	// Load the YAML configuration
	log.Println("Loading configuration from YAML file...")
	cfg, err := config.LoadConfig(*configFile)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	log.Println("Configuration loaded successfully.")

	// Process scan targets and images
	for _, target := range cfg.ScanTargets {
		log.Printf("Processing registry: %s", target.Registry)

		// Initialize the appropriate registry
		log.Printf("Initializing registry: %s", target.Registry)
		reg, err := registry.NewRegistry(target.Registry)
		if err != nil {
			log.Fatalf("Failed to initialize registry: %v", err)
		}

		// Handle Docker login if credentials are provided
		log.Printf("Logging in to registry: %s", target.Registry)
		if err := reg.Login(target); err != nil {
			log.Fatalf("Registry login failed: %v", err)
		}
		log.Println("Successfully logged in.")

		// Iterate over the images and scan them
		for _, image := range target.Images {
			for _, tag := range image.Tags {
				imageWithTag := fmt.Sprintf("%s:%s", image.Name, tag)
				log.Printf("Pulling and scanning image: %s", imageWithTag)

				if err := reg.PullImage(imageWithTag); err != nil {
					log.Fatalf("Failed to pull Docker image: %v", err)
				}
				log.Printf("Successfully pulled image: %s", imageWithTag)

				// Initialize the scanner (Trivy)
				trivyScanner := scanner.NewTrivyScanner()

				// Scan the Docker image
				log.Printf("Scanning image: %s", imageWithTag)
				rawReport, vulnerabilities, err := trivyScanner.ScanImage(imageWithTag)
				if err != nil {
					log.Fatalf("Error scanning Docker image: %v", err)
				}
				log.Println("Scan completed successfully.")

				// Filter vulnerabilities based on scan policy
				filteredVulnerabilities := filterVulnerabilities(vulnerabilities, image.ScanPolicy.Vulnerabilities)

				// Notify the user based on the policy
				for _, notifyConfig := range image.ScanPolicy.Notify {
					if notifyConfig.Tool == "Linear" {
						linearNotifier := notifier.NewLinearNotifier()
						err := linearNotifier.SendNotification(filteredVulnerabilities, notifier.NotifyConfig{
							APIKey:    notifyConfig.APIKey,
							ProjectID: notifyConfig.ProjectID,
							Title:     notifyConfig.Title,
							Priority:  notifyConfig.Priority,
							Assignee:  notifyConfig.Assignee,
							Label:     notifyConfig.Label,
							DueDate:   notifyConfig.DueDate,
						})
						if err != nil {
							log.Printf("Failed to send notification: %v", err)
						}
					}
					// Add other notifiers here example jira
				}

				// Save report to file
				restructuredImageName := strings.ReplaceAll(image.Name, "/", "_") // Replace slashes with underscores
				reportFilePath := fmt.Sprintf("%s/%s_%s_report.json", cfg.Website.FilesPath, restructuredImageName, tag)
				log.Printf("Saving report to: %s", reportFilePath)
				err = os.WriteFile(reportFilePath, []byte(rawReport), 0644)
				if err != nil {
					log.Fatalf("Failed to save report: %v", err)
				}
				log.Printf("Report saved successfully at: %s", reportFilePath)
			}
		}
	}

	// Start web servers to serve reports
	log.Println("Starting public and private web servers...")
	go web.StartPublicServer(cfg.Website.FilesPath, cfg.Website.Public.Port)
	go web.StartPrivateServer(cfg.Website.FilesPath, cfg.Website.Private.Port)

	// Graceful shutdown handling
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)
	<-shutdown

	log.Println("Shutting down goKakashi gracefully...")
}

func filterVulnerabilities(vulnerabilities []notifier.Vulnerability, levels []string) []notifier.Vulnerability {
	var filtered []notifier.Vulnerability
	for _, v := range vulnerabilities {
		for _, level := range levels {
			if v.Severity == level {
				filtered = append(filtered, v)
			}
		}
	}
	return filtered
}
