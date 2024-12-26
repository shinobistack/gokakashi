package agenttasks

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/shinobistack/gokakashi/ent"
	"github.com/shinobistack/gokakashi/ent/agenttasks"
	"github.com/swaggest/usecase/status"
)

type UpdateAgentTaskRequest struct {
	ID      uuid.UUID `path:"id"`
	AgentID int       `path:"agent_id"`
	Status  string    `json:"status"`
	// Todo: Should the created_AT be updated to time.now whenever an update call is made?
}

type UpdateAgentTaskResponse struct {
	ID     uuid.UUID `json:"id"`
	Status string    `json:"status"`
}

func UpdateAgentTask(client *ent.Client) func(ctx context.Context, req UpdateAgentTaskRequest, res *UpdateAgentTaskResponse) error {
	return func(ctx context.Context, req UpdateAgentTaskRequest, res *UpdateAgentTaskResponse) error {
		if req.ID == uuid.Nil || req.Status == "" {
			return status.Wrap(errors.New("invalid ID or Status"), status.InvalidArgument)
		}

		// Fetch the task and validate the AgentID
		task, err := client.AgentTasks.Query().
			Where(agenttasks.ID(req.ID)).
			Only(ctx)
		if err != nil {
			if ent.IsNotFound(err) {
				return status.Wrap(errors.New("agent task not found"), status.NotFound)
			}
			return status.Wrap(err, status.Internal)
		}

		if task.AgentID != req.AgentID {
			return status.Wrap(errors.New("agent ID mismatch"), status.InvalidArgument)
		}

		task, err = client.AgentTasks.UpdateOneID(req.ID).
			SetStatus(req.Status).
			Save(ctx)
		if err != nil {
			return status.Wrap(err, status.Internal)
		}

		res.ID = task.ID
		res.Status = task.Status
		return nil

	}
}
