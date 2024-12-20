package integrationtype

import (
	"context"
	_ "github.com/mattn/go-sqlite3"
	"github.com/shinobistack/gokakashi/ent/enttest"
	"github.com/shinobistack/gokakashi/internal/restapi/v1/integrationtype"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetIntegrationType_ValidID(t *testing.T) {
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	defer client.Close()

	// Seed data
	_, _ = client.IntegrationType.Create().
		SetID("linear").
		SetDisplayName("Linear Integration").
		Save(context.Background())

	// Test case: Valid ID
	req := integrationtype.GetIntegrationTypeRequests{ID: "linear"}
	res := &integrationtype.GetIntegrationTypeResponse{}
	handler := integrationtype.GetIntegrationType(client)
	err := handler(context.Background(), req, res)

	assert.NoError(t, err)
	assert.Equal(t, "linear", res.ID)
	assert.Equal(t, "Linear Integration", res.DisplayName)
}

func TestGetIntegrationType_NonExistentID(t *testing.T) {
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	defer client.Close()

	req := integrationtype.GetIntegrationTypeRequests{ID: "nonexistent"}
	res := &integrationtype.GetIntegrationTypeResponse{}
	handler := integrationtype.GetIntegrationType(client)
	err := handler(context.Background(), req, res)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "integration type not found")
}

func TestGetIntegrationType_InvalidIDFormat(t *testing.T) {
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	defer client.Close()

	req := integrationtype.GetIntegrationTypeRequests{ID: "Invalid ID!"}
	res := &integrationtype.GetIntegrationTypeResponse{}
	handler := integrationtype.GetIntegrationType(client)
	err := handler(context.Background(), req, res)

	req = integrationtype.GetIntegrationTypeRequests{ID: "Inv*ID"}
	handler = integrationtype.GetIntegrationType(client)
	err = handler(context.Background(), req, res)

	assert.Error(t, err)
}

func TestListIntegrationTypes_ValidRequest(t *testing.T) {
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	defer client.Close()

	// Seed data
	_, _ = client.IntegrationType.Create().
		SetID("linear").
		SetDisplayName("Linear Integration").
		Save(context.Background())
	_, _ = client.IntegrationType.Create().
		SetID("jira").
		SetDisplayName("Jira Integration").
		Save(context.Background())

	req := struct{}{}
	var res []integrationtype.GetIntegrationTypeResponse
	handler := integrationtype.ListIntegrationType(client)
	err := handler(context.Background(), req, &res)

	assert.NoError(t, err)
	assert.Len(t, res, 2)
}

func TestListIntegrations_EmptyDB(t *testing.T) {
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	defer client.Close()

	// Prepare response
	var res []integrationtype.GetIntegrationTypeResponse

	// Execute ListIntegrations handler
	handler := integrationtype.ListIntegrationType(client)
	err := handler(context.Background(), struct{}{}, &res)

	// Validate response
	assert.NoError(t, err)
	assert.Len(t, res, 0)
}
