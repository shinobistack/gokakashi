// This package will soon be deprecated with introduction of V2 Agents
package v1

import (
	"errors"
	"slices"
)

type Status string

const (
	Connected      Status = "connected"
	ScanInProgress Status = "scan_in_progress"
	Disconnected   Status = "disconnected"
)

var ErrInvalidAgentStatus = errors.New("invalid agent status")

func validStatuses() []Status {
	return []Status{
		Connected,
		ScanInProgress,
		Disconnected,
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
