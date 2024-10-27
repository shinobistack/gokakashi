package scan

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/ashwiniag/goKakashi/internal/imgscan"
	"github.com/swaggest/usecase/status"
)

type GetRequest struct {
	ScanID string `path:"scan_id"`
}

type GetResponse struct {
	ScanID     string   `json:"scanID"`
	Status     string   `json:"status"`
	ReportURLs []string `json:"report_url,omitempty"` // Optional field
}

func Get(_ context.Context, req GetRequest, res *GetResponse) error {
	_, s, err := imgscan.GetScanStatus(req.ScanID)
	if err != nil {
		return status.Wrap(errors.New("scan id not found"), status.InvalidArgument)
	}

	res.ScanID = req.ScanID
	res.Status = s

	// If the status is completed, add the report URL
	// Todo: to provide correct file path
	if res.Status == string(imgscan.StatusCompleted) {
		reportFilePath := fmt.Sprintf("/tmp/%s.json", req.ScanID)
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
					res.ReportURLs = reportURLs
				}
			}
		}
	}

	return nil
}
