package scans

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/shinobistack/gokakashi/ent"
	"github.com/swaggest/usecase/status"
)

type UpdateScanRequest struct {
	ID     uuid.UUID       `path:"id"`
	Status *string         `json:"status"`
	Report json.RawMessage `json:"report,omitempty"`
}

type UpdateScanResponse struct {
	ID     uuid.UUID `json:"id"`
	Status string    `json:"status"`
}

func UpdateScan(client *ent.Client) func(ctx context.Context, req UpdateScanRequest, res *UpdateScanResponse) error {
	return func(ctx context.Context, req UpdateScanRequest, res *UpdateScanResponse) error {
		if req.ID == uuid.Nil {
			return status.Wrap(errors.New("invalid UUID: cannot be nil"), status.InvalidArgument)
		}

		_, err := client.Scans.Get(ctx, req.ID)
		if err != nil {
			if ent.IsNotFound(err) {
				return status.Wrap(errors.New("scan not found"), status.NotFound)
			}
			return status.Wrap(fmt.Errorf("unexpected error: %v", err), status.Internal)
		}

		update := client.Scans.UpdateOneID(req.ID)

		if req.Status != nil {
			update.SetStatus(*req.Status)
		}

		if len(req.Report) > 0 {
			update.SetReport(req.Report)
		}

		updatedScan, err := update.Save(ctx)
		if err != nil {
			return status.Wrap(fmt.Errorf("failed to update scan: %v", err), status.Internal)
		}
		res.ID = updatedScan.ID
		res.Status = updatedScan.Status
		return nil
	}
}
