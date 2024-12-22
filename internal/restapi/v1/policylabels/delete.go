package policylabels

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/shinobistack/gokakashi/ent"
	"github.com/shinobistack/gokakashi/ent/policylabels"
	"github.com/swaggest/usecase/status"
)

type DeletePolicyLabelRequest struct {
	PolicyID uuid.UUID `path:"policy_id"`
	Key      string    `path:"key"`
	Value    string    `query:"value"`
}

type DeletePolicyLabelResponse struct {
	Status string `json:"status"`
}

func DeletePolicyLabel(client *ent.Client) func(ctx context.Context, req DeletePolicyLabelRequest, res *DeletePolicyLabelResponse) error {
	return func(ctx context.Context, req DeletePolicyLabelRequest, res *DeletePolicyLabelResponse) error {
		// Validate inputs
		if req.PolicyID == uuid.Nil {
			return status.Wrap(errors.New("invalid Policy ID: cannot be nil"), status.InvalidArgument)
		}
		if req.Key == "" {
			return status.Wrap(errors.New("invalid Key: cannot be empty"), status.InvalidArgument)
		}
		if req.Value == "" {
			return status.Wrap(errors.New("invalid Value: cannot be empty"), status.InvalidArgument)
		}

		// Check if the label exists
		exists, err := client.PolicyLabels.Query().
			Where(
				policylabels.PolicyID(req.PolicyID),
				policylabels.Key(req.Key),
				policylabels.Value(req.Value),
			).
			Exist(ctx)

		if err != nil {
			return status.Wrap(err, status.Internal)
		}

		if !exists {
			return status.Wrap(errors.New("label not found"), status.NotFound)
		}

		// Delete the label
		_, err = client.PolicyLabels.Delete().
			Where(
				policylabels.PolicyID(req.PolicyID),
				policylabels.Key(req.Key),
				policylabels.Value(req.Value),
			).
			Exec(ctx)

		if err != nil {
			return status.Wrap(err, status.Internal)
		}

		res.Status = "deleted"
		return nil
	}
}
