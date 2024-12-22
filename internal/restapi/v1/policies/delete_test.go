package policies_test

import (
	"context"
	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
	"github.com/shinobistack/gokakashi/ent/enttest"
	policyent "github.com/shinobistack/gokakashi/ent/policies"
	"github.com/shinobistack/gokakashi/ent/schema"
	"github.com/shinobistack/gokakashi/internal/restapi/v1/policies"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDeletePolicy_ValidDeletion(t *testing.T) {
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	defer client.Close()

	// Create a policy
	policy := client.Policies.Create().
		SetName("to-be-deleted-test-policy").
		SetImage(schema.Image{Registry: "example-registry", Name: "example-name", Tags: []string{"v1.0"}}).
		SaveX(context.Background())
	req := policies.DeletePolicyRequest{ID: policy.ID}
	res := &policies.DeletePolicyResponse{}
	err := policies.DeletePolicy(client)(context.Background(), req, res)

	assert.NoError(t, err)
	assert.Equal(t, policy.ID, res.ID)

	// Verify deletion
	exists := client.Policies.Query().Where(policyent.ID(policy.ID)).ExistX(context.Background())
	assert.False(t, exists)
}

func TestDeletePolicy_InvalidUUID(t *testing.T) {
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	defer client.Close()

	req := policies.DeletePolicyRequest{ID: uuid.Nil}
	res := &policies.DeletePolicyResponse{}
	err := policies.DeletePolicy(client)(context.Background(), req, res)

	assert.Error(t, err)
}

func TestDeletePolicy_NonExistentID(t *testing.T) {
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	defer client.Close()

	req := policies.DeletePolicyRequest{ID: uuid.New()}
	res := &policies.DeletePolicyResponse{}
	err := policies.DeletePolicy(client)(context.Background(), req, res)

	assert.Error(t, err)
}

// ToDo: TestDeletePolicy_WithDependentRecords
