package agents

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/shinobistack/gokakashi/ent/enttest"
	"github.com/shinobistack/gokakashi/ent"
	"github.com/shinobistack/gokakashi/internal/restapi/v2/io"
	"github.com/stretchr/testify/require"
)

type mockV2AgentTasksQuery struct {
	tasks []*ent.V2AgentTasks
	total int
	err  error
}

func (q *mockV2AgentTasksQuery) Where(...interface{}) *mockV2AgentTasksQuery { return q }
func (q *mockV2AgentTasksQuery) Clone() *mockV2AgentTasksQuery               { return q }
func (q *mockV2AgentTasksQuery) Count(ctx context.Context) (int, error)      { return q.total, q.err }
func (q *mockV2AgentTasksQuery) Order(...interface{}) *mockV2AgentTasksQuery { return q }
func (q *mockV2AgentTasksQuery) Limit(int) *mockV2AgentTasksQuery           { return q }
func (q *mockV2AgentTasksQuery) Offset(int) *mockV2AgentTasksQuery          { return q }
func (q *mockV2AgentTasksQuery) All(ctx context.Context) ([]*ent.V2AgentTasks, error) {
	if q.err != nil {
		return nil, q.err
	}
	return q.tasks, nil
}

type mockV2AgentTasksClient struct {
	query *mockV2AgentTasksQuery
}

func (c *mockV2AgentTasksClient) Query() *mockV2AgentTasksQuery { return c.query }

type mockEntClient struct {
	v2AgentTasks *mockV2AgentTasksClient
}

func (c *mockEntClient) V2AgentTasks() *mockV2AgentTasksClient { return c.v2AgentTasks }

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
			AgentID: agentID,
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
			AgentID: agentID,
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
