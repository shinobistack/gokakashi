package scans

import (
	"context"
	"encoding/json"
	"errors"
	"log"

	"github.com/google/uuid"
	"github.com/shinobistack/gokakashi/ent"
	"github.com/shinobistack/gokakashi/ent/integrations"
	"github.com/shinobistack/gokakashi/ent/policies"
	"github.com/shinobistack/gokakashi/ent/schema"
	"github.com/swaggest/usecase/status"
)

type CreateScanRequest struct {
	PolicyID       uuid.UUID             `json:"policy_id"`
	Image          string                `json:"image"`
	Scanner        string                `json:"scanner"`
	IntegrationID  *uuid.UUID            `json:"integration_id,omitempty"`
	Notify         []schema.Notify       `json:"notify"`
	Status         string                `json:"status"`
	Report         json.RawMessage       `json:"report,omitempty"`
	Labels         []schema.CommonLabels `json:"labels,omitempty"`
	ScannerOptions map[string]string     `json:"scanner_options,omitempty"`
}

type CreateScanResponse struct {
	ID     uuid.UUID `json:"id"`
	Status string    `json:"status"`
}

func CreateScan(client *ent.Client) func(ctx context.Context, req CreateScanRequest, res *CreateScanResponse) error {
	return func(ctx context.Context, req CreateScanRequest, res *CreateScanResponse) error {
		if req.PolicyID == uuid.Nil || req.Image == "" || req.Status == "" {
			return status.Wrap(errors.New("missing required fields"), status.InvalidArgument)
		}

		// Validate notify structure
		for _, notify := range req.Notify {
			if notify.To == "" {
				return status.Wrap(errors.New("notify 'to' field is required"), status.InvalidArgument)
			}
			if notify.When == "" {
				return status.Wrap(errors.New("notify 'when' field is required"), status.InvalidArgument)
			}
		}

		// Check if PolicyID exists before creating the scan.
		exists, err := client.Policies.Query().
			Where(policies.ID(req.PolicyID)).
			Exist(ctx)
		if err != nil {
			return status.Wrap(err, status.Internal)
		}
		if !exists {
			return status.Wrap(errors.New("policy not found"), status.NotFound)
		}

		// Check if IntegrationID exist before creating the scan.
		if req.IntegrationID != nil {
			integrationExists, err := client.Integrations.Query().
				Where(integrations.ID(*req.IntegrationID)).
				Exist(ctx)
			if err != nil {
				return status.Wrap(err, status.Internal)
			}
			if !integrationExists {
				return status.Wrap(errors.New("integration not found"), status.NotFound)
			}
		}

		tx, err := client.Tx(ctx)
		if err != nil {
			return status.Wrap(err, status.Internal)
		}

		// Create the scan
		scanCreate := tx.Scans.Create().
			SetPolicyID(req.PolicyID).
			SetImage(req.Image).
			SetScanner(req.Scanner).
			SetNotify(req.Notify).
			SetStatus(req.Status).
			SetReport(req.Report)

		if req.IntegrationID != nil {
			scanCreate.SetIntegrationID(*req.IntegrationID)
		}
		if len(req.ScannerOptions) > 0 {
			scanCreate.SetScannerOptions(req.ScannerOptions)
		}

		scan, err := scanCreate.Save(ctx)
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				log.Printf("rollback failed: %v\n", rollbackErr)
			}
			return status.Wrap(err, status.Internal)
		}

		// Save agent labels
		if len(req.Labels) > 0 {
			bulk := make([]*ent.ScanLabelsCreate, len(req.Labels))
			for i, label := range req.Labels {
				bulk[i] = tx.ScanLabels.Create().
					SetScanID(scan.ID).
					SetKey(label.Key).
					SetValue(label.Value)
			}

			if _, err := tx.ScanLabels.CreateBulk(bulk...).Save(ctx); err != nil {
				if rollbackErr := tx.Rollback(); rollbackErr != nil {
					log.Printf("rollback failed: %v\n", rollbackErr)
				}
				return status.Wrap(err, status.Internal)
			}
		}

		if err := tx.Commit(); err != nil {
			return status.Wrap(err, status.Internal)
		}

		*res = CreateScanResponse{
			ID:     scan.ID,
			Status: scan.Status,
		}
		return nil
	}
}
