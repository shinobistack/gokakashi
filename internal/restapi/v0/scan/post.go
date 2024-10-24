package scan

import (
	"encoding/json"
	_ "encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	_ "strings"
	"time"

	"github.com/ashwiniag/goKakashi/internal/imgscan"
	"github.com/ashwiniag/goKakashi/pkg/config"
	_ "github.com/ashwiniag/goKakashi/pkg/utils"
	"github.com/scriptnull/jsonseal"
	"golang.org/x/exp/maps"
)

type PostRequest struct {
	Image    string `json:"image"`
	Severity string `json:"severity"`
	Publish  string `json:"publish"`
}

type PostResponse struct {
	ScanID string `json:"scan_id"`
	Status string `json:"status"`
}

type PostHandler struct {
	Websites map[string]config.Website
}

var (
	ErrNotFound error = errors.New("not found")
)

func (req *PostRequest) Validate() error {
	var check jsonseal.CheckGroup

	check.Field("image").Check(func() error {
		if req.Image == "" {
			return ErrNotFound
		}
		return nil
	})

	check.Field("severity").Check(func() error {
		allowedSev := map[string]struct{}{"HIGH": {}, "CRITICAL": {}}
		validSev := true
		for _, sev := range strings.Split(req.Severity, "websites,") {
			if _, exists := allowedSev[sev]; !exists {
				validSev = false
				break
			}
		}
		if !validSev {
			return fmt.Errorf("severity must be %s", strings.Join(maps.Keys(allowedSev), ","))
		}
		return nil
	})

	check.Field("publish").Check(func() error {
		if req.Publish == "" {
			return ErrNotFound
		}
		return nil
	})

	return check.Validate()
}

func (ph *PostHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var req PostRequest

	err := jsonseal.NewDecoder(r.Body).DecodeValidate(&req)
	if err != nil {
		http.Error(w, jsonErrorResponse(jsonseal.JSONFormat(err)), http.StatusBadRequest)
		return
	}

	scanID := generateScanID()
	imgscan.UpdateScanStatus(scanID, imgscan.StatusQueued)

	log.Printf("Initiating scan for image %s with severity %s", req.Image, req.Severity)
	go imgscan.RunScan(scanID, req.Image, req.Severity, req.Publish, ph.Websites)

	response := PostResponse{
		ScanID: scanID,
		Status: string(imgscan.StatusQueued),
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Println("Error responding json", err)
		return
	}
}

// Generate a unique scan ID based on the current timestamp
func generateScanID() string {
	return fmt.Sprintf("scan-%d", time.Now().UnixNano())
}

// jsonErrorResponse creates a unified JSON error response
func jsonErrorResponse(message string) string {
	response := map[string]string{"error": message}
	jsonResponse, _ := json.Marshal(response) // Ignore error for simplicity
	return string(jsonResponse)
}
