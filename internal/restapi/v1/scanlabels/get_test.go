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

func TestListScanLabels_Valid(t *testing.T) {
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	defer client.Close()
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

	// Create a test scan
	scan := client.Scans.Create().
		SetPolicyID(policy.ID).
		SetImage("example-image:latest").
		SetScanner(policy.Scanner).
		SetStatus("scan_pending").
		SetIntegrationID(integrations.ID).
		SaveX(context.Background())

	// Add labels
	client.ScanLabels.Create().
		SetScanID(scan.ID).
		SetKey("env").
		SetValue("prod").
		SaveX(context.Background())

	client.ScanLabels.Create().
		SetScanID(scan.ID).
		SetKey("version").
		SetValue("v1.0").
		SaveX(context.Background())

	req := scanlabels.ListScanLabelsRequest{
		ScanID: scan.ID,
	}
	res := &scanlabels.ListScanLabelsResponse{}

	err := scanlabels.ListScanLabels(client)(context.Background(), req, res)

	assert.NoError(t, err)
	assert.Len(t, res.Labels, 2)
}

func TestListScanLabels_NoLabels(t *testing.T) {
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	defer client.Close()

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

	// Create a test scan
	scan := client.Scans.Create().
		SetPolicyID(policy.ID).
		SetImage("example-image:latest").
		SetScanner(policy.Scanner).
		SetStatus("scan_pending").
		SetIntegrationID(integrations.ID).
		SaveX(context.Background())

	req := scanlabels.ListScanLabelsRequest{
		ScanID: scan.ID,
	}
	res := &scanlabels.ListScanLabelsResponse{}

	err := scanlabels.ListScanLabels(client)(context.Background(), req, res)

	assert.NoError(t, err)
	assert.Len(t, res.Labels, 0)
}

func TestGetScanLabel_Valid(t *testing.T) {
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	defer client.Close()
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

	// Create a test scan
	scan := client.Scans.Create().
		SetPolicyID(policy.ID).
		SetImage("example-image:latest").
		SetScanner("trivy").
		SetStatus("scan_pending").
		SetIntegrationID(integrations.ID).
		SaveX(context.Background())

	// Add a label
	client.ScanLabels.Create().
		SetScanID(scan.ID).
		SetKey("env").
		SetValue("prod").
		SaveX(context.Background())

	req := scanlabels.GetScanLabelRequest{
		ScanID: scan.ID,
		Key:    "env",
	}
	res := &scanlabels.GetScanLabelResponse{}

	err := scanlabels.GetScanLabel(client)(context.Background(), req, res)

	assert.NoError(t, err)
	assert.Equal(t, "env", res.Key)
	assert.Equal(t, "prod", res.Value)
}

func TestGetScanLabel_NotFound(t *testing.T) {
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	defer client.Close()
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

	// Create a test scan
	scan := client.Scans.Create().
		SetPolicyID(policy.ID).
		SetImage("example-image:latest").
		SetScanner("trivy").
		SetStatus("scan_pending").
		SetIntegrationID(integrations.ID).
		SaveX(context.Background())

	req := scanlabels.GetScanLabelRequest{
		ScanID: scan.ID,
		Key:    "nonexistent-key",
	}
	res := &scanlabels.GetScanLabelResponse{}

	err := scanlabels.GetScanLabel(client)(context.Background(), req, res)

	assert.Error(t, err)
}

func TestGetScanLabel_InvalidScanID(t *testing.T) {
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	defer client.Close()

	req := scanlabels.GetScanLabelRequest{
		ScanID: uuid.Nil, // Invalid ScanID
		Key:    "env",
	}
	res := &scanlabels.GetScanLabelResponse{}

	err := scanlabels.GetScanLabel(client)(context.Background(), req, res)

	assert.Error(t, err)
}
