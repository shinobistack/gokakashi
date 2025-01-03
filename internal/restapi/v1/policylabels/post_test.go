package policylabels_test

import (
	"context"
	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
	"github.com/shinobistack/gokakashi/ent/enttest"
	"github.com/shinobistack/gokakashi/ent/schema"
	"github.com/shinobistack/gokakashi/internal/restapi/v1/policylabels"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreatePolicyLabel_Valid(t *testing.T) {
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	defer client.Close()

	policy := client.Policies.Create().
		SetName("to-be-deleted-test-policy").
		SetImage(schema.Image{Registry: "example-registry", Name: "example-name", Tags: []string{"v1.0"}}).
		SetScanner("trivy").
		SaveX(context.Background())

	req := policylabels.CreatePolicyLabelRequest{
		PolicyID: policy.ID,
		Key:      "env",
		Value:    "prod",
	}
	res := &policylabels.CreatePolicyLabelResponse{}

	err := policylabels.CreatePolicyLabel(client)(context.Background(), req, res)

	assert.NoError(t, err)
	assert.Equal(t, req.PolicyID, res.ID)
	assert.Len(t, res.Labels, 1)
	assert.Equal(t, req.Key, res.Labels[0].Key)
	assert.Equal(t, req.Value, res.Labels[0].Value)

}

func TestCreatePolicyLabel_Invalid(t *testing.T) {
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	defer client.Close()

	req := policylabels.CreatePolicyLabelRequest{PolicyID: uuid.Nil, Key: "", Value: ""}
	res := &policylabels.CreatePolicyLabelResponse{}

	err := policylabels.CreatePolicyLabel(client)(context.Background(), req, res)

	assert.Error(t, err)
}

func TestCreatePolicyLabel_Duplicate(t *testing.T) {
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

	req := policylabels.CreatePolicyLabelRequest{
		PolicyID: policy.ID,
		Key:      "env",
		Value:    "prod",
	}
	res := &policylabels.CreatePolicyLabelResponse{}

	err := policylabels.CreatePolicyLabel(client)(context.Background(), req, res)

	assert.Error(t, err)
}
