package policies_test

import (
	"context"
	_ "github.com/mattn/go-sqlite3"
	"github.com/shinobistack/gokakashi/ent/enttest"
	"github.com/shinobistack/gokakashi/ent/schema"
	"github.com/shinobistack/gokakashi/internal/restapi/v1/policies"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreatepolicy_InvalidPolicyNameFormat(t *testing.T) {
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	defer client.Close()

	req := policies.CreatePolicyRequest{
		Name: "tTest-policy",
		Image: schema.Image{
			Registry: "example-registry",
			Name:     "example-name",
			Tags:     []string{"v1.0"},
		},
		Scanner: "trivy",
		Trigger: map[string]interface{}{"type": "cron", "schedule": "0 0 * * *"},
		Check: &schema.Check{
			Condition: "sev.high > 0",
			Notify:    []string{"team-slack"},
		},
	}
	res := &policies.CreatePolicyResponse{}
	err := policies.CreatePolicy(client)(context.Background(), req, res)
	assert.Error(t, err)

	req = policies.CreatePolicyRequest{
		Name: "#test policy",
		Image: schema.Image{
			Registry: "example-registry",
			Name:     "example-name",
			Tags:     []string{"v1.0"},
		},
		Scanner: "trivy",
		Trigger: map[string]interface{}{"type": "cron", "schedule": "0 0 * * *"},
		Check: &schema.Check{
			Condition: "sev.high > 0",
			Notify:    []string{"team-slack"},
		},
	}
	err = policies.CreatePolicy(client)(context.Background(), req, res)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid id format")
}

func TestCreatePolicy_ValidInput(t *testing.T) {
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	defer client.Close()

	req := policies.CreatePolicyRequest{
		Name: "test-policy",
		Image: schema.Image{
			Registry: "example-registry",
			Name:     "example-name",
			Tags:     []string{"v1.0"},
		},
		Scanner: "trivy",
		Trigger: map[string]interface{}{"type": "cron", "schedule": "0 0 * * *"},
		Check: &schema.Check{
			Condition: "sev.high > 0",
			Notify:    []string{"team-slack"},
		},
	}
	res := &policies.CreatePolicyResponse{}
	err := policies.CreatePolicy(client)(context.Background(), req, res)

	assert.NoError(t, err)
	assert.NotNil(t, res.ID)
	assert.Equal(t, "created", res.Status)
}

func TestCreatePolicy_MissingFields(t *testing.T) {
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	defer client.Close()

	req := policies.CreatePolicyRequest{
		Name: "incomplete-policy",
		Image: schema.Image{
			Registry: "",
			Name:     "example-name",
			Tags:     []string{},
		},
		Scanner: "trivy",
		Trigger: nil,
		Check:   nil,
	}
	res := &policies.CreatePolicyResponse{}
	err := policies.CreatePolicy(client)(context.Background(), req, res)

	assert.Error(t, err)
}

func TestCreatePolicy_DuplicateName(t *testing.T) {
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	defer client.Close()

	// Seed database
	req := policies.CreatePolicyRequest{
		Name: "duplicate-policy",
		Image: schema.Image{
			Registry: "example-registry",
			Name:     "example-name",
			Tags:     []string{"v1.0"},
		},
		Scanner: "trivy",
		Trigger: map[string]interface{}{"type": "cron", "schedule": "0 0 * * *"},
		Check: &schema.Check{
			Condition: "sev.high > 0",
			Notify:    []string{"team-slack"},
		},
	}
	res := &policies.CreatePolicyResponse{}
	handler := policies.CreatePolicy(client)
	err := handler(context.Background(), req, res)
	assert.NoError(t, err)

	// Test case: Duplicate name
	err = handler(context.Background(), req, res)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "policy with the same name already exists")
}
