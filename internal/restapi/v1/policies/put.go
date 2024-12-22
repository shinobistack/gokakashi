package policies

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/shinobistack/gokakashi/ent"
	"github.com/shinobistack/gokakashi/ent/schema"
	"github.com/swaggest/usecase/status"
)

type UpdatePolicyRequest struct {
	ID      uuid.UUID               `path:"id"`
	Name    *string                 `json:"name"`
	Image   *schema.Image           `json:"image"`
	Trigger *map[string]interface{} `json:"trigger"`
	Check   *schema.Check           `json:"check"`
}

type UpdatePolicyResponse struct {
	ID     uuid.UUID `json:"id"`
	Status string    `json:"status"`
}

func UpdatePolicy(client *ent.Client) func(ctx context.Context, req UpdatePolicyRequest, res *GetPolicyResponse) error {
	return func(ctx context.Context, req UpdatePolicyRequest, res *GetPolicyResponse) error {
		// Validate ID
		if req.ID == uuid.Nil {
			return status.Wrap(errors.New("invalid UUID: cannot be nil"), status.InvalidArgument)
		}

		// Fetch the policy
		policy, err := client.Policies.Get(ctx, req.ID)
		if err != nil {
			if ent.IsNotFound(err) {
				return status.Wrap(errors.New("policy not found"), status.NotFound)
			}
			return status.Wrap(fmt.Errorf("unexpected error: %v", err), status.Internal)
		}

		// Start updating fields if provided
		// ToDo: Allows to update the objects same as present in db. Not required, see whats efficient.
		update := client.Policies.UpdateOne(policy)

		if req.Name != nil {
			update.SetName(*req.Name)
		}

		if req.Image != nil {
			update.SetImage(*req.Image)
		}

		if req.Trigger != nil {
			update.SetTrigger(*req.Trigger)
		}

		if req.Check != nil {
			update.SetCheck(*req.Check)
		}

		// Save updates
		updatedPolicy, err := update.Save(ctx)
		if err != nil {
			return status.Wrap(fmt.Errorf("failed to update policy: %v", err), status.Internal)
		}

		res.ID = updatedPolicy.ID
		res.Name = policy.Name
		res.Image = policy.Image
		res.Trigger = policy.Trigger
		res.Check = convertToPointer(policy.Check)
		return nil
	}
}
