package health

import (
	"context"
	"fmt"
	"github.com/shinobistack/gokakashi/ent"
	"syscall"
	"time"
)

// HealthCheckRequest defines the request structure for the health check.
type HealthCheckRequest struct{}

// HealthCheckResponse defines the response structure for the health check.
type HealthCheckResponse struct {
	Status    string `json:"status"`
	Database  string `json:"database"`
	DiskSpace string `json:"disk_space"`
	Uptime    string `json:"uptime"`
}

// HealthCheck provides the health status of the application.
func HealthCheck(client *ent.Client, startTime time.Time) func(ctx context.Context, req HealthCheckRequest, res *HealthCheckResponse) error {
	return func(ctx context.Context, req HealthCheckRequest, res *HealthCheckResponse) error {
		*res = HealthCheckResponse{}

		// Check database connectivity
		if !client.Policies.Query().Limit(1).ExistX(ctx) {
			res.Database = "disconnected"
			res.Status = "unhealthy"
		} else {
			res.Database = "connected"
		}

		// Check disk space
		diskFree, diskTotal, err := getDiskSpace("/")
		if err != nil {
			res.DiskSpace = "unavailable"
			res.Status = "unhealthy"
		} else {
			res.DiskSpace = fmt.Sprintf("Free: %.2fGB / Total: %.2fGB", float64(diskFree)/1024, float64(diskTotal)/1024)
		}

		// Maybe Check environment variables

		//Calculate uptime
		res.Uptime = time.Since(startTime).String()

		// Finalize status
		if res.Status == "" {
			res.Status = "healthy"
		}

		return nil
	}
}

// Todo: this to remove, did out of thinking if this helps
func getDiskSpace(path string) (freeMB, totalMB int64, err error) {
	var stat syscall.Statfs_t
	if err := syscall.Statfs(path, &stat); err != nil {
		return 0, 0, err
	}

	// Calculate free and total space in MB
	freeMB = int64((stat.Bavail * uint64(stat.Bsize)) / (1024 * 1024))
	totalMB = int64((stat.Blocks * uint64(stat.Bsize)) / (1024 * 1024))
	return freeMB, totalMB, nil
}
