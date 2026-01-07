package config

import (
    "os"
    "github.com/joho/godotenv"
)

type Config struct {
    Port         string
    DatabaseURL  string
    RedisURL     string
    JWTSecret    string
    Environment  string
}

func Load() (*Config, error) {
    godotenv.Load()

    return &Config{
        Port:        getEnv("PORT", "8080"),
        DatabaseURL: getEnv("DATABASE_URL", ""),
        RedisURL:    getEnv("REDIS_URL", "localhost:6379"),
        JWTSecret:   getEnv("JWT_SECRET", "your-secret-key"),
        Environment: getEnv("ENVIRONMENT", "development"),
    }, nil
}

func getEnv(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
}