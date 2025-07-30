package io

import (
	"github.com/google/uuid"
	"github.com/shinobistack/gokakashi/internal/scan/v2"
)

type ScanCreateRequest struct {
	Image  string            `json:"image"`
	Labels map[string]string `json:"labels"`
}

type ScanCreateResponse struct {
	ID     uuid.UUID   `json:"id"`
	Status scan.Status `json:"status"`
}
