package agents

import (
	"context"

	"github.com/shinobistack/gokakashi/ent"
	"github.com/shinobistack/gokakashi/ent/v2agenttasks"
	"github.com/shinobistack/gokakashi/internal/agent/task"
	"github.com/shinobistack/gokakashi/internal/restapi/v2/io"
	"github.com/swaggest/usecase/status"
)

func ListAgentTasks(client *ent.Client) func(ctx context.Context, req io.AgentTaskListRequest, res *io.AgentTaskListResponse) error {
	return func(ctx context.Context, req io.AgentTaskListRequest, res *io.AgentTaskListResponse) error {
		q := client.V2AgentTasks.Query().Where(v2agenttasks.AgentID(req.AgentID))
		if req.Status != nil {
			q = q.Where(v2agenttasks.Status(*req.Status))
		}

		// Pagination defaults
		page := req.Page
		perPage := req.PerPage
		if page <= 0 {
			page = 1
		}
		if perPage <= 0 {
			perPage = 25
		}
		if perPage > 100 {
			perPage = 100
		}
		offset := (page - 1) * perPage

		// Get total count
		total, err := q.Clone().Count(ctx)
		if err != nil {
			return status.Wrap(err, status.Internal)
		}

		tasks, err := q.Order(ent.Asc(v2agenttasks.FieldCreatedAt)).Limit(perPage).Offset(offset).All(ctx)
		if err != nil {
			return status.Wrap(err, status.Internal)
		}

		*res = io.AgentTaskListResponse{
			Tasks:      make([]io.AgentTask, len(tasks)),
			Pagination: io.Pagination{Page: page, PerPage: perPage},
			Total:      total,
		}

		for i, t := range tasks {
			res.Tasks[i] = io.AgentTask{
				ID:        t.ID,
				ScanID:    t.ScanID,
				AgentID:   t.AgentID,
				Status:    task.Status(t.Status),
				CreatedAt: t.CreatedAt,
				UpdatedAt: t.UpdatedAt,
			}
		}

		return nil
	}
}
