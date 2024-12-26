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

func TestCreateAgentTask_Valid(t *testing.T) {
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	defer client.Close()

	policy := client.Policies.Create().
		SetName("to-be-deleted-test-policy").
		SetImage(schema.Image{Registry: "example-registry", Name: "example-name", Tags: []string{"v1.0"}}).
		SaveX(context.Background())

	scan := client.Scans.Create().
		SetPolicyID(policy.ID).
		SetImage("example-image:latest").
		SetStatus("scan_pending").
		SaveX(context.Background())

	agent := client.Agents.Create().
		SetStatus("connected").
		SaveX(context.Background())

	// Create test request
	req := agenttasks.CreateAgentTaskRequest{
		AgentID: agent.ID,
		ScanID:  scan.ID,
		Status:  "pending",
	}
	res := &agenttasks.CreateAgentTaskResponse{}

	err := agenttasks.CreateAgentTask(client)(context.Background(), req, res)

	assert.NoError(t, err)
	assert.Equal(t, req.Status, res.Status)
	assert.NotZero(t, res.ID)
}

func TestCreateAgentTask_MissingFields(t *testing.T) {
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	defer client.Close()

	req := agenttasks.CreateAgentTaskRequest{
		AgentID: 0, // Missing AgentID
		ScanID:  uuid.Nil,
		Status:  "",
	}
	res := &agenttasks.CreateAgentTaskResponse{}

	err := agenttasks.CreateAgentTask(client)(context.Background(), req, res)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "missing required fields")
}
