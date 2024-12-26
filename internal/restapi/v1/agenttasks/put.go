package agenttasks

import (
	"context"
	"errors"
	"github.com/shinobistack/gokakashi/ent"
	"github.com/swaggest/usecase/status"
)

type UpdateAgentTaskRequest struct {
	ID     int    `path:"id"`
	Status string `json:"status"`
}

type UpdateAgentTaskResponse struct {
	ID     int    `json:"id"`
	Status string `json:"status"`
}

func UpdateAgentTask(client *ent.Client) func(ctx context.Context, req UpdateAgentTaskRequest, res *UpdateAgentTaskResponse) error {
	return func(ctx context.Context, req UpdateAgentTaskRequest, res *UpdateAgentTaskResponse) error {
		if req.ID <= 0 || req.Status == "" {
			return status.Wrap(errors.New("invalid ID or Status"), status.InvalidArgument)
		}

		task, err := client.AgentTasks.UpdateOneID(req.ID).
			SetStatus(req.Status).
			Save(ctx)
		if err != nil {
			if ent.IsNotFound(err) {
				return status.Wrap(errors.New("agent task not found"), status.NotFound)
			}
			return status.Wrap(err, status.Internal)
		}

		res.ID = task.ID
		res.Status = task.Status
		return nil
	}
}
