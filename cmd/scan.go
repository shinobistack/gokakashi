package cmd

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/google/uuid"
	"github.com/shinobistack/gokakashi/ent/schema"
	"github.com/shinobistack/gokakashi/internal/helper"
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
	image       string
	policyName  string
	scanID      string
	scanTimeout string
)

func scanImage(cmd *cobra.Command, args []string) {
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

	policy, err := fetchPolicyByName(cmd.Context())
	if err != nil {
		log.Fatalf("Failed to fetch policy: %v", err)
	}

	triggerType, ok := policy.Trigger["type"].(string)
	if !ok {
		log.Fatalf("Error: Unable to determine trigger type for policy %s", policyName)
	}

	switch triggerType {
	case "ci":
		// No image field (local images)
		if policy.Image.Registry == "" && policy.Image.Name == "" && policy.Image.Tags == nil {
			handleNotifyAndScan(cmd.Context(), policy, nil, parsedLabels)
			return
		}

		// Only registry specified
		if policy.Image.Registry != "" && policy.Image.Name == "" && policy.Image.Tags == nil {
			integrationID, err := fetchIntegrationByName(cmd.Context(), policy.Image.Registry)
			if err != nil {
				log.Fatalf("Failed to fetch integration: %v", err)
			}
			handleNotifyAndScan(cmd.Context(), policy, integrationID, parsedLabels)
			return
		}
		log.Fatalf("Unsupported CI configuration for policy %s", policyName)
	default:
		log.Println("Trigger type is Cron. Processing scheduled scan... (will support soon)")

		integrationID, err := fetchIntegrationByName(cmd.Context(), policy.Image.Registry)
		if err != nil {
			log.Fatalf("Failed to fetch integration: %v", err)
		}

		handleNotifyAndScan(cmd.Context(), policy, integrationID, parsedLabels)
	}

}

func handleNotifyAndScan(ctx context.Context, policy *policies.GetPolicyResponse, integrationID *uuid.UUID, labels []schema.CommonLabels) {
	if policy.Notify != nil && len(*policy.Notify) > 0 {
		var formattedNotifies []schema.Notify

		for _, notify := range *policy.Notify {
			notifyDetails, err := fetchIntegrationByName(ctx, notify.To)
			if err != nil {
				log.Fatalf("Failed to fetch integration: %v", err)
			}

			formattedNotifies = append(formattedNotifies, schema.Notify{
				To:          notifyDetails.String(),
				When:        notify.When,
				Format:      notify.Format,
				Fingerprint: notify.Fingerprint,
			})
		}

		_, err := postScanDetails(ctx, policy.ID, policy.Scanner, integrationID, labels, formattedNotifies)
		if err != nil {
			log.Fatalf("Failed to post scan requests")
		}
	} else {
		log.Println("No notify found in policy. Skipping...")
	}
}

func postScanDetails(ctx context.Context, policyID uuid.UUID, scanner string, integrationID *uuid.UUID, labels []schema.CommonLabels, notify []schema.Notify) (*scans.CreateScanResponse, error) {
	reqBody := scans.CreateScanRequest{
		PolicyID:      policyID,
		Image:         image,
		Scanner:       scanner,
		IntegrationID: integrationID,
		Notify:        notify,
		Status:        "scan_pending",
		Labels:        labels,
	}

	// Add timeout if specified
	if scanTimeout != "" {
		reqBody.ScannerOptions = map[string]string{
			"timeout": scanTimeout,
		}
	}

	reqBodyJSON, _ := json.Marshal(reqBody)
	url := helper.ConstructURL(server, "/api/v1/scans")

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(reqBodyJSON))
	if err != nil {
		log.Fatalf("Failed to create scan request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := ctx.Value(httpClientKey{}).(*client.Client).Do(req)
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
	return &response, nil
}

func getScanStatus(cmd *cobra.Command, args []string) {
	ctx := cmd.Context()
	httpClient := ctx.Value(httpClientKey{}).(*client.Client)

	if server == "" || token == "" || scanID == "" {
		log.Fatalf("Error: Missing required inputs. Please provide --server, --token, and --scanID.")
	}

	url := helper.ConstructURL(server, fmt.Sprintf("/api/v1/scans/%s", scanID))

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

func fetchPolicyByName(ctx context.Context) (*policies.GetPolicyResponse, error) {
	url := helper.ConstructURL(server, "/api/v1/policies") + fmt.Sprintf("?name=%s", policyName)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create policy fetch request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := ctx.Value(httpClientKey{}).(*client.Client).Do(req)
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

func fetchIntegrationByName(ctx context.Context, integrationName string) (*uuid.UUID, error) {
	url := helper.ConstructURL(server, "/api/v1/integrations") + fmt.Sprintf("?name=%s", integrationName)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create integration fetch request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := ctx.Value(httpClientKey{}).(*client.Client).Do(req)
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

	return &integration[0].ID, nil
}
