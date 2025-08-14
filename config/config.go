package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

func LoadEnv() {
	_ = godotenv.Load()
	_ = godotenv.Load(filepath.Join("..", ".env"))

	appEnv := os.Getenv("APP_ENV")
	if appEnv == "" {
		appEnv = "dev"
	}
	fmt.Printf("Using APP_ENV = %s\n", appEnv)
	_ = godotenv.Overload(fmt.Sprintf(".env.%s", appEnv))
	_ = godotenv.Overload(filepath.Join("..", fmt.Sprintf(".env.%s", appEnv)))
}

func GetEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

func GetRequiredEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("Required environment variable %s is not set", key)
	}
	return value
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	Schema   string
	SSLMode  string
}

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

func (db *DatabaseConfig) GetConnectionString() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s search_path=%s",
		db.Host, db.Port, db.User, db.Password, db.Name, db.SSLMode, db.Schema)
}

type AppConfig struct {
	Port     string
	Database *DatabaseConfig
}

func LoadAppConfig() *AppConfig {
	LoadEnv()

	return &AppConfig{
		Port:     GetEnv("APP_PORT", "8080"),
		Database: LoadDatabaseConfig(),
	}
}
