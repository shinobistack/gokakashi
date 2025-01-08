package scannotify

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/shinobistack/gokakashi/ent"
	"github.com/swaggest/usecase/status"
)

type GetScanNotifyRequest struct {
	ID *uuid.UUID `query:"id"`
}

type GetScanNotifyResponse struct {
	ID     uuid.UUID `json:"id"`
	ScanID uuid.UUID `json:"scan_id"`
	Hash   string    `json:"hash"`
}

func GetScanNotify(client *ent.Client) func(ctx context.Context, req GetScanNotifyRequest, res *[]GetScanNotifyResponse) error {
	return func(ctx context.Context, req GetScanNotifyRequest, res *[]GetScanNotifyResponse) error {
		if req.ID != nil {
			notification, err := client.ScanNotify.Get(ctx, *req.ID)
			if err != nil {
				if ent.IsNotFound(err) {
					return status.Wrap(errors.New("notification not found"), status.NotFound)
				}
				return status.Wrap(err, status.Internal)
			}

			*res = []GetScanNotifyResponse{
				{
					ID:     notification.ID,
					ScanID: notification.ScanID,
					Hash:   notification.Hash,
				},
			}
			return nil
		}

		notifications, err := client.ScanNotify.Query().All(ctx)
		if err != nil {
			return status.Wrap(err, status.Internal)
		}

		for _, notification := range notifications {
			*res = append(*res, GetScanNotifyResponse{
				ID:     notification.ID,
				ScanID: notification.ScanID,
				Hash:   notification.Hash,
			})
		}

		return nil
	}
}
