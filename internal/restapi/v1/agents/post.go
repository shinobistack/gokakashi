package agents

import (
	"context"
	"github.com/shinobistack/gokakashi/ent"
	"github.com/swaggest/usecase/status"
)

type CreateAgentRequest struct {
	Status string `json:"status"`
}

type CreateAgentResponse struct {
	ID     int    `json:"id"`
	Status string `json:"status"`
}

func CreateAgent(client *ent.Client) func(ctx context.Context, req CreateAgentRequest, res *CreateAgentResponse) error {
	return func(ctx context.Context, req CreateAgentRequest, res *CreateAgentResponse) error {
		// ToDo: defaults to "connected". Revisit and check what would be more apt
		if req.Status == "" {
			req.Status = "connected"
		}

		agent, err := client.Agents.Create().
			SetStatus(req.Status).
			Save(ctx)
		if err != nil {
			return status.Wrap(err, status.Internal)
		}

		res.ID = agent.ID
		res.Status = agent.Status
		return nil
	}
}
