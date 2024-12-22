package integrations

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/shinobistack/gokakashi/ent"
	"github.com/shinobistack/gokakashi/ent/integrations"
	"github.com/swaggest/usecase/status"
)

type UpdateIntegrationRequest struct {
	ID     string                  `path:"id"`
	Name   *string                 `json:"name"`
	Type   *string                 `json:"type"`
	Config *map[string]interface{} `json:"config"`
}

type UpdateIntegrationResponse struct {
	ID     string `json:"id"`
	Status string `json:"status"`
}

func UpdateIntegration(client *ent.Client) func(ctx context.Context, req UpdateIntegrationRequest, res *GetIntegrationResponse) error {
	return func(ctx context.Context, req UpdateIntegrationRequest, res *GetIntegrationResponse) error {
		uid, err := uuid.Parse(req.ID)
		if err != nil {
			return status.Wrap(fmt.Errorf("invalid UUID format: %v", err), status.InvalidArgument)
		}

		// Check if integration exists
		exists, err := client.Integrations.Query().Where(integrations.ID(uid)).Exist(ctx)
		if err != nil {
			return status.Wrap(fmt.Errorf("unexpected database error: %v", err), status.Internal)
		}
		if !exists {
			return status.Wrap(errors.New("integration not found"), status.NotFound)
		}

		update := client.Integrations.UpdateOneID(uid)
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

		res.ID = integration.ID.String()
		res.Name = integration.Name
		res.Type = integration.Type
		res.Config = integration.Config

		return nil
	}
}
