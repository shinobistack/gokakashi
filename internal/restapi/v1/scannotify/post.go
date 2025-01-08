package scannotify

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/shinobistack/gokakashi/ent"
	"github.com/shinobistack/gokakashi/ent/scannotify"
	"github.com/swaggest/usecase/status"
)

type CreateScanNotifyRequest struct {
	ScanID uuid.UUID `json:"scan_id"`
	Hash   string    `json:"hash"`
}

type CreateScanNotifyResponse struct {
	ScanID uuid.UUID `json:"scan_id"`
	Hash   string    `json:"hash"`
}

func CreateScanNotify(client *ent.Client) func(ctx context.Context, req CreateScanNotifyRequest, res *CreateScanNotifyResponse) error {
	return func(ctx context.Context, req CreateScanNotifyRequest, res *CreateScanNotifyResponse) error {
		// Validate input
		if req.ScanID == uuid.Nil || req.Hash == "" {
			return status.Wrap(errors.New("invalid input: missing fields"), status.InvalidArgument)
		}

		// Check if the notification already exists
		exists, err := client.ScanNotify.Query().
			Where(scannotify.ScanID(req.ScanID), scannotify.Hash(req.Hash)).
			Exist(ctx)
		if err != nil {
			return status.Wrap(err, status.Internal)
		}
		if exists {
			return status.Wrap(errors.New("notification already exists"), status.AlreadyExists)
		}

		// Create the notification
		notification, err := client.ScanNotify.
			Create().
			SetScanID(req.ScanID).
			SetHash(req.Hash).
			Save(ctx)
		if err != nil {
			return status.Wrap(err, status.Internal)
		}

		// Populate the response
		res.ScanID = notification.ScanID
		res.Hash = notification.Hash

		return nil
	}
}
