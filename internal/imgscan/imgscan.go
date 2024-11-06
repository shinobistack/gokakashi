package imgscan

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/shinobistack/gokakashi/internal/config/v0"
	"github.com/shinobistack/gokakashi/pkg/scanner"
	"github.com/shinobistack/gokakashi/pkg/utils"
)

var (
	// Store scan_id and its status
	scanStatusStore = make(map[string]string)
	statusMutex     = &sync.Mutex{}
)

type ScanStatus string

const (
	StatusQueued     ScanStatus = "queued"
	StatusInProgress ScanStatus = "in-progress"
	StatusCompleted  ScanStatus = "completed"
	StatusFailed     ScanStatus = "failed"
)

func RunScan(scanID string, image string, severity string, publishTarget string, websites map[string]config.Website) {
	UpdateScanStatus(scanID, StatusInProgress)

	// Simulating long-running process
	//time.Sleep(10 * time.Second)
	imageConfig := config.Image{Name: image}

	// Initialize the scanner (Trivy)
	// Todo: to abstract the scanner. By defaults to Trivy scanning maybe
	trivyScanner := scanner.NewTrivyScanner()

	log.Printf("API Scanning image: %s with severity: %s", image, severity)

	// Use the existing ScanImage function
	report, _, err := trivyScanner.ScanImage(image, strings.Split(severity, ","))
	if err != nil {
		log.Printf("Error scanning image: %v", err)
		UpdateScanStatus(scanID, StatusFailed)
		// TODO: Save the error result
		_, err := saveScanStatus(scanID, string(StatusFailed), nil, "")
		if err != nil {
			log.Println("Error saving scan status", err)
			return
		}
		return
	}
	// On scan completion, update status and save the report
	UpdateScanStatus(scanID, StatusCompleted)

	reportFilePaths, err := utils.SaveScanReport(imageConfig, scanID, report, websites, publishTarget)
	if err != nil {
		log.Printf("Error saving report: %v", err)
		return
	}
	// Get the website configuration for the publish target
	websiteConfig := websites[publishTarget]

	// Save the scan result for future retrieval
	_, err = saveScanStatus(scanID, string(StatusCompleted), reportFilePaths, websiteConfig.ConfiguredDomain)
	if err != nil {
		log.Println("Error saving scan status", err)
		return
	}

	// ToDo: save the report if visibility=private|public and status=completed
}

func saveScanStatus(scanID, status string, reportFilePaths []string, configuredDomain string) ([]string, error) {
	resultFilePath := fmt.Sprintf("/tmp/%s.json", scanID)
	// Build the report URLs for each saved report
	var reportURLs []string
	for _, filePath := range reportFilePaths {
		// Generate the URL from hostname, port, and the file path
		reportURL := fmt.Sprintf("https://%s/reports/%s", configuredDomain, filepath.Base(filePath))
		reportURLs = append(reportURLs, reportURL)
	}
	// Build the scan result
	scanResult := map[string]interface{}{
		"scanID":      scanID,
		"status":      status,
		"report_urls": reportURLs, // List of report URLs
	}
	// Save the result in a JSON file
	jsonData, err := json.Marshal(scanResult)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal scan result: %v", err)
	}

	err = os.WriteFile(resultFilePath, jsonData, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to save result: %v", err)
	}

	log.Printf("Temporary scan result saved to %s", resultFilePath)
	return nil, nil
}

// Update the scan status in a thread-safe manner
func UpdateScanStatus(scanID string, status ScanStatus) {
	statusMutex.Lock()
	defer statusMutex.Unlock()
	scanStatusStore[scanID] = string(status)
}

func GetScanStatus(scanID string) (string, string, error) {
	statusMutex.Lock()
	defer statusMutex.Unlock()

	// Check if the scan ID exists in the in-memory status store
	if status, exists := scanStatusStore[scanID]; exists {
		return scanID, status, nil
	}

	// If not found in memory, check the temporary file
	filePath := fmt.Sprintf("/tmp/%s.json", scanID)
	if _, err := os.Stat(filePath); err == nil {
		// File exists, read the status from the file
		fileData, err := os.ReadFile(filePath)
		if err == nil {
			var result map[string]string
			if json.Unmarshal(fileData, &result) == nil {
				return result["scanID"], result["status"], nil
			}
		}
	}

	// Return an error if scan ID not found
	return "", "", fmt.Errorf("scan ID not found")
}
