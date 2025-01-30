package agentlabels

import (
	"context"
	"errors"
	"github.com/shinobistack/gokakashi/ent"
	"github.com/shinobistack/gokakashi/ent/agentlabels"
	"github.com/swaggest/usecase/status"
)

// ToDo: To have a requests to create labels in bulk?
type CreateAgentLabelRequest struct {
	AgentID int    `path:"agent_id"`
	Key     string `json:"key"`
	Value   string `json:"value"`
}

type CreateAgentLabelResponse struct {
	AgentID int    `path:"agent_id"`
	Key     string `json:"key"`
	Value   string `json:"value"`
}

func CreateAgentLabel(client *ent.Client) func(ctx context.Context, req CreateAgentLabelRequest, res *CreateAgentLabelResponse) error {
	return func(ctx context.Context, req CreateAgentLabelRequest, res *CreateAgentLabelResponse) error {
		if req.AgentID < 0 || req.Key == "" || req.Value == "" {
			return status.Wrap(errors.New("invalid input: missing fields"), status.InvalidArgument)
		}

		// Check if the label already exists
		exists, err := client.AgentLabels.Query().
			Where(agentlabels.AgentID(req.AgentID), agentlabels.Key(req.Key)).
			Exist(ctx)
		if err != nil {
			return status.Wrap(err, status.Internal)
		}
		if exists {
			return status.Wrap(errors.New("label already exists"), status.AlreadyExists)
		}

		label, err := client.AgentLabels.Create().
			SetAgentID(req.AgentID).
			SetKey(req.Key).
			SetValue(req.Value).
			Save(ctx)
		if err != nil {
			return status.Wrap(err, status.Internal)
		}

		res.AgentID = label.AgentID
		res.Key = label.Key
		res.Value = label.Value

		return nil
	}
}
