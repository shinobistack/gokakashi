package policies_test

import (
	"context"
	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
	"github.com/shinobistack/gokakashi/ent/enttest"
	"github.com/shinobistack/gokakashi/ent/schema"
	"github.com/shinobistack/gokakashi/internal/restapi/v1/policies"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUpdatePolicy_ValidUpdate(t *testing.T) {
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	defer client.Close()

	// Create a policy
	policy := client.Policies.Create().
		SetName("test-policy").
		SetImage(schema.Image{Registry: "example-registry", Name: "example-name", Tags: []string{"v1.0"}}).
		SaveX(context.Background())

	req := policies.UpdatePolicyRequest{
		ID:   policy.ID,
		Name: strPtr("updated-test-policy"),
	}
	var res policies.GetPolicyResponse
	err := policies.UpdatePolicy(client)(context.Background(), req, &res)

	assert.NoError(t, err)
	assert.Equal(t, policy.ID, res.ID)

	// Verify update
	updatedPolicy := client.Policies.GetX(context.Background(), policy.ID)
	assert.Equal(t, "updated-test-policy", updatedPolicy.Name)
}

func TestUpdatePolicy_InvalidUUID(t *testing.T) {
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	defer client.Close()

	req := policies.UpdatePolicyRequest{
		ID: uuid.Nil,
	}
	var res policies.GetPolicyResponse
	err := policies.UpdatePolicy(client)(context.Background(), req, &res)

	assert.Error(t, err)
}

func TestUpdatePolicy_NonExistentID(t *testing.T) {
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	defer client.Close()

	req := policies.UpdatePolicyRequest{
		ID: uuid.New(),
	}
	var res policies.GetPolicyResponse
	err := policies.UpdatePolicy(client)(context.Background(), req, &res)

	assert.Error(t, err)
}

// ToDo: implement no changes in policy when no to be changed objects are passed
//func TestUpdatePolicy_NoChanges(t *testing.T) {
//	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
//	defer client.Close()
//
//	// Create a policy
//	policy := client.Policies.Create().
//		SetName("test-policy").
//		SetImage(schema.Image{Registry: "example-registry", Name: "example-name", Tags: []string{"v1.0"}}).
//		SaveX(context.Background())
//
//	req := policies.UpdatePolicyRequest{
//		ID: policy.ID,
//	}
//	var res policies.GetPolicyResponse
//	err := policies.UpdatePolicy(client)(context.Background(), req, &res)
//
//	assert.NoError(t, err)
//	assert.Equal(t, "No Change Policy", client.Policies.GetX(context.Background(), policy.ID).Name)
//}

func strPtr(s string) *string {
	return &s
}
