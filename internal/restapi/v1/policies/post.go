package policies

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/scriptnull/jsonseal"
	"github.com/swaggest/usecase/status"
)

type PostRequest struct {
	Name string `json:"name"`
}

type PostResponse struct {
	ID uuid.UUID `json:"id"`
}

var (
	ErrNotFound error = errors.New("not found")
)

func (req *PostRequest) Validate() error {
	var check jsonseal.CheckGroup

	check.Field("name").Check(func() error {
		if req.Name == "" {
			return ErrNotFound
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
