package policies

import (
	"context"
	"errors"
	"slices"

	"github.com/google/uuid"
	"github.com/scriptnull/jsonseal"
	"github.com/swaggest/usecase/status"
)

type PostRequest struct {
	Name    string  `json:"name"`
	Trigger Trigger `json:"trigger"`
}

type Trigger struct {
	Type TriggerType `json:"type"`
}

type TriggerType string

var (
	Cron TriggerType = "cron"
	CI   TriggerType = "ci"

	allowedTriggerTypes = []TriggerType{
		Cron,
		CI,
	}
)

func (t TriggerType) Valid() bool {
	return slices.Contains(allowedTriggerTypes, t)
}

type PostResponse struct {
	ID      uuid.UUID `json:"id"`
	Name    string    `json:"name"`
	Trigger Trigger   `json:"trigger"`
}

var (
	ErrNotFound           error = errors.New("not found")
	ErrInvalidTriggerType error = errors.New("invalid trigger type")
)

func (req *PostRequest) Validate() error {
	var check jsonseal.CheckGroup

	check.Field("name").Check(func() error {
		if req.Name == "" {
			return ErrNotFound
		}
		return nil
	})

	check.Field("trigger.type").Check(func() error {
		if !req.Trigger.Type.Valid() {
			return ErrInvalidTriggerType
		}

		return nil
	})

	return check.Validate()
}

func Post(_ context.Context, req PostRequest, res *PostResponse) error {
	err := req.Validate()
	if err != nil {
		return status.Wrap(err, status.InvalidArgument)
	}

	res.ID = uuid.New()

	return nil
}
