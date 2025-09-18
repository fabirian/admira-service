package config

import (
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	AdsAPIURL   string
	CrmAPIURL   string
	SinkURL     string
	SinkSecret  string
	Port        string
	Timeout     time.Duration
	MaxRetries  int
	RetryDelay  time.Duration
}

func LoadConfig() (*Config, error) {
	// Cargar .env si existe
	godotenv.Load()

	return &Config{
		AdsAPIURL:   getEnv("ADS_API_URL", "https://api.mocky.io/v3/ads-uuid"),
		CrmAPIURL:   getEnv("CRM_API_URL", "https://api.mocky.io/v3/crm-uuid"),
		SinkURL:     getEnv("SINK_URL", ""),
		SinkSecret:  getEnv("SINK_SECRET", "admira_secret_example"),
		Port:        getEnv("PORT", "8080"),
		Timeout:     30 * time.Second,
		MaxRetries:  3,
		RetryDelay:  1 * time.Second,
	}, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}