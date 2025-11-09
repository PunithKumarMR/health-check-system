package config

import (
    "fmt"
    "os"
    "strconv"
    "time"
)

type Config struct {
    Database DatabaseConfig
    App      AppConfig
}

type DatabaseConfig struct {
    Host     string
    Port     string
    User     string
    Password string
    Database string
}

type AppConfig struct {
    Environment         string
    LogLevel            string
    MaxConcurrentChecks int
    PollInterval        time.Duration
    MaxWait             time.Duration
}

func Load() (*Config, error) {
    cfg := &Config{
        Database: DatabaseConfig{
            Host:     getEnv("DB_HOST", "localhost"),
            Port:     getEnv("DB_PORT", "3306"),
            User:     getEnv("DB_USER", "root"),
            Password: getEnv("DB_PASSWORD", ""),
            Database: getEnv("DB_NAME", "mito_inventory"),
        },
        App: AppConfig{
            Environment:         getEnv("ENVIRONMENT", "development"),
            LogLevel:            getEnv("LOG_LEVEL", "INFO"),
            MaxConcurrentChecks: getEnvInt("MAX_CONCURRENT_CHECKS", 50),
            PollInterval:        getEnvDuration("HC_POLL_INTERVAL", 30*time.Second),
            MaxWait:             getEnvDuration("HC_MAX_WAIT", 80*time.Minute),
        },
    }

    if cfg.Database.Password == "" {
        return nil, fmt.Errorf("DB_PASSWORD is required")
    }

    return cfg, nil
}

func getEnv(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
    if value := os.Getenv(key); value != "" {
        if intVal, err := strconv.Atoi(value); err == nil {
            return intVal
        }
    }
    return defaultValue
}

func getEnvDuration(key string, defaultValue time.Duration) time.Duration {
    if value := os.Getenv(key); value != "" {
        if duration, err := time.ParseDuration(value); err == nil {
            return duration
        }
    }
    return defaultValue
}
