package scan

import (
	"errors"
	"slices"
)

type Status string

const (
	Pending          Status = "scan_pending"
	InProgress       Status = "scan_in_progress"
	NotifyPending    Status = "notify_pending"
	NotifyInProgress Status = "notify_in_progress"
	Success          Status = "success"
	Error            Status = "error"
)

var (
	ErrInvalidScanStatus = errors.New("invalid scan status")
)

func validStatuses() []Status {
	return []Status{
		Pending,
		InProgress,
		NotifyPending,
		NotifyInProgress,
		Success,
		Error,
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
	if !slices.Contains(Statuses(), s) {
		return ErrInvalidScanStatus
	}
	return nil
}
