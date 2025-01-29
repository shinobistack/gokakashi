package agents

import (
	"context"
	"errors"
	"github.com/shinobistack/gokakashi/ent"
	"github.com/swaggest/usecase/status"
	"log"
	"time"
)

type UpdateAgentRequest struct {
	ID     int    `path:"id"`
	Status string `json:"status"`
}

type UpdateAgentResponse struct {
	ID     int    `json:"id"`
	Status string `json:"status"`
}

type UpdateAgentHeartbeatRequest struct {
	ID int `path:"id"`
}

type UpdateAgentHeartbeatResponse struct{}

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

func UpdateAgentHeartbeat(client *ent.Client) func(ctx context.Context, req UpdateAgentHeartbeatRequest, res *UpdateAgentHeartbeatResponse) error {
	return func(ctx context.Context, req UpdateAgentHeartbeatRequest, res *UpdateAgentHeartbeatResponse) error {
		if req.ID <= 0 {
			return status.Wrap(errors.New("invalid ID"), status.InvalidArgument)
		}
		agent, err := client.Agents.UpdateOneID(req.ID).
			SetLastHeartbeat(time.Now()).
			Save(ctx)
		if err != nil {
			if ent.IsNotFound(err) {
				return status.Wrap(errors.New("agent not found"), status.NotFound)
			}
			log.Printf("Failed to update heartbeat for agent ID %d: %v", req.ID, err)
			return status.Wrap(err, status.Internal)
		}
		log.Printf("API: Successfully updated heartbeat for agent ID %d at %v", agent.ID, agent.LastHeartbeat.Format(time.RFC3339))
		return nil

	}
}
