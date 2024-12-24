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

type CreateIntegrationRequest struct {
	Name   string                 `json:"name"`
	Type   string                 `json:"type"`
	Config map[string]interface{} `json:"config"`
}

type CreateIntegrationResponse struct {
	ID     uuid.UUID `json:"id"`
	Status string    `json:"status"`
}

func CreateIntegration(client *ent.Client) func(ctx context.Context, req CreateIntegrationRequest, res *CreateIntegrationResponse) error {
	return func(ctx context.Context, req CreateIntegrationRequest, res *CreateIntegrationResponse) error {
		// Validate required fields
		if req.Name == "" || req.Type == "" {
			return status.Wrap(errors.New("missing required fields: name and/or type"), status.InvalidArgument)
		}

		// Check for duplicate name
		exists, err := client.Integrations.Query().
			Where(integrations.Name(req.Name)).
			Exist(ctx)
		if err != nil {
			return status.Wrap(fmt.Errorf("failed to check for duplicate integration name: %v", err), status.Internal)
		}
		if exists {
			return status.Wrap(errors.New("integration with the same name already exists"), status.AlreadyExists)
		}

		integration, err := client.Integrations.Create().
			SetName(req.Name).
			SetType(req.Type).
			SetConfig(req.Config).
			Save(ctx)
		if err != nil {
			return status.Wrap(fmt.Errorf("failed to create integration: %v", err), status.Internal)
		}

		res.ID = integration.ID
		res.Status = "created"
		return nil

	}

}
