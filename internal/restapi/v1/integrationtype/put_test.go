package integrationtype_test

import (
	"context"
	"github.com/shinobistack/gokakashi/internal/restapi/v1/integrationtype"
	"testing"

	"github.com/shinobistack/gokakashi/ent/enttest"
	"github.com/stretchr/testify/assert"
)

func TestUpdateIntegrationType_ValidUpdate(t *testing.T) {
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	defer client.Close()

	_, _ = client.IntegrationType.Create().
		SetID("linear").
		SetDisplayName("Linear Integration").
		Save(context.Background())

	req := integrationtype.UpdateIntegrationTypeRequest{
		ID:          "linear",
		DisplayName: ptr("Updated Linear Integration"),
	}
	res := &integrationtype.GetIntegrationTypeResponse{}
	handler := integrationtype.UpdateIntegrationType(client)
	err := handler(context.Background(), req, res)

	assert.NoError(t, err)
	assert.Equal(t, "linear", res.ID)
	assert.Equal(t, "Updated Linear Integration", res.DisplayName)
}

func TestUpdateIntegrationType_InvalidIDFormat(t *testing.T) {
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	defer client.Close()

	req := integrationtype.UpdateIntegrationTypeRequest{
		ID:          "Invalid ID!",
		DisplayName: ptr("Invalid Format"),
	}
	res := &integrationtype.GetIntegrationTypeResponse{}
	handler := integrationtype.UpdateIntegrationType(client)
	err := handler(context.Background(), req, res)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid id format")
}

func TestUpdateIntegrationType_NonExistentID(t *testing.T) {
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	defer client.Close()

	req := integrationtype.UpdateIntegrationTypeRequest{
		ID:          "nonexistent",
		DisplayName: ptr("Non-existent Integration"),
	}
	res := &integrationtype.GetIntegrationTypeResponse{}
	handler := integrationtype.UpdateIntegrationType(client)
	err := handler(context.Background(), req, res)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "integration type not found")
}

func TestUpdateIntegrationType_NoFieldsToUpdate(t *testing.T) {
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	defer client.Close()

	req := integrationtype.UpdateIntegrationTypeRequest{
		ID:          "linear",
		DisplayName: nil,
	}
	var res integrationtype.GetIntegrationTypeResponse
	handler := integrationtype.UpdateIntegrationType(client)
	err := handler(context.Background(), req, &res)

	assert.Error(t, err)
}

func ptr(value string) *string {
	return &value
}
