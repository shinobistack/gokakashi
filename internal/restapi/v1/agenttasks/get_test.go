package agenttasks_test

import (
	"context"
	"github.com/shinobistack/gokakashi/ent/schema"
	"github.com/shinobistack/gokakashi/internal/restapi/v1/agenttasks"
	"testing"
	"time"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
	"github.com/shinobistack/gokakashi/ent/enttest"
	"github.com/stretchr/testify/assert"
)

func TestGetAgentTask_Valid(t *testing.T) {
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	defer client.Close()

	policy := client.Policies.Create().
		SetName("to-be-deleted-test-policy").
		SetImage(schema.Image{Registry: "example-registry", Name: "example-name", Tags: []string{"v1.0"}}).
		SetScanner("trivy").
		SaveX(context.Background())

	integrations := client.Integrations.Create().
		SetName("integration").
		SetType("linear").
		SetConfig(map[string]interface{}{"key": "value"}).
		SaveX(context.Background())

	scan := client.Scans.Create().
		SetPolicyID(policy.ID).
		SetImage("example-image:latest").
		SetScanner(policy.Scanner).
		SetStatus("scan_pending").
		SetIntegrationID(integrations.ID).
		SaveX(context.Background())

	agent := client.Agents.Create().
		SetStatus("connected").
		SaveX(context.Background())

	task := client.AgentTasks.Create().
		SetAgentID(agent.ID).
		SetScanID(scan.ID).
		SetStatus("pending").
		SaveX(context.Background())

	req := agenttasks.GetAgentTaskRequest{ID: task.ID}
	res := &agenttasks.GetAgentTaskResponse{}

	err := agenttasks.GetAgentTask(client)(context.Background(), req, res)

	assert.NoError(t, err)
	assert.Equal(t, task.Status, res.Status)
	assert.Equal(t, task.ID, res.ID)
	assert.Equal(t, task.AgentID, res.AgentID)
	assert.Equal(t, task.ScanID, res.ScanID)
}

func TestGetAgentTask_NotFound(t *testing.T) {
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	defer client.Close()

	req := agenttasks.GetAgentTaskRequest{ID: uuid.New()} // Non-existent ID
	res := &agenttasks.GetAgentTaskResponse{}

	err := agenttasks.GetAgentTask(client)(context.Background(), req, res)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "agent task not found")
}

func TestListAgentTasks_Valid(t *testing.T) {
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	defer client.Close()

	policy := client.Policies.Create().
		SetName("to-be-deleted-test-policy").
		SetImage(schema.Image{Registry: "example-registry", Name: "example-name", Tags: []string{"v1.0"}}).
		SetScanner("trivy").
		SaveX(context.Background())

	integrations := client.Integrations.Create().
		SetName("integration").
		SetType("linear").
		SetConfig(map[string]interface{}{"key": "value"}).
		SaveX(context.Background())

	scan := client.Scans.Create().
		SetPolicyID(policy.ID).
		SetImage("example-image:latest").
		SetStatus("scan_pending").
		SetScanner(policy.Scanner).
		SetIntegrationID(integrations.ID).
		SaveX(context.Background())

	agent := client.Agents.Create().
		SetStatus("connected").
		SaveX(context.Background())

	client.AgentTasks.Create().
		SetAgentID(agent.ID).
		SetScanID(scan.ID).
		SetStatus("pending").
		SaveX(context.Background())
	client.AgentTasks.Create().
		SetAgentID(agent.ID).
		SetScanID(scan.ID).
		SetStatus("pending").
		SaveX(context.Background())

	req := agenttasks.ListAgentTasksRequest{agent.ID}
	res := []agenttasks.GetAgentTaskResponse{}

	err := agenttasks.ListAgentTasksByAgentID(client)(context.Background(), req, &res)

	assert.NoError(t, err)
	assert.Equal(t, 2, len(res))

}

func TestListAgentTasks_OrderedByCreatedAt(t *testing.T) {
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	defer client.Close()

	policy := client.Policies.Create().
		SetName("to-be-deleted-test-policy").
		SetImage(schema.Image{Registry: "example-registry", Name: "example-name", Tags: []string{"v1.0"}}).
		SetScanner("trivy").
		SaveX(context.Background())

	integrations := client.Integrations.Create().
		SetName("integration").
		SetType("linear").
		SetConfig(map[string]interface{}{"key": "value"}).
		SaveX(context.Background())

	scan := client.Scans.Create().
		SetPolicyID(policy.ID).
		SetImage("example-image:latest").
		SetStatus("scan_pending").
		SetScanner(policy.Scanner).
		SetIntegrationID(integrations.ID).
		SaveX(context.Background())

	agent := client.Agents.Create().
		SetStatus("connected").
		SaveX(context.Background())

	task1 := client.AgentTasks.Create().
		SetAgentID(agent.ID).
		SetScanID(scan.ID).
		SetStatus("pending").
		SetCreatedAt(time.Now().Add(-10 * time.Minute)). // Older task
		SaveX(context.Background())

	task2 := client.AgentTasks.Create().
		SetAgentID(agent.ID).
		SetScanID(scan.ID).
		SetStatus("pending").
		SetCreatedAt(time.Now()). // Newer task
		SaveX(context.Background())

	req := agenttasks.ListAgentTasksRequest{AgentID: agent.ID}
	res := []agenttasks.GetAgentTaskResponse{}

	err := agenttasks.ListAgentTasksByAgentID(client)(context.Background(), req, &res)

	assert.NoError(t, err)
	assert.Equal(t, 2, len(res))
	assert.Equal(t, task1.ID, res[0].ID) // Oldest task first
	assert.Equal(t, task2.ID, res[1].ID) // Newest task last
}
