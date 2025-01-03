package agents_test

import (
	"context"
	"github.com/shinobistack/gokakashi/internal/restapi/v1/agents"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/shinobistack/gokakashi/ent/enttest"
	"github.com/stretchr/testify/assert"
)

func TestUpdateAgent_Valid(t *testing.T) {
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	defer client.Close()

	agent := client.Agents.Create().
		SetStatus("connected").
		SaveX(context.Background())

	req := agents.UpdateAgentRequest{
		ID:     agent.ID,
		Status: "scan_in_progress",
	}
	res := &agents.UpdateAgentResponse{}

	err := agents.UpdateAgent(client)(context.Background(), req, res)

	assert.NoError(t, err)
	assert.Equal(t, req.Status, res.Status)
}

func TestUpdateAgent_NotFound(t *testing.T) {
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	defer client.Close()

	req := agents.UpdateAgentRequest{
		ID:     9999, // Non-existent ID
		Status: "scan_in_progress",
	}
	res := &agents.UpdateAgentResponse{}

	err := agents.UpdateAgent(client)(context.Background(), req, res)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "agent not found")
}
