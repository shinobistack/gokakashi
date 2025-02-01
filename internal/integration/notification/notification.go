package notification

import (
	"context"
)

type Notifier interface {
	Notify(context.Context) error
}

type IntegrationType string

var Linear IntegrationType = "linear"
