package notifier

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/shinobistack/gokakashi/pkg/scanner/v1"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/shinobistack/gokakashi/internal/integration/notification"
	"github.com/shinobistack/gokakashi/internal/parser"
	"github.com/shinobistack/gokakashi/internal/restapi/v1/integrations"
	"github.com/shinobistack/gokakashi/internal/restapi/v1/scans"
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

func Start(server string, port int, token string, interval time.Duration) {
	log.Println("Starting the periodic notify execution...")
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for range ticker.C {
		NotifyProcess(server, port, token)
	}

}

func NotifyProcess(server string, port int, token string) {
	scans, err := fetchPendingScans(server, port, token, "notify_pending")
	if err != nil {
		log.Printf("Notifier: Error fetching pending notify: %v", err)
	}

	if len(scans) == 0 {
		log.Println("Notifier: No pending notify to execute.")
		return
	}

	for _, scan := range scans {
		for _, notify := range *scan.Notify {
			scanner, err := scanner.NewScanner(scan.Scanner)
			if err != nil {
				log.Printf("Notifier: Unsupported scanner tool: %s", scan.Scanner)
				continue
			}
			// Evaluate the 'when' condition using ReportParser
			matched, severities, err := parser.ReportParser(notify.When, &scan)
			if err != nil {
				log.Printf("Error evaluating notify.when:%v, %v", notify.When, err)
				continue
			}

			if matched {
				// Fetch integration details
				// Todo: To define separate schema for scans table to take notify.to as UUID and update all dependent APIs
				parsedNotifyToUUID, err := uuid.Parse(notify.To)
				if err != nil {
					log.Printf("Notifier: invalid UUID string: %v", err)
				}

				integration, err := fetchIntegrationDetails(server, port, token, parsedNotifyToUUID)
				if err != nil {
					log.Printf("Error fetching integration details: %v", err)
					continue
				}

				filteredVulnerabilities, err := scanner.FormatReportForNotify(scan.Report, severities, scan.Image)
				if err != nil {
					log.Printf("Notifier: Error formatting report for notify: %v", err)
					continue
				}

				if len(filteredVulnerabilities) == 0 {
					log.Printf("Notifier: no vulnerabilities found for scanID: %s and image: %s. Skipping creation of issues", scan.ID, scan.Image)
					err = updateScanStatus(server, port, token, scan.ID, "success")
					if err != nil {
						log.Printf("Notifier: Error updating scan status: %v", err)
					}
					return
				}

				// Generate a hash and check/save
				var hash string
				if notify.Fingerprint != "" {
					hash, err = scanner.GenerateFingerprint(scan.Image, scan.Report, notify.Fingerprint)
					if err != nil {
						log.Printf("Notifier: Error generating fingerprint using CEL: %v", err)
						continue
					}
				} else {
					vulnerabilityEntries := scanner.ConvertVulnerabilities(filteredVulnerabilities)
					hash = scanner.GenerateHash(scan.Image, vulnerabilityEntries)
				}

				saved, err := CheckAndSaveHash(server, port, token, scan.ID, hash)
				if err != nil {
					log.Printf("Notifier: Error checking or saving hash: %v", err)
					continue
				}

				if saved {
					var n notification.Notifier
					switch notification.IntegrationType(integration.Type) {
					case notification.Linear:
						linearIssue, err := constructLinearIssue(integration.Config)
						if err != nil {
							log.Printf("Notifier: Error constructing linear issue: %v", err)
							continue
						}
						linearIssue.Title = fmt.Sprintf("%s - %s", scan.Image, linearIssue.Title)
						linearIssue.Description = scanner.FormatVulnerabilityReport(scan.Image, filteredVulnerabilities)
						n = linearIssue
					default:
						log.Printf("Notifier: Error creating notifier: unknown notifier %s", integration.Type)
						continue
					}
					err = n.Notify(context.TODO())
					if err != nil {
						log.Printf("Notifier: Error sending notification: %v", err)
					} else {
						// Update scan status
						log.Println("to do")
						err = updateScanStatus(server, port, token, scan.ID, "success")
						if err != nil {
							log.Printf("Notifier: Error updating scan status: %v", err)
						}
					}
				}
			}
			if !matched {
				log.Printf("Notifier: Condition not matched for scanID: %s and image: %s. Updating status to success.", scan.ID, scan.Image)
				err = updateScanStatus(server, port, token, scan.ID, "success")
				if err != nil {
					log.Printf("Notifier: Failed to update status for scanID: %s: %v", scan.ID, err)
				}
				continue
			}
		}
	}
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

	var scans []scans.GetScanResponse
	if err := json.NewDecoder(resp.Body).Decode(&scans); err != nil {
		return nil, fmt.Errorf("failed to decode scans response: %w", err)
	}

	return scans, nil
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

func CheckAndSaveHash(server string, port int, token string, scanID uuid.UUID, hash string) (bool, error) {
	// Construct the API URL
	url := constructURL(server, port, "/api/v1/scannotify")

	// Prepare the request payload
	payload := map[string]interface{}{
		"scan_id": scanID,
		"hash":    hash,
	}
	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return false, fmt.Errorf("failed to marshal hash payload: %w", err)
	}

	// Create the API request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payloadJSON))
	if err != nil {
		return false, fmt.Errorf("failed to create request for CheckAndSaveHash: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	// Execute the request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return false, fmt.Errorf("failed to execute request for CheckAndSaveHash: %w", err)
	}
	defer resp.Body.Close()

	// Handle the response
	if resp.StatusCode == http.StatusConflict {
		// Hash already exists
		return false, nil
	} else if resp.StatusCode == http.StatusCreated || resp.StatusCode == http.StatusOK {
		// Hash successfully saved
		return true, nil
	} else {
		return false, fmt.Errorf("unexpected server response: %d", resp.StatusCode)
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
