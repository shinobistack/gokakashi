package agentlabels

import (
	"context"
	"errors"
	"fmt"
	"github.com/shinobistack/gokakashi/ent"
	"github.com/shinobistack/gokakashi/ent/agentlabels"
	"github.com/swaggest/usecase/status"
)

type UpdateAgentLabelRequest struct {
	AgentID int          `path:"agent_id"`
	Labels  []AgentLabel `json:"labels"`
	Key     *string      `json:"key"`
	Value   *string      `json:"value"`
}

type UpdateAgentLabelResponse struct {
	Labels []AgentLabel `json:"labels"`
}

func UpdateAgentLabel(client *ent.Client) func(ctx context.Context, req UpdateAgentLabelRequest, res *UpdateAgentLabelResponse) error {
	return func(ctx context.Context, req UpdateAgentLabelRequest, res *UpdateAgentLabelResponse) error {
		if req.AgentID < 0 {
			return status.Wrap(errors.New("invalid UUID: cannot be nil"), status.InvalidArgument)
		}

		if req.Key != nil && req.Value != nil {
			// Update a specific label
			_, err := client.AgentLabels.Update().
				Where(agentlabels.AgentID(req.AgentID), agentlabels.Key(*req.Key)).
				SetValue(*req.Value).
				Save(ctx)
			if err != nil {
				if ent.IsNotFound(err) {
					return status.Wrap(errors.New("label not found"), status.NotFound)
				}
				return status.Wrap(err, status.Internal)
			}

			updatedLabel, err := client.AgentLabels.Query().
				Where(agentlabels.AgentID(req.AgentID), agentlabels.Key(*req.Key)).
				Only(ctx)
			if err != nil {
				if ent.IsNotFound(err) {
					return status.Wrap(errors.New("label not found after update"), status.NotFound)
				}
				return status.Wrap(err, status.Internal)
			}

			// Return the updated label
			res.Labels = []AgentLabel{
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
			_, err = tx.AgentLabels.Delete().
				Where(agentlabels.AgentID(req.AgentID)).
				Exec(ctx)
			if err != nil {
				if rollbackErr := tx.Rollback(); rollbackErr != nil {
					fmt.Printf("rollback failed: %v\n", rollbackErr)
				}
				return status.Wrap(err, status.Internal)
			}

			// Add new labels
			bulk := make([]*ent.AgentLabelsCreate, len(req.Labels))
			for i, label := range req.Labels {
				bulk[i] = tx.AgentLabels.Create().
					SetAgentID(req.AgentID).
					SetKey(label.Key).
					SetValue(label.Value)
			}

			createdLabels, err := tx.AgentLabels.CreateBulk(bulk...).Save(ctx)
			if err != nil {
				if rollbackErr := tx.Rollback(); rollbackErr != nil {
					fmt.Printf("rollback failed: %v\n", rollbackErr)
				}
				return status.Wrap(err, status.Internal)
			}

			err = tx.Commit()
			if err != nil {
				return status.Wrap(err, status.Internal)
			}

			// Build the response
			res.Labels = make([]AgentLabel, len(createdLabels))
			for i, label := range createdLabels {
				res.Labels[i] = AgentLabel{
					Key:   label.Key,
					Value: label.Value,
				}
			}
			return nil
		}

		return status.Wrap(errors.New("invalid request: either key-value or labels must be provided"), status.InvalidArgument)
	}
}
