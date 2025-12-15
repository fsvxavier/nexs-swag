package utils

import "time"

// Helper represents a utility helper
type Helper struct {
	Name      string    `json:"name" example:"Logger"`
	Version   string    `json:"version" example:"1.0.0"`
	CreatedAt time.Time `json:"created_at"`
}

// HealthCheck represents health check response
type HealthCheck struct {
	Status  string `json:"status" example:"healthy"`
	Uptime  int64  `json:"uptime" example:"3600"`
	Version string `json:"version" example:"1.0.0"`
}

// GetHelpers retrieves available helpers
// @Summary Get helpers (pkg level)
// @Description Retrieves list of available utility helpers
// @Tags helpers
// @Produce json
// @Success 200 {array} Helper
// @Router /pkg/helpers [get]
func GetHelpers() {}

// HealthCheckHandler returns health status
// @Summary Health check endpoint
// @Description Returns the health status of the application
// @Tags system
// @Produce json
// @Success 200 {object} HealthCheck
// @Router /health [get]
func HealthCheckHandler() {}
