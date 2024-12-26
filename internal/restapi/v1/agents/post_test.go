package agents_test

import (
	"context"
	"github.com/shinobistack/gokakashi/internal/restapi/v1/agents"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/shinobistack/gokakashi/ent/enttest"

	"github.com/stretchr/testify/assert"
)

func TestCreateAgent_Valid(t *testing.T) {
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	defer client.Close()

	req := agents.CreateAgentRequest{
		Status: "connected",
	}
	res := &agents.CreateAgentResponse{}

	err := agents.CreateAgent(client)(context.Background(), req, res)

	assert.NoError(t, err)
	assert.Equal(t, req.Status, res.Status)
	assert.NotZero(t, res.ID)
}

//func TestCreateAgent_MissingFields(t *testing.T) {
//	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
//	defer client.Close()
//
//	req := agents.CreateAgentRequest{
//		Status: "", // Missing Status
//	}
//	res := &agents.CreateAgentResponse{}
//
//	err := agents.CreateAgent(client)(context.Background(), req, res)
//
//	assert.Error(t, err)
//	assert.Contains(t, err.Error(), "missing required fields")
//}
