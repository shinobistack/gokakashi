package client

import (
	"context"

	"github.com/shinobistack/gokakashi/internal/restapi/v2/io"
)

type AgentService service

type AgentRegisterRequest io.AgentRegisterRequest
type AgentRegisterResponse io.AgentRegisterResponse

func (s *AgentService) Register(ctx context.Context, _ *AgentRegisterRequest) (*AgentRegisterResponse, error) {
	req, err := s.client.NewRequest("POST", "api/v2/agents", nil)
	if err != nil {
		return nil, err
	}

	var resp AgentRegisterResponse
	if err := s.client.Do(ctx, req, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}
