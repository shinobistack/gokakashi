package health

import (
	"encoding/json"
	"fmt"
	"github.com/shinobistack/gokakashi/internal/restapi/v1/agents"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// Todo: Put the common helper function like normaliseserver and constructURL into single util.

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

func monitorAgentHeartbeats(server string, port int, token string, interval time.Duration) {
	ticker := time.NewTicker(61 * time.Second) // Run every 61 seconds
	defer ticker.Stop()

	for range ticker.C {
		now := time.Now()

		// Fetch all connected agents
		agentList, err := fetchConnectedAgents(server, port, token, "scan_in_progress")
		if err != nil {
			log.Printf("Failed to fetch connected agents: %v", err)
			continue
		}

		// Check heartbeats and update status if stale
		for _, agent := range agentList {
			if now.Sub(agent.LastHeartbeat) > interval {
				log.Printf("Agent %d is stale (last heartbeat: %v). Marking as disconnected.", agent.ID, agent.LastHeartbeat)
				err := updateSoftDeleteAgent(server, port, token, agent.ID)
				if err != nil {
					log.Printf("Failed to mark agent %d as disconnected: %v", agent.ID, err)
				}
			}
		}
	}
}

func fetchConnectedAgents(server string, port int, token, status string) ([]agents.GetAgentResponse, error) {
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

func updateSoftDeleteAgent(server string, port int, token string, agentID int) error {
	url := constructURL(server, port, fmt.Sprintf("/api/v1/agents/%d", agentID))

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("server responded with status: %d", resp.StatusCode)
	}

	return nil
}
