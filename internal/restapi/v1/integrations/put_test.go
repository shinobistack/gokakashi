package integrations_test

import (
	"context"
	"github.com/shinobistack/gokakashi/internal/restapi/v1/integrations"
	"testing"

	"github.com/google/uuid"
	"github.com/shinobistack/gokakashi/ent/enttest"
	"github.com/stretchr/testify/assert"
)

func TestUpdateIntegration(t *testing.T) {
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	defer client.Close()

	// Seed the database
	id := uuid.New()
	_, err := client.Integrations.Create().
		SetID(id).
		SetName("Old Integration").
		SetType("linear").
		SetConfig(map[string]interface{}{"key": "value"}).
		Save(context.Background())
	assert.NoError(t, err)

	// Test case: Valid update
	req := integrations.UpdateIntegrationRequest{
		ID:   id,
		Name: stringPointer("Updated Integration"),
	}
	var res integrations.GetIntegrationResponse
	handler := integrations.UpdateIntegration(client)
	err = handler(context.Background(), req, &res)

	assert.NoError(t, err)
	assert.NoError(t, err)
	assert.Equal(t, "Updated Integration", res.Name)
	assert.Equal(t, "linear", res.Type)
	assert.Equal(t, "value", res.Config["key"])

	// Verify database
	integration, err := client.Integrations.Get(context.Background(), id)
	assert.NoError(t, err)
	assert.Equal(t, "Updated Integration", integration.Name)
}

func TestUpdateIntegration_InvalidUUID(t *testing.T) {
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	defer client.Close()
	invalidUUID := "c60c700e-3dd9-4059-8372-f77235"
	parsedUUID, err := uuid.Parse(invalidUUID)

	req := integrations.UpdateIntegrationRequest{
		ID:   parsedUUID,
		Name: stringPointer("Invalid Update"),
	}
	var res integrations.GetIntegrationResponse

	handler := integrations.UpdateIntegration(client)
	err = handler(context.Background(), req, &res)

	assert.Error(t, err)
}

func TestUpdateIntegration_NonExistentID(t *testing.T) {
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	defer client.Close()

	req := integrations.UpdateIntegrationRequest{
		ID:   uuid.New(),
		Name: stringPointer("Non-existent ID"),
	}
	var res integrations.GetIntegrationResponse

	handler := integrations.UpdateIntegration(client)
	err := handler(context.Background(), req, &res)

	assert.Error(t, err)
}

func TestUpdateIntegration_NoChanges(t *testing.T) {
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	defer client.Close()

	// Seed the database
	id := uuid.New()
	_, err := client.Integrations.Create().
		SetID(id).
		SetName("Integration").
		SetType("linear").
		SetConfig(map[string]interface{}{"key": "value"}).
		Save(context.Background())
	assert.NoError(t, err)

	// Test case: No changes
	req := integrations.UpdateIntegrationRequest{
		ID: id,
	}
	var res integrations.GetIntegrationResponse

	handler := integrations.UpdateIntegration(client)
	err = handler(context.Background(), req, &res)

	assert.NoError(t, err)
	assert.Equal(t, "Integration", res.Name)
}

func stringPointer(s string) *string {
	return &s
}
