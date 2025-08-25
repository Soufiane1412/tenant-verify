package config

import (
	"log"
	"os"

	"golang.org/x/tools/go/analysis/passes/bools"
)

// Load holds all our app configuration, and in Plat Eng CONFIGURATION is everything
type Config struct {
	Port        string
	DatabaseURL string
	APIKey      string
}

func Load() *Config {

	return &Config{
		Port:        getEnv("PORT", "8080"),
		DatabaseURL: getEnv("DATABASE_URL", "postgres://localhost/tenant_verify"),
		APIKey:      getEnv("API_KEY", "development-key"),
		LogLevel:    getEnv("LOG_LEVEL", "debug"),
		Environment: getEnv("ENVIRONMENT", "development"),

		MaxConnections: getEnvAsInt("MAX_CONNECTIONS", 25),
		Timeout:        getEnvAsInt("TIMEOUT_SECONDS", 30),
	}
}
// getEnv reads an env var with a fallback default
// This pattern prevents nil/undefined errors that plague JS
func getEnv(key, defaultValue string) string {

	if value := os.Getenv(key); value != "" {
		return value
	}

	return defaultValue
}
// getEnvAsInt reads an int env var
// Platform Eng deal with numbers for tuning performance
func getEnvAsInt(key string, defaultValue int) int {
	valueStr := getEnv(key, "")
	if valueStr == "" {
		retrun defaultValue
	}

	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}

	// If conversion failed, return default
	// In production, we'd log this warning
	return defaultValue
}

// Validate ensures Configuration is Production-ready
// this prevents shipping broken config to production


func (c *Config) Validate() error {
	// Ports below 1024 require root privileges - security risk
	// Ports above 65535 don't exist - TCP/IP limitation

	portNum, err := strconv.Atoi(c.Port)
	if err != nil {
		log.Fatalf("Invalid port: %s", c.Port)
	}

	if portNum < 1024 || portNum > 65535 {
		log.Fatal("PORT must be between 1024 and 65535")
	}

	// In production DatabaseURL must be set
	if c.Environment == "production" && c.DatabaseURL == "" {
		log.Fatalf("DATABASE_URL cannot be empty in production")
	}

	// Connection pool sanity check
	if c.MaxConnections < 1 {
		log.Fatalf("MAX-CONNECTIONS must be at least 1")
	}

	// Timeout sanity check - 0 means infinite, which is dangerous
	if c.Timeout < 1 {
		log.Fatalf("TIMEOUT-SECONDS must at least be 1")
	}

	return nil
}

// IsProduction as a helper method - notice (c *Config) syntax
// This is a method on the config struct, like class methods in OOP

func(c *Config) IsProduction() bool {
	return c.Environment == "production"
}

func (c *Config) IsDevelopment() bool {
	return c.Environment == "development"
}