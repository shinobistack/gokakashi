package scans

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/shinobistack/gokakashi/ent"
	"github.com/swaggest/usecase/status"
)

type DeleteScanRequest struct {
	ID uuid.UUID `path:"id"`
}

type DeleteScanResponse struct {
	ID     uuid.UUID `json:"id"`
	Status string    `json:"status"`
}

func DeleteScan(client *ent.Client) func(ctx context.Context, req DeleteScanRequest, res *DeleteScanResponse) error {
	return func(ctx context.Context, req DeleteScanRequest, res *DeleteScanResponse) error {
		// Validate ID
		if req.ID == uuid.Nil {
			return status.Wrap(errors.New("invalid UUID: cannot be nil"), status.InvalidArgument)
		}

		err := client.Scans.DeleteOneID(req.ID).Exec(ctx)
		if err != nil {
			return status.Wrap(err, status.Internal)
		}

		res.ID = req.ID
		res.Status = "deleted"
		return nil
	}
}
