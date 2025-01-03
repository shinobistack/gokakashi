package agents

import (
	"context"
	"errors"
	"github.com/shinobistack/gokakashi/ent"
	"github.com/shinobistack/gokakashi/ent/agents"
	"github.com/swaggest/usecase/status"
	"time"
)

type CreateAgentRequest struct {
	Status string `json:"status"`
}

type CreateAgentResponse struct {
	ID     int    `json:"id"`
	Status string `json:"status"`
}
type RegisterAgentRequest struct {
	Server    string `json:"server"`
	Token     string `json:"token"`
	Workspace string `json:"workspace"`
	Status    string `json:"status"`
	Name      string `json:"name,omitempty"`
}

type RegisterAgentResponse struct {
	ID        int    `json:"id"`
	Status    string `json:"status"`
	LastSeen  string `json:"last_seen"`
	Workspace string `json:"workspace"`
	Name      string `json:"name"`
	Server    string `json:"server"`
}

// ToDo: To remove this function and write test case for RegisterAgent
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

func RegisterAgent(client *ent.Client) func(ctx context.Context, req RegisterAgentRequest, res *RegisterAgentResponse) error {
	return func(ctx context.Context, req RegisterAgentRequest, res *RegisterAgentResponse) error {
		// Validate input
		if req.Server == "" || req.Workspace == "" {
			return status.Wrap(errors.New("missing required fields"), status.InvalidArgument)
		}

		// Check if the agent already exists
		existingAgent, err := client.Agents.Query().
			Where(
				agents.Server(req.Server),
				agents.Name(req.Name),
			).Only(ctx)

		if err == nil && existingAgent != nil {
			// Update existing agent
			updatedAgent, err := client.Agents.UpdateOne(existingAgent).
				SetWorkspace(req.Workspace).
				SetLastSeen(time.Now()).
				SetStatus("connected").
				Save(ctx)
			if err != nil {
				return status.Wrap(err, status.Internal)
			}

			*res = RegisterAgentResponse{
				ID:        updatedAgent.ID,
				Status:    updatedAgent.Status,
				LastSeen:  updatedAgent.LastSeen.Format(time.RFC3339),
				Workspace: updatedAgent.Workspace,
				Name:      updatedAgent.Name,
				Server:    updatedAgent.Server,
			}
			return nil
		} else if !ent.IsNotFound(err) {
			return status.Wrap(err, status.Internal)
		}

		// Register new agent
		newAgent, err := client.Agents.Create().
			SetName(req.Name).
			SetServer(req.Server).
			SetWorkspace(req.Workspace).
			SetStatus("connected").
			SetLastSeen(time.Now()).
			Save(ctx)
		if err != nil {
			return status.Wrap(err, status.Internal)
		}

		*res = RegisterAgentResponse{
			ID:        newAgent.ID,
			Status:    newAgent.Status,
			LastSeen:  newAgent.LastSeen.Format(time.RFC3339),
			Workspace: newAgent.Workspace,
			Name:      newAgent.Name,
			Server:    newAgent.Server,
		}
		return nil
	}
}
