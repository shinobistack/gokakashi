package integrationtype_test

import (
	"context"
	"github.com/shinobistack/gokakashi/internal/restapi/v1/integrationtype"
	"testing"

	"github.com/shinobistack/gokakashi/ent/enttest"
	"github.com/stretchr/testify/assert"
)

func TestCreateIntegrationType_ValidInput(t *testing.T) {
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	defer client.Close()

	req := integrationtype.CreateIntegrationTypeRequest{
		ID:          "linear",
		DisplayName: "Linear Integration",
	}
	res := &integrationtype.GetIntegrationTypeResponse{}
	handler := integrationtype.CreateIntegrationType(client)
	err := handler(context.Background(), req, res)

	assert.NoError(t, err)
	assert.Equal(t, "linear", res.ID)
	assert.Equal(t, "Linear Integration", res.DisplayName)
}

func TestCreateIntegrationType_MissingFields(t *testing.T) {
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	defer client.Close()

	req := integrationtype.CreateIntegrationTypeRequest{}
	res := &integrationtype.GetIntegrationTypeResponse{}
	handler := integrationtype.CreateIntegrationType(client)
	err := handler(context.Background(), req, res)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "missing required fields")
}

func TestCreateIntegrationType_InvalidIDFormat(t *testing.T) {
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	defer client.Close()

	req := integrationtype.CreateIntegrationTypeRequest{
		ID:          "Invalid ID!",
		DisplayName: "Valid Format ",
	}
	res := &integrationtype.GetIntegrationTypeResponse{}
	handler := integrationtype.CreateIntegrationType(client)
	err := handler(context.Background(), req, res)
	assert.Error(t, err)

	req = integrationtype.CreateIntegrationTypeRequest{
		ID:          "invalid*#",
		DisplayName: "Valid Format ",
	}
	handler = integrationtype.CreateIntegrationType(client)
	err = handler(context.Background(), req, res)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid id format")
}

func TestCreateIntegrationType_DuplicateID(t *testing.T) {
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	defer client.Close()

	_, _ = client.IntegrationType.Create().
		SetID("linear").
		SetDisplayName("Linear Integration").
		Save(context.Background())

	req := integrationtype.CreateIntegrationTypeRequest{
		ID:          "linear",
		DisplayName: "Linear Integration",
	}
	res := &integrationtype.GetIntegrationTypeResponse{}
	handler := integrationtype.CreateIntegrationType(client)
	err := handler(context.Background(), req, res)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "integration type already exists")
}
