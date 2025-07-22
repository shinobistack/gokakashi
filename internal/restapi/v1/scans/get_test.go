package scans_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
	"github.com/shinobistack/gokakashi/ent/enttest"
	"github.com/shinobistack/gokakashi/ent/schema"
	"github.com/shinobistack/gokakashi/internal/restapi/v1/scans"
	"github.com/stretchr/testify/assert"
)

func TestListScans_Valid(t *testing.T) {
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	defer client.Close()

	policy := client.Policies.Create().
		SetName("test-policy").
		SetImage(schema.Image{Registry: "test-registry", Name: "test-name", Tags: []string{"v1.0"}}).
		SetScanner("trivy").
		SaveX(context.Background())
	integrations := client.Integrations.Create().
		SetName("Integration 1").
		SetType("docker-hub").
		SetConfig(map[string]interface{}{"key": "value1"}).
		SaveX(context.Background())

	client.Scans.Create().
		SetPolicyID(policy.ID).
		SetImage("test-image-1").
		SetStatus("scan_pending").
		SetScanner(policy.Scanner).
		SetIntegrationID(integrations.ID).
		SaveX(context.Background())

	client.Scans.Create().
		SetPolicyID(policy.ID).
		SetImage("test-image-2").
		SetScanner(policy.Scanner).
		SetStatus("success").
		SetIntegrationID(integrations.ID).
		SaveX(context.Background())

	req := scans.ListScanRequest{Status: "", Name: "", Page: 1, PerPage: 100}
	res := scans.ListScansResponse{}
	err := scans.ListScans(client)(context.Background(), req, &res)

	assert.NoError(t, err)
	assert.Equal(t, 2, len(res.Scans))
	assert.Equal(t, 1, res.Page)
	assert.Equal(t, 100, res.PerPage)
	assert.Equal(t, 2, res.Total)
	assert.Equal(t, 1, res.TotalPages)

	// Test default perPage (should be 30)
	resDefault := scans.ListScansResponse{}
	reqDefault := scans.ListScanRequest{Status: "", Name: "", Page: 1, PerPage: 0}
	err = scans.ListScans(client)(context.Background(), reqDefault, &resDefault)
	assert.NoError(t, err)
	assert.Equal(t, 30, resDefault.PerPage)

	// Test perPage > 100 is capped at 100
	resMax := scans.ListScansResponse{}
	reqMax := scans.ListScanRequest{Status: "", Name: "", Page: 1, PerPage: 999}
	err = scans.ListScans(client)(context.Background(), reqMax, &resMax)
	assert.NoError(t, err)
	assert.Equal(t, 100, resMax.PerPage)
}

func TestGetScan_Valid(t *testing.T) {
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	defer client.Close()

	policy := client.Policies.Create().
		SetName("test-policy").
		SetImage(schema.Image{Registry: "test-registry", Name: "test-name", Tags: []string{"v1.0"}}).
		SetScanner("trivy").
		SaveX(context.Background())
	integrations := client.Integrations.Create().
		SetName("Integration 1").
		SetType("docker-hub").
		SetConfig(map[string]interface{}{"key": "value1"}).
		SaveX(context.Background())

	scan := client.Scans.Create().
		SetPolicyID(policy.ID).
		SetImage("test-image").
		SetScanner(policy.Scanner).
		SetStatus("scan_pending").
		SetIntegrationID(integrations.ID).
		SaveX(context.Background())

	res := scans.GetScanResponse{}
	err := scans.GetScan(client)(context.Background(), scans.GetScanRequest{ID: scan.ID}, &res)

	assert.NoError(t, err)
	assert.Equal(t, scan.ID, res.ID)
	assert.Equal(t, "test-image", res.Image)
}

func TestGetScan_NonExistentID(t *testing.T) {
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	defer client.Close()

	res := scans.GetScanResponse{}
	err := scans.GetScan(client)(context.Background(), scans.GetScanRequest{ID: uuid.New()}, &res)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "scan not found")
}

func TestGetScan_InvalidID(t *testing.T) {
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	defer client.Close()

	res := scans.GetScanResponse{}
	err := scans.GetScan(client)(context.Background(), scans.GetScanRequest{ID: uuid.Nil}, &res)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid UUID: cannot be nil")
}
