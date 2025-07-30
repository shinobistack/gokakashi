package agents

import (
	"context"
	"time"

	"github.com/shinobistack/gokakashi/ent"
	agentstatus "github.com/shinobistack/gokakashi/internal/agent/status/v2"
	"github.com/shinobistack/gokakashi/internal/restapi/v2/io"
	"github.com/swaggest/usecase/status"
)

func Heartbeat(client *ent.Client) func(ctx context.Context, req io.AgentHeartbeatRequest, res *io.AgentHeartbeatResponse) error {
	return func(ctx context.Context, req io.AgentHeartbeatRequest, res *io.AgentHeartbeatResponse) error {
		agent, err := client.V2Agents.UpdateOneID(req.ID).
			SetLastHeartbeatAt(time.Now()).
			Save(ctx)
		if err != nil {
			return status.Wrap(err, status.Internal)
		}

		*res = io.AgentHeartbeatResponse{
			ID:              agent.ID,
			Status:          agentstatus.Status(agent.Status),
			LastHeartbeatAt: agent.LastHeartbeatAt,
			CreatedAt:       agent.CreatedAt,
			UpdatedAt:       agent.UpdatedAt,
		}
		return nil
	}
}
