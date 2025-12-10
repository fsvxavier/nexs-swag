package internal

// Config is internal configuration.
type Config struct {
	APIKey string `json:"api_key"`
	Secret string `json:"secret"`
}

// GetConfig returns internal config
// @Summary Get internal config
// @Tags internal
// @Success 200 {object} Config
// @Router /internal/config [get].
func GetConfig() {}
