package agents

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/shinobistack/gokakashi/ent"
	"github.com/shinobistack/gokakashi/ent/agents"
	"github.com/shinobistack/gokakashi/ent/schema"
	"github.com/swaggest/usecase/status"
	"log"
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
	Server    string                `json:"server"`
	Token     string                `json:"token"`
	Workspace string                `json:"workspace"`
	Status    string                `json:"status"`
	Name      string                `json:"name,omitempty"`
	Labels    []schema.CommonLabels `json:"labels,omitempty"`
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
		if req.Server == "" {
			return status.Wrap(errors.New("missing required fields"), status.InvalidArgument)
		}

		if req.Name == "" {
			req.Name = fmt.Sprintf("agent-%s", uuid.New().String()[:6])
		}
		if req.Workspace == "" {
			req.Workspace = fmt.Sprintf("/tmp/%s", req.Name)
		}

		// Check if the agent already exists
		existingAgent, err := client.Agents.Query().
			Where(
				agents.Server(req.Server),
				agents.Name(req.Name),
			).Only(ctx)

		if err == nil && existingAgent != nil {
			return status.Wrap(errors.New("agent with the same name already exists on this server"), status.AlreadyExists)
		} else if !ent.IsNotFound(err) {
			return status.Wrap(err, status.Internal)
		}

		// Enforce workspace uniqueness
		workspaceConflict, err := client.Agents.Query().
			Where(
				agents.Workspace(req.Workspace),
			).Only(ctx)
		if err == nil && workspaceConflict != nil {
			return status.Wrap(errors.New("workspace is already in use by another agent"), status.InvalidArgument)
		} else if !ent.IsNotFound(err) {
			return status.Wrap(err, status.Internal)
		}

		tx, err := client.Tx(ctx)
		if err != nil {
			return status.Wrap(err, status.Internal)
		}

		// Register the new agent
		newAgent, err := tx.Agents.Create().
			SetName(req.Name).
			SetServer(req.Server).
			SetWorkspace(req.Workspace).
			SetStatus("connected").
			SetLastSeen(time.Now()).
			SetLastHeartbeat(time.Now()).
			Save(ctx)
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				log.Printf("rollback failed: %v\n", rollbackErr)
			}
			return status.Wrap(err, status.Internal)
		}

		// Save agent labels
		if len(req.Labels) > 0 {
			bulk := make([]*ent.AgentLabelsCreate, len(req.Labels))
			for i, label := range req.Labels {
				bulk[i] = tx.AgentLabels.Create().
					SetAgentsID(newAgent.ID).
					SetKey(label.Key).
					SetValue(label.Value)
			}

			if _, err := tx.AgentLabels.CreateBulk(bulk...).Save(ctx); err != nil {
				if rollbackErr := tx.Rollback(); rollbackErr != nil {
					log.Printf("rollback failed: %v\n", rollbackErr)
				}
				return status.Wrap(err, status.Internal)
			}
		}

		if err := tx.Commit(); err != nil {
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
