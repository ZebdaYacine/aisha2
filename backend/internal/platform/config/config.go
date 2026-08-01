package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	AppName            string
	Environment        string
	APIPort            string
	ShutdownTimeout    time.Duration
	DatabaseURL        string
	DatabaseMaxOpen    int32
	DatabaseMaxIdle    int32
	RedisURL           string
	MinIOEndpoint      string
	MinIOAccessKey     string
	MinIOSecretKey     string
	MinIOUseSSL        bool
	MinIOPublicBucket  string
	MinIOPrivateBucket string
	MinIOArtisanBucket string
	AllowedOrigins     []string
}

func Load() (Config, error) {
	maxOpen, err := positiveInt32("DATABASE_MAX_OPEN_CONNS", 20)
	if err != nil {
		return Config{}, err
	}
	maxIdle, err := positiveInt32("DATABASE_MAX_IDLE_CONNS", 5)
	if err != nil {
		return Config{}, err
	}
	useSSL, err := strconv.ParseBool(env("MINIO_USE_SSL", "false"))
	if err != nil {
		return Config{}, fmt.Errorf("MINIO_USE_SSL must be true or false: %w", err)
	}
	cfg := Config{
		AppName:            env("APP_NAME", "AISHA"),
		Environment:        env("APP_ENV", "development"),
		APIPort:            env("API_PORT", "8080"),
		ShutdownTimeout:    10 * time.Second,
		DatabaseURL:        env("DATABASE_URL", "postgres://aisha:aisha_dev@localhost:5432/aisha?sslmode=disable"),
		DatabaseMaxOpen:    maxOpen,
		DatabaseMaxIdle:    maxIdle,
		RedisURL:           env("REDIS_URL", "redis://localhost:6379/0"),
		MinIOEndpoint:      env("MINIO_ENDPOINT", "localhost:9000"),
		MinIOAccessKey:     env("MINIO_ACCESS_KEY", "aisha"),
		MinIOSecretKey:     env("MINIO_SECRET_KEY", "aisha_minio_dev"),
		MinIOUseSSL:        useSSL,
		MinIOPublicBucket:  env("MINIO_PUBLIC_BUCKET", "product-public"),
		MinIOPrivateBucket: env("MINIO_PRIVATE_BUCKET", "product-private"),
		MinIOArtisanBucket: env("MINIO_ARTISAN_BUCKET", "artisan-private"),
		AllowedOrigins:     splitCSV(env("CORS_ALLOWED_ORIGINS", "http://localhost:3000")),
	}
	if cfg.Environment == "production" {
		if strings.Contains(cfg.DatabaseURL, "aisha_dev") || cfg.MinIOSecretKey == "aisha_minio_dev" {
			return Config{}, fmt.Errorf("development credentials are forbidden in production")
		}
	}
	return cfg, nil
}

func positiveInt32(key string, fallback int32) (int32, error) {
	value, err := strconv.ParseInt(env(key, strconv.FormatInt(int64(fallback), 10)), 10, 32)
	if err != nil || value < 1 {
		return 0, fmt.Errorf("%s must be a positive integer", key)
	}
	return int32(value), nil
}

func splitCSV(value string) []string {
	values := strings.Split(value, ",")
	result := make([]string, 0, len(values))
	for _, item := range values {
		if trimmed := strings.TrimSpace(item); trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

func env(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
