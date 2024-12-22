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

func TestListPolicies_EmptyDatabase(t *testing.T) {
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	defer client.Close()

	req := policies.ListPoliciesRequest{}
	res := &[]policies.GetPolicyResponse{}
	err := policies.ListPolicies(client)(context.Background(), req, res)

	assert.NoError(t, err)
	assert.Equal(t, 0, len(*res))
}

func TestListPolicies(t *testing.T) {
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	defer client.Close()

	// Create test data
	client.Policies.Create().
		SetName("test-policy1").
		SetImage(schema.Image{Registry: "example-registry", Name: "example-name", Tags: []string{"v1.0"}}).
		Save(context.Background())
	client.Policies.Create().
		SetName("test-policy2").
		SetImage(schema.Image{Registry: "example-registry", Name: "example-name", Tags: []string{"v1.0"}}).
		Save(context.Background())

	req := policies.ListPoliciesRequest{}
	res := &[]policies.GetPolicyResponse{}
	err := policies.ListPolicies(client)(context.Background(), req, res)

	assert.NoError(t, err)
	assert.Equal(t, 2, len(*res))
}

func TestGetPolicy_ValidID(t *testing.T) {
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	defer client.Close()

	// Create a policy
	policy := client.Policies.Create().
		SetName("test-policy").
		SetImage(schema.Image{Registry: "example-registry", Name: "example-name", Tags: []string{"v1.0"}}).
		SaveX(context.Background())

	req := policies.GetPolicyRequests{ID: policy.ID}
	res := &policies.GetPolicyResponse{}
	err := policies.GetPolicy(client)(context.Background(), req, res)

	assert.NoError(t, err)
	assert.Equal(t, policy.Name, res.Name)
}

func TestGetPolicy_InvalidUUID(t *testing.T) {
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	defer client.Close()

	req := policies.GetPolicyRequests{ID: uuid.Nil}
	res := &policies.GetPolicyResponse{}
	err := policies.GetPolicy(client)(context.Background(), req, res)

	assert.Error(t, err)
}

func TestGetPolicy_NonExistentID(t *testing.T) {
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	defer client.Close()

	req := policies.GetPolicyRequests{ID: uuid.New()}
	res := &policies.GetPolicyResponse{}
	err := policies.GetPolicy(client)(context.Background(), req, res)

	assert.Error(t, err)
}
