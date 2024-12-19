package integrations

import (
	"context"
	"fmt"
	"github.com/shinobistack/gokakashi/ent"
	"github.com/swaggest/usecase/status"
)

type CreateIntegrationRequest struct {
	Name   string                 `json:"name"`
	Type   string                 `json:"type"`
	Config map[string]interface{} `json:"config"`
}

type CreateIntegrationResponse struct {
	ID     string `json:"id"`
	Status string `json:"status"`
}

func CreateIntegration(client *ent.Client) func(ctx context.Context, req CreateIntegrationRequest, res *CreateIntegrationResponse) error {
	return func(ctx context.Context, req CreateIntegrationRequest, res *CreateIntegrationResponse) error {
		integration, err := client.Integrations.Create().
			SetName(req.Name).
			SetType(req.Type).
			SetConfig(req.Config).
			Save(ctx)
		if err != nil {
			return status.Wrap(fmt.Errorf("failed to create integration: %v", err), status.Internal)
		}

		res.ID = integration.ID.String()
		res.Status = "created"
		return nil

	}

}
