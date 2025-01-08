package scannotify

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/shinobistack/gokakashi/ent"
	"github.com/shinobistack/gokakashi/ent/scannotify"
	"github.com/swaggest/usecase/status"
)

// DeleteScanNotifyRequest defines the input for deleting a ScanNotify.
type DeleteScanNotifyRequest struct {
	ScanID uuid.UUID `path:"scan_id"`
}

type DeleteScanNotifyResponse struct {
	Status string `json:"status"`
}

func DeleteScanNotify(client *ent.Client) func(ctx context.Context, req DeleteScanNotifyRequest, res *DeleteScanNotifyResponse) error {
	return func(ctx context.Context, req DeleteScanNotifyRequest, res *DeleteScanNotifyResponse) error {
		if req.ScanID == uuid.Nil {
			return status.Wrap(errors.New("invalid input: missing scan_id"), status.InvalidArgument)
		}

		_, err := client.ScanNotify.
			Delete().
			Where(scannotify.ScanID(req.ScanID)).
			Exec(ctx)
		if err != nil {
			if ent.IsNotFound(err) {
				return status.Wrap(errors.New("notification not found"), status.NotFound)
			}
			return status.Wrap(err, status.Internal)
		}
		res.Status = "deleted"
		return nil
	}
}
