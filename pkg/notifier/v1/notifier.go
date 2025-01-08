package notifier

import (
	"fmt"
)

func NewNotifier(tool string) (Notifier, error) {
	switch tool {
	case "linear":
		return &LinearNotifier{}, nil
	// Add other tools like Jira here
	//case "jira":
	//	return &JiraNotifier{}, nil
	default:
		return nil, fmt.Errorf("unsupported notifier tool: %s", tool)
	}
}
