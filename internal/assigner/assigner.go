package assigner

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/shinobistack/gokakashi/internal/restapi/v1/agents"
	"github.com/shinobistack/gokakashi/internal/restapi/v1/agenttasks"
	"github.com/shinobistack/gokakashi/internal/restapi/v1/scans"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
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

func StartAssigner(server string, port int, token string, interval time.Duration) {
	log.Println("Starting the periodic task assigner...")
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for range ticker.C {
		AssignTasks(server, port, token)
	}

}

func AssignTasks(server string, port int, token string) {
	log.Println("Assigner now begins assigning your scans")
	// Step 1: Fetch scans needing assignment
	pendingScans, err := fetchPendingScans(server, port, token, "scan_pending")
	if err != nil {
		log.Printf("Error fetching pending scans: %v", err)
		return
	}

	if len(pendingScans) == 0 {
		log.Println("No pending scans to assign.")
		return
	}

	// Step 2: Fetch available agents
	availableAgents, err := fetchAvailableAgents(server, port, token, "connected")
	if err != nil {
		log.Printf("Error fetching available agents: %v", err)
		return
	}

	if len(availableAgents) == 0 {
		log.Println("No agents available for assignment.")
		return
	}

	// log.Printf("Agents are available: %v", availableAgents)

	// Step 3: Assign scans to agents
	// ToDo: to explore task assignment for better efficiency
	for i, scan := range pendingScans {
		// Check if scan is already assigned
		if isScanAssigned(server, port, token, scan.ID) {
			log.Printf("Scan ID %s is already assigned. Skipping.", scan.ID)
			continue
		}

		// Select agent using round-robin
		agent := availableAgents[i%len(availableAgents)]
		if err := createAgentTask(server, port, token, agent.ID, scan.ID); err != nil {
			log.Printf("Failed to assign scan %s to agent %d: %v", scan.ID, agent.ID, err)
		} else {
			log.Printf("Successfully assigned scan %s to agent %d", scan.ID, agent.ID)
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

func fetchAvailableAgents(server string, port int, token, status string) ([]agents.GetAgentResponse, error) {
	url := constructURL(server, port, "/api/v1/agents") + fmt.Sprintf("?status=%s", status)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("server responded with status: %d", resp.StatusCode)
	}

	var agents []agents.GetAgentResponse
	if err := json.NewDecoder(resp.Body).Decode(&agents); err != nil {
		return nil, err
	}

	return agents, nil
}

func isScanAssigned(server string, port int, token string, scanID uuid.UUID) bool {
	url := constructURL(server, port, "/api/v1/agents/tasks") + fmt.Sprintf("?scan_id=%s", scanID)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Printf("Error checking scan assignment: %v", err)
		return false
	}

	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("Error checking scan assignment: %v", err)
		return false
	}
	defer resp.Body.Close()

	return resp.StatusCode == http.StatusOK
}

func createAgentTask(server string, port int, token string, agentID int, scanID uuid.UUID) error {
	url := constructURL(server, port, fmt.Sprintf("/api/v1/agents/%d/tasks", agentID))

	reqBody := agenttasks.CreateAgentTaskRequest{
		AgentID:   agentID,
		ScanID:    scanID,
		Status:    "pending",
		CreatedAt: time.Now(),
	}

	reqBodyJSON, _ := json.Marshal(reqBody)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(reqBodyJSON))
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("server responded with status: %d", resp.StatusCode)
	}

	return nil
}
