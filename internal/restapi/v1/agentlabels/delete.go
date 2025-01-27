package agentlabels

import (
	"context"
	"errors"
	"github.com/shinobistack/gokakashi/ent"
	"github.com/shinobistack/gokakashi/ent/agentlabels"
	"github.com/swaggest/usecase/status"
)

type DeleteAgentLabelRequest struct {
	AgentID int    `path:"agent_id"`
	Key     string `path:"key"`
}

type DeleteAgentLabelResponse struct {
	Status string `json:"status"`
}

func DeleteAgentLabel(client *ent.Client) func(ctx context.Context, req DeleteAgentLabelRequest, res *DeleteAgentLabelResponse) error {
	return func(ctx context.Context, req DeleteAgentLabelRequest, res *DeleteAgentLabelResponse) error {
		if req.AgentID < 0 {
			return status.Wrap(errors.New("invalid UUID: cannot be nil"), status.InvalidArgument)
		}
		if req.Key == "" {
			return status.Wrap(errors.New("invalid key: cannot be nil"), status.InvalidArgument)
		}
		// Check if the label exists
		exists, err := client.AgentLabels.Query().
			Where(agentlabels.AgentID(req.AgentID), agentlabels.Key(req.Key)).
			Exist(ctx)
		if err != nil {
			return status.Wrap(err, status.Internal)
		}
		if !exists {
			return status.Wrap(errors.New("label not found"), status.NotFound)
		}

		// Delete the label
		_, err = client.AgentLabels.Delete().
			Where(agentlabels.AgentID(req.AgentID), agentlabels.Key(req.Key)).
			Exec(ctx)
		if err != nil {
			return status.Wrap(err, status.Internal)
		}

		res.Status = "deleted"
		return nil
	}
}
