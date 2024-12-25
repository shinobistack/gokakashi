package scanlabels

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/shinobistack/gokakashi/ent"
	"github.com/shinobistack/gokakashi/ent/scanlabels"
	"github.com/swaggest/usecase/status"
)

type CreateScanLabelRequest struct {
	ScanID uuid.UUID `path:"scan_id"`
	Key    string    `json:"key"`
	Value  string    `json:"value"`
}

type CreateScanLabelResponse struct {
	ScanID uuid.UUID `path:"scan_id"`
	Key    string    `json:"key"`
	Value  string    `json:"value"`
}

func CreateScanLabel(client *ent.Client) func(ctx context.Context, req CreateScanLabelRequest, res *CreateScanLabelResponse) error {
	return func(ctx context.Context, req CreateScanLabelRequest, res *CreateScanLabelResponse) error {
		if req.ScanID == uuid.Nil || req.Key == "" || req.Value == "" {
			return status.Wrap(errors.New("invalid input: missing fields"), status.InvalidArgument)
		}

		// Check if the label already exists
		exists, err := client.ScanLabels.Query().
			Where(scanlabels.ScanID(req.ScanID), scanlabels.Key(req.Key)).
			Exist(ctx)
		if err != nil {
			return status.Wrap(err, status.Internal)
		}
		if exists {
			return status.Wrap(errors.New("label already exists"), status.AlreadyExists)
		}

		label, err := client.ScanLabels.Create().
			SetScanID(req.ScanID).
			SetKey(req.Key).
			SetValue(req.Value).
			Save(ctx)
		if err != nil {
			return status.Wrap(err, status.Internal)
		}

		res.ScanID = label.ScanID
		res.Key = label.Key
		res.Value = label.Value

		return nil
	}
}
