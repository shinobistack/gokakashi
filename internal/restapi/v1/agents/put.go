package agents

import (
	"context"
	"errors"
	"github.com/shinobistack/gokakashi/ent"
	"github.com/swaggest/usecase/status"
)

type UpdateAgentRequest struct {
	ID     int    `path:"id"`
	Status string `json:"status"`
}

type UpdateAgentResponse struct {
	ID     int    `json:"id"`
	Status string `json:"status"`
}

func UpdateAgent(client *ent.Client) func(ctx context.Context, req UpdateAgentRequest, res *UpdateAgentResponse) error {
	return func(ctx context.Context, req UpdateAgentRequest, res *UpdateAgentResponse) error {
		if req.ID <= 0 || req.Status == "" {
			return status.Wrap(errors.New("invalid ID or Status"), status.InvalidArgument)
		}

		agent, err := client.Agents.UpdateOneID(req.ID).
			SetStatus(req.Status).
			Save(ctx)
		if err != nil {
			if ent.IsNotFound(err) {
				return status.Wrap(errors.New("agent not found"), status.NotFound)
			}
			return status.Wrap(err, status.Internal)
		}

		res.ID = agent.ID
		res.Status = agent.Status
		return nil
	}
}
