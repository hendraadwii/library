package config

import (
	"log"
	"os"
	"strconv"
	"time"
)

// DBConfig holds the database configuration
type DBConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
}

// AuthConfig holds the authentication configuration
type AuthConfig struct {
	AccessTokenSecret  string
	RefreshTokenSecret string
	AccessTokenExpiry  time.Duration
	RefreshTokenExpiry time.Duration
	Issuer             string
}

// Config holds the application configuration
type Config struct {
	Database     DBConfig
	Auth         AuthConfig
	Server       ServerConfig
	Environment  string
	RateLimiting RateLimitConfig
}

// ServerConfig holds the server configuration
type ServerConfig struct {
	Port         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	AllowOrigins []string
}

// RateLimitConfig holds rate limiting configuration
type RateLimitConfig struct {
	MaxRequests  int
	PerDuration  time.Duration
	CleanupEvery time.Duration
}

// LoadConfig loads the configuration from environment variables
func LoadConfig() *Config {
	config := &Config{
		Database: DBConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnvAsInt("DB_PORT", 3306),
			User:     getEnv("DB_USER", "root"),
			Password: getEnv("DB_PASSWORD", ""),
			DBName:   getEnv("DB_NAME", "library_db"),
		},
		Auth: AuthConfig{
			AccessTokenSecret:  getEnv("JWT_ACCESS_SECRET", "your-access-secret-key"),
			RefreshTokenSecret: getEnv("JWT_REFRESH_SECRET", "your-refresh-secret-key"),
			AccessTokenExpiry:  getEnvAsDuration("JWT_ACCESS_EXPIRY", 15*time.Minute),
			RefreshTokenExpiry: getEnvAsDuration("JWT_REFRESH_EXPIRY", 7*24*time.Hour),
			Issuer:             getEnv("JWT_ISSUER", "library-api"),
		},
		Server: ServerConfig{
			Port:         getEnv("SERVER_PORT", "8081"),
			ReadTimeout:  getEnvAsDuration("SERVER_READ_TIMEOUT", 10*time.Second),
			WriteTimeout: getEnvAsDuration("SERVER_WRITE_TIMEOUT", 10*time.Second),
			AllowOrigins: getEnvAsSlice("CORS_ALLOW_ORIGINS", []string{"http://localhost:8081", "http://localhost:8080", "http://localhost:8082", "http://localhost:3000"}),
		},
		Environment: getEnv("ENV", "development"),
		RateLimiting: RateLimitConfig{
			MaxRequests:  getEnvAsInt("RATE_LIMIT_MAX", 5),
			PerDuration:  getEnvAsDuration("RATE_LIMIT_DURATION", time.Minute),
			CleanupEvery: getEnvAsDuration("RATE_LIMIT_CLEANUP", 5*time.Minute),
		},
	}

	return config
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// getEnvAsInt gets an environment variable as an integer or returns a default value
func getEnvAsInt(key string, defaultValue int) int {
	valueStr := getEnv(key, "")
	if valueStr == "" {
		return defaultValue
	}

	value, err := strconv.Atoi(valueStr)
	if err != nil {
		log.Printf("Warning: Invalid value for %s, using default: %v\n", key, err)
		return defaultValue
	}
	return value
}

// getEnvAsDuration gets an environment variable as a duration or returns a default value
func getEnvAsDuration(key string, defaultValue time.Duration) time.Duration {
	valueStr := getEnv(key, "")
	if valueStr == "" {
		return defaultValue
	}

	value, err := time.ParseDuration(valueStr)
	if err != nil {
		log.Printf("Warning: Invalid duration for %s, using default: %v\n", key, err)
		return defaultValue
	}
	return value
}

// getEnvAsSlice gets an environment variable as a slice or returns a default value
func getEnvAsSlice(key string, defaultValue []string) []string {
	valueStr := getEnv(key, "")
	if valueStr == "" {
		return defaultValue
	}
	return []string{valueStr} // This is simplified, in a real application you'd probably parse a comma-separated list
}
