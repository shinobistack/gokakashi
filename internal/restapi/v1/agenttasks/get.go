package agenttasks

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/shinobistack/gokakashi/ent"
	"github.com/swaggest/usecase/status"
	"time"
)

type GetAgentTaskRequest struct {
	ID int `path:"id"`
}

type GetAgentTaskResponse struct {
	ID        int       `json:"id"`
	AgentID   int       `json:"agent_id"`
	ScanID    uuid.UUID `json:"scan_id"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

type ListAgentTasksResponse struct {
	AgentTasks []GetAgentTaskResponse `json:"agent_tasks"`
}

func ListAgentTasks(client *ent.Client) func(ctx context.Context, req interface{}, res *[]GetAgentTaskResponse) error {
	return func(ctx context.Context, req interface{}, res *[]GetAgentTaskResponse) error {
		tasks, err := client.AgentTasks.Query().All(ctx)
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

func GetAgentTask(client *ent.Client) func(ctx context.Context, req GetAgentTaskRequest, res *GetAgentTaskResponse) error {
	return func(ctx context.Context, req GetAgentTaskRequest, res *GetAgentTaskResponse) error {
		if req.ID <= 0 {
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
