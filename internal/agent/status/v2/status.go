package agent

import (
	"errors"
	"slices"
)

type Status string

/*
         disconnected       +----> lost
             |              |
             |              |
             v              |
  +----> connected ---------+
  |          |              |
  |          |              |
  |          v              |
  +----- scan_in_progress   +----> exited
*/

const (
	Connected      Status = "connected"
	ScanInProgress Status = "scan_in_progress"
	Disconnected   Status = "disconnected"
	Lost           Status = "lost"
	Exited         Status = "exited"
)

var ErrInvalidAgentStatus = errors.New("invalid agent status")

func validStatuses() []Status {
	return []Status{
		Connected,
		ScanInProgress,
		Disconnected,
		Lost,
		Exited,
	}
}

func Statuses() []string {
	statuses := validStatuses()
	result := make([]string, len(statuses))
	for i, s := range statuses {
		result[i] = string(s)
	}
	return result
}

func ValidateStatus(s string) error {
	if !slices.Contains(validStatuses(), Status(s)) {
		return ErrInvalidAgentStatus
	}
	return nil
}
