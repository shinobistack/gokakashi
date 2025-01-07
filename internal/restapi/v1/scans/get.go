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
	ID            uuid.UUID        `json:"id"`
	PolicyID      uuid.UUID        `json:"policy_id"`
	Image         string           `json:"image"`
	Scanner       string           `json:"scanner"`
	IntegrationID uuid.UUID        `json:"integration_id"`
	Status        string           `json:"status"`
	Notify        *[]schema.Notify `json:"notify"`
	Report        json.RawMessage  `json:"report,omitempty"`
}

type ListScanRequest struct {
	Status string `query:"status"`
}

type GetScanRequest struct {
	ID uuid.UUID `path:"id"`
}

func ListScans(client *ent.Client) func(ctx context.Context, req ListScanRequest, res *[]GetScanResponse) error {
	return func(ctx context.Context, req ListScanRequest, res *[]GetScanResponse) error {
		query := client.Scans.Query()

		if req.Status != "" {
			query = query.Where(scans.Status(req.Status))
		}

		scanResults, err := query.All(ctx)
		if err != nil {
			return status.Wrap(errors.New("failed to fetch scan details"), status.Internal)
		}

		*res = make([]GetScanResponse, len(scanResults))
		for i, scan := range scanResults {
			(*res)[i] = GetScanResponse{
				ID:            scan.ID,
				PolicyID:      scan.PolicyID,
				Image:         scan.Image,
				Scanner:       scan.Scanner,
				IntegrationID: scan.IntegrationID,
				Status:        scan.Status,
				Notify:        convertToPointer(scan.Notify),
				Report:        scan.Report,
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

		res.ID = scan.ID
		res.PolicyID = scan.PolicyID
		res.IntegrationID = scan.IntegrationID
		res.Image = scan.Image
		res.Scanner = scan.Scanner
		res.Status = scan.Status
		res.Notify = convertToPointer(scan.Notify)
		res.Report = scan.Report
		return nil
	}
}

func convertToPointer(data []schema.Notify) *[]schema.Notify {
	return &data
}
