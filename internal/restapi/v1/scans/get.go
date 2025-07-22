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
	Status  string `query:"status"`
	Name    string `query:"name"`
	Page    int    `query:"page"`
	PerPage int    `query:"per_page"`
}

type GetScanRequest struct {
	ID     uuid.UUID             `path:"id"`
	Labels []schema.CommonLabels `json:"labels,omitempty"`
}

type ListScansResponse struct {
	Scans      []GetScanResponse `json:"scans"`
	Page       int               `json:"page"`
	PerPage    int               `json:"per_page"`
	Total      int               `json:"total"`
	TotalPages int               `json:"total_pages"`
}

func ListScans(client *ent.Client) func(ctx context.Context, req ListScanRequest, res *ListScansResponse) error {
	return func(ctx context.Context, req ListScanRequest, res *ListScansResponse) error {
		baseQuery := client.Scans.Query()
		if req.Status != "" {
			baseQuery = baseQuery.Where(scans.Status(req.Status))
		}

		// Count total matching records (without pagination)
		total, err := baseQuery.Clone().Count(ctx)
		if err != nil {
			return status.Wrap(errors.New("failed to count scans"), status.Internal)
		}

		// Pagination logic
		page := req.Page
		perPage := req.PerPage
		if page < 1 {
			page = 1
		}
		if perPage < 1 {
			perPage = 30
		}
		if perPage > 100 {
			perPage = 100
		}
		offset := (page - 1) * perPage

		query := baseQuery.Clone().WithScanLabels().Limit(perPage).Offset(offset)
		scanResults, err := query.All(ctx)
		if err != nil {
			return status.Wrap(errors.New("failed to fetch scan details"), status.Internal)
		}

		resp := ListScansResponse{
			Scans:      make([]GetScanResponse, len(scanResults)),
			Page:       page,
			PerPage:    perPage,
			Total:      total,
			TotalPages: (total + perPage - 1) / perPage,
		}
		for i, scan := range scanResults {
			labels := mapScanLabels(scan.Edges.ScanLabels)
			resp.Scans[i] = GetScanResponse{
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
		*res = resp
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
