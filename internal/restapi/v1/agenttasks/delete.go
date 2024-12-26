package agenttasks

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/shinobistack/gokakashi/ent"
	"github.com/swaggest/usecase/status"
)

type DeleteAgentTaskRequest struct {
	ID      uuid.UUID `path:"id"`
	AgentID int       `path:"agent_id"`
}

type DeleteAgentTaskResponse struct {
	ID     uuid.UUID `json:"id"`
	Status string    `json:"status"`
}

func DeleteAgentTask(client *ent.Client) func(ctx context.Context, req DeleteAgentTaskRequest, res *DeleteAgentTaskResponse) error {
	return func(ctx context.Context, req DeleteAgentTaskRequest, res *DeleteAgentTaskResponse) error {
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

		err = client.AgentTasks.DeleteOne(task).Exec(ctx)
		if err != nil {
			return status.Wrap(err, status.Internal)
		}

		res.ID = task.ID
		res.Status = "deleted"
		return nil
	}
}
