package cmd

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/shinobistack/gokakashi/internal/http/client"
	"github.com/shinobistack/gokakashi/internal/restapi/v1/integrations"
	"github.com/shinobistack/gokakashi/internal/restapi/v1/policies"
	"github.com/shinobistack/gokakashi/internal/restapi/v1/scans"
	"github.com/spf13/cobra"
)

var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Scan a container image",
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

var scanImageCmd = &cobra.Command{
	Use:   "image",
	Short: "Trigger a scan for a container image",
	Run:   scanImage,
}

var scanStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Get the status of a scan",
	Run:   getScanStatus,
}

var (
	image      string
	policyName string
	scanID     string
)

func normalizeServer(server string) string {
	if !strings.HasPrefix(server, "http://") && !strings.HasPrefix(server, "https://") {
		server = "http://" + server
	}
	return server
}

func constructURL(server string, path string) string {
	base := normalizeServer(server)
	u, err := url.Parse(base)
	if err != nil {
		log.Fatalf("Invalid server URL: %s", base)
	}
	u.Path = path
	return u.String()
}

func scanImage(cmd *cobra.Command, args []string) {
	ctx := cmd.Context()
	httpClient := ctx.Value(httpClientKey{}).(*client.Client)

	if policyName == "" {
		log.Fatalf("Error: --policy is required to specify the policy name.")
	}
	if image == "" {
		log.Fatalf("Error: --image is required to specify the container image.")
	}

	parsedLabels, err := parseLabels(labels)
	if err != nil {
		log.Fatalf("Error parsing labels: %v", err)
	}

	policy, err := fetchPolicyByName(ctx, httpClient, policyName)
	if err != nil {
		log.Fatalf("Failed to fetch policy: %v", err)
	}
	log.Println(policy)

	integration, err := fetchIntegrationByName(ctx, httpClient, policy.Image.Registry)
	if err != nil {
		log.Fatalf("Failed to fetch integration: %v", err)
	}
	log.Println(integration)

	reqBody := map[string]interface{}{
		"policy_id":      policy.ID,
		"image":          image,
		"scanner":        policy.Scanner,
		"integration_id": integration.ID,
		"notify":         policy.Notify,
		"status":         "scan_pending",
		"labels":         parsedLabels,
	}
	reqBodyJSON, _ := json.Marshal(reqBody)
	url := constructURL(server, "/api/v1/scans")

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(reqBodyJSON))
	if err != nil {
		log.Fatalf("Failed to create scan request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := httpClient.Do(req)
	if err != nil {
		log.Fatalf("Failed to send scan request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Server responded with status: %d", resp.StatusCode)
	}

	var response scans.CreateScanResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		log.Fatalf("Failed to decode scan response: %v", err)
	}

	log.Printf("Scan triggered successfully! Scan ID: %s, status: %s", response.ID, response.Status)
}

func getScanStatus(cmd *cobra.Command, args []string) {
	ctx := cmd.Context()
	httpClient := ctx.Value(httpClientKey{}).(*client.Client)

	if server == "" || token == "" || scanID == "" {
		log.Fatalf("Error: Missing required inputs. Please provide --server, --token, and --scanID.")
	}

	url := constructURL(server, fmt.Sprintf("/api/v1/scans/%s", scanID))

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatalf("Failed to create status request: %v", err)
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		log.Fatalf("Failed to send status request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Server responded with status: %d", resp.StatusCode)
	}

	var response struct {
		Status string `json:"status"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		log.Fatalf("Failed to decode status response: %v", err)
	}

	log.Printf("Scan status: %s", response.Status)
}

func fetchPolicyByName(ctx context.Context, httpClient *client.Client, policyName string) (*policies.GetPolicyResponse, error) {
	url := constructURL(server, "/api/v1/policies") + fmt.Sprintf("?name=%s", policyName)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create policy fetch request: %w", err)
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch policy: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("server responded with status: %d", resp.StatusCode)
	}

	var policy []policies.GetPolicyResponse
	if err := json.NewDecoder(resp.Body).Decode(&policy); err != nil {
		return nil, fmt.Errorf("failed to decode policy response: %w", err)
	}

	return &policy[0], nil
}

func fetchIntegrationByName(ctx context.Context, httpClient *client.Client, integrationName string) (*integrations.GetIntegrationResponse, error) {
	url := constructURL(server, "/api/v1/integrations") + fmt.Sprintf("?name=%s", integrationName)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create integration fetch request: %w", err)
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch integration: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("server responded with status: %d", resp.StatusCode)
	}

	var integration []integrations.GetIntegrationResponse
	if err := json.NewDecoder(resp.Body).Decode(&integration); err != nil {
		return nil, fmt.Errorf("failed to decode integration response: %w", err)
	}

	return &integration[0], nil
}
