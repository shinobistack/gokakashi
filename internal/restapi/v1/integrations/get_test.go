package integrations_test

import (
	"context"
	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
	"github.com/shinobistack/gokakashi/ent/enttest"
	"github.com/shinobistack/gokakashi/internal/restapi/v1/integrations"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetIntegration(t *testing.T) {
	// client := enttest.Open(t, "postgres", "host=localhost port=5432 user=postgres password=secret dbname=testdb sslmode=disable")
	// This requires a running DB with pre-configured testdb

	// sqlite3 lightweight, file based db faster to test
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	defer client.Close()

	// Seed the DB
	id := uuid.New()
	_, err := client.Integrations.Create().
		SetID(id).
		SetName("Test Integration").
		SetType("linear").
		SetConfig(map[string]interface{}{"Key": "value"}).
		Save(context.Background())
	assert.NoError(t, err)

	// Test case: Valid ID
	req := integrations.GetIntegrationRequests{ID: id.String()}
	res := &integrations.GetIntegrationResponse{}
	handler := integrations.GetIntegration(client)
	err = handler(context.Background(), req, res)

	assert.NoError(t, err)
	assert.Equal(t, "Test Integration", res.Name)
	assert.Equal(t, "linear", res.Type)

	// Test case: Invalid UUID
	req = integrations.GetIntegrationRequests{ID: "invalid-uuid"}
	err = handler(context.Background(), req, res)
	assert.Error(t, err)

	// Test case: Non-existent ID
	req = integrations.GetIntegrationRequests{ID: uuid.New().String()}
	err = handler(context.Background(), req, res)
	assert.Error(t, err)

}

func TestListIntegrations(t *testing.T) {
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	defer client.Close()

	// Seed the db
	_, err := client.Integrations.Create().
		SetName("Integration 1").
		SetType("linear").
		SetConfig(map[string]interface{}{"key": "value1"}).
		Save(context.Background())
	assert.NoError(t, err)

	_, err = client.Integrations.Create().
		SetName("Integration 2").
		SetType("jira").
		SetConfig(map[string]interface{}{"key": "value2"}).
		Save(context.Background())
	assert.NoError(t, err)

	// Test case: List integrations
	req := struct{}{}
	var res []integrations.GetIntegrationResponse
	handler := integrations.ListIntegrations(client)
	err = handler(context.Background(), req, &res)

	assert.NoError(t, err)
	assert.Len(t, res, 2)
	assert.Equal(t, "Integration 1", res[0].Name)
	assert.Equal(t, "Integration 2", res[1].Name)
}

func TestListIntegrations_EmptyDB(t *testing.T) {
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	defer client.Close()

	// Prepare response
	var res []integrations.GetIntegrationResponse

	// Execute ListIntegrations handler
	handler := integrations.ListIntegrations(client)
	err := handler(context.Background(), struct{}{}, &res)

	// Validate response
	assert.NoError(t, err)
	assert.Len(t, res, 0)
}
