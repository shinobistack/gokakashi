package integrations_test

import (
	"context"
	"github.com/shinobistack/gokakashi/internal/restapi/v1/integrations"
	"testing"

	"github.com/shinobistack/gokakashi/ent/enttest"
	"github.com/stretchr/testify/assert"
)

func TestCreateIntegration_ValidInput(t *testing.T) {
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	defer client.Close()

	// Test case: Valid input
	req := integrations.CreateIntegrationRequest{
		Name:   "Valid Integration",
		Type:   "linear",
		Config: map[string]interface{}{"key": "value"},
	}
	res := &integrations.CreateIntegrationResponse{}
	handler := integrations.CreateIntegration(client)
	err := handler(context.Background(), req, res)

	assert.NoError(t, err)
	assert.NotEmpty(t, res.ID)
	assert.Equal(t, "created", res.Status)

	// Verify database
	integration, err := client.Integrations.Get(context.Background(), res.ID)
	assert.NoError(t, err)
	assert.Equal(t, "Valid Integration", integration.Name)
	assert.Equal(t, "linear", integration.Type)
	assert.Equal(t, "value", integration.Config["key"])
}

func TestCreateIntegration_DuplicateName(t *testing.T) {
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	defer client.Close()

	// Seed database
	req := integrations.CreateIntegrationRequest{
		Name:   "Duplicate Integration",
		Type:   "jira",
		Config: map[string]interface{}{"key": "value"},
	}
	res := &integrations.CreateIntegrationResponse{}
	handler := integrations.CreateIntegration(client)
	err := handler(context.Background(), req, res)
	assert.NoError(t, err)

	// Test case: Duplicate name
	err = handler(context.Background(), req, res)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "integration with the same name already exists")
}

func TestCreateIntegration_MissingFields(t *testing.T) {
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	defer client.Close()

	// Test case: Missing fields
	req := integrations.CreateIntegrationRequest{
		Name:   "",
		Type:   "",
		Config: nil,
	}
	res := &integrations.CreateIntegrationResponse{}
	handler := integrations.CreateIntegration(client)
	err := handler(context.Background(), req, res)

	assert.Error(t, err)
}

// TODO: To intorduce this test case post writing Integration Types with right test case condition
//func TestCreateIntegration_InvalidType(t *testing.T) {
//	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
//	defer client.Close()
//
//	// Test case: Invalid type
//	req := integrations.CreateIntegrationRequest{
//		Name:   "Invalid Type Integration",
//		Type:   "nonexistent-type",
//		Config: map[string]interface{}{"key": "value"},
//	}
//	res := &integrations.CreateIntegrationResponse{}
//	handler := integrations.CreateIntegration(client)
//	err := handler(context.Background(), req, res)
//
//	assert.Error(t, err)
//}

func TestCreateIntegration_EmptyConfig(t *testing.T) {
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	defer client.Close()

	// Test case: Empty config
	req := integrations.CreateIntegrationRequest{
		Name:   "Empty Config Integration",
		Type:   "linear",
		Config: map[string]interface{}{},
	}
	res := &integrations.CreateIntegrationResponse{}
	handler := integrations.CreateIntegration(client)
	err := handler(context.Background(), req, res)

	assert.NoError(t, err)
	assert.NotEmpty(t, res.ID)

	// Verify database
	integration, err := client.Integrations.Get(context.Background(), res.ID)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(integration.Config))
}
