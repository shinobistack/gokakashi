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
	ID     string                  `path:"id"`
	Name   *string                 `json:"name"`
	Type   *string                 `json:"type"`
	Config *map[string]interface{} `json:"config"`
}

type UpdateIntegrationResponse struct {
	ID     string `json:"id"`
	Status string `json:"status"`
}

func UpdateIntegration(client *ent.Client) func(ctx context.Context, req UpdateIntegrationRequest, res *UpdateIntegrationResponse) error {
	return func(ctx context.Context, req UpdateIntegrationRequest, res *UpdateIntegrationResponse) error {
		uid, err := uuid.Parse(req.ID)
		if err != nil {
			return status.Wrap(errors.New(fmt.Sprintf("invalid UUID format: %v", err)), status.InvalidArgument)
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
			return status.Wrap(errors.New(fmt.Sprintf("failed to update integration: %v", err)), status.Internal)
		}

		res.ID = integration.ID.String()
		res.Status = "updated"

		return nil
	}
}
