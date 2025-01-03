// Example test file for Scans - post_test.go
package scans_test

import (
	"context"
	"github.com/shinobistack/gokakashi/ent/schema"
	"testing"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
	"github.com/shinobistack/gokakashi/ent/enttest"
	"github.com/shinobistack/gokakashi/internal/restapi/v1/scans"
	"github.com/stretchr/testify/assert"
)

func TestCreateScan_ValidInput(t *testing.T) {
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	defer client.Close()

	policy := client.Policies.Create().
		SetName("to-be-deleted-test-policy").
		SetImage(schema.Image{Registry: "example-registry", Name: "example-name", Tags: []string{"v1.0"}}).
		SaveX(context.Background())

	req := scans.CreateScanRequest{
		PolicyID: policy.ID,
		Image:    "dockerhub/nginx:latest",
		Check:    schema.Check{Condition: "severity > high", Notify: []string{"email"}},
		Status:   "scan_pending",
		Report:   "",
	}
	res := &scans.CreateScanResponse{}

	err := scans.CreateScan(client)(context.Background(), req, res)

	assert.NoError(t, err)
	assert.NotNil(t, res.ID)
	assert.Equal(t, "scan_pending", res.Status)
}

func TestCreateScan_MissingFields(t *testing.T) {
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	defer client.Close()
	policy := client.Policies.Create().
		SetName("to-be-deleted-test-policy").
		SetImage(schema.Image{Registry: "example-registry", Name: "example-name", Tags: []string{"v1.0"}}).
		SaveX(context.Background())

	req := scans.CreateScanRequest{
		PolicyID: policy.ID,
		Image:    "",
		Check:    schema.Check{Condition: "severity > high", Notify: []string{"email"}},
		Status:   "scan_pending",
		Report:   "",
	}
	res := &scans.CreateScanResponse{}

	err := scans.CreateScan(client)(context.Background(), req, res)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "missing required fields")
}

func TestCreateScan_InvalidPolicyID(t *testing.T) {
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	defer client.Close()

	req := scans.CreateScanRequest{
		PolicyID: uuid.Nil,
		Image:    "dockerhub/nginx:latest",
		Check:    schema.Check{Condition: "severity > high", Notify: []string{"email"}},
		Status:   "scan_pending",
		Report:   "",
	}
	res := &scans.CreateScanResponse{}

	err := scans.CreateScan(client)(context.Background(), req, res)

	assert.Error(t, err)
}
