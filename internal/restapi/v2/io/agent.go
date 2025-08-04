// io defines the input and output types for the API.
package io

import (
	"time"

	"github.com/google/uuid"
	agent "github.com/shinobistack/gokakashi/internal/agent/status/v2"
	"github.com/shinobistack/gokakashi/internal/agent/task"
)

type Pagination struct {
	Page    int `query:"page"`
	PerPage int `query:"per_page"`
}

type AgentRegisterRequest struct{}

type AgentRegisterResponse struct {
	ID uuid.UUID `json:"id"`

	Status          agent.Status `json:"status"`
	LastHeartbeatAt time.Time    `json:"last_heartbeat_at"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type AgentHeartbeatRequest struct {
	ID uuid.UUID `path:"agent_id"`
}

type AgentHeartbeatResponse struct {
	ID              uuid.UUID    `json:"id"`
	Status          agent.Status `json:"status"`
	LastHeartbeatAt time.Time    `json:"last_heartbeat_at"`
	CreatedAt       time.Time    `json:"created_at"`
	UpdatedAt       time.Time    `json:"updated_at"`
}

type AgentTask struct {
	ID uuid.UUID `json:"id"`

	ScanID  uuid.UUID   `json:"scan_id"`
	AgentID uuid.UUID   `json:"agent_id"`
	Status  task.Status `json:"status"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type AgentTaskListRequest struct {
	AgentID uuid.UUID `path:"agent_id"`
	Status  *string   `query:"status"`

	Pagination
}

type AgentTaskListResponse struct {
	Tasks      []AgentTask `json:"tasks"`
	Pagination Pagination   `json:"pagination"`
	Total      int          `json:"total"`
}
