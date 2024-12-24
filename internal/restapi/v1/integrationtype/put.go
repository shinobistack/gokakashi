package integrationtype

import (
	"context"
	"errors"
	"fmt"
	"github.com/shinobistack/gokakashi/ent"
	"github.com/swaggest/usecase/status"
)

type UpdateIntegrationTypeRequest struct {
	ID          string  `path:"id"`
	DisplayName *string `json:"display_name"`
}

func UpdateIntegrationType(client *ent.Client) func(ctx context.Context, req UpdateIntegrationTypeRequest, res *GetIntegrationTypeResponse) error {
	return func(ctx context.Context, req UpdateIntegrationTypeRequest, res *GetIntegrationTypeResponse) error {
		if !isValidID(req.ID) {
			return status.Wrap(errors.New("invalid id format"), status.InvalidArgument) // 400 Bad Request
		}

		update := client.IntegrationType.UpdateOneID(req.ID)
		if req.DisplayName != nil {
			update = update.SetDisplayName(*req.DisplayName)
		}

		it, err := update.Save(ctx)
		if err != nil {
			if ent.IsNotFound(err) {
				return status.Wrap(errors.New("integration type not found"), status.NotFound)
			}
			return status.Wrap(fmt.Errorf("failed to update integration type: %v", err), status.Internal)
		}

		res.ID = it.ID
		res.DisplayName = it.DisplayName
		return nil
	}
}
