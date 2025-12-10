package main

import (
	"encoding/json"
	"net/http"
)

// This is the MAIN file with general API info
// All other files should NOT have these annotations

// @title E-Commerce API
// @version 2.0
// @description Complete e-commerce API with products and orders
// @termsOfService http://example.com/terms/

// @contact.name API Support
// @contact.url http://example.com/support
// @contact.email support@example.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host api.example.com
// @BasePath /v2

func main() {
	http.HandleFunc("/v2/health", HealthCheck)
	http.ListenAndServe(":8080", nil)
}

// HealthCheck returns API health status
// @Summary Health check
// @Description Check if API is running
// @Tags system
// @Produce json
// @Success 200 {object} map[string]string
// @Router /health [get]
func HealthCheck(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}
