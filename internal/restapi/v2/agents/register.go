package agents

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/shinobistack/gokakashi/ent"
	"github.com/shinobistack/gokakashi/internal/agent"
	"github.com/swaggest/usecase/status"
)

type RegisterAgentRequest struct{}

type RegisterAgentResponse struct {
	ID uuid.UUID `json:"id"`

	Status          agent.Status `json:"status"`
	LastHeartbeatAt time.Time    `json:"last_heartbeat_at"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func Register(client *ent.Client) func(ctx context.Context, req RegisterAgentRequest, res *RegisterAgentResponse) error {
	return func(ctx context.Context, req RegisterAgentRequest, res *RegisterAgentResponse) error {
		newAgent, err := client.V2Agents.Create().
			SetStatus(string(agent.Disconnected)).
			Save(ctx)
		if err != nil {
			return status.Wrap(err, status.Internal)
		}

		*res = RegisterAgentResponse{
			ID:              newAgent.ID,
			Status:          agent.Status(newAgent.Status),
			LastHeartbeatAt: newAgent.LastHeartbeatAt,
			CreatedAt:       newAgent.CreatedAt,
			UpdatedAt:       newAgent.UpdatedAt,
		}
		return nil
	}
}
