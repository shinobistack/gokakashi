package policies

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/shinobistack/gokakashi/ent"
	"github.com/swaggest/usecase/status"
)

type DeletePolicyRequest struct {
	ID uuid.UUID `path:"id"`
}

type DeletePolicyResponse struct {
	ID     uuid.UUID `json:"id"`
	Status string    `json:"status"`
}

func DeletePolicy(client *ent.Client) func(ctx context.Context, req DeletePolicyRequest, res *DeletePolicyResponse) error {
	return func(ctx context.Context, req DeletePolicyRequest, res *DeletePolicyResponse) error {
		// Validate ID
		if req.ID == uuid.Nil {
			return status.Wrap(errors.New("invalid UUID: cannot be nil"), status.InvalidArgument)
		}

		// Check if the policy exists
		_, err := client.Policies.Get(ctx, req.ID)
		if err != nil {
			if ent.IsNotFound(err) {
				return status.Wrap(errors.New("policy not found"), status.NotFound)
			}
			return status.Wrap(fmt.Errorf("unexpected error: %v", err), status.Internal)
		}

		// Delete the policy
		err = client.Policies.DeleteOneID(req.ID).Exec(ctx)
		if err != nil {
			return status.Wrap(fmt.Errorf("failed to delete policy: %v", err), status.Internal)
		}

		res.ID = req.ID
		res.Status = "deleted"
		return nil
	}
}
