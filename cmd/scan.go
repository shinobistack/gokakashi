package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/shinobistack/gokakashi/internal/restapi/v1/integrations"
	"github.com/shinobistack/gokakashi/internal/restapi/v1/policies"
	"github.com/shinobistack/gokakashi/internal/restapi/v1/scans"
	"github.com/spf13/cobra"
	"log"
	"net/http"
	"net/url"
	"strings"
)

var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Scan a container image",
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
	// server and token are Defined in agent.go
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

func scanImage(cmd *cobra.Command, args []string) {

	if server == "" || token == "" || image == "" || policyName == "" {
		log.Fatalf("Error: Missing required inputs. Please provide --server, --token, and --workspace.")
	}

	policy, err := fetchPolicyByName(server, 0, token, policyName)
	if err != nil {
		log.Fatalf("Failed to fetch policy: %v", err)
	}
	log.Println(policy)

	integration, err := fetchIntegrationByName(server, 0, token, policy.Image.Registry)
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
	}
	log.Println(reqBody)

	reqBodyJSON, _ := json.Marshal(reqBody)
	url := constructURL(server, 0, "/api/v1/scans")

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(reqBodyJSON))

	if err != nil {
		log.Fatalf("Failed to create scan request: %v", err)
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf("Failed to send scan request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Server responded with status: %d", resp.StatusCode)
	}
	log.Println(resp)
	var response scans.CreateScanResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		log.Fatalf("Failed to decode scan response: %v", err)
	}

	log.Printf("Scan triggered successfully! Scan ID: %s, status: %s", response.ID, response.Status)
}

func getScanStatus(cmd *cobra.Command, args []string) {
	if server == "" || token == "" || scanID == "" {
		log.Fatalf("Error: Missing required inputs. Please provide --server, --token, and --workspace.")
	}

	url := constructURL(server, 0, fmt.Sprintf("/api/v1/scans/%s", scanID))

	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		log.Fatalf("Failed to create status request: %v", err)
	}
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := http.DefaultClient.Do(req)
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

func fetchPolicyByName(server string, port int, token string, policyName string) (*policies.GetPolicyResponse, error) {
	url := constructURL(server, port, "/api/v1/policies") + fmt.Sprintf("?name=%s", policyName)

	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return nil, fmt.Errorf("failed to create policy fetch request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := http.DefaultClient.Do(req)
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

func fetchIntegrationByName(server string, port int, token string, integrationName string) (*integrations.GetIntegrationResponse, error) {
	url := constructURL(server, port, "/api/v1/integrations") + fmt.Sprintf("?name=%s", integrationName)

	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return nil, fmt.Errorf("failed to create integration fetch request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := http.DefaultClient.Do(req)
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
