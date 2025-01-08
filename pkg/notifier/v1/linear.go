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
type Report struct {
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

// //
type LinearNotifier struct{}

func NewLinearNotifier() *LinearNotifier {
	return &LinearNotifier{}
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

func (ln *LinearNotifier) CreateIssue(image string, vulnerabilities []Vulnerability, config NotificationConfig) error {
	url := "https://api.linear.app/graphql"

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

	// Construct the description from vulnerabilities
	description := FormatVulnerabilityReport(image, vulnerabilities)

	// Concatenate the scan.image and issue_title for the Linear issue title
	issueTitle := fmt.Sprintf("%s - %s", image, config.Title)

	// Prepare variables for GraphQL mutation
	variables := map[string]interface{}{
		"input": map[string]interface{}{
			"title":       issueTitle,
			"description": description,
			"teamId":      config.TeamID,
			"projectId":   config.ProjectID,
			"priority":    config.Priority,
			"assigneeId":  config.Assignee,
			"stateId":     config.StateID,
			"dueDate":     config.DueDate,
		},
	}

	// Prepare GraphQL request body
	requestBody := GraphQLRequest{
		Query:     mutation,
		Variables: variables,
	}

	jsonPayload, err := json.Marshal(requestBody)
	if err != nil {
		return fmt.Errorf("failed to marshal GraphQL request: %v", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return fmt.Errorf("failed to create HTTP request: %v", err)
	}

	req.Header.Set("Authorization", config.APIKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request to Linear: %v", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read Linear API response: %v", err)
	}

	log.Printf("Linear API Response: %s", string(bodyBytes))

	var graphQLResponse GraphQLResponse
	err = json.Unmarshal(bodyBytes, &graphQLResponse)
	if err != nil {
		return fmt.Errorf("failed to parse Linear API response: %v", err)
	}

	if len(graphQLResponse.Errors) > 0 {
		for _, apiError := range graphQLResponse.Errors {
			log.Printf("Linear API error: %s", apiError.Message)
		}
		return fmt.Errorf("issue creation in Linear failed")
	}

	log.Println("Issue created successfully in Linear.")
	return nil
}

func FormatVulnerabilityReport(image string, vulnerabilities []Vulnerability) string {
	var buffer bytes.Buffer

	// Add image information
	buffer.WriteString(fmt.Sprintf("Image: %s\n\n", image))

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
