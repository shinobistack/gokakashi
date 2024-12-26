package scanlabels

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/shinobistack/gokakashi/ent"
	"github.com/shinobistack/gokakashi/ent/scanlabels"
	"github.com/swaggest/usecase/status"
)

type UpdateScanLabelRequest struct {
	ScanID uuid.UUID   `path:"scan_id"`
	Labels []ScanLabel `json:"labels"`
	Key    *string     `json:"key"`
	Value  *string     `json:"value"`
}

type UpdateScanLabelResponse struct {
	Labels []ScanLabel `json:"labels"`
}

func UpdateScanLabel(client *ent.Client) func(ctx context.Context, req UpdateScanLabelRequest, res *UpdateScanLabelResponse) error {
	return func(ctx context.Context, req UpdateScanLabelRequest, res *UpdateScanLabelResponse) error {
		if req.ScanID == uuid.Nil {
			return status.Wrap(errors.New("invalid UUID: cannot be nil"), status.InvalidArgument)
		}

		if req.Key != nil && req.Value != nil {
			// Update a specific label
			_, err := client.ScanLabels.Update().
				Where(scanlabels.ScanID(req.ScanID), scanlabels.Key(*req.Key)).
				SetValue(*req.Value).
				Save(ctx)
			if err != nil {
				if ent.IsNotFound(err) {
					return status.Wrap(errors.New("label not found"), status.NotFound)
				}
				return status.Wrap(err, status.Internal)
			}

			updatedLabel, err := client.ScanLabels.Query().
				Where(scanlabels.ScanID(req.ScanID), scanlabels.Key(*req.Key)).
				Only(ctx)
			if err != nil {
				if ent.IsNotFound(err) {
					return status.Wrap(errors.New("label not found after update"), status.NotFound)
				}
				return status.Wrap(err, status.Internal)
			}

			// Return the updated label
			res.Labels = []ScanLabel{
				{
					Key:   updatedLabel.Key,
					Value: updatedLabel.Value,
				},
			}
			return nil
		}

		if req.Labels != nil {
			// Replace all labels
			tx, err := client.Tx(ctx)
			if err != nil {
				return status.Wrap(err, status.Internal)
			}

			// Delete existing labels
			_, err = tx.ScanLabels.Delete().
				Where(scanlabels.ScanID(req.ScanID)).
				Exec(ctx)
			if err != nil {
				if rollbackErr := tx.Rollback(); rollbackErr != nil {
					fmt.Printf("rollback failed: %v\n", rollbackErr)
				}
				return status.Wrap(err, status.Internal)
			}

			// Add new labels
			bulk := make([]*ent.ScanLabelsCreate, len(req.Labels))
			for i, label := range req.Labels {
				bulk[i] = tx.ScanLabels.Create().
					SetScanID(req.ScanID).
					SetKey(label.Key).
					SetValue(label.Value)
			}

			createdLabels, err := tx.ScanLabels.CreateBulk(bulk...).Save(ctx)
			if err != nil {
				if rollbackErr := tx.Rollback(); rollbackErr != nil {
					fmt.Printf("rollback failed: %v\n", rollbackErr)
				}
				return status.Wrap(err, status.Internal)
			}

			err = tx.Commit()
			if err != nil {
				return status.Wrap(err, status.Internal)
			}

			// Build the response
			res.Labels = make([]ScanLabel, len(createdLabels))
			for i, label := range createdLabels {
				res.Labels[i] = ScanLabel{
					Key:   label.Key,
					Value: label.Value,
				}
			}
			return nil
		}

		return status.Wrap(errors.New("invalid request: either key-value or labels must be provided"), status.InvalidArgument)
	}
}
