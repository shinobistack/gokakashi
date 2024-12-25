package scanlabels

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/shinobistack/gokakashi/ent"
	"github.com/shinobistack/gokakashi/ent/scanlabels"
	"github.com/swaggest/usecase/status"
)

type DeleteScanLabelRequest struct {
	ScanID uuid.UUID `path:"scan_id"`
	Key    string    `path:"key"`
}

type DeleteScanLabelResponse struct {
	Status string `json:"status"`
}

func DeleteScanLabel(client *ent.Client) func(ctx context.Context, req DeleteScanLabelRequest, res *DeleteScanLabelResponse) error {
	return func(ctx context.Context, req DeleteScanLabelRequest, res *DeleteScanLabelResponse) error {
		if req.ScanID == uuid.Nil {
			return status.Wrap(errors.New("invalid UUID: cannot be nil"), status.InvalidArgument)
		}
		if req.Key == "" {
			return status.Wrap(errors.New("invalid key: cannot be nil"), status.InvalidArgument)
		}
		_, err := client.ScanLabels.Delete().
			Where(scanlabels.ScanID(req.ScanID), scanlabels.Key(req.Key)).
			Exec(ctx)
		if err != nil {
			return status.Wrap(err, status.Internal)
		}

		res.Status = "deleted"
		return nil
	}
}
