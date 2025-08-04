package scans

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/shinobistack/gokakashi/ent"
	"github.com/shinobistack/gokakashi/internal/restapi/v2/io"
	"github.com/shinobistack/gokakashi/internal/scan/v2"
	"github.com/swaggest/usecase/status"
)

func GetScan(client *ent.Client) func(ctx context.Context, req io.ScanGetRequest, res *io.ScanGetResponse) error {
	return func(ctx context.Context, req io.ScanGetRequest, res *io.ScanGetResponse) error {
		if req.ID == uuid.Nil {
			return status.Wrap(errors.New("invalid ID: cannot be nil"), status.InvalidArgument)
		}

		s, err := client.V2Scans.Get(ctx, req.ID)
		if err != nil {
			if ent.IsNotFound(err) {
				return status.Wrap(errors.New("scan not found"), status.NotFound)
			}
			return status.Wrap(fmt.Errorf("unexpected error: %v", err), status.Internal)
		}

		*res = io.ScanGetResponse{
			Scan: io.Scan{
				ID:        s.ID,
				Image:     s.Image,
				Labels:    s.Labels,
				Status:    scan.Status(s.Status),
				CreatedAt: s.CreatedAt,
				UpdatedAt: s.UpdatedAt,
			},
		}
		return nil
	}
}
