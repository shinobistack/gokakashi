package scans

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/shinobistack/gokakashi/ent"
	"github.com/shinobistack/gokakashi/ent/scans"
	"github.com/shinobistack/gokakashi/ent/schema"
	"github.com/swaggest/usecase/status"
)

type GetScanResponse struct {
	ID             uuid.UUID             `json:"id"`
	PolicyID       uuid.UUID             `json:"policy_id"`
	Image          string                `json:"image"`
	Scanner        string                `json:"scanner"`
	IntegrationID  *uuid.UUID            `json:"integration_id"`
	Status         string                `json:"status"`
	Labels         []schema.CommonLabels `json:"labels,omitempty"`
	Notify         *[]schema.Notify      `json:"notify,omitempty"`
	Report         json.RawMessage       `json:"report,omitempty"`
	ScannerOptions map[string]string     `json:"scanner_options,omitempty"`
}

type ListScanRequest struct {
	Status string `query:"status"`
	Name   string `query:"name"`
}

type GetScanRequest struct {
	ID     uuid.UUID             `path:"id"`
	Labels []schema.CommonLabels `json:"labels,omitempty"`
}

func ListScans(client *ent.Client) func(ctx context.Context, req ListScanRequest, res *[]GetScanResponse) error {
	return func(ctx context.Context, req ListScanRequest, res *[]GetScanResponse) error {
		query := client.Scans.Query().WithScanLabels()

		if req.Status != "" {
			query = query.Where(scans.Status(req.Status))
		}

		scanResults, err := query.All(ctx)
		if err != nil {
			return status.Wrap(errors.New("failed to fetch scan details"), status.Internal)
		}

		*res = make([]GetScanResponse, len(scanResults))
		for i, scan := range scanResults {
			labels := mapScanLabels(scan.Edges.ScanLabels)

			(*res)[i] = GetScanResponse{
				ID:             scan.ID,
				PolicyID:       scan.PolicyID,
				Image:          scan.Image,
				Scanner:        scan.Scanner,
				IntegrationID:  scan.IntegrationID,
				Status:         scan.Status,
				Notify:         &scan.Notify,
				Report:         scan.Report,
				Labels:         labels,
				ScannerOptions: scan.ScannerOptions,
			}
		}
		return nil
	}
}

func GetScan(client *ent.Client) func(ctx context.Context, req GetScanRequest, res *GetScanResponse) error {
	return func(ctx context.Context, req GetScanRequest, res *GetScanResponse) error {
		if req.ID == uuid.Nil {
			return status.Wrap(errors.New("invalid UUID: cannot be nil"), status.InvalidArgument)
		}

		scan, err := client.Scans.Get(ctx, req.ID)
		if err != nil {
			if ent.IsNotFound(err) {
				return status.Wrap(errors.New("scan not found"), status.NotFound)
			}
			return status.Wrap(fmt.Errorf("unexpected error: %v", err), status.Internal)
		}

		labels := mapScanLabels(scan.Edges.ScanLabels)

		*res = GetScanResponse{
			ID:             scan.ID,
			PolicyID:       scan.PolicyID,
			Image:          scan.Image,
			Scanner:        scan.Scanner,
			IntegrationID:  scan.IntegrationID,
			Status:         scan.Status,
			Notify:         &scan.Notify,
			Report:         scan.Report,
			Labels:         labels,
			ScannerOptions: scan.ScannerOptions,
		}

		return nil
	}
}

func mapScanLabels(labels []*ent.ScanLabels) []schema.CommonLabels {
	if labels == nil {
		return nil
	}

	mapped := make([]schema.CommonLabels, len(labels))
	for i, label := range labels {
		mapped[i] = schema.CommonLabels{
			Key:   label.Key,
			Value: label.Value,
		}
	}
	return mapped
}
