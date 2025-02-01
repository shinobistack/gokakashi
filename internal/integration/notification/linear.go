package notification

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type LinearIntegration struct {
	APIKey    string `json:"api_key"`
	ProjectID string `json:"project_id"`
	TeamID    string `json:"team_id"`
}

func (l *LinearIntegration) Type() IntegrationType {
	return Linear
}

type LinearIssue struct {
	// Config refers to the linear configuration
	// with which the issue could be created
	Config *LinearIntegration

	Title       string `json:"issue_title"`
	Description string `json:"description"`
	Priority    int    `json:"issue_priority"`
	Assignee    string `json:"issue_assignee_id"`
	StateID     string `json:"issue_state_id"`
	DueDate     string `json:"issue_due_date"`
}

func (l *LinearIssue) Notify(ctx context.Context) error {
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

	variables := map[string]interface{}{
		"input": map[string]interface{}{
			"title":       l.Title,
			"description": l.Description,
			"teamId":      l.Config.TeamID,
			"projectId":   l.Config.ProjectID,
			"priority":    l.Priority,
			"assigneeId":  l.Assignee,
			"stateId":     l.StateID,
			"dueDate":     l.DueDate,
		},
	}

	requestBody := graphQLRequest{
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

	req.Header.Set("Authorization", l.Config.APIKey)
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

	var graphQLResponse graphQLResponse
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

type graphQLRequest struct {
	Query     string                 `json:"query"`
	Variables map[string]interface{} `json:"variables"`
}

type graphQLResponse struct {
	Data   interface{} `json:"data"`
	Errors []struct {
		Message string `json:"message"`
	} `json:"errors"`
}
