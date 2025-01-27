package agentlabels

import (
	"context"
	"errors"
	"github.com/shinobistack/gokakashi/ent"
	"github.com/shinobistack/gokakashi/ent/agentlabels"
	"github.com/swaggest/usecase/status"
)

type ListAgentLabelsRequest struct {
	AgentID int      `path:"agent_id"`
	Keys    []string `query:"key"`
	Page    int      `query:"page"`
	Limit   int      `query:"limit"`
}

type ListAgentLabelsResponse struct {
	Labels []AgentLabel `json:"labels"`
}

type AgentLabel struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type GetAgentLabelRequest struct {
	AgentID int    `path:"agent_id"`
	Key     string `path:"key"`
}

type GetAgentLabelResponse struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func ListAgentLabels(client *ent.Client) func(ctx context.Context, req ListAgentLabelsRequest, res *ListAgentLabelsResponse) error {
	return func(ctx context.Context, req ListAgentLabelsRequest, res *ListAgentLabelsResponse) error {
		if req.AgentID < 0 {
			return status.Wrap(errors.New("invalid UUID: cannot be nil"), status.InvalidArgument)
		}
		query := client.AgentLabels.Query().Where(agentlabels.AgentID(req.AgentID))

		// Filter by keys if provided
		if len(req.Keys) > 0 {
			query = query.Where(agentlabels.KeyIn(req.Keys...))
		}

		page := req.Page
		limit := req.Limit
		if page < 1 {
			page = 1
		}
		if limit < 1 || limit > 100 {
			limit = 10
		}
		offset := (page - 1) * limit
		query = query.Offset(offset).Limit(limit)

		labels, err := query.All(ctx)
		if err != nil {
			return status.Wrap(err, status.Internal)
		}

		res.Labels = make([]AgentLabel, len(labels))
		for i, label := range labels {
			res.Labels[i] = AgentLabel{
				Key:   label.Key,
				Value: label.Value,
			}
		}
		return nil
	}
}

func GetAgentLabel(client *ent.Client) func(ctx context.Context, req GetAgentLabelRequest, res *GetAgentLabelResponse) error {
	return func(ctx context.Context, req GetAgentLabelRequest, res *GetAgentLabelResponse) error {
		// Validate inputs
		if req.AgentID < 0 || req.Key == "" {
			return status.Wrap(errors.New("invalid Scan ID or Key"), status.InvalidArgument)
		}

		label, err := client.AgentLabels.Query().
			Where(agentlabels.AgentID(req.AgentID), agentlabels.Key(req.Key)).
			Only(ctx)
		if err != nil {
			return status.Wrap(err, status.Internal)
		}

		res.Key = label.Key
		res.Value = label.Value
		return nil
	}
}
