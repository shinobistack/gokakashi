package assigner

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/shinobistack/gokakashi/ent/schema"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/shinobistack/gokakashi/internal/restapi/v1/agents"
	"github.com/shinobistack/gokakashi/internal/restapi/v1/agenttasks"
	"github.com/shinobistack/gokakashi/internal/restapi/v1/scans"
)

// Assigns scanID to available Agents

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
	log.Println("Starting the periodic task assigner...")
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	// Todo: Maybe we can do batch processing
	for range ticker.C {
		AssignTasks(server, port, token)
	}
}

// Global agent index for round-robin
var globalAgentIndex int

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
		log.Printf("Unassignable scans: %d scans pending without agents.", len(pendingScans))
		return
	}

	// log.Printf("Agents are available: %v", availableAgents)

	// Step 3: Assign scans to agents
	// ToDo: to explore task assignment for better efficiency
	for _, scan := range pendingScans {
		// Check if scan is already assigned
		if isScanAssigned(server, port, token, scan.ID) {
			log.Printf("Scan ID %s is already assigned. Skipping.", scan.ID)
			continue
		}

		// Step 3a: Filter agents by matching labels
		matchingAgents := filterAgentsByLabels(availableAgents, scan.Labels)

		// Step 3b: Assign using round-robin
		var agent agents.GetAgentResponse
		if len(matchingAgents) > 0 {
			agent = selectAgentRoundRobin(matchingAgents)
		} else {
			// Fallback to any connected agent
			agent = selectAgentRoundRobin(availableAgents)
		}

		// Assign the scan to the selected agent
		if assignTaskToAgent(server, port, token, agent, scan) {
			log.Printf("Successfully assigned scan %s to agent %d", scan.ID, agent.ID)
		} else {
			log.Printf("Failed to assign scan %s. It will be retried in the next cycle.", scan.ID)
		}

		//// Select agent using round-robin
		//agent := availableAgents[i%len(availableAgents)]
		//if err := createAgentTask(server, port, token, agent.ID, scan.ID); err != nil {
		//	log.Printf("Failed to assign scan %s to agent %d: %v", scan.ID, agent.ID, err)
		//} else {
		//	log.Printf("Successfully assigned scan %s to agent %d", scan.ID, agent.ID)
		//}

	}
}

// Selects the next agent in a round-robin fashion
func selectAgentRoundRobin(agents []agents.GetAgentResponse) agents.GetAgentResponse {
	// Select the agent at the current global index
	agent := agents[globalAgentIndex%len(agents)]

	// Update the index for the next assignment
	globalAgentIndex = (globalAgentIndex + 1) % len(agents)
	return agent
}

func filterAgentsByLabels(agentList []agents.GetAgentResponse, scanLabels []schema.CommonLabels) []agents.GetAgentResponse {
	var matchingAgents []agents.GetAgentResponse
	for _, agent := range agentList {
		if labelsMatch(agent.Labels, scanLabels) {
			matchingAgents = append(matchingAgents, agent)
		}
	}
	return matchingAgents
}

// Matches atleast one label? Todo: Maybe having maritial labels would benefit? Because a scan have many labels for filtering?
//func partialLabelsMatch(agentLabels, scanLabels []schema.CommonLabels) bool {
//	for _, scanLabel := range scanLabels {
//		for _, agentLabel := range agentLabels {
//			if scanLabel.Key == agentLabel.Key && scanLabel.Value == agentLabel.Value {
//				return true // Return true as soon as one label matches
//			}
//		}
//	}
//	return false // Return false if no labels match
//}

// Matches all labels?
func labelsMatch(agentLabels, scanLabels []schema.CommonLabels) bool {
	for _, scanLabel := range scanLabels {
		matchFound := false
		for _, agentLabel := range agentLabels {
			if scanLabel.Key == agentLabel.Key && scanLabel.Value == agentLabel.Value {
				matchFound = true
				break
			}
		}
		if !matchFound {
			return false // If any scan label doesn’t match, return false
		}
	}
	return true
}

func assignTaskToAgent(server string, port int, token string, agent agents.GetAgentResponse, scan scans.GetScanResponse) bool {
	if err := createAgentTask(server, port, token, agent.ID, scan.ID); err != nil {
		log.Printf("Failed to assign scan %s to agent %d: %v", scan.ID, agent.ID, err)
		return false
	}
	log.Printf("Successfully assigned scan %s to agent %d", scan.ID, agent.ID)
	return true
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

	// Check if the status is not OK
	if resp.StatusCode != http.StatusOK {
		log.Printf("Unexpected status code %d when checking scan assignment", resp.StatusCode)
		return false
	}

	// Parse the response body
	var tasks []agenttasks.GetAgentTaskResponse
	if err := json.NewDecoder(resp.Body).Decode(&tasks); err != nil {
		log.Printf("Failed to parse response body: %v", err)
		return false
	}

	// If the tasks list is empty, the scan is not assigned
	return len(tasks) > 0
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

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("server responded with status: %d", resp.StatusCode)
	}

	return nil
}
