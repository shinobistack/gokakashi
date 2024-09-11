package notifier

type Notifier interface {
	SendNotification(vulnerabilities []Vulnerability, config NotifyConfig) error
}

type NotifyConfig struct {
	APIKey    string
	ProjectID string
	Title     string
	Priority  int
	Assignee  string
	Label     string
	DueDate   string
	TeamID    string
	StateID   string
}
