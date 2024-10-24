package scan

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/ashwiniag/goKakashi/internal/imgscan"
	"github.com/gorilla/mux"
)

type GetResponse struct {
	ScanID     string   `json:"scanID"`
	Status     string   `json:"status"`
	ReportURLs []string `json:"report_url,omitempty"` // Optional field
}

func GetHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	scanID := vars["scan_id"]

	log.Printf("Get scan status for scanID %s", scanID)
	_, status, err := imgscan.GetScanStatus(scanID)
	if err != nil {
		http.Error(w, jsonErrorResponse(fmt.Sprintf("Scan ID %s not found", scanID)), http.StatusNotFound)
		return
	}

	// Create the response
	response := GetResponse{
		ScanID: scanID,
		Status: status,
	}

	// If the status is completed, add the report URL
	// Todo: to provide correct file path
	if status == string(imgscan.StatusCompleted) {
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
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Println("Error responsing json", err)
		return
	}
}
