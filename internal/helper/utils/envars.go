package utils

import (
	"os"

	"github.com/joho/godotenv"
)

func logFailure(key string, sources []string) {
	logger := NewLogger()
	logger.Errorw("Failed to get key-value pair", "failed sources", sources, "key", key)
}

func GetEnv(key string) string {
	var failedSources []string

	if value := getOSEnv(key, &failedSources); value != "" {
		return value
	}

	if value := getDotEnv(key, &failedSources); value != "" {
		return value
	}

	logFailure(key, failedSources)
	return ""
}

func getOSEnv(key string, failedSources *[]string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	*failedSources = append(*failedSources, "os")
	return ""
}

func getDotEnv(key string, failedSources *[]string) string {
	if err := godotenv.Load("../../.env"); err != nil {
		*failedSources = append(*failedSources, ".env")
		return ""
	}
	return getOSEnv(key, failedSources)
}
