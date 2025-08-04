package scans

import (
	"context"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"

	"github.com/google/uuid"
	"github.com/shinobistack/gokakashi/ent/enttest"
	"github.com/shinobistack/gokakashi/internal/restapi/v2/io"
	"github.com/shinobistack/gokakashi/internal/scan/v2"
	"github.com/stretchr/testify/require"
)

func TestGetScan_Success(t *testing.T) {
	db := enttest.Open(t, "sqlite3", "file:get_scan_success?mode=memory&cache=shared&_fk=1")
	defer db.Close()

	ctx := context.Background()
	labels := map[string]string{"key1": "label1", "key2": "label2"}
	// Create a scan directly in the DB
	saved, err := db.V2Scans.Create().
		SetImage("alpine:latest").
		SetLabels(labels).
		SetStatus(string(scan.Pending)).
		Save(ctx)
	require.NoError(t, err)

	handler := GetScan(db)
	request := io.ScanGetRequest{ID: saved.ID}
	var response io.ScanGetResponse

	err = handler(ctx, request, &response)
	require.NoError(t, err)
	require.Equal(t, saved.ID, response.ID)
	require.Equal(t, saved.Image, response.Image)
	require.Equal(t, saved.Labels, response.Labels)
	require.Equal(t, scan.Status(saved.Status), response.Status)
	require.WithinDuration(t, saved.CreatedAt, response.CreatedAt, 2*time.Second)
	require.WithinDuration(t, saved.UpdatedAt, response.UpdatedAt, 2*time.Second)
}

func TestGetScan_NotFound(t *testing.T) {
	db := enttest.Open(t, "sqlite3", "file:get_scan_notfound?mode=memory&cache=shared&_fk=1")
	defer db.Close()

	handler := GetScan(db)
	ctx := context.Background()
	request := io.ScanGetRequest{ID: uuid.New()}
	var response io.ScanGetResponse

	err := handler(ctx, request, &response)
	require.Error(t, err)
	require.Contains(t, err.Error(), "not found")
}

func TestGetScan_InvalidID(t *testing.T) {
	db := enttest.Open(t, "sqlite3", "file:get_scan_invalidid?mode=memory&cache=shared&_fk=1")
	defer db.Close()

	handler := GetScan(db)
	ctx := context.Background()
	request := io.ScanGetRequest{ID: uuid.Nil}
	var response io.ScanGetResponse

	err := handler(ctx, request, &response)
	require.Error(t, err)
	require.Contains(t, err.Error(), "invalid ID")
}
