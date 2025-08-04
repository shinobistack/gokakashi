package client

import (
	"context"

	"github.com/shinobistack/gokakashi/internal/restapi/v2/io"
)

type AgentService service

type AgentRegisterRequest io.AgentRegisterRequest
type AgentRegisterResponse io.AgentRegisterResponse

type AgentHeartbeatRequest io.AgentHeartbeatRequest
type AgentHeartbeatResponse io.AgentHeartbeatResponse

type AgentTaskListRequest io.AgentTaskListRequest
type AgentTaskListResponse io.AgentTaskListResponse

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

// Heartbeat sends a heartbeat for the agent with the given ID.
func (s *AgentService) Heartbeat(ctx context.Context, reqData *AgentHeartbeatRequest) (*AgentHeartbeatResponse, error) {
	req, err := s.client.NewRequest("PATCH", "api/v2/agents/"+reqData.ID.String()+"/heartbeat", nil)
	if err != nil {
		return nil, err
	}

	var resp AgentHeartbeatResponse
	if err := s.client.Do(ctx, req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

func (s *AgentService) ListTasks(ctx context.Context, reqData *AgentTaskListRequest) (*AgentTaskListResponse, error) {
	req, err := s.client.NewRequest("GET", "api/v2/agents/"+reqData.AgentID.String()+"/tasks", nil)
	if err != nil {
		return nil, err
	}

	var resp AgentTaskListResponse
	if err := s.client.Do(ctx, req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
