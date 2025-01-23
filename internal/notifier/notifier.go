package notifier

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/shinobistack/gokakashi/internal/parser"
	"github.com/shinobistack/gokakashi/internal/restapi/v1/integrations"
	"github.com/shinobistack/gokakashi/internal/restapi/v1/scans"
	"github.com/shinobistack/gokakashi/pkg/notifier/v1"
)

type HashEntry struct {
	Image           string              `json:"image"`           // Image name
	Vulnerabilities []VulnerabilityData `json:"vulnerabilities"` // Detailed vulnerability data
	Hash            string              `json:"hash"`            // Generated hash for the entry
}

type VulnerabilityData struct {
	VulnerabilityID  string `json:"vulnerability_id"`  // The CVE or vulnerability ID
	Severity         string `json:"severity"`          // Severity level (e.g., Critical, High)
	InstalledVersion string `json:"installed_version"` // Version of the package installed
	FixedVersion     string `json:"fixed_version"`     // Version where the vulnerability is fixed (if available)
}

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
		log.Printf("Error fetching pending notify: %v", err)
	}

	if len(scans) == 0 {
		log.Println("No pending notify to execute.")
		return
	}

	for _, scan := range scans {
		for _, notify := range *scan.Notify {
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
					log.Printf("invalid UUID string: %v", err)
				}
				log.Println(parsedNotifyToUUID)

				integration, err := fetchIntegrationDetails(server, port, token, parsedNotifyToUUID)
				if err != nil {
					log.Printf("Error fetching integration details: %v", err)
					continue
				}

				notifier, err := notifier.NewNotifier(integration.Type)
				if err != nil {
					log.Printf("Error creating notifier: %v", err)
					continue
				}

				// Todo: Better way to deal handle.
				notifierConfig, err := MapToNotificationConfig(integration.Config)
				if err != nil {
					log.Printf("Error parsing notification config: %v", err)
					continue
				}

				filteredVulnerabilities, err := formatReportForNotify(scan.Report, severities, scan.Image)
				if err != nil {
					log.Printf("Error formatting report for notify: %v", err)
					continue
				}

				if len(filteredVulnerabilities) == 0 {
					log.Printf("no vulnerabilities found for scanID: %s and image: %s. Skipping creation of issues", scan.ID, scan.Image)
					err = updateScanStatus(server, port, token, scan.ID, "success")
					if err != nil {
						log.Printf("Error updating scan status: %v", err)
					}
					return
				}

				vulnerabilityEntries := ConvertVulnerabilities(filteredVulnerabilities)

				// Generate a hash and check/save
				hash := GenerateHash(scan.Image, vulnerabilityEntries)
				saved, err := CheckAndSaveHash(server, port, token, scan.ID, hash)
				if err != nil {
					log.Printf("Error checking or saving hash: %v", err)
					continue
				}

				if saved {
					err = notifier.CreateIssue(scan.Image, filteredVulnerabilities, notifierConfig)
					if err != nil {
						log.Printf("Error sending notification: %v", err)
					} else {
						// Update scan status
						log.Println("to do")
						err = updateScanStatus(server, port, token, scan.ID, "success")
						if err != nil {
							log.Printf("Error updating scan status: %v", err)
						}
					}
				}
			}
			if !matched {
				log.Printf("Condition not matched for scanID: %s and image: %s. Updating status to success.", scan.ID, scan.Image)
				err = updateScanStatus(server, port, token, scan.ID, "success")
				if err != nil {
					log.Printf("Failed to update status for scanID: %s: %v", scan.ID, err)
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

func formatReportForNotify(scanReport json.RawMessage, severities []string, scanImage string) ([]notifier.Vulnerability, error) {
	var report notifier.Report
	err := json.Unmarshal(scanReport, &report)
	if err != nil {
		log.Printf("Error failed to parse scan report: %v", err)
	}

	// Load the vulnerabilities from scans.report
	var vulnerabilities []notifier.Vulnerability
	for _, result := range report.Results {
		vulnerabilities = append(vulnerabilities, result.Vulnerabilities...)
	}

	filteredVulnerabilities := FilterVulnerabilitiesBySeverity(vulnerabilities, severities)

	return filteredVulnerabilities, nil
}

func FilterVulnerabilitiesBySeverity(vulnerabilities []notifier.Vulnerability, severityLevels []string) []notifier.Vulnerability {
	var filtered []notifier.Vulnerability
	for _, v := range vulnerabilities {
		for _, level := range severityLevels {
			if v.Severity == level {
				filtered = append(filtered, v)
			}
		}
	}
	return filtered
}

func ConvertVulnerabilities(filteredVulnerabilities []notifier.Vulnerability) []string {
	var vulnerabilityEntries []string
	for _, v := range filteredVulnerabilities {
		data := VulnerabilityData{
			VulnerabilityID:  v.VulnerabilityID,
			Severity:         v.Severity,
			InstalledVersion: v.InstalledVersion,
			FixedVersion:     v.FixedVersion,
		}
		entry := fmt.Sprintf("%s_%s_%s_%s", data.VulnerabilityID, data.Severity, data.InstalledVersion, data.FixedVersion)
		vulnerabilityEntries = append(vulnerabilityEntries, entry)
	}
	return vulnerabilityEntries
}

func GenerateHash(image string, vulnerabilities []string) string {
	data := fmt.Sprintf("%s_%s", image, strings.Join(vulnerabilities, "_"))
	hash := sha256.New()
	hash.Write([]byte(data))
	return hex.EncodeToString(hash.Sum(nil))
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

func MapToNotificationConfig(config map[string]interface{}) (notifier.NotificationConfig, error) {
	// Ensure all fields are present and valid
	apiKey, ok := config["api_key"].(string)
	if !ok {
		return notifier.NotificationConfig{}, fmt.Errorf("missing or invalid api_key")
	}

	projectID, _ := config["project_id"].(string)
	teamID, _ := config["team_id"].(string)
	title, _ := config["issue_title"].(string)
	priority, _ := config["issue_priority"].(int)
	assignee, _ := config["issue_assignee_id"].(string)
	stateID, _ := config["issue_state_id"].(string)
	dueDate, _ := config["issue_due_date"].(string)

	return notifier.NotificationConfig{
		APIKey:    apiKey,
		ProjectID: projectID,
		TeamID:    teamID,
		Title:     title,
		Priority:  priority,
		Assignee:  assignee,
		StateID:   stateID,
		DueDate:   dueDate,
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
