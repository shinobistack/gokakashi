package agenttasks

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/shinobistack/gokakashi/ent"
	"github.com/swaggest/usecase/status"
	"time"
)

type CreateAgentTaskRequest struct {
	AgentID   int       `json:"agent_id"`
	ScanID    uuid.UUID `json:"scan_id"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at,omitempty"`
}

type CreateAgentTaskResponse struct {
	ID     int    `json:"id"`
	Status string `json:"status"`
}

func CreateAgentTask(client *ent.Client) func(ctx context.Context, req CreateAgentTaskRequest, res *CreateAgentTaskResponse) error {
	return func(ctx context.Context, req CreateAgentTaskRequest, res *CreateAgentTaskResponse) error {
		if req.AgentID <= 0 || req.ScanID == uuid.Nil || req.Status == "" {
			return status.Wrap(errors.New("missing required fields"), status.InvalidArgument)
		}

		task, err := client.AgentTasks.Create().
			SetAgentID(req.AgentID).
			SetScanID(req.ScanID).
			SetStatus(req.Status).
			SetCreatedAt(time.Now()).
			Save(ctx)

		if err != nil {
			return status.Wrap(err, status.Internal)
		}

		res.ID = task.ID
		res.Status = task.Status
		return nil
	}
}
