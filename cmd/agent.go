package cmd

// Task Fetching:
// Execute Tasks
// Publish Results

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/shinobistack/gokakashi/ent/schema"
	"github.com/shinobistack/gokakashi/internal/helper"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/shinobistack/gokakashi/internal/http/client"
	"github.com/shinobistack/gokakashi/internal/restapi/v1/agents"
	"github.com/shinobistack/gokakashi/internal/restapi/v1/agenttasks"
	"github.com/shinobistack/gokakashi/internal/restapi/v1/integrations"
	"github.com/shinobistack/gokakashi/internal/restapi/v1/scans"
	"github.com/shinobistack/gokakashi/pkg/registry/v1"
	"github.com/shinobistack/gokakashi/pkg/scanner/v1"
	"github.com/spf13/cobra"
)

type httpClientKey struct{}

var agentCmd = &cobra.Command{
	Use:   "agent",
	Short: "Manage agents for GoKakashi",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if token == "" {
			log.Fatalf("Error: missing required flag --token")
		}

		if server == "" {
			log.Fatalf("Error: missing required flag --server")
		}

		headers := make(map[string]string)
		cfClientID := os.Getenv("CF_ACCESS_CLIENT_ID")
		cfClientSecret := os.Getenv("CF_ACCESS_CLIENT_SECRET")
		if cfClientID != "" && cfClientSecret != "" {
			headers["CF-Access-Client-Id"] = cfClientID
			headers["CF-Access-Client-Secret"] = cfClientSecret
		} else if cfClientSecret != "" {
			fmt.Println("Warning: ignoring CF_ACCESS_CLIENT_SECRET because CF_ACCESS_CLIENT_ID is not set")
		} else if cfClientID != "" {
			fmt.Println("Warning: ignoring CF_ACCESS_CLIENT_ID because CF_ACCESS_CLIENT_SECRET is not set")
		}

		httpClient := client.New(
			client.WithToken(token),
			client.WithHeaders(headers),
		)

		ctx := context.WithValue(context.Background(), httpClientKey{}, httpClient)
		cmd.SetContext(ctx)
	},
}

var agentStartCmd = &cobra.Command{
	Use:   "start",
	Short: "Register an agent and start polling for tasks",
	Run:   agentRegister,
}

var agentStopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Deregister an agent gracefully",
	Run:   agentDeRegister,
}

var (
	server       string
	token        string
	workspace    string
	name         string
	id           int
	chidori      bool
	labels       string
	singleStrike bool
)

