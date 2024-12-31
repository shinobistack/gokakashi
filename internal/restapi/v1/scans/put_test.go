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

func TestUpdateScan_Valid(t *testing.T) {
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

	// Create a test scan
	scan := client.Scans.Create().
		SetPolicyID(policy.ID).
		SetImage("example-image:latest").
		SetStatus("scan_pending").
		SetScanner(policy.Scanner).
		SetIntegrationID(integrations.ID).
		SaveX(context.Background())

	req := scans.UpdateScanRequest{
		ID:     scan.ID,
		Status: strPtr("scan_in_progress"),
		Report: "https://reports.server.com/scan/123",
	}

	res := &scans.UpdateScanResponse{}
	err := scans.UpdateScan(client)(context.Background(), req, res)

	assert.NoError(t, err)
	assert.Equal(t, req.ID, res.ID)
	assert.Equal(t, "scan_in_progress", res.Status)
}

func TestUpdateScan_MissingFields(t *testing.T) {
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	defer client.Close()

	policy := client.Policies.Create().
		SetName("test-policy").
		SetImage(schema.Image{Registry: "test-registry", Name: "test-name", Tags: []string{"v1.0"}}).
		SetScanner("trivy").
		SaveX(context.Background())

	req := scans.UpdateScanRequest{
		ID: policy.ID,
		// Missing required fields
	}

	res := &scans.UpdateScanResponse{}
	err := scans.UpdateScan(client)(context.Background(), req, res)

	assert.Error(t, err)
}

func TestUpdateScan_InvalidScanID(t *testing.T) {
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	defer client.Close()

	req := scans.UpdateScanRequest{
		ID:     uuid.Nil,
		Status: strPtr("in_progress"),
	}

	res := &scans.UpdateScanResponse{}
	err := scans.UpdateScan(client)(context.Background(), req, res)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid UUID")
}

func TestUpdateScan_NonExistentScanID(t *testing.T) {
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	defer client.Close()

	req := scans.UpdateScanRequest{
		ID:     uuid.New(), // Non-existent ID
		Status: strPtr("in_progress"),
	}

	res := &scans.UpdateScanResponse{}
	err := scans.UpdateScan(client)(context.Background(), req, res)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

func strPtr(s string) *string {
	return &s
}
