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

func TestDeleteScan_Valid(t *testing.T) {
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	defer client.Close()

	// Create a test scan
	policy := client.Policies.Create().
		SetName("test-policy").
		SetImage(schema.Image{Registry: "test-registry", Name: "test-name", Tags: []string{"v1.0"}}).
		SetScanner("trivy").
		SaveX(context.Background())
	integrations := client.Integrations.Create().
		SetName("integration").
		SetType("linear").
		SetConfig(map[string]interface{}{"key": "value"}).
		SaveX(context.Background())
	scan := client.Scans.Create().
		SetPolicyID(policy.ID).
		SetImage("example-image:latest").
		SetStatus("scan_pending").
		SetScanner(policy.Scanner).
		SetIntegrationID(integrations.ID).
		SaveX(context.Background())

	req := scans.DeleteScanRequest{
		ID: scan.ID,
	}

	res := &scans.DeleteScanResponse{}
	err := scans.DeleteScan(client)(context.Background(), req, res)

	assert.NoError(t, err)
	assert.Equal(t, req.ID, res.ID)
	assert.Equal(t, "deleted", res.Status)
}

func TestDeleteScan_InvalidScanID(t *testing.T) {
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	defer client.Close()

	req := scans.DeleteScanRequest{
		ID: uuid.Nil, // Invalid ID
	}

	res := &scans.DeleteScanResponse{}
	err := scans.DeleteScan(client)(context.Background(), req, res)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid UUID")
}

func TestDeleteScan_NonExistentScanID(t *testing.T) {
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	defer client.Close()

	req := scans.DeleteScanRequest{
		ID: uuid.New(), // Non-existent ID
	}

	res := &scans.DeleteScanResponse{}
	err := scans.DeleteScan(client)(context.Background(), req, res)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

// Todo: to use when we do transaction operation
//func TestDeleteScan_ErrorDuringDeletion(t *testing.T) {
//	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
//	defer client.Close()
//
//	req := scans.DeleteScanRequest{
//		ID: uuid.New(),
//	}
//
//	res := &scans.DeleteScanResponse{}
//	err := scans.DeleteScan(client)(context.Background(), req, res)
//
//	assert.Error(t, err)
//	assert.Contains(t, err.Error(), "failed to delete")
//}
