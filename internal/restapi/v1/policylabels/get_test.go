package policylabels_test

import (
	"context"
	"github.com/shinobistack/gokakashi/ent/enttest"
	"github.com/shinobistack/gokakashi/ent/schema"
	"github.com/shinobistack/gokakashi/internal/restapi/v1/policylabels"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestListPolicyLabels_Valid(t *testing.T) {
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

	req := policylabels.ListPolicyLabelsRequest{PolicyID: policy.ID}
	res := &policylabels.ListPolicyLabelsResponse{}

	err := policylabels.ListPolicyLabels(client)(context.Background(), req, res)

	assert.NoError(t, err)
	assert.Len(t, res.Labels, 1)
	assert.Equal(t, "env", res.Labels[0].Key)
	assert.Equal(t, "prod", res.Labels[0].Value)
}

func TestListPolicyLabels_NoLabels(t *testing.T) {
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	defer client.Close()

	policy := client.Policies.Create().
		SetName("to-be-deleted-test-policy").
		SetImage(schema.Image{Registry: "example-registry", Name: "example-name", Tags: []string{"v1.0"}}).
		SetScanner("trivy").
		SaveX(context.Background())
	req := policylabels.ListPolicyLabelsRequest{PolicyID: policy.ID}
	res := &policylabels.ListPolicyLabelsResponse{}

	err := policylabels.ListPolicyLabels(client)(context.Background(), req, res)

	assert.NoError(t, err)
	assert.Empty(t, res.Labels)
}

func TestGetPolicyLabel_SpecificLabel(t *testing.T) {
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

	req := policylabels.GetPolicyLabelRequest{PolicyID: policy.ID, Key: "team"}
	res := &policylabels.GetPolicyLabelResponse{}

	err := policylabels.GetPolicyLabel(client)(context.Background(), req, res)

	assert.NoError(t, err)
	assert.Equal(t, "team", res.Key)
	assert.Equal(t, "gokakashi", res.Value)

}

func TestGetPolicyLabel_NonExistentLabel(t *testing.T) {
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

	req := policylabels.GetPolicyLabelRequest{PolicyID: policy.ID, Key: "env"}
	res := &policylabels.GetPolicyLabelResponse{}

	err := policylabels.GetPolicyLabel(client)(context.Background(), req, res)

	assert.Error(t, err)
}
