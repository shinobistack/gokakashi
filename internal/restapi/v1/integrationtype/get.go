package integrationtype

import (
	"context"
	"errors"
	"github.com/shinobistack/gokakashi/ent"
	"github.com/swaggest/usecase/status"
)

type GetIntegrationTypeRequests struct {
	ID string `path:"id"`
}

type GetIntegrationTypeResponse struct {
	ID          string `json:"id"`
	DisplayName string `json:"display_name"`
}

func GetIntegrationType(client *ent.Client) func(ctx context.Context, req GetIntegrationTypeRequests, res *GetIntegrationTypeResponse) error {
	return func(ctx context.Context, req GetIntegrationTypeRequests, res *GetIntegrationTypeResponse) error {
		it, err := client.IntegrationType.Get(ctx, req.ID)
		if err != nil {
			if ent.IsNotFound(err) {
				return status.Wrap(errors.New("integration type not found"), status.NotFound)
			}
			return status.Wrap(err, status.Internal)
		}

		res.ID = it.ID
		res.DisplayName = it.DisplayName
		return nil

	}

}

func ListIntegrationType(client *ent.Client) func(ctx context.Context, req struct{}, res *[]GetIntegrationTypeResponse) error {
	return func(ctx context.Context, req struct{}, res *[]GetIntegrationTypeResponse) error {
		its, err := client.IntegrationType.Query().All(ctx)
		if err != nil {
			return status.Wrap(errors.New("failed to fetch integration types"), status.Internal)
		}

		*res = make([]GetIntegrationTypeResponse, len(its))
		for i, it := range its {
			(*res)[i] = GetIntegrationTypeResponse{
				ID:          it.ID,
				DisplayName: it.DisplayName,
			}
		}
		return nil
	}
}
