package scans

import (
	"context"
	"errors"

	"github.com/shinobistack/gokakashi/ent"
	"github.com/shinobistack/gokakashi/internal/restapi/v2/io"
	"github.com/shinobistack/gokakashi/internal/scan/v2"
	"github.com/swaggest/usecase/status"
)

func CreateScan(client *ent.Client) func(ctx context.Context, req io.ScanCreateRequest, res *io.ScanCreateResponse) error {
	return func(ctx context.Context, req io.ScanCreateRequest, res *io.ScanCreateResponse) error {
		if req.Image == "" {
			return status.Wrap(errors.New("missing image field"), status.InvalidArgument)
		}

		scanCreate := client.V2Scans.Create().
			SetStatus(string(scan.Pending)).
			SetImage(req.Image)

		if len(req.Labels) > 0 {
			scanCreate.SetLabels(req.Labels)
		}

		s, err := scanCreate.Save(ctx)
		if err != nil {
			return status.Wrap(err, status.Internal)
		}

		*res = io.ScanCreateResponse{
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
