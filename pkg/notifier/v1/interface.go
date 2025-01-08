package notifier

type NotificationConfig struct {
	APIKey    string `json:"api_key"`
	ProjectID string `json:"project_id"`
	TeamID    string `json:"team_id"`
	Title     string `json:"issue_title"`
	Priority  int    `json:"issue_priority"`
	Assignee  string `json:"issue_assignee_id"`
	StateID   string `json:"issue_state_id"`
	DueDate   string `json:"issue_due_date"`
}

type Notifier interface {
	CreateIssue(image string, vulnerabilities []Vulnerability, config NotificationConfig) error
}
