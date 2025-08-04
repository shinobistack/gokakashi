package agents

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/shinobistack/gokakashi/ent/enttest"
	"github.com/shinobistack/gokakashi/internal/restapi/v2/io"
	"github.com/stretchr/testify/require"
)

func TestListAgentTasks(t *testing.T) {
	now := time.Now()

	t.Run("no tasks", func(t *testing.T) {
		client := enttest.Open(t, "sqlite3", "file:agenttasklist_notasks?mode=memory&cache=shared&_fk=1")
		defer client.Close()
		ctx := context.Background()
		agentID := uuid.MustParse("00000000-0000-0000-0000-000000000001")
		// Create agent only, no tasks
		_, err := client.V2Agents.Create().
			SetID(agentID).
			SetStatus("disconnected").
			Save(ctx)
		require.NoError(t, err)

		req := io.AgentTaskListRequest{
			AgentID:    agentID,
			Pagination: io.Pagination{Page: 1, PerPage: 10},
		}
		var res io.AgentTaskListResponse
		err = ListAgentTasks(client)(ctx, req, &res)
		require.NoError(t, err)
		require.Equal(t, 0, len(res.Tasks))
		require.Equal(t, 0, res.Total)
		require.Equal(t, 1, res.Pagination.Page)
		require.Equal(t, 10, res.Pagination.PerPage)
	})

	t.Run("one task", func(t *testing.T) {
		client := enttest.Open(t, "sqlite3", "file:agenttasklist_onetask?mode=memory&cache=shared&_fk=1")
		defer client.Close()
		ctx := context.Background()
		agentID := uuid.MustParse("00000000-0000-0000-0000-000000000001")
		scanID := uuid.MustParse("00000000-0000-0000-0000-000000000010")
		taskID := uuid.MustParse("00000000-0000-0000-0000-000000000042")
		// Create agent
		_, err := client.V2Agents.Create().
			SetID(agentID).
			SetStatus("disconnected").
			Save(ctx)
		require.NoError(t, err)
		// Create task
		_, err = client.V2AgentTasks.Create().
			SetID(taskID).
			SetScanID(scanID).
			SetAgentID(agentID).
			SetStatus("pending").
			SetCreatedAt(now).
			SetUpdatedAt(now).
			Save(ctx)
		require.NoError(t, err)

		req := io.AgentTaskListRequest{
			AgentID:    agentID,
			Pagination: io.Pagination{Page: 1, PerPage: 10},
		}
		var res io.AgentTaskListResponse
		err = ListAgentTasks(client)(ctx, req, &res)
		require.NoError(t, err)
		require.Equal(t, 1, len(res.Tasks))
		require.Equal(t, 1, res.Total)
		task := res.Tasks[0]
		require.Equal(t, taskID, task.ID)
		require.Equal(t, scanID, task.ScanID)
		require.Equal(t, agentID, task.AgentID)
		require.Equal(t, "pending", string(task.Status))
		require.WithinDuration(t, now, task.CreatedAt, time.Second)
		require.WithinDuration(t, now, task.UpdatedAt, time.Second)
		require.Equal(t, 1, res.Pagination.Page)
		require.Equal(t, 10, res.Pagination.PerPage)
	})
}
