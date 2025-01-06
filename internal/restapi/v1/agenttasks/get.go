package agenttasks

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/shinobistack/gokakashi/ent"
	"github.com/shinobistack/gokakashi/ent/agenttasks"
	"github.com/swaggest/usecase/status"
	"time"
)

type GetAgentTaskRequest struct {
	ID      uuid.UUID `path:"id"`
	AgentID *int      `path:"agent_id"`
}

type GetAgentTaskResponse struct {
	ID        uuid.UUID `json:"id"`
	AgentID   int       `json:"agent_id"`
	ScanID    uuid.UUID `json:"scan_id"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}
type ListAgentTasksRequest struct {
	AgentID *int       `path:"agent_id"`
	ScanID  *uuid.UUID `query:"scan_id"`
	Status  string     `query:"status"`
}

type ListAgentTasksQueryRequest struct {
	AgentID *int       `query:"agent_id"`
	ScanID  *uuid.UUID `query:"scan_id"`
	Status  string     `query:"status"`
}

type ListAgentTasksByScanIDRequest struct {
	ScanID uuid.UUID `query:"scan_id"`
}

func GetAgentTask(client *ent.Client) func(ctx context.Context, req GetAgentTaskRequest, res *GetAgentTaskResponse) error {
	return func(ctx context.Context, req GetAgentTaskRequest, res *GetAgentTaskResponse) error {
		if req.ID == uuid.Nil {
			return status.Wrap(errors.New("invalid ID"), status.InvalidArgument)
		}

		task, err := client.AgentTasks.Get(ctx, req.ID)
		if err != nil {
			if ent.IsNotFound(err) {
				return status.Wrap(errors.New("agent task not found"), status.NotFound)
			}
			return status.Wrap(err, status.Internal)
		}

		res.ID = task.ID
		res.AgentID = task.AgentID
		res.ScanID = task.ScanID
		res.Status = task.Status
		res.CreatedAt = task.CreatedAt
		return nil
	}
}

func ListAgentTasksByAgentID(client *ent.Client) func(ctx context.Context, req ListAgentTasksRequest, res *[]GetAgentTaskResponse) error {
	return func(ctx context.Context, req ListAgentTasksRequest, res *[]GetAgentTaskResponse) error {
		if req.AgentID == nil || *req.AgentID <= 0 {
			return status.Wrap(errors.New("invalid agent ID"), status.InvalidArgument)
		}

		// Query builder
		query := client.AgentTasks.Query().
			Where(agenttasks.AgentID(*req.AgentID)).
			Order(ent.Asc(agenttasks.FieldCreatedAt)) // Order by created_at ASC

		// Filter by status if provided
		if req.Status != "" {
			query = query.Where(agenttasks.Status((req.Status)))
		}

		tasks, err := query.All(ctx)
		if err != nil {
			return status.Wrap(err, status.Internal)
		}

		*res = make([]GetAgentTaskResponse, len(tasks))
		for i, task := range tasks {
			(*res)[i] = GetAgentTaskResponse{
				ID:        task.ID,
				AgentID:   task.AgentID,
				ScanID:    task.ScanID,
				Status:    task.Status,
				CreatedAt: task.CreatedAt,
			}
		}
		return nil
	}
}

func ListAgentTasks(client *ent.Client) func(ctx context.Context, req ListAgentTasksQueryRequest, res *[]GetAgentTaskResponse) error {
	return func(ctx context.Context, req ListAgentTasksQueryRequest, res *[]GetAgentTaskResponse) error {
		query := client.AgentTasks.Query()

		// Filter by agent ID if provided
		if req.AgentID != nil {
			query = query.Where(agenttasks.AgentID(*req.AgentID))
		}

		// Filter by scan ID if provided
		if req.ScanID != nil && *req.ScanID != uuid.Nil {
			query = query.Where(agenttasks.ScanID(*req.ScanID))
		}

		// Filter by status if provided
		if req.Status != "" {
			query = query.Where(agenttasks.Status(req.Status))
		}

		// Execute query
		tasks, err := query.Order(ent.Asc(agenttasks.FieldCreatedAt)).All(ctx)
		if err != nil {
			return status.Wrap(err, status.Internal)
		}

		// Populate response
		*res = make([]GetAgentTaskResponse, len(tasks))
		for i, task := range tasks {
			(*res)[i] = GetAgentTaskResponse{
				ID:        task.ID,
				AgentID:   task.AgentID,
				ScanID:    task.ScanID,
				Status:    task.Status,
				CreatedAt: task.CreatedAt,
			}
		}
		return nil
	}
}
