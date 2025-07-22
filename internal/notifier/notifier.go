package notifier

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/shinobistack/gokakashi/internal/parser"
	"github.com/shinobistack/gokakashi/internal/restapi/v1/scannotify"
	"github.com/shinobistack/gokakashi/pkg/scanner/v1"

	"github.com/google/uuid"
	"github.com/shinobistack/gokakashi/internal/integration/notification"
	"github.com/shinobistack/gokakashi/internal/restapi/v1/integrations"
	"github.com/shinobistack/gokakashi/internal/restapi/v1/scans"
	"golang.org/x/sync/singleflight"
)

func normalizeServer(server string) string {
	if !strings.HasPrefix(server, "http://") && !strings.HasPrefix(server, "https://") {
		server = "http://" + server // Default to HTTP
	}
	return server
}

func constructURL(server string, port int, path string) string {
	base := normalizeServer(server)
	u, err := url.Parse(base)
	if err != nil {
		log.Fatalf("Invalid server URL: %s", base)
	}
	if u.Port() == "" {
		u.Host = fmt.Sprintf("%s:%d", u.Host, port)
	}
	u.Path = path
	return u.String()
}

var notifyGroup = &singleflight.Group{}

func Start(server string, port int, token string, interval time.Duration) {
	log.Println("Starting the periodic notify execution...")
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for range ticker.C {
		_, err, _ := notifyGroup.Do("notify", func() (interface{}, error) {
			return nil, NotifyProcess(server, port, token)
		})
		if err != nil {
			log.Printf("Error in notification process: %v", err)
		}
	}
}

func NotifyProcess(server string, port int, token string) error {
	scans, err := fetchPendingScans(server, port, token, "notify_pending")
	if err != nil {
		log.Printf("Error fetching pending notify: %v", err)
		return err
	}

	if len(scans) == 0 {
		log.Println("No pending notify to execute.")
		return nil
	}

	log.Println("Found", len(scans), "scans to notify.")
	for _, scan := range scans {
		log.Println("Processing scan ID:", scan.ID)
		status, err := processScan(server, port, token, scan)
		if err != nil {
			log.Printf("Error processing scan ID: %s, %v", scan.ID, err)
		}
		err = updateScanStatus(server, port, token, scan.ID, status)
		if err != nil {
			log.Printf("Failed to update status for scanID: %s: %v", scan.ID, err)
		}
		log.Println("Updated status for scanID:", scan.ID, "to", status)
	}
	log.Println("Processed", len(scans), "scans.")

	return nil
}

func processScan(server string, port int, token string, scan scans.GetScanResponse) (string, error) {
	for _, notify := range *scan.Notify {
		scanner, err := scanner.NewScanner(scan.Scanner)
		if err != nil {
			log.Printf("Unsupported scanner tool: %s", scan.Scanner)
			return "error", err
		}

		// Evaluate the 'when' condition using ReportParser
		matched, severities, err := parser.ReportParser(notify.When, &scan)
		if err != nil {
			log.Printf("Error evaluating notify.when:%v, %v", notify.When, err)
			return "error", err
		}

		if matched {
			// Fetch integration details
			// Todo: To define separate schema for scans table to take notify.to as UUID and update all dependent APIs
			parsedNotifyToUUID, err := uuid.Parse(notify.To)
			if err != nil {
				log.Printf("invalid UUID string: %v", err)
				return "error", err
			}

			integration, err := fetchIntegrationDetails(server, port, token, parsedNotifyToUUID)
			if err != nil {
				log.Printf("Error fetching integration details: %v", err)
				return "error", err
			}

			filteredVulnerabilities, err := scanner.FormatReportForNotify(scan.Report, severities, scan.Image)
			if err != nil {
				log.Printf("Error formatting report for notify: %v", err)
				return "error", err
			}

			if len(filteredVulnerabilities) == 0 {
				log.Printf("no vulnerabilities found for scanID: %s and image: %s. Skipping creation of issues", scan.ID, scan.Image)
				return "success", nil
			}

			// Generate a hash and check/save
			var hash string
			if notify.Fingerprint != "" {
				fingerprint, err := scanner.GenerateFingerprint(scan.Image, scan.Report, notify.Fingerprint)
				if err != nil {
					log.Printf("Error generating fingerprint using CEL: %v", err)
					return "error", err
				}
				hash = scanner.GenerateFingerprintHash(fingerprint)
			} else {
				vulnerabilityEntries := scanner.ConvertVulnerabilities(filteredVulnerabilities)
				hash = scanner.GenerateDefaultHash(scan.Image, vulnerabilityEntries)
			}

			occurrences, err := fetchHashCount(server, port, token, hash)
			if err != nil {
				log.Printf("Error fetching occurances of hash: %v", err)
				return "error", err
			}

			if occurrences == nil || occurrences.Count == 0 {
				err := saveHash(server, port, token, scan.ID, hash)
				if err != nil {
					log.Printf("Error saving hash: %v", err)
					return "error", err
				}

				var n notification.Notifier
				switch notification.IntegrationType(integration.Type) {
				case notification.Linear:
					linearIssue, err := constructLinearIssue(integration.Config)
					if err != nil {
						return "error", errors.New("error constructing linear issue")
					}
					linearIssue.Title = fmt.Sprintf("%s - %s", scan.Image, linearIssue.Title)
					linearIssue.Description = scanner.FormatVulnerabilityReport(scan.Image, filteredVulnerabilities)
					n = linearIssue
				default:
					return "error", fmt.Errorf("unknown notifer: %s", integration.Type)
				}
				err = n.Notify(context.TODO())
				if err != nil {
					log.Printf("Error sending notification: %v", err)
				} else {
					return "success", nil
				}
			} else {
				log.Printf("Linear issue exists for image: %s", scan.Image)
				return "success", nil
			}
		}
		if !matched {
			log.Printf("Condition not matched for scanID: %s and image: %s", scan.ID, scan.Image)
			return "success", nil
		}
	}

	return "error", nil
}

