package config

// Settings represents application configuration
type Settings struct {
	APIKey      string `json:"api_key" example:"secret_key_123"`
	Debug       bool   `json:"debug" example:"true"`
	Port        int    `json:"port" example:"8080"`
	DatabaseURL string `json:"database_url" example:"postgresql://localhost/db"`
}

// DatabaseConfig represents database configuration
type DatabaseConfig struct {
	Host     string `json:"host" example:"localhost"`
	Port     int    `json:"port" example:"5432"`
	User     string `json:"user" example:"admin"`
	Password string `json:"password" example:"secret"`
	DBName   string `json:"db_name" example:"mydb"`
}

// GetSettings retrieves application settings
// @Summary Get settings (config package)
// @Description Retrieves current application configuration settings
// @Tags config
// @Produce json
// @Success 200 {object} Settings
// @Router /config/settings [get]
func GetSettings() {}

// GetDatabaseConfig retrieves database configuration
// @Summary Get database config
// @Description Retrieves database connection configuration
// @Tags config
// @Produce json
// @Success 200 {object} DatabaseConfig
// @Router /config/database [get]
func GetDatabaseConfig() {}
