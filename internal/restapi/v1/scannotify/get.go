package scannotify

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/shinobistack/gokakashi/ent"
	"github.com/shinobistack/gokakashi/ent/scannotify"
	"github.com/swaggest/usecase/status"
)

type GetScanNotifyRequest struct {
	ID   *uuid.UUID `query:"id"`
	Hash *string    `query:"hash"`
}

type GetScanNotifyItem struct {
	ID     uuid.UUID `json:"id"`
	ScanID uuid.UUID `json:"scan_id"`
	Hash   string    `json:"hash"`
}

type GetScanNotifyResponse struct {
	Count int                 `json:"count"`
	Data  []GetScanNotifyItem `json:"data"`
}

func GetScanNotify(client *ent.Client) func(ctx context.Context, req GetScanNotifyRequest, res *GetScanNotifyResponse) error {
	return func(ctx context.Context, req GetScanNotifyRequest, res *GetScanNotifyResponse) error {
		if res == nil {
			return status.Wrap(errors.New("response cannot be nil"), status.InvalidArgument)
		}

		// Initialize empty response
		*res = GetScanNotifyResponse{
			Count: 0,
			Data:  []GetScanNotifyItem{},
		}

		// Query by ID
		if req.ID != nil {
			notification, err := client.ScanNotify.Get(ctx, *req.ID)
			if err != nil {
				if ent.IsNotFound(err) {
					return status.Wrap(errors.New("notification not found"), status.NotFound)
				}
				return status.Wrap(err, status.Internal)
			}

			res.Data = append(res.Data, GetScanNotifyItem{
				ID:     notification.ID,
				ScanID: notification.ScanID,
				Hash:   notification.Hash,
			})
			res.Count = 1
			return nil
		}

		// Query by Hash
		if req.Hash != nil {
			notifications, err := client.ScanNotify.Query().
				Where(scannotify.HashEQ(*req.Hash)).
				All(ctx)

			if err != nil {
				return status.Wrap(err, status.Internal)
			}

			if len(notifications) == 0 {
				return status.Wrap(errors.New("no notifications found for the given hash"), status.NotFound)
			}

			for _, notification := range notifications {
				res.Data = append(res.Data, GetScanNotifyItem{
					ID:     notification.ID,
					ScanID: notification.ScanID,
					Hash:   notification.Hash,
				})
			}

			res.Count = len(notifications)
			return nil
		}

		// Query all notifications
		notifications, err := client.ScanNotify.Query().All(ctx)
		if err != nil {
			return status.Wrap(err, status.Internal)
		}

		for _, notification := range notifications {
			res.Data = append(res.Data, GetScanNotifyItem{
				ID:     notification.ID,
				ScanID: notification.ScanID,
				Hash:   notification.Hash,
			})
		}

		res.Count = len(notifications)
		return nil
	}
}
