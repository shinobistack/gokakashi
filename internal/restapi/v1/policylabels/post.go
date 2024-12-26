package policylabels

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/shinobistack/gokakashi/ent"
	"github.com/shinobistack/gokakashi/ent/policylabels"
	"github.com/swaggest/usecase/status"
)

// ToDo: To have a requests to create labels in bulk?
type CreatePolicyLabelRequest struct {
	PolicyID uuid.UUID `path:"policy_id"`
	Key      string    `json:"key"`
	Value    string    `json:"value"`
}

type CreatePolicyLabelResponse struct {
	ID     uuid.UUID     `json:"policy_id"`
	Labels []PolicyLabel `json:"labels"`
}

// ToDO: When extra fields are passed?
func CreatePolicyLabel(client *ent.Client) func(ctx context.Context, req CreatePolicyLabelRequest, res *CreatePolicyLabelResponse) error {
	return func(ctx context.Context, req CreatePolicyLabelRequest, res *CreatePolicyLabelResponse) error {
		// Validate inputs
		if req.PolicyID == uuid.Nil {
			return status.Wrap(errors.New("invalid Policy ID"), status.InvalidArgument)
		}
		if req.Key == "" || req.Value == "" {
			return status.Wrap(errors.New("key and value must not be empty"), status.InvalidArgument)
		}

		// Ensure the policy exists
		_, err := client.Policies.Get(ctx, req.PolicyID)
		if err != nil {
			if ent.IsNotFound(err) {
				return status.Wrap(errors.New("policy not found"), status.NotFound)
			}
			return status.Wrap(err, status.Internal)
		}
		// Validate of if label exists
		exists, _ := client.PolicyLabels.Query().
			Where(policylabels.PolicyID(req.PolicyID), policylabels.Key(req.Key)).
			Exist(ctx)
		if err != nil {
			return status.Wrap(err, status.Internal)
		}
		if exists {
			return status.Wrap(errors.New("label with the same key already exists"), status.AlreadyExists)
		}

		// Create the policy label
		// ToDo: append the key
		label, err := client.PolicyLabels.Create().
			SetPolicyID(req.PolicyID).
			SetKey(req.Key).
			SetValue(req.Value).
			Save(ctx)
		if err != nil {
			return status.Wrap(err, status.Internal)
		}

		res.ID = label.PolicyID
		res.Labels = []PolicyLabel{
			{
				Key:   label.Key,
				Value: label.Value,
			},
		}
		return nil
	}
}