func agentDeRegister(cmd *cobra.Command, args []string) {
	id, _ := cmd.Flags().GetInt("id")
	name, _ := cmd.Flags().GetString("name")
	chidori, _ := cmd.Flags().GetBool("chidori")

	if name == "" && id == 0 {
		log.Fatalf("Error: Either --name or --id must be provided")
	}

	queryParams := url.Values{}
	if id != 0 {
		queryParams.Add("id", fmt.Sprintf("%d", id))
	}
	if name != "" {
		queryParams.Add("name", name)
	}
	if chidori {
		queryParams.Add("chidori", "true")
	}

	url := fmt.Sprintf("%s/api/v1/agents?%s", server, queryParams.Encode())

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		log.Fatalf("Failed to create deregistration request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	httpClient := cmd.Context().Value(httpClientKey{}).(*client.Client)
	resp, err := httpClient.Do(req)
	if err != nil {
		log.Fatalf("Failed to deregister the agent: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		log.Fatalf("Failed to deregister the agent. Status: %d, Response: %s", resp.StatusCode, string(body))
	}

	var response agents.DeleteAgentResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		log.Fatalf("Failed to parse deregistration response: %v", err)
	}

	log.Printf("Agent successfully deregistered. ID: %d, Status: %s", response.ID, response.Status)
}

//ToDo: for any table status which results to error should we upload err message or just status error

func agentRegister(cmd *cobra.Command, args []string) {

	parsedLabels, err := parseLabels(labels)
	if err != nil {
		log.Fatalf("Error parsing labels: %v", err)
	}

	// Register the agent
	agent, err := registerAgent(cmd.Context(), server, token, workspace, name, parsedLabels)
	if err != nil {
		log.Fatalf("Failed to register the agent: %v", err)
	}

	log.Printf("Agent registered successfully! Agent ID: %d, Name: %s, Workspace: %s", agent.ID, agent.Name, agent.Workspace)

	if err := sendAgentHeartbeat(cmd.Context(), server, token, agent.ID); err != nil {
		log.Printf("Failed to send final heartbeat for agent %d: %v", agent.ID, err)
	}

	//// Start sending heartbeats in a separate goroutine, thinking how to handle stale goroutines then?
	//go startHeartbeat(cmd.Context(), server, token, agent.ID, 30*time.Second)

	// Update Agent status
	if err := updateAgentStatus(cmd.Context(), server, token, agent.ID, "scan_in_progress"); err != nil {
		log.Printf("Failed to update scan status to 'error': %v", err)
	}

	// Ephemeral agents
	if singleStrike {
		log.Printf("Ephemeral agent registered. Starting in sometime...")
		// Todo: Dynamic Delay: Instead of hardcoding 5 seconds, make it configurable
		time.Sleep(20 * time.Second)
		if err := sendAgentHeartbeat(cmd.Context(), server, token, agent.ID); err != nil {
			log.Printf("Failed to send final heartbeat for agent %d: %v", agent.ID, err)
		}
		executeEphemeraTasks(cmd.Context(), server, token, agent.ID, agent.Workspace)
	} else {
		// Long-lived agents keep polling
		pollTasks(cmd.Context(), server, token, agent.ID, agent.Workspace)
	}
}

//func startHeartbeat(ctx context.Context, server, token string, agentID int, interval time.Duration) {
//	ticker := time.NewTicker(interval)
//	defer ticker.Stop()
//
//	for {
//		select {
//		case <-ticker.C:
//			// Send heartbeat
//			if err := sendAgentHeartbeat(ctx, server, token, agentID); err != nil {
//				log.Printf("Failed to send heartbeat for agent %d: %v", agentID, err)
//			}
//		case <-ctx.Done():
//			log.Printf("Stopping heartbeat for agent %d", agentID)
//			return
//		}
//	}
//}

func sendAgentHeartbeat(ctx context.Context, server, token string, agentID int) error {
	path := fmt.Sprintf("/api/v1/agents/%d/heartbeat", agentID)
	url := helper.ConstructURL(server, path)

	req, err := http.NewRequest("PUT", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create heartbeat request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := ctx.Value(httpClientKey{}).(*client.Client).Do(req)
	if err != nil {
		return fmt.Errorf("failed to send heartbeat: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("server responded with status: %d, response: %s", resp.StatusCode, body)
	}

	log.Printf("Heartbeat successfully sent for agent ID %d", agentID)

	return nil
}

func parseLabels(labels string) ([]schema.CommonLabels, error) {
	if labels == "" {
		return nil, nil
	}

	// Split the labels string into key=value pairs
	pairs := strings.Split(labels, ",")
	parsedLabels := make([]schema.CommonLabels, len(pairs))

	for i, pair := range pairs {
		kv := strings.Split(pair, "=")
		if len(kv) != 2 || kv[0] == "" || kv[1] == "" {
			return nil, fmt.Errorf("invalid label format: %s (expected key=value)", pair)
		}
		parsedLabels[i] = schema.CommonLabels{
			Key:   strings.TrimSpace(kv[0]),
			Value: strings.TrimSpace(kv[1]),
		}
	}

	return parsedLabels, nil
}

func registerAgent(ctx context.Context, server, token, workspace, name string, labels []schema.CommonLabels) (agents.RegisterAgentResponse, error) {
	var response agents.RegisterAgentResponse
	reqBody := agents.RegisterAgentRequest{
		Server:    server,
		Token:     token,
		Workspace: workspace,
		Name:      name,
		Labels:    labels,
	}
	reqBodyJSON, _ := json.Marshal(reqBody)

	url := helper.ConstructURL(server, "/api/v1/agents")
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(reqBodyJSON))
	if err != nil {
		return response, fmt.Errorf("failed to create registration request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := ctx.Value(httpClientKey{}).(*client.Client).Do(req)
	if err != nil {
		return response, fmt.Errorf("failed to send registration request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return response, fmt.Errorf("server responded with status: %d", resp.StatusCode)
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return response, fmt.Errorf("failed to decode registration response: %w", err)
	}
	return response, nil
}

func updateAgentStatus(ctx context.Context, server, token string, agentID int, status string) error {
	reqBody := agents.UpdateAgentRequest{
		ID:     agentID,
		Status: status,
	}
	reqBodyJSON, _ := json.Marshal(reqBody)

	path := fmt.Sprintf("/api/v1/agents/%d", agentID)
	url := helper.ConstructURL(server, path)
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(reqBodyJSON))

	if err != nil {
		return fmt.Errorf("failed to create agent status update request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := ctx.Value(httpClientKey{}).(*client.Client).Do(req)
	if err != nil {
		return fmt.Errorf("failed to update agent status: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("server responded with status: %d", resp.StatusCode)
	}

	return nil
}

// Todo: Better way to handle polling for some time for ephemeral agent's task

func executeEphemeraTasks(ctx context.Context, server, token string, agentID int, workspace string) {
	// todo: retry logic for sometime

	if err := sendAgentHeartbeat(ctx, server, token, agentID); err != nil {
		log.Printf("Failed to send heartbeat for agent %d: %v", agentID, err)
	}

	tasks, err := fetchTasks(ctx, server, token, agentID, "pending", 1)
	if err != nil {
		log.Printf("Error fetching task: %v", err)
		return
	}

	if len(tasks) == 0 {
		log.Println("No pending tasks available for the ephemeral agent.")
		// Deregister the agent if no tasks are available
		deregisterEphemeralAgent(ctx, agentID, server, token, true)
		return
	}

	task := tasks[0]

	// Mark task as "in_progress"
	err = updateAgentTaskStatus(ctx, server, token, task.ID, agentID, "in_progress")
	if err != nil {
		log.Printf("Failed to update task status to 'in_progress': %v", err)
		return
	}

	// Process the task
	if err := sendAgentHeartbeat(ctx, server, token, agentID); err != nil {
		log.Printf("Failed to send heartbeat for agent %d: %v", agentID, err)
	}
	processTask(ctx, server, token, task, workspace, agentID)

	// Mark task as "complete"
	err = updateAgentTaskStatus(ctx, server, token, task.ID, agentID, "complete")
	if err != nil {
		log.Printf("Failed to update agent_task status to 'in_progress': %v", err)
		return
	}

	err = updateAgentStatus(ctx, server, token, agentID, "disconnected")
	if err != nil {
		log.Printf("Failed to update scan status to 'disconnected': %v", err)
	}

	if err := sendAgentHeartbeat(ctx, server, token, agentID); err != nil {
		log.Printf("Failed to send heartbeat for agent %d: %v", agentID, err)
	}

	log.Printf("Ephemeral agent %d completed its task and will now exit.", agentID)
	deregisterEphemeralAgent(ctx, agentID, server, token, true)
	log.Printf("Ephemeral agent %d shutting down.", agentID)
	os.Exit(0)

}

func deregisterEphemeralAgent(ctx context.Context, agentID int, server, token string, chidori bool) {
	cmd := &cobra.Command{}
	cmd.SetContext(ctx)

	cmd.Flags().Int("id", agentID, "Agent ID")
	cmd.Flags().String("server", server, "Server URL")
	cmd.Flags().String("token", token, "API Token")
	cmd.Flags().Bool("chidori", chidori, "Trigger hard delete")

	args := []string{}
	if agentID != 0 {
		args = append(args, fmt.Sprintf("--id=%d", agentID))
	}
	if chidori {
		args = append(args, "--chidori")
	}

	err := cmd.ParseFlags(args)
	if err != nil {
		log.Fatalf("Failed to parse flags: %v", err)
	}
	agentDeRegister(cmd, args)
}

func pollTasks(ctx context.Context, server, token string, agentID int, workspace string) {
	// Heartbeat frequency
	heartbeatInterval := 30 * time.Second
	lastHeartbeatSent := time.Now()

	for {

		now := time.Now()

		// Send heartbeat if due
		// alive as long as it's polling, so heartbeat is sent when loop restarts.
		if now.Sub(lastHeartbeatSent) >= heartbeatInterval {
			if err := sendAgentHeartbeat(ctx, server, token, agentID); err != nil {
				log.Printf("Failed to send heartbeat for agent %d: %v", agentID, err)
			} else {
				lastHeartbeatSent = now
			}
		}

		// Process only tasks with status "pending" in the order returned (created_at ASC)
		tasks, err := fetchTasks(ctx, server, token, agentID, "pending", 0)
		if err != nil {
			log.Printf("Error fetching tasks: %v", err)
			time.Sleep(10 * time.Second)
			continue
		}

		if len(tasks) == 0 {
			log.Println("No pending tasks. Retrying after 10 seconds.")
			time.Sleep(10 * time.Second)
			continue
		}

		for _, task := range tasks {
			// Update task status to "in_progress"
			err := updateAgentTaskStatus(ctx, server, token, task.ID, agentID, "in_progress")
			if err != nil {
				log.Printf("Failed to update agent_task status to 'in_progress': %v", err)
				return
			}
			if err := sendAgentHeartbeat(ctx, server, token, agentID); err != nil {
				log.Printf("Failed to send heartbeat for agent %d: %v", agentID, err)
			} else {
				lastHeartbeatSent = time.Now()
			}
			processTask(ctx, server, token, task, workspace, agentID)
			continue
		}
		// Todo: Polling interval time decide
		// Sleep for a defined interval
		time.Sleep(10 * time.Second)
	}
}

func fetchTasks(ctx context.Context, server, token string, agentID int, status string, limit int) ([]agenttasks.GetAgentTaskResponse, error) {
	url := helper.ConstructURL(server, fmt.Sprintf("/api/v1/agents/%d/tasks", agentID)) + fmt.Sprintf("?status=%s&limit=%d", status, limit)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create task polling request: %w", err)
	}

	resp, err := ctx.Value(httpClientKey{}).(*client.Client).Do(req)
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

func updateAgentTaskStatus(ctx context.Context, server, token string, taskID uuid.UUID, agentID int, status string) error {
	reqBody := agenttasks.UpdateAgentTaskRequest{
		ID:      taskID,
		AgentID: intPtr(agentID),
		Status:  strPtr(status),
	}

	reqBodyJSON, _ := json.Marshal(reqBody)

	path := fmt.Sprintf("/api/v1/agents/%d/tasks/%s", agentID, taskID)
	url := helper.ConstructURL(server, path)
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(reqBodyJSON))

	if err != nil {
		return fmt.Errorf("failed to create task status update request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := ctx.Value(httpClientKey{}).(*client.Client).Do(req)
	if err != nil {
		return fmt.Errorf("failed to update task status: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("server responded with status: %d", resp.StatusCode)
	}

	return nil
}

func processTask(ctx context.Context, server, token string, task agenttasks.GetAgentTaskResponse, workspace string, agentID int) {
	// Step 1: Fetch scan details
	scan, err := fetchScan(ctx, server, token, task.ScanID)
	if err != nil {
		log.Printf("Failed to fetch scan details: %v", err)
		return
	}

	// if scan.IntegrationID not empty then execute fetchIntegration and authenticateAndPullImage
	if scan.IntegrationID != nil {
		// Step 2: Fetch integration details
		integration, err := fetchIntegration(ctx, server, token, scan.IntegrationID)
		if err != nil {
			log.Printf("Failed to fetch integration details: %v", err)
			//return
		}

		// Step 3: Authenticate and pull the image
		if err := authenticateAndPullImage(scan.Image, integration); err != nil {
			log.Printf("Failed to authenticate or pull image: %v", err)
			//return
		}
	}

	err = updateScanStatus(ctx, server, token, scan.ID, "scan_in_progress")
	if err != nil {
		log.Printf("Failed to update scan status to 'scan_in_progress': %v", err)
	}
	// Step 4: Perform the scan
	// severityLevels := []string{"HIGH", "CRITICAL"}
	reportPath, err := performScan(scan.Image, scan.Scanner)
	if err != nil {
		log.Printf("Failed to perform scan: %v", err)
		if err := updateScanStatus(ctx, server, token, scan.ID, "error"); err != nil {
			log.Printf("Failed to update scan status to 'error': %v", err)
		}
		return
	}

	// Step 5: Upload the scan report
	if err := uploadReport(ctx, server, token, scan.ID, reportPath); err != nil {
		log.Printf("Failed to upload scan report: %v", err)
		if err := updateScanStatus(ctx, server, token, scan.ID, "error"); err != nil {
			log.Printf("Failed to update scan status to 'error': %v", err)
		}
		return
	}

	// step 6: Verify scans.Notify field exist
	// Todo: if exists update the status to notify_pending else complete
	if scan.Notify == nil {
		log.Printf("No notify specified for scan ID: %s", scan.ID)
		if err := updateScanStatus(ctx, server, token, scan.ID, "success"); err != nil {
			log.Printf("Failed to update scan status to 'success': %v", err)
		}
	} else {
		err = updateScanStatus(ctx, server, token, scan.ID, "notify_pending")
		if err != nil {
			log.Printf("Failed to update scan status to 'scan_in_progress': %v", err)
		}
	}

	if err := updateAgentTaskStatus(ctx, server, token, task.ID, agentID, "complete"); err != nil {
		log.Printf("Failed to update agent_task status to 'complete': %v", err)
	}

	log.Printf("AgentTaskID completed successfully: %v", task.ID)
}

func updateScanStatus(ctx context.Context, server, token string, scanID uuid.UUID, status string) error {
	reqBody := scans.UpdateScanRequest{
		ID:     scanID,
		Status: strPtr(status),
	}
	reqBodyJSON, _ := json.Marshal(reqBody)

	path := fmt.Sprintf("/api/v1/scans/%s", scanID)
	url := helper.ConstructURL(server, path)
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(reqBodyJSON))

	if err != nil {
		return fmt.Errorf("failed to create scan status update request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := ctx.Value(httpClientKey{}).(*client.Client).Do(req)
	if err != nil {
		return fmt.Errorf("failed to update scan status: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("server responded with status: %d", resp.StatusCode)
	}

	return nil
}

func fetchScan(ctx context.Context, server, token string, scanID uuid.UUID) (*scans.GetScanResponse, error) {
	path := fmt.Sprintf("/api/v1/scans/%s", scanID)
	url := helper.ConstructURL(server, path)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create scan request: %w", err)
	}

	resp, err := ctx.Value(httpClientKey{}).(*client.Client).Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch scan details: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("server responded with status: %d", resp.StatusCode)
	}

	var scan scans.GetScanResponse
	if err := json.NewDecoder(resp.Body).Decode(&scan); err != nil {
		return nil, fmt.Errorf("failed to decode scan response: %w", err)
	}

	return &scan, nil
}

func fetchIntegration(ctx context.Context, server, token string, integrationID *uuid.UUID) (*integrations.GetIntegrationResponse, error) {
	path := fmt.Sprintf("/api/v1/integrations/%s", integrationID)
	url := helper.ConstructURL(server, path)
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return nil, fmt.Errorf("failed to create integration fetch request: %w", err)
	}

	resp, err := ctx.Value(httpClientKey{}).(*client.Client).Do(req)
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

func authenticateAndPullImage(image string, integration *integrations.GetIntegrationResponse) error {
	// Select the registry integration based on type
	log.Printf("Registry: %s | Image: %s | Step: Authentication started", integration.Type, image)
	reg, err := registry.NewRegistry(integration.Type, integration.Config)
	if err != nil {
		return err
	}

	// Authenticate
	if err := reg.Authenticate(); err != nil {
		return fmt.Errorf("authentication failed: %w", err)
	}

	// Pull image
	if err := reg.PullImage(image); err != nil {
		return fmt.Errorf("failed to pull image: %w", err)
	}

	return nil
}

func performScan(image, scannerType string) (string, error) {
	// Initialize the scanner using the factory function.
	scanner, err := scanner.NewScanner(scannerType)
	if err != nil {
		return "", fmt.Errorf("failed to initialize scanner: %w", err)
	}

	// Perform scan
	// Todo: to add feature for severity args or tool args.
	reportPath, err := scanner.Scan(image, nil)
	if err != nil {
		return "", fmt.Errorf("scan failed: %w", err)
	}

	return reportPath, nil
}

func uploadReport(ctx context.Context, server, token string, scanID uuid.UUID, reportPath string) error {
	report, err := os.ReadFile(reportPath)
	if err != nil {
		return fmt.Errorf("failed to read report file: %w", err)
	}

	reqBody := scans.UpdateScanRequest{
		ID:     scanID,
		Report: json.RawMessage(report),
	}
	reqBodyJSON, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("failed to marshal request body: %w", err)
	}

	path := fmt.Sprintf("/api/v1/scans/%s", scanID)
	url := helper.ConstructURL(server, path)
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(reqBodyJSON))

	if err != nil {
		return fmt.Errorf("failed to create report upload request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := ctx.Value(httpClientKey{}).(*client.Client).Do(req)
	if err != nil {
		return fmt.Errorf("failed to upload scan report: %w", err)
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

func intPtr(i int) *int {
	return &i
}
