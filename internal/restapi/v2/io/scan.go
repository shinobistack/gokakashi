package io

import (
	"time"

	"github.com/google/uuid"
	"github.com/shinobistack/gokakashi/internal/scan/v2"
)

type Scan struct {
	ID     uuid.UUID         `json:"id"`
	Image  string            `json:"image"`
	Labels map[string]string `json:"labels"`
	Status scan.Status       `json:"status"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ScanCreateRequest struct {
	Image  string            `json:"image"`
	Labels map[string]string `json:"labels"`
}

type ScanCreateResponse struct {
	Scan
}

type ScanGetRequest struct {
	ID uuid.UUID `path:"id"`
}

type ScanGetResponse struct {
	Scan
}
