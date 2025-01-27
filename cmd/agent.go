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
}

var agentStartCmd = &cobra.Command{
	Use:   "start",
	Short: "Register an agent and start polling for tasks",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if token == "" {
			log.Fatalf("Error: missing required flag --token")
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
	Run: agentRegister,
}

var agentStopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Deregister an agent gracefully",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if token == "" {
			log.Fatalf("Error: missing required flag --token")
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
	Run: agentDeRegister,
}

var (
	server    string
	token     string
	workspace string
	name      string
	id        int
	chidori   bool
	labels    string
)

func agentDeRegister(cmd *cobra.Command, args []string) {
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
	// Validate inputs
	if server == "" || token == "" {
		log.Fatalf("Error: Missing required inputs. Please provide --server, --token.")
	}

	parsedLabels, err := parseLabels(labels)
	if err != nil {
		log.Fatalf("Error parsing labels: %v", err)
	}

	// Register the agent
	agentID, err := registerAgent(cmd.Context(), server, token, workspace, name, parsedLabels)
	if err != nil {
		log.Fatalf("Failed to register the agent: %v", err)
	}

	// ToDo: to display auto assigned name
	log.Printf("Agent registered successfully! Agent ID: %d, Name: %s, Workspace: %s", agentID, name, workspace)

	// Start polling for tasks
	pollTasks(cmd.Context(), server, token, agentID, workspace)
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

func registerAgent(ctx context.Context, server, token, workspace, name string, labels []schema.CommonLabels) (int, error) {
	reqBody := agents.RegisterAgentRequest{
		Server:    server,
		Token:     token,
		Workspace: workspace,
		Name:      name,
		Labels:    labels,
	}
	reqBodyJSON, _ := json.Marshal(reqBody)

	url := constructURL(server, "/api/v1/agents")
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(reqBodyJSON))
	if err != nil {
		return 0, fmt.Errorf("failed to create registration request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := ctx.Value(httpClientKey{}).(*client.Client).Do(req)
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

func pollTasks(ctx context.Context, server, token string, agentID int, workspace string) {
	for {
		// Process only tasks with status "pending" in the order returned (created_at ASC)
		tasks, err := fetchTasks(ctx, server, token, agentID, "pending")
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
			processTask(ctx, server, token, task, workspace, agentID)
			continue
		}
		// Todo: Polling interval time decide
		// Sleep for a defined interval
		time.Sleep(10 * time.Second)
	}
}

func fetchTasks(ctx context.Context, server, token string, agentID int, status string) ([]agenttasks.GetAgentTaskResponse, error) {
	url := constructURL(server, fmt.Sprintf("/api/v1/agents/%d/tasks", agentID)) + fmt.Sprintf("?status=%s", status)

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
	url := constructURL(server, path)
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

	// Step 2: Fetch integration details
	integration, err := fetchIntegration(ctx, server, token, scan.IntegrationID)
	if err != nil {
		log.Printf("Failed to fetch integration details: %v", err)
		return
	}

	// Step 3: Authenticate and pull the image
	if err := authenticateAndPullImage(scan.Image, integration); err != nil {
		log.Printf("Failed to authenticate or pull image: %v", err)
		return
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
	if scan.Notify == nil || len(*scan.Notify) == 0 {
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
	url := constructURL(server, path)
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
	url := constructURL(server, path)
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

func fetchIntegration(ctx context.Context, server, token string, integrationID uuid.UUID) (*integrations.GetIntegrationResponse, error) {
	path := fmt.Sprintf("/api/v1/integrations/%s", integrationID)
	url := constructURL(server, path)
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
	url := constructURL(server, path)
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