func fetchPendingScans(server string, port int, token, status string) ([]scans.GetScanResponse, error) {
	url := constructURL(server, port, "/api/v1/scans") + fmt.Sprintf("?status=%s", status)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request for pending scans: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("server responded with status: %d", resp.StatusCode)
	}

	var listResp scans.ListScansResponse
	if err := json.NewDecoder(resp.Body).Decode(&listResp); err != nil {
		return nil, fmt.Errorf("failed to decode scans response: %w", err)
	}

	return listResp.Scans, nil
}

func fetchIntegrationDetails(server string, port int, token string, integrationID uuid.UUID) (*integrations.GetIntegrationResponse, error) {
	url := constructURL(server, port, "/api/v1/integrations") + fmt.Sprintf("/%s", integrationID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create integration fetch request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch integration details: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("server responded with status: %d", resp.StatusCode)
	}

	var integration integrations.GetIntegrationResponse
	if err := json.NewDecoder(resp.Body).Decode(&integration); err != nil {
		return nil, fmt.Errorf("failed to decode integration response: %w", err)
	}

	return &integration, nil
}

func fetchHashCount(server string, port int, token string, hash string) (*scannotify.GetScanNotifyResponse, error) {
	url := constructURL(server, port, "/api/v1/scannotify") + fmt.Sprintf("?hash=%s", hash)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch hash details: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request for CheckAndSaveHash: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("server responded with status: %d", resp.StatusCode)
	}

	var notifications scannotify.GetScanNotifyResponse
	if err := json.NewDecoder(resp.Body).Decode(&notifications); err != nil {
		return nil, fmt.Errorf("failed to decode notifications response: %w", err)
	}

	return &notifications, nil
}

func saveHash(server string, port int, token string, scanID uuid.UUID, hash string) error {
	// Construct the API URL
	url := constructURL(server, port, "/api/v1/scannotify")

	// Prepare the request payload
	payload := map[string]interface{}{
		"scan_id": scanID,
		"hash":    hash,
	}
	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal hash payload: %w", err)
	}

	// Create the API request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payloadJSON))
	if err != nil {
		return fmt.Errorf("failed to create request for CheckAndSaveHash: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	// Execute the request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute request for CheckAndSaveHash: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusCreated || resp.StatusCode == http.StatusOK {
		// Hash successfully saved
		return nil
	} else {
		return fmt.Errorf("unexpected server response: %d", resp.StatusCode)
	}
}

func constructLinearIssue(config map[string]interface{}) (*notification.LinearIssue, error) {
	// Ensure all fields are present and valid
	apiKey, ok := config["api_key"].(string)
	if !ok {
		return nil, fmt.Errorf("missing or invalid api_key")
	}
	projectID, _ := config["project_id"].(string)
	teamID, _ := config["team_id"].(string)

	title, _ := config["issue_title"].(string)
	priority, _ := config["issue_priority"].(int)
	assignee, _ := config["issue_assignee_id"].(string)
	stateID, _ := config["issue_state_id"].(string)
	dueDate, _ := config["issue_due_date"].(string)

	return &notification.LinearIssue{
		Config: &notification.LinearIntegration{
			APIKey:    apiKey,
			ProjectID: projectID,
			TeamID:    teamID,
		},
		Title:    title,
		Priority: priority,
		Assignee: assignee,
		StateID:  stateID,
		DueDate:  dueDate,
	}, nil
}

func updateScanStatus(server string, port int, token string, scanID uuid.UUID, status string) error {
	url := constructURL(server, port, fmt.Sprintf("/api/v1/scans/%s", scanID))

	reqBody := scans.UpdateScanRequest{
		ID:     scanID,
		Status: strPtr(status),
	}
	reqBodyJSON, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("failed to marshal scan status update request: %w", err)
	}

	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(reqBodyJSON))
	if err != nil {
		return fmt.Errorf("failed to create scan status update request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to update scan status: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("server responded with status: %d", resp.StatusCode)
	}

	return nil
}

func strPtr(s string) *string {
	return &s
}
