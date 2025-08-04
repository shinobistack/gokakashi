package agents

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/shinobistack/gokakashi/ent"
	"github.com/shinobistack/gokakashi/ent/v2agenttasks"
	"github.com/shinobistack/gokakashi/internal/restapi/v2/io"
)

// mockEntClient and helpers would be implemented here or imported from a testutil package.
// For brevity, we use a minimal stub approach. Replace with proper ent test setup as needed.

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
	tests := []struct {
		name   string
		query  *mockV2AgentTasksQuery
		req    io.AgentTaskListRequest
		want   io.AgentTaskListResponse
		wantErr bool
	}{
		{
			name:  "no tasks",
			query: &mockV2AgentTasksQuery{tasks: []*ent.V2AgentTasks{}, total: 0},
			req:   io.AgentTaskListRequest{AgentID: 1, Page: 1, PerPage: 10},
			want: io.AgentTaskListResponse{
				Tasks:      []io.AgentTask{},
				Pagination: io.Pagination{Page: 1, PerPage: 10},
				Total:      0,
			},
			wantErr: false,
		},
		{
			name: "one task",
			query: &mockV2AgentTasksQuery{
				tasks: []*ent.V2AgentTasks{{
					ID:        42,
					ScanID:    10,
					AgentID:   1,
					Status:    "pending",
					CreatedAt: now,
					UpdatedAt: now,
				}},
				total: 1,
			},
			req: io.AgentTaskListRequest{AgentID: 1, Page: 1, PerPage: 10},
			want: io.AgentTaskListResponse{
				Tasks: []io.AgentTask{{
					ID:        42,
					ScanID:    10,
					AgentID:   1,
					Status:    "pending",
					CreatedAt: now,
					UpdatedAt: now,
				}},
				Pagination: io.Pagination{Page: 1, PerPage: 10},
				Total:      1,
			},
			wantErr: false,
		},
		{
			name:   "db error",
			query:  &mockV2AgentTasksQuery{err: errors.New("db fail")},
			req:    io.AgentTaskListRequest{AgentID: 1, Page: 1, PerPage: 10},
			want:   io.AgentTaskListResponse{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := &mockEntClient{v2AgentTasks: &mockV2AgentTasksClient{query: tt.query}}
			handler := ListAgentTasks((*ent.Client)(client)) // type cast for compatibility
			var res io.AgentTaskListResponse
			err := handler(context.Background(), tt.req, &res)
			if (err != nil) != tt.wantErr {
				t.Errorf("error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && len(res.Tasks) != len(tt.want.Tasks) {
				t.Errorf("got %d tasks, want %d", len(res.Tasks), len(tt.want.Tasks))
			}
			// Additional deep checks can be added here for Tasks, Pagination, etc.
		})
	}
}
