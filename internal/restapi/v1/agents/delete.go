package agents

import (
	"context"
	"errors"
	"github.com/shinobistack/gokakashi/ent"
	"github.com/swaggest/usecase/status"
)

type DeleteAgentRequest struct {
	ID int `path:"id"`
}

type DeleteAgentResponse struct {
	ID     int    `json:"id"`
	Status string `json:"status"`
}

func DeleteAgent(client *ent.Client) func(ctx context.Context, req DeleteAgentRequest, res *DeleteAgentResponse) error {
	return func(ctx context.Context, req DeleteAgentRequest, res *DeleteAgentResponse) error {
		if req.ID <= 0 {
			return status.Wrap(errors.New("invalid ID"), status.InvalidArgument)
		}

		agent, err := client.Agents.Get(ctx, req.ID)
		if err != nil {
			if ent.IsNotFound(err) {
				return status.Wrap(errors.New("agent not found"), status.NotFound)
			}
			return status.Wrap(err, status.Internal)
		}

		err = client.Agents.DeleteOne(agent).Exec(ctx)
		if err != nil {
			return status.Wrap(err, status.Internal)
		}

		res.ID = agent.ID
		res.Status = "deleted"
		return nil
	}
}
