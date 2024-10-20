package application

import "os"

// Config struct that holds the configuration of the server
type Config struct {
	Host             string
	Port             string
	Environment      string
	DatabaseHost     string
	DatabasePort     string
	DatabaseName     string
	DatabaseUser     string
	DatabasePassword string
}

// LoadConfig loads the configuration from the Environment variables
func LoadConfig() *Config {
	return &Config{
		Host:             getEnvOrDefault("HOST", "0.0.0.0"),
		Port:             getEnvOrDefault("PORT", "8080"),
		Environment:      getEnvOrDefault("ENVIRONMENT", "development"),
		DatabaseHost:     os.Getenv("DATABASE_HOST"),
		DatabasePort:     os.Getenv("DATABASE_PORT"),
		DatabaseName:     os.Getenv("DATABASE_NAME"),
		DatabaseUser:     os.Getenv("DATABASE_USER"),
		DatabasePassword: os.Getenv("DATABASE_PASSWORD"),
	}
}

func getEnvOrDefault(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}