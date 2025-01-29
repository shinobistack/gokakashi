package agents

import (
	"context"
	"errors"
	"fmt"
	"github.com/shinobistack/gokakashi/ent"
	"github.com/shinobistack/gokakashi/ent/agents"
	"github.com/shinobistack/gokakashi/ent/schema"
	"github.com/swaggest/usecase/status"
	"time"
)

type GetAgentRequest struct {
	ID int `path:"id"`
}

type ListAgentsRequest struct{}

type GetAgentResponse struct {
	ID            int                   `json:"id"`
	Status        string                `json:"status"`
	Name          string                `json:"name,omitempty"`
	Labels        []schema.CommonLabels `json:"labels,omitempty"`
	LastHeartbeat time.Time             `json:"last_heartbeat"`
}

type ListAgentsResponse struct {
	Agents []GetAgentResponse `json:"agents"`
}

type PollAgentsRequest struct {
	Status string `query:"status"`
	Name   string `query:"name"`
}

type PollAgentsResponse struct {
	ID        int                   `json:"id"`
	Status    string                `json:"status"`
	LastSeen  string                `json:"last_seen"`
	Workspace string                `json:"workspace"`
	Name      string                `json:"name"`
	Server    string                `json:"server"`
	Labels    []schema.CommonLabels `json:"labels,omitempty"`
}

// ToDo: to remove list agents and write test case for pollagents

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

		// Fetch the agent by ID with its labels
		agent, err := client.Agents.Query().
			Where(agents.ID(req.ID)).
			WithAgentLabels().
			Only(ctx)
		if err != nil {
			if ent.IsNotFound(err) {
				return status.Wrap(errors.New("agent not found"), status.NotFound)
			}
			return status.Wrap(fmt.Errorf("unexpected error: %v", err), status.Internal)
		}

		labels := mapAgentLabels(agent.Edges.AgentLabels)

		*res = GetAgentResponse{
			ID:            agent.ID,
			Status:        agent.Status,
			Name:          agent.Name,
			Labels:        labels,
			LastHeartbeat: agent.LastHeartbeat,
		}
		return nil
	}
}

func PollAgents(client *ent.Client) func(ctx context.Context, req PollAgentsRequest, res *[]PollAgentsResponse) error {
	return func(ctx context.Context, req PollAgentsRequest, res *[]PollAgentsResponse) error {
		query := client.Agents.Query().WithAgentLabels()

		if req.Status != "" {
			query = query.Where(agents.Status(req.Status))
		}

		agentsList, err := query.All(ctx)
		if err != nil {
			return status.Wrap(err, status.Internal)
		}

		*res = make([]PollAgentsResponse, len(agentsList))
		for i, agent := range agentsList {
			labels := mapAgentLabels(agent.Edges.AgentLabels)

			(*res)[i] = PollAgentsResponse{
				ID:        agent.ID,
				Status:    agent.Status,
				LastSeen:  agent.LastSeen.Format(time.RFC3339), // ToDo: which format is needed?
				Workspace: agent.Workspace,
				Name:      agent.Name,
				Server:    agent.Server,
				Labels:    labels,
			}
		}
		return nil
	}
}

func mapAgentLabels(labels []*ent.AgentLabels) []schema.CommonLabels {
	if labels == nil {
		return nil
	}

	mapped := make([]schema.CommonLabels, len(labels))
	for i, label := range labels {
		mapped[i] = schema.CommonLabels{
			Key:   label.Key,
			Value: label.Value,
		}
	}
	return mapped
}
