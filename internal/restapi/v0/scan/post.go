package scan

import (
	"context"
	_ "encoding/json"
	"errors"
	"fmt"
	"strings"
	_ "strings"
	"time"

	"github.com/scriptnull/jsonseal"
	"github.com/shinobistack/gokakashi/internal/config/v0"
	"github.com/shinobistack/gokakashi/internal/imgscan"
	_ "github.com/shinobistack/gokakashi/pkg/utils"
	"github.com/swaggest/usecase/status"
	"golang.org/x/exp/maps"
)

type PostRequest struct {
	Image    string `json:"image"`
	Severity string `json:"severity"`
	Publish  string `json:"publish"`
}

type PostResponse struct {
	ScanID string             `json:"scan_id"`
	Status imgscan.ScanStatus `json:"status"`
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

func Post(_ context.Context, req PostRequest, res *PostResponse) error {
	err := req.Validate()
	if err != nil {
		return status.Wrap(err, status.InvalidArgument)
	}

	res.ScanID = generateScanID()
	res.Status = imgscan.StatusQueued
	imgscan.UpdateScanStatus(res.ScanID, res.Status)

	return nil
}

// Generate a unique scan ID based on the current timestamp
func generateScanID() string {
	return fmt.Sprintf("scan-%d", time.Now().UnixNano())
}
