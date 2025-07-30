package scans

import (
	"context"
	"testing"

	_ "github.com/mattn/go-sqlite3"

	"github.com/shinobistack/gokakashi/ent/enttest"
	"github.com/shinobistack/gokakashi/internal/restapi/v2/io"
	"github.com/shinobistack/gokakashi/internal/scan/v2"
	"github.com/stretchr/testify/require"
)

func TestCreateScan_Success(t *testing.T) {
	db := enttest.Open(t, "sqlite3", "file:create_scan_success?mode=memory&cache=shared&_fk=1")
	defer db.Close()

	handler := CreateScan(db)
	ctx := context.Background()
	labels := map[string]string{"key1": "label1", "key2": "label2"}
	request := io.ScanCreateRequest{
		Image:  "alpine:latest",
		Labels: labels,
	}
	var response io.ScanCreateResponse

	err := handler(ctx, request, &response)
	require.NoError(t, err)
	require.NotZero(t, response.ID)
	require.Equal(t, scan.Pending, response.Status)

	saved, err := db.V2Scans.Get(ctx, response.ID)
	require.NoError(t, err)
	require.Equal(t, request.Image, saved.Image)
	require.Equal(t, labels, saved.Labels)
}

func TestCreateScan_MissingImage(t *testing.T) {
	db := enttest.Open(t, "sqlite3", "file:create_scan_missing_image?mode=memory&cache=shared&_fk=1")
	defer db.Close()

	handler := CreateScan(db)
	ctx := context.Background()
	request := io.ScanCreateRequest{
		Image:  "",
		Labels: map[string]string{"key1": "label1"},
	}
	var response io.ScanCreateResponse

	err := handler(ctx, request, &response)
	require.Error(t, err)
	require.Contains(t, err.Error(), "missing image field")
}

func TestCreateScan_DBError(t *testing.T) {
	// Simulate a DB error by closing the DB before calling handler
	db := enttest.Open(t, "sqlite3", "file:create_scan_dberror?mode=memory&cache=shared&_fk=1")
	db.Close()

	handler := CreateScan(db)
	ctx := context.Background()
	request := io.ScanCreateRequest{
		Image:  "alpine:latest",
	}
	var response io.ScanCreateResponse

	err := handler(ctx, request, &response)
	require.Error(t, err)
}
