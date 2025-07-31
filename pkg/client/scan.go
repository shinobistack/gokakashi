package client

import (
	"context"

	"github.com/shinobistack/gokakashi/internal/restapi/v2/io"
)

type ScanService service

type ScanCreateRequest io.ScanCreateRequest
type ScanCreateResponse io.ScanCreateResponse

type ScanGetRequest io.ScanGetRequest
type ScanGetResponse io.ScanGetResponse

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

func (s *ScanService) Get(ctx context.Context, reqData *ScanGetRequest) (*ScanGetResponse, error) {
	req, err := s.client.NewRequest("GET", "api/v2/scans/"+reqData.ID.String(), nil)
	if err != nil {
		return nil, err
	}

	var resp ScanGetResponse
	if err := s.client.Do(ctx, req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
