package task

import (
	"errors"
	"slices"
)

type Status string

const (
	Pending    Status = "pending"
	InProgress Status = "in_progress"
	Complete   Status = "complete"
)

var ErrInvalidStatus = errors.New("invalid status")

func validStatuses() []Status {
	return []Status{
		Pending,
		InProgress,
		Complete,
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
		return ErrInvalidStatus
	}
	return nil
}
