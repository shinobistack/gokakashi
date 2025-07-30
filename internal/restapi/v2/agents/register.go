package agents

import (
	"context"

	"github.com/shinobistack/gokakashi/ent"
	agent "github.com/shinobistack/gokakashi/internal/agent/status/v2"
	"github.com/shinobistack/gokakashi/internal/restapi/v2/io"
	"github.com/swaggest/usecase/status"
)

func Register(client *ent.Client) func(ctx context.Context, req io.AgentRegisterRequest, res *io.AgentRegisterResponse) error {
	return func(ctx context.Context, req io.AgentRegisterRequest, res *io.AgentRegisterResponse) error {
		newAgent, err := client.V2Agents.Create().
			SetStatus(string(agent.Disconnected)).
			Save(ctx)
		if err != nil {
			return status.Wrap(err, status.Internal)
		}

		*res = io.AgentRegisterResponse{
			ID:              newAgent.ID,
			Status:          agent.Status(newAgent.Status),
			LastHeartbeatAt: newAgent.LastHeartbeatAt,
			CreatedAt:       newAgent.CreatedAt,
			UpdatedAt:       newAgent.UpdatedAt,
		}
		return nil
	}
}
