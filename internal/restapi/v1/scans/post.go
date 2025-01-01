package scans

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"github.com/shinobistack/gokakashi/ent"
	"github.com/shinobistack/gokakashi/ent/integrations"
	"github.com/shinobistack/gokakashi/ent/policies"
	"github.com/shinobistack/gokakashi/ent/scans"
	"github.com/shinobistack/gokakashi/ent/schema"
	"github.com/swaggest/usecase/status"
)

type CreateScanRequest struct {
	PolicyID uuid.UUID `json:"policy_id"`
	// ToDo: To think if the image stored would be single registery/image:tag.
	Image         string    `json:"image"`
	Scanner       string    `json:"scanner"`
	IntegrationID uuid.UUID `json:"integration_id"`
	//ToDo: Similarly to think if the check should have the evaluate conditions. How would the notify work.
	Check schema.Check `json:"check"`
	// ToDo: can we pre-define the values for scan status that can be used?
	Status string `json:"status"`
	// ToDo: Just the report URL or status of public or private
	Report json.RawMessage `json:"report,omitempty"`
}

type CreateScanResponse struct {
	ID     uuid.UUID `json:"id"`
	Status string    `json:"status"`
}

// ToDo: The server would pick up and filter and give the single image, evaluated to check condition and notify

func CreateScan(client *ent.Client) func(ctx context.Context, req CreateScanRequest, res *CreateScanResponse) error {
	return func(ctx context.Context, req CreateScanRequest, res *CreateScanResponse) error {
		// Validate
		if req.PolicyID == uuid.Nil || req.Image == "" || req.Status == "" {
			return status.Wrap(errors.New("missing required fields"), status.InvalidArgument)
		}
		// ToDo: To think if we need to check image duplicate for scans? IF image already scheduled for scans?
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
		// Check if IntegrationID existence
		integrationExists, err := client.Integrations.Query().
			Where(integrations.ID(req.IntegrationID)).
			Exist(ctx)
		if err != nil {
			return status.Wrap(err, status.Internal)
		}
		if !integrationExists {
			return status.Wrap(errors.New("integration not found"), status.NotFound)
		}
		// Check for duplicate scans for the same image under a policy
		duplicate, err := client.Scans.Query().
			Where(scans.PolicyID(req.PolicyID), scans.Image(req.Image)).
			Exist(ctx)
		if err != nil {
			return status.Wrap(err, status.Internal)
		}
		if duplicate {
			return status.Wrap(errors.New("scan already scheduled for this image"), status.AlreadyExists)
		}

		// Create the scan
		scan, err := client.Scans.Create().
			SetPolicyID(req.PolicyID).
			SetImage(req.Image).
			SetScanner(req.Scanner).
			SetCheck(req.Check).
			SetStatus(req.Status).
			SetIntegrationID(req.IntegrationID).
			SetReport(req.Report).
			Save(ctx)

		if err != nil {
			return status.Wrap(err, status.Internal)
		}

		res.ID = scan.ID
		res.Status = "scan_pending"
		return nil
	}
}
