package notifier

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

// TrivyReport represents the overall Trivy scan report
type TrivyReport struct {
	ArtifactName string   `json:"ArtifactName"`
	Results      []Result `json:"Results"`
}

// Result represents the result field in Trivy output
type Result struct {
	Target          string          `json:"Target"`
	Type            string          `json:"Type"`
	Vulnerabilities []Vulnerability `json:"Vulnerabilities"`
}

// Vulnerability represents a vulnerability in the scan results
type Vulnerability struct {
	VulnerabilityID  string `json:"VulnerabilityID"`
	PkgName          string `json:"PkgName"`
	Severity         string `json:"Severity"`
	InstalledVersion string `json:"InstalledVersion"`
	FixedVersion     string `json:"FixedVersion"`
	Title            string `json:"Title"`
	Description      string `json:"Description"`
	PrimaryURL       string `json:"PrimaryURL"`
	Status           string `json:"Status"`
}

// GraphQLRequest struct for Linear GraphQL API
type GraphQLRequest struct {
	Query     string                 `json:"query"`
	Variables map[string]interface{} `json:"variables"`
}

// GraphQLResponse represents the GraphQL response structure
type GraphQLResponse struct {
	Data   interface{} `json:"data"`
	Errors []struct {
		Message string `json:"message"`
	} `json:"errors"`
}

// LinearNotifier is an example of a notifier for Linear
type LinearNotifier struct{}

func NewLinearNotifier() *LinearNotifier {
	return &LinearNotifier{}
}

// SendNotification sends the vulnerability report to Linear via GraphQL
func (ln *LinearNotifier) SendNotification(report TrivyReport, vulnerabilities []Vulnerability, config NotifyConfig) error {
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
	description := formatVulnerabilityReport(report, vulnerabilities)

	// Prepare variables for the GraphQL mutation
	variables := map[string]interface{}{
		"input": map[string]interface{}{
			"title":       config.Title,
			"description": description,
			"teamId":      config.TeamID,    // Ensure this is a valid UUID
			"projectId":   config.ProjectID, // Ensure this is a valid UUID, if needed
			"priority":    config.Priority,
			"assigneeId":  config.Assignee, // Ensure correct assignee ID
			"stateId":     config.StateID,  // Ensure correct state ID
			"dueDate":     config.DueDate,
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

	// Read and log the full response body
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %v", err)
	}
	bodyString := string(bodyBytes)
	log.Printf("Linear API Response: %s", bodyString)

	// Check if the response contains errors
	var graphQLResponse GraphQLResponse
	err = json.Unmarshal(bodyBytes, &graphQLResponse)
	if err != nil {
		return fmt.Errorf("failed to parse GraphQL response: %v", err)
	}

	// Check for errors in the GraphQL response
	if len(graphQLResponse.Errors) > 0 {
		for _, apiError := range graphQLResponse.Errors {
			log.Printf("Linear API error: %s", apiError.Message)
		}
		return fmt.Errorf("linear issue creation failed due to API errors")
	}

	log.Println("Successfully created issue in Linear using GraphQL API.")
	return nil
}

// formatVulnerabilityReport converts vulnerabilities to a description for the issue
func formatVulnerabilityReport(report TrivyReport, vulnerabilities []Vulnerability) string {
	var buffer bytes.Buffer

	// Add image information
	buffer.WriteString(fmt.Sprintf("Image: %s\n\n", report.ArtifactName))

	// Iterate over vulnerabilities and format them in the simplified format
	for _, vuln := range vulnerabilities {
		buffer.WriteString(fmt.Sprintf("Library: %s\n", vuln.PkgName))
		buffer.WriteString(fmt.Sprintf("Vulnerability: %s\n", vuln.VulnerabilityID))
		buffer.WriteString(fmt.Sprintf("Severity: %s\n", vuln.Severity))
		buffer.WriteString(fmt.Sprintf("Status: %s\n", vuln.Status))
		buffer.WriteString(fmt.Sprintf("Installed Version: %s\n", vuln.InstalledVersion))
		buffer.WriteString(fmt.Sprintf("Fixed Version: %s\n", vuln.FixedVersion))
		buffer.WriteString(fmt.Sprintf("Title: %s\n", vuln.Title))
		if vuln.PrimaryURL != "" {
			buffer.WriteString(fmt.Sprintf("More details: %s\n", vuln.PrimaryURL))
		}
		buffer.WriteString("\n") // Add a line break between vulnerabilities
	}

	return buffer.String()
}
