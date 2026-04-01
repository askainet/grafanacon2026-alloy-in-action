package config

import (
	"os"
	"strconv"
)

type Config struct {
	// Server
	Port string

	// OpenTelemetry
	OTelEndpoint string
	ServiceName  string

	// Logging
	LogLevel string

	// S3
	S3Endpoint string

	// Loki
	LokiEndpoint string

	// Tempo
	TempoEndpoint string

	// Mimir
	MimirEndpoint string

	// Alloy
	AlloyEndpoint string

	// Mission configurations
	Mission1MaxCardinality int
	Mission1GrowthRate     int
	Mission4SecretMessage  string
}

func Load() (*Config, error) {
	return &Config{
		Port:         getEnv("PORT", "8080"),
		OTelEndpoint: getEnv("OTEL_ENDPOINT", "alloy:4318"),
		ServiceName:  getEnv("SERVICE_NAME", "mission-control"),
		LogLevel:     "debug",
		S3Endpoint:    getEnv("S3_ENDPOINT", "http://localstack:4566"),
		LokiEndpoint:  getEnv("LOKI_ENDPOINT", "http://loki:3100"),
		TempoEndpoint: getEnv("TEMPO_ENDPOINT", "http://tempo:3200"),
		MimirEndpoint: getEnv("MIMIR_ENDPOINT", "http://mimir:9009"),
		AlloyEndpoint: getEnv("ALLOY_ENDPOINT", "http://alloy:12345"),

		Mission1MaxCardinality: getEnvInt("MISSION1_MAX_CARDINALITY", 10000),
		Mission1GrowthRate:     getEnvInt("MISSION1_GROWTH_RATE", 50),
		Mission4SecretMessage:  getEnv("MISSION4_SECRET", "ALPHA_BRAVO_CHARLIE_DELTA_ECHO"),
	}, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}
