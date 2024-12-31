package agenttasks_test

import (
	"context"
	"github.com/shinobistack/gokakashi/ent/schema"
	"github.com/shinobistack/gokakashi/internal/restapi/v1/agenttasks"
	"testing"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
	"github.com/shinobistack/gokakashi/ent/enttest"
	"github.com/stretchr/testify/assert"
)

func TestDeleteAgentTask_Valid(t *testing.T) {
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

	task := client.AgentTasks.Create().
		SetAgentID(agent.ID).
		SetScanID(scan.ID).
		SetStatus("pending").
		SaveX(context.Background())

	req := agenttasks.DeleteAgentTaskRequest{ID: task.ID}
	res := &agenttasks.DeleteAgentTaskResponse{}

	err := agenttasks.DeleteAgentTask(client)(context.Background(), req, res)

	assert.NoError(t, err)
	assert.Equal(t, task.ID, res.ID)
	assert.Equal(t, "deleted", res.Status)
}

func TestDeleteAgentTask_NotFound(t *testing.T) {
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	defer client.Close()

	req := agenttasks.DeleteAgentTaskRequest{ID: uuid.New()} // Non-existent ID
	res := &agenttasks.DeleteAgentTaskResponse{}

	err := agenttasks.DeleteAgentTask(client)(context.Background(), req, res)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "agent task not found")
}
