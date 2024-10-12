package api

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
)

type ScanStatusResponse struct {
	ScanID     string   `json:"scanID"`
	Status     string   `json:"status"`
	ReportURLs []string `json:"report_url,omitempty"` // Optional field
}

func StatusHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	scanID := vars["scan_id"]
	log.Printf("Get scan status for scanID %s", scanID)
	// Retrieve the status of the scan
	_, status, err := getScanStatus(scanID)
	if err != nil {
		http.Error(w, jsonErrorResponse(fmt.Sprintf("Scan ID %s not found", scanID)), http.StatusNotFound)
		return
	}

	// Create the response
	response := ScanStatusResponse{
		ScanID: scanID,
		Status: status,
	}

	// If the status is completed, add the report URL
	// Todo: to provide correct file path
	if status == string(StatusCompleted) {
		reportFilePath := fmt.Sprintf("/tmp/%s.json", scanID)
		fileData, err := os.ReadFile(reportFilePath)
		if err == nil {
			var result map[string]interface{}
			if json.Unmarshal(fileData, &result) == nil {
				if urls, ok := result["report_urls"].([]interface{}); ok && len(urls) > 0 {
					// Convert []interface{} to []string
					reportURLs := make([]string, len(urls))
					for i, url := range urls {
						if strURL, ok := url.(string); ok {
							reportURLs[i] = strURL
						}
					}
					// Set all URLs in the response
					response.ReportURLs = reportURLs
				}
			}
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func getScanStatus(scanID string) (string, string, error) {
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
