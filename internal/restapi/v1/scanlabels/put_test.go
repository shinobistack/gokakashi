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

func TestUpdateScanLabel_Valid(t *testing.T) {
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
		SetScanner(policy.Scanner).
		SetStatus("scan_pending").
		SetIntegrationID(integrations.ID).
		SaveX(context.Background())

	// Add a label
	client.ScanLabels.Create().
		SetScanID(scan.ID).
		SetKey("env").
		SetValue("prod").
		SaveX(context.Background())

	req := scanlabels.UpdateScanLabelRequest{
		ScanID: scan.ID,
		Key:    strPtr("env"),
		Value:  strPtr("dev"),
	}
	res := &scanlabels.UpdateScanLabelResponse{}

	err := scanlabels.UpdateScanLabel(client)(context.Background(), req, res)

	assert.NoError(t, err)
	assert.Equal(t, "dev", res.Labels[0].Value)
}

func TestUpdateScanLabel_MissingFields(t *testing.T) {
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

	scan := client.Scans.Create().
		SetPolicyID(policy.ID).
		SetImage("example-image:latest").
		SetScanner(policy.Scanner).
		SetStatus("scan_pending").
		SetIntegrationID(integrations.ID).
		SaveX(context.Background())

	client.ScanLabels.Create().
		SetScanID(scan.ID).
		SetKey("env").
		SetValue("prod").
		SaveX(context.Background())

	req := scanlabels.UpdateScanLabelRequest{
		ScanID: scan.ID,
		Key:    strPtr(""),
		Value:  strPtr(""),
	}
	res := &scanlabels.UpdateScanLabelResponse{}

	err := scanlabels.UpdateScanLabel(client)(context.Background(), req, res)

	assert.Error(t, err)
}

func TestUpdateScanLabel_LabelNotFound(t *testing.T) {
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

	scan := client.Scans.Create().
		SetPolicyID(policy.ID).
		SetImage("example-image:latest").
		SetStatus("scan_pending").
		SetScanner(policy.Scanner).
		SetIntegrationID(integrations.ID).
		SaveX(context.Background())

	client.ScanLabels.Create().
		SetScanID(scan.ID).
		SetKey("env").
		SetValue("prod").
		SaveX(context.Background())

	req := scanlabels.UpdateScanLabelRequest{
		ScanID: scan.ID,
		Key:    strPtr("nonexistent-key"),
		Value:  strPtr("new-value"),
	}
	res := &scanlabels.UpdateScanLabelResponse{}

	err := scanlabels.UpdateScanLabel(client)(context.Background(), req, res)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "label not found")
}

func TestUpdateScanLabel_InvalidScanID(t *testing.T) {
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	defer client.Close()

	req := scanlabels.UpdateScanLabelRequest{
		ScanID: uuid.Nil, // Invalid ScanID
		Key:    strPtr("env"),
		Value:  strPtr("new-value"),
	}
	res := &scanlabels.UpdateScanLabelResponse{}

	err := scanlabels.UpdateScanLabel(client)(context.Background(), req, res)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid UUID")
}

func strPtr(s string) *string {
	return &s
}
