package policylabels

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/shinobistack/gokakashi/ent"
	"github.com/shinobistack/gokakashi/ent/policylabels"
	"github.com/swaggest/usecase/status"
)

type UpdatePolicyLabelsRequest struct {
	PolicyID uuid.UUID     `path:"policy_id"`
	Labels   []PolicyLabel `json:"labels"`
	Key      *string       `json:"key"`
	Value    *string       `json:"value"`
}

type UpdatePolicyLabelsResponse struct {
	Labels []PolicyLabel `json:"labels"`
}

func UpdatePolicyLabels(client *ent.Client) func(ctx context.Context, req UpdatePolicyLabelsRequest, res *UpdatePolicyLabelsResponse) error {
	return func(ctx context.Context, req UpdatePolicyLabelsRequest, res *UpdatePolicyLabelsResponse) error {
		// Validate PolicyID
		if req.PolicyID == uuid.Nil {
			return status.Wrap(errors.New("invalid Policy ID"), status.InvalidArgument)
		}

		if req.Key != nil && req.Value != nil {
			// Update a specific label
			_, err := client.PolicyLabels.Update().
				Where(policylabels.PolicyID(req.PolicyID), policylabels.Key(*req.Key)).
				SetValue(*req.Value).
				Save(ctx)
			if err != nil {
				if ent.IsNotFound(err) {
					return status.Wrap(errors.New("label not found"), status.NotFound)
				}
				return status.Wrap(err, status.Internal)
			}

			updatedLabel, err := client.PolicyLabels.Query().
				Where(policylabels.PolicyID(req.PolicyID), policylabels.Key(*req.Key)).
				Only(ctx)
			if err != nil {
				if ent.IsNotFound(err) {
					return status.Wrap(errors.New("label not found after update"), status.NotFound)
				}
				return status.Wrap(err, status.Internal)
			}

			// Return the updated label
			res.Labels = []PolicyLabel{
				{
					Key:   updatedLabel.Key,
					Value: updatedLabel.Value,
				},
			}
			return nil
		}

		if req.Labels != nil {
			// Replace all labels
			tx, err := client.Tx(ctx)
			if err != nil {
				return status.Wrap(err, status.Internal)
			}

			// Delete existing labels
			_, err = tx.PolicyLabels.Delete().
				Where(policylabels.PolicyID(req.PolicyID)).
				Exec(ctx)
			if err != nil {
				if rollbackErr := tx.Rollback(); rollbackErr != nil {
					fmt.Printf("rollback failed: %v\n", rollbackErr)
				}
				return status.Wrap(err, status.Internal)
			}

			// Add new labels
			bulk := make([]*ent.PolicyLabelsCreate, len(req.Labels))
			for i, label := range req.Labels {
				bulk[i] = tx.PolicyLabels.Create().
					SetPolicyID(req.PolicyID).
					SetKey(label.Key).
					SetValue(label.Value)
			}

			createdLabels, err := tx.PolicyLabels.CreateBulk(bulk...).Save(ctx)
			if err != nil {
				tx.Rollback()
				return status.Wrap(err, status.Internal)
			}

			err = tx.Commit()
			if err != nil {
				return status.Wrap(err, status.Internal)
			}

			// Build the response
			res.Labels = make([]PolicyLabel, len(createdLabels))
			for i, label := range createdLabels {
				res.Labels[i] = PolicyLabel{
					Key:   label.Key,
					Value: label.Value,
				}
			}
			return nil
		}

		return status.Wrap(errors.New("invalid request: either key-value or labels must be provided"), status.InvalidArgument)
	}
}
