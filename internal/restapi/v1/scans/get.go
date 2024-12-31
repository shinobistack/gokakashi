package scans

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/shinobistack/gokakashi/ent"
	"github.com/shinobistack/gokakashi/ent/schema"
	"github.com/swaggest/usecase/status"
)

type GetScanResponse struct {
	ID            uuid.UUID    `json:"id"`
	PolicyID      uuid.UUID    `json:"policy_id"`
	Image         string       `json:"image"`
	Scanner       string       `json:"scanner"`
	IntegrationID uuid.UUID    `json:"integration_id"`
	Status        string       `json:"status"`
	Check         schema.Check `json:"check"`
	Report        interface{}  `json:"report"`
}

type ListScanRequest struct{}

type GetScanRequest struct {
	ID uuid.UUID `path:"id"`
}

func ListScans(client *ent.Client) func(ctx context.Context, req ListScanRequest, res *[]GetScanResponse) error {
	return func(ctx context.Context, req ListScanRequest, res *[]GetScanResponse) error {
		scans, err := client.Scans.Query().All(ctx)
		if err != nil {
			return status.Wrap(errors.New("failed to fetch scan details"), status.Internal)
		}

		*res = make([]GetScanResponse, len(scans))
		for i, scan := range scans {
			(*res)[i] = GetScanResponse{
				ID:            scan.ID,
				PolicyID:      scan.PolicyID,
				Image:         scan.Image,
				Scanner:       scan.Scanner,
				IntegrationID: scan.IntegrationID,
				Status:        scan.Status,
				Check:         scan.Check,
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
		res.Check = scan.Check
		res.Report = scan.Report
		return nil
	}
}
