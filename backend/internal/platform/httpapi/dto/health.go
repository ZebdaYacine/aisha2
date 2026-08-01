package dto

import "github.com/aisha-platform/aisha/backend/internal/platform/health"

type HealthResponse struct {
	Status   string            `json:"status"`
	Services map[string]string `json:"services,omitempty"`
}

func HealthResponseFrom(result health.Result) HealthResponse {
	return HealthResponse{Status: result.Status, Services: result.Services}
}
