package scanlabels_test

import (
	"context"
	"github.com/shinobistack/gokakashi/ent/schema"
	"testing"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
	"github.com/shinobistack/gokakashi/ent/enttest"
	"github.com/shinobistack/gokakashi/internal/restapi/v1/scanlabels"
	"github.com/stretchr/testify/assert"
)

func TestCreateScanLabel_Valid(t *testing.T) {
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	defer client.Close()
	policy := client.Policies.Create().
		SetName("test-policy").
		SetImage(schema.Image{Registry: "test-registry", Name: "test-name", Tags: []string{"v1.0"}}).
		SetScanner("trivy").
		SaveX(context.Background())
	// Create a test scan
	scan := client.Scans.Create().
		SetPolicyID(policy.ID).
		SetImage("example-image:latest").
		SetScanner(policy.Scanner).
		SetStatus("scan_pending").
		SaveX(context.Background())

	req := scanlabels.CreateScanLabelRequest{
		ScanID: scan.ID,
		Key:    "env",
		Value:  "prod",
	}
	res := &scanlabels.CreateScanLabelResponse{}

	err := scanlabels.CreateScanLabel(client)(context.Background(), req, res)

	assert.NoError(t, err)
	assert.Equal(t, req.ScanID, res.ScanID)
	assert.Equal(t, req.Key, res.Key)
	assert.Equal(t, req.Value, res.Value)
}

func TestCreateScanLabel_MissingFields(t *testing.T) {
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	defer client.Close()

	policy := client.Policies.Create().
		SetName("test-policy").
		SetImage(schema.Image{Registry: "test-registry", Name: "test-name", Tags: []string{"v1.0"}}).
		SetScanner("trivy").
		SaveX(context.Background())

	req := scanlabels.CreateScanLabelRequest{
		ScanID: policy.ID,
		Key:    "",
		Value:  "",
	}
	res := &scanlabels.CreateScanLabelResponse{}

	err := scanlabels.CreateScanLabel(client)(context.Background(), req, res)

	assert.Error(t, err)
}

func TestCreateScanLabel_InvalidScanID(t *testing.T) {
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	defer client.Close()

	req := scanlabels.CreateScanLabelRequest{
		ScanID: uuid.Nil, // Invalid ID
		Key:    "env",
		Value:  "prod",
	}
	res := &scanlabels.CreateScanLabelResponse{}

	err := scanlabels.CreateScanLabel(client)(context.Background(), req, res)

	assert.Error(t, err)
}

func TestCreateScanLabel_DuplicateKey(t *testing.T) {
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	defer client.Close()

	policy := client.Policies.Create().
		SetName("test-policy").
		SetImage(schema.Image{Registry: "test-registry", Name: "test-name", Tags: []string{"v1.0"}}).
		SetScanner("trivy").
		SaveX(context.Background())

	// Create a test scan
	scan := client.Scans.Create().
		SetPolicyID(policy.ID).
		SetImage("example-image:latest").
		SetStatus("scan_pending").
		SetScanner(policy.Scanner).
		SaveX(context.Background())

	// Create a label
	client.ScanLabels.Create().
		SetScanID(scan.ID).
		SetKey("env").
		SetValue("prod").
		SaveX(context.Background())

	req := scanlabels.CreateScanLabelRequest{
		ScanID: scan.ID,
		Key:    "env",
		Value:  "staging",
	}
	res := &scanlabels.CreateScanLabelResponse{}

	err := scanlabels.CreateScanLabel(client)(context.Background(), req, res)

	assert.Error(t, err)
}

func TestCreateScanLabel_NonExistentScanID(t *testing.T) {
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	defer client.Close()

	req := scanlabels.CreateScanLabelRequest{
		ScanID: uuid.New(), // Non-existent ID
		Key:    "env",
		Value:  "prod",
	}
	res := &scanlabels.CreateScanLabelResponse{}

	err := scanlabels.CreateScanLabel(client)(context.Background(), req, res)

	assert.Error(t, err)
}
