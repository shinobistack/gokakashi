package agenttasks

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/shinobistack/gokakashi/ent"
	"github.com/shinobistack/gokakashi/ent/agents"
	"github.com/shinobistack/gokakashi/ent/agenttasks"
	"github.com/shinobistack/gokakashi/ent/scans"
	"github.com/swaggest/usecase/status"
	"time"
)

type CreateAgentTaskRequest struct {
	AgentID   int       `path:"agent_id"`
	ScanID    uuid.UUID `json:"scan_id"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at,omitempty"`
}

type CreateAgentTaskResponse struct {
	ID     uuid.UUID `json:"id"`
	Status string    `json:"status"`
}

func CreateAgentTask(client *ent.Client) func(ctx context.Context, req CreateAgentTaskRequest, res *CreateAgentTaskResponse) error {
	return func(ctx context.Context, req CreateAgentTaskRequest, res *CreateAgentTaskResponse) error {
		if req.AgentID <= 0 || req.ScanID == uuid.Nil || req.Status == "" {
			return status.Wrap(errors.New("missing required fields"), status.InvalidArgument)
		}

		// Check if the agent exists
		agentExists, err := client.Agents.Query().Where(agents.ID(req.AgentID)).Exist(ctx)
		if err != nil {
			return status.Wrap(fmt.Errorf("failed to check agent existence: %w", err), status.Internal)
		}
		if !agentExists {
			return status.Wrap(errors.New("agent not found"), status.NotFound)
		}

		// Check if the scan exists
		scanExists, err := client.Scans.Query().Where(scans.ID(req.ScanID)).Exist(ctx)
		if err != nil {
			return status.Wrap(fmt.Errorf("failed to check scan existence: %w", err), status.Internal)
		}
		if !scanExists {
			return status.Wrap(errors.New("scan not found"), status.NotFound)
		}

		// Ensure the same scan ID isn't already assigned to another task
		existingTask, err := client.AgentTasks.Query().
			Where(agenttasks.ScanID(req.ScanID)).
			First(ctx)
		if err != nil && !ent.IsNotFound(err) {
			return status.Wrap(fmt.Errorf("failed to check for existing tasks: %w", err), status.Internal)
		}
		if existingTask != nil {
			return status.Wrap(errors.New("scan ID is already assigned to another agent"), status.InvalidArgument)
		}

		// Create the agent task
		task, err := client.AgentTasks.Create().
			SetAgentID(req.AgentID).
			SetScanID(req.ScanID).
			SetStatus(req.Status).
			SetCreatedAt(time.Now()).
			Save(ctx)

		if err != nil {
			// Handle foreign key constraint errors
			if ent.IsConstraintError(err) {
				return status.Wrap(errors.New("constraint violation: ensure valid agent ID and scan ID"), status.InvalidArgument)
			}
			return status.Wrap(fmt.Errorf("failed to create agent task: %w", err), status.Internal)
		}

		// Populate the response
		res.ID = task.ID
		res.Status = task.Status
		return nil
	}
}
