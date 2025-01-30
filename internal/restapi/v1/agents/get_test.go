package agents_test

import (
	"context"
	"github.com/shinobistack/gokakashi/internal/restapi/v1/agents"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/shinobistack/gokakashi/ent/enttest"
	"github.com/stretchr/testify/assert"
)

func TestGetAgent_Valid(t *testing.T) {
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	defer client.Close()

	agent := client.Agents.Create().
		SetStatus("connected").
		SaveX(context.Background())

	req := agents.GetAgentRequest{ID: agent.ID}
	res := &agents.GetAgentResponse{}

	err := agents.GetAgent(client)(context.Background(), req, res)

	assert.NoError(t, err)
	assert.Equal(t, agent.Status, res.Status)
	assert.Equal(t, agent.ID, res.ID)
}

func TestGetAgent_NotFound(t *testing.T) {
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	defer client.Close()

	req := agents.GetAgentRequest{ID: 9999} // Non-existent ID
	res := &agents.GetAgentResponse{}

	err := agents.GetAgent(client)(context.Background(), req, res)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "agent not found")
}

func TestListAgents_Valid(t *testing.T) {
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	defer client.Close()

	client.Agents.Create().
		SetStatus("connected").
		SaveX(context.Background())

	client.Agents.Create().
		SetStatus("scan_in_progress").
		SaveX(context.Background())

	req := agents.ListAgentsRequest{}
	res := []agents.GetAgentResponse{}

	err := agents.ListAgents(client)(context.Background(), req, &res)

	assert.NoError(t, err)
	assert.Equal(t, 2, len(res))

}
