package agents

import (
	"context"

	"github.com/shinobistack/gokakashi/ent"
	"github.com/shinobistack/gokakashi/internal/agent"
	"github.com/shinobistack/gokakashi/internal/restapi/v2/io"
	"github.com/swaggest/usecase/status"
)

func Register(client *ent.Client) func(ctx context.Context, req io.RegisterAgentRequest, res *io.RegisterAgentResponse) error {
	return func(ctx context.Context, req io.RegisterAgentRequest, res *io.RegisterAgentResponse) error {
		newAgent, err := client.V2Agents.Create().
			SetStatus(string(agent.Disconnected)).
			Save(ctx)
		if err != nil {
			return status.Wrap(err, status.Internal)
		}

		*res = io.RegisterAgentResponse{
			ID:              newAgent.ID,
			Status:          agent.Status(newAgent.Status),
			LastHeartbeatAt: newAgent.LastHeartbeatAt,
			CreatedAt:       newAgent.CreatedAt,
			UpdatedAt:       newAgent.UpdatedAt,
		}
		return nil
	}
}
