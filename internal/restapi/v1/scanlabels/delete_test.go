package scanlabels_test

import (
	"context"
	"github.com/shinobistack/gokakashi/ent/schema"
	"testing"

	"github.com/google/uuid"
	"github.com/shinobistack/gokakashi/ent/enttest"
	"github.com/shinobistack/gokakashi/internal/restapi/v1/scanlabels"
	"github.com/stretchr/testify/assert"
)

func TestDeleteScanLabel_Valid(t *testing.T) {
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	defer client.Close()
	policy := client.Policies.Create().
		SetName("test-policy").
		SetImage(schema.Image{Registry: "test-registry", Name: "test-name", Tags: []string{"v1.0"}}).
		SaveX(context.Background())

	scan := client.Scans.Create().
		SetPolicyID(policy.ID).
		SetImage("example-image:latest").
		SetStatus("scan_pending").
		SaveX(context.Background())

	// Add a label
	client.ScanLabels.Create().
		SetScanID(scan.ID).
		SetKey("env").
		SetValue("prod").
		SaveX(context.Background())

	req := scanlabels.DeleteScanLabelRequest{
		ScanID: scan.ID,
		Key:    "env",
	}
	res := &scanlabels.DeleteScanLabelResponse{}

	err := scanlabels.DeleteScanLabel(client)(context.Background(), req, res)

	assert.NoError(t, err)
	assert.Equal(t, "deleted", res.Status)
}

func TestDeleteScanLabel_MissingFields(t *testing.T) {
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	defer client.Close()

	policy := client.Policies.Create().
		SetName("test-policy").
		SetImage(schema.Image{Registry: "test-registry", Name: "test-name", Tags: []string{"v1.0"}}).
		SaveX(context.Background())

	scan := client.Scans.Create().
		SetPolicyID(policy.ID).
		SetImage("example-image:latest").
		SetStatus("scan_pending").
		SaveX(context.Background())

	req := scanlabels.DeleteScanLabelRequest{
		ScanID: scan.ID, // Missing ScanID
		Key:    "",
	}
	res := &scanlabels.DeleteScanLabelResponse{}

	err := scanlabels.DeleteScanLabel(client)(context.Background(), req, res)

	assert.Error(t, err)
}

func TestDeleteScanLabel_LabelNotFound(t *testing.T) {
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	defer client.Close()
	policy := client.Policies.Create().
		SetName("test-policy").
		SetImage(schema.Image{Registry: "test-registry", Name: "test-name", Tags: []string{"v1.0"}}).
		SaveX(context.Background())

	// Create a test scan
	scan := client.Scans.Create().
		SetPolicyID(policy.ID).
		SetImage("example-image:latest").
		SetStatus("scan_pending").
		SaveX(context.Background())

	client.ScanLabels.Create().
		SetScanID(scan.ID).
		SetKey("env").
		SetValue("prod").
		SaveX(context.Background())

	req := scanlabels.DeleteScanLabelRequest{
		ScanID: scan.ID,
		Key:    "nonexistent-key",
	}
	res := &scanlabels.DeleteScanLabelResponse{}

	err := scanlabels.DeleteScanLabel(client)(context.Background(), req, res)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "label not found")
}

func TestDeleteScanLabel_InvalidScanID(t *testing.T) {
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	defer client.Close()

	req := scanlabels.DeleteScanLabelRequest{
		ScanID: uuid.Nil, // Invalid ScanID
		Key:    "env",
	}
	res := &scanlabels.DeleteScanLabelResponse{}

	err := scanlabels.DeleteScanLabel(client)(context.Background(), req, res)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid UUID")
}
