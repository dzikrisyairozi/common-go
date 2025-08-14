package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

// LoadEnv loads environment variables from .env file
func LoadEnv(envPath ...string) {
	var err error

	if len(envPath) > 0 {
		// Load specific .env file path
		err = godotenv.Load(envPath[0])
	} else {
		// Try to load .env from current directory
		err = godotenv.Load()
		if err != nil {
			// Try to load from parent directory (for workspace setup)
			err = godotenv.Load(filepath.Join("..", ".env"))
		}
	}

	if err != nil {
		log.Println("No .env file found, using system environment variables")
	}
}

// GetEnv gets environment variable with fallback
func GetEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

// GetRequiredEnv gets required environment variable, panics if not found
func GetRequiredEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("Required environment variable %s is not set", key)
	}
	return value
}

// DatabaseConfig represents database configuration
type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	Schema   string
	SSLMode  string
}

// LoadDatabaseConfig loads database configuration from environment variables
func LoadDatabaseConfig() *DatabaseConfig {
	return &DatabaseConfig{
		Host:     GetEnv("DB_HOST", "localhost"),
		Port:     GetEnv("DB_PORT", "5432"),
		User:     GetRequiredEnv("DB_USER"),
		Password: GetRequiredEnv("DB_PASSWORD"),
		Name:     GetRequiredEnv("DB_NAME"),
		Schema:   GetEnv("DB_SCHEMA", "public"),
		SSLMode:  GetEnv("DB_SSL_MODE", "disable"),
	}
}

// GetConnectionString returns PostgreSQL connection string
func (db *DatabaseConfig) GetConnectionString() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s search_path=%s",
		db.Host, db.Port, db.User, db.Password, db.Name, db.SSLMode, db.Schema)
}

// AppConfig represents general application configuration
type AppConfig struct {
	Port      string
	JWTSecret string
	Database  *DatabaseConfig
}

// LoadAppConfig loads application configuration
func LoadAppConfig() *AppConfig {
	LoadEnv()

	return &AppConfig{
		Port:      GetEnv("APP_PORT", "8080"),
		JWTSecret: GetRequiredEnv("JWT_SECRET"),
		Database:  LoadDatabaseConfig(),
	}
}
