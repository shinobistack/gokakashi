package api

import (
	"encoding/json"
	_ "encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	_ "strings"
	"sync"
	"time"

	"github.com/ashwiniag/goKakashi/pkg/config"
	"github.com/ashwiniag/goKakashi/pkg/scanner"
	"github.com/ashwiniag/goKakashi/pkg/utils"
	_ "github.com/ashwiniag/goKakashi/pkg/utils"
)

var (
	// Store scan_id and its status
	scanStatusStore = make(map[string]string)
	statusMutex     = &sync.Mutex{}
)

type ScanResponse struct {
	ScanID string `json:"scan_id"`
	Status string `json:"status"`
	// Result
}

type ScanStatus string

const (
	StatusQueued     ScanStatus = "queued"
	StatusInProgress ScanStatus = "in-progress"
	StatusCompleted  ScanStatus = "completed"
	StatusFailed     ScanStatus = "failed"
)

// Generate a unique scan ID based on the current timestamp
func generateScanID() string {
	return fmt.Sprintf("scan-%d", time.Now().UnixNano())
}

// Create a unified JSON error response
func jsonErrorResponse(message string) string {
	response := map[string]string{"error": message}
	jsonResponse, _ := json.Marshal(response) // Ignore error for simplicity
	return string(jsonResponse)
}

// Update the scan status in a thread-safe manner
func updateScanStatus(scanID string, status ScanStatus) {
	statusMutex.Lock()
	defer statusMutex.Unlock()
	scanStatusStore[scanID] = string(status)
}

// StartSingleImageScan POST /api/v0/scan?image=<>&severity=<>&publish=<>
func StartScan(w http.ResponseWriter, r *http.Request, websites map[string]config.Website) {
	image := r.URL.Query().Get("image")
	if image == "" {
		http.Error(w, jsonErrorResponse("Image is missing"), http.StatusBadRequest)
		return
	}

	// Ensure it can take severity is either HIGH or CRITICAL or HIGH,CRITICAL
	severity := r.URL.Query().Get("severity")
	if severity == "" || !(severity == "HIGH" || severity == "CRITICAL" || severity == "HIGH,CRITICAL") {
		http.Error(w, jsonErrorResponse("Severity must be HIGH, CRITICAL, or HIGH,CRITICAL"), http.StatusBadRequest)
		return
	}
	// Publish to mentioned website.ReportSubDir
	publishTarget := r.URL.Query().Get("publish")
	if publishTarget == "" {
		http.Error(w, jsonErrorResponse("publish field is missing, report will not be saved"), http.StatusBadRequest)
		return
	}

	// Unique scan ID
	scanID := generateScanID()
	updateScanStatus(scanID, StatusQueued)

	log.Printf("Initiating scan for image %s with severity %s", image, severity)

	// Start the scan asynchronously
	go runScan(scanID, image, severity, publishTarget, websites)

	response := ScanResponse{
		ScanID: scanID,
		Status: string(StatusQueued),
	}
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Println("Error responding json", err)
		return
	}
}

func runScan(scanID string, image string, severity string, publishTarget string, websites map[string]config.Website) {
	updateScanStatus(scanID, StatusInProgress)

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
		updateScanStatus(scanID, StatusFailed)
		// TODO: Save the error result
		_, err := saveScanStatus(scanID, string(StatusFailed), nil, "")
		if err != nil {
			log.Println("Error saving scan status", err)
			return
		}
		return
	}
	// On scan completion, update status and save the report
	updateScanStatus(scanID, StatusCompleted)

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
