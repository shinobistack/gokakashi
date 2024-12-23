package policylabels_test

import (
	"context"
	"github.com/shinobistack/gokakashi/ent/enttest"
	"github.com/shinobistack/gokakashi/ent/schema"
	"github.com/shinobistack/gokakashi/internal/restapi/v1/policylabels"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDeletePolicyLabel_Valid(t *testing.T) {
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	defer client.Close()

	policy := client.Policies.Create().
		SetName("to-be-deleted-test-policy").
		SetImage(schema.Image{Registry: "example-registry", Name: "example-name", Tags: []string{"v1.0"}}).
		SaveX(context.Background())

	client.PolicyLabels.Create().
		SetPolicyID(policy.ID).
		SetKey("team").
		SetValue("gokakashi").
		SaveX(context.Background())

	req := policylabels.DeletePolicyLabelRequest{
		PolicyID: policy.ID,
		Key:      "team",
		Value:    "gokakashi",
	}
	res := &policylabels.DeletePolicyLabelResponse{}

	err := policylabels.DeletePolicyLabel(client)(context.Background(), req, res)

	assert.NoError(t, err)
	assert.Equal(t, "deleted", res.Status)
}

func TestDeletePolicyLabel_NonExistentLabel(t *testing.T) {
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	defer client.Close()

	policy := client.Policies.Create().
		SetName("to-be-deleted-test-policy").
		SetImage(schema.Image{Registry: "example-registry", Name: "example-name", Tags: []string{"v1.0"}}).
		SaveX(context.Background())

	client.PolicyLabels.Create().
		SetPolicyID(policy.ID).
		SetKey("team").
		SetValue("gokakashi").
		SaveX(context.Background())

	req := policylabels.DeletePolicyLabelRequest{
		PolicyID: policy.ID,
		Key:      "dev",
		Value:    "role",
	}
	res := &policylabels.DeletePolicyLabelResponse{}

	err := policylabels.DeletePolicyLabel(client)(context.Background(), req, res)

	assert.Error(t, err)
}
