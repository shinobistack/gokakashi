package scanlabels

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/shinobistack/gokakashi/ent"
	"github.com/shinobistack/gokakashi/ent/scanlabels"
	"github.com/swaggest/usecase/status"
)

type ListScanLabelsRequest struct {
	ScanID uuid.UUID `path:"scan_id"`
	Keys   []string  `query:"key"`
}

type ListScanLabelsResponse struct {
	Labels []ScanLabel `json:"labels"`
}

type ScanLabel struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type GetScanLabelRequest struct {
	ScanID uuid.UUID `path:"scan_id"`
	Key    string    `path:"key"`
}

type GetScanLabelResponse struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func ListScanLabels(client *ent.Client) func(ctx context.Context, req ListScanLabelsRequest, res *ListScanLabelsResponse) error {
	return func(ctx context.Context, req ListScanLabelsRequest, res *ListScanLabelsResponse) error {
		if req.ScanID == uuid.Nil {
			return status.Wrap(errors.New("invalid UUID: cannot be nil"), status.InvalidArgument)
		}
		query := client.ScanLabels.Query().Where(scanlabels.ScanID(req.ScanID))

		// Filter by keys if provided
		if len(req.Keys) > 0 {
			query = query.Where(scanlabels.KeyIn(req.Keys...))
		}

		labels, err := query.All(ctx)
		if err != nil {
			return status.Wrap(err, status.Internal)
		}

		res.Labels = make([]ScanLabel, len(labels))
		for i, label := range labels {
			res.Labels[i] = ScanLabel{
				Key:   label.Key,
				Value: label.Value,
			}
		}
		return nil
	}
}

func GetScanLabel(client *ent.Client) func(ctx context.Context, req GetScanLabelRequest, res *GetScanLabelResponse) error {
	return func(ctx context.Context, req GetScanLabelRequest, res *GetScanLabelResponse) error {
		// Validate inputs
		if req.ScanID == uuid.Nil || req.Key == "" {
			return status.Wrap(errors.New("invalid Scan ID or Key"), status.InvalidArgument)
		}

		// Query the label
		label, err := client.ScanLabels.Query().
			Where(scanlabels.ScanID(req.ScanID), scanlabels.Key(req.Key)).
			Only(ctx)
		if err != nil {
			if ent.IsNotFound(err) {
				return status.Wrap(errors.New("label not found"), status.NotFound)
			}
			return status.Wrap(err, status.Internal)
		}

		res.Key = label.Key
		res.Value = label.Value
		return nil
	}
}
