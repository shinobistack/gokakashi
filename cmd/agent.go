package cmd

// Task Fetching:
// Execute Tasks
// Publish Results

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/shinobistack/gokakashi/internal/restapi/v1/agents"
	"github.com/shinobistack/gokakashi/internal/restapi/v1/agenttasks"
	"github.com/spf13/cobra"
	"log"
	"net/http"
	"time"
)

var agentCmd = &cobra.Command{
	Use:   "agent",
	Short: "Manage agents for GoKakashi",
}

var agentStartCmd = &cobra.Command{
	Use:   "start",
	Short: "Register an agent and start polling for tasks",
	Run:   agentRegister,
}

var (
	server    string
	token     string
	workspace string
)

func agentRegister(cmd *cobra.Command, args []string) {
	// Validate inputs
	if server == "" || token == "" || workspace == "" {
		log.Fatalf("Error: Missing required inputs. Please provide --server, --token, and --workspace.")
	}

	// log.Printf("Server: %s, Token: %s, Workspace: %s", server, token, workspace)

	// Register the agent
	agentID, err := registerAgent(server, token, workspace)
	if err != nil {
		log.Fatalf("Failed to register the agent: %v", err)
	}

	log.Printf("Agent registered successfully! Agent ID: %d", agentID)

	// Start polling for tasks
	pollTasks(server, token, agentID, workspace)
}

func registerAgent(server, token, workspace string) (int, error) {
	reqBody := agents.RegisterAgentRequest{
		Server:    server,
		Token:     token,
		Workspace: workspace,
	}
	reqBodyJSON, _ := json.Marshal(reqBody)

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/v1/agents", server), bytes.NewBuffer(reqBodyJSON))
	if err != nil {
		return 0, fmt.Errorf("failed to create registration request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, fmt.Errorf("failed to send registration request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("server responded with status: %d", resp.StatusCode)
	}

	var response agents.RegisterAgentResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return 0, fmt.Errorf("failed to decode registration response: %w", err)
	}
	return response.ID, nil
}

func pollTasks(server, token string, agentID int, workspace string) {
	log.Printf("WORKSPCAE: %v", workspace)
	for {
		tasks, err := fetchTasks(server, token, agentID)
		log.Println(tasks)
		if err != nil {
			log.Printf("Error fetching tasks: %v", err)
			time.Sleep(10 * time.Second)
			continue
		}

		for _, task := range tasks {
			log.Printf("Executing task: %v", task)
			log.Printf("TASK execution logic can be added here: %v", task)
			// Task execution logic
		}
		// Todo: Polling interval time decide
		// Sleep for a defined interval
		time.Sleep(10 * time.Second)
	}
}

func fetchTasks(server, token string, agentID int) ([]agenttasks.GetAgentTaskResponse, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/v1/agents/%d/tasks", server, agentID), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create task polling request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send task polling request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("server responded with status: %d", resp.StatusCode)
	}

	var tasks []agenttasks.GetAgentTaskResponse
	if err := json.NewDecoder(resp.Body).Decode(&tasks); err != nil {
		return nil, fmt.Errorf("failed to decode task polling response: %w", err)
	}

	return tasks, nil
}
