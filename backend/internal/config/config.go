package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Env            string
	APIPort        int
	RelayerPort    int
	DatabaseURL    string
	RedisURL       string
	TelcoHMACKey   string
	RelayerPrivKey string
}

func Load() (*Config, error) {
	// Load .env file if it exists (silently fails in production)
	_ = godotenv.Load()

	apiPort, err := strconv.Atoi(getEnvOrDefault("API_PORT", "8080"))
	if err != nil {
		return nil, fmt.Errorf("invalid API_PORT: %w", err)
	}

	relayerPort, err := strconv.Atoi(getEnvOrDefault("RELAYER_PORT", "8081"))
	if err != nil {
		return nil, fmt.Errorf("invalid RELAYER_PORT: %w", err)
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		return nil, fmt.Errorf("DATABASE_URL is required")
	}

	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		return nil, fmt.Errorf("REDIS_URL is required")
	}

	return &Config{
		Env:            getEnvOrDefault("APP_ENV", "development"),
		APIPort:        apiPort,
		RelayerPort:    relayerPort,
		DatabaseURL:    dbURL,
		RedisURL:       redisURL,
		TelcoHMACKey:   os.Getenv("TELCO_HMAC_KEY"),
		RelayerPrivKey: os.Getenv("RELAYER_PRIVATE_KEY"),
	}, nil
}

func getEnvOrDefault(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}