package notifier

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

// Vulnerability represents a vulnerability in the scan results
type Vulnerability struct {
	VulnerabilityID string `json:"VulnerabilityID"`
	PkgName         string `json:"PkgName"`
	Severity        string `json:"Severity"`
	Description     string `json:"Description"`
}

// Result represents the result field in Trivy output
type Result struct {
	Target          string          `json:"Target"`
	Type            string          `json:"Type"`
	Vulnerabilities []Vulnerability `json:"Vulnerabilities"`
}

// TrivyReport represents the overall Trivy scan report
type TrivyReport struct {
	ArtifactName string   `json:"ArtifactName"`
	Results      []Result `json:"Results"`
}

// GraphQLRequest struct for Linear GraphQL API
type GraphQLRequest struct {
	Query     string                 `json:"query"`
	Variables map[string]interface{} `json:"variables"`
}

// LinearNotifier is an example of a notifier that sends data to Linear's API
type LinearNotifier struct{}

func NewLinearNotifier() *LinearNotifier {
	return &LinearNotifier{}
}

// SendNotification sends the vulnerability report to Linear via GraphQL
func (ln *LinearNotifier) SendNotification(vulnerabilities []Vulnerability, config NotifyConfig) error {
	url := "https://api.linear.app/graphql"

	// Prepare GraphQL mutation
	mutation := `
		mutation CreateIssue($input: IssueCreateInput!) {
			issueCreate(input: $input) {
				success
				issue {
					id
					title
					priority 
					state { id name } 
					assignee { name } 
					dueDate
				}
			}
		}
	`

	// Construct issue description from vulnerabilities
	description := formatVulnerabilityReport(vulnerabilities)

	// Prepare variables for the GraphQL mutation
	variables := map[string]interface{}{
		"input": map[string]interface{}{
			"title":       config.Title,
			"description": description,
			"teamId":      config.TeamID,    // Ensure this is a valid UUID
			"projectId":   config.ProjectID, // Ensure this is a valid UUID, if needed
			"priority":    config.Priority,
			"assigneeId":  config.Assignee, // Ensure this is a valid UUID
			"stateId":     config.StateID,  // Ensure this is a valid UUID
			"dueDate":     config.DueDate,  // Ensure correct date format (YYYY-MM-DD)
		},
	}

	// Create GraphQL request
	requestBody := GraphQLRequest{
		Query:     mutation,
		Variables: variables,
	}

	// Convert request body to JSON
	jsonPayload, err := json.Marshal(requestBody)
	if err != nil {
		return fmt.Errorf("failed to marshal GraphQL request: %v", err)
	}

	// Prepare HTTP request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return fmt.Errorf("failed to create new HTTP request: %v", err)
	}

	// Set headers
	req.Header.Set("Authorization", config.APIKey)
	req.Header.Set("Content-Type", "application/json")

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send GraphQL request: %v", err)
	}
	defer resp.Body.Close()

	// Check the response status
	if resp.StatusCode != http.StatusOK {
		// Read and log the response body for debugging
		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		bodyString := string(bodyBytes)
		return fmt.Errorf("failed to create issue: status code %d, response: %s", resp.StatusCode, bodyString)
	}

	log.Println("Successfully created issue in Linear using GraphQL API.")
	return nil
}

// formatVulnerabilityReport converts vulnerabilities to a description for the issue
func formatVulnerabilityReport(vulnerabilities []Vulnerability) string {
	report := "Detected Vulnerabilities:\n"
	for _, vuln := range vulnerabilities {
		report += fmt.Sprintf("ID: %s, Severity: %s, Package: %s\n", vuln.VulnerabilityID, vuln.Severity, vuln.PkgName)
	}
	return report
}
