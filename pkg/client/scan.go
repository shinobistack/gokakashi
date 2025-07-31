package client

import (
	"context"

	"github.com/shinobistack/gokakashi/internal/restapi/v2/io"
)

type ScanService service

type ScanCreateRequest io.ScanCreateRequest
type ScanCreateResponse io.ScanCreateResponse

func (s *ScanService) Create(ctx context.Context, reqData *ScanCreateRequest) (*ScanCreateResponse, error) {
	req, err := s.client.NewRequest("POST", "api/v2/scans", reqData)
	if err != nil {
		return nil, err
	}

	var resp ScanCreateResponse
	if err := s.client.Do(ctx, req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
