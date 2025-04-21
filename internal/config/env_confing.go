package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
}

type AppConfig struct {
	Port       string
	Rest       string
	GRPC       string
	Prometheus string
	JWTSecret  string
}

type Config struct {
	DB  DBConfig
	App AppConfig
}

func getEnvOrDefault(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}

func Load() *Config {
	_ = godotenv.Load()

	return &Config{
		DB: DBConfig{
			Host:     getEnvOrDefault("DB_HOST", "localhost"),
			Port:     getEnvOrDefault("DB_PORT", "5432"),
			User:     getEnvOrDefault("DB_USER", "postgres"),
			Password: getEnvOrDefault("DB_PASSWORD", "password"),
			Name:     getEnvOrDefault("DB_NAME", "pvz"),
		},
		App: AppConfig{
			Rest:       getEnvOrDefault("APP_PORT", "8080"),
			GRPC:       getEnvOrDefault("APP_GRPC_PORT", "3000"),
			Prometheus: getEnvOrDefault("APP_PROMETHEUS_PORT", "9000"),
			JWTSecret:  getEnvOrDefault("JWT_SECRET", "jwt-secret"),
		},
	}
}

func (c *DBConfig) DSN() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		c.User, c.Password, c.Host, c.Port, c.Name,
	)
}
