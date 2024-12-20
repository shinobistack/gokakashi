package integrations

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/shinobistack/gokakashi/ent"
	"github.com/swaggest/usecase/status"
)

type GetIntegrationRequests struct {
	ID string `path:"id"`
}

type GetIntegrationResponse struct {
	ID     string                 `json:"id"`
	Name   string                 `json:"name"`
	Type   string                 `json:"type"`
	Config map[string]interface{} `json:"config"`
}

type ListIntegrationResponse struct {
	Integrations []GetIntegrationResponse `json:"integrations"`
}

func GetIntegration(client *ent.Client) func(ctx context.Context, req GetIntegrationRequests, res *GetIntegrationResponse) error {
	return func(ctx context.Context, req GetIntegrationRequests, res *GetIntegrationResponse) error {
		// Convert string to uuid.UUID
		uid, err := uuid.Parse(req.ID)
		if err != nil {
			return status.Wrap(errors.New("invalid UUID format"), status.InvalidArgument)
		}

		// Fetch integration by ID
		integration, err := client.Integrations.Get(ctx, uid)
		if err != nil {
			if ent.IsNotFound(err) {
				return status.Wrap(errors.New("integration not found"), status.NotFound)
			}
			return status.Wrap(fmt.Errorf("unexpected error: %v", err), status.Internal)
		}

		res.ID = integration.ID.String()
		res.Name = integration.Name
		res.Type = integration.Type
		res.Config = integration.Config

		return nil

	}
}

func ListIntegrations(client *ent.Client) func(ctx context.Context, req struct{}, res *[]GetIntegrationResponse) error {
	return func(ctx context.Context, req struct{}, res *[]GetIntegrationResponse) error {
		integrations, err := client.Integrations.Query().All(ctx)
		if err != nil {
			return status.Wrap(errors.New("failed to fetch integrations"), status.Internal)
		}

		*res = make([]GetIntegrationResponse, len(integrations))
		for i, integration := range integrations {
			(*res)[i] = GetIntegrationResponse{
				ID:     integration.ID.String(),
				Name:   integration.Name,
				Type:   integration.Type,
				Config: integration.Config,
			}
		}
		return nil
	}
}
