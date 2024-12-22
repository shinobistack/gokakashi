package integrations

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/shinobistack/gokakashi/ent"
	"github.com/swaggest/usecase/status"
)

type UpdateIntegrationRequest struct {
	// ToDo: convert to UUID
	ID     uuid.UUID               `path:"id"`
	Name   *string                 `json:"name"`
	Type   *string                 `json:"type"`
	Config *map[string]interface{} `json:"config"`
}

type UpdateIntegrationResponse struct {
	ID     uuid.UUID `json:"id"`
	Status string    `json:"status"`
}

func UpdateIntegration(client *ent.Client) func(ctx context.Context, req UpdateIntegrationRequest, res *GetIntegrationResponse) error {
	return func(ctx context.Context, req UpdateIntegrationRequest, res *GetIntegrationResponse) error {
		// Validate ID
		if req.ID == uuid.Nil {
			return status.Wrap(errors.New("invalid UUID: cannot be nil"), status.InvalidArgument)
		}

		// Check if integration exists
		integrationID, err := client.Integrations.Get(ctx, req.ID)
		if err != nil {
			if ent.IsNotFound(err) {
				return status.Wrap(errors.New("policy not found"), status.NotFound)
			}
			return status.Wrap(fmt.Errorf("unexpected error: %v", err), status.Internal)
		}

		update := client.Integrations.UpdateOne(integrationID)
		if req.Name != nil {
			update = update.SetName(*req.Name)
		}
		if req.Config != nil {
			update = update.SetConfig(*req.Config)
		}

		integration, err := update.Save(ctx)
		if err != nil {
			return status.Wrap(fmt.Errorf("failed to update integration: %v", err), status.Internal)
		}

		res.ID = integration.ID
		res.Name = integration.Name
		res.Type = integration.Type
		res.Config = integration.Config

		return nil
	}
}
