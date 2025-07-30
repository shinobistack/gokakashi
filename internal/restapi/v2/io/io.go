// io defines the input and output types for the API.
package io

import (
	"time"

	"github.com/google/uuid"
	agent "github.com/shinobistack/gokakashi/internal/agent/status/v2"
)

type AgentRegisterRequest struct{}

type AgentRegisterResponse struct {
	ID uuid.UUID `json:"id"`

	Status          agent.Status `json:"status"`
	LastHeartbeatAt time.Time    `json:"last_heartbeat_at"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
