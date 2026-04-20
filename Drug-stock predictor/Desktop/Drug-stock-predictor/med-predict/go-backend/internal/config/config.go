package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	// Server
	Port        int
	GinMode     string
	FrontendURL string

	// Database
	DBHost     string
	DBPort     int
	DBName     string
	DBUser     string
	DBPassword string
	DBSSLMode  string

	// JWT
	JWTSecret string

	// Logging
	LogLevel string
	LogDir   string

	// AI Services
	AnthropicAPIKey string
	OpenAIAPIKey    string

	// Notifications
	TwilioSID          string
	TwilioAuthToken    string
	TwilioWhatsAppNum  string
	MailgunDomain      string
	MailgunAPIKey      string

	// Environment
	Env string
}

// Load reads configuration from .env file and environment variables
func Load() (*Config, error) {
	// Load .env file (non-fatal if missing)
	_ = godotenv.Load()

	cfg := &Config{
		Port:        getEnvInt("PORT", 8080),
		GinMode:     getEnv("GIN_MODE", "debug"),
		FrontendURL: getEnv("FRONTEND_URL", "http://localhost:5173"),

		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnvInt("DB_PORT", 5432),
		DBName:     getEnv("DB_NAME", "medpredict"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", "postgres"),
		DBSSLMode:  getEnv("DB_SSL_MODE", "disable"),

		JWTSecret: getEnv("JWT_SECRET", "change-this-in-production"),

		LogLevel: getEnv("LOG_LEVEL", "info"),
		LogDir:   getEnv("LOG_DIR", "logs"),

		AnthropicAPIKey: getEnv("ANTHROPIC_API_KEY", ""),
		OpenAIAPIKey:    getEnv("OPENAI_API_KEY", ""),

		TwilioSID:         getEnv("TWILIO_SID", ""),
		TwilioAuthToken:   getEnv("TWILIO_AUTH_TOKEN", ""),
		TwilioWhatsAppNum: getEnv("TWILIO_WHATSAPP_NUMBER", ""),
		MailgunDomain:     getEnv("MAILGUN_DOMAIN", ""),
		MailgunAPIKey:     getEnv("MAILGUN_API_KEY", ""),

		Env: getEnv("ENV", "development"),
	}

	return cfg, nil
}

// GetDSN returns PostgreSQL connection string
func (c *Config) GetDSN() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.DBHost,
		c.DBPort,
		c.DBUser,
		c.DBPassword,
		c.DBName,
		c.DBSSLMode,
	)
}

func getEnv(key, defaultVal string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultVal
}

func getEnvInt(key string, defaultVal int) int {
	if value := os.Getenv(key); value != "" {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return defaultVal
}
