package scanlabels

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/shinobistack/gokakashi/ent"
	"github.com/shinobistack/gokakashi/ent/scanlabels"
	"github.com/shinobistack/gokakashi/ent/scans"
	"github.com/swaggest/usecase/status"
)

// ToDo: To have a requests to create labels in bulk?
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

		// Check if the ScanID exists
		exists, err := client.Scans.Query().
			Where(scans.ID(req.ScanID)).
			Exist(ctx)
		if err != nil {
			return status.Wrap(err, status.Internal)
		}
		if !exists {
			return status.Wrap(errors.New("scan not found"), status.NotFound)
		}

		// Check if the label already exists
		labelExists, err := client.ScanLabels.Query().
			Where(scanlabels.ScanID(req.ScanID), scanlabels.Key(req.Key)).
			Exist(ctx)
		if err != nil {
			return status.Wrap(err, status.Internal)
		}
		if labelExists {
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
