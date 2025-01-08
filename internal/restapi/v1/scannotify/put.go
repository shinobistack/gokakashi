package scannotify

import (
	"context"
	"errors"
	"github.com/shinobistack/gokakashi/ent/scannotify"

	"github.com/google/uuid"
	"github.com/shinobistack/gokakashi/ent"
	"github.com/swaggest/usecase/status"
)

type UpdateScanNotifyRequest struct {
	ID   uuid.UUID `path:"scan_id"`
	Hash string    `json:"hash"`
}

type UpdateScanNotifyResponse struct {
	ID     uuid.UUID `json:"id"`
	ScanID uuid.UUID `json:"scan_id"`
	Hash   string    `json:"hash"`
}

func UpdateScanNotify(client *ent.Client) func(ctx context.Context, req UpdateScanNotifyRequest, res *UpdateScanNotifyResponse) error {
	return func(ctx context.Context, req UpdateScanNotifyRequest, res *UpdateScanNotifyResponse) error {
		if req.ID == uuid.Nil || req.Hash == "" {
			return status.Wrap(errors.New("invalid input: missing fields"), status.InvalidArgument)
		}

		// Query the notification by scan_id
		notification, err := client.ScanNotify.
			Query().
			Where(scannotify.ScanID(req.ID)).
			Only(ctx) // Expect exactly one result
		if err != nil {
			if ent.IsNotFound(err) {
				return status.Wrap(errors.New("not found"), status.NotFound)
			}
			return status.Wrap(err, status.Internal)
		}

		// Update the hash
		updatedNotification, err := client.ScanNotify.
			UpdateOne(notification).
			SetHash(req.Hash).
			Save(ctx)
		if err != nil {
			return status.Wrap(err, status.Internal)
		}

		// Populate the response
		res.ID = updatedNotification.ID
		res.ScanID = updatedNotification.ScanID
		res.Hash = updatedNotification.Hash

		return nil
	}
}
