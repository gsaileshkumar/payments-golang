package main

import (
	"os"
)

type Config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
}

func GetEnvOrDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func LoadConfig() *Config {
	return &Config{
		DBHost:     GetEnvOrDefault("DB_HOST", "localhost"),
		DBPort:     GetEnvOrDefault("DB_PORT", "5432"),
		DBUser:     GetEnvOrDefault("DB_USER", "postgres"),
		DBPassword: GetEnvOrDefault("DB_PASSWORD", "postgres"),
		DBName:     GetEnvOrDefault("DB_NAME", "payments"),
	}
}
