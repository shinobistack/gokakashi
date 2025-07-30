package agents

import (
	"context"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/shinobistack/gokakashi/ent/enttest"
	"github.com/shinobistack/gokakashi/internal/agent"
	"github.com/stretchr/testify/assert"
)

func TestRegisterAgent_Integration_Success(t *testing.T) {
	client := enttest.Open(t, "sqlite3", "file:enttest?mode=memory&cache=shared&_fk=1")
	defer client.Close()

	handler := Register(client)
	var res RegisterAgentResponse
	err := handler(context.Background(), RegisterAgentRequest{}, &res)
	assert.NoError(t, err)

	dbAgent, err := client.V2Agents.Get(context.Background(), res.ID)
	assert.NoError(t, err)
	assert.Equal(t, dbAgent.ID, res.ID)
	assert.Equal(t, string(agent.Disconnected), dbAgent.Status)
	assert.Equal(t, agent.Status(dbAgent.Status), res.Status)
	assert.WithinDuration(t, dbAgent.CreatedAt, res.CreatedAt, time.Second)
	assert.WithinDuration(t, dbAgent.UpdatedAt, res.UpdatedAt, time.Second)
	assert.WithinDuration(t, dbAgent.LastHeartbeatAt, res.LastHeartbeatAt, time.Second)
}

func TestRegisterAgent_Integration_DBError(t *testing.T) {
	client := enttest.Open(t, "sqlite3", "file:enttest2?mode=memory&cache=shared&_fk=1")
	client.Close() // force DB error

	handler := Register(client)
	var res RegisterAgentResponse
	err := handler(context.Background(), RegisterAgentRequest{}, &res)
	assert.Error(t, err)
}
