package scans

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/shinobistack/gokakashi/ent"
	"github.com/shinobistack/gokakashi/ent/scanlabels"
	"github.com/swaggest/usecase/status"
	"log"
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

		tx, err := client.Tx(ctx)
		if err != nil {
			return status.Wrap(err, status.Internal)
		}

		// Delete associated labels
		_, err = tx.ScanLabels.Delete().
			Where(scanlabels.ScanID(req.ID)).
			Exec(ctx)

		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				log.Printf("rollback failed: %v\n", rollbackErr)
			}
			return status.Wrap(fmt.Errorf("failed to delete associated labels for scan ID %s: %w", req.ID, err), status.Internal)
		}

		err = tx.Scans.DeleteOneID(req.ID).Exec(ctx)
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				log.Printf("rollback failed: %v\n", rollbackErr)
			}
			return status.Wrap(fmt.Errorf("failed to delete scan: %w", err), status.Internal)
		}

		if err := tx.Commit(); err != nil {
			return status.Wrap(fmt.Errorf("failed to commit transaction: %w", err), status.Internal)
		}

		*res = DeleteScanResponse{
			ID:     req.ID,
			Status: "deleted",
		}
		return nil
	}
}
