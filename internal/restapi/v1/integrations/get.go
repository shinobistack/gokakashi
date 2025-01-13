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

type GetIntegrationRequests struct {
	ID uuid.UUID `path:"id"`
}

type GetIntegrationResponse struct {
	ID     uuid.UUID              `json:"id"`
	Name   string                 `json:"name"`
	Type   string                 `json:"type"`
	Config map[string]interface{} `json:"config"`
}

type ListGetIntegrationRequests struct {
	Name string `query:"name"`
}

type ListIntegrationResponse struct {
	Integrations []GetIntegrationResponse `json:"integrations"`
}

func GetIntegration(client *ent.Client) func(ctx context.Context, req GetIntegrationRequests, res *GetIntegrationResponse) error {
	return func(ctx context.Context, req GetIntegrationRequests, res *GetIntegrationResponse) error {
		// Validate ID
		if req.ID == uuid.Nil {
			return status.Wrap(errors.New("invalid UUID: cannot be nil"), status.InvalidArgument)
		}

		// Fetch integration by ID
		integration, err := client.Integrations.Get(ctx, req.ID)
		if err != nil {
			if ent.IsNotFound(err) {
				return status.Wrap(errors.New("integration not found"), status.NotFound)
			}
			return status.Wrap(fmt.Errorf("unexpected error: %v", err), status.Internal)
		}

		res.ID = integration.ID
		res.Name = integration.Name
		res.Type = integration.Type
		res.Config = integration.Config

		return nil

	}
}

func ListIntegrations(client *ent.Client) func(ctx context.Context, req ListGetIntegrationRequests, res *[]GetIntegrationResponse) error {
	return func(ctx context.Context, req ListGetIntegrationRequests, res *[]GetIntegrationResponse) error {
		query := client.Integrations.Query()
		if req.Name != "" {
			query = query.Where(integrations.Name(req.Name))
		}

		integrations, err := query.All(ctx)
		if err != nil {
			return status.Wrap(errors.New("failed to fetch integrations"), status.Internal)
		}

		*res = make([]GetIntegrationResponse, len(integrations))
		for i, integration := range integrations {
			(*res)[i] = GetIntegrationResponse{
				ID:     integration.ID,
				Name:   integration.Name,
				Type:   integration.Type,
				Config: integration.Config,
			}
		}
		return nil
	}
}
