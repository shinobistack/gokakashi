package policylabels_test

import (
	"context"
	"github.com/shinobistack/gokakashi/ent/enttest"
	"github.com/shinobistack/gokakashi/ent/schema"
	"github.com/shinobistack/gokakashi/internal/restapi/v1/policylabels"
	"github.com/stretchr/testify/assert"
	"testing"
)

// Todo: Certainly can better this logic to check for single label update
func TestUpdatePolicyLabel_SingleLabel(t *testing.T) {
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	defer client.Close()

	policy := client.Policies.Create().
		SetName("to-be-deleted-test-policy").
		SetImage(schema.Image{Registry: "example-registry", Name: "example-name", Tags: []string{"v1.0"}}).
		SetScanner("trivy").
		SaveX(context.Background())

	client.PolicyLabels.Create().
		SetPolicyID(policy.ID).
		SetKey("team").
		SetValue("gokakashi").
		SaveX(context.Background())

	client.PolicyLabels.Create().
		SetPolicyID(policy.ID).
		SetKey("env").
		SetValue("prod").
		SaveX(context.Background())

	req := policylabels.UpdatePolicyLabelsRequest{
		PolicyID: policy.ID,
		Key:      StringPointer("env"),
		Value:    StringPointer("dev"),
	}
	res := &policylabels.UpdatePolicyLabelsResponse{}

	err := policylabels.UpdatePolicyLabels(client)(context.Background(), req, res)
	assert.NoError(t, err)
	assert.Equal(t, "dev", res.Labels[0].Value)
}

func TestUpdatePolicyLabel_AllLabels(t *testing.T) {
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	defer client.Close()

	policy := client.Policies.Create().
		SetName("to-be-deleted-test-policy").
		SetImage(schema.Image{Registry: "example-registry", Name: "example-name", Tags: []string{"v1.0"}}).
		SetScanner("trivy").
		SaveX(context.Background())

	client.PolicyLabels.Create().
		SetPolicyID(policy.ID).
		SetKey("env").
		SetValue("prod").
		SaveX(context.Background())
	client.PolicyLabels.Create().
		SetPolicyID(policy.ID).
		SetKey("team").
		SetValue("gokakashi").
		SaveX(context.Background())

	req := policylabels.UpdatePolicyLabelsRequest{
		PolicyID: policy.ID,
		Labels: []policylabels.PolicyLabel{
			{Key: "updated", Value: "all"},
		},
	}
	res := &policylabels.UpdatePolicyLabelsResponse{}

	err := policylabels.UpdatePolicyLabels(client)(context.Background(), req, res)

	assert.NoError(t, err)
	assert.Equal(t, "all", res.Labels[0].Value)
}

func StringPointer(s string) *string {
	return &s
}
