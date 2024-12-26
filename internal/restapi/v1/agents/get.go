package agents

import (
	"context"
	"errors"
	"github.com/shinobistack/gokakashi/ent"
	"github.com/swaggest/usecase/status"
)

type GetAgentRequest struct {
	ID int `path:"id"`
}

type ListAgentsRequest struct{}

type GetAgentResponse struct {
	ID     int    `json:"id"`
	Status string `json:"status"`
}

type ListAgentsResponse struct {
	Agents []GetAgentResponse `json:"agents"`
}

func ListAgents(client *ent.Client) func(ctx context.Context, req ListAgentsRequest, res *[]GetAgentResponse) error {
	return func(ctx context.Context, req ListAgentsRequest, res *[]GetAgentResponse) error {
		agentsList, err := client.Agents.Query().All(ctx)
		if err != nil {
			return status.Wrap(err, status.Internal)
		}

		*res = make([]GetAgentResponse, len(agentsList))
		for i, agent := range agentsList {
			(*res)[i] = GetAgentResponse{
				ID:     agent.ID,
				Status: agent.Status,
			}
		}
		return nil
	}
}

func GetAgent(client *ent.Client) func(ctx context.Context, req GetAgentRequest, res *GetAgentResponse) error {
	return func(ctx context.Context, req GetAgentRequest, res *GetAgentResponse) error {
		// ToDo: As of now agentID is not an UUID. Hence <=0. Revisit and see what would be better
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

		res.ID = agent.ID
		res.Status = agent.Status
		return nil
	}
}
